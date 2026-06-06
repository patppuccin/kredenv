package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/helpers"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

const helpImportCmd = "Import secrets from a file into the kredenv store"

var (
	flagImportOverwrite   bool
	flagImportNoKredsfile bool
	flagImportNamespaces  []string
)

var importCmd = &cobra.Command{
	Use:           "import <file>",
	Short:         helpImportCmd,
	Long:          banner(helpImportCmd),
	GroupID:       "secrets",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			termactions.Log().Error("Expected exactly one argument: path to the file to import")
			os.Exit(1)
		}

		filePath := args[0]

		data, err := os.ReadFile(filePath)
		if err != nil {
			termactions.Log().Error("Could not read file: " + err.Error())
			os.Exit(1)
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		base := filepath.Base(filePath)

		isEnvFile := strings.Contains(base, ".env")
		isStructured := slices.Contains([]string{".json", ".yaml", ".yml", ".toml"}, ext)
		if !isEnvFile && !isStructured {
			termactions.LogGroup().Error(
				"Files with extension "+ext+" are not supported",
				"Supported formats: .env, .env.<namespace>, .json, .yaml, .yml, .toml}",
			)
			os.Exit(1)
		}

		if isEnvFile && len(flagImportNamespaces) > 1 {
			termactions.Log().Error("Only one namespace can be assigned to an env file")
			os.Exit(1)
		}

		var grouped map[string]map[string]string

		switch ext {
		case ".json":
			if err := json.Unmarshal(data, &grouped); err != nil {
				termactions.Log().Error("Could not parse the JSON document: " + err.Error())
				os.Exit(1)
			}
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(data, &grouped); err != nil {
				termactions.Log().Error("Could not parse the YAML document: " + err.Error())
				os.Exit(1)
			}
		case ".toml":
			if _, err := toml.Decode(string(data), &grouped); err != nil {
				termactions.Log().Error("Could not parse the TOML document: " + err.Error())
				os.Exit(1)
			}
		default:
			var envParseErr []string
			flat := map[string]string{}
			for i, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) != 2 {
					envParseErr = append(envParseErr, fmt.Sprintf("line %d: invalid env line: %q", i+1, line))
					continue
				}
				flat[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}

			if len(envParseErr) > 0 {
				termactions.LogGroup().Error("Could not parse env file", envParseErr...)
				os.Exit(1)
			}

			ns := ""
			parts := strings.Split(base, ".")
			if parts[0] == "" && len(parts) > 2 {
				ns = sanitizeNamespace(strings.Join(parts[2:], "."))
			}
			if len(flagImportNamespaces) > 0 {
				ns = flagImportNamespaces[0]
			}

			grouped = map[string]map[string]string{ns: flat}
		}

		if len(grouped) == 0 {
			termactions.Log().Warn("No secrets found in file")
			return
		}

		// prompt for decryption password if encrypted values exist
		decryptPassword, err := getDecryptionPassword(grouped)
		if err != nil {
			termactions.Log().Error("Could not read password: " + err.Error())
			os.Exit(1)
		}

		if decryptPassword != "" {
			if err := decryptValues(grouped, decryptPassword); err != nil {
				termactions.Log().Error(err.Error())
				os.Exit(1)
			}
		}

		// open store
		password, err := auth.Retrieve()
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}

		s, err := store.Open(password)
		if err != nil {
			termactions.Log().Error("Could not open store")
			os.Exit(1)
		}
		defer s.Close()

		if !flagImportNoKredsfile {
			if err := updateOrCreateKredsfile(grouped); err != nil {
				termactions.Log().Error("Could not update .kredsfile: " + err.Error())
				os.Exit(1)
			}
		}

		imported, skipped, misses := storeInStore(s, grouped)

		if len(misses) > 0 {
			termactions.LogGroup().Error("Failed to store some keys", misses...)
		}

		termactions.LogGroup().Info(
			"Import complete",
			"Imported "+fmt.Sprintf("%d", imported)+" keys",
			"Skipped "+fmt.Sprintf("%d", skipped)+" keys",
			"Failed to store "+fmt.Sprintf("%d", len(misses))+" keys",
		)
	},
}

func getDecryptionPassword(grouped map[string]map[string]string) (string, error) {
	for _, secrets := range grouped {
		for _, value := range secrets {
			if strings.HasPrefix(value, "enc:") {
				pwd, err := termactions.Secret().
					WithLabel("Enter decryption password").
					WithValidator(func(s string) (string, bool) {
						if strings.TrimSpace(s) == "" {
							return "password cannot be empty", false
						}
						return "", true
					}).
					Render()
				if err != nil {
					return "", fmt.Errorf("could not read password")
				}
				return pwd, nil
			}
		}
	}
	return "", nil
}

func decryptValues(grouped map[string]map[string]string, password string) error {
	for ns, secrets := range grouped {
		for key, value := range secrets {
			if strings.HasPrefix(value, "enc:") {
				decrypted, err := helpers.Decrypt(value[4:], password)
				if err != nil {
					return fmt.Errorf("could not decrypt value for %s:%s: %w", ns, key, err)
				}
				grouped[ns][key] = string(decrypted)
			}
		}
	}
	return nil
}

func storeInStore(s *store.Store, grouped map[string]map[string]string) (imported, skipped int, misses []string) {
	for ns, secrets := range grouped {
		for key, value := range secrets {
			var storeKey string
			if ns == "" || ns == "_default" {
				storeKey = key
			} else {
				storeKey = ns + ":" + key
			}

			// check if key exists
			if _, err := s.Get(storeKey); err == nil && !flagImportOverwrite {
				skipped++
				continue
			}

			if err := s.Set(storeKey, value); err != nil {
				misses = append(misses, fmt.Sprintf("Cannot store key '%s': %s", storeKey, err.Error()))
				continue
			}
			imported++
		}
	}
	return
}

func updateOrCreateKredsfile(grouped map[string]map[string]string) error {
	newFile := false
	path, err := spec.Locate()
	if err != nil || path == "" {
		newFile = true
		path = filepath.Join(".", ".kredsfile")
	}

	var existing []spec.Secret
	if _, err := os.Stat(path); err == nil {
		kf, errs := spec.Parse(path)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, e := range errs {
				errMsgs[i] = e.Error()
			}
			return fmt.Errorf("could not parse existing .kredsfile: %s", strings.Join(errMsgs, ", "))
		}
		existing = kf.Secrets
	}

	existingKeys := map[string]bool{}
	for _, s := range existing {
		existingKeys[s.Key] = true
	}

	var lines []string

	if newFile {
		lines = append(lines, "# File created by kredenv")
		lines = append(lines, "# recurse to 0 # uncomment this to recurse secrets")
		lines = append(lines, "")
		lines = append(lines, "# autoload # uncomment this to autoload secrets (accepts 'on', 'off', or 'for <namespace>')")
	}

	for ns, secrets := range grouped {
		for key := range secrets {
			var storeKey string
			if ns == "" || ns == "_default" {
				storeKey = key
			} else {
				storeKey = ns + ":" + key
			}

			if existingKeys[storeKey] && !flagImportOverwrite {
				continue
			}

			if ns == "" || ns == "_default" {
				lines = append(lines, "needs "+key)
			} else {
				lines = append(lines, fmt.Sprintf("needs %s as %s", storeKey, key))
			}
		}
	}

	if len(lines) == 0 {
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open .kredsfile: %w", err)
	}
	defer f.Close()

	if len(existing) > 0 {
		if _, err := fmt.Fprintln(f); err != nil {
			return err
		}
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(f, line); err != nil {
			return err
		}
	}

	if newFile {
		termactions.Log().Success("Created .kredsfile at " + path)
	} else {
		termactions.Log().Success("Updated .kredsfile at " + path)
	}
	return nil
}

func init() {
	importCmd.Flags().SortFlags = false
	importCmd.Flags().BoolVar(&flagImportOverwrite, "overwrite", false, "Overwrite existing keys if they already exist in the store")
	importCmd.Flags().BoolVar(&flagImportNoKredsfile, "no-kredsfile", false, "Skip updating or creating a .kredsfile")
	importCmd.Flags().StringArrayVarP(&flagImportNamespaces, "namespaces", "n", []string{}, "Import one or more specific namespaces from the file")
}
