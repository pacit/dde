package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Project config file name
const ProjectJsonFileName = "proj.json"

// Project config
type ProjectJson struct {
	// Project git repository address - to clone sources
	GitRepo string `json:"gitRepo"`
	// Project custom scripts
	Scripts ProjectJsonScripts `json:"scripts"`
	// Docker images Which can be used as this project image.
	// They can have different names, can exists in different repositories, but must have the same version
	DockerImages []ProjectJsonDockerImage `json:"dockerImages"`
	// Project properties to use in templates
	Properties map[string]string `json:"properties"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (pj *ProjectJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	pj.GitRepo = ctmpl.CompileStringOrDie(ctx, pj.GitRepo, props)
	(&pj.Scripts.Build).CompileTemplatesOrDie(ctx, props)
	for i, di := range pj.DockerImages {
		(&di).CompileTemplatesOrDie(ctx, props)
		pj.DockerImages[i] = di
	}
}
