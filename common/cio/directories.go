package cio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Gets all file names in a directory
func GetFileNamesInDir(ctx common.DCtx, dirPath string) ([]string, error) {
	names := []string{}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		clog.Error(ctx, err, "Error reading dir", dirPath)
		return []string{}, err
	}
	for _, e := range entries {
		if e.Type().IsRegular() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// Gets all file names in a directory.
//
// On error - empty list is returned
func GetFileNamesInDirSilent(ctx common.DCtx, dirPath string) []string {
	names, err := GetFileNamesInDir(ctx, dirPath)
	if err != nil {
		clog.Warning(ctx, "Silent error while get file names in dir", dirPath, fmt.Sprintf("%v", err))
		return []string{}
	}
	return names
}

// Check if directory contains a file
func DirHasFileSilent(ctx common.DCtx, dirPath string, filename string) bool {
	files := GetFileNamesInDirSilent(ctx, dirPath)
	return common.StringSliceContains(files, filename)
}

// Gets directories names which contains a file with the given name
func GetDirNamesContainsFile(rootPath string, filename string) ([]string, error) {
	names := []string{}
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return names, err
	}
	for _, e := range entries {
		if e.IsDir() {
			if FileExists(filepath.Join(rootPath, e.Name(), filename)) {
				names = append(names, e.Name())
			}
		}
	}
	return names, nil
}

// Checks if file exists
func FileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		return true
	}
	return false
}
