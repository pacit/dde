package modelc

import "github.com/pacit/dde/common"

// Docker resources declaration
type DockerResourcesJson struct {
	// Network declarations
	Networks []DockerResourceNetwork `json:"networks"`
	// Secret declarations
	Secrets []DockerResourceSecret `json:"secrets"`
	// Config declarations
	Configs []DockerResourceConfig `json:"configs"`
	// Volume declarations
	Volumes []DockerResourceVolume `json:"volumes"`
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerResourcesJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	for i, n := range dr.Networks {
		(&n).CompileTemplatesOrDie(ctx, props)
		dr.Networks[i] = n
	}
	for i, s := range dr.Secrets {
		(&s).CompileTemplatesOrDie(ctx, props)
		dr.Secrets[i] = s
	}
	for i, c := range dr.Configs {
		(&c).CompileTemplatesOrDie(ctx, props)
		dr.Configs[i] = c
	}
	for i, v := range dr.Volumes {
		(&v).CompileTemplatesOrDie(ctx, props)
		dr.Volumes[i] = v
	}
}

// Checking any resource declaration exists
//
// Returns false when no resource declaration found
func (dr DockerResourcesJson) HasResources() bool {
	return len(dr.Configs) > 0 ||
		len(dr.Networks) > 0 ||
		len(dr.Secrets) > 0 ||
		len(dr.Volumes) > 0
}
