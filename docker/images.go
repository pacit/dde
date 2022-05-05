package docker

import (
	"fmt"
	"strings"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/model"
)

// Docker image
type DockerImage struct {
	// Image name
	Name string
	// Image version
	Version string
}

// Loads available docker images using `docker image ls` command
func LoadAvailableImages(ctx common.DCtx) common.CmdRes {
	clog.Trace(ctx, "LoadAvailableImages()")
	images := []DockerImage{}
	res := cbash.Call(ctx, "docker image ls --format \"{{.Repository}} {{.Tag}}\" | grep -v \"<none>\"")
	if res.Err != nil {
		return res
	}
	for _, line := range strings.Split(res.StdOut, "\n") {
		if len(line) > 0 {
			split := strings.Split(line, " ")
			images = append(images, DockerImage{
				Name:    split[0],
				Version: split[1],
			})
		}
	}
	Available.Images = images
	return res
}

// Checks availability of docker image in a remote repository
func CheckImageExistsInRepo(ctx common.DCtx, repo model.WorkspaceJsonDockerRepo, name string, version string) bool {
	nameWithoutRepoAddr := strings.Replace(name, repo.Address, "", 1)
	command := fmt.Sprintf("curl --silent -I -u \"%s:%s\" %s://%s/v2%s/manifests/%s | grep \"200 OK\"",
		repo.Username, repo.Password, repo.Protocol, repo.Address, nameWithoutRepoAddr, version,
	)
	res := cbash.Call(ctx, command)
	if res.Err != nil {
		return false
	}
	return len(res.StdOut) > 5
}
