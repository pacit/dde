package engine

import (
	"fmt"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cio"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/jobs"
	"github.com/pacit/dde/model"
)

// Project engine
type Project struct {
	// Project configuration (from json file)
	Cfg model.ProjectJson
	// Pointer to a workspace (parent)
	Workspace *Workspace
	// Project name (dir name)
	Name string
	// Flag indicates that project has already defined job to prepare image
	HasImageJob bool
}

// Creates new Project engine
func NewProject(ctx common.DCtx, dirPath string, w *Workspace) (*Project, error) {
	clog.Trace(ctx, fmt.Sprintf("NewProject(%s)", dirPath))
	proj := Project{
		Cfg:       model.ProjectJson{},
		Workspace: w,
		Name:      filepath.Base(dirPath),
	}
	if err := cio.ReadJsonFile(ctx, filepath.Join(dirPath, model.ProjectJsonFileName), &proj.Cfg); err != nil {
		clog.Panic(ctx, err, 120, "Cannot load project config file", dirPath)
		return nil, err
	}
	proj.setDefaults(ctx)
	return &proj, nil
}

// Gets image names which can be used as this project images
func (p *Project) GetImageNames() []string {
	names := []string{p.Name}
	for _, di := range p.Cfg.DockerImages {
		names = append(names, di.ImageName)
	}
	return names
}

// Converts Project engine to info object used by jobs
func (p *Project) ToJobProjectInfo(v string) jobs.JobProjectInfo {
	return jobs.JobProjectInfo{
		Name:              p.Name,
		Cfg:               p.Cfg,
		Version:           v,
		VersionDockerSafe: VersionDockerSafe(v),
		VersionDotnetSafe: VersionDotnetSafe(v),
	}
}

// Creates Project properties to use in templates
func (p *Project) GetTemplateProps(ctx common.DCtx, compile bool) map[string]string {
	props := p.Workspace.GetTemplateProps(ctx, false)
	// builtin props
	props["ddeProjName"] = p.Name
	ver := p.Workspace.VersionMap[p.Name]
	props["ddeProjVersion"] = ver
	props["ddeProjVersionDockerSafe"] = VersionDockerSafe(ver)
	props["ddeProjVersionDotnetSafe"] = VersionDotnetSafe(ver)
	if compile {
		return prepareAndCompileTemplateProperties(ctx, props, p.Cfg.Properties)
	} else {
		return prepareTemplateProperties(ctx, props, p.Cfg.Properties)
	}
}

// Sets project config defaults
func (p *Project) setDefaults(ctx common.DCtx) {
	if len(p.Cfg.Scripts.Build.Cmd) == 0 && len(p.Cfg.Scripts.Build.File) == 0 {
		p.Cfg.Scripts.Build.Cmd = "docker build --add-host host.docker.internal:host-gateway -q -t {{.ddeProjName}}:{{.ddeProjVersionDockerSafe}} ."
	}
}
