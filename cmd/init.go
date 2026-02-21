package cmd

import (
	"os"
	"path/filepath"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpInitCmd = "Initializes a minimal .kredsfile in the current directory"

var (
	flagInitOverwrite bool
)

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         helpInitCmd,
	Long:          console.Banner(helpInitCmd),
	Args:          cobra.NoArgs,
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			console.Error("Could not determine current directory")
			os.Exit(1)
		}
		target := filepath.Join(cwd, ".kredsfile")
		if _, err := os.Stat(target); err == nil && !flagInitOverwrite {
			console.Error("File already exists")
			os.Exit(1)
		}
		if err := os.WriteFile(target, []byte(kredsfile.MinimalTemplate), 0644); err != nil {
			console.Error("Could not write .kredsfile")
			os.Exit(1)
		}
		console.Success("Initialized at " + target)
	},
}

func init() {
	initCmd.Flags().SortFlags = false
	initCmd.Flags().BoolVar(&flagInitOverwrite, "overwrite", false, "Overwrite existing .kredsfile")
}
