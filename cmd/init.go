package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/consts"
	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/spec"
	"github.com/patppuccin/kredenv/utils/store"
	"github.com/spf13/cobra"
)

const helpInitCmd = "Initialize a .kredsfile and optionally fill in missing secrets"

var (
	flagInitOverwrite bool
	flagInitFile      string
	flagInitNamespace string
	flagInitNoSetup   bool
)

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         helpInitCmd,
	Long:          console.Banner(helpInitCmd),
	GroupID:       "setup",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {

		if flagInitFile == "" {
			flagInitFile = ".kredsfile"
		}

		target, err := filepath.Abs(flagInitFile)
		if err != nil {
			console.Error("Could not resolve " + flagInitFile + ": " + err.Error())
			os.Exit(1)
		}

		if !strings.HasSuffix(filepath.Base(target), ".kredsfile") {
			console.Error("Kredsfile manifest must end in .kredsfile, got: " + filepath.Base(target))
			os.Exit(1)
		}

		fileExists := false
		if _, err := os.Stat(target); err == nil {
			fileExists = true
		}

		if fileExists && flagInitOverwrite {
			if err := os.WriteFile(target, []byte(spec.MinimalTemplate), 0644); err != nil {
				console.Error("Could not overwrite .kredsfile")
				os.Exit(1)
			}
			console.Success("Overwritten at " + target)
		} else if !fileExists {
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				console.Error("Could not create directories for: " + target)
				os.Exit(1)
			}
			if err := os.WriteFile(target, []byte(spec.MinimalTemplate), 0644); err != nil {
				console.Error("Could not write .kredsfile")
				os.Exit(1)
			}
			console.Success("Initialized at " + target)
		} else {
			console.Info("Using existing .kredsfile at " + target)
		}

		if flagInitNoSetup {
			return
		}

		if !isatty.IsTerminal(os.Stdin.Fd()) {
			return
		}

		kf, errs := spec.Parse(target)
		if len(errs) > 0 {
			return
		}

		if len(kf.Secrets) == 0 {
			console.Warn("No secrets found in .kredsfile")
		}

		password, err := auth.Retrieve()
		if err != nil {
			console.Warn("Could not open store: run '" + consts.AppName + " setup' first to set up")
			return
		}

		s, err := store.Open(password)
		if err != nil {
			console.Warn("Could not open store")
			return
		}
		defer s.Close()

		ns := flagInitNamespace
		if ns == "" {
			ns = kf.AutoloadNamespace
		}
		nsLabel := ""
		if ns != "" {
			nsLabel = " (namespace: " + ns + ")"
		}

		stored, skipped, alreadySet := []string{}, []string{}, []string{}

		for _, secret := range kf.Secrets {
			if ns != "" {
				if !strings.HasPrefix(secret.Key, ns+":") {
					continue
				}
			} else {
				if strings.Contains(secret.Key, ":") {
					continue
				}
			}

			if _, err := s.Get(secret.Key); err == nil {
				alreadySet = append(alreadySet, secret.Key)
				continue
			}

			label := secret.Key
			if secret.Optional {
				label += " (optional)"
			}

			value, err := console.PromptSecret("Enter value for " + label)
			if err != nil {
				console.Error("Could not read input")
				os.Exit(1)
			}

			if value == "" {
				skipped = append(skipped, secret.Key)
				continue
			}

			if err := s.Set(secret.Key, value); err != nil {
				console.Error("Could not store " + secret.Key)
				os.Exit(1)
			}

			stored = append(stored, secret.Key)
		}

		if len(stored) == 0 && len(skipped) == 0 {
			console.Success("All secrets already set" + nsLabel)
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

		console.InfoGroup(
			"Setup complete"+nsLabel,
			fmtSetupGroup(stored, "Stored"),
			fmtSetupGroup(skipped, "Skipped"),
			fmtSetupGroup(alreadySet, "Already set"),
		)

	},
}

func init() {
	initCmd.Flags().SortFlags = false
	initCmd.Flags().BoolVar(&flagInitOverwrite, "overwrite", false, "Overwrite existing .kredsfile")
	initCmd.Flags().BoolVar(&flagInitNoSetup, "no-setup", false, "Skip the interactive secret prompting after init")
	initCmd.Flags().StringVarP(&flagInitFile, "file", "f", ".kredsfile", "Path to the kredsfile (must end in .kredsfile)")
	initCmd.Flags().StringVarP(&flagInitNamespace, "namespace", "n", "", "Fill in secrets only for a specific namespace")
}
