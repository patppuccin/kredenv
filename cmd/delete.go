package cmd

import "github.com/spf13/cobra"

const helpDeleteCmd = "Deletes one or more keys from the keyring"

var deleteCmd = &cobra.Command{
	Use:           "delete <key> [keys...]",
	Short:         helpDeleteCmd,
	Long:          banner(helpDeleteCmd),
	Args:          cobra.MinimumNArgs(1),
	GroupID:       "keyring",
	Aliases:       []string{"del"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	deleteCmd.Flags().SortFlags = false
}
