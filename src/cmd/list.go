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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			termactions.Log().Error("No arguments expected, got " + strconv.Itoa(len(args)))
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

		if flagListAll {
			listFromStore(s, flagListNamespace)
		} else {
			listFromKredsfile(s, flagListNamespace)
		}
	},
}

func listFromKredsfile(s *store.Store, ns string) {
	path, err := spec.Locate()
	if err != nil {
		termactions.Log().Error(err.Error())
		os.Exit(1)
	}
	if path == "" {
		termactions.Log().Warn("No kredsfile.yaml found")
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
				lookupHits = append(lookupHits, fmt.Sprintf("%s → %s = %s", secret.Alias, secret.VaultKey(), value))
			} else {
				lookupHits = append(lookupHits, fmt.Sprintf("%s → %s", secret.Alias, secret.VaultKey()))
			}
		} else {
			lookupMisses = append(lookupMisses, fmt.Sprintf("%s → %s = <not set>", secret.Alias, secret.VaultKey()))
		}
	}

	if len(lookupHits) > 0 {
		termactions.LogGroup().Info("The following keys were found", lookupHits...)
	}
	if len(lookupMisses) > 0 {
		termactions.LogGroup().Warn("The following keys were not found", lookupMisses...)
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

	msgs := make([]string, 0, len(data))
	for key, value := range data {
		if ns != "" && !strings.HasPrefix(key, ns+":") {
			continue
		}
		if flagListShowValues {
			msgs = append(msgs, fmt.Sprintf("%s = %s", key, value))
		} else {
			msgs = append(msgs, key)
		}
	}

	if len(msgs) == 0 && ns != "" {
		termactions.Log().Warn("No secrets found for namespace: " + ns)
		return
	}

	termactions.LogGroup().Info("Store secrets", msgs...)
}

func init() {
	listCmd.Flags().SortFlags = false
	listCmd.Flags().BoolVar(&flagListShowValues, "show-values", false, "Show secret values (use with caution)")
	listCmd.Flags().BoolVarP(&flagListAll, "all", "a", false, "List all secrets in the store")
	listCmd.Flags().StringVarP(&flagListNamespace, "namespace", "n", "", "Filter by namespace")
}
