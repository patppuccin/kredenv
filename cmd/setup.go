package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpSetupCmd = "Finds missing env vars and prompts to be stored in the keyring"

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         helpSetupCmd,
	Long:          console.Banner(helpSetupCmd),
	Args:          cobra.NoArgs,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
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

		stored := 0
		skipped := 0
		alreadySet := 0

		for _, secret := range kf.Secrets {
			if keyring.Exists(secret.Key) {
				alreadySet++
				continue
			}

			fmt.Printf("Enter value for %q: ", secret.Key)
			var value string
			fmt.Scan(&value)
			fmt.Println()

			if value == "" {
				console.Warn("Skipping: " + secret.Key)
				skipped++
				continue
			}

			if err := keyring.Set(secret.Key, value); err != nil {
				console.Error("Could not store " + secret.Key + ": " + err.Error())
				os.Exit(1)
			}

			stored++
			console.Info("Stored: " + secret.Key)
		}

		console.InfoGroup("Setup Complete", []string{
			"Stored:     " + strconv.Itoa(stored),
			"Skipped:    " + strconv.Itoa(skipped),
			"Already Set:" + strconv.Itoa(alreadySet),
		})
	},
}

func init() {
	setupCmd.Flags().SortFlags = false
}
