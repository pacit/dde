package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
)

// Job that clones project's git repository
type JobPrepareImageGitClone struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
}

// Job name
func (j *JobPrepareImageGitClone) JobName() string {
	return JobNamePrepImgGitClone
}

// Base job data
func (j *JobPrepareImageGitClone) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareImageGitClone) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", j.Project.Name, "Clone git repo")
}

// Job environment name
func (j *JobPrepareImageGitClone) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobPrepareImageGitClone) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareImageGitClone) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobPrepareImageGitClone) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobPrepareImageGitClone) GetJobDescription() string {
	return "Clone git repo"
}

// Job implementation
func (j *JobPrepareImageGitClone) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	projWorkingSrcPath := cpath.ProjWorkingDirSrc(j.Project.Name)
	if len(j.Project.Cfg.GitRepo) > 0 {
		cmd := "cd " + projWorkingSrcPath +
			" && git clone --recurse-submodules --depth 1 --shallow-submodules -b " + j.Project.Version + " " + j.Project.Cfg.GitRepo
		results.Append(cbash.Call(ctx, cmd))
	}
	return results
}
