package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/engine"
)

// Parses command line arguments
func parseArguments(ctx common.DCtx) engine.CliArgs {
	args := os.Args[1:]
	command := ""
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	} else {
		clog.Panic(ctx, errors.New("no arguments"), 111, "No arguments")
	}
	ddeFS := flag.NewFlagSet("dde", flag.ExitOnError)
	argLL := ddeFS.String("ll", "", "Log level")
	argT := ddeFS.Int("t", 4, "Threads number")
	argE := ddeFS.String("e", "", "An Environment name, or a list of names - separated by [space]")
	argS := ddeFS.String("s", "", "A Service name, or a list of names - separated by [space]")
	argES := ddeFS.String("es", "", "list of {env}:{srv} - separated by [space]")
	argVF := ddeFS.String("vf", "", "Version file name")
	argV := ddeFS.String("v", "", "Single version value to use")
	argFB := ddeFS.Bool("fb", false, "Force build projects")
	argFR := ddeFS.Bool("fr", false, "Force redeploy services")
	argFP := ddeFS.Bool("fp", false, "Force prepare environments")
	argRE := ddeFS.Bool("re", false, "Rm environment if empty")
	argRV := ddeFS.Bool("rv", false, "Rm volumes")
	argDR := ddeFS.Bool("dr", false, "Dry run")
	ddeFS.Parse(args)

	out := engine.CliArgs{
		Command:         command,
		LogLevels:       argStr(argLL),
		Threads:         argInt(argT),
		Environments:    argArr(argE),
		Services:        argArr(argS),
		EnvSvr:          argArr(argES),
		VersionFile:     argStr(argVF),
		Version:         argStr(argV),
		ForceBuild:      argBool(argFB),
		ForceRedeploy:   argBool(argFR),
		ForcePrepareEnv: argBool(argFP),
		ForceRmEnv:      argBool(argRE),
		ForceRmVolumes:  argBool(argRV),
		DryRun:          argBool(argDR),
		OtherArgs:       ddeFS.Args(),
	}
	if len(out.EnvSvr) > 0 && (len(out.Environments) > 0 || len(out.Services) > 0) {
		clog.Panic(ctx, errors.New("illegal arguments"), 111, "Cannot use -es with -e and -s")
	}
	return out
}

// Parses command line argument as a string array
func argArr(arrStr *string) []string {
	items := []string{}
	if arrStr != nil && len(*arrStr) > 0 {
		items = strings.Split(*arrStr, " ")
	}
	return items
}

// Parses command line argument as a string value
func argStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Parses command line argument as an int value
func argInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

// Parses command line argument as a bool value
func argBool(b *bool) bool {
	return b != nil && *b
}
