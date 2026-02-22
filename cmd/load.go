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

var loadCmd = &cobra.Command{
	Use:           "load [-- <command>]",
	Short:         helpLoadCmd,
	Long:          console.Banner(helpLoadCmd),
	Args:          cobra.ArbitraryArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("__KREDENV_BIN") == "" {
			console.Warn("Shell hook not detected, run '" + consts.AppName + " hook <shell>' to set up")
			return
		}

		loaded := os.Getenv("KREDENV_LOADED")
		if loaded == "" {
			console.Warn("No secrets currently loaded")
			return
		}

		keys := strings.Split(loaded, ",")
		console.InfoGroup(
			strconv.Itoa(len(keys))+" secrets loaded",
			keys,
		)
	},
}

func init() {
	loadCmd.Flags().SortFlags = false
}
