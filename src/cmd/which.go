package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const (
	helpWhichCmd         = "Show paths to kredenv-managed files"
	helpWhichManifestCmd = "Path to the kredsfile.yaml in scope"
	helpWhichStoreCmd    = "Path to the encrypted secrets store"
	helpWhichCredsCmd    = "Path to the credentials file or keyring in use"
)

var whichCmd = &cobra.Command{
	Use:           "which",
	Short:         helpWhichCmd,
	Long:          banner(helpWhichCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var whichManifestCmd = &cobra.Command{
	Use:           "manifest",
	Short:         helpWhichManifestCmd,
	Long:          banner(helpWhichManifestCmd),
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		kp, err := spec.Locate()
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}
		if kp == "" {
			termactions.Log().Warn("No kredsfile.yaml found in scope")
			os.Exit(1)
		}
		fmt.Println(kp)
	},
}

var whichStoreCmd = &cobra.Command{
	Use:           "store",
	Short:         helpWhichStoreCmd,
	Long:          banner(helpWhichStoreCmd),
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		p, err := store.Path()
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(p)
	},
}

var whichCredsCmd = &cobra.Command{
	Use:           "creds",
	Short:         helpWhichCredsCmd,
	Long:          banner(helpWhichCredsCmd),
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		p, err := auth.Path()
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(p)
	},
}

func init() {
	whichCmd.Flags().SortFlags = false
	whichCmd.AddCommand(whichManifestCmd)
	whichCmd.AddCommand(whichStoreCmd)
	whichCmd.AddCommand(whichCredsCmd)
}
