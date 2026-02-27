package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpInitCmd = "Initializes a minimal .kredsfile in the current directory"

var (
	flagInitOverwrite bool
	flagInitFile      string
)

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         helpInitCmd,
	Long:          console.Banner(helpInitCmd),
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		target := flagInitFile

		if !filepath.IsAbs(target) {
			cwd, err := os.Getwd()
			if err != nil {
				console.Error("Could not determine current directory")
				os.Exit(1)
			}
			target = filepath.Join(cwd, target)
		}

		if !strings.HasSuffix(filepath.Base(target), ".kredsfile") {
			console.Error("Kredsfile manifest must end in .kredsfile, got: " + filepath.Base(target))
			os.Exit(1)
		}

		if _, err := os.Stat(target); err == nil && !flagInitOverwrite {
			console.Error("File already exists: " + target)
			os.Exit(1)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			console.Error("Could not create directories for: " + target)
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
	initCmd.Flags().StringVarP(&flagInitFile, "file", "f", ".kredsfile", "Path to the kredsfile (must end in .kredsfile)")
}
