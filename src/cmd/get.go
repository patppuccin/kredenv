package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpGetCmd = "Retrieve a secret from the kredenv store"

var (
	flagGetNamespace string
)

var getCmd = &cobra.Command{
	Use:           "get <key>",
	Short:         helpGetCmd,
	Long:          banner(helpGetCmd),
	GroupID:       "secrets",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			termactions.Log().Error("Expected exactly one argument, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
		// fail if both colon syntax and --namespace flag are used
		if flagGetNamespace != "" && strings.Contains(args[0], ":") {
			termactions.Log().Error("Cannot use both 'namespace:key' syntax and --namespace (-n) flag")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if flagGetNamespace != "" {
			key = flagGetNamespace + ":" + key
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

		value, err := s.Get(key)
		if err != nil {
			termactions.Log().Error("Could not retrieve " + key)
			os.Exit(1)
		}

		fmt.Println(value)
	},
}

func init() {
	getCmd.Flags().SortFlags = false
	getCmd.Flags().StringVarP(&flagGetNamespace, "namespace", "n", "", "Get key from a specific namespace")
}
