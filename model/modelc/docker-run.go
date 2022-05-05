package modelc

import (
	"fmt"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Docker service create configuration
type DockerRunJson struct {
	// Service name
	Name string `json:"name"`
	// Hostname
	Hostname string `json:"hostname"`
	// Number of replicas
	Replicas int `json:"replicas"`
	// Service secrets
	Secrets []DockerRunSecret `json:"secrets"`
	// Service configs
	Configs []DockerRunConfig `json:"configs"`
	// Service networks
	Networks []string `json:"networks"`
	// Service mounts
	Mounts []DockerRunMount `json:"mounts"`
	// Service published ports
	PublishPorts []DockerRunPublish `json:"publishPorts"`
	// Service Environment variables
	Envs []string `json:"envs"`
}

// Creates docker command to run service (docker service create [params])
func (r DockerRunJson) GetCommand(projName string, version string) string {
	cmd := "docker service create"
	cmd += " --name \"" + r.Name + "\""
	cmd += " --hostname \"" + r.Hostname + "\""
	cmd += " --replicas \"" + fmt.Sprint(r.Replicas) + "\""
	for _, n := range r.Networks {
		cmd += " --network \"" + n + "\""
	}
	for _, m := range r.Mounts {
		cmd += " --mount \"type=" + m.Type + ",source=" + m.Source + ",destination=" + m.Destination + "\""
	}
	for _, s := range r.Secrets {
		cmd += " --secret \"source=" + s.Source
		if len(s.Target) > 0 {
			cmd += ",target=" + s.Target
		}
		cmd += "\""
	}
	for _, c := range r.Configs {
		cmd += " --config \"source=" + c.Source + ",target=" + c.Target + "\""
	}
	for _, p := range r.PublishPorts {
		cmd += " --publish published=" + p.Published + ",target=" + p.Target
		if len(p.Protocol) > 0 {
			cmd += ",protocol=" + p.Protocol
		}
		if len(p.Mode) > 0 {
			cmd += ",mode=" + p.Mode
		}
	}
	for _, n := range r.Envs {
		cmd += " --env \"" + n + "\""
	}
	cmd += " " + projName + ":" + version
	return cmd
}

// It compiles values which are go templates (replaces placeholders with values from properties)
func (dr *DockerRunJson) CompileTemplatesOrDie(ctx common.DCtx, props map[string]string) {
	dr.Name = ctmpl.CompileStringOrDie(ctx, dr.Name, props)
	dr.Hostname = ctmpl.CompileStringOrDie(ctx, dr.Hostname, props)
	for i, s := range dr.Secrets {
		(&s).CompileTemplatesOrDie(ctx, props)
		dr.Secrets[i] = s
	}
	for i, c := range dr.Configs {
		(&c).CompileTemplatesOrDie(ctx, props)
		dr.Configs[i] = c
	}
	for i, m := range dr.Mounts {
		(&m).CompileTemplatesOrDie(ctx, props)
		dr.Mounts[i] = m
	}
	for i, p := range dr.PublishPorts {
		(&p).CompileTemplatesOrDie(ctx, props)
		dr.PublishPorts[i] = p
	}

	ctmpl.CompileStringArrOrDie(ctx, dr.Networks, props)
	ctmpl.CompileStringArrOrDie(ctx, dr.Envs, props)
}
