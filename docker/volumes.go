package docker

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Loads available docker volumes using `docker volume ls` command
func LoadAvailableVolumes(ctx common.DCtx) common.CmdRes {
	clog.Trace(ctx, "LoadAvailableVolumes()")
	res, volumes := loadAvailableStringArr(ctx, "docker volume ls --format \"{{.Name}}\"")
	Available.Volumes = volumes
	return res
}
