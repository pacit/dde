package docker

import (
	"strconv"
	"strings"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/clog"
)

// Docker service
type DockerService struct {
	// Service name
	ServiceName string
	// Docker image name
	ImageName string
	// Image version
	ImageVersion string
	// Number of replicas that runs
	ReplicasRuns int
	// Number of replicas that should run
	ReplicasAll int
}

// Loads running docker services using `docker service ls` command
func LoadAvailableServices(ctx common.DCtx) common.CmdRes {
	clog.Trace(ctx, "LoadAvailableServices()")
	services := []DockerService{}
	res := cbash.Call(ctx, "docker service ls --format \"{{.Name}} {{.Replicas}} {{.Image}}\"")
	if res.Err != nil {
		return res
	}
	for _, line := range strings.Split(res.StdOut, "\n") {
		if len(line) > 0 {
			split := strings.Split(line, " ")
			replicasSplit := strings.Split(split[1], "/")
			replRuns, err := strconv.Atoi(replicasSplit[0])
			if err != nil {
				res.Err = err
				return res
			}
			replAll, err := strconv.Atoi(replicasSplit[1])
			if err != nil {
				res.Err = err
				return res
			}
			imageSplit := strings.Split(split[2], ":")
			services = append(services, DockerService{
				ServiceName:  split[0],
				ImageName:    imageSplit[0],
				ImageVersion: imageSplit[1],
				ReplicasRuns: replRuns,
				ReplicasAll:  replAll,
			})
		}
	}
	Available.Services = services
	return res
}
