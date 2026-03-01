package cmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/store"
	"github.com/spf13/cobra"
)

const helpSetCmd = "Store a secret in the kredenv store"

var (
	flagSetNamespace string
)

var setCmd = &cobra.Command{
	Use:           "set <key> [value]",
	Short:         helpSetCmd,
	Long:          console.Banner(helpSetCmd),
	GroupID:       "secrets",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			console.Error("Expected at least one argument, got 0")
			os.Exit(1)
		case 1, 2: // valid
		default:
			console.Error("Expected at most two arguments, got " + fmt.Sprintf("%d", len(args)))
			os.Exit(1)
		}

		key := args[0]
		if flagSetNamespace != "" {
			key = flagSetNamespace + ":" + key
		}

		var value string
		if len(args) == 2 {
			value = args[1]
		} else {
			if !isatty.IsTerminal(os.Stdin.Fd()) {
				console.Error("No value provided — interactive set requires a terminal to prompt")
				os.Exit(1)
			}

			var err error
			value, err = console.PromptSecret("Value for " + key + ": ")
			if err != nil {
				console.Error("Could not read input")
				os.Exit(1)
			}
		}

		if value == "" {
			console.Warn("Value cannot be empty")
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

		if err := s.Set(key, value); err != nil {
			console.Error("Could not store " + key)
			os.Exit(1)
		}

		console.Success("Stored " + key)
	},
}

func init() {
	setCmd.Flags().SortFlags = false
	setCmd.Flags().StringVarP(&flagSetNamespace, "namespace", "n", "", "Store key under a specific namespace")
}
