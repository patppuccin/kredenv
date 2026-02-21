package cmd

import (
	"fmt"
	"os"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpUnloadCmd = "Unloads the .kredsfile from the environment"

var unloadCmd = &cobra.Command{
	Use:           "unload",
	Short:         helpUnloadCmd,
	Long:          console.Banner(helpUnloadCmd),
	Args:          cobra.NoArgs,
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

		// TODO: emit unset statements for shell to eval
		for _, secret := range kf.Secrets {
			fmt.Printf("unset %s\n", secret.Alias)
		}
	},
}

func init() {
	unloadCmd.Flags().SortFlags = false
}
