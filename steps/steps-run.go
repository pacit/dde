package steps

import (
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/jobs"
)

// Runs all existings jobs
func RunJobs(ctx common.DCtx, threads int) {
	clog.Debug(ctx, "RunJobs()")
	allJobsResponsePresent := make(chan bool)
	go listenJobsResponse(ctx, allJobsResponsePresent)
	go jobs.StartRunners(ctx, threads)
runJobsLoop:
	for {
		j := getJobToDo()
		if j != nil {
			clog.Debug(ctx, "Send job to do: "+j.Base().Id)
			jobs.ReqCH <- j
		} else if !jobsToDoExists(ctx) {
			break runJobsLoop
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	clog.Debug(ctx, "Stop sending jobs to do")
	<-allJobsResponsePresent
	jobs.StopAllRunners(ctx)
}

// Listens for a job responses
func listenJobsResponse(ctx common.DCtx, allDone chan bool) {
	for i := 1; i <= len(Jobs); i++ {
		j := <-jobs.RespCH
		clog.Debug(ctx, "Receive job response", j.Base().Id, j.Base().State)
		if j.Base().State == "ERROR" {
			cancelDepJobs(ctx, j.Base().Id)
		}
	}
	clog.Trace(ctx, "All jobs response present")
	allDone <- true
}

// Finds the job to do. It omit jobs waiting for other jobs
func getJobToDo() jobs.IJob {
	for _, j := range Jobs {
		if len(j.Base().State) == 0 {
			if len(j.Base().DependsOnId) == 0 {
				j.Base().State = jobs.JobStateRunning
				return j
			} else {
				dependsDone := true
				for _, dj := range Jobs {
					if common.StringSliceContains(j.Base().DependsOnId, dj.Base().Id) {
						if dj.Base().State != jobs.JobStateDone {
							dependsDone = false
							break
						}
					}
				}
				if dependsDone {
					j.Base().State = jobs.JobStateRunning
					return j
				}
			}
		}
	}
	return nil
}

// Checks if exists job to do
func jobsToDoExists(ctx common.DCtx) bool {
	for _, j := range Jobs {
		if len(j.Base().State) == 0 { // no status
			return true
		}
	}
	clog.Debug(ctx, "There's no job to do")
	return false
}

// Cancels all jobs depends on job with given id
func cancelDepJobs(ctx common.DCtx, id string) {
	for _, j := range Jobs {
		if common.StringSliceContains(j.Base().DependsOnId, id) {
			if len(j.Base().State) == 0 { // no status
				clog.Debug(ctx, "Cancel job", j.Base().Id, "cause job err,cancel:", id)
				j.Base().State = "CANCELED"
				jobs.RespCH <- j
				cancelDepJobs(ctx, j.Base().Id)
			}
		}

	}
}
