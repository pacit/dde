package jobs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/common/ctmpl"
)

// Job that creates eorkspace's working directory
type JobRunCustomScriptCreateWrkFiles struct {
	// Base job data
	BaseJ BaseJob
	// Workspace info
	Workspace JobWorkspaceInfo
	// Properties to use in templates
	Properties map[string]string
}

// Job name
func (j *JobRunCustomScriptCreateWrkFiles) JobName() string {
	return JobNameRunCustomScriptCreateWrkFiles
}

// Base job data
func (j *JobRunCustomScriptCreateWrkFiles) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobRunCustomScriptCreateWrkFiles) ShortInfo() string {
	return fmt.Sprintf("[%-21s] - %s", "runs", "Create workspace files")
}

// Job environment name
func (j *JobRunCustomScriptCreateWrkFiles) GetEnvName() string {
	return ""
}

// Job service name
func (j *JobRunCustomScriptCreateWrkFiles) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobRunCustomScriptCreateWrkFiles) GetProjName() string {
	return ""
}

// Job project version
func (j *JobRunCustomScriptCreateWrkFiles) GetProjVer() string {
	return ""
}

// Job dscription
func (j *JobRunCustomScriptCreateWrkFiles) GetJobDescription() string {
	return "Create workspace files"
}

// Job implementation
func (j *JobRunCustomScriptCreateWrkFiles) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	wrkWorkingDirPath := cpath.WrkWorkingDir()
	wrkDefinitionDirPath := cpath.WrkDefinitionDir()
	// copy files
	for _, dirRelPath := range j.Workspace.Cfg.TemplateFileDirs {
		fromDir := filepath.Clean(filepath.Join(wrkDefinitionDirPath, dirRelPath))
		toDir := filepath.Clean(filepath.Join(wrkWorkingDirPath, dirRelPath))
		// create dirs
		err := os.MkdirAll(toDir, os.ModePerm)
		if err != nil {
			results.Append(common.CmdRes{
				StdErr: "Cannot create wrk dir: " + toDir,
				Err:    err,
			})
			return results
		}
		results.Append(cbash.Call(ctx, "/bin/cp -rf "+fromDir+"/* "+toDir+"/"))
		if results.Err != nil {
			return results
		}
	}
	for _, dirRelPath := range j.Workspace.Cfg.TemplateFileDirs {
		//execute templates
		tmplDir := filepath.Clean(filepath.Join(wrkWorkingDirPath, dirRelPath))
		tmplOut, err := ctmpl.CompileTmplFilesInDir(ctx, tmplDir, j.Properties)
		results.Append(common.CmdRes{
			StdOut: tmplOut,
			Err:    err,
		})
		if results.Err != nil {
			return results
		}
	}

	return results
}
