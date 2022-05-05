package ctmpl

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Compiles string go template
//
// On error - exit application
func CompileStringOrDie(ctx common.DCtx, tmpl string, params interface{}) string {
	if len(tmpl) == 0 || !strings.ContainsRune(tmpl, '{') {
		return tmpl
	}
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		clog.Error(ctx, err, "ERROR parse template string", tmpl)
		os.Exit(123)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, params)
	if err != nil {
		clog.Error(ctx, err, "ERROR execute template string", tmpl)
		clog.Error(ctx, nil, "Used params: ", fmt.Sprintf("%v", params))
		os.Exit(123)
	}
	out := buff.String()
	return out
}

// Compiles array of string go templates (in place)
//
// On error - exit application
func CompileStringArrOrDie(ctx common.DCtx, arr []string, params interface{}) {
	for idx, v := range arr {
		arr[idx] = CompileStringOrDie(ctx, v, params)
	}
}

// Compile text file with a go template content to destination file
//
// File name must starts with 'tmpl-' prefix
// Compiled file has the same name except prefix 'tmpl-'
func CompileTmplFile(ctx common.DCtx, tmplPath string, params interface{}) (string, error) {
	appendLog := &clog.AppendLog{}
	clog.Trace(ctx, "CompileTmplFile()", tmplPath)
	// parse template
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		clog.ErrorAndAppend(ctx, appendLog, err, "ERROR parsing template file", tmplPath)
		return appendLog.Txt, err
	}
	// create destination file (without 'tmpl-' prefix)
	tmplName := filepath.Base(tmplPath)
	tmplDir := filepath.Dir(tmplPath)
	filePath := filepath.Join(tmplDir, tmplName[5:])
	destinationFile, err := os.Create(filePath)
	if err != nil {
		clog.ErrorAndAppend(ctx, appendLog, err, "ERROR creating destination file for template", tmplPath, filePath)
		return appendLog.Txt, err
	}
	defer destinationFile.Close()

	// set script permissions
	// it is a dev environment - so there is no need to do it more complicated
	if strings.HasSuffix(filePath, ".sh") {
		err = destinationFile.Chmod(0777)
		if err != nil {
			clog.ErrorAndAppend(ctx, appendLog, err, "ERROR setting run permission (tmpl destination file)", filePath)
			return appendLog.Txt, err
		}
	}
	// execute template
	err = t.Execute(destinationFile, params)
	if err != nil {
		clog.ErrorAndAppend(ctx, appendLog, err, "ERROR execute template file", tmplPath)
		clog.ErrorAndAppend(ctx, appendLog, nil, "Used params: ", fmt.Sprintf("%v", params))
		return appendLog.Txt, err
	}
	return appendLog.Txt, nil
}

// Compiles all template files in a directory. Template file is a file which name starts with 'tmpl-' prefix
func CompileTmplFilesInDir(ctx common.DCtx, dirPath string, params interface{}) (string, error) {
	appendLog := &clog.AppendLog{}
	filePaths, err := getTmplFilePathsInDir(ctx, dirPath)
	if err != nil {
		clog.ErrorAndAppend(ctx, appendLog, err, "ERROR get template files paths in dir", dirPath)
		return appendLog.Txt, err
	}
	for _, tmplFilePath := range filePaths {
		txt, err := CompileTmplFile(ctx, tmplFilePath, params)
		appendLog.Txt += txt
		if err != nil {
			return appendLog.Txt, err
		}
	}
	return appendLog.Txt, nil
}

// Gets all template file names in a directory. Template file is a file which name starts with 'tmpl-' prefix
func getTmplFilePathsInDir(ctx common.DCtx, dirPath string) ([]string, error) {
	paths := []string{}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		clog.Error(ctx, err, "Error get tmpl file paths in dir", dirPath)
		return paths, err
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "tmpl-") {
			paths = append(paths, filepath.Join(dirPath, e.Name()))
		}
	}
	return paths, nil
}
