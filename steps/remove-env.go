package steps

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
)

// Creates jobs of remove environment step
//
// For all required environments
func AddRemoveEnvSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddRemoveEnvSteps()")
	for _, env := range wrk.Environments {
		clog.Trace(ctx, "AddRemoveEnvSteps()", env.Name)
		envInfo := env.ToJobEnvironmentInfo()
		allEnvJobs := getJobIds(jobs.JobDependency{EnvName: env.Name})
		if len(env.Cfg.Scripts.Cleanup) > 0 {
			// recreate working dir - and execute script templates
			Jobs = append(Jobs, &jobs.JobPrepareEnvCleanDir{
				BaseJ: jobs.BaseJob{
					Id:   jobs.GetNewId(),
					Type: jobs.JobTypeRmEnv,
				},
				Environment: envInfo,
			})
			Jobs = append(Jobs, &jobs.JobPrepareEnvCreateWorkFiles{
				BaseJ: jobs.BaseJob{
					Id: jobs.GetNewId(),
					DependsOnJobs: []jobs.JobDependency{
						{JobName: jobs.JobNamePrepEnvCleanWrkDir, EnvName: env.Name},
					},
					Type: jobs.JobTypeRmEnv,
				},
				Environment: envInfo,
				Properties:  env.GetTemplateProps(ctx, true),
			})
			dependsOn := []string{}
			dependsOn = append(dependsOn, allEnvJobs...)
			dependsOn = append(dependsOn, getJobIds(jobs.JobDependency{JobName: jobs.JobNamePrepEnvCreateWrkFiles, EnvName: env.Name})...)

			for i, script := range env.Cfg.Scripts.Cleanup {
				jId := jobs.GetNewId()
				Jobs = append(Jobs, &jobs.JobExecuteScript{
					BaseJ: jobs.BaseJob{
						Id:          jId,
						DependsOnId: dependsOn,
						Type:        jobs.JobTypeRmEnv,
					},
					Environment:    &envInfo,
					ScriptCfg:      script,
					Description:    "Exe cleanup script " + fmt.Sprintf("%v", i+1),
					HostWorkingDir: cpath.EnvWorkingDir(env.Name),
				})
				dependsOn = []string{jId}
			}
		}

		// remove docker resources
		if env.Cfg.DockerResources.HasResources() {
			Jobs = append(Jobs, &jobs.JobRmEnvRmDockerResources{
				BaseJ: jobs.BaseJob{
					Id: jobs.GetNewId(),
					DependsOnJobs: []jobs.JobDependency{
						{EnvName: env.Name},
					},
					Type: jobs.JobTypeRmEnv,
				},
				Environment:    envInfo,
				ForceRmVolumes: wrk.CliArgs.ForceRmVolumes,
			})
		}
	}
	return nil
}
