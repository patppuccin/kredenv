package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/shells"
	"github.com/spf13/cobra"
)

const helpHookCmd = "Emits a shell hook script"

var usageInstructions = `
Supported shells: ` + shells.Names() + `

Run the appropriate command once to install the hook:

 - ` + strings.Join(shells.SetupCmds(), "\n - ")

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

		sh, ok := shells.Get(args[0])
		if !ok {
			console.Error("Unsupported shell: " + args[0])
			console.Info("Supported: " + shells.Names())
			os.Exit(1)
		}
		fmt.Println(sh.Hook)
	},
}

func init() {
	hookCmd.Flags().SortFlags = false
}
