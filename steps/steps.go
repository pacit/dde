package steps

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
)

// All jobs
var Jobs = []jobs.IJob{}

// Creates jobs of all required steps
func AddAllSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddAllSteps()")
	// run custom script
	if wrk.CliArgs.Command == "runs" {
		if err := AddCustomScriptSteps(ctx, wrk); err != nil {
			return err
		}
	}
	// prepare Envs
	if wrk.CliArgs.Command == "prepare" || wrk.CliArgs.ForcePrepareEnv {
		if err := AddPreapreEnvSteps(ctx, wrk); err != nil {
			return err
		}
	}
	// prepare images
	if common.StringSliceContains([]string{"build", "deploy"}, wrk.CliArgs.Command) {
		// force deploy is checking inside
		if err := AddBuildSteps(ctx, wrk); err != nil {
			return err
		}
	}
	// deploy/redeploy services
	if wrk.CliArgs.Command == "deploy" {
		if wrk.CliArgs.ForceRedeploy {
			if err := AddRemoveSrvSteps(ctx, wrk); err != nil {
				return err
			}
		}
		if err := AddDeploySrvSteps(ctx, wrk); err != nil {
			return err
		}
	}
	// remove services/environments
	if wrk.CliArgs.Command == "rm" {
		// remove services
		if err := AddRemoveSrvSteps(ctx, wrk); err != nil {
			return err
		}
		// remove env
		if wrk.CliArgs.ForceRmEnv {
			if err := AddRemoveEnvSteps(ctx, wrk); err != nil {
				return err
			}
		}
	}
	calculateJobDependencies()
	return nil
}

// Prints table with jobs to do
func PrintJobsTableToDo(ctx common.DCtx) {
	t := table.NewWriter()
	tableBytes := new(bytes.Buffer)
	t.SetOutputMirror(tableBytes)
	t.AppendHeader(table.Row{"#", "Type", "Work", "Env", "Service", "Project", "Version", "Depends"})
	for _, j := range Jobs {
		t.AppendRow([]interface{}{
			j.Base().Id,
			j.Base().Type,
			j.GetJobDescription(),
			j.GetEnvName(),
			j.GetSrvName(),
			j.GetProjName(),
			j.GetProjVer(),
			j.Base().DependsOnId,
		})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	clog.Info(ctx, "Jobs to do:\n"+tableBytes.String())
}

// Prints table with jobs state
func PrintJobsTableState(ctx common.DCtx) {
	t := table.NewWriter()
	tableBytes := new(bytes.Buffer)
	t.SetOutputMirror(tableBytes)
	t.AppendHeader(table.Row{"#", "Type", "Work", "Env", "Service", "Project", "Version", "State", "Time"})
	for _, j := range Jobs {
		timeStr := fmt.Sprintf("%v", j.Base().TimeDuration.Truncate(time.Millisecond))
		if !strings.HasSuffix(timeStr, "ms") {
			timeStr = timeStr + " "
		}
		t.AppendRow([]interface{}{
			j.Base().Id,
			j.Base().Type,
			j.GetJobDescription(),
			j.GetEnvName(),
			j.GetSrvName(),
			j.GetProjName(),
			j.GetProjVer(),
			j.Base().State,
			fmt.Sprintf("%10s", timeStr),
		})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	clog.Info(ctx, "Jobs state:\n"+tableBytes.String())
}

// Prints errors of all jobs with state "ERROR"
func PrintErrors(ctx common.DCtx) bool {
	isErr := false
	for _, j := range Jobs {
		if j.Base().State == "ERROR" {
			isErr = true
			clog.Error(ctx, nil, "Job error", jobs.ShortJobInfo(j))
			clog.Error(ctx, nil, "stdout", j.Base().Result.ErrRes.StdOut)
			clog.Error(ctx, nil, "stderr", j.Base().Result.ErrRes.StdErr)
			clog.Error(ctx, j.Base().Result.Err, "error", j.Base().Id)
		}
	}
	return isErr
}

// Creates array of identifiers. Omit empty values.
func depends(ids ...string) []string {
	outIds := []string{}
	for _, id := range ids {
		if len(id) > 0 {
			outIds = append(outIds, id)
		}
	}
	return outIds
}

// Find jobs identifiers by giver predicates
func getJobIds(dep jobs.JobDependency) []string {
	ids := []string{}
	for _, j := range Jobs {
		if len(dep.JobName) > 0 && j.JobName() != dep.JobName {
			continue
		}
		if len(dep.JobType) > 0 && j.Base().Type != dep.JobType {
			continue
		}
		if len(dep.EnvName) > 0 && j.GetEnvName() != dep.EnvName {
			continue
		}
		if len(dep.SrvName) > 0 && j.GetSrvName() != dep.SrvName {
			continue
		}
		if len(dep.ProjName) > 0 && j.GetProjName() != dep.ProjName {
			continue
		}
		ids = append(ids, j.Base().Id)
	}
	return ids
}

// Maps job dependencies to an array of job identifiers
func calculateJobDependencies() {
	// map DependsOnJobs to DependsOnId
	for _, j := range Jobs {
		for _, dep := range j.Base().DependsOnJobs {
			ids := getJobIds(dep)
			for _, id := range ids {
				if !common.StringSliceContains(j.Base().DependsOnId, id) {
					j.Base().DependsOnId = append(j.Base().DependsOnId, id)
				}
			}
		}
	}
	// map DependsOnSrv to DependsOnId
	for _, j := range Jobs {
		for _, dep := range j.Base().DependsOnSrv {
			split := strings.Split(dep, ":")
			ids := getJobIds(jobs.JobDependency{
				EnvName: split[0],
				SrvName: split[1],
				JobName: jobs.JobNameDeploySrvWait,
			})
			for _, id := range ids {
				if !common.StringSliceContains(j.Base().DependsOnId, id) {
					j.Base().DependsOnId = append(j.Base().DependsOnId, id)
				}
			}
		}
	}
}
