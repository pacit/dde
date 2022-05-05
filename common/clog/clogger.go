package clog

import (
	"fmt"
	"os"
	"strings"

	"github.com/pacit/dde/common"
)

var (
	printTrace   = false
	printDebug   = false
	printInfo    = false
	printWarning = false
	PrintBaner   = true
)

const (
	ClogRED     = "\033[0;31m"
	ClogGREEN   = "\033[0;32m"
	ClogYELLOW  = "\033[0;33m"
	ClogBLUE    = "\033[0;34m"
	ClogMAGENTA = "\033[0;35m"
	ClogCYAN    = "\033[0;36m"
	ClogDEFAULT = "\033[0;39m"
	ClogDEBUG   = "\033[38;5;24m"
	clogTRACE   = "\033[38;5;240m"
)

// String wrapper to append text
type AppendLog struct {
	Txt string
}

// Log item to print on console
type LogItem struct {
	Color  string
	Ctx    common.DCtx
	Prefix string
	Msgs   []string
}

// Channel with log items stream
var logItemCH = make(chan LogItem)

// init logger
//
// Starts listening on log items to print (in a new thread)
func init() {
	go getMsgsAndPrintToConsole()
}

// Wait for log item and print it
// It blocks thread - so mus be used in a new thread
func getMsgsAndPrintToConsole() {
	for {
		logItem := <-logItemCH
		printColor(logItem)
	}
}

// Configures logger using provided log level
func Configure(lstr string) {
	if len(lstr) == 0 {
		printInfo = true
		printWarning = true
		return
	}
	levels := strings.Split(lstr, "|")
	for _, lvl := range levels {
		switch lvl {
		case "TRACE":
			printTrace = true
			printDebug = true
			printInfo = true
			printWarning = true
		case "DEBUG":
			printDebug = true
			printInfo = true
			printWarning = true
		case "INFO":
			printInfo = true
			printWarning = true
		case "WARNING":
			printWarning = true
		case "NOBANNER":
			PrintBaner = false
		}
	}
}

// Prints Error log and exit application with provided exit code
func Panic(ctx common.DCtx, err error, exitCode int, msgs ...string) {
	Error(ctx, err, msgs...)
	os.Exit(exitCode)
}

// Prints Error log
func Error(ctx common.DCtx, err error, msgs ...string) {
	msgs2 := append(msgs, fmt.Sprintf("%v", err))
	logItemCH <- LogItem{
		Color:  ClogRED,
		Ctx:    ctx,
		Prefix: "ERR",
		Msgs:   msgs2,
	}
}

// Prints Error log and append printed text to provided append object
func ErrorAndAppend(ctx common.DCtx, append *AppendLog, err error, msgs ...string) {
	appendLog(append, msgs)
	Error(ctx, err, msgs...)
}

// Prints Warning log
func Warning(ctx common.DCtx, msgs ...string) {
	if printWarning {
		logItemCH <- LogItem{
			Color:  ClogYELLOW,
			Ctx:    ctx,
			Prefix: "WAR",
			Msgs:   msgs,
		}
	}
}

// Prints Warning log and append printed text to provided append object
func WarningAndAppend(ctx common.DCtx, append *AppendLog, msgs ...string) {
	if printWarning {
		appendLog(append, msgs)
		Warning(ctx, msgs...)
	}
}

// Prints Info log
func Info(ctx common.DCtx, msgs ...string) {
	if printInfo {
		logItemCH <- LogItem{
			Color:  ClogCYAN,
			Ctx:    ctx,
			Prefix: "INF",
			Msgs:   msgs,
		}
	}
}

// Prints Info log and append printed text to provided append object
func InfoAndAppend(ctx common.DCtx, append *AppendLog, msgs ...string) {
	if printInfo {
		appendLog(append, msgs)
		Info(ctx, msgs...)
	}
}

// Prints Debug log
func Debug(ctx common.DCtx, msgs ...string) {
	if printDebug {
		logItemCH <- LogItem{
			Color:  ClogDEBUG,
			Ctx:    ctx,
			Prefix: "DEB",
			Msgs:   msgs,
		}
	}
}

// Prints Debug log and append printed text to provided append object
func DebugAndAppend(ctx common.DCtx, append *AppendLog, msgs ...string) {
	if printDebug {
		appendLog(append, msgs)
		Debug(ctx, msgs...)
	}
}

// Prints Trace log
func Trace(ctx common.DCtx, msgs ...string) {
	if printTrace {
		logItemCH <- LogItem{
			Color:  clogTRACE,
			Ctx:    ctx,
			Prefix: "TRA",
			Msgs:   msgs,
		}
	}
}

// Prints Trace log and append printed text to provided append object
func TraceAndAppend(ctx common.DCtx, append *AppendLog, msgs ...string) {
	if printTrace {
		appendLog(append, msgs)
		Trace(ctx, msgs...)
	}
}

// Prints Test log
//
// It should be used only in development
func TestLog(ctx common.DCtx, msgs ...string) {
	logItemCH <- LogItem{
		Color:  ClogMAGENTA,
		Ctx:    ctx,
		Prefix: "TST",
		Msgs:   msgs,
	}
}

// Creates text to print on console
//
// It is prefixed using context. All messages are in one line, separated by '|'
func buildLogEntry(ctx common.DCtx, msgs []string) string {
	logEntry := ctx.ThreadId
	if len(logEntry) == 0 {
		logEntry = "main"
	}
	logEntry += "|"
	if len(ctx.JobId) > 0 {
		logEntry += ctx.JobId
	} else {
		logEntry += "-"
	}
	logEntry += "|"
	if len(ctx.EnvName) > 0 {
		logEntry += ctx.EnvName
	} else {
		logEntry += "-"
	}
	logEntry += "|"
	if len(ctx.SrvName) > 0 {
		logEntry += ctx.SrvName
	} else {
		logEntry += "-"
	}
	logEntry += "|"
	if len(ctx.ProjName) > 0 {
		logEntry += ctx.ProjName
	} else {
		logEntry += "-"
	}
	logEntry += "|>> "
	for i, msg := range msgs {
		if i > 0 {
			logEntry += "|"
		}
		logEntry += msg
	}
	return logEntry
}

// Appends messages in a one line to provided text in append object
func appendLog(append *AppendLog, msgs []string) {
	for i, msg := range msgs {
		if i > 0 {
			append.Txt += "|"
		}
		append.Txt += msg
	}
	append.Txt += "\n"
}

// Prints log item on console
func printColor(item LogItem) {
	startColor(item.Color)
	if len(item.Prefix) > 0 {
		fmt.Print(item.Prefix + "|")
	}
	logItemTxt := buildLogEntry(item.Ctx, item.Msgs)
	fmt.Print(logItemTxt)
	endColor()
	fmt.Print("\n")
}

// Starts new color on console
func startColor(color string) {
	fmt.Print(color)
}

// Clear color on console. Restore default color
func endColor() {
	fmt.Print("\033[0m")
}
