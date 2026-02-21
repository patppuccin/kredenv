package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpDeleteCmd = "Deletes one or more keys from the keyring"

var deleteCmd = &cobra.Command{
	Use:           "delete <key> [keys...]",
	Short:         helpDeleteCmd,
	Long:          console.Banner(helpDeleteCmd),
	Args:          cobra.MinimumNArgs(1),
	GroupID:       "keyring",
	Aliases:       []string{"del"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		for _, key := range args {
			if err := keyring.Delete(key); err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
			console.Info("Deleted Key: " + key)
		}
	},
}

func init() {
	deleteCmd.Flags().SortFlags = false
}
