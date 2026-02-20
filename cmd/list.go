package cmd

import "github.com/spf13/cobra"

const helpListCmd = "Lists keys from the local .kredsfile or the keyring"

var (
	flagListAll        bool
	flagListShowValues bool
)

var listCmd = &cobra.Command{
	Use:           "list",
	Short:         helpListCmd,
	Long:          banner(helpListCmd),
	Args:          cobra.NoArgs,
	GroupID:       "keyring",
	Aliases:       []string{"ls"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	listCmd.Flags().SortFlags = false
	listCmd.Flags().BoolVarP(&flagListAll, "all", "a", false, "List all keys in the keyring instead")
	listCmd.Flags().BoolVar(&flagListShowValues, "show-values", false, "Show secret values (use with caution)")
}
