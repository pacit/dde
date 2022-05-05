package cio

import (
	"fmt"

	"github.com/magiconair/properties"
	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Reads text file with a .properties format to a map
func ReadPropertiesFileAsMap(ctx common.DCtx, filePath string) (map[string]string, error) {
	props, err := properties.LoadFile(filePath, properties.UTF8)
	if err != nil {
		return map[string]string{}, err
	}
	return props.Map(), nil
}

// Reads text file with a .properties format to a map
//
// On error - empty map is returned
func ReadPropertiesFileAsMapSilent(ctx common.DCtx, filePath string) map[string]string {
	propsMap, err := ReadPropertiesFileAsMap(ctx, filePath)
	if err != nil {
		clog.Warning(ctx, "Silent error while reading properties file", filePath, fmt.Sprintf("%v", err))
		return map[string]string{}
	}
	return propsMap
}
