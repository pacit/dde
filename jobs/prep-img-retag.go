package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
)

// Job that tags docker image with a different name as a project image
type JobPrepareImageRetag struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Image name to tag as a project image
	FromName string
}

// Job name
func (j *JobPrepareImageRetag) JobName() string {
	return JobNamePrepImgRetag
}

// Base job data
func (j *JobPrepareImageRetag) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareImageRetag) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", j.Project.Name, "Retag existing image")
}

// Job environment name
func (j *JobPrepareImageRetag) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobPrepareImageRetag) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareImageRetag) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobPrepareImageRetag) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobPrepareImageRetag) GetJobDescription() string {
	return "Retag existing image"
}

// Job implementation
func (j *JobPrepareImageRetag) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	cmd := "docker image tag " + j.FromName + ":" + j.Project.VersionDockerSafe + " " + j.Project.Name + ":" + j.Project.VersionDockerSafe
	results.Append(cbash.Call(ctx, cmd))
	return results
}
