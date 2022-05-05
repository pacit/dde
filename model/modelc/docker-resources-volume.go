package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker volume declaration
type DockerResourceVolume struct {
	// Volume name
	Name string `json:"name"`
	// Volume driver
	Driver string `json:"driver"`
	// Volume options
	Options []string `json:"options"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerResourceVolume) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Name = ctmpl.CompileStringOrDie(ctx, dr.Name, props)
	dr.Driver = ctmpl.CompileStringOrDie(ctx, dr.Driver, props)
	ctmpl.CompileStringArrOrDie(ctx, dr.Options, props)
}
