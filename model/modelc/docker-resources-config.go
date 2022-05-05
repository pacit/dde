package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker config declaration.
//
// Shoud use one of: Value | File
type DockerResourceConfig struct {
	// Config name
	Name string `json:"name"`
	// Config value, if it's a simple string value (config file content)
	Value string `json:"value"`
	// Config file path, if it's a config file
	File string `json:"file"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerResourceConfig) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Name = ctmpl.CompileStringOrDie(ctx, dr.Name, props)
	dr.Value = ctmpl.CompileStringOrDie(ctx, dr.Value, props)
	dr.File = ctmpl.CompileStringOrDie(ctx, dr.File, props)
}
