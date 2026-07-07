package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpListCmd = "List secrets from the kredsfile.yaml or the store"

var (
	flagListAll        bool
	flagListLoaded     bool
	flagListShowValues bool
	flagListNamespace  string
)

var listCmd = &cobra.Command{
	Use:           "list",
	Short:         helpListCmd,
	Long:          banner(helpListCmd),
	GroupID:       "secrets",
	Aliases:       []string{"ls"},
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}
		if flagListAll && flagListLoaded {
			termactions.LogGroup().Error("Invalid flag combination", "Only one of --all or --loaded may be used at a time")
			os.Exit(1)
		}
		if flagListLoaded && flagListNamespace != "" {
			termactions.Log().Warn("Flag --namespace is ignored in --loaded mode")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if flagListLoaded {
			listFromShell()
			return
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

		if flagListAll {
			listFromStore(s, flagListNamespace)
		} else {
			listFromKredsfile(s, flagListNamespace)
		}
	},
}

func listFromShell() {
	if configured := os.Getenv("__KREDENV_BIN"); configured == "" {
		termactions.Log().Warn("kredenv is not configured for current shell")
		return
	}

	loaded := os.Getenv("__KREDENV_LOADED_VARS")
	if loaded == "" {
		termactions.Log().Warn("No kredenv secrets loaded in current shell session")
		return
	}

	keys := strings.Split(loaded, ",")
	entries := make([]string, 0, len(keys))

	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if flagListShowValues {
			entries = append(entries, fmt.Sprintf("%s = %s", key, os.Getenv(key)))
		} else {
			entries = append(entries, key)
		}
	}

	termactions.LogGroup().Info(fmt.Sprintf("%d secrets loaded in current shell session", len(keys)), entries...)
}

func listFromKredsfile(s *store.Store, ns string) {
	path, err := spec.Locate()
	if err != nil {
		termactions.Log().Error(err.Error())
		os.Exit(1)
	}
	if path == "" {
		termactions.Log().Warn("No kredsfile.yaml found in current scope")
		os.Exit(1)
	}

	kf, errs := spec.Parse(path)
	if len(errs) > 0 {
		errMsgs := make([]string, len(errs))
		for i, e := range errs {
			errMsgs[i] = e.Error()
		}
		termactions.LogGroup().Error("Failed to parse "+path, errMsgs...)
		os.Exit(1)
	}

	if ns == "" && kf.AutoloadNamespace != "" {
		ns = kf.AutoloadNamespace
	}

	var lookupHits []string
	var lookupMisses []string

	for _, secret := range kf.Secrets {
		if ns != "" {
			if secret.Namespace != ns {
				continue
			}
		} else {
			if secret.Namespace != "" {
				continue
			}
		}

		value, err := s.Get(secret.VaultKey())
		if err == nil {
			if flagListShowValues {
				lookupHits = append(lookupHits, fmt.Sprintf("%s = %s", secret.Key, value))
			} else {
				lookupHits = append(lookupHits, secret.Key)
			}
		} else {
			lookupMisses = append(lookupMisses, fmt.Sprintf("%s = <not set>", secret.Key))
		}
	}

	hitsHeader := fmt.Sprintf("%d secrets in scope", len(lookupHits))
	if ns != "" {
		hitsHeader += fmt.Sprintf(" (namespace: %s)", ns)
	}

	missHeader := fmt.Sprintf("%d secrets found, but not set in store", len(lookupMisses))
	if ns != "" {
		missHeader += fmt.Sprintf(" (namespace: %s)", ns)
	}

	if len(lookupHits) > 0 {
		termactions.LogGroup().Info(hitsHeader, lookupHits...)
	}
	if len(lookupMisses) > 0 {
		termactions.LogGroup().Warn(missHeader, lookupMisses...)
	}
}

func listFromStore(s *store.Store, ns string) {
	data, err := s.List()
	if err != nil {
		termactions.Log().Error(err.Error())
		os.Exit(1)
	}
	if len(data) == 0 {
		termactions.Log().Warn("No secrets found in store")
		return
	}

	entries := make([]string, 0, len(data))
	for key, value := range data {
		if ns != "" && !strings.HasPrefix(key, ns+":") {
			continue
		}
		if flagListShowValues {
			entries = append(entries, fmt.Sprintf("%s = %s", key, value))
		} else {
			entries = append(entries, key)
		}
	}

	if len(entries) == 0 && ns != "" {
		termactions.Log().Warn("No secrets found for namespace: " + ns)
		return
	}

	termactions.LogGroup().Info(fmt.Sprintf("%d secrets found in store", len(entries)), entries...)
}

func init() {
	listCmd.Flags().SortFlags = false
	listCmd.Flags().BoolVar(&flagListShowValues, "show-values", false, "Show secret values (use with caution)")
	listCmd.Flags().BoolVar(&flagListAll, "all", false, "List all secrets in the store")
	listCmd.Flags().BoolVar(&flagListLoaded, "loaded", false, "List all secrets loaded in the current shell")
	listCmd.Flags().StringVarP(&flagListNamespace, "namespace", "n", "", "Filter by namespace")
}
