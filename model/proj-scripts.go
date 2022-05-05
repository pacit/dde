package model

import "github.com/pacit/dde/model/modelc"

// Project custom scripts config
type ProjectJsonScripts struct {
	// Script to build project image
	Build modelc.ScriptJson `json:"build"`
}
