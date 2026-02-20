package cmd

import "github.com/spf13/cobra"

const helpGetCmd = "Prints the value of a key from the keyring"

var getCmd = &cobra.Command{
	Use:           "get <key>",
	Short:         helpGetCmd,
	Long:          banner(helpGetCmd),
	Args:          cobra.ExactArgs(1),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	getCmd.Flags().SortFlags = false
}
