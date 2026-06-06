package cmd

import (
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpWhichCmd = "Prints the path to the .kredsfile that will be used"

var whichCmd = &cobra.Command{
	Use:           "which",
	Short:         helpWhichCmd,
	Long:          banner(helpWhichCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
		kp, err := spec.Locate()
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}
		if kp == "" {
			termactions.Log().Warn("Could not locate .kredsfile")
			os.Exit(1)
		}
		termactions.Log().Info("Located kredsfile at: " + kp)
	},
}

func init() {
	whichCmd.Flags().SortFlags = false
}
