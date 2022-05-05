package jobs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
)

// Job that cleans service working directory
type JobDeployServiceCleanDir struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
}

// Job name
func (j *JobDeployServiceCleanDir) JobName() string {
	return JobNameDeploySrvCleanDir
}

// Base job data
func (j *JobDeployServiceCleanDir) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobDeployServiceCleanDir) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Clean srv wrk dir")
}

// Job environment name
func (j *JobDeployServiceCleanDir) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobDeployServiceCleanDir) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobDeployServiceCleanDir) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobDeployServiceCleanDir) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobDeployServiceCleanDir) GetJobDescription() string {
	return "Clean srv wrk dir"
}

// Job implementation
func (j *JobDeployServiceCleanDir) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	srvWorkingDirPath := cpath.SrvWorkingDir(j.Service.EnvName, j.Service.SrvName)
	// clear service working dir
	results.Append(cbash.Call(ctx, "rm -rf "+srvWorkingDirPath))
	if results.Err != nil {
		return results
	}

	// create dirs to run service
	err := os.MkdirAll(filepath.Join(srvWorkingDirPath), os.ModePerm)
	if err != nil {
		results.Append(common.CmdRes{
			StdErr: "Error while creating srv working dir: " + srvWorkingDirPath,
			Err:    err,
		})
		return results
	}
	return results
}
