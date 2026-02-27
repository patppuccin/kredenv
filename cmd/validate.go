package cmd

import (
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpValidateCmd = "Validates .kredsfile syntax"

var validateCmd = &cobra.Command{
	Use:           "validate [file]",
	Short:         helpValidateCmd,
	Long:          console.Banner(helpValidateCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		var kp string
		var err error

		switch len(args) {
		case 0:
			kp, err = kredsfile.Locate()
			if err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
			if kp == "" {
				console.Error("No .kredsfile found")
				os.Exit(1)
			}
		case 1:
			kp = args[0]
		default:
			console.Error("Expected at most one argument, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		if _, err := os.Stat(kp); os.IsNotExist(err) {
			console.Error("No file found at " + kp)
			os.Exit(1)
		}

		kf, errs := kredsfile.Parse(kp)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, err := range errs {
				errMsgs[i] = err.Error()
			}
			console.ErrorGroup("Failed to parse .kredsfile", errMsgs)
			os.Exit(1)
		}

		if len(kf.Secrets) == 0 {
			console.Warn("No secrets declared in " + kp)
			os.Exit(0)
		}

		console.Success("Valid .kredsfile with " + strconv.Itoa(len(kf.Secrets)) + " secrets at " + kp)

	},
}

func init() {
	validateCmd.Flags().SortFlags = false
}
