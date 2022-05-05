package steps

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
)

// Creates jobs of prepare environment step
//
// For all required environment
func AddPreapreEnvSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddPreapreEnvSteps()")
	for _, env := range wrk.Environments {
		clog.Trace(ctx, "AddPreapreEnvSteps()", env.Name)
		envInfo := env.ToJobEnvironmentInfo()
		// PrepareEnvCleanDir
		Jobs = append(Jobs, &jobs.JobPrepareEnvCleanDir{
			BaseJ: jobs.BaseJob{
				Id:   jobs.GetNewId(),
				Type: jobs.JobTypePrepare,
			},
			Environment: envInfo,
		})
		// PrepareEnvCreateWorkFiles
		Jobs = append(Jobs, &jobs.JobPrepareEnvCreateWorkFiles{
			BaseJ: jobs.BaseJob{
				Id: jobs.GetNewId(),
				DependsOnJobs: []jobs.JobDependency{
					{JobName: jobs.JobNamePrepEnvCleanWrkDir, EnvName: env.Name},
				},
				Type: jobs.JobTypePrepare,
			},
			Environment: envInfo,
			Properties:  env.GetTemplateProps(ctx, true),
		})
		// create docker resources
		if env.Cfg.DockerResources.HasResources() {
			Jobs = append(Jobs, &jobs.JobPrepareEnvCreateDockerResources{
				BaseJ: jobs.BaseJob{
					Id: jobs.GetNewId(),
					DependsOnJobs: []jobs.JobDependency{
						{JobName: jobs.JobNamePrepEnvCreateWrkFiles, EnvName: env.Name},
					},
					Type: jobs.JobTypePrepare,
				},
				Environment: envInfo,
			})
		}

		if len(env.Cfg.Scripts.Prepare) > 0 {
			// first script dependency
			dependsOn := depends(
				getJobIds(jobs.JobDependency{JobName: jobs.JobNamePrepEnvCreateWrkFiles, EnvName: env.Name})...,
			)
			dependsOn = append(
				dependsOn,
				getJobIds(jobs.JobDependency{JobName: jobs.JobNamePrepEnvCreateDockerRes, EnvName: env.Name})...,
			)
			for i, script := range env.Cfg.Scripts.Prepare {
				jId := jobs.GetNewId()
				Jobs = append(Jobs, &jobs.JobExecuteScript{
					BaseJ: jobs.BaseJob{
						Id:          jId,
						DependsOnId: dependsOn,
						Type:        jobs.JobTypePrepare,
					},
					Environment:    &envInfo,
					ScriptCfg:      script,
					Description:    "Exe prepare script " + fmt.Sprintf("%v", i+1),
					HostWorkingDir: cpath.EnvWorkingDir(env.Name),
				})
				dependsOn = []string{jId}
			}
		}
	}
	return nil
}
