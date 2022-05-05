package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cpath"
)

// Job that creates environment's docker resources
type JobPrepareEnvCreateDockerResources struct {
	// Base job data
	BaseJ BaseJob
	// Environment info
	Environment JobEnvironmentInfo
}

// Job name
func (j *JobPrepareEnvCreateDockerResources) JobName() string {
	return JobNamePrepEnvCreateDockerRes
}

// Base job data
func (j *JobPrepareEnvCreateDockerResources) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareEnvCreateDockerResources) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Environment.Name, "*", "Create docker resources")
}

// Job environment name
func (j *JobPrepareEnvCreateDockerResources) GetEnvName() string {
	return j.Environment.Name
}

// Job service name
func (j *JobPrepareEnvCreateDockerResources) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareEnvCreateDockerResources) GetProjName() string {
	return ""
}

// Job project version
func (j *JobPrepareEnvCreateDockerResources) GetProjVer() string {
	return ""
}

// Job dscription
func (j *JobPrepareEnvCreateDockerResources) GetJobDescription() string {
	return "Create docker resources"
}

// Job implementation
func (j *JobPrepareEnvCreateDockerResources) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	envWorkingDirPath := cpath.EnvWorkingDir(j.Environment.Name)
	results.AppendMulti(createDockerResources_networks(ctx, j.Environment.Cfg.DockerResources.Networks, false))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(createDockerResources_volumes(ctx, j.Environment.Cfg.DockerResources.Volumes, false))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(createDockerResources_secrets(ctx, j.Environment.Cfg.DockerResources.Secrets, envWorkingDirPath, false))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(createDockerResources_configs(ctx, j.Environment.Cfg.DockerResources.Configs, envWorkingDirPath, false))
	return results
}
