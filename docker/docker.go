package docker

import (
	"strings"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/clog"
)

// Available docker resources
type DockerAvailableResources struct {
	// Available config names
	Configs []string
	// Available docker images (name, version)
	Images []DockerImage
	// Available docker network names
	Networks []string
	// Available docker secret names
	Secrets []string
	// Available docker volume names
	Volumes []string
	// Running docker services
	Services []DockerService
}

// Available docker resources
var Available = &DockerAvailableResources{}

// Prints error logs
func printLoadError(ctx common.DCtx, res common.CmdRes, msg string) {
	clog.Error(ctx, nil, res.StdOut)
	clog.Error(ctx, nil, res.StdErr)
	clog.Error(ctx, res.Err, "Error loading available "+msg)
}

// Loads all available docker resources
func LoadAllAvailableResources(ctx common.DCtx) error {
	clog.Trace(ctx, "LoadAllAvailableResources()")
	res := LoadAvailableConfigs(ctx)
	if res.Err != nil {
		printLoadError(ctx, res, "configs")
		return res.Err
	}
	res = LoadAvailableSecrets(ctx)
	if res.Err != nil {
		printLoadError(ctx, res, "secrets")
		return res.Err
	}
	res = LoadAvailableNetworks(ctx)
	if res.Err != nil {
		printLoadError(ctx, res, "networks")
		return res.Err
	}
	res = LoadAvailableVolumes(ctx)
	if res.Err != nil {
		printLoadError(ctx, res, "volumes")
		return res.Err
	}
	res = LoadAvailableServices(ctx)
	if res.Err != nil {
		printLoadError(ctx, res, "services")
		return res.Err
	}
	res = LoadAvailableImages(ctx)
	if res.Err != nil {
		printLoadError(ctx, res, "images")
		return res.Err
	}
	return nil
}

// Checks availability of docker image
func (r *DockerAvailableResources) HasImage(name string, version string) bool {
	for _, ai := range Available.Images {
		if ai.Name == name && ai.Version == version {
			return true
		}
	}
	return false
}

// Checks availability of docker service
func (r *DockerAvailableResources) HasServiceName(name string) bool {
	for _, s := range Available.Services {
		if s.ServiceName == name {
			return true
		}
	}
	return false
}

// Checks availability of docker image in a provided version
func (r *DockerAvailableResources) HasServiceInVersion(name string, version string) bool {
	for _, s := range Available.Services {
		if s.ServiceName == name && s.ImageVersion == version {
			return true
		}
	}
	return false
}

// Runs bash command and convert output to a string array of lines
func loadAvailableStringArr(ctx common.DCtx, cmd string) (common.CmdRes, []string) {
	arr := []string{}
	res := cbash.Call(ctx, cmd)
	if res.Err != nil {
		return res, arr
	}
	for _, line := range strings.Split(res.StdOut, "\n") {
		if len(line) > 0 {
			arr = append(arr, strings.Trim(line, " "))
		}
	}
	return res, arr
}
