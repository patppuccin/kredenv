package cmd

import "github.com/spf13/cobra"

const helpWhichCmd = "Prints the path to the .kredsfile that will be used"

var whichCmd = &cobra.Command{
	Use:           "which",
	Short:         helpWhichCmd,
	Long:          banner(helpWhichCmd),
	Args:          cobra.NoArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	whichCmd.Flags().SortFlags = false
}
