package cmd

import (
	"os"
	"strings"

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
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			termactions.Log().Error("Expected at least one argument, got 0")
			os.Exit(1)
		}
		if flagDeleteNamespace != "" {
			for _, arg := range args {
				if strings.Contains(arg, ":") {
					termactions.Log().Error("Cannot use both 'namespace:key' syntax and --namespace (-n) flag")
					os.Exit(1)
				}
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
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
				termactions.Log().Error("Could not delete " + key + ": " + err.Error())
				continue
			}
			termactions.Log().Success("Deleted " + key)
		}
	},
}

func init() {
	deleteCmd.Flags().SortFlags = false
	deleteCmd.Flags().StringVarP(&flagDeleteNamespace, "namespace", "n", "", "Delete key from a specific namespace")
}
