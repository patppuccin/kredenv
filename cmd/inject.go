package cmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpInjectCmd = "Emits export statements for the shell hook (internal use)"

var injectCmd = &cobra.Command{
	Use:           "inject",
	Short:         helpInjectCmd,
	GroupID:       "env",
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			console.Error("inject is for shell hook use only, use `kredenv load` instead")
			return
		}

		path, err := kredsfile.Locate()
		if err != nil || path == "" {
			return // silent exit, hook should not break the shell
		}

		kf, errs := kredsfile.Parse(path)
		if len(errs) > 0 {
			return // silent exit, hook should not break the shell
		}

		for _, secret := range kf.Secrets {
			value, err := keyring.Get(secret.Key)
			if err != nil {
				continue
			}
			fmt.Printf("export %s=%q\n", secret.Alias, value)
		}
	},
}

func init() {
	injectCmd.Flags().SortFlags = false
}
