package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/common/ctmpl"
)

// Job that creates service's working directory
type JobDeployServiceCreateWorkFiles struct {
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
func (j *JobDeployServiceCreateWorkFiles) JobName() string {
	return JobNameDeploySrvCreateWrkFiles
}

// Base job data
func (j *JobDeployServiceCreateWorkFiles) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobDeployServiceCreateWorkFiles) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Create wrk files")
}

// Job environment name
func (j *JobDeployServiceCreateWorkFiles) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobDeployServiceCreateWorkFiles) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobDeployServiceCreateWorkFiles) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobDeployServiceCreateWorkFiles) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobDeployServiceCreateWorkFiles) GetJobDescription() string {
	return "Create wrk files"
}

// Job implementation
func (j *JobDeployServiceCreateWorkFiles) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	srvWorkingDirPath := cpath.SrvWorkingDir(j.Service.EnvName, j.Service.SrvName)
	srvDefinitionDirPath := cpath.SrvDefinitionDir(j.Service.EnvName, j.Service.SrvName)

	results.Append(cbash.Call(ctx, "/bin/cp -f "+srvDefinitionDirPath+"/* "+srvWorkingDirPath+"/"))
	if results.Err != nil {
		return results
	}
	//execute templates
	tmplOut, err := ctmpl.CompileTmplFilesInDir(ctx, srvWorkingDirPath, j.Properties)
	results.Append(common.CmdRes{
		StdOut: tmplOut,
		Err:    err,
	})
	return results
}
