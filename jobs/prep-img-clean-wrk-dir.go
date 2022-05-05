package jobs

import (
	"fmt"
	"os"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
)

// Job that cleans project's working directory
type JobPrepareImageCleanDir struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
}

// Job name
func (j *JobPrepareImageCleanDir) JobName() string {
	return JobNamePrepImgCleanWrkDir
}

// Base job data
func (j *JobPrepareImageCleanDir) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareImageCleanDir) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", j.Project.Name, "Clean proj wrk dir")
}

// Job environment name
func (j *JobPrepareImageCleanDir) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobPrepareImageCleanDir) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareImageCleanDir) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobPrepareImageCleanDir) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobPrepareImageCleanDir) GetJobDescription() string {
	return "Clean proj wrk dir"
}

// Job implementation
func (j *JobPrepareImageCleanDir) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	projWorkingDirPath := cpath.ProjWorkingDir(j.Project.Name)
	// clear project working dir
	results.Append(cbash.Call(ctx, "rm -rf "+projWorkingDirPath))

	// create dirs to build image
	projWorkingSrcPath := cpath.ProjWorkingDirSrc(j.Project.Name)
	err := os.MkdirAll(projWorkingSrcPath, os.ModePerm)
	if err != nil {
		results.Append(common.CmdRes{
			StdErr: "Cannot create project src dir: " + projWorkingSrcPath,
			Err:    err,
		})
	}
	return results
}
