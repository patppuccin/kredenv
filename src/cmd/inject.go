package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpInjectCmd = "Emits export statements for the shell hook (internal use)"

var (
	flagInjectNamespace string
	flagInjectFormat    string
)

var injectCmd = &cobra.Command{
	Use:           "inject",
	Short:         helpInjectCmd,
	GroupID:       "env",
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			termactions.Log().Error("inject is for shell hook use only, use `kredenv load` instead")
			return
		}

		if !slices.Contains(consts.SupportedInjectFormats, flagInjectFormat) {
			return // silent exit - don't break shell hook
		}

		path, err := spec.Locate()
		if err != nil || path == "" {
			return
		}

		kf, errs := spec.Parse(path)
		if len(errs) > 0 {
			return
		}

		if !kf.Autoload {
			return
		}

		password, err := auth.Retrieve()
		if err != nil {
			return // silent exit — don't break the shell if no master password
		}

		s, err := store.Open(password)
		if err != nil {
			return
		}
		defer s.Close()

		ns := kf.AutoloadNamespace
		if flagInjectNamespace != "" {
			ns = flagInjectNamespace
		}

		resolved := map[string]string{}

		// inject namespace as env var for tracking
		if ns != "" {
			resolved["__KREDENV_LOADED_NS"] = ns
		}

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
			if err != nil {
				continue
			}

			resolved[secret.Key] = value
		}

		switch flagInjectFormat {
		case "dotenv":
			for k, v := range resolved {
				fmt.Printf("%s=%s\n", k, v)
			}
		case "json":
			data, _ := json.Marshal(resolved)
			fmt.Println(string(data))
		}
	},
}

func init() {
	injectCmd.Flags().SortFlags = false
	injectCmd.Flags().StringVar(&flagInjectNamespace, "namespace", "", "Inject keys from a specific namespace")
	injectCmd.Flags().StringVar(&flagInjectFormat, "format", "dotenv", "Output format ("+strings.Join(consts.SupportedInjectFormats, ", ")+")")
}
