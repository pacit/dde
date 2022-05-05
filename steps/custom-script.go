package steps

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
	"github.com/pacit/dde/model"
)

// Creates jobs of run custom script step
//
// First it checks if it needed
func AddCustomScriptSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddCustomScriptSteps()")
	if len(wrk.CliArgs.OtherArgs) > 0 {
		scrName := wrk.CliArgs.OtherArgs[0]
		args := wrk.CliArgs.OtherArgs[1:]
		rs := wrk.GetRunScriptByName(scrName)
		if rs != nil {
			if len(args) > 0 {
				for _, arg := range args {
					for _, script := range rs.Scripts {
						script.Args = append(script.Args, arg)
					}
				}
			}
			addCustomScriptSteps(ctx, wrk, *rs)
		}
	}
	return nil
}

// Creates jobs of run custom script step
func addCustomScriptSteps(ctx common.DCtx, wrk *engine.Workspace, runScript model.RunScriptJson) error {
	// before run (update projects)
	if len(runScript.BeforeRun) > 0 {
		for _, beforeRun := range runScript.BeforeRun {
			if len(beforeRun.UpdateProject) > 0 {
				proj := wrk.GetProjectByName(ctx, beforeRun.UpdateProject)
				if proj != nil {
					projInfo := proj.ToJobProjectInfo(wrk.VersionMap[proj.Name])
					// PrepareImageCleanDir
					Jobs = append(Jobs, &jobs.JobPrepareImageCleanDir{
						BaseJ: jobs.BaseJob{
							Id:   jobs.GetNewId(),
							Type: jobs.JobTypeRunScript,
						},
						Project: projInfo,
					})
					// PrepareImageCreateWorkFiles
					Jobs = append(Jobs, &jobs.JobPrepareImageCreateWorkFiles{
						BaseJ: jobs.BaseJob{
							Id: jobs.GetNewId(),
							DependsOnJobs: []jobs.JobDependency{
								{JobName: jobs.JobNamePrepImgCleanWrkDir, ProjName: proj.Name},
							},
							Type: jobs.JobTypeRunScript,
						},
						Project:    projInfo,
						Properties: proj.GetTemplateProps(ctx, true),
					})
					// clone git repo
					Jobs = append(Jobs, &jobs.JobPrepareImageGitClone{
						BaseJ: jobs.BaseJob{
							Id:   jobs.GetNewId(),
							Type: jobs.JobTypeRunScript,
							DependsOnJobs: []jobs.JobDependency{
								{JobName: jobs.JobNamePrepImgCleanWrkDir, ProjName: proj.Name},
							},
						},
						Project: projInfo,
					})
				}
			}
		}
	}

	if len(wrk.Cfg.TemplateFileDirs) > 0 {
		Jobs = append(Jobs, &jobs.JobRunCustomScriptCleanDir{
			BaseJ: jobs.BaseJob{
				Id:   jobs.GetNewId(),
				Type: jobs.JobTypeRunScript,
			},
		})
		Jobs = append(Jobs, &jobs.JobRunCustomScriptCreateWrkFiles{
			BaseJ: jobs.BaseJob{
				Id:   jobs.GetNewId(),
				Type: jobs.JobTypeRunScript,
				DependsOnJobs: []jobs.JobDependency{
					{JobName: jobs.JobNameRunCustomScriptCleanDir},
				},
			},
			Workspace:  wrk.ToJobEnvironmentInfo(),
			Properties: wrk.GetTemplateProps(ctx, true),
		})
	}

	dependsOn := getJobIds(jobs.JobDependency{JobName: jobs.JobNamePrepImgGitClone})
	dependsOn = append(dependsOn, getJobIds(jobs.JobDependency{JobName: jobs.JobNameRunCustomScriptCreateWrkFiles})...)
	if len(runScript.Scripts) > 0 {
		for _, s := range runScript.Scripts {
			sId := jobs.GetNewId()
			Jobs = append(Jobs, &jobs.JobExecuteScript{
				BaseJ: jobs.BaseJob{
					Id:          sId,
					Type:        jobs.JobTypeRunScript,
					DependsOnId: dependsOn,
				},
				HostWorkingDir: cpath.WrkWorkingDir(),
				ScriptCfg:      s,
				Description:    "Run custom script " + runScript.Name,
			})
			dependsOn = []string{sId}
		}
	}

	return nil
}
