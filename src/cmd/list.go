package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/console"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/spf13/cobra"
)

const helpListCmd = "List secrets from the .kredsfile or the store"

var (
	flagListAll        bool
	flagListShowValues bool
	flagListNamespace  string
)

var listCmd = &cobra.Command{
	Use:           "list",
	Short:         helpListCmd,
	Long:          console.Banner(helpListCmd),
	GroupID:       "secrets",
	Aliases:       []string{"ls"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("No arguments expected, got " + strconv.Itoa(len(args)))
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
		console.Error(err.Error())
		os.Exit(1)
	}
	if path == "" {
		console.Warn("No .kredsfile found")
		os.Exit(1)
	}

	kf, errs := spec.Parse(path)
	if len(errs) > 0 {
		errMsgs := make([]string, len(errs))
		for i, e := range errs {
			errMsgs[i] = e.Error()
		}
		console.ErrorGroup("Failed to parse "+path, errMsgs...)
		os.Exit(1)
	}

	var lookupHits []string
	var lookupMisses []string

	for _, secret := range kf.Secrets {
		if ns != "" {
			if !strings.HasPrefix(secret.Key, ns+":") {
				continue
			}
		} else {
			if strings.Contains(secret.Key, ":") {
				continue
			}
		}

		value, err := s.Get(secret.Key)
		if err == nil {
			if flagListShowValues {
				lookupHits = append(lookupHits, fmt.Sprintf("%s → %s = %s", secret.Alias, secret.Key, value))
			} else {
				lookupHits = append(lookupHits, fmt.Sprintf("%s → %s", secret.Alias, secret.Key))
			}
		} else if !secret.Optional {
			lookupMisses = append(lookupMisses, fmt.Sprintf("%s → %s = <not set>", secret.Alias, secret.Key))
		}
	}

	if len(lookupHits) > 0 {
		console.InfoGroup("The following keys were found", lookupHits...)
	}
	if len(lookupMisses) > 0 {
		console.WarnGroup("The following keys were not found", lookupMisses...)
	}
}

func listFromStore(s *store.Store, ns string) {
	data, err := s.List()
	if err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	if len(data) == 0 {
		console.Warn("No secrets found in store")
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
		console.Warn("No secrets found for namespace: " + ns)
		return
	}

	console.InfoGroup("Store secrets", msgs...)
}

func init() {
	listCmd.Flags().SortFlags = false
	listCmd.Flags().BoolVar(&flagListShowValues, "show-values", false, "Show secret values (use with caution)")
	listCmd.Flags().BoolVarP(&flagListAll, "all", "a", false, "List all secrets in the store")
	listCmd.Flags().StringVarP(&flagListNamespace, "namespace", "n", "", "Filter by namespace")
}
