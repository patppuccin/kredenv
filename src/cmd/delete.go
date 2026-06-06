package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpDeleteCmd = "Delete one or more secrets from the kredenv store"

var (
	flagDeleteNamespace string
)

var deleteCmd = &cobra.Command{
	Use:           "delete <key> [keys...]",
	Short:         helpDeleteCmd,
	Long:          banner(helpDeleteCmd),
	GroupID:       "secrets",
	Aliases:       []string{"del"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			termactions.Log().Error("Expected at least one argument, got 0")
			os.Exit(1)
		}

		password, err := auth.Retrieve()
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}

		s, err := store.Open(password)
		if err != nil {
			termactions.Log().Error("Could not open store")
			os.Exit(1)
		}
		defer s.Close()

		for _, key := range args {
			if flagDeleteNamespace != "" {
				key = flagDeleteNamespace + ":" + key
			}
			if err := s.Delete(key); err != nil {
				termactions.Log().Error(err.Error())
				os.Exit(1)
			}
			termactions.Log().Info("Deleted " + key)
		}
	},
}

func init() {
	deleteCmd.Flags().SortFlags = false
	deleteCmd.Flags().StringVarP(&flagDeleteNamespace, "namespace", "n", "", "Delete key from a specific namespace")
}
