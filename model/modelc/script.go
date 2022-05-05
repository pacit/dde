package modelc

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Script run configuration
type ScriptJson struct {
	// Command to run
	Cmd string `json:"cmd"`
	// Script file path
	File string `json:"file"`
	// Where the file exists
	FileIsIn string `json:"fileIsIn"` // host|container
	// Where to run command/script file
	RunIn string `json:"runIn"`
	// File script arguments
	Args []string `json:"args"`
}

// Creates command to run script
func (s ScriptJson) GetCommand(workingDir string) string {
	cmd := "cd " + workingDir + " && "
	if len(s.Cmd) > 0 {
		cmd += s.Cmd
	} else if len(s.File) > 0 {
		cmd += s.File
	}
	if len(s.Args) > 0 {
		for _, a := range s.Args {
			cmd += " " + a
		}
	}
	return cmd
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (s *ScriptJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	ctmpl.CompileStringArrOrDie(ctx, s.Args, props)
	s.Cmd = ctmpl.CompileStringOrDie(ctx, s.Cmd, props)
	s.File = ctmpl.CompileStringOrDie(ctx, s.File, props)
	s.FileIsIn = ctmpl.CompileStringOrDie(ctx, s.FileIsIn, props)
	s.RunIn = ctmpl.CompileStringOrDie(ctx, s.RunIn, props)
}
