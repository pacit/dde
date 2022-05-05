package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
)

// Job that removes service's docker service
type JobRmServiceRmService struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
}

// Job name
func (j *JobRmServiceRmService) JobName() string {
	return JobNameRmSrvRmService
}

// Base job data
func (j *JobRmServiceRmService) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobRmServiceRmService) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Rm service")
}

// Job environment name
func (j *JobRmServiceRmService) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobRmServiceRmService) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobRmServiceRmService) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobRmServiceRmService) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobRmServiceRmService) GetJobDescription() string {
	return "Rm service"
}

// Job implementation
func (j *JobRmServiceRmService) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	serviceName := j.Service.EnvName + "_" + j.Service.SrvName
	cmd := "docker service rm " + serviceName
	results.Append(cbash.Call(ctx, cmd))
	return results
}
