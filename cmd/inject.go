package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/spec"
	"github.com/patppuccin/kredenv/utils/store"
	"github.com/spf13/cobra"
)

const helpInjectCmd = "Emits export statements for the shell hook (internal use)"

var (
	flagInjectNamespace string
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
			console.Error("inject is for shell hook use only, use `kredenv load` instead")
			return
		}

		path, err := spec.Locate()
		if err != nil || path == "" {
			return
		}

		kf, errs := spec.Parse(path)
		if len(errs) > 0 {
			return
		}

		if kf.AutoloadOff {
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
			if err != nil {
				continue
			}
			fmt.Printf("export %s=%q\n", secret.Alias, value)
		}
	},
}

func init() {
	injectCmd.Flags().SortFlags = false
	injectCmd.Flags().StringVar(&flagInjectNamespace, "namespace", "", "Inject keys from a specific namespace")
}
