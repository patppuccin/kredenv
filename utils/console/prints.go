package console

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/patppuccin/kredenv/consts"
)

var writer io.Writer = os.Stdout

var (
	fmtSuccess = func(content string) string { return color.New(color.FgGreen).Sprint(content) }
	fmtDebug   = func(content string) string { return color.New(color.FgHiBlack).Sprint(content) }
	fmtInfo    = func(content string) string { return color.New(color.FgBlue).Sprint(content) }
	fmtWarn    = func(content string) string { return color.New(color.FgYellow).Sprint(content) }
	fmtError   = func(content string) string { return color.New(color.FgRed).Sprint(content) }
	fmtKV      = func(content string) string { return color.New(color.FgWhite).Sprint(content) }
	fmtDesc    = func(content string) string { return color.New(color.FgGreen).Sprint(content) }
)

func Banner(msg string) string { return consts.AppBanner + "\n" + fmtDesc(msg) + "\n" }

func Success(msg string) { logSingle(fmtSuccess("(✓)"), msg) }
func Debug(msg string)   { logSingle(fmtDebug("(~)"), msg) }
func Info(msg string)    { logSingle(fmtInfo("(i)"), msg) }
func Warn(msg string)    { logSingle(fmtWarn("(!)"), msg) }
func Error(msg string)   { logSingle(fmtError("(x)"), msg) }

func SuccessGroup(title string, msgs []string) { logGroup(fmtSuccess("SUCCESS:"), title, msgs) }
func DebugGroup(title string, msgs []string)   { logGroup(fmtDebug("DEBUG:"), title, msgs) }
func InfoGroup(title string, msgs []string)    { logGroup(fmtInfo("INFO:"), title, msgs) }
func WarnGroup(title string, msgs []string)    { logGroup(fmtWarn("WARN:"), title, msgs) }
func ErrorGroup(title string, msgs []string)   { logGroup(fmtError("ERROR:"), title, msgs) }

func logSingle(pfx, msg string) { fmt.Fprintf(writer, "%s %s\n", pfx, msg) }
func logGroup(pfx, title string, msgs []string) {
	fmt.Fprintf(writer, "%s %s\n", pfx, title)
	for _, msg := range msgs {
		fmt.Fprintf(writer, "  %s\n", msg)
	}
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
