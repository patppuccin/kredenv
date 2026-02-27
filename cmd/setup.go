package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpSetupCmd = "Finds missing env vars and prompts to be stored in the keyring"

var (
	flagSetupNamespace string
)

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         helpSetupCmd,
	Long:          console.Banner(helpSetupCmd),
	Args:          cobra.NoArgs,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			console.Error("setup requires a terminal to prompt for values")
			os.Exit(1)
		}

		path, err := kredsfile.Locate()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
		if path == "" {
			console.Warn("no .kredsfile found")
			os.Exit(1)
		}

		kf, errs := kredsfile.Parse(path)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, err := range errs {
				errMsgs[i] = err.Error()
			}
			console.ErrorGroup("Failed to parse "+path, errMsgs)
			os.Exit(1)
		}

		stored, skipped, alreadySet := []string{}, []string{}, []string{}

		reader := bufio.NewReader(os.Stdin)

		ns := flagSetupNamespace
		if ns == "" {
			ns = kf.AutoloadNamespace
		}
		nsLabel := ""
		if ns != "" {
			nsLabel = " (namespace: " + ns + ")"
		}

		for _, secret := range kf.Secrets {
			if ns != "" {
				if !strings.HasPrefix(secret.Key, ns+":") {
					continue
				}
			} else {
				if strings.Contains(secret.Key, ":") {
					continue
				}
			}

			if keyring.Exists(secret.Key) {
				alreadySet = append(alreadySet, secret.Key)
				continue
			}

			fmt.Printf("\nEnter value for %q: ", secret.Key)
			input, err := reader.ReadString('\n')
			if err != nil {
				console.Error("Could not read input")
				os.Exit(1)
			}
			value := strings.TrimSpace(input)

			if value == "" {
				skipped = append(skipped, secret.Key)
				continue
			}

			if err := keyring.Set(secret.Key, value); err != nil {
				console.Error("Could not store " + secret.Key)
				os.Exit(1)
			}

			stored = append(stored, secret.Key)
		}

		if len(stored) == 0 && len(skipped) == 0 {
			console.Success("All secrets already set in keyring" + nsLabel)
			return
		}

		joinOrNone := func(s []string) string {
			if len(s) == 0 {
				return "none"
			}
			return strings.Join(s, ", ")
		}

		console.InfoGroup("Setup complete"+nsLabel, []string{
			"Stored      :" + joinOrNone(stored),
			"Skipped     :" + joinOrNone(skipped),
			"Already set :" + joinOrNone(alreadySet),
		})
	},
}

func init() {
	setupCmd.Flags().SortFlags = false
	setupCmd.Flags().StringVarP(&flagSetupNamespace, "namespace", "n", "", "Setup keys only for a specific namespace")
}
