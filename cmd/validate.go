package cmd

import "github.com/spf13/cobra"

const helpValidateCmd = "Validates .kredsfile syntax"

var validateCmd = &cobra.Command{
	Use:           "validate [file]",
	Short:         helpValidateCmd,
	Long:          banner(helpValidateCmd),
	Args:          cobra.MaximumNArgs(1),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	validateCmd.Flags().SortFlags = false
}
