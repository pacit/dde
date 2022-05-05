package jobs

import (
	"fmt"
	"time"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/model"
)

const (
	JobTypeImage     = "Prepare image"
	JobTypeDeploy    = "Deploy service"
	JobTypeRemove    = "Remove service"
	JobTypeRunScript = "Script run"
	JobTypePrepare   = "Prepare env"
	JobTypeRmEnv     = "Remove env"

	JobStateRunning  = "RUNNING"
	JobStateDone     = "DONE"
	JobStateError    = "ERROR"
	JobStateCanceled = "CANCELED"

	JobNameDeploySrvCleanDir             = "deploy-srv-clean-wrk-dir"
	JobNameDeploySrvCreateDockerRes      = "deploy-srv-create-docker-res"
	JobNameDeploySrvCreateWrkFiles       = "deploy-srv-create-wrk-files"
	JobNameDeploySrvRun                  = "deploy-srv-run"
	JobNameDeploySrvWait                 = "deploy-srv-wait"
	JobNameExecuteScript                 = "execute-script"
	JobNamePrepEnvCleanWrkDir            = "prep-env-clean-wrk-dir"
	JobNamePrepEnvCreateDockerRes        = "prep-env-create-docker-res"
	JobNamePrepEnvCreateWrkFiles         = "prep-env-create-wrk-files"
	JobNamePrepImgBuild                  = "prep-img-build"
	JobNamePrepImgCleanWrkDir            = "prep-img-clean-wrk-dir"
	JobNamePrepImgCreateWrkFiles         = "prep-img-create-wrk-files"
	JobNamePrepImgExists                 = "prep-img-exists"
	JobNamePrepImgGitClone               = "prep-img-git-clone"
	JobNamePrepImgPull                   = "prep-img-pull"
	JobNamePrepImgRetag                  = "prep-img-retag"
	JobNameRmSrvRmDockerRes              = "rm-srv-rm-docker-res"
	JobNameRmSrvRmService                = "rm-srv-rm-service"
	JobNameRmSrvWait                     = "rm-srv-wait"
	JobNameRmEnvRmDockerRes              = "rm-env-rm-docker-res"
	JobNameRmEnvCleanup                  = "rm-env-cleanup"
	JobNameRunCustomScriptCleanDir       = "runs-clean-wrk-dir"
	JobNameRunCustomScriptCreateWrkFiles = "runs-create-wrk-files"
)

var (
	idSeq = 0
)

// Workspace info
type JobWorkspaceInfo struct {
	// Workspace root directory path
	Dir string
	// Workspace configuration
	Cfg model.WorkspaceJson
}

// Project info
type JobProjectInfo struct {
	// Project name
	Name string
	// Project version
	Version string
	// Project version (safe to use as a docker image tag)
	VersionDockerSafe string
	// Project version (safe to use as a dotnet assembly version)
	VersionDotnetSafe string
	// Project configuration
	Cfg model.ProjectJson
}

// Service info
type JobServiceInfo struct {
	// Environment name
	EnvName string
	// Service name
	SrvName string
	// Service configuration
	Cfg model.ServiceJson
}

// Environment info
type JobEnvironmentInfo struct {
	// Environment name
	Name string
	// Environment configuration
	Cfg model.EnvJson
}

// Job interface. All jobs must implements this.
type IJob interface {
	// Base job data
	Base() *BaseJob
	// Job short info
	ShortInfo() string
	// Job implementation
	DoIt(ctx common.DCtx) common.CmdResMulti
	// Job name
	JobName() string
	// Job project name
	GetProjName() string
	// Job project version
	GetProjVer() string
	// Job service name
	GetSrvName() string
	// Job environment name
	GetEnvName() string
	// Job description
	GetJobDescription() string
}

// Short job info
//
// Job.ShortInfo supplemented with job id and type
func ShortJobInfo(j IJob) string {
	out := fmt.Sprintf("Job-%-3s (%-14s)", j.Base().Id, j.Base().Type)
	out += j.ShortInfo()
	if len(j.Base().DependsOnId) > 0 {
		out += " {dependsOn:" + fmt.Sprintf("%v", j.Base().DependsOnId) + "}"
	}
	return out
}

// Job dependency
type JobDependency struct {
	// Dependant job's name
	JobName string
	// Dependant environment's name
	EnvName string
	// Dependant service's name
	SrvName string
	// Dependant project's name
	ProjName string
	// Dependant job's type
	JobType string
}

// Base job data
type BaseJob struct {
	// Job identifier
	Id string // seq number
	// Job type (one of constants: JobType*)
	Type string // JobType*
	// Dependant job's identifiers
	DependsOnId []string // Id
	// Dependant jobs
	//
	// It is mapped at the end - to the list: DependsOnId
	DependsOnJobs []JobDependency
	// Dependant services. Format: {env}:{srv}
	DependsOnSrv []string
	// Job state
	State string // "", RUNNING, End states: DONE|ERROR|CANCELED
	// Job result
	Result common.CmdResMulti
	// Job start time
	StartTime time.Time
	// Job time spent
	TimeDuration time.Duration
}

// Creates identifier for a new job
func GetNewId() string {
	idSeq += 1
	return fmt.Sprintf("J%v", idSeq)
}
