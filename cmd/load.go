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

const helpLoadCmd = "Resolves the .kredsfile and injects it into the environment"

var (
	flagLoadOnce bool
)

var loadCmd = &cobra.Command{
	Use:           "load [-- <command>]",
	Short:         helpLoadCmd,
	Long:          console.Banner(helpLoadCmd),
	Args:          cobra.ArbitraryArgs,
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := kredsfile.Locate()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
		if path == "" {
			console.Warn("No .kredsfile found")
			os.Exit(1)
		}

		kf, errs := kredsfile.Parse(path)
		if len(errs) > 0 {
			if len(errs) == 1 {
				console.Error(errs[0].Error())
				os.Exit(1)
			}
			errMsgs := make([]string, len(errs))
			for i, e := range errs {
				errMsgs[i] = e.Error()
			}
			console.ErrorGroup("Failed to parse "+path, errMsgs)
			os.Exit(1)
		}

		resolved := map[string]string{}
		var missingRequired []string

		for _, secret := range kf.Secrets {
			value, err := keyring.Get(secret.Key)
			if err != nil {
				if !secret.Optional {
					missingRequired = append(missingRequired, secret.Key)
				} else {
					console.Warn("Missing optional secret: " + secret.Alias)
				}
				continue
			}
			resolved[secret.Alias] = value
		}

		if len(missingRequired) > 0 {
			console.ErrorGroup("Missing required secrets", missingRequired)
			os.Exit(1)
		}

		if flagLoadOnce || len(args) > 0 {
			if len(args) == 0 {
				console.Error("--once requires a command: kredenv load --once -- <command>")
				os.Exit(1)
			}
			runWith(args, resolved)
			return
		}

		for k, v := range resolved {
			fmt.Printf("export %s=%q\n", k, v)
		}
	},
}

// runWith execs a command with resolved secrets injected as env vars
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
	loadCmd.Flags().SortFlags = false
	loadCmd.Flags().BoolVar(&flagLoadOnce, "once", false, "Inject secrets for this command only")
}
