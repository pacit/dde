package jobs

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/common/ctmpl"
)

// Job that creates project's working directory
type JobPrepareImageCreateWorkFiles struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Properties to use in templates
	Properties map[string]string
}

// Job name
func (j *JobPrepareImageCreateWorkFiles) JobName() string {
	return JobNamePrepImgCreateWrkFiles
}

// Base job data
func (j *JobPrepareImageCreateWorkFiles) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareImageCreateWorkFiles) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", j.Project.Name, "Create wrk files")
}

// Job environment name
func (j *JobPrepareImageCreateWorkFiles) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobPrepareImageCreateWorkFiles) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareImageCreateWorkFiles) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobPrepareImageCreateWorkFiles) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobPrepareImageCreateWorkFiles) GetJobDescription() string {
	return "Create wrk files"
}

// Job implementation
func (j *JobPrepareImageCreateWorkFiles) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	projWorkingDirPath := cpath.ProjWorkingDir(j.Project.Name)
	projDefinitionDirPath := cpath.ProjDefinitionDir(j.Project.Name)
	results.Append(cbash.Call(ctx, "/bin/cp -f "+projDefinitionDirPath+"/* "+projWorkingDirPath+"/"))
	if results.Err != nil {
		return results
	}

	//execute templates
	tmplOut, err := ctmpl.CompileTmplFilesInDir(ctx, projWorkingDirPath, j.Properties)
	results.Append(common.CmdRes{
		StdOut: tmplOut,
		Err:    err,
	})
	return results
}
