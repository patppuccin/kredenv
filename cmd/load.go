package cmd

import "github.com/spf13/cobra"

const helpLoadCmd = "Resolves the .kredsfile and injects it into the environment"

var (
	flagLoadTransient bool
)

var loadCmd = &cobra.Command{
	Use:           "load [-- <command>]",
	Short:         helpLoadCmd,
	Long:          banner(helpLoadCmd),
	Args:          cobra.ArbitraryArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	loadCmd.Flags().SortFlags = false
	loadCmd.Flags().BoolVar(&flagLoadTransient, "once", false, "Inject secrets for this command only")
}
