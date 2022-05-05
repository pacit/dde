package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
	"github.com/pacit/dde/model/modelc"
)

// Custom script configuration
type RunScriptJson struct {
	// Custom script name - use it in `dde runs [name]` command
	Name string `json:"name"`
	// Before actions
	BeforeRun []RunScriptBeforeRunJson `json:"beforeRun"`
	// Scripts to run in order
	Scripts []modelc.ScriptJson `json:"scripts"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (rs *RunScriptJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	rs.Name = ctmpl.CompileStringOrDie(ctx, rs.Name, props)
	for i, b := range rs.BeforeRun {
		(&b).CompileTemplatesOrDie(ctx, props)
		rs.BeforeRun[i] = b
	}
	for i, s := range rs.Scripts {
		(&s).CompileTemplatesOrDie(ctx, props)
		rs.Scripts[i] = s
	}
}
