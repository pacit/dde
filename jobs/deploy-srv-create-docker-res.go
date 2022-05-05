package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cpath"
)

// Job that creates service's docker resources
type JobDeployServiceCreateDockerResources struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
	// Force reload available docker's resources and check before creation
	ForceCheckResourceExists bool
}

// Job name
func (j *JobDeployServiceCreateDockerResources) JobName() string {
	return JobNameDeploySrvCreateDockerRes
}

// Base job data
func (j *JobDeployServiceCreateDockerResources) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobDeployServiceCreateDockerResources) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Create docker resources")
}

// Job environment name
func (j *JobDeployServiceCreateDockerResources) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobDeployServiceCreateDockerResources) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobDeployServiceCreateDockerResources) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobDeployServiceCreateDockerResources) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobDeployServiceCreateDockerResources) GetJobDescription() string {
	return "Create docker resources"
}

// Job implementation
func (j *JobDeployServiceCreateDockerResources) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	srvWorkingDirPath := cpath.SrvWorkingDir(j.Service.EnvName, j.Service.SrvName)

	results.AppendMulti(createDockerResources_networks(ctx, j.Service.Cfg.DockerResources.Networks, j.ForceCheckResourceExists))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(createDockerResources_volumes(ctx, j.Service.Cfg.DockerResources.Volumes, j.ForceCheckResourceExists))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(createDockerResources_secrets(ctx, j.Service.Cfg.DockerResources.Secrets, srvWorkingDirPath, j.ForceCheckResourceExists))
	if results.Err != nil {
		return results
	}
	results.AppendMulti(createDockerResources_configs(ctx, j.Service.Cfg.DockerResources.Configs, srvWorkingDirPath, j.ForceCheckResourceExists))

	return results
}
