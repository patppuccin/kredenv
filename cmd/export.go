package cmd

import "github.com/spf13/cobra"

const helpExportCmd = "Exports secrets from the keyring to stdout or a file"

var (
	flagExportAll     bool
	flagExportFormat  string
	flagExportOutput  string
	flagExportEncrypt bool
)

var exportCmd = &cobra.Command{
	Use:           "export",
	Short:         helpExportCmd,
	Long:          banner(helpExportCmd),
	Args:          cobra.NoArgs,
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	exportCmd.Flags().SortFlags = false
	exportCmd.Flags().BoolVar(&flagExportAll, "all", false, "Export all keys in the keyring")
	exportCmd.Flags().StringVarP(&flagExportFormat, "format", "f", "env", "Export format (env, json, yaml, toml)")
	exportCmd.Flags().StringVarP(&flagExportOutput, "output", "o", "", "Output file path (defaults to stdout)")
	exportCmd.Flags().BoolVar(&flagExportEncrypt, "encrypt", false, "Encrypt the exported file")
}
