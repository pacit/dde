package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
	"github.com/pacit/dde/model/modelc"
)

// Service configuration file name
const ServiceJsonFileName = "srv.json"

// Service configuration
type ServiceJson struct {
	// Project name which image will be used to run service
	Project string `json:"project"`
	// Service custom scripts for run, check, before/after run actions
	Scripts ServiceJsonScripts `json:"scripts"`
	// Service's docker resources
	DockerResources modelc.DockerResourcesJson `json:"dockerResources"`
	// Service run configuration. It is used when no custom script to run exists.
	DockerRun modelc.DockerRunJson `json:"dockerRun"`
	// Services must run before. When no service found - error occured.
	//
	// The format is: {envName}:{srvName}
	DependsOnSrv []string `json:"dependsOnSrv"`
	// Custom timeout to wait for service runs [s].
	WaitForServiceTimeoutS int `json:"waitForServiceTimeoutS"`
	// Service properties to use in templates
	Properties map[string]string `json:"properties"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (sj *ServiceJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	sj.Project = ctmpl.CompileStringOrDie(ctx, sj.Project, props)
	(&sj.Scripts.IsRunning).CompileTemplatesOrDie(ctx, props)
	for i, b := range sj.Scripts.BeforeRun {
		(&b).CompileTemplatesOrDie(ctx, props)
		sj.Scripts.BeforeRun[i] = b
	}
	(&sj.Scripts.Run).CompileTemplatesOrDie(ctx, props)
	for i, a := range sj.Scripts.AfterRun {
		(&a).CompileTemplatesOrDie(ctx, props)
		sj.Scripts.AfterRun[i] = a
	}
	(&sj.DockerResources).CompileTemplatesOrDie(ctx, props)
	(&sj.DockerRun).CompileTemplatesOrDie(ctx, props)
	ctmpl.CompileStringArrOrDie(ctx, sj.DependsOnSrv, props)
}
