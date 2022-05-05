package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker run param '--config' declaration
type DockerRunConfig struct {
	// Config source
	Source string `json:"source"`
	// Config target
	Target string `json:"target"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerRunConfig) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Source = ctmpl.CompileStringOrDie(ctx, dr.Source, props)
	dr.Target = ctmpl.CompileStringOrDie(ctx, dr.Target, props)
}
