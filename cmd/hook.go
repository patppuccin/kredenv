package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/hooks"
	"github.com/spf13/cobra"
)

const helpHookCmd = "Emits a shell hook script"

var usageInstructions = `
Supported shells: ` + hooks.Names() + `

Run the appropriate command once to install the hook:

 - ` + strings.Join(hooks.SetupCmds(), "\n - ")

var hookCmd = &cobra.Command{
	Use:           "hook <shell>",
	Short:         helpHookCmd,
	Long:          console.Banner(helpHookCmd) + usageInstructions,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			console.Error("Expected a shell name as the only argument")
			os.Exit(1)
		}

		sh, ok := hooks.Get(args[0])
		if !ok {
			console.Error("The shell '" + args[0] + "' is not supported")
			console.Info("Supported shells are " + hooks.Names())
			os.Exit(1)
		}
		// fmt.Println(sh.Hook)
		fmt.Print(strings.ReplaceAll(sh.Hook, "\r\n", "\n"))
	},
}

func init() {
	hookCmd.Flags().SortFlags = false
}
