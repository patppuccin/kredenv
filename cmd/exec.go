package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpExecCmd = "Executes a command with secrets injected into its environment"

var execCmd = &cobra.Command{
	Use:                "exec <command> [args...]",
	Short:              helpExecCmd,
	Long:               console.Banner(helpExecCmd),
	GroupID:            "env",
	SilenceUsage:       true,
	SilenceErrors:      true,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			console.Error("No command provided to execute")
			os.Exit(1)
		}

		path, err := kredsfile.Locate()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		kf, errs := kredsfile.Parse(path)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, e := range errs {
				errMsgs[i] = e.Error()
			}
			console.ErrorGroup("Failed to parse .kredsfile", errMsgs)
			os.Exit(1)
		}

		resolved := map[string]string{}
		var missingRequired []string

		for _, secret := range kf.Secrets {
			value, err := keyring.Get(secret.Key)
			if err != nil {
				if !secret.Optional {
					missingRequired = append(missingRequired, secret.Key)
				}
				continue
			}
			resolved[secret.Alias] = value
		}

		if len(missingRequired) > 0 {
			console.ErrorGroup("Missing required secrets", missingRequired)
			os.Exit(1)
		}

		runWith(args, resolved)
	},
}

func runWith(args []string, secrets map[string]string) {
	c := exec.Command(args[0], args[1:]...)
	c.Env = os.Environ()
	for k, v := range secrets {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		console.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	execCmd.Flags().SortFlags = false
}
