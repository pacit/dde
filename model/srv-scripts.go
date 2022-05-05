package model

import "github.com/pacit/dde/model/modelc"

// Service custom scripts config
type ServiceJsonScripts struct {
	// Script to check service is running
	IsRunning modelc.ScriptJson `json:"isRunning"`
	// Scripts to run before service in order
	BeforeRun []modelc.ScriptJson `json:"beforeRun"`
	// Script to run service
	Run modelc.ScriptJson `json:"run"`
	// Scripts to run after service in order
	AfterRun []modelc.ScriptJson `json:"afterRun"`
}
