package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Action to run before custom script run
type RunScriptBeforeRunJson struct {
	// Project name to update sources from git repository
	UpdateProject string `json:"updateProject"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (rs *RunScriptBeforeRunJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	rs.UpdateProject = ctmpl.CompileStringOrDie(ctx, rs.UpdateProject, props)

}
