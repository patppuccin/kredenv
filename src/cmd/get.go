package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/console"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/spf13/cobra"
)

const helpGetCmd = "Retrieve a secret from the kredenv store"

var (
	flagGetNamespace string
)

var getCmd = &cobra.Command{
	Use:           "get <key>",
	Short:         helpGetCmd,
	Long:          console.Banner(helpGetCmd),
	GroupID:       "secrets",
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

		value, err := s.Get(key)
		if err != nil {
			console.Error("Could not retrieve " + key)
			os.Exit(1)
		}

		fmt.Println(value)
	},
}

func init() {
	getCmd.Flags().SortFlags = false
	getCmd.Flags().StringVarP(&flagGetNamespace, "namespace", "n", "", "Get key from a specific namespace")
}
