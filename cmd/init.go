package cmd

import (
	"github.com/spf13/cobra"
)

const helpInitCmd = "Initializes a minimal .kredsfile in the current directory"

var (
	flagInitOverwrite bool
)

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         helpInitCmd,
	Long:          banner(helpInitCmd),
	Args:          cobra.NoArgs,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	initCmd.Flags().SortFlags = false
	initCmd.Flags().BoolVar(&flagInitOverwrite, "overwrite", false, "Overwrite existing .kredsfile")
}
