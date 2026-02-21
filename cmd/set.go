package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpSetCmd = "Stores a key in the keyring, prompts if value not provided"

var setCmd = &cobra.Command{
	Use:           "set <key> [value]",
	Short:         helpSetCmd,
	Long:          console.Banner(helpSetCmd),
	Args:          cobra.RangeArgs(1, 2),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var value string

		if len(args) == 2 {
			value = args[1]
		} else {
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
}
