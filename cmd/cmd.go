package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/consts"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/spf13/cobra"
)

// type cliLogger struct{}

// func (c cliLogger) Debug(msg string) {
// 	os.Stdout.WriteString(color.New(color.FgHiBlack).Sprint("[~]") + " " + msg + "\n")
// }

// func (c cliLogger) Info(msg string) {
// 	os.Stdout.WriteString(color.New(color.FgBlue).Sprint("[i]") + " " + msg + "\n")
// }

// func (c cliLogger) Warn(msg string) {
// 	os.Stdout.WriteString(color.New(color.FgYellow).Sprint("[!]") + " " + msg + "\n")
// }

// func (c cliLogger) Error(msg string) {
// 	os.Stdout.WriteString(color.New(color.FgRed).Sprint("[x]") + " " + msg + "\n")
// }

// func (c cliLogger) DebugMulti(title string, msgs []string) {
// 	prefix := color.New(color.FgHiBlack).Sprint("[~]")
// 	os.Stdout.WriteString(prefix + " " + title + "\n")
// 	for _, msg := range msgs {
// 		os.Stdout.WriteString("    • " + msg + "\n")
// 	}
// }

// func (c cliLogger) InfoMulti(title string, msgs []string) {
// 	prefix := color.New(color.FgBlue).Sprint("[i]")
// 	os.Stdout.WriteString(prefix + " " + title + "\n")
// 	for _, msg := range msgs {
// 		os.Stdout.WriteString("    • " + msg + "\n")
// 	}
// }

// func (c cliLogger) WarnMulti(title string, msgs []string) {
// 	prefix := color.New(color.FgYellow).Sprint("[!]")
// 	os.Stdout.WriteString(prefix + " " + title + "\n")
// 	for _, msg := range msgs {
// 		os.Stdout.WriteString("    • " + msg + "\n")
// 	}
// }

// func (c cliLogger) ErrorMulti(title string, msgs []string) {
// 	prefix := color.New(color.FgRed).Sprint("[x]")
// 	os.Stdout.WriteString(prefix + " " + title + "\n")
// 	for _, msg := range msgs {
// 		os.Stdout.WriteString("    • " + msg + "\n")
// 	}
// }

// var log = cliLogger{}

// func banner(msg string) string {
// 	var sb strings.Builder
// 	sb.WriteString(consts.AppBanner)
// 	sb.WriteString("\n")
// 	sb.WriteString(color.New(color.FgBlue).Sprint(msg))
// 	sb.WriteString("\n")
// 	return sb.String()
// }

var KredEnvCmd = &cobra.Command{
	Use:           consts.AppName,
	Short:         consts.AppDesc,
	Long:          console.Banner(consts.AppDesc),
	Version:       consts.AppVersion,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.EnableCommandSorting = false
	KredEnvCmd.CompletionOptions.DisableDefaultCmd = true

	KredEnvCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Help()
		os.Stdout.WriteString("\n")
		console.Error(err.Error())
		os.Stdout.WriteString("\n")
		os.Exit(1)
		return nil
	})

	KredEnvCmd.AddGroup(&cobra.Group{ID: "setup", Title: "Setup Commands:"})
	KredEnvCmd.AddCommand(initCmd)
	KredEnvCmd.AddCommand(setupCmd)
	KredEnvCmd.AddCommand(hookCmd)

	KredEnvCmd.AddGroup(&cobra.Group{ID: "env", Title: "Environment Commands:"})
	KredEnvCmd.AddCommand(loadCmd)
	KredEnvCmd.AddCommand(unloadCmd)
	KredEnvCmd.AddCommand(whichCmd)
	KredEnvCmd.AddCommand(validateCmd)

	KredEnvCmd.AddGroup(&cobra.Group{ID: "keyring", Title: "Keyring Commands:"})
	KredEnvCmd.AddCommand(getCmd)
	KredEnvCmd.AddCommand(setCmd)
	KredEnvCmd.AddCommand(deleteCmd)
	KredEnvCmd.AddCommand(listCmd)
	KredEnvCmd.AddCommand(exportCmd)
	KredEnvCmd.AddCommand(importCmd)
}
