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
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/helpers"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

const helpExportCmd = "Export secrets from the secrets store to stdout or a file"

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
	Long:          banner(helpExportCmd),
	GroupID:       "secrets",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments supported for export command")
			os.Exit(1)
		}

		if !slices.Contains(consts.SupportedExportFormats, flagExportFormat) {
			termactions.LogGroup().Error(
				"Format '"+flagExportFormat+"' is not supported",
				"Supported formats: "+strings.Join(consts.SupportedExportFormats, ", "),
			)
			os.Exit(1)
		}

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

		groupedSecrets, errs := collectGroupedSecrets(s)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, err := range errs {
				errMsgs[i] = err.Error()
			}
			termactions.LogGroup().Error("Failed to collect secrets", errMsgs...)
			os.Exit(1)
		}

		if len(groupedSecrets) == 0 {
			termactions.Log().Warn("No secrets to export")
			return
		}

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
				termactions.Log().Warn("No secrets found for the specified namespaces")
				return
			}
		}

		var encPassword string
		if flagExportEncrypt {
			encPassword, err = termactions.Secret().
				WithLabel("Enter encryption password").
				WithValidator(func(s string) (string, bool) {
					if strings.TrimSpace(s) == "" {
						return "encryption password cannot be empty", false
					}
					return "", true
				}).
				Render()
			if err != nil {
				if err == termactions.ErrInterrupted {
					termactions.Log().Warn("Export aborted by user")
					os.Exit(1)
				}
				termactions.Log().Error(err.Error())
				os.Exit(1)
			}
		}

		termactions.Log().Warn("Export contains secrets — do not commit to version control")

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
					termactions.Log().Error("Could not encrypt value for " + k)
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
			termactions.Log().Error("Could not resolve " + flagExportOutput + ": " + err.Error())
			os.Exit(1)
		}

		info, err := os.Stat(absOutputPath)
		if err != nil && !os.IsNotExist(err) {
			termactions.Log().Error("Could not identify path: " + absOutputPath)
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
					termactions.Log().Error("Could not encrypt value for " + k)
					os.Exit(1)
				}
				grouped[ns][k] = "enc:" + encrypted
			}
		}
	}

	if secrets, ok := grouped[""]; ok {
		if _, conflict := grouped["_default"]; conflict {
			termactions.LogGroup().Error(
				"Unable to export namespace '_default'",
				"The namespace '_default' conflicts with flat keys",
				"Rename your namespace to prevent namespace collisions",
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
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}
		output = string(data) + "\n"
	case "yaml":
		data, err := yaml.Marshal(grouped)
		if err != nil {
			termactions.Log().Error(err.Error())
			os.Exit(1)
		}
		output = string(data) + "\n"
	case "toml":
		data, err := toml.Marshal(grouped)
		if err != nil {
			termactions.Log().Error(err.Error())
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
		termactions.Log().Error("Could not resolve " + flagExportOutput + ": " + err.Error())
		os.Exit(1)
	}

	info, err := os.Stat(absOutputPath)
	if err != nil && !os.IsNotExist(err) {
		termactions.Log().Error("Could not identify path: " + absOutputPath)
		os.Exit(1)
	}

	if isDir := (err == nil && info.IsDir()) || (os.IsNotExist(err) && filepath.Ext(absOutputPath) == ""); isDir {
		writeFile(filepath.Join(absOutputPath, "creds."+flagExportFormat), output)
		return
	}

	ext := strings.TrimPrefix(filepath.Ext(absOutputPath), ".")
	if ext != flagExportFormat {
		termactions.Log().Error("File extension ." + ext + " does not match format " + flagExportFormat)
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
			ns, keyName, _ := strings.Cut(key, ":")
			if keyName == "" {
				// no colon — flat key, ns is actually the key name
				keyName = ns
				ns = ""
			}
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
		return nil, []error{fmt.Errorf("no kredsfile.yaml found")}
	}

	kf, errs := spec.Parse(path)
	if len(errs) > 0 {
		return nil, errs
	}

	var collectErrs []error
	for _, secret := range kf.Secrets {
		value, err := s.Get(secret.VaultKey())
		if err != nil {
			collectErrs = append(collectErrs, fmt.Errorf("missing required key: %s", secret.VaultKey()))
			continue
		}
		ns := secret.Namespace
		if _, ok := grouped[ns]; !ok {
			grouped[ns] = map[string]string{}
		}
		grouped[ns][secret.Key] = value
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
		termactions.Log().Error("Could not resolve " + path + ": " + err.Error())
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		termactions.Log().Error("Could not create directory: " + err.Error())
		os.Exit(1)
	}
	if err := os.WriteFile(absPath, []byte(content), 0600); err != nil {
		termactions.Log().Error("Could not write " + absPath + ": " + err.Error())
		os.Exit(1)
	}
	termactions.Log().Success("Exported to " + absPath)
}

func init() {
	exportCmd.Flags().SortFlags = false
	exportCmd.Flags().BoolVar(&flagExportAll, "all", false, "Export all secrets in the store")
	exportCmd.Flags().BoolVar(&flagExportEncrypt, "encrypt", false, "Encrypt secret values with a password")
	exportCmd.Flags().StringVarP(&flagExportFormat, "format", "f", "env", "Export format (env, json, yaml, toml)")
	exportCmd.Flags().StringVarP(&flagExportOutput, "output", "o", "", "Output file path (defaults to stdout)")
	exportCmd.Flags().StringArrayVarP(&flagExportNamespaces, "namespaces", "n", []string{}, "Export secrets from specific namespaces (repeatable)")
}
