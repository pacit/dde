package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/model"
)

// Job that pulls project's docker image from remote repository
type JobPrepareImagePullImage struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Docker image description
	Image *model.ProjectJsonDockerImage
}

// Job name
func (j *JobPrepareImagePullImage) JobName() string {
	return JobNamePrepImgPull
}

// Base job data
func (j *JobPrepareImagePullImage) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareImagePullImage) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", j.Project.Name, "Pull image from "+j.Image.RepoId)
}

// Job environment name
func (j *JobPrepareImagePullImage) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobPrepareImagePullImage) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareImagePullImage) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobPrepareImagePullImage) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobPrepareImagePullImage) GetJobDescription() string {
	return "Pull image from " + j.Image.RepoId
}

// Job implementation
func (j *JobPrepareImagePullImage) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	cmd := "docker image pull " + j.Image.ImageName + ":" + j.Project.VersionDockerSafe
	results.Append(cbash.Call(ctx, cmd))
	return results
}
