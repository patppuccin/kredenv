package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/crypto"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

const helpExportCmd = "Exports secrets from the keyring to stdout or a file"

var supportedExportFormats = []string{"env", "json", "yaml", "toml"}

var (
	flagExportAll        bool
	flagExportFormat     string
	flagExportOutput     string
	flagExportEncrypt    bool
	flagExportNamespaces []string
)

var exportCmd = &cobra.Command{
	Use:           "export",
	Short:         helpExportCmd,
	Long:          console.Banner(helpExportCmd),
	GroupID:       "keyring",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("No arguments supported for export command")
			os.Exit(1)
		}

		if !slices.Contains(supportedExportFormats, flagExportFormat) {
			console.ErrorGroup(
				"Format '"+flagExportFormat+"' is not supported",
				[]string{"Supported formats: " + strings.Join(supportedExportFormats, ", ")},
			)
			os.Exit(1)
		}

		groupedSecrets, errs := collectGroupedSecrets()
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, err := range errs {
				errMsgs[i] = err.Error()
			}
			console.ErrorGroup("Failed to collect secrets", errMsgs)
			os.Exit(1)
		}

		if len(groupedSecrets) == 0 {
			console.Warn("No secrets to export")
			return
		}

		// Filter only the required namespaces (if specified)
		if len(flagExportNamespaces) > 0 {
			nsSet := map[string]bool{}
			for _, ns := range flagExportNamespaces {
				nsSet[ns] = true
			}
			for ns := range groupedSecrets {
				if !nsSet[ns] {
					delete(groupedSecrets, ns)
				}
			}
			if len(groupedSecrets) == 0 {
				console.Warn("No secrets found for the specified namespaces")
				return
			}
		}

		// Prompt password once if encrypting
		var password string
		if flagExportEncrypt {
			var err error
			password, err = promptPassword()
			if err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
		}

		console.Warn("Export contains secrets - do not commit to version control")

		switch flagExportFormat {
		case "env":
			exportEnv(groupedSecrets, password)
		default:
			exportStructured(groupedSecrets, password)
		}
	},
}

func exportEnv(grouped map[string]map[string]string, password string) {
	for ns, secrets := range grouped {
		var sb strings.Builder
		for k, v := range secrets {
			if password != "" {
				encrypted, err := crypto.Encrypt([]byte(v), password)
				if err != nil {
					console.Error("Could not encrypt value for " + k)
					os.Exit(1)
				}
				v = "enc:" + encrypted
			}
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
		output := sb.String()

		if flagExportOutput == "" {
			fmt.Print(output)
			continue
		}

		absOutputPath, err := filepath.Abs(flagExportOutput)
		if err != nil {
			console.Error("Could not resolve " + flagExportOutput + ": " + err.Error())
			os.Exit(1)
		}

		info, err := os.Stat(absOutputPath)
		if err != nil && !os.IsNotExist(err) {
			console.Error("Could not identify path: " + absOutputPath)
			os.Exit(1)
		}

		isDir := (err == nil && info.IsDir()) || (os.IsNotExist(err) && filepath.Ext(absOutputPath) == "")
		if !isDir {
			absOutputPath = filepath.Dir(absOutputPath)
		}

		envFilename := ".env"
		if ns != "" {
			envFilename = ".env." + sanitizeNamespace(ns)
		}
		writeFile(filepath.Join(absOutputPath, envFilename), output)
	}
}

// Writes a single file with namespaces as top-level keys
func exportStructured(grouped map[string]map[string]string, password string) {
	if password != "" {
		for ns, secrets := range grouped {
			for k, v := range secrets {
				encrypted, err := crypto.Encrypt([]byte(v), password)
				if err != nil {
					console.Error("Could not encrypt value for " + k)
					os.Exit(1)
				}
				grouped[ns][k] = "enc:" + encrypted
			}
		}
	}

	if secrets, ok := grouped[""]; ok {
		if _, conflict := grouped["_default"]; conflict {
			console.ErrorGroup(
				"Unable to export namespace '_default",
				[]string{
					"The namespace '_default' conflicts with flat keys",
					"Rename your namespace to prevent namespace collisions",
				},
			)
			console.Error("Cannot export: namespace 'default' as it conflicts with flat keys")
			os.Exit(1)
		}
		grouped["_default"] = secrets
		delete(grouped, "")
	}

	var output string

	switch flagExportFormat {
	case "json":
		data, err := json.MarshalIndent(grouped, "", "  ")
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
		output = string(data) + "\n"
	case "yaml":
		data, err := yaml.Marshal(grouped)
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
		output = string(data) + "\n"
	case "toml":
		data, err := toml.Marshal(grouped)
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}
		output = string(data) + "\n"
	}

	if flagExportOutput == "" {
		fmt.Print(output)
		return
	}

	absOutputPath, err := filepath.Abs(flagExportOutput)
	if err != nil {
		console.Error("Could not resolve " + flagExportOutput + ": " + err.Error())
		os.Exit(1)
	}

	info, err := os.Stat(absOutputPath)
	if err != nil && !os.IsNotExist(err) {
		console.Error("Could not identify path: " + absOutputPath)
		os.Exit(1)
	}

	// If it is a directory, write to dir-path/creds.<format> file
	if isDir := (err == nil && info.IsDir()) || (os.IsNotExist(err) && filepath.Ext(absOutputPath) == ""); isDir {
		writeFile(filepath.Join(absOutputPath, "creds."+flagExportFormat), output)
		return
	}

	// If it is a file, error if the extension does not match the format
	ext := strings.TrimPrefix(filepath.Ext(absOutputPath), ".")
	if ext != flagExportFormat {
		console.Error("File extension ." + ext + " does not match format " + flagExportFormat)
		os.Exit(1)
	}

	writeFile(absOutputPath, output)
}

func collectGroupedSecrets() (map[string]map[string]string, []error) {
	grouped := map[string]map[string]string{}

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
			ns, keyName := kredsfile.SplitNamespacedKey(key)
			if _, ok := grouped[ns]; !ok {
				grouped[ns] = map[string]string{}
			}
			grouped[ns][keyName] = value
		}
		return grouped, nil
	}

	path, err := kredsfile.Locate()
	if err != nil {
		return nil, []error{err}
	}
	if path == "" {
		return nil, []error{fmt.Errorf("no .kredsfile found")}
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
				collectErrs = append(collectErrs, fmt.Errorf("missing required key: %s", secret.Key))
			}
			continue
		}
		ns, _ := kredsfile.SplitNamespacedKey(secret.Key)
		if _, ok := grouped[ns]; !ok {
			grouped[ns] = map[string]string{}
		}
		grouped[ns][secret.Alias] = value
	}

	if len(collectErrs) > 0 {
		return nil, collectErrs
	}

	return grouped, nil
}

func promptPassword() (string, error) {
	fmt.Print("Enter encryption password: ")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("could not read password")
	}
	password = strings.TrimSpace(password)
	fmt.Println()
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	return password, nil
}

func sanitizeNamespace(ns string) string {
	if ns == "" {
		return ""
	}
	var sb strings.Builder
	for _, r := range ns {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|', '\x00':
			sb.WriteRune('-')
		default:
			sb.WriteRune(r)
		}
	}
	return strings.TrimSpace(sb.String())
}

func writeFile(path, content string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		console.Error("Could not resolve " + path + ": " + err.Error())
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		console.Error("Could not create directory: " + err.Error())
		os.Exit(1)
	}
	if err := os.WriteFile(absPath, []byte(content), 0600); err != nil {
		console.Error("Could not write " + absPath + ": " + err.Error())
		os.Exit(1)
	}
	console.Success("Exported to " + absPath)
}

func init() {
	exportCmd.Flags().SortFlags = false
	exportCmd.Flags().BoolVar(&flagExportAll, "all", false, "Export all keys in the keyring")
	exportCmd.Flags().BoolVar(&flagExportEncrypt, "encrypt", false, "Encrypt secret values with a password")
	exportCmd.Flags().StringVarP(&flagExportFormat, "format", "f", "env", "Export format (env, json, yaml, toml)")
	exportCmd.Flags().StringVarP(&flagExportOutput, "output", "o", "", "Output file path (defaults to stdout)")
	exportCmd.Flags().StringArrayVarP(&flagExportNamespaces, "namespaces", "n", []string{}, "Export keys from specific namespaces (repeatable)")
}
