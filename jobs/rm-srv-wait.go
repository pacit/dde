package jobs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/clog"
)

// Job - that waits while a service has removed
type JobRmServiceWait struct {
	// Base job data
	BaseJ BaseJob
	// Project info
	Project JobProjectInfo
	// Service info
	Service JobServiceInfo
}

// Job name
func (j *JobRmServiceWait) JobName() string {
	return JobNameRmSrvWait
}

// Base job data
func (j *JobRmServiceWait) Base() *BaseJob {
	return &j.BaseJ
}

// Short job info
func (j *JobRmServiceWait) ShortInfo() string {
	return fmt.Sprintf("[%6s:%-14s] - %s", j.Service.EnvName, j.Service.SrvName, "Wait for srv remove")
}

// Job environment name
func (j *JobRmServiceWait) GetEnvName() string {
	return j.Service.EnvName
}

// Job service name
func (j *JobRmServiceWait) GetSrvName() string {
	return j.Service.SrvName
}

// Job project name
func (j *JobRmServiceWait) GetProjName() string {
	return j.Project.Name
}

// Job project version
func (j *JobRmServiceWait) GetProjVer() string {
	return j.Project.Version
}

// Job dscription
func (j *JobRmServiceWait) GetJobDescription() string {
	return "Wait for srv remove"
}

// Job implementation
func (j *JobRmServiceWait) DoIt(ctx common.DCtx) common.CmdResMulti {
	results := common.CmdResMulti{}
	serviceName := j.Service.EnvName + "_" + j.Service.SrvName
	serviceRemoved := false
	// wait for service remove
	cmd := "docker service ls -f name=" + serviceName + " | grep \"" + serviceName + " \" | wc -l"
	for i := 0; i < 60; i++ {
		result := cbash.Call(ctx, cmd)
		results.Append(result)
		if strings.Trim(result.StdOut, " \r\n\t") == "0" && result.Err == nil {
			serviceRemoved = true
			break
		}
		time.Sleep(time.Second)
		clog.Trace(ctx, "Wait for service rm", j.Service.EnvName+":"+j.Service.SrvName, fmt.Sprintf("%v", i))
	}
	if serviceRemoved {
		// wait for container remove
		cmd := "docker container ls -f \"name=" + serviceName + ".*\" | grep " + serviceName + " | wc -l"
		for i := 0; i < 60; i++ {
			result := cbash.Call(ctx, cmd)
			results.Append(result)
			if strings.Trim(result.StdOut, " \r\n\t") == "0" && result.Err == nil {
				time.Sleep(3 * time.Second)
				return results
			}
			time.Sleep(time.Second)
			clog.Trace(ctx, "Wait for container rm", j.Service.EnvName+":"+j.Service.SrvName, fmt.Sprintf("%v", i))
		}
	}
	errMsg := "Service (or container) " + j.Service.EnvName + ":" + j.Service.SrvName + " - not removed in " + strconv.Itoa(60) + "s"
	results.Append(common.CmdRes{
		StdErr: errMsg,
		Err:    errors.New(errMsg),
	})
	return results
}
