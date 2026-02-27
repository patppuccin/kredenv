package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpGetCmd = "Prints the value of a key from the keyring"

var (
	flagGetNamespace string
)

var getCmd = &cobra.Command{
	Use:           "get <key>",
	Short:         helpGetCmd,
	Long:          console.Banner(helpGetCmd),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			console.Error("Expected exactly one argument, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		key := args[0]
		if flagGetNamespace != "" {
			key = flagGetNamespace + ":" + key
		}

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
	getCmd.Flags().StringVarP(&flagGetNamespace, "namespace", "n", "", "Get keys from a specific namespace")
}
