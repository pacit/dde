package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
)

// Job that removes service's docker resources
type JobRmServiceRmDockerResources struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
	// Flag indicates if docker's volumes should be removed
	ForceRmVolumes bool
}

// Job name
func (j *JobRmServiceRmDockerResources) JobName() string {
	return JobNameRmSrvRmDockerRes
}

// Base job data
func (j *JobRmServiceRmDockerResources) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobRmServiceRmDockerResources) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Rm docker resources")
}

// Job environment name
func (j *JobRmServiceRmDockerResources) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobRmServiceRmDockerResources) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobRmServiceRmDockerResources) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobRmServiceRmDockerResources) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobRmServiceRmDockerResources) GetJobDescription() string {
	return "Rm docker resources"
}

// Job implementation
func (j *JobRmServiceRmDockerResources) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	results.AppendMulti(rmDockerResources_networks(ctx, j.Service.Cfg.DockerResources.Networks))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(rmDockerResources_secrets(ctx, j.Service.Cfg.DockerResources.Secrets))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(rmDockerResources_configs(ctx, j.Service.Cfg.DockerResources.Configs))
	if results.Err != nil {
		return results
	}
	if j.ForceRmVolumes {
		results.AppendMulti(rmDockerResources_volumes(ctx, j.Service.Cfg.DockerResources.Volumes))
	}
	return results
}
