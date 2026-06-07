package cmd

import (
	"os"

	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpUnloadCmd = "Unloads the " + consts.AppName + " secrets from the environment"

var unloadCmd = &cobra.Command{
	Use:           "unload",
	Short:         helpUnloadCmd,
	Long:          banner(helpUnloadCmd),
	Args:          cobra.NoArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("__KREDENV_BIN") == "" {
			termactions.Log().Error("Shell hook not detected, run '" + consts.AppName + " hook <shell>' to set up")
			return
		}

		if os.Getenv("KREDENV_LOADED_VARS") != "" {
			termactions.Log().Error("Unload failed, secrets still present in session")
			os.Exit(1)
		}

		termactions.Log().Success("Secrets unloaded from session")
	},
}

func init() {
	unloadCmd.Flags().SortFlags = false
}
