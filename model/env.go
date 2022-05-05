package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
	"github.com/pacit/dde/model/modelc"
)

// Environment config file name
const EnvJsonFileName = "env.json"

// Environment config
type EnvJson struct {
	// Environment custom scripts config
	Scripts EnvJsonScripts `json:"scripts"`
	// Environment docker resources config
	DockerResources modelc.DockerResourcesJson `json:"dockerResources"`
	// Parent docker run configuration for all environment services
	DockerRun modelc.DockerRunJson `json:"dockerRun"`
	// Directories contains template files `tmpl-*`
	TemplateFileDirs []string `json:"templateFileDirs"`
	// Environment properties to use in templates
	Properties map[string]string `json:"properties"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (ej *EnvJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	for i, p := range ej.Scripts.Prepare {
		(&p).CompileTemplatesOrDie(ctx, props)
		ej.Scripts.Prepare[i] = p
	}
	for i, c := range ej.Scripts.Cleanup {
		(&c).CompileTemplatesOrDie(ctx, props)
		ej.Scripts.Cleanup[i] = c
	}
	(&ej.DockerResources).CompileTemplatesOrDie(ctx, props)
	(&ej.DockerRun).CompileTemplatesOrDie(ctx, props)
	ctmpl.CompileStringArrOrDie(ctx, ej.TemplateFileDirs, props)
}
