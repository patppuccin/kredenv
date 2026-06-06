package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

var KredEnvCmd = &cobra.Command{
	Use:           consts.AppName,
	Short:         consts.AppDesc,
	Long:          banner(consts.AppDesc),
	Version:       consts.AppVersion,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func banner(msg string) string {
	return consts.AppBanner + "\n" +
		color.New(color.FgBlue).Sprint(msg) + "\n"
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
	KredEnvCmd.AddCommand(loadCmd)
	KredEnvCmd.AddCommand(unloadCmd)
	KredEnvCmd.AddCommand(execCmd)
	KredEnvCmd.AddCommand(whichCmd)
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
