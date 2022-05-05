package steps

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
)

// Creates jobs of remove service step
//
// For all required steps
func AddRemoveSrvSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddRemoveSrvSteps()")
	for _, env := range wrk.Environments {
		for _, srv := range env.Services {
			clog.Trace(ctx, "AddRemoveSrvSteps()", srv.Environment.Name, srv.Name)
			proj := wrk.GetProjectByName(ctx, srv.Cfg.Project)
			projInfo := proj.ToJobProjectInfo(wrk.VersionMap[proj.Name])
			srvInfo := srv.ToJobServiceInfo()
			if srv.NeedRmService(ctx) {
				Jobs = append(Jobs, &jobs.JobRmServiceRmService{
					BaseJ: jobs.BaseJob{
						Id:   jobs.GetNewId(),
						Type: jobs.JobTypeRemove,
					},
					Project: projInfo,
					Service: srvInfo,
				})
				Jobs = append(Jobs, &jobs.JobRmServiceWait{
					BaseJ: jobs.BaseJob{
						Id:   jobs.GetNewId(),
						Type: jobs.JobTypeRemove,
						DependsOnJobs: []jobs.JobDependency{
							{JobName: jobs.JobNameRmSrvRmService, EnvName: srv.Environment.Name, SrvName: srv.Name},
						},
					},
					Project: projInfo,
					Service: srvInfo,
				})
			}
			if srv.Cfg.DockerResources.HasResources() {
				Jobs = append(Jobs, &jobs.JobRmServiceRmDockerResources{
					BaseJ: jobs.BaseJob{
						Id:   jobs.GetNewId(),
						Type: jobs.JobTypeRemove,
						DependsOnJobs: []jobs.JobDependency{
							{JobName: jobs.JobNameRmSrvWait, EnvName: srv.Environment.Name, SrvName: srv.Name},
						},
					},
					Project:        projInfo,
					Service:        srvInfo,
					ForceRmVolumes: wrk.CliArgs.ForceRmVolumes,
				})
			}
		}
	}
	return nil
}
