package cmd

import "github.com/spf13/cobra"

const helpSetCmd = "Stores a key in the keyring, prompts if value not provided"

var setCmd = &cobra.Command{
	Use:           "set <key> [value]",
	Short:         helpSetCmd,
	Long:          banner(helpSetCmd),
	Args:          cobra.RangeArgs(1, 2),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	setCmd.Flags().SortFlags = false
}
