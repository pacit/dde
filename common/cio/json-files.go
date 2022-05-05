package cio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Reads text file with a JSON content to provided interface
func ReadJsonFile(ctx common.DCtx, filePath string, jsonObj interface{}) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, jsonObj)
}

// Reads text file with a JSON content to provided interface
//
// On error - do nothing. Provided object won't be filled
func ReadJsonFileSilent(ctx common.DCtx, filePath string, jsonObj interface{}) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		clog.Warning(ctx, "Silent error while reading json file", filePath, fmt.Sprintf("%v", err))
		return
	}
	if len(bytes) == 0 {
		return
	}
	err = json.Unmarshal(bytes, jsonObj)
	if err != nil {
		clog.Warning(ctx, "Silent error while unmarshal json file", filePath, fmt.Sprintf("%v", err))
	}
}

// Reads text file with a JSON content to provided map
func ReadJsonFileAsMap(ctx common.DCtx, filePath string) (map[string]string, error) {
	outMap := make(map[string]string)
	err := ReadJsonFile(ctx, filePath, &outMap)
	if err != nil {
		clog.Error(ctx, err, "Error reading json file as a map", filePath)
		return map[string]string{}, err
	}
	return outMap, nil
}

// Reads text file with a JSON content to provided map
//
// On error - empty map will be returned
func ReadJsonFileAsMapSilent(ctx common.DCtx, filePath string) map[string]string {
	outMap := make(map[string]string)
	ReadJsonFileSilent(ctx, filePath, &outMap)
	return outMap
}
