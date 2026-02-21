package cmd

import (
	"fmt"
	"os"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/shells"
	"github.com/spf13/cobra"
)

const helpHookCmd = "Emits a shell hook script (supported: bash, zsh, fish, powershell, nushell)"

var hookCmd = &cobra.Command{
	Use:           "hook <shell>",
	Short:         helpHookCmd,
	Long:          console.Banner(helpHookCmd),
	Args:          cobra.ExactArgs(1),
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		shell, ok := shells.Get(args[0])
		if !ok {
			console.Error("Unsupported shell: " + args[0] + ", (supported: " + shells.Names() + ")")
			os.Exit(1)
		}
		fmt.Println(shell.Hook())
	},
}

func init() {
	hookCmd.Flags().SortFlags = false
}
