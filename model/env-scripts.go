package model

import "github.com/pacit/dde/model/modelc"

// Environment custom scripts config
type EnvJsonScripts struct {
	// Scripts to run in `dde prepare` command. Runs in order
	Prepare []modelc.ScriptJson `json:"prepare"`
	// Scripts to run in `dde rm -re` command. Runs in order
	Cleanup []modelc.ScriptJson `json:"cleanup"`
}
