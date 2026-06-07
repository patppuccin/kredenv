package cmd

import (
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpValidateCmd = "Validates kredsfile.yaml syntax"

var validateCmd = &cobra.Command{
	Use:           "validate [file]",
	Short:         helpValidateCmd,
	Long:          banner(helpValidateCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		var kp string
		var err error

		switch len(args) {
		case 0:
			kp, err = spec.Locate()
			if err != nil {
				termactions.Log().Error(err.Error())
				os.Exit(1)
			}
			if kp == "" {
				termactions.Log().Error("No kredsfile.yaml found")
				os.Exit(1)
			}
		case 1:
			kp = args[0]
		default:
			termactions.Log().Error("Expected at most one argument, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		if _, err := os.Stat(kp); os.IsNotExist(err) {
			termactions.Log().Error("No file found at " + kp)
			os.Exit(1)
		}

		kf, errs := spec.Parse(kp)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, err := range errs {
				errMsgs[i] = err.Error()
			}
			termactions.LogGroup().Error("Failed to parse kredsfile.yaml", errMsgs...)
			os.Exit(1)
		}

		if len(kf.Secrets) == 0 {
			termactions.Log().Warn("No secrets declared in " + kp)
			return
		}

		termactions.Log().Success("Valid kredsfile.yaml with " + strconv.Itoa(len(kf.Secrets)) + " secrets at " + kp)

	},
}

func init() {
	validateCmd.Flags().SortFlags = false
}
