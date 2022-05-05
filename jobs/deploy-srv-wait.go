package jobs

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/model/modelc"
)

// Job - that waits for a service has started
type JobDeployServiceWait struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
	// Max time in seconds to wait for a service
	Timeout int
}

// Job name
func (j *JobDeployServiceWait) JobName() string {
	return JobNameDeploySrvWait
}

// Base job data
func (j *JobDeployServiceWait) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobDeployServiceWait) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Wait for service")
}

// Job environment name
func (j *JobDeployServiceWait) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobDeployServiceWait) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobDeployServiceWait) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobDeployServiceWait) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobDeployServiceWait) GetJobDescription() string {
	return "Wait for service"
}

// Job implementation
func (j *JobDeployServiceWait) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	srvWorkingDirPath := cpath.SrvWorkingDir(j.Service.EnvName, j.Service.SrvName)
	isCmd := len(j.Service.Cfg.Scripts.IsRunning.Cmd) > 0
	isFile := len(j.Service.Cfg.Scripts.IsRunning.File) > 0
	scriptJob := &JobExecuteScript{
		HostWorkingDir: srvWorkingDirPath,
		ScriptCfg:      j.Service.Cfg.Scripts.IsRunning,
		Project:        &j.Project,
		Service:        &j.Service,
		Description:    "Wait for service",
	}
	if !isCmd && !isFile {
		cmd := `ok=$(docker service ls -f name=` + j.Service.EnvName + "_" + j.Service.SrvName + ` | grep "` + j.Project.Name + `:` + j.Project.VersionDockerSafe +
			`" | wc -l | xargs echo -n); if [ "$ok" == "1" ]; then exit 0; fi; exit 123`
		scriptJob = &JobExecuteScript{
			HostWorkingDir: srvWorkingDirPath,
			ScriptCfg:      modelc.ScriptJson{Cmd: cmd},
			Project:        &j.Project,
			Service:        &j.Service,
			Description:    "Wait for service",
		}
	}

	for i := 0; i < j.Timeout; i++ {
		res := scriptJob.DoIt(ctx)
		results.AppendMulti(res)
		if res.Err == nil {
			results.Err = nil
			results.ErrRes = common.CmdRes{}
			return results
		}
		time.Sleep(time.Second)
		clog.Trace(ctx, "Wait for service", j.Service.EnvName+":"+j.Service.SrvName, fmt.Sprintf("%v", i))
	}
	errMsg := "Service " + j.Service.EnvName + ":" + j.Service.SrvName + " - not respond in " + strconv.Itoa(j.Timeout) + "s"
	results.Append(common.CmdRes{
		StdErr: errMsg,
		Err:    errors.New(errMsg),
	})
	return results
}
