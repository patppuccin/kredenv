package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/store"
	"github.com/spf13/cobra"
)

const helpDeleteCmd = "Delete one or more secrets from the kredenv store"

var (
	flagDeleteNamespace string
)

var deleteCmd = &cobra.Command{
	Use:           "delete <key> [keys...]",
	Short:         helpDeleteCmd,
	Long:          console.Banner(helpDeleteCmd),
	GroupID:       "secrets",
	Aliases:       []string{"del"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			console.Error("Expected at least one argument, got 0")
			os.Exit(1)
		}

		password, err := auth.Retrieve()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		s, err := store.Open(password)
		if err != nil {
			console.Error("Could not open store")
			os.Exit(1)
		}
		defer s.Close()

		for _, key := range args {
			if flagDeleteNamespace != "" {
				key = flagDeleteNamespace + ":" + key
			}
			if err := s.Delete(key); err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
			console.Info("Deleted " + key)
		}
	},
}

func init() {
	deleteCmd.Flags().SortFlags = false
	deleteCmd.Flags().StringVarP(&flagDeleteNamespace, "namespace", "n", "", "Delete key from a specific namespace")
}
