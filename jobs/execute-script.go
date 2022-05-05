package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/model/modelc"
)

// Job that executes a script
type JobExecuteScript struct {
	// Base job data
	BaseJ BaseJob
	// Project info (if script refers to a project)
	Project *JobProjectInfo
	// Service info (if script refers to a project)
	Service *JobServiceInfo
	// Environment info (if script refers to a project)
	Environment *JobEnvironmentInfo
	// Script working directory
	HostWorkingDir string
	// Script configuration
	ScriptCfg modelc.ScriptJson
	// Script description (job description)
	Description string
}

// Job name
func (j *JobExecuteScript) JobName() string {
	return JobNameExecuteScript
}

// Base job data
func (j *JobExecuteScript) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobExecuteScript) ShortInfo() string {
	if j.Environment != nil {
		return fmt.Sprintf("[%6s:%-14s] - %s", j.Environment.Name, "*", j.Description)
	} else if j.Service != nil {
		return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, j.Description)
	} else if j.Project != nil {
		return fmt.Sprintf("[%-21s] - %s", j.Project.Name, j.Description)
	}
	return fmt.Sprintf("[%6s:%-14s] - %s", "-", "-", j.Description)
}

// Job environment name
func (j *JobExecuteScript) GetEnvName() string {
	if j.Environment != nil {
		return j.Environment.Name
	} else if j.Service != nil {
		return j.Service.EnvName
	}
	return ""
}

// Job service name
func (j *JobExecuteScript) GetSrvName() string {
	if j.Service != nil {
		return j.Service.SrvName
	}
	return ""
}

// Job project name
func (j *JobExecuteScript) GetProjName() string {
	if j.Project != nil {
		return j.Project.Name
	}
	return ""
}

// Job project version
func (j *JobExecuteScript) GetProjVer() string {
	if j.Project != nil {
		return j.Project.Version
	}
	return ""
}

// Job dscription
func (j *JobExecuteScript) GetJobDescription() string {
	return j.Description
}

// Job implementation
func (j *JobExecuteScript) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(j.ScriptCfg.RunIn) == 0 {
		// run on host
		command := j.ScriptCfg.GetCommand(j.HostWorkingDir)
		results.Append(cbash.Call(ctx, command))
	} else {
		// run in container
		panic("not implemented jet")
	}
	return results
}
