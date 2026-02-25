package console

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/patppuccin/kredenv/consts"
)

var dimmed = color.New(color.FgHiBlack)

type level struct {
	color *color.Color
	tag   string
}

func (l level) log(msg string) {
	fmt.Fprintln(os.Stdout, l.color.Sprint(l.tag), msg)
}

func (l level) group(title string, msgs []string) {
	fmt.Fprintln(os.Stdout, l.color.Sprint(l.tag), title)
	for _, msg := range msgs {
		fmt.Fprintln(os.Stdout, " ", dimmed.Sprint("•"), msg)
	}
}

var (
	lvlSuccess = level{color.New(color.FgGreen), "[+]"}
	lvlDebug   = level{color.New(color.FgHiBlack), "[~]"}
	lvlInfo    = level{color.New(color.FgCyan), "[i]"}
	lvlWarn    = level{color.New(color.FgYellow), "[!]"}
	lvlError   = level{color.New(color.FgRed), "[x]"}
)

func Success(msg string) { lvlSuccess.log(msg) }
func Debug(msg string)   { lvlDebug.log(msg) }
func Info(msg string)    { lvlInfo.log(msg) }
func Warn(msg string)    { lvlWarn.log(msg) }
func Error(msg string)   { lvlError.log(msg) }

func SuccessGroup(title string, msgs []string) { lvlSuccess.group(title, msgs) }
func DebugGroup(title string, msgs []string)   { lvlDebug.group(title, msgs) }
func InfoGroup(title string, msgs []string)    { lvlInfo.group(title, msgs) }
func WarnGroup(title string, msgs []string)    { lvlWarn.group(title, msgs) }
func ErrorGroup(title string, msgs []string)   { lvlError.group(title, msgs) }

func Banner(msg string) string {
	return consts.AppBanner + "\n" + color.New(color.FgCyan).Sprint(msg) + "\n"
}
