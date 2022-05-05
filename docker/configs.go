package docker

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Loads available docker configs
//
// Using command `docker config ls`
func LoadAvailableConfigs(ctx common.DCtx) common.CmdRes {
	clog.Trace(ctx, "LoadAvailableConfigs()")
	res, configs := loadAvailableStringArr(ctx, "docker config ls --format \"{{.Name}}\"")
	Available.Configs = configs
	return res
}
