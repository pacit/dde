package model

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker image reference
type ProjectJsonDockerImage struct {
	// repository name from wrk config
	RepoId string `json:"repoId"`
	// image name
	ImageName string `json:"imageName"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (di *ProjectJsonDockerImage) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	di.RepoId = ctmpl.CompileStringOrDie(ctx, di.RepoId, props)
	di.ImageName = ctmpl.CompileStringOrDie(ctx, di.ImageName, props)
}
