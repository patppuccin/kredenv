package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpSetCmd = "Store a secret in the kredenv store"

var (
	flagSetNamespace string
)

var setCmd = &cobra.Command{
	Use:           "set <key> [value]",
	Short:         helpSetCmd,
	Long:          banner(helpSetCmd),
	GroupID:       "secrets",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			termactions.Log().Error("Expected at least one argument, got 0")
			os.Exit(1)
		case 1, 2: // valid
		default:
			termactions.Log().Error("Expected at most two arguments, got " + fmt.Sprintf("%d", len(args)))
			os.Exit(1)
		}
		if flagSetNamespace != "" && strings.Contains(args[0], ":") {
			termactions.Log().Error("Cannot use both 'namespace:key' syntax and --namespace (-n) flag")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if flagSetNamespace != "" {
			key = flagSetNamespace + ":" + key
		}

		var value string
		if len(args) == 2 {
			value = args[1]
		} else {
			if !isatty.IsTerminal(os.Stdin.Fd()) {
				termactions.Log().Error("No value provided — interactive set requires a terminal to prompt")
				os.Exit(1)
			}

			var err error
			value, err = termactions.Secret().WithLabel("Value for " + key).Render()
			if err != nil {
				if err == termactions.ErrInterrupted {
					termactions.Log().Warn("Set action aborted by user")
					os.Exit(1)
				}
				termactions.Log().Error("Could not read user input")
				os.Exit(1)
			}
		}

		if value == "" {
			termactions.Log().Warn("Value cannot be empty")
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

		if err := s.Set(key, value); err != nil {
			termactions.Log().Error("Could not store " + key)
			os.Exit(1)
		}

		termactions.Log().Success("Stored " + key)
	},
}

func init() {
	setCmd.Flags().SortFlags = false
	setCmd.Flags().StringVarP(&flagSetNamespace, "namespace", "n", "", "Store key under a specific namespace")
}
