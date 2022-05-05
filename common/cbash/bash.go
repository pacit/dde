package cbash

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/pacit/dde/common"
	"github.com/pacit/dde/common/clog"
)

// Call bash command
func Call(ctx common.DCtx, cmd string) common.CmdRes {
	clog.Debug(ctx, "cbash.Call", cmd)
	command := exec.Command("bash", "-c", cmd)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	stdoutTraceProgress := &customWriter{Ctx: ctx, TypePrefix: "stdout"}
	stderrTraceProgress := &customWriter{Ctx: ctx, TypePrefix: "stderr"}
	command.Stdout = io.MultiWriter(stdout, stdoutTraceProgress)
	command.Stderr = io.MultiWriter(stderr, stderrTraceProgress)
	if err := command.Start(); err != nil {
		clog.Error(ctx, err, "Error calling bash command")
		return common.CmdRes{
			StdOut: stdout.String(),
			StdErr: stderr.String(),
			Err:    err,
		}
	}

	if err := command.Wait(); err != nil {
		clog.Error(ctx, err, "Error waiting bash command")
		return common.CmdRes{
			StdOut: stdout.String(),
			StdErr: stderr.String(),
			Err:    err,
		}
	}
	return common.CmdRes{
		StdOut: stdout.String(),
		StdErr: stderr.String(),
	}
}

// Custom writer
//
// Prints Trace logs
type customWriter struct {
	TypePrefix string
	Ctx        common.DCtx
	bufferStr  string
}

// Prints Trace logs
func (p *customWriter) Write(b []byte) (n int, err error) {
	p.bufferStr += string(b)
	lines := strings.Split(p.bufferStr, "\n")
	for _, line := range lines {
		clog.Trace(p.Ctx, p.TypePrefix, line)
	}
	return len(b), nil
}
