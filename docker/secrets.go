package docker

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Loads available docker secrets using `docker secret ls` command
func LoadAvailableSecrets(ctx common.DCtx) common.CmdRes {
	clog.Trace(ctx, "LoadAvailableSecrets()")
	res, secrets := loadAvailableStringArr(ctx, "docker secret ls --format \"{{.Name}}\"")
	Available.Secrets = secrets
	return res
}
