package jobs

import (
	"fmt"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/common/ctmpl"
)

// Job that creates environment's working directory
type JobPrepareEnvCreateWorkFiles struct {
	// Base job data
	BaseJ BaseJob
	// Environment info
	Environment JobEnvironmentInfo
	// Properties to use in templates
	Properties map[string]string
}

// Job name
func (j *JobPrepareEnvCreateWorkFiles) JobName() string {
	return JobNamePrepEnvCreateWrkFiles
}

// Base job data
func (j *JobPrepareEnvCreateWorkFiles) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobPrepareEnvCreateWorkFiles) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Environment.Name, "*", "Create env files")
}

// Job environment name
func (j *JobPrepareEnvCreateWorkFiles) GetEnvName() string {
	return j.Environment.Name
}

// Job service name
func (j *JobPrepareEnvCreateWorkFiles) GetSrvName() string {
	return ""
}

// Job project name
func (j *JobPrepareEnvCreateWorkFiles) GetProjName() string {
	return ""
}

// Job project version
func (j *JobPrepareEnvCreateWorkFiles) GetProjVer() string {
	return ""
}

// Job dscription
func (j *JobPrepareEnvCreateWorkFiles) GetJobDescription() string {
	return "Create env files"
}

// Job implementation
func (j *JobPrepareEnvCreateWorkFiles) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	envWorkingDirPath := cpath.EnvWorkingDir(j.Environment.Name)
	envDefinitionDirPath := cpath.EnvDefinitionDir(j.Environment.Name)
	results.Append(cbash.Call(ctx, "/bin/cp -rf "+envDefinitionDirPath+"/* "+envWorkingDirPath+"/"))
	if results.Err != nil {
		return results
	}
	if !common.StringSliceContains(j.Environment.Cfg.TemplateFileDirs, ".") {
		j.Environment.Cfg.TemplateFileDirs = append(j.Environment.Cfg.TemplateFileDirs, ".")
	}
	for _, dirRelPath := range j.Environment.Cfg.TemplateFileDirs {
		//execute templates
		tmplDir := filepath.Clean(filepath.Join(envWorkingDirPath, dirRelPath))
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
