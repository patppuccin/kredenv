package cmd

import (
	"fmt"
	"os"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpGetCmd = "Prints the value of a key from the keyring"

var getCmd = &cobra.Command{
	Use:           "get <key>",
	Short:         helpGetCmd,
	Long:          console.Banner(helpGetCmd),
	Args:          cobra.ExactArgs(1),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value, err := keyring.Get(key)
		if err != nil {
			console.Error("Could not retrieve key " + key)
			os.Exit(1)
		}
		fmt.Println(value)
	},
}

func init() {
	getCmd.Flags().SortFlags = false
}
