package jobs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
)

// Job that cleans environment working directory
type JobPrepareEnvCleanDir struct {
	// Base job data
	BaseJ BaseJob
	// Environment info
	Environment JobEnvironmentInfo
}

// Job name
func (j *JobPrepareEnvCleanDir) JobName() string {
	return JobNamePrepEnvCleanWrkDir
}

// Base job data
func (j *JobPrepareEnvCleanDir) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareEnvCleanDir) ShortInfo() string {
	return fmt.Sprintf("[%s:%s] - %s", j.Environment.Name, "*", "Clean env dir")
}

// Job environment name
func (j *JobPrepareEnvCleanDir) GetEnvName() string {
	return j.Environment.Name
}

// Job service name
func (j *JobPrepareEnvCleanDir) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareEnvCleanDir) GetProjName() string {
	return ""
}

// Job project version
func (j *JobPrepareEnvCleanDir) GetProjVer() string {
	return ""
}

// Job dscription
func (j *JobPrepareEnvCleanDir) GetJobDescription() string {
	return "Clean env dir"
}

// Job implementation
func (j *JobPrepareEnvCleanDir) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	envWorkingDirPath := cpath.EnvWorkingDir(j.Environment.Name)
	// clear environment working dir
	results.Append(cbash.Call(ctx, "rm -rf "+envWorkingDirPath))
	if results.Err != nil {
		return results
	}
	// create dirs to run service
	err := os.MkdirAll(filepath.Join(envWorkingDirPath), os.ModePerm)
	if err != nil {
		results.Append(common.CmdRes{
			StdErr: "Cannot create env working dir: " + envWorkingDirPath,
			Err:    err,
		})
	}
	return results
}
