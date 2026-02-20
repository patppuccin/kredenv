package cmd

import "github.com/spf13/cobra"

const helpImportCmd = "Imports secrets from a file into the keyring"

var (
	flagImportOverwrite bool
)

var importCmd = &cobra.Command{
	Use:           "import <file>",
	Short:         helpImportCmd,
	Long:          banner(helpImportCmd),
	Args:          cobra.ExactArgs(1),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	importCmd.Flags().SortFlags = false
	importCmd.Flags().BoolVar(&flagImportOverwrite, "overwrite", false, "Overwrite existing keys if they already exist in the keyring")
}
