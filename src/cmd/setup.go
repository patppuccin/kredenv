package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/console"
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/helpers"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/spf13/cobra"
)

const helpSetupCmd = "Initialize kredenv on this machine"

var (
	flagSetupOverwrite bool
	flagSetupNuke      bool
)

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         helpSetupCmd,
	Long:          console.Banner(helpSetupCmd),
	GroupID:       "setup",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("The '" + cmd.CommandPath() + "' command does not accept any arguments")
			os.Exit(1)
		}

		if flagSetupOverwrite && flagSetupNuke {
			console.ErrorGroup(
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
					console.Error(err.Error())
					os.Exit(1)
				}
				if err := freshSetup(""); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}

			case flagSetupOverwrite:
				if err := migrateSetup(existingPassword, false); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}

			default:
				console.WarnGroup(
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
					console.Error(err.Error())
					os.Exit(1)
				}
				if err := freshSetup(""); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}

			default:
				console.Info("Partial configuration detected, recovering...")
				if err := freshSetup(existingPassword); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}
			}

		case !hasValidPassword && hasStore:
			switch {
			case flagSetupNuke:
				if err := nukeSetup(); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}
				if err := freshSetup(""); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}

			case flagSetupOverwrite:
				console.InfoBlock(
					"Secrets store found but login credentials are missing",
					"The original password is required to recover existing secrets",
					"To start fresh instead, run '"+cmd.CommandPath()+" --nuke'",
				)
				if !console.PromptConfirm("Do you remember the original password?") {
					console.WarnGroup("Aborting setup",
						"To start fresh, run '"+cmd.CommandPath()+" --nuke'",
					)
					os.Exit(1)
				}
				if err := migrateSetup("", true); err != nil {
					console.Error(err.Error())
					os.Exit(1)
				}

			default:
				console.ErrorGroup(
					"Secrets store found but login credentials are missing",
					"To recover, run '"+cmd.CommandPath()+" --overwrite'",
					"To wipe everything and start fresh, run '"+cmd.CommandPath()+" --nuke'",
				)
				os.Exit(1)
			}

		case !hasValidPassword && !hasStore:
			if flagSetupNuke || flagSetupOverwrite {
				console.WarnGroup("Nothing to delete or overwrite, kredenv is not configured yet",
					"Proceeding to a fresh setup",
				)
			}
			if err := freshSetup(""); err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
		}
	},
}

func nukeSetup() error {
	confirmed := console.PromptConfirm("This will delete all kredenv configuration and secrets. Are you sure?")
	if !confirmed {
		return fmt.Errorf("aborted while deleting kredenv configuration")
	}

	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return auth.Delete() // ensure login credential is deleted even if directory does not exist
	}

	if err := os.RemoveAll(rootDir); err != nil {
		return fmt.Errorf("could not remove directory %s", rootDir)
	}

	if err := auth.Delete(); err != nil {
		return err
	}

	console.Success("Configuration and secrets wiped successfully")
	return nil
}

func freshSetup(withPassword string) error {

	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return err
	}

	encFilePath := filepath.Join(rootDir, consts.EncFileName)

	if withPassword == "" {
		console.InfoBlock("Create master password",
			"Encrypts all secrets and credentials",
			"Stored at "+encFilePath,
			"Keep it safe. Secrets cannot be recovered if forgotten",
		)

		withPassword, err = console.PromptAndConfirmPassword("Enter master password", "Confirm master password")
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(rootDir, 0700); err != nil {
		return fmt.Errorf("could not create directory %s", rootDir)
	}

	if err := os.WriteFile(encFilePath, []byte{}, 0600); err != nil {
		return fmt.Errorf("could not create store file %s", encFilePath)
	}

	if err := auth.Store(withPassword); err != nil {
		return fmt.Errorf("could not store master password")
	}

	console.Success("Configured successfully")
	return nil
}

func migrateSetup(existingPassword string, recovery bool) error {
	if recovery {
		var err error
		existingPassword, err = console.PromptSecret("Enter existing/old password")
		if err != nil || existingPassword == "" {
			return fmt.Errorf("password cannot be empty")
		}
		if _, err := store.Open(existingPassword); err != nil {
			return fmt.Errorf("incorrect password")
		}
	}

	console.InfoBlock("Change master password",
		"All existing secrets will be re-encrypted with the new password",
		"Keep it safe. Secrets cannot be recovered if forgotten",
	)

	newPassword, err := console.PromptAndConfirmPassword("Enter new master password", "Confirm new master password")
	if err != nil {
		return err
	}

	if err := store.Migrate(existingPassword, newPassword); err != nil {
		return err
	}

	if err := auth.Store(newPassword); err != nil {
		return fmt.Errorf("could not store master password")
	}

	console.Success("Re-configured successfully")
	return nil
}

func init() {
	setupCmd.Flags().SortFlags = false
	setupCmd.Flags().BoolVar(&flagSetupOverwrite, "overwrite", false, "Re-encrypt secrets with a new master password")
	setupCmd.Flags().BoolVar(&flagSetupNuke, "nuke", false, "Wipe all kredenv configuration and secrets")
}
