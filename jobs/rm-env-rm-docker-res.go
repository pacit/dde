package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
)

// Job that removes environment's docker resources
type JobRmEnvRmDockerResources struct {
	// Base job data
	BaseJ BaseJob
	// Environment info
	Environment JobEnvironmentInfo
	// Flag indicates if docker's volumes should be removed
	ForceRmVolumes bool
}

// Job name
func (j *JobRmEnvRmDockerResources) JobName() string {
	return JobNameRmEnvRmDockerRes
}

// Base job data
func (j *JobRmEnvRmDockerResources) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobRmEnvRmDockerResources) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Environment.Name, "*", "Rm docker resources")
}

// Job environment name
func (j *JobRmEnvRmDockerResources) GetEnvName() string {
	return j.Environment.Name
}

// Job service name
func (j *JobRmEnvRmDockerResources) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobRmEnvRmDockerResources) GetProjName() string {
	return ""
}

// Job project version
func (j *JobRmEnvRmDockerResources) GetProjVer() string {
	return ""
}

// Job dscription
func (j *JobRmEnvRmDockerResources) GetJobDescription() string {
	return "Rm docker resources"
}

// Job implementation
func (j *JobRmEnvRmDockerResources) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	results.AppendMulti(rmDockerResources_networks(ctx, j.Environment.Cfg.DockerResources.Networks))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(rmDockerResources_secrets(ctx, j.Environment.Cfg.DockerResources.Secrets))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(rmDockerResources_configs(ctx, j.Environment.Cfg.DockerResources.Configs))
	if results.Err != nil {
		return results
	}
	if j.ForceRmVolumes {
		results.AppendMulti(rmDockerResources_volumes(ctx, j.Environment.Cfg.DockerResources.Volumes))
	}
	return results
}
