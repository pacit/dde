package engine

import (
	"regexp"
	"strings"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/ctmpl"
)

// Converts version text to a safe text which can be used as a docker's image version
func VersionDockerSafe(v string) string {
	out := strings.ReplaceAll(v, "/", "-")
	out = strings.ReplaceAll(out, " ", "")
	return out
}

// Converts version text to a safe text which can be used as a dotnet's assembly version
func VersionDotnetSafe(v string) string {
	out := strings.ReplaceAll(v, "/", "")
	out = strings.ReplaceAll(out, " ", "")
	validVer := regexp.MustCompile(`^[0-9]+(\.[0-9]+)*$`)
	if validVer.MatchString(out) {
		return v
	} else {
		return "0.0.0"
	}
}

// Prepares properties and compile this which are templates
func prepareAndCompileTemplateProperties(ctx common.DCtx, parent map[string]string, own map[string]string) map[string]string {
	props := make(map[string]string)
	for k, v := range parent {
		props[k] = v
	}
	for k, v := range own {
		props[k] = v
	}
	for i := 0; i < 3; i++ {
		for k, v := range props {
			props[k] = ctmpl.CompileStringOrDie(ctx, v, props)
		}
	}
	return props
}

// Combines properties. Gets parent properties and overwrites with own properties
func prepareTemplateProperties(ctx common.DCtx, parent map[string]string, own map[string]string) map[string]string {
	props := make(map[string]string)
	for k, v := range parent {
		props[k] = v
	}
	for k, v := range own {
		props[k] = v
	}
	return props
}
