package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/crypto"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/spf13/cobra"
)

const helpImportCmd = "Imports secrets from a file into the keyring"

var (
	flagImportOverwrite bool
)

var importCmd = &cobra.Command{
	Use:           "import <file>",
	Short:         helpImportCmd,
	Long:          console.Banner(helpImportCmd),
	Args:          cobra.ExactArgs(1),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		data, err := os.ReadFile(filePath)
		if err != nil {
			console.Error("could not read file: " + err.Error())
			os.Exit(1)
		}

		// decrypt if needed
		content := string(data)
		if isEncrypted(content) {
			fmt.Print("enter decryption password: ")
			var password string
			fmt.Scan(&password)
			fmt.Println()

			plaintext, err := crypto.Decrypt(content, password)
			if err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
			content = string(plaintext)
		}

		// detect format and parse
		secrets, err := parseImport(filePath, content)
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		if len(secrets) == 0 {
			console.Warn("no secrets found in file")
			return
		}

		// store in keyring
		imported := 0
		skipped := 0
		for key, value := range secrets {
			if keyring.Exists(key) && !flagImportOverwrite {
				console.Warn("skipping existing key: " + key)
				skipped++
				continue
			}
			if err := keyring.Set(key, value); err != nil {
				console.Error("could not store key " + key + ": " + err.Error())
				os.Exit(1)
			}
			imported++
		}

		console.Info(fmt.Sprintf("imported %d keys, skipped %d", imported, skipped))
	},
}

// isEncrypted checks if content looks like a base64 encrypted payload
func isEncrypted(content string) bool {
	content = strings.TrimSpace(content)
	// base64 only contains these chars and no spaces or newlines mid-content
	for _, c := range content {
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=", c) {
			return false
		}
	}
	return true
}

// parseImport detects format from extension then content and parses accordingly
func parseImport(filePath, content string) (map[string]string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return parseJSON(content)
	case ".yaml", ".yml":
		return parseYAML(content)
	case ".toml":
		return parseTOML(content)
	case ".env", "":
		return parseENV(content)
	default:
		// sniff content
		return sniffAndParse(content)
	}
}

// sniffAndParse attempts to detect format from content
func sniffAndParse(content string) (map[string]string, error) {
	content = strings.TrimSpace(content)

	if strings.HasPrefix(content, "{") {
		return parseJSON(content)
	}
	if strings.Contains(content, " = ") {
		return parseTOML(content)
	}
	if strings.Contains(content, ": ") {
		return parseYAML(content)
	}
	// fallback to env
	return parseENV(content)
}

func parseJSON(content string) (map[string]string, error) {
	var secrets map[string]string
	if err := json.Unmarshal([]byte(content), &secrets); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}
	return secrets, nil
}

func parseENV(content string) (map[string]string, error) {
	secrets := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env line: %q", line)
		}
		secrets[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return secrets, nil
}

func parseYAML(content string) (map[string]string, error) {
	secrets := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid yaml line: %q", line)
		}
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		secrets[strings.TrimSpace(parts[0])] = value
	}
	return secrets, nil
}

func parseTOML(content string) (map[string]string, error) {
	secrets := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid toml line: %q", line)
		}
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		secrets[strings.TrimSpace(parts[0])] = value
	}
	return secrets, nil
}

func init() {
	importCmd.Flags().SortFlags = false
	importCmd.Flags().BoolVar(&flagImportOverwrite, "overwrite", false, "Overwrite existing keys if they already exist in the keyring")
}
