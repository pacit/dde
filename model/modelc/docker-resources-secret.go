package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker secret declaration.
//
// Shoud use one of: Value | File
type DockerResourceSecret struct {
	// Secret name
	Name string `json:"name"`
	// Secret driver
	Driver string `json:"driver"`
	// Secret value (file content)
	Value string `json:"value"`
	// Secret file path
	File string `json:"file"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerResourceSecret) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Name = ctmpl.CompileStringOrDie(ctx, dr.Name, props)
	dr.Value = ctmpl.CompileStringOrDie(ctx, dr.Value, props)
	dr.File = ctmpl.CompileStringOrDie(ctx, dr.File, props)
	dr.Driver = ctmpl.CompileStringOrDie(ctx, dr.Driver, props)
}
