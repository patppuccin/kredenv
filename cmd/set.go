package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpSetCmd = "Stores a key in the keyring, prompts if value not provided"

var (
	flagSetNamespace string
)

var setCmd = &cobra.Command{
	Use:           "set <key> [value]",
	Short:         helpSetCmd,
	Long:          console.Banner(helpSetCmd),
	GroupID:       "keyring",
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
				console.Error("No value provided - interactive set requires a terminal to prompt")
				os.Exit(1)
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Value for %s: ", key)
			input, err := reader.ReadString('\n')
			if err != nil {
				console.Error("Could not read input")
				os.Exit(1)
			}
			value = strings.TrimSpace(input)
		}

		if value == "" {
			console.Warn("Value cannot be empty")
			os.Exit(1)
		}

		if err := keyring.Set(key, value); err != nil {
			console.Error("Could not store " + key)
			os.Exit(1)
		}

		console.Success("Stored " + key + " in keyring")
	},
}

func init() {
	setCmd.Flags().SortFlags = false
	setCmd.Flags().StringVarP(&flagSetNamespace, "namespace", "n", "", "Store key under a specific namespace")
}
