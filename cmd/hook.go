package cmd

import (
	"github.com/spf13/cobra"
)

const helpHookCmd = "Emits a shell hook script (supported: bash, zsh, fish, powershell, nushell)"

var hookCmd = &cobra.Command{
	Use:           "hook <shell>",
	Short:         helpHookCmd,
	Long:          banner(helpHookCmd),
	Args:          cobra.ExactArgs(1),
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	hookCmd.Flags().SortFlags = false
}
