package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker network declaration.
type DockerResourceNetwork struct {
	// Network name
	Name string `json:"name"`
	// Network driver: bridge|overlay
	Driver string `json:"driver"`
	// Network subnet mask. eg. 172.28.0.0/16
	Subnet string `json:"subnet"`
	// Network gateway address eg. 172.28.5.254
	Gateway string `json:"gateway"`
	// Network IP addresses range eg. 172.28.5.0/24
	IpRange string `json:"ipRange"`
	// Enable manual container attachment
	Attachable bool `json:"attachable"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerResourceNetwork) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Name = ctmpl.CompileStringOrDie(ctx, dr.Name, props)
	dr.Driver = ctmpl.CompileStringOrDie(ctx, dr.Driver, props)
	dr.Subnet = ctmpl.CompileStringOrDie(ctx, dr.Subnet, props)
	dr.Gateway = ctmpl.CompileStringOrDie(ctx, dr.Gateway, props)
	dr.IpRange = ctmpl.CompileStringOrDie(ctx, dr.IpRange, props)

}
