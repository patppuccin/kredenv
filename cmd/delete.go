package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpDeleteCmd = "Deletes one or more keys from the keyring"

var (
	flagDeleteNamespace string
)

var deleteCmd = &cobra.Command{
	Use:           "delete <key> [keys...]",
	Short:         helpDeleteCmd,
	Long:          console.Banner(helpDeleteCmd),
	GroupID:       "keyring",
	Aliases:       []string{"del"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			console.Error("Expected at least one argument, got 0")
			os.Exit(1)
		}
		for _, key := range args {
			if flagDeleteNamespace != "" {
				key = flagDeleteNamespace + ":" + key
			}
			if err := keyring.Delete(key); err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
			console.Info("Deleted key: " + key)
		}
	},
}

func init() {
	deleteCmd.Flags().SortFlags = false
	deleteCmd.Flags().StringVarP(&flagDeleteNamespace, "namespace", "n", "", "Delete key from a specific namespace")
}
