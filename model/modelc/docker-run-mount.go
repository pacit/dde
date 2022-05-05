package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker run param '--mount' declaration
type DockerRunMount struct {
	// Mount type
	Type string `json:"type"`
	// Mount source
	Source string `json:"source"`
	// Mount destination
	Destination string `json:"destination"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerRunMount) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Source = ctmpl.CompileStringOrDie(ctx, dr.Source, props)
	dr.Destination = ctmpl.CompileStringOrDie(ctx, dr.Destination, props)
	dr.Type = ctmpl.CompileStringOrDie(ctx, dr.Type, props)
}
