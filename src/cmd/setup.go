package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/helpers"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpSetupCmd = "Initialize " + consts.AppName + " on this machine"

var (
	flagSetupOverwrite bool
	flagSetupNuke      bool
)

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         helpSetupCmd,
	Long:          banner(helpSetupCmd),
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("The '" + cmd.CommandPath() + "' command does not accept any arguments")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if flagSetupOverwrite && flagSetupNuke {
			termactions.LogGroup().Error(
				"Cannot use --overwrite and --nuke at the same time",
				"To re-configure, run '"+cmd.CommandPath()+" --overwrite'",
				"To wipe everything, run '"+cmd.CommandPath()+" --nuke'",
			)
			os.Exit(1)
		}

		existingPassword, _ := auth.Retrieve()
		hasValidPassword := false
		if existingPassword != "" {
			if _, err := store.Open(existingPassword); err == nil {
				hasValidPassword = true
			}
		}
		hasStore := store.Exists()

		switch {
		case hasValidPassword && hasStore:
			switch {
			case flagSetupNuke:
				if err := nukeSetup(); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}
				fmt.Println()
				if err := freshSetup(""); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}

			case flagSetupOverwrite:
				if err := migrateSetup(existingPassword, false); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}

			default:
				termactions.LogGroup().Warn(
					consts.AppName+" is already configured and functional",
					"To re-configure, run '"+cmd.CommandPath()+" --overwrite'",
					"To wipe everything, run '"+cmd.CommandPath()+" --nuke'",
				)
				os.Exit(1)
			}

		case hasValidPassword && !hasStore:
			switch {
			case flagSetupNuke:
				if err := nukeSetup(); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}
				fmt.Println()
				if err := freshSetup(""); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}

			default:
				termactions.LogGroup().Warn(
					"No vault found, but master password exists",
					"A fresh empty vault will be created",
					"Run '"+consts.AppName+" set' or '"+consts.AppName+" init' to populate",
				)
				fmt.Println()
				if err := freshSetup(existingPassword); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}
			}

		case !hasValidPassword && hasStore:
			switch {
			case flagSetupNuke:
				if err := nukeSetup(); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}
				fmt.Println()
				if err := freshSetup(""); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}

			default:
				termactions.LogBlock().Info(
					"Vault found but master password is missing",
					"The original password is required to recover existing secrets",
					"To start fresh instead, run '"+cmd.CommandPath()+" --nuke'",
				)
				confirmed, err := termactions.Confirm().
					WithLabel("Do you remember the original password?").
					WithDefault(false).
					Render()
				if err != nil || !confirmed {
					termactions.Log().Warn("Setup aborted, to start fresh, run '" + cmd.CommandPath() + " --nuke'")
					os.Exit(1)
				}
				if err := migrateSetup("", true); err != nil {
					termactions.Log().Error(helpers.CapitalizeErrMsg(err))
					os.Exit(1)
				}
			}

		case !hasValidPassword && !hasStore:
			if flagSetupNuke || flagSetupOverwrite {
				termactions.LogGroup().Warn(
					"Nothing to delete or overwrite, "+consts.AppName+" is not configured yet",
					"Proceeding to a fresh setup",
				)
				fmt.Println()
			}
			if err := freshSetup(""); err != nil {
				termactions.Log().Error(helpers.CapitalizeErrMsg(err))
				os.Exit(1)
			}
		}
	},
}

func nukeSetup() error {
	confirmed, err := termactions.Confirm().
		WithLabel("This will delete all " + consts.AppName + " configuration and secrets. Are you sure?").
		WithDefault(false).
		Render()
	if err != nil || !confirmed {
		return fmt.Errorf("aborted — " + consts.AppName + " vault and configuration will not be deleted")
	}

	termactions.Log().Warn("Deletion confirmed — removing vault and configuration")

	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return auth.Delete()
	}

	if err := os.RemoveAll(rootDir); err != nil {
		return fmt.Errorf("could not remove directory %s", rootDir)
	}

	if err := auth.Delete(); err != nil {
		return err
	}

	termactions.Log().Success("Vault and configuration wiped successfully")
	return nil
}

func freshSetup(withPassword string) error {
	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return err
	}

	encFilePath := filepath.Join(rootDir, consts.EncFileName)

	if withPassword == "" {
		termactions.LogBlock().Info("Create master password",
			"Encrypts all secrets stored in the vault",
			"Stored at "+encFilePath,
			"Keep it safe — secrets cannot be recovered if forgotten",
		)

		withPassword, err = promptAndConfirmPassword("Enter master password", "Confirm master password")
		if err != nil {
			return err
		}

		termactions.Log().Success("Master password received — will be used to encrypt secrets")
	}

	if err := os.MkdirAll(rootDir, 0700); err != nil {
		return fmt.Errorf("could not create directory %s", rootDir)
	}

	if err := os.WriteFile(encFilePath, []byte{}, 0600); err != nil {
		return fmt.Errorf("could not create vault file %s", encFilePath)
	}

	if err := auth.Store(withPassword); err != nil {
		return fmt.Errorf("could not store master password")
	}

	termactions.Log().Success("Secrets vault created successfully")
	return nil
}

func migrateSetup(existingPassword string, recovery bool) error {
	if recovery {
		pwd, err := promptPassword("Enter existing master password")
		if err != nil {
			return err
		}
		existingPassword = pwd
	}

	s, err := store.Open(existingPassword)
	if err != nil {
		return fmt.Errorf("Incorrect password — retry or run '" + consts.AppName + " setup --nuke' to start fresh")
	}

	data, _ := s.List()

	if recovery {
		if len(data) != 0 {
			termactions.Log().Success("Existing master password verified")
			fmt.Println()
		} else {
			termactions.Log().Warn("No secrets found in store — proceeding to a fresh setup")
			fmt.Println()
		}
	}

	label := "Change master password"
	var bodyLines []string
	if len(data) != 0 {
		bodyLines = []string{
			"All existing secrets will be re-encrypted with the new password",
			"Keep it safe — secrets cannot be recovered if forgotten",
		}
	} else {
		label = "Create master password"
		bodyLines = []string{
			"Keep it safe — secrets cannot be recovered if forgotten",
		}
	}

	termactions.LogBlock().Info(label, bodyLines...)

	newPassword, err := promptAndConfirmPassword("Enter new master password", "Confirm new master password")
	if err != nil {
		return err
	}

	termactions.Log().Success("New master password received — will be used to re-encrypt secrets")

	if err := store.Migrate(existingPassword, newPassword); err != nil {
		return err
	}

	if err := auth.Store(newPassword); err != nil {
		return fmt.Errorf("Could not store master password")
	}

	termactions.Log().Success("Vault re-encrypted successfully")
	return nil
}

func promptPassword(label string) (string, error) {
	pwd, err := termactions.Secret().
		WithLabel(label).
		WithValidator(func(s string) (string, bool) {
			if strings.TrimSpace(s) == "" {
				return "master password cannot be empty", false
			}
			return "", true
		}).
		Render()

	if err != nil {
		if err == termactions.ErrInterrupted {
			return "", fmt.Errorf("password prompt aborted by user")
		}
		return "", err
	}

	return pwd, nil
}

func promptAndConfirmPassword(prompt, confirmPrompt string) (string, error) {
	pwd, err := promptPassword(prompt)
	if err != nil {
		return "", err
	}

	confirm, err := promptPassword(confirmPrompt)
	if err != nil {
		return "", err
	}

	if pwd != confirm {
		return "", errors.New("passwords do not match")
	}

	return pwd, nil
}

func init() {
	setupCmd.Flags().SortFlags = false
	setupCmd.Flags().BoolVar(&flagSetupOverwrite, "overwrite", false, "Re-encrypt vault secrets with a new master password")
	setupCmd.Flags().BoolVar(&flagSetupNuke, "nuke", false, "Wipe all "+consts.AppName+" configuration and secrets")
}
