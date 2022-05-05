package jobs

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Jobs to do stream
var ReqCH = make(chan IJob, 100)

// Job's responses stream
var RespCH = make(chan IJob, 100)

// How many jobs are processing now
var RunningJobs = 0

// Channel to send stop event to all runners
var stopCH = make(chan bool, 100)

// Number of working runners
var runnersNo = 0

// Starts runners (new threads)
func StartRunners(ctx common.DCtx, n int) {
	clog.Trace(ctx, "StartRunners", strconv.Itoa(n))
	for i := 1; i <= n; i++ {
		go goRunner(ctx, i)
	}
	runnersNo = n
}

// Stops all working runners
func StopAllRunners(ctx common.DCtx) {
	clog.Trace(ctx, "StopAllRunners()")
	for i := 0; i < runnersNo; i++ {
		stopCH <- true
	}
	runnersNo = 0
}

// Runner logic (loop and get job to do)
func goRunner(rCtx common.DCtx, rn int) {
	rCtx.ThreadId = "R-" + fmt.Sprint(rn)
	clog.Debug(rCtx, "Start")
el:
	for {
		select {
		case <-stopCH:
			break el
		case j := <-ReqCH:
			RunningJobs += 1
			clog.Debug(rCtx, "Gets job", j.Base().Id)
			jCtx := createJobCtx(rCtx, j)
			j = processJob(jCtx, j)
			RespCH <- j
			RunningJobs -= 1
		}
	}
	clog.Debug(rCtx, "Stop")
}

// Creates job context
func createJobCtx(ctx common.DCtx, j IJob) common.DCtx {
	return common.DCtx{
		ThreadId: ctx.ThreadId,
		JobId:    j.Base().Id,
		SrvName:  j.GetSrvName(),
		ProjName: j.GetProjName(),
		EnvName:  j.GetEnvName(),
	}
}

// Processing a job
func processJob(ctx common.DCtx, j IJob) IJob {
	startTime := time.Now()
	j.Base().StartTime = startTime
	pingStopCH := make(chan bool)
	go runJobPingLog(ctx, j.ShortInfo(), startTime, pingStopCH)
	result := j.DoIt(ctx)
	j.Base().Result = result

	if result.Err != nil {
		clog.Error(ctx, result.Err, "Job error while processing")
		j.Base().State = JobStateError
	} else {
		j.Base().State = JobStateDone
	}
	pingStopCH <- true
	clog.Debug(ctx, "finish job", j.Base().Id, j.Base().State)
	j.Base().TimeDuration = time.Since(startTime)
	return j
}

// Prints a log message that a runner processing a job - once a 5 sec
func runJobPingLog(ctx common.DCtx, jobDesc string, startTime time.Time, stopCH chan bool) {
el:
	for {
		select {
		case <-stopCH:
			break el
		default:
			time.Sleep(100 * time.Millisecond)
			timeS := time.Since(startTime).Seconds()
			if math.Remainder(timeS, 5) == 0 {
				clog.Info(ctx, "running", jobDesc, fmt.Sprintf("%v", time.Since(startTime).Truncate(time.Millisecond)))
			}
		}
	}
}
