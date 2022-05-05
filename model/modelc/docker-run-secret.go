package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker run param '--secret' declaration
type DockerRunSecret struct {
	// Secret source
	Source string `json:"source"`
	// Secret target
	Target string `json:"target"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerRunSecret) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Source = ctmpl.CompileStringOrDie(ctx, dr.Source, props)
	dr.Target = ctmpl.CompileStringOrDie(ctx, dr.Target, props)
}
