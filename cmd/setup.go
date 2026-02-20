package cmd

import "github.com/spf13/cobra"

const helpSetupCmd = "Finds missing env vars and prompts to be stored in the keyring"

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         helpSetupCmd,
	Long:          banner(helpSetupCmd),
	Args:          cobra.NoArgs,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	setupCmd.Flags().SortFlags = false
}
