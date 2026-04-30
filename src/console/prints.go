package console

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/patppuccin/kredenv/src/consts"
)

var writer io.Writer = os.Stdout

var (
	fmtSuccess = func(txt string) string { return color.New(color.FgGreen).Sprint(txt) }
	fmtDebug   = func(txt string) string { return color.New(color.FgHiBlack).Sprint(txt) }
	fmtInfo    = func(txt string) string { return color.New(color.FgBlue).Sprint(txt) }
	fmtWarn    = func(txt string) string { return color.New(color.FgYellow).Sprint(txt) }
	fmtError   = func(txt string) string { return color.New(color.FgRed).Sprint(txt) }
	fmtKV      = func(txt string) string { return color.New(color.FgWhite).Sprint(txt) }
	fmtDesc    = func(txt string) string { return color.New(color.FgGreen).Sprint(txt) }
)

func VSpacer(lines int) { fmt.Fprint(writer, strings.Repeat("\n", lines)) }

func Banner(msg string) string { return consts.AppBanner + "\n" + fmtDesc(msg) + "\n" }

func Success(msg string) { logSingle(fmtSuccess("(✓)"), msg) }
func Debug(msg string)   { logSingle(fmtDebug("(~)"), msg) }
func Info(msg string)    { logSingle(fmtInfo("(i)"), msg) }
func Warn(msg string)    { logSingle(fmtWarn("(!)"), msg) }
func Error(msg string)   { logSingle(fmtError("(x)"), msg) }

func SuccessGroup(title string, msgs ...string) { logGroup(fmtSuccess("SUCCESS:"), title, msgs...) }
func DebugGroup(title string, msgs ...string)   { logGroup(fmtDebug("DEBUG:"), title, msgs...) }
func InfoGroup(title string, msgs ...string)    { logGroup(fmtInfo("INFO:"), title, msgs...) }
func WarnGroup(title string, msgs ...string)    { logGroup(fmtWarn("WARN:"), title, msgs...) }
func ErrorGroup(title string, msgs ...string)   { logGroup(fmtError("ERROR:"), title, msgs...) }

func SuccessBlock(title string, msgs ...string) { logBlock(fmtSuccess(title), msgs...) }
func DebugBlock(title string, msgs ...string)   { logBlock(fmtDebug(title), msgs...) }
func InfoBlock(title string, msgs ...string)    { logBlock(fmtInfo(title), msgs...) }
func WarnBlock(title string, msgs ...string)    { logBlock(fmtWarn(title), msgs...) }
func ErrorBlock(title string, msgs ...string)   { logBlock(fmtError(title), msgs...) }

func logSingle(pfx, msg string) { fmt.Fprintf(writer, "%s %s\n", pfx, msg) }

func logGroup(pfx, title string, msgs ...string) {
	fmt.Fprintf(writer, "%s %s\n", pfx, title)
	for _, msg := range msgs {
		fmt.Fprintf(writer, "  %s\n", msg)
	}
}
func logBlock(title string, msgs ...string) {
	fmt.Fprintf(writer, "%s\n", title)
	fmt.Fprintln(writer)
	for _, msg := range msgs {
		fmt.Fprintf(writer, "%s\n", msg)
	}
	fmt.Fprintln(writer)
}

type KV struct {
	key   string
	value string
}

func KeyOnly(key string) KV         { return KV{key: key} }
func KeyValue(key, value string) KV { return KV{key: key, value: value} }

func PrintKVs(title string, entries ...KV) {
	fmt.Fprintf(writer, "%s %s\n", fmtKV("KV Store:"), title)
	for _, e := range entries {
		if e.value == "" {
			fmt.Fprintf(writer, "  %s\n", e.key)
		} else {
			fmt.Fprintf(writer, "  %s = %s\n", e.key, e.value)
		}
	}
}
