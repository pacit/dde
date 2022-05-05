package jobs

import (
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cbash"
	"github.com/pacit/dde/common/clog"
	"github.com/pacit/dde/docker"
	"github.com/pacit/dde/model/modelc"
)

// Creates docker's volumes
func createDockerResources_volumes(ctx common.DCtx, volumes []modelc.DockerResourceVolume, forceCheck bool) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(volumes) > 0 {
		if forceCheck {
			clog.Debug(ctx, "Force check resource available: volumes")
			results.Append(docker.LoadAvailableVolumes(ctx))
			if results.Err != nil {
				return results
			}
		}
		for _, v := range volumes {
			if !common.StringSliceContains(docker.Available.Volumes, v.Name) {
				cmd := "docker volume create"
				if len(v.Driver) > 0 {
					cmd += " --driver " + v.Driver
				}
				if len(v.Options) > 0 {
					for _, opt := range v.Options {
						cmd += " --opt " + opt
					}
				}
				cmd += " " + v.Name
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "volume "+v.Name+" - already exists")
			}
		}
	}
	return results
}

// Removes docker's volumes
func rmDockerResources_volumes(ctx common.DCtx, volumes []modelc.DockerResourceVolume) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(volumes) > 0 {
		for _, v := range volumes {
			if common.StringSliceContains(docker.Available.Volumes, v.Name) {
				cmd := "docker volume rm " + v.Name
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "volume "+v.Name+" - not exists (omit rm)")
			}
		}
	}
	return results
}

// Creates docker's networks
func createDockerResources_networks(ctx common.DCtx, networks []modelc.DockerResourceNetwork, forceCheck bool) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(networks) > 0 {
		if forceCheck {
			clog.Debug(ctx, "Force check resource available: networks")
			results.Append(docker.LoadAvailableNetworks(ctx))
			if results.Err != nil {
				return results
			}
		}
		for _, n := range networks {
			if !common.StringSliceContains(docker.Available.Networks, n.Name) {
				cmd := "docker network create"
				if len(n.Driver) > 0 {
					cmd += " --driver " + n.Driver
				}
				if len(n.Gateway) > 0 {
					cmd += " --gateway " + n.Gateway
				}
				if len(n.IpRange) > 0 {
					cmd += " --ip-range " + n.IpRange
				}
				if len(n.Subnet) > 0 {
					cmd += " --subnet " + n.Subnet
				}
				if n.Attachable {
					cmd += " --attachable"
				}
				cmd += " " + n.Name
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "network "+n.Name+" - already exists")
			}
		}
	}
	return results
}

// Removes docker's networks
func rmDockerResources_networks(ctx common.DCtx, networks []modelc.DockerResourceNetwork) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(networks) > 0 {
		for _, n := range networks {
			if common.StringSliceContains(docker.Available.Networks, n.Name) {
				cmd := "docker network rm " + n.Name
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "network "+n.Name+" - not exists (omit rm)")
			}
		}
	}
	return results
}

// Creates docker's secrets
func createDockerResources_secrets(ctx common.DCtx, secrets []modelc.DockerResourceSecret, wrk string, forceCheck bool) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(secrets) > 0 {
		if forceCheck {
			clog.Debug(ctx, "Force check resource available: secrets")
			results.Append(docker.LoadAvailableSecrets(ctx))
			if results.Err != nil {
				return results
			}
		}
		for _, s := range secrets {
			if !common.StringSliceContains(docker.Available.Secrets, s.Name) {
				cmd := "cd " + wrk + " && "
				if len(s.Value) > 0 {
					cmd += "printf \"" + s.Value + "\" | docker secret create"
				} else {
					cmd += "docker secret create"
				}
				if len(s.Driver) > 0 {
					cmd += " --driver " + s.Driver
				}
				cmd += " " + s.Name
				if len(s.File) > 0 {
					cmd += " " + s.File
				} else {
					cmd += " -"
				}
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "secret "+s.Name+" - already exists")
			}
		}
	}
	return results
}

// Removes docker's secrets
func rmDockerResources_secrets(ctx common.DCtx, secrets []modelc.DockerResourceSecret) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(secrets) > 0 {
		for _, s := range secrets {
			if common.StringSliceContains(docker.Available.Secrets, s.Name) {
				cmd := "docker secret rm " + s.Name
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "secret "+s.Name+" - not exists (omit rm)")
			}
		}
	}
	return results
}

// Creates docker's configs
func createDockerResources_configs(ctx common.DCtx, configs []modelc.DockerResourceConfig, wrk string, forceCheck bool) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(configs) > 0 {
		if forceCheck {
			clog.Debug(ctx, "Force check resource available: configs")
			results.Append(docker.LoadAvailableConfigs(ctx))
			if results.Err != nil {
				return results
			}
		}
		for _, c := range configs {
			if !common.StringSliceContains(docker.Available.Configs, c.Name) {
				cmd := "cd " + wrk + " && "
				if len(c.Value) > 0 {
					cmd += "printf \"" + c.Value + "\" | docker config create"
				} else {
					cmd += "docker config create"
				}
				cmd += " " + c.Name
				if len(c.File) > 0 {
					cmd += " " + c.File
				} else {
					cmd += " -"
				}
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "config "+c.Name+" - already exists")
			}
		}
	}
	return results
}

// Removes docker's configs
func rmDockerResources_configs(ctx common.DCtx, configs []modelc.DockerResourceConfig) common.CmdResMulti {
	results := common.CmdResMulti{}
	if len(configs) > 0 {
		for _, c := range configs {
			if common.StringSliceContains(docker.Available.Configs, c.Name) {
				cmd := "docker config rm " + c.Name
				results.Append(cbash.Call(ctx, cmd))
				if results.Err != nil {
					return results
				}
			} else {
				clog.Debug(ctx, "config "+c.Name+" - not exists (omit rm)")
			}
		}
	}
	return results
}
