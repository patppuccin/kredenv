package cmd

import (
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpWhichCmd = "Prints the path to the .kredsfile that will be used"

var whichCmd = &cobra.Command{
	Use:           "which",
	Short:         helpWhichCmd,
	Long:          console.Banner(helpWhichCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
		kp, err := kredsfile.Locate()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
		if kp == "" {
			console.Warn("Could not locate .kredsfile")
			os.Exit(1)
		}
		console.Info("Located kredsfile at: " + kp)
	},
}

func init() {
	whichCmd.Flags().SortFlags = false
}
