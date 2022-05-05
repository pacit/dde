package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

const WorkspaceJsonFileName = "workspace.json"

// Workspace configuration
type WorkspaceJson struct {
	// Default env names
	//
	// Used when -e/-es param is not set
	DefaultEnvs []string `json:"defaultEnvs"`
	// Default version file
	//
	// Used when -v/-vf param is not used
	DefaultVerFile string `json:"defaultVerFile"`
	// Custom scripts configs
	//
	// Script can be run using 'dde runs' command
	RunScripts []RunScriptJson `json:"runScripts"`
	// Docker repos
	DockerRepos []WorkspaceJsonDockerRepo `json:"dockerRepos"`
	// Path where 'tmpl-' files exists
	TemplateFileDirs []string          `json:"templateFileDirs"`
	Properties       map[string]string `json:"properties"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (wj *WorkspaceJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	ctmpl.CompileStringArrOrDie(ctx, wj.DefaultEnvs, props)
	wj.DefaultVerFile = ctmpl.CompileStringOrDie(ctx, wj.DefaultVerFile, props)
	for i, rs := range wj.RunScripts {
		(&rs).CompileTemplatesOrDie(ctx, props)
		wj.RunScripts[i] = rs
	}
	for i, dr := range wj.DockerRepos {
		(&dr).CompileTemplatesOrDie(ctx, props)
		wj.DockerRepos[i] = dr
	}
	ctmpl.CompileStringArrOrDie(ctx, wj.TemplateFileDirs, props)
}
