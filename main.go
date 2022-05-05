package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/docker"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/steps"
)

// Application entry point
func main() {
	ctx := common.DCtx{}
	startTime := time.Now()
	args := parseArguments(ctx)
	clog.Configure(args.LogLevels)
	if clog.PrintBaner {
		clog.Info(ctx, "**********************************")
		clog.Info(ctx, "*  Docker Dev Environment 0.0.1  *")
		clog.Info(ctx, "**********************************\n")
	}

	wrk := createNewEngine(ctx)
	setParams(ctx, wrk, args)
	if err := wrk.LoadSelectedEnvironments(ctx); err != nil {
		clog.Panic(ctx, err, 120, "Error loading selected environments")
	}
	if err := wrk.LoadSelectedServices(ctx); err != nil {
		clog.Panic(ctx, err, 120, "Error loading selected Services")
	}
	if err := wrk.LoadRequiredProjects(ctx); err != nil {
		clog.Panic(ctx, err, 120, "Error loading required projects")
	}
	if err := wrk.CreateVersionsMap(ctx); err != nil {
		clog.Panic(ctx, err, 120, "Error creating versions map")
	}
	if err := docker.LoadAllAvailableResources(ctx); err != nil {
		clog.Panic(ctx, err, 120, "Error loading dorker available resources")
	}
	if err := wrk.CompileTemplatesInCfgJsons(ctx); err != nil {
		clog.Panic(ctx, err, 120, "Error compile templates in cfg jsons")
	}
	if err := steps.AddAllSteps(ctx, wrk); err != nil {
		clog.Panic(ctx, err, 120, "Error creating steps")
	}
	steps.PrintJobsTableToDo(ctx)

	if !args.DryRun {
		steps.RunJobs(ctx, wrk.CliArgs.Threads)
		steps.PrintJobsTableState(ctx)
		isErr := steps.PrintErrors(ctx)
		if isErr {
			clog.Warning(ctx, "Error run in "+fmt.Sprintf("%v", time.Since(startTime).Truncate(time.Millisecond)))
		} else {
			clog.Info(ctx, "Success run in "+fmt.Sprintf("%v", time.Since(startTime).Truncate(time.Millisecond)))
		}
	} else {
		clog.Info(ctx, "Done dry run in "+fmt.Sprintf("%v", time.Since(startTime).Truncate(time.Millisecond)))
	}

	time.Sleep(500 * time.Millisecond)
}

// Creates new workspace engine
func createNewEngine(ctx common.DCtx) *engine.Workspace {
	wrk, err := engine.New(ctx)
	if err != nil {
		clog.Panic(ctx, err, 119, "ERROR loading workspace")
	}
	return wrk
}

// Sets command line arguments and calculates required environments and services
func setParams(ctx common.DCtx, w *engine.Workspace, args engine.CliArgs) {
	w.CliArgs = args
	if len(w.CliArgs.VersionFile) == 0 {
		if len(w.Cfg.DefaultVerFile) > 0 {
			w.CliArgs.VersionFile = w.Cfg.DefaultVerFile
		}
	}
	setSelectedEnvironments(ctx, w, args)
	setSelectedServices(ctx, w, args)
}

// Calculates required environments
func setSelectedEnvironments(ctx common.DCtx, w *engine.Workspace, args engine.CliArgs) {
	if len(args.EnvSvr) > 0 {
		// environments used in -es list
		for _, es := range args.EnvSvr {
			split := strings.Split(es, ":")
			if len(split) == 2 {
				w.SelectedEnvironments = append(w.SelectedEnvironments, split[0])
			}
		}
	} else if len(args.Environments) > 0 {
		// list environments from arg -e
		w.SelectedEnvironments = args.Environments
	} else if len(w.Cfg.DefaultEnvs) > 0 {
		// default list from wrk.json
		w.SelectedEnvironments = w.Cfg.DefaultEnvs
	} else {
		// all environments
		all, err := w.GetAllEnvNames()
		if err != nil {
			clog.Panic(ctx, err, 120, "Cannot read environment names")
		}
		w.SelectedEnvironments = all
	}
	clog.Info(ctx, "Selected environments", fmt.Sprintf("%v", w.SelectedEnvironments))
}

// Calculates required services
func setSelectedServices(ctx common.DCtx, w *engine.Workspace, args engine.CliArgs) {
	if len(args.EnvSvr) > 0 {
		// services used in -es list
		w.SelectedServices = args.EnvSvr
	} else if len(args.Services) > 0 {
		// list services from arg -s (only in selected environments)
		allNames := w.GetAllServiceNamesInSelectedEnvs(ctx)
		w.SelectedServices = []string{}
		for _, name := range allNames {
			envName := strings.Split(name, ":")[0]
			if common.StringSliceContains(w.SelectedEnvironments, envName) {
				w.SelectedServices = append(w.SelectedServices, name)
			}
		}
	} else {
		// all services in selected enviromnemts
		w.SelectedServices = w.GetAllServiceNamesInSelectedEnvs(ctx)
	}
	for _, name := range w.SelectedServices {
		srvName := strings.Split(name, ":")[1]
		if !common.StringSliceContains(w.SelectedServicesOnlyNames, srvName) {
			w.SelectedServicesOnlyNames = append(w.SelectedServicesOnlyNames, srvName)
		}
	}
	clog.Info(ctx, "Selected services (only names)", fmt.Sprintf("%v", w.SelectedServicesOnlyNames))
	clog.Info(ctx, "Selected services", fmt.Sprintf("%v", w.SelectedServices))
}
