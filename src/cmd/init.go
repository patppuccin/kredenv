package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpInitCmd = "Initialize a kredsfile.yaml and optionally fill in missing secrets"

var (
	flagInitOverwrite bool
	flagInitFile      string
	flagInitNamespace string
	flagInitNoSetup   bool
)

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         helpInitCmd,
	Long:          banner(helpInitCmd),
	GroupID:       "setup",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {

		target, err := filepath.Abs(flagInitFile)
		if err != nil {
			termactions.Log().Error("Could not resolve " + flagInitFile + ": " + err.Error())
			os.Exit(1)
		}

		if filepath.Base(target) != "kredsfile.yaml" {
			termactions.Log().Error("Kredsfile manifest must be named kredsfile.yaml, got: " + filepath.Base(target))
			os.Exit(1)
		}

		fileExists := false
		if _, err := os.Stat(target); err == nil {
			fileExists = true
		}

		if fileExists && flagInitOverwrite {
			if err := os.WriteFile(target, []byte(spec.MinimalTemplate), 0644); err != nil {
				termactions.Log().Error("Could not overwrite kredsfile.yaml")
				os.Exit(1)
			}
			termactions.Log().Success("Overwritten at " + target)
		} else if !fileExists {
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				termactions.Log().Error("Could not create directories for: " + target)
				os.Exit(1)
			}
			if err := os.WriteFile(target, []byte(spec.MinimalTemplate), 0644); err != nil {
				termactions.Log().Error("Could not write manifest template to " + target)
				os.Exit(1)
			}
			termactions.Log().Success("Initialized kredsfile manifest at " + target)
		} else {
			termactions.Log().Info("Using existing manifest at " + target)
		}

		if flagInitNoSetup {
			return
		}

		if !isatty.IsTerminal(os.Stdin.Fd()) {
			return
		}

		kf, errs := spec.Parse(target)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, e := range errs {
				errMsgs[i] = e.Error()
			}
			termactions.LogGroup().Error("Failed to parse "+target, errMsgs...)
			os.Exit(1)
		}

		if len(kf.Secrets) == 0 {
			termactions.Log().Warn("No secrets declared in kredsfile.yaml")
			return
		}

		authPasswd, err := auth.Retrieve()
		if err != nil {
			termactions.Log().Warn("Could not open auth store: run '" + consts.AppName + " setup' first to set up")
			return
		}

		s, err := store.Open(authPasswd)
		if err != nil {
			termactions.Log().Warn("Could not open the secrets store")
			return
		}
		defer s.Close()

		ns := flagInitNamespace
		nsLabel := ""
		if ns != "" {
			nsLabel = " (namespace: " + ns + ")"
		}

		stored, skipped, alreadySet := []string{}, []string{}, []string{}

		for _, secret := range kf.Secrets {
			if ns != "" && secret.Namespace != ns {
				continue
			}

			if _, err := s.Get(secret.VaultKey()); err == nil {
				alreadySet = append(alreadySet, secret.VaultKey())
				continue
			}

			value, err := termactions.Secret().WithLabel("Enter value for " + secret.VaultKey()).Render()
			if err != nil {
				if err == termactions.ErrInterrupted {
					termactions.Log().Warn("Initialization aborted by user")
					os.Exit(1)
				}
				termactions.Log().Error("Could not read input")
				os.Exit(1)
			}

			if value == "" {
				skipped = append(skipped, secret.VaultKey())
				continue
			}

			if err := s.Set(secret.VaultKey(), value); err != nil {
				termactions.Log().Error("Could not store " + secret.VaultKey())
				os.Exit(1)
			}

			stored = append(stored, secret.VaultKey())
			termactions.Log().Success("Stored " + secret.VaultKey())
		}

		if len(stored) == 0 && len(skipped) == 0 {
			termactions.Log().Success("All secrets already set" + nsLabel)
			return
		}

		fmtSetupGroup := func(keys []string, label string) string {
			count := len(keys)
			noun := "secrets"
			if count == 1 {
				noun = "secret"
			}
			if count == 0 {
				return fmt.Sprintf("%s 0 secrets", label)
			}
			return fmt.Sprintf("%s %d %s: %s", label, count, noun, strings.Join(keys, ", "))
		}

		fmt.Println()
		termactions.LogGroup().Info(
			"Setup complete"+nsLabel,
			fmtSetupGroup(stored, "Stored"),
			fmtSetupGroup(skipped, "Skipped"),
			fmtSetupGroup(alreadySet, "Already set"),
		)

	},
}

func init() {
	initCmd.Flags().SortFlags = false
	initCmd.Flags().BoolVar(&flagInitOverwrite, "overwrite", false, "Overwrite existing kredsfile.yaml")
	initCmd.Flags().BoolVar(&flagInitNoSetup, "no-setup", false, "Skip the interactive secret prompting after init")
	initCmd.Flags().StringVarP(&flagInitFile, "file", "f", "kredsfile.yaml", "Path to the kredenv manifest file")
	initCmd.Flags().StringVarP(&flagInitNamespace, "namespace", "n", "", "Fill in secrets only for a specific namespace")
}
