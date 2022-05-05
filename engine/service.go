package engine

import (
	"fmt"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cio"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/docker"
	"github.com/pacit/dde/jobs"
	"github.com/pacit/dde/model"
)

// Service engine
type Service struct {
	// Service Configuration (from json file)
	Cfg model.ServiceJson
	// Pointer to an Environment (parent)
	Environment *Environment
	// Pointer to a workspace (parent parent)
	Workspace *Workspace
	// Service name
	Name string
	// FIXME: wtf????
	CachedImageToPull *model.ProjectJsonDockerImage
}

// Creates new Service engine
func NewService(ctx common.DCtx, dirPath string, env *Environment) (*Service, error) {
	clog.Trace(ctx, fmt.Sprintf("NewService(%s)", dirPath))
	srv := &Service{
		Cfg:         model.ServiceJson{},
		Name:        filepath.Base(dirPath),
		Environment: env,
		Workspace:   env.Workspace,
	}
	if err := cio.ReadJsonFile(ctx, filepath.Join(dirPath, model.ServiceJsonFileName), &srv.Cfg); err != nil {
		clog.Panic(ctx, err, 120, "Cannot load service config file", dirPath)
		return nil, err
	}
	srv.setDefaults(ctx)
	return srv, nil
}

// Checks if project image must be built
func (s *Service) NeedBuildImage(ctx common.DCtx) bool {
	if s.Workspace.CliArgs.ForceBuild || s.Workspace.CliArgs.Command == "build" {
		// if -fb flag present
		return true
	}
	needVersion := s.Workspace.VersionMap[s.Cfg.Project]
	needVersionDockerSafe := VersionDockerSafe(needVersion)
	if docker.Available.HasImage(s.Cfg.Project, needVersionDockerSafe) {
		// if local image exists - do not build
		return false
	}
	proj := s.Workspace.GetProjectByName(ctx, s.Cfg.Project)
	if proj != nil {
		for _, imgName := range proj.GetImageNames() {
			if docker.Available.HasImage(imgName, needVersionDockerSafe) {
				// if exists any image from remote repo
				return false
			}
		}
	}
	return true
}

// Check remote repositories and search project images
func (s *Service) GetImageToPull(ctx common.DCtx) *model.ProjectJsonDockerImage {
	needVersion := s.Workspace.VersionMap[s.Cfg.Project]
	proj := s.Workspace.GetProjectByName(ctx, s.Cfg.Project)
	if proj != nil {
		for _, di := range proj.Cfg.DockerImages {
			repoCfg := s.Workspace.GetDockerRepo(di.RepoId)
			available := docker.CheckImageExistsInRepo(ctx, *repoCfg, di.ImageName, needVersion)
			clog.Debug(ctx, "CheckImageExistsInRepo", di.ImageName, needVersion, fmt.Sprintf("%v", available))
			if available {
				s.CachedImageToPull = &di
				return &di
			}
		}
	}
	return nil
}

// Creates docker service full name
func (s *Service) GetDockerServiceName() string {
	return s.Environment.Name + "_" + s.Name
}

// Checks if service needs to be removed
func (s *Service) NeedRmService(ctx common.DCtx) bool {
	dockerServiceName := s.GetDockerServiceName()
	needVersion := s.Workspace.VersionMap[s.Cfg.Project]
	needVersionDockerSafe := VersionDockerSafe(needVersion)
	if docker.Available.HasServiceName(dockerServiceName) {
		// if service is running
		if s.Workspace.CliArgs.Command == "rm" {
			// rm command
			clog.Trace(ctx, fmt.Sprintf("NeedRmService(%s)=%v, (%s)", s.Name, true, "rm command"))
			return true
		} else if s.Workspace.CliArgs.ForceRedeploy {
			// if -fr flag is present
			clog.Trace(ctx, fmt.Sprintf("NeedRmService(%s)=%v, (%s)", s.Name, true, "-fr flag"))
			return true
		} else if !docker.Available.HasServiceInVersion(dockerServiceName, needVersionDockerSafe) {
			// if service runs in incorrect version
			clog.Trace(ctx, fmt.Sprintf("NeedRmService(%s)=%v, (%s)", s.Name, true, "runs in incorrect version"))
			return true
		}
	} else {
		clog.Trace(ctx, fmt.Sprintf("NeedRmService(%s)=%v, (%s)", s.Name, false, "not running"))
		return false
	}
	clog.Trace(ctx, fmt.Sprintf("NeedRmService(%s)=%v, (%s)", s.Name, false, "unknown"))
	return false
}

// Check if the service needs to be deployed
func (s *Service) NeedDeploy() bool {
	dockerServiceName := s.GetDockerServiceName()
	needVersion := s.Workspace.VersionMap[s.Cfg.Project]
	needVersionDockerSafe := VersionDockerSafe(needVersion)
	if s.Workspace.CliArgs.ForceRedeploy {
		// if -fr flag is present
		return true
	} else if docker.Available.HasServiceInVersion(dockerServiceName, needVersionDockerSafe) {
		// if service runs in correct version
		return false
	}
	return true
}

// Converts Service engine to info object used by jobs
func (s *Service) ToJobServiceInfo() jobs.JobServiceInfo {
	return jobs.JobServiceInfo{
		EnvName: s.Environment.Name,
		SrvName: s.Name,
		Cfg:     s.Cfg,
	}
}

// Creates Service properties to use in templates
func (s *Service) GetTemplateProps(ctx common.DCtx, compile bool) map[string]string {
	props := s.Environment.GetTemplateProps(ctx, false)
	// builtin props
	props["ddeSrvName"] = s.Name
	props["ddeProjName"] = s.Cfg.Project
	props["ddeSrvVersion"] = s.Workspace.VersionMap[s.Cfg.Project]
	props["ddeSrvVersionDockerSafe"] = VersionDockerSafe(props["ddeSrvVersion"])
	props["ddeSrvVersionDotnetSafe"] = VersionDotnetSafe(props["ddeSrvVersion"])
	if compile {
		return prepareAndCompileTemplateProperties(ctx, props, s.Cfg.Properties)
	} else {
		return prepareTemplateProperties(ctx, props, s.Cfg.Properties)
	}
}

// Sets service config defaults
func (s *Service) setDefaults(ctx common.DCtx) {
	pr := s.Environment.Cfg.DockerRun
	cr := &s.Cfg.DockerRun
	if len(s.Cfg.DockerRun.Name) == 0 {
		cr.Name = s.Environment.Name + "_" + s.Name
	}
	if len(s.Cfg.DockerRun.Hostname) == 0 {
		cr.Hostname = s.Environment.Name + "_" + s.Name
	}
	if s.Cfg.DockerRun.Replicas == 0 {
		cr.Replicas = 1
	}
	cr.Configs = append(cr.Configs, pr.Configs...)
	cr.Secrets = append(cr.Secrets, pr.Secrets...)
	cr.Networks = append(cr.Networks, pr.Networks...)
	cr.Envs = append(cr.Envs, pr.Envs...)
	cr.Mounts = append(cr.Mounts, pr.Mounts...)

}
