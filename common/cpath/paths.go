package cpath

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/cio"
	"github.com/pacit/dde/common/clog"
)

// Workspace root dir path
var wrkRootPath = ""

// Inits common paths package
//
// Sets workspace root dir path
func init() {
	ctx := common.DCtx{}
	searchFileName := "workspace.json"
	wd, err := os.Getwd()
	if err != nil {
		clog.Error(ctx, err, "Cannot get working dir path")
	}
	root, err := findRootDir(ctx, wd, searchFileName)
	if err != nil {
		clog.Error(ctx, err, "Cannot find wrk root dir")
	}
	wrkRootPath = root
	clog.Info(ctx, "Wrk root path", wrkRootPath)
}

// Search for a directory which contains file. Search in parent directories
func findRootDir(ctx common.DCtx, dir string, searchFileName string) (string, error) {
	if cio.DirHasFileSilent(ctx, dir, searchFileName) {
		return dir, nil
	} else {
		parent := filepath.Dir(dir)
		if len(parent) > 1 {
			return findRootDir(ctx, parent, searchFileName)
		} else {
			return "", errors.New("root wrk dir not found in a parent tree")
		}
	}
}

// Gets absolute path
//
// If given path is absolute - returns it. If given path is relative - returns full path where as root, workspace root dir will be used
func GetAbsPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Clean(filepath.Join(wrkRootPath, p))
}

// Workspace definition directory path
func WrkDefinitionDir() string {
	return wrkRootPath
}

// Workspace definition - environments directory path
func WrkDefinitionEnvsDir() string {
	return filepath.Join(wrkRootPath, "envs")
}

// Workspace definition - projects directory path
func WrkDefinitionProjectsDir() string {
	return filepath.Join(wrkRootPath, "projects")
}

// Workspace definition - version files directory path
func WrkVersionsDir() string {
	return filepath.Join(wrkRootPath, "versions")
}

// Workspace working directory path
func WrkWorkingDir() string {
	return filepath.Join(wrkRootPath, "nogit", "ww")
}

// Service definition directory path
func SrvDefinitionDir(envName string, srvName string) string {
	return filepath.Join(wrkRootPath, "envs", envName, srvName)
}

// Service working directory path
func SrvWorkingDir(envName string, srvName string) string {
	return filepath.Join(wrkRootPath, "nogit", "ws", envName, srvName)
}

// Environment definition directory path
func EnvDefinitionDir(envName string) string {
	return filepath.Join(wrkRootPath, "envs", envName)
}

// Environment working directory path
func EnvWorkingDir(envName string) string {
	return filepath.Join(wrkRootPath, "nogit", "we", envName)
}

// Project definition directory path
func ProjDefinitionDir(projName string) string {
	return filepath.Join(wrkRootPath, "projects", projName)
}

// Project working directory path
func ProjWorkingDir(projName string) string {
	return filepath.Join(wrkRootPath, "nogit", "wp", projName)
}

// Project source path (from git repo) in project working directory
func ProjWorkingDirSrc(projName string) string {
	return filepath.Join(wrkRootPath, "nogit", "wp", projName, "src")
}
