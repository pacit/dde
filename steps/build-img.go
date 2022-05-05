package steps

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/engine"
	"github.com/pacit/dde/jobs"
)

// Creates jobs of build image step
//
// For all required projects
func AddBuildSteps(ctx common.DCtx, wrk *engine.Workspace) error {
	clog.Trace(ctx, "AddBuildSteps")
	for _, env := range wrk.Environments {
		for _, srv := range env.Services {
			proj := wrk.GetProjectByName(ctx, srv.Cfg.Project)
			if !buildImageJobAlreadyExists(ctx, proj.Name, wrk.VersionMap[proj.Name]) && srv.NeedBuildImage(ctx) {
				if err := addBuildStep(ctx, srv, wrk, proj); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Creates jobs of build image step - for one project
func addBuildStep(ctx common.DCtx, srv *engine.Service, wrk *engine.Workspace, proj *engine.Project) error {
	projInfo := proj.ToJobProjectInfo(wrk.VersionMap[proj.Name])
	// PrepareImageCleanDir
	Jobs = append(Jobs, &jobs.JobPrepareImageCleanDir{
		BaseJ: jobs.BaseJob{
			Id:   jobs.GetNewId(),
			Type: jobs.JobTypeImage,
		},
		Project: projInfo,
	})
	// PrepareImageCreateWorkFiles
	Jobs = append(Jobs, &jobs.JobPrepareImageCreateWorkFiles{
		BaseJ: jobs.BaseJob{
			Id: jobs.GetNewId(),
			DependsOnJobs: []jobs.JobDependency{
				{JobName: jobs.JobNamePrepImgCleanWrkDir, ProjName: projInfo.Name},
			},
			Type: jobs.JobTypeImage,
		},
		Project:    projInfo,
		Properties: proj.GetTemplateProps(ctx, true),
	})
	// if project has git repo
	if len(proj.Cfg.GitRepo) > 0 {
		// PrepareImageGitClone
		Jobs = append(Jobs, &jobs.JobPrepareImageGitClone{
			BaseJ: jobs.BaseJob{
				Id: jobs.GetNewId(),
				DependsOnJobs: []jobs.JobDependency{
					{JobName: jobs.JobNamePrepImgCleanWrkDir, ProjName: projInfo.Name},
				},
				Type: jobs.JobTypeImage,
			},
			Project: projInfo,
		})
	}
	// PrepareImageBuild
	Jobs = append(Jobs, &jobs.JobPrepareImageBuild{
		BaseJ: jobs.BaseJob{
			Id: jobs.GetNewId(),
			DependsOnJobs: []jobs.JobDependency{
				{JobName: jobs.JobNamePrepImgCreateWrkFiles, ProjName: projInfo.Name},
				{JobName: jobs.JobNamePrepImgGitClone, ProjName: projInfo.Name},
			},
			Type: jobs.JobTypeImage,
		},
		Project: projInfo,
	})
	return nil
}

// Checks if there is already exists a build job for a project
func buildImageJobAlreadyExists(ctx common.DCtx, projName string, projVer string) bool {
	for _, j := range Jobs {
		if j.JobName() == jobs.JobNamePrepImgBuild && j.GetProjName() == projName && j.GetProjVer() == projVer {
			return true
		}
	}
	return false
}
