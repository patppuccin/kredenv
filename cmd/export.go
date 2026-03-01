package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/helpers"
	"github.com/patppuccin/kredenv/utils/spec"
	"github.com/patppuccin/kredenv/utils/store"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

const helpExportCmd = "Export secrets from the secrets store to stdout or a file"

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
	GroupID:       "secrets",
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

		password, err := auth.Retrieve()
		if err != nil {
			console.Error(err.Error())
			os.Exit(1)
		}

		s, err := store.Open(password)
		if err != nil {
			console.Error("Could not open store")
			os.Exit(1)
		}
		defer s.Close()

		groupedSecrets, errs := collectGroupedSecrets(s)
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

		// Filter only the required namespaces if specified
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

		// Prompt for encryption password if encrypting
		var encPassword string
		if flagExportEncrypt {
			encPassword, err = console.PromptSecret("Enter encryption password: ")
			if err != nil {
				console.Error(err.Error())
				os.Exit(1)
			}
			if encPassword == "" {
				console.Error("Encryption password cannot be empty")
				os.Exit(1)
			}
		}

		console.Warn("Export contains secrets — do not commit to version control")

		switch flagExportFormat {
		case "env":
			exportEnv(groupedSecrets, encPassword)
		default:
			exportStructured(groupedSecrets, encPassword)
		}
	},
}

func exportEnv(grouped map[string]map[string]string, password string) {
	for ns, secrets := range grouped {
		var sb strings.Builder
		for k, v := range secrets {
			if password != "" {
				encrypted, err := helpers.Encrypt([]byte(v), password)
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

func exportStructured(grouped map[string]map[string]string, password string) {
	if password != "" {
		for ns, secrets := range grouped {
			for k, v := range secrets {
				encrypted, err := helpers.Encrypt([]byte(v), password)
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
				"Unable to export namespace '_default'",
				[]string{
					"The namespace '_default' conflicts with flat keys",
					"Rename your namespace to prevent namespace collisions",
				},
			)
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

	if isDir := (err == nil && info.IsDir()) || (os.IsNotExist(err) && filepath.Ext(absOutputPath) == ""); isDir {
		writeFile(filepath.Join(absOutputPath, "creds."+flagExportFormat), output)
		return
	}

	ext := strings.TrimPrefix(filepath.Ext(absOutputPath), ".")
	if ext != flagExportFormat {
		console.Error("File extension ." + ext + " does not match format " + flagExportFormat)
		os.Exit(1)
	}

	writeFile(absOutputPath, output)
}

func collectGroupedSecrets(s *store.Store) (map[string]map[string]string, []error) {
	grouped := map[string]map[string]string{}

	if flagExportAll {
		data, err := s.List()
		if err != nil {
			return nil, []error{err}
		}
		for key, value := range data {
			ns, keyName := spec.SplitNamespacedKey(key)
			if _, ok := grouped[ns]; !ok {
				grouped[ns] = map[string]string{}
			}
			grouped[ns][keyName] = value
		}
		return grouped, nil
	}

	path, err := spec.Locate()
	if err != nil {
		return nil, []error{err}
	}
	if path == "" {
		return nil, []error{fmt.Errorf("no .kredsfile found")}
	}

	kf, errs := spec.Parse(path)
	if len(errs) > 0 {
		return nil, errs
	}

	var collectErrs []error
	for _, secret := range kf.Secrets {
		value, err := s.Get(secret.Key)
		if err != nil {
			if !secret.Optional {
				collectErrs = append(collectErrs, fmt.Errorf("missing required key: %s", secret.Key))
			}
			continue
		}
		ns, _ := spec.SplitNamespacedKey(secret.Key)
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
	exportCmd.Flags().BoolVar(&flagExportAll, "all", false, "Export all secrets in the store")
	exportCmd.Flags().BoolVar(&flagExportEncrypt, "encrypt", false, "Encrypt secret values with a password")
	exportCmd.Flags().StringVarP(&flagExportFormat, "format", "f", "env", "Export format (env, json, yaml, toml)")
	exportCmd.Flags().StringVarP(&flagExportOutput, "output", "o", "", "Output file path (defaults to stdout)")
	exportCmd.Flags().StringArrayVarP(&flagExportNamespaces, "namespaces", "n", []string{}, "Export secrets from specific namespaces (repeatable)")
}
