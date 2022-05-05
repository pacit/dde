package cio

import (
	"fmt"
	"io/ioutil"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Reads text file to a string
func ReadTextFile(ctx common.DCtx, filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		clog.Error(ctx, err, "Error reading text file", filePath)
		return "", err
	}
	return string(bytes), nil
}

// Reads text file to a string
//
// On error - empty string is returned
func ReadTextFileSilent(ctx common.DCtx, filePath string) string {
	txt, err := ReadTextFile(ctx, filePath)
	if err != nil {
		clog.Warning(ctx, "Silent error while reading text file", filePath, fmt.Sprintf("%v", err))
	}
	return txt
}
