package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cio"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/jobs"
	"github.com/pacit/dde/model"
)

// Command line arguments
type CliArgs struct {
	// Command - first argument after `dde`
	Command string
	// Other arguments which are not parsed as dde command arguments. They are used as custom script arguments
	OtherArgs []string
	// Log levels argument value
	LogLevels string
	// Number of job concurrent runners (threads)
	Threads int
	// Environment names to select
	Environments []string
	// Service names to select
	Services []string
	// Selected services in environments. Format: '{envName}:{srvName}'
	EnvSvr []string
	// File name with versions (git tags/branches)
	VersionFile string
	// Single version to use instead of version file
	Version string
	// Forced build project's images
	ForceBuild bool
	// Forced redeploy services
	ForceRedeploy bool
	// Forced prepare environments. Runs environment's prepare sctipts
	ForcePrepareEnv bool
	// Removes environments resources and runs clean scripts
	ForceRmEnv bool
	// Removes docker volumes declared by services and environments
	ForceRmVolumes bool
	// Just calculate jobs to do. Print jobs table and exit without jobs run
	DryRun bool
}

// Workspace engine
type Workspace struct {
	// Workspace configuration (from json file)
	Cfg model.WorkspaceJson
	// Workspace's environments (only selected by cli arguments)
	Environments []*Environment
	// Workspace's projects (only selected by cli arguments)
	Projects []*Project
	// Selected environment's names
	SelectedEnvironments []string
	// Selected services. Format: {envName}:{srvName}
	SelectedServices []string
	// Selected service's names
	SelectedServicesOnlyNames []string
	// Map with project's versions
	VersionMap map[string]string
	// Command line arguments
	CliArgs CliArgs
}

// Creates new Workspace engine
func New(ctx common.DCtx) (*Workspace, error) {
	wrk := &Workspace{
		Cfg: model.WorkspaceJson{
			Properties: make(map[string]string),
		},
		VersionMap: make(map[string]string),
	}
	if err := loadWorkspaceCfg(ctx, &wrk.Cfg); err != nil {
		return nil, err
	}
	wrk.Cfg.CompileTemplatesOrDie(ctx, wrk.GetTemplateProps(ctx, true))
	wrk.loadCfgPasswords(ctx)
	return wrk, nil
}

// Finds project engine by project name
func (w *Workspace) GetProjectByName(ctx common.DCtx, name string) *Project {
	for _, proj := range w.Projects {
		if proj.Name == name {
			clog.Trace(ctx, "GetProjectByName", name, "="+cpath.ProjDefinitionDir(proj.Name))
			return proj
		}
	}
	clog.Trace(ctx, "GetProjectByName", name, "=nil")
	return nil
}

// Loads Workspace configuration (from json file)
func loadWorkspaceCfg(ctx common.DCtx, cfg *model.WorkspaceJson) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	wrkCfgFilePath := filepath.Join(wd, model.WorkspaceJsonFileName)
	return cio.ReadJsonFile(ctx, wrkCfgFilePath, cfg)
}

// Loads passwords defined in json configuration as files
func (w *Workspace) loadCfgPasswords(ctx common.DCtx) {
	for _, repo := range w.Cfg.DockerRepos {
		if len(repo.UsernameFile) > 0 {
			repo.Username = cio.ReadTextFileSilent(ctx, cpath.GetAbsPath(repo.UsernameFile))
		}
		if len(repo.PasswordFile) > 0 {
			repo.Password = cio.ReadTextFileSilent(ctx, cpath.GetAbsPath(repo.PasswordFile))
		}
	}
}

// Gets all environment's names
func (w *Workspace) GetAllEnvNames() ([]string, error) {
	return cio.GetDirNamesContainsFile(cpath.WrkDefinitionEnvsDir(), "env.json")
}

// Gets all service's names in selected environments
func (w *Workspace) GetAllServiceNamesInSelectedEnvs(ctx common.DCtx) []string {
	names := []string{}
	envsDirPath := cpath.WrkDefinitionEnvsDir()
	for _, envName := range w.SelectedEnvironments {
		envDirPath := filepath.Join(envsDirPath, envName)
		allNames, err := cio.GetDirNamesContainsFile(envDirPath, "srv.json")
		if err != nil {
			clog.Panic(ctx, err, 120, "Error getting service names")
		}
		// filter service templates
		for _, name := range allNames {
			fullName := envName + ":" + name
			if !common.StringSliceContains(names, fullName) { // distinct
				names = append(names, fullName)
			}
		}
	}
	return names
}

// Loads selected environments
func (w *Workspace) LoadSelectedEnvironments(ctx common.DCtx) error {
	clog.Trace(ctx, "LoadSelectedEnvironments()")
	envsDirPath := cpath.WrkDefinitionEnvsDir()
	for _, envName := range w.SelectedEnvironments {
		envDir := filepath.Join(envsDirPath, envName)
		env, err := NewEnvironment(ctx, envDir, w)
		if err != nil {
			return err
		}
		w.Environments = append(w.Environments, env)
	}
	return nil
}

// Loads selected services
func (w *Workspace) LoadSelectedServices(ctx common.DCtx) error {
	clog.Trace(ctx, "wrk.LoadSelectedServices()")
	for _, env := range w.Environments {
		err := env.LoadSelectedServices(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Loads required projects
func (w *Workspace) LoadRequiredProjects(ctx common.DCtx) error {
	clog.Trace(ctx, "wrk.LoadRequiredProjects()")
	projectsDir := cpath.WrkDefinitionProjectsDir()
	for _, env := range w.Environments {
		clog.Trace(ctx, "Load projects for env: "+env.Name)
		for _, srv := range env.Services {
			clog.Trace(ctx, "Load project for service: "+srv.Name)
			if len(srv.Cfg.Project) > 0 {
				proj := w.FindProject(srv.Cfg.Project)
				if proj == nil {
					proj, err := NewProject(ctx, filepath.Join(projectsDir, srv.Cfg.Project), w)
					if err != nil {
						return err
					}
					w.Projects = append(w.Projects, proj)
				}
			}
		}
	}
	if w.CliArgs.Command == "runs" {
		clog.Debug(ctx, "Load projects for custom scripts")
		for _, rs := range w.Cfg.RunScripts {
			for _, br := range rs.BeforeRun {
				if len(br.UpdateProject) > 0 {
					// load project
					proj := w.FindProject(br.UpdateProject)
					if proj == nil {
						proj, err := NewProject(ctx, filepath.Join(projectsDir, br.UpdateProject), w)
						if err != nil {
							return err
						}
						w.Projects = append(w.Projects, proj)
					}
				}
			}
		}
	}
	return nil
}

// Creates versions map
func (w *Workspace) CreateVersionsMap(ctx common.DCtx) error {
	clog.Trace(ctx, "wrk.CreateVersionsMap()")
	if len(w.CliArgs.Version) > 0 {
		clog.Info(ctx, "Use single version for all - "+w.CliArgs.Version)
		for _, proj := range w.Projects {
			w.VersionMap[proj.Name] = w.CliArgs.Version
		}
	} else if len(w.CliArgs.VersionFile) > 0 {
		verFilePath := getVerFilePath(ctx, cpath.WrkVersionsDir(), w.CliArgs.VersionFile)
		clog.Info(ctx, "Use versions from file", w.CliArgs.VersionFile, verFilePath)
		verFileMap := loadVersionsFromFile(ctx, verFilePath)
		defVer := verFileMap["_default"]
		for _, proj := range w.Projects {
			projVer := verFileMap[proj.Name]

			if len(projVer) > 0 {
				w.VersionMap[proj.Name] = projVer
			} else if len(defVer) > 0 {
				w.VersionMap[proj.Name] = defVer
			} else {
				panic(errors.New("No version for project: " + proj.Name))
			}
		}
	} else {
		clog.Info(ctx, "No version provided - using 'master' for all")
		for _, proj := range w.Projects {
			w.VersionMap[proj.Name] = "master"
		}
	}
	clog.Info(ctx, "Use versions:")
	for k, v := range w.VersionMap {
		clog.Info(ctx, fmt.Sprintf("%30s -> %s", k, v))
	}
	return nil
}

// Compiles all go templates in configurations
func (w *Workspace) CompileTemplatesInCfgJsons(ctx common.DCtx) error {
	w.Cfg.CompileTemplatesOrDie(ctx, w.GetTemplateProps(ctx, true))
	for _, p := range w.Projects {
		p.Cfg.CompileTemplatesOrDie(ctx, p.GetTemplateProps(ctx, true))
	}
	for _, e := range w.Environments {
		e.Cfg.CompileTemplatesOrDie(ctx, e.GetTemplateProps(ctx, true))
		for _, s := range e.Services {
			s.Cfg.CompileTemplatesOrDie(ctx, s.GetTemplateProps(ctx, true))
		}
	}
	return nil
}

// Finds environment engine by environment name
func (w *Workspace) FindEnvironment(name string) *Environment {
	for _, env := range w.Environments {
		if env.Name == name {
			return env
		}
	}
	return nil
}

// Finds project engine by project name
func (w *Workspace) FindProject(name string) *Project {
	for _, proj := range w.Projects {
		if proj.Name == name {
			return proj
		}
	}
	return nil
}

// Converts Workspace engine to info object used by jobs
func (w *Workspace) ToJobEnvironmentInfo() jobs.JobWorkspaceInfo {
	return jobs.JobWorkspaceInfo{
		Dir: cpath.WrkDefinitionDir(),
		Cfg: w.Cfg,
	}
}

// Finds custom script to run by script name
func (w *Workspace) GetRunScriptByName(n string) *model.RunScriptJson {
	for _, rs := range w.Cfg.RunScripts {
		if rs.Name == n {
			return &rs
		}
	}
	return nil
}

// Finds docker repo config by name
func (w *Workspace) GetDockerRepo(name string) *model.WorkspaceJsonDockerRepo {
	for _, r := range w.Cfg.DockerRepos {
		if r.Name == name {
			return &r
		}
	}
	return nil
}

// Creates Workspace properties to use in templates
func (w *Workspace) GetTemplateProps(ctx common.DCtx, compile bool) map[string]string {
	props := make(map[string]string)
	// builtin props
	props["ddeWrkDir"] = cpath.WrkDefinitionDir()
	now := time.Now()
	props["ddeRunDateYear4"] = now.Format("2006")
	props["ddeRunDateYear2"] = now.Format("06")
	props["ddeRunDateMonth2"] = now.Format("01")
	props["ddeRunDateMonthS"] = now.Format("Jan")
	props["ddeRunDateMonthF"] = now.Format("January")
	props["ddeRunDateDay2"] = now.Format("02")
	props["ddeRunDateDay1"] = now.Format("2")
	props["ddeRunTimeHours12"] = now.Format("03")
	props["ddeRunTimeHours24"] = now.Format("15")
	props["ddeRunTimeMinutes2"] = now.Format("04")
	props["ddeRunTimeMinutes1"] = now.Format("4")
	props["ddeRunTimeSeconds2"] = now.Format("05")
	props["ddeRunTimeSeconds1"] = now.Format("5")
	if compile {
		return prepareAndCompileTemplateProperties(ctx, props, w.Cfg.Properties)
	} else {
		return prepareTemplateProperties(ctx, props, w.Cfg.Properties)
	}
}

// Gets version file name from command line attribute '-vf' value
func getVerFilePath(ctx common.DCtx, verFilesPath string, verName string) string {
	verFileNames := cio.GetFileNamesInDirSilent(ctx, verFilesPath)
	// check full name exists
	name := verName
	if common.StringSliceContains(verFileNames, name) {
		return filepath.Clean(filepath.Join(verFilesPath, name))
	}
	// check properties files
	name = verName + ".properties"
	if common.StringSliceContains(verFileNames, name) {
		return filepath.Clean(filepath.Join(verFilesPath, name))
	}
	// check json files
	name = verName + ".json"
	if common.StringSliceContains(verFileNames, name) {
		return filepath.Clean(filepath.Join(verFilesPath, name))
	}
	// kaboom
	panic("There is no version file with name '" + name + "'")
}

// Reads versions file (.properties or .json) into a map
//
// '_parent' property is recognized as a parent file
func loadVersionsFromFile(ctx common.DCtx, path string) map[string]string {
	clog.Info(ctx, "Read versions from file", path)
	vMap := map[string]string{}
	if strings.HasSuffix(path, ".json") {
		vMap = cio.ReadJsonFileAsMapSilent(ctx, path)
	} else if strings.HasSuffix(path, ".properties") {
		vMap = cio.ReadPropertiesFileAsMapSilent(ctx, path)
	}
	if val, ok := vMap["_parent"]; ok {
		parentPath := filepath.Dir(path)
		parentPath = filepath.Clean(filepath.Join(parentPath, val))
		parentVMap := loadVersionsFromFile(ctx, parentPath)
		for k, v := range vMap {
			parentVMap[k] = v
		}
		return parentVMap
	}
	return vMap
}
