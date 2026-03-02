package cmd

import (
	"os"
	"path/filepath"

	"github.com/patppuccin/kredenv/consts"
	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/helpers"
	"github.com/patppuccin/kredenv/utils/store"
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
			console.ErrorGroup("Cannot use --overwrite and --nuke at the same time",
				[]string{
					"If you want to re-configure, run '" + cmd.CommandPath() + " --overwrite'",
					"If you want to wipe everything (including secrets), run '" + cmd.CommandPath() + " --nuke'",
				},
			)
			os.Exit(1)
		}

		rootDir, err := helpers.GetRootDir()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		_, rootErr := os.Stat(rootDir)
		rootExists := rootErr == nil

		if rootExists {
			if flagSetupNuke {
				confirmed, err := console.PromptConfirm("This will delete all kredenv configuration and secrets. Are you sure?")
				if err != nil || !confirmed {
					console.Error("Aborted")
					os.Exit(1)
				}
				if err := os.RemoveAll(rootDir); err != nil {
					console.Error("Could not remove " + rootDir)
					os.Exit(1)
				}
				console.Success("Previous kredenv configuration wiped successfully")
				console.Info("Run '" + cmd.CommandPath() + "' again to configure with a new password")
				return
			}

			existingPassword, err := auth.Retrieve()

			if existingPassword != "" && !flagSetupOverwrite {
				console.WarnGroup(
					consts.AppName+" is already configured and functional",
					[]string{
						"If you want to re-configure, run '" + cmd.CommandPath() + " --overwrite'",
						"If you want to wipe everything (including secrets), run '" + cmd.CommandPath() + " --nuke'",
					},
				)
				os.Exit(1)
			}

			if err != nil {
				if !flagSetupOverwrite {
					console.ErrorGroup(consts.AppName+" is already configured but unable to authenticate",
						[]string{
							"If previous password is known and you want to re-configure, run '" + cmd.CommandPath() + " --overwrite'",
							"If you want to wipe everything (including secrets), run '" + cmd.CommandPath() + " --nuke'",
						},
					)
					os.Exit(1)
				}

				knowsPwd, promptErr := console.PromptConfirm("Could not retrieve existing password. Do you remember it?")
				if promptErr != nil || !knowsPwd {
					console.Error("Use --nuke to wipe everything and start fresh")
					os.Exit(1)
				}

				existingPassword, promptErr = console.PromptSecret("Enter existing password")
				if promptErr != nil || existingPassword == "" {
					console.Error("Password cannot be empty")
					os.Exit(1)
				}
			}

			newPassword, err := console.PromptAndConfirmPassword("Enter new master password", "Confirm new master password")
			if err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}

			if err := store.Migrate(existingPassword, newPassword); err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}

			if err := auth.Store(newPassword); err != nil {
				console.Error("Could not store master password")
				os.Exit(1)
			}

			console.Success("kredenv re-configured successfully")
			return
		}

		// If crossed this point, kredenv is to be freshly configured

		password, err := console.PromptAndConfirmPassword("Enter master password", "Confirm master password")
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		if err := os.MkdirAll(rootDir, 0700); err != nil {
			console.Error("Could not create " + rootDir)
			os.Exit(1)
		}

		encFilePath := filepath.Join(rootDir, consts.EncFileName)
		if err := os.WriteFile(encFilePath, []byte{}, 0600); err != nil {
			console.Error("Could not create store file")
			os.Exit(1)
		}

		if err := auth.Store(password); err != nil {
			console.Error("Could not store master password")
			os.Exit(1)
		}

		console.Success("kredenv configured successfully")

	},
}

func init() {
	setupCmd.Flags().SortFlags = false
	setupCmd.Flags().BoolVar(&flagSetupOverwrite, "overwrite", false, "Re-encrypt secrets with a new master password")
	setupCmd.Flags().BoolVar(&flagSetupNuke, "nuke", false, "Wipe all kredenv configuration and secrets")
}
