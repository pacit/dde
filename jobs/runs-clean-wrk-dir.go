package jobs

import (
	"fmt"
	"os"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
)

// Job that cleans workspace's working directory
type JobRunCustomScriptCleanDir struct {
	// Base job data
	BaseJ BaseJob
}

// Job name
func (j *JobRunCustomScriptCleanDir) JobName() string {
	return JobNameRunCustomScriptCleanDir
}

// Base job data
func (j *JobRunCustomScriptCleanDir) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobRunCustomScriptCleanDir) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", "runs", "Clean workspace wrk dir")
}

// Job environment name
func (j *JobRunCustomScriptCleanDir) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobRunCustomScriptCleanDir) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobRunCustomScriptCleanDir) GetProjName() string {
	return ""
}

// Job project version
func (j *JobRunCustomScriptCleanDir) GetProjVer() string {
	return ""
}

// Job dscription
func (j *JobRunCustomScriptCleanDir) GetJobDescription() string {
	return "Clean workspace wrk dir"
}

// Job implementation
func (j *JobRunCustomScriptCleanDir) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	wrkWorkingDirPath := cpath.WrkWorkingDir()
	// clear workspace working dir
	results.Append(cbash.Call(ctx, "rm -rf "+wrkWorkingDirPath))
	if results.Err != nil {
		return results
	}

	// create dirs
	err := os.MkdirAll(wrkWorkingDirPath, os.ModePerm)
	if err != nil {
		results.Append(common.CmdRes{
			StdErr: "Cannot create wrk working dir: " + wrkWorkingDirPath,
			Err:    err,
		})
	}
	return results
}
