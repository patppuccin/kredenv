package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

func banner(msg string) string {
	return consts.AppBanner + "\n" +
		color.New(color.FgBlue).Sprint(msg) + "\n"
}

func versionString() string {
	date := consts.BuildDate
	if t, err := time.Parse(time.RFC3339, consts.BuildDate); err == nil {
		date = t.UTC().Format("02 Jan 2006 15:04 UTC")
	}
	return fmt.Sprintf("%s (commit: %s, built: %s)", consts.AppVersion, consts.BuildCommit, date)
}

var KredEnvCmd = &cobra.Command{
	Use:           consts.AppName,
	Short:         consts.AppDesc,
	Long:          banner(consts.AppDesc),
	Version:       versionString(),
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
		termactions.Log().Error(err.Error())
		os.Stdout.WriteString("\n")
		os.Exit(1)
		return nil
	})

	KredEnvCmd.AddGroup(&cobra.Group{ID: "setup", Title: "Setup Commands:"})
	KredEnvCmd.AddCommand(setupCmd)
	KredEnvCmd.AddCommand(initCmd)
	KredEnvCmd.AddCommand(hookCmd)

	KredEnvCmd.AddGroup(&cobra.Group{ID: "env", Title: "Environment Commands:"})
	KredEnvCmd.AddCommand(whichCmd)
	KredEnvCmd.AddCommand(loadCmd)
	KredEnvCmd.AddCommand(unloadCmd)
	KredEnvCmd.AddCommand(execCmd)
	KredEnvCmd.AddCommand(validateCmd)
	KredEnvCmd.AddCommand(injectCmd)

	KredEnvCmd.AddGroup(&cobra.Group{ID: "secrets", Title: "Secrets Commands:"})
	KredEnvCmd.AddCommand(setCmd)
	KredEnvCmd.AddCommand(getCmd)
	KredEnvCmd.AddCommand(listCmd)
	KredEnvCmd.AddCommand(deleteCmd)
	KredEnvCmd.AddCommand(exportCmd)
	KredEnvCmd.AddCommand(importCmd)
}
