package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/consts"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/spf13/cobra"
)

const helpLoadCmd = "Loads the secrets from the .kredsfile in scope into the environment"

var (
	flagLoadNamespace string
)

var loadCmd = &cobra.Command{
	Use:           "load",
	Short:         helpLoadCmd,
	Long:          console.Banner(helpLoadCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		if os.Getenv("__KREDENV_BIN") == "" {
			console.Warn("Shell hook not detected, run '" + consts.AppName + " hook <shell>' to set up")
			return
		}

		loaded := os.Getenv("KREDENV_LOADED_VARS")
		if loaded == "" {
			console.Warn("No secrets currently loaded")
			return
		}

		keys := strings.Split(loaded, ",")
		title := strconv.Itoa(len(keys)) + " secrets loaded"
		if flagLoadNamespace != "" {
			title += " (namespace: " + flagLoadNamespace + ")"
		}

		console.InfoGroup(title, keys...)
	},
}

func init() {
	loadCmd.Flags().SortFlags = false
	loadCmd.Flags().StringVarP(&flagLoadNamespace, "namespace", "n", "", "Load keys from a specific namespace")
}
