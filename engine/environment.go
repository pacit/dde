package engine

import (
	"fmt"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cio"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/jobs"
	"github.com/pacit/dde/model"
)

// Environment engine
type Environment struct {
	// Pointer to a Workspace (parent)
	Workspace *Workspace
	// Environment configuration (from json file)
	Cfg model.EnvJson
	// Environment name (dir name)
	Name string
	// Environment's services
	Services []*Service
}

// Creates new Environment engine
func NewEnvironment(ctx common.DCtx, dirPath string, wrk *Workspace) (*Environment, error) {
	clog.Trace(ctx, fmt.Sprintf("NewEnvironment(%s)", dirPath))
	env := &Environment{
		Workspace: wrk,
		Cfg: model.EnvJson{
			Properties: make(map[string]string),
		},
		Name: filepath.Base(dirPath),
	}
	if err := cio.ReadJsonFile(ctx, filepath.Join(dirPath, model.EnvJsonFileName), &env.Cfg); err != nil {
		clog.Panic(ctx, err, 120, "Cannot load environment config file", dirPath)
		return nil, err
	}
	return env, nil
}

// Converts Environment engine to info object used by jobs
func (e *Environment) ToJobEnvironmentInfo() jobs.JobEnvironmentInfo {
	return jobs.JobEnvironmentInfo{
		Name: e.Name,
		Cfg:  e.Cfg,
	}
}

// Loads Environment services which are selected to use
func (e *Environment) LoadSelectedServices(ctx common.DCtx) error {
	envDir := cpath.EnvDefinitionDir(e.Name)
	clog.Trace(ctx, fmt.Sprintf("[env:%s] LoadSelectedServices()", e.Name))
	serviceNames, err := cio.GetDirNamesContainsFile(envDir, "srv.json")
	if err != nil {
		return err
	}
	for _, srvName := range serviceNames {
		fullName := e.Name + ":" + srvName
		if common.StringSliceContains(e.Workspace.SelectedServices, fullName) {
			srv, err := NewService(ctx, filepath.Join(envDir, srvName), e)
			if err != nil {
				return err
			}
			e.Services = append(e.Services, srv)
		}
	}
	return nil
}

// Creates Environment properties to use in templates
func (e *Environment) GetTemplateProps(ctx common.DCtx, compile bool) map[string]string {
	props := e.Workspace.GetTemplateProps(ctx, false)
	// builtin props
	props["ddeEnvName"] = e.Name
	if compile {
		return prepareAndCompileTemplateProperties(ctx, props, e.Cfg.Properties)
	} else {
		return prepareTemplateProperties(ctx, props, e.Cfg.Properties)
	}
}
