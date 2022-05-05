package docker

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Loads available docker networks using `docker network ls` command
func LoadAvailableNetworks(ctx common.DCtx) common.CmdRes {
	clog.Trace(ctx, "LoadAvailableNetworks()")
	res, networks := loadAvailableStringArr(ctx, "docker network ls --format \"{{.Name}}\"")
	Available.Networks = networks
	return res
}
