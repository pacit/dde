package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker run param '--publish' declaration
type DockerRunPublish struct {
	// Host port number
	Published string `json:"published"`
	// Container port number
	Target string `json:"target"`
	// Protocol
	Protocol string `json:"protocol"`
	// Publish mode
	Mode string `json:"mode"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerRunPublish) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Published = ctmpl.CompileStringOrDie(ctx, dr.Published, props)
	dr.Target = ctmpl.CompileStringOrDie(ctx, dr.Target, props)
	dr.Protocol = ctmpl.CompileStringOrDie(ctx, dr.Protocol, props)
	dr.Mode = ctmpl.CompileStringOrDie(ctx, dr.Mode, props)
}
