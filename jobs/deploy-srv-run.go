package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/model/modelc"
)

// Job that runs service
type JobDeployServiceRun struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
	// Properties to use in templates
	Properties map[string]string
}

// Job name
func (j *JobDeployServiceRun) JobName() string {
	return JobNameDeploySrvRun
}

// Base job data
func (j *JobDeployServiceRun) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobDeployServiceRun) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Run service")
}

// Job environment name
func (j *JobDeployServiceRun) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobDeployServiceRun) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobDeployServiceRun) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobDeployServiceRun) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobDeployServiceRun) GetJobDescription() string {
	return "Run service"
}

// Job implementation
func (j *JobDeployServiceRun) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	srvWorkingDirPath := cpath.SrvWorkingDir(j.Service.EnvName, j.Service.SrvName)
	isCmd := len(j.Service.Cfg.Scripts.Run.Cmd) > 0
	clog.Trace(ctx, fmt.Sprint(j.Properties))
	j.Service.Cfg.CompileTemplatesOrDie(ctx, j.Properties)
	clog.Trace(ctx, fmt.Sprint(j.Service.Cfg))
	isFile := len(j.Service.Cfg.Scripts.Run.File) > 0
	if isCmd || isFile {
		scriptJob := &JobExecuteScript{
			HostWorkingDir: srvWorkingDirPath,
			ScriptCfg:      j.Service.Cfg.Scripts.Run,
			Project:        &j.Project,
			Service:        &j.Service,
			Description:    "Run service",
		}
		results.AppendMulti(scriptJob.DoIt(ctx))
	} else {
		// use builtin run command
		cmd := j.Service.Cfg.DockerRun.GetCommand(j.Project.Name, j.Project.VersionDockerSafe)
		scriptJob := &JobExecuteScript{
			HostWorkingDir: srvWorkingDirPath,
			ScriptCfg:      modelc.ScriptJson{Cmd: cmd},
			Project:        &j.Project,
			Service:        &j.Service,
			Description:    "Run service",
		}
		results.AppendMulti(scriptJob.DoIt(ctx))
	}
	return results
}
