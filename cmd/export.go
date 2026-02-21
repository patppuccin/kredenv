package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/crypto"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpExportCmd = "Exports secrets from the keyring to stdout or a file"

var (
	flagExportAll     bool
	flagExportFormat  string
	flagExportOutput  string
	flagExportEncrypt bool
)

var exportCmd = &cobra.Command{
	Use:           "export",
	Short:         helpExportCmd,
	Long:          console.Banner(helpExportCmd),
	Args:          cobra.NoArgs,
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		// collect keys to export
		secrets, errs := collectSecrets()
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, err := range errs {
				errMsgs[i] = err.Error()
			}
			console.ErrorGroup("Failed to collect secrets", errMsgs)
			os.Exit(1)
		}

		if len(secrets) == 0 {
			console.Warn("no secrets to export")
			return
		}

		// format output
		output, err := formatSecrets(secrets, flagExportFormat)
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		// encrypt if requested
		if flagExportEncrypt {
			output, err = encrypt(output)
			if err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
		}

		// write to file or stdout
		if flagExportOutput == "" {
			fmt.Print(output)
			return
		}

		if err := os.WriteFile(flagExportOutput, []byte(output), 0600); err != nil {
			console.Error("could not write file: " + err.Error())
			os.Exit(1)
		}

		console.Info("exported to " + flagExportOutput)
	},
}

func collectSecrets() (map[string]string, []error) {
	secrets := map[string]string{}

	if flagExportAll {
		keys, err := keyring.List()
		if err != nil {
			return nil, []error{err}
		}
		for _, key := range keys {
			value, err := keyring.Get(key)
			if err != nil {
				continue
			}
			secrets[key] = value
		}
		return secrets, nil
	}

	// default — use active .kredsfile
	path, err := kredsfile.Locate()
	if err != nil {
		return nil, []error{err}
	}
	if path == "" {
		return nil, []error{fmt.Errorf("No .kredsfile found")}
	}

	kf, errs := kredsfile.Parse(path)
	if len(errs) > 0 {
		return nil, errs
	}

	var collectErrs []error
	for _, secret := range kf.Secrets {
		value, err := keyring.Get(secret.Key)
		if err != nil {
			if !secret.Optional {
				collectErrs = append(collectErrs, fmt.Errorf("Missing required key: %s", secret.Key))
			}
			continue
		}
		secrets[secret.Alias] = value
	}

	if len(collectErrs) > 0 {
		return nil, collectErrs
	}

	return secrets, nil
}

// formatSecrets formats the secrets map into the requested format
func formatSecrets(secrets map[string]string, format string) (string, error) {
	switch format {
	case "env":
		var sb strings.Builder
		for k, v := range secrets {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
		return sb.String(), nil

	case "json":
		data, err := json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data) + "\n", nil

	case "yaml":
		var sb strings.Builder
		for k, v := range secrets {
			sb.WriteString(fmt.Sprintf("%s: %q\n", k, v))
		}
		return sb.String(), nil

	case "toml":
		var sb strings.Builder
		for k, v := range secrets {
			sb.WriteString(fmt.Sprintf("%s = %q\n", k, v))
		}
		return sb.String(), nil

	default:
		return "", fmt.Errorf("unsupported format %q, supported: env, json, yaml, toml", format)
	}
}

// encrypt prompts for a password and encrypts the output
func encrypt(data string) (string, error) {
	fmt.Print("enter encryption password: ")
	var password string
	fmt.Scan(&password)
	fmt.Println()

	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	encrypted, err := crypto.Encrypt([]byte(data), password)
	if err != nil {
		return "", fmt.Errorf("could not encrypt: %w", err)
	}

	return encrypted, nil
}

func init() {
	exportCmd.Flags().SortFlags = false
	exportCmd.Flags().BoolVar(&flagExportAll, "all", false, "Export all keys in the keyring")
	exportCmd.Flags().StringVarP(&flagExportFormat, "format", "f", "env", "Export format (env, json, yaml, toml)")
	exportCmd.Flags().StringVarP(&flagExportOutput, "output", "o", "", "Output file path (defaults to stdout)")
	exportCmd.Flags().BoolVar(&flagExportEncrypt, "encrypt", false, "Encrypt the exported file")
}
