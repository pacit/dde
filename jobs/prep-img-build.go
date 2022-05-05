package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
)

// Job that builds project's docker image
type JobPrepareImageBuild struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
}

// Job name
func (j *JobPrepareImageBuild) JobName() string {
	return JobNamePrepImgBuild
}

// Base job data
func (j *JobPrepareImageBuild) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareImageBuild) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", j.Project.Name, "Build image")
}

// Job environment name
func (j *JobPrepareImageBuild) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobPrepareImageBuild) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareImageBuild) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobPrepareImageBuild) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobPrepareImageBuild) GetJobDescription() string {
	return "Build image"
}

// Job implementation
func (j *JobPrepareImageBuild) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	projWorkingDirPath := cpath.ProjWorkingDir(j.Project.Name)
	command := j.Project.Cfg.Scripts.Build.GetCommand(projWorkingDirPath)
	results.Append(cbash.Call(ctx, command))
	return results
}
