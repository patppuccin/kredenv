package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpLoadCmd = "Loads the secrets from the .kredsfile in scope into the environment"

var (
	flagLoadNamespace string
)

var loadCmd = &cobra.Command{
	Use:           "load",
	Short:         helpLoadCmd,
	Long:          banner(helpLoadCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		if os.Getenv("__KREDENV_BIN") == "" {
			termactions.Log().Warn("Shell hook not detected, run '" + consts.AppName + " hook <shell>' to set up")
			return
		}

		loaded := os.Getenv("KREDENV_LOADED_VARS")
		if loaded == "" {
			termactions.Log().Warn("No secrets currently loaded")
			return
		}

		keys := strings.Split(loaded, ",")
		title := strconv.Itoa(len(keys)) + " secrets loaded"
		if flagLoadNamespace != "" {
			title += " (namespace: " + flagLoadNamespace + ")"
		}

		termactions.LogGroup().Info(title, keys...)
	},
}

func init() {
	loadCmd.Flags().SortFlags = false
	loadCmd.Flags().StringVarP(&flagLoadNamespace, "namespace", "n", "", "Load keys from a specific namespace")
}
