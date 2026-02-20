package cmd

import "github.com/spf13/cobra"

const helpUnloadCmd = "Unloads the .kredsfile from the environment"

var unloadCmd = &cobra.Command{
	Use:           "unload",
	Short:         helpUnloadCmd,
	Long:          banner(helpUnloadCmd),
	Args:          cobra.NoArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	unloadCmd.Flags().SortFlags = false
}
