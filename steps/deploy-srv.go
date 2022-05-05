package steps

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/common/cpath"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
)

// Creates jobs of deploy service step
//
// For all required services
func AddDeploySrvSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddDeploySrvSteps()")
	for _, env := range wrk.Environments {
		for _, srv := range env.Services {
			clog.Trace(ctx, "AddDeploySrvSteps()", srv.Environment.Name, srv.Name)
			proj := wrk.GetProjectByName(ctx, srv.Cfg.Project)
			projInfo := proj.ToJobProjectInfo(wrk.VersionMap[proj.Name])
			srvInfo := srv.ToJobServiceInfo()

			if !srv.NeedDeploy() {
				clog.Trace(ctx, "Service is already running:", env.Name, srv.Name, "img:", projInfo.Name, projInfo.VersionDockerSafe)
				continue
			}

			// clean wrk dir
			Jobs = append(Jobs, &jobs.JobDeployServiceCleanDir{
				BaseJ: jobs.BaseJob{
					Id:   jobs.GetNewId(),
					Type: jobs.JobTypeDeploy,
				},
				Project: projInfo,
				Service: srvInfo,
			})

			// create wrk files
			Jobs = append(Jobs, &jobs.JobDeployServiceCreateWorkFiles{
				BaseJ: jobs.BaseJob{
					Id: jobs.GetNewId(),
					DependsOnJobs: []jobs.JobDependency{
						{JobName: jobs.JobNameDeploySrvCleanDir, EnvName: env.Name, SrvName: srv.Name},
					},
					Type: jobs.JobTypeDeploy,
				},
				Project:    projInfo,
				Service:    srvInfo,
				Properties: srv.GetTemplateProps(ctx, true),
			})

			// creating docker resources (if any defined)
			if srv.Cfg.DockerResources.HasResources() {
				Jobs = append(Jobs, &jobs.JobDeployServiceCreateDockerResources{
					BaseJ: jobs.BaseJob{
						Id: jobs.GetNewId(),
						DependsOnJobs: []jobs.JobDependency{
							{JobName: jobs.JobNameDeploySrvCreateWrkFiles, EnvName: env.Name, SrvName: srv.Name},
							{JobName: jobs.JobNameRmSrvRmDockerRes, EnvName: env.Name, SrvName: srv.Name},
						},
						Type: jobs.JobTypeDeploy,
					},
					Project:                  projInfo,
					Service:                  srvInfo,
					ForceCheckResourceExists: wrk.CliArgs.ForceRedeploy,
				})
			}

			dependsOn := depends(
				getJobIds(jobs.JobDependency{JobName: jobs.JobNameDeploySrvCreateWrkFiles, EnvName: env.Name, SrvName: srv.Name})...,
			)
			dependsOn = append(
				dependsOn,
				getJobIds(
					jobs.JobDependency{JobName: jobs.JobNameDeploySrvCreateDockerRes, EnvName: env.Name, SrvName: srv.Name},
				)...,
			)
			// before scripts
			if len(srv.Cfg.Scripts.BeforeRun) > 0 {
				for i, script := range srv.Cfg.Scripts.BeforeRun {
					jId := jobs.GetNewId()
					Jobs = append(Jobs, &jobs.JobExecuteScript{
						BaseJ: jobs.BaseJob{
							Id:          jId,
							DependsOnId: dependsOn,
							Type:        jobs.JobTypeDeploy,
						},
						Project:        &projInfo,
						Service:        &srvInfo,
						ScriptCfg:      script,
						Description:    "Exe before run script " + fmt.Sprintf("%v", i+1),
						HostWorkingDir: cpath.SrvWorkingDir(env.Name, srv.Name),
					})
					dependsOn = []string{jId}
				}
			}

			// Run service
			Jobs = append(Jobs, &jobs.JobDeployServiceRun{
				BaseJ: jobs.BaseJob{
					Id:           jobs.GetNewId(),
					DependsOnId:  dependsOn,
					DependsOnSrv: srv.Cfg.DependsOnSrv,
					DependsOnJobs: []jobs.JobDependency{
						{JobName: jobs.JobNamePrepImgBuild, ProjName: proj.Name},
						{JobName: jobs.JobNamePrepImgPull, ProjName: proj.Name},
						{JobName: jobs.JobNamePrepImgRetag, ProjName: proj.Name},
						{JobName: jobs.JobNameRmSrvWait, EnvName: env.Name, SrvName: srv.Name},
						{JobName: jobs.JobNameRmSrvRmDockerRes, EnvName: env.Name, SrvName: srv.Name},
						{JobName: jobs.JobNamePrepEnvCreateDockerRes, EnvName: env.Name},
						{JobName: jobs.JobNameExecuteScript, JobType: jobs.JobTypePrepare, EnvName: env.Name},
					},
					Type: jobs.JobTypeDeploy,
				},
				Project:    projInfo,
				Service:    srvInfo,
				Properties: srv.GetTemplateProps(ctx, true),
			})

			// wait for service
			waitId := jobs.GetNewId()
			timeout := 60
			if srv.Cfg.WaitForServiceTimeoutS > 0 {
				timeout = srv.Cfg.WaitForServiceTimeoutS
			}
			Jobs = append(Jobs, &jobs.JobDeployServiceWait{
				BaseJ: jobs.BaseJob{
					Id: waitId,
					DependsOnJobs: []jobs.JobDependency{
						{JobName: jobs.JobNameDeploySrvRun, EnvName: env.Name, SrvName: srv.Name},
					},
					Type: jobs.JobTypeDeploy,
				},
				Project: projInfo,
				Service: srvInfo,
				Timeout: timeout,
			})

			dependsOn = depends(
				getJobIds(jobs.JobDependency{JobName: jobs.JobNameDeploySrvWait, EnvName: env.Name, SrvName: srv.Name})...,
			)
			// after scripts
			if len(srv.Cfg.Scripts.AfterRun) > 0 {
				for i, script := range srv.Cfg.Scripts.AfterRun {
					jId := jobs.GetNewId()
					Jobs = append(Jobs, &jobs.JobExecuteScript{
						BaseJ: jobs.BaseJob{
							Id:          jId,
							DependsOnId: dependsOn,
							Type:        jobs.JobTypeDeploy,
						},
						Project:        &projInfo,
						Service:        &srvInfo,
						ScriptCfg:      script,
						Description:    "Exe after run script " + fmt.Sprintf("%v", i+1),
						HostWorkingDir: cpath.SrvWorkingDir(env.Name, srv.Name),
					})
					dependsOn = []string{jId}
				}
			}
		}
	}
	return nil
}
