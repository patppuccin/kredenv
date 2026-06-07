package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patppuccin/kredenv/src/hooks"
	"github.com/patppuccin/termactions"
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
	Long:          banner(helpHookCmd) + usageInstructions,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			termactions.Log().Error("Expected a shell name as the only argument")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		sh, ok := hooks.Get(args[0])
		if !ok {
			termactions.Log().Error("The shell '" + args[0] + "' is not supported")
			termactions.Log().Info("Supported shells are " + hooks.Names())
			os.Exit(1)
		}
		fmt.Print(strings.ReplaceAll(sh.Hook, "\r\n", "\n"))
	},
}

func init() {
	hookCmd.Flags().SortFlags = false
}
