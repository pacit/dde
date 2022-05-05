package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker repository config
type WorkspaceJsonDockerRepo struct {
	// Config name, repo name, id - to use in project configs
	Name string `json:"name"`
	// Protocol: http|https
	Protocol string `json:"protocol"`
	// Repository address or IP
	Address string `json:"address"`
	// Path to file contains repository username
	UsernameFile string `json:"usernameFile"`
	// Path to file contains repository password
	PasswordFile string `json:"passwordFile"`
	// Username readed from file `UsernameFile`. It not exists in config json.
	Username string `json:"-"`
	// Password readed from file `PasswordFile`. It not exists in config json.
	Password string `json:"-"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *WorkspaceJsonDockerRepo) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Name = ctmpl.CompileStringOrDie(ctx, dr.Name, props)
	dr.Protocol = ctmpl.CompileStringOrDie(ctx, dr.Protocol, props)
	dr.Address = ctmpl.CompileStringOrDie(ctx, dr.Address, props)
	dr.UsernameFile = ctmpl.CompileStringOrDie(ctx, dr.UsernameFile, props)
	dr.PasswordFile = ctmpl.CompileStringOrDie(ctx, dr.PasswordFile, props)
	dr.Username = ctmpl.CompileStringOrDie(ctx, dr.Username, props)
	dr.Password = ctmpl.CompileStringOrDie(ctx, dr.Password, props)
}
