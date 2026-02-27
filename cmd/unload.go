package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/consts"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/spf13/cobra"
)

const helpUnloadCmd = "Unloads the " + consts.AppName + " secrets from the environment"

var unloadCmd = &cobra.Command{
	Use:           "unload",
	Short:         helpUnloadCmd,
	Long:          console.Banner(helpUnloadCmd),
	Args:          cobra.NoArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("__KREDENV_BIN") == "" {
			console.Warn("Shell hook not detected, run '" + consts.AppName + " hook <shell>' to set up")
			return
		}

		if os.Getenv("KREDENV_LOADED_VARS") != "" {
			console.Error("Unload failed, secrets still present in session")
			os.Exit(1)
		}

		console.Success("Secrets unloaded from session")
	},
}

func init() {
	unloadCmd.Flags().SortFlags = false
}
