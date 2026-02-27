package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/keyring"
	"github.com/patppuccin/kredenv/utils/kredsfile"
	"github.com/spf13/cobra"
)

const helpListCmd = "Lists keys from the local .kredsfile or the keyring"

var (
	flagListAll        bool
	flagListShowValues bool
	flagListNamespace  string
)

var listCmd = &cobra.Command{
	Use:           "list",
	Short:         helpListCmd,
	Long:          console.Banner(helpListCmd),
	GroupID:       "keyring",
	Aliases:       []string{"ls"},
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.Error("No arguments expected, got " + strconv.Itoa(len(args)))
			os.Exit(1)
		}

		if flagListAll {
			listFromKeyring(flagListNamespace)
		} else {
			listFromKredsfile(flagListNamespace)
		}
	},
}

func listFromKredsfile(ns string) {
	path, err := kredsfile.Locate()
	if err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	if path == "" {
		console.Warn("No .kredsfile found")
		os.Exit(1)
	}

	kf, errs := kredsfile.Parse(path)
	if len(errs) > 0 {
		if len(errs) == 1 {
			console.Error(errs[0].Error())
			os.Exit(1)
		}
		errMsgs := make([]string, len(errs))
		for i, e := range errs {
			errMsgs[i] = e.Error()
		}
		console.ErrorGroup("Failed to parse "+path, errMsgs)
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

		value, err := keyring.Get(secret.Key)
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
		console.InfoGroup("The following keys were found", lookupHits)
	}
	if len(lookupMisses) > 0 {
		console.WarnGroup("The following keys were not found", lookupMisses)
	}
}

func listFromKeyring(ns string) {
	keys, err := keyring.List()
	if err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	if len(keys) == 0 {
		console.Warn("No keys found in keyring")
		return
	}

	msgs := make([]string, 0, len(keys))
	for _, key := range keys {
		if ns != "" {
			if !strings.HasPrefix(key, ns+":") {
				continue
			}
		}

		if flagListShowValues {
			value, err := keyring.Get(key)
			if err != nil {
				msgs = append(msgs, fmt.Sprintf("%s = <not set>", key))
				continue
			}
			msgs = append(msgs, fmt.Sprintf("%s = %s", key, value))
		} else {
			msgs = append(msgs, key)
		}
	}

	if len(msgs) == 0 && ns != "" {
		console.Warn("No keys found for namespace: " + ns)
		return
	}

	console.InfoGroup("Keyring Keys", msgs)
}

func init() {
	listCmd.Flags().SortFlags = false
	listCmd.Flags().BoolVar(&flagListShowValues, "show-values", false, "Show secret values (use with caution)")
	listCmd.Flags().BoolVarP(&flagListAll, "all", "a", false, "List all keys in the keyring instead")
	listCmd.Flags().StringVarP(&flagListNamespace, "namespace", "n", "", "Filter by namespace")
}
