package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/patppuccin/kredenv/utils/auth"
	"github.com/patppuccin/kredenv/utils/console"
	"github.com/patppuccin/kredenv/utils/spec"
	"github.com/patppuccin/kredenv/utils/store"
	"github.com/spf13/cobra"
)

const helpExecCmd = "Execute a command with secrets injected into its environment"

var execExamples = `  # Execute a command with secrets from the default namespace
  kredenv exec -- node server.js

  # Execute a command with secrets from a specific namespace
  kredenv exec -n staging -- terraform apply --auto-approve
  kredenv exec -n production -- kubectl apply -f deploy.yaml`

var (
	flagExecNamespace string
)

var execCmd = &cobra.Command{
	Use:           "exec -- <command> [args...]",
	Short:         helpExecCmd,
	Long:          console.Banner(helpExecCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Example:       execExamples,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			console.Error("No command provided to execute")
			os.Exit(1)
		}

		path, err := spec.Locate()
		if err != nil || path == "" {
			console.Error("No .kredsfile found")
			os.Exit(1)
		}

		kf, errs := spec.Parse(path)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, e := range errs {
				errMsgs[i] = e.Error()
			}
			console.ErrorGroup("Failed to parse .kredsfile", errMsgs)
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

		ns := flagExecNamespace
		if ns == "" {
			ns = kf.AutoloadNamespace
		}

		resolved := map[string]string{}
		var missingRequired []string

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
				if !secret.Optional {
					missingRequired = append(missingRequired, secret.Key)
				}
				continue
			}
			resolved[secret.Alias] = value
		}

		if len(missingRequired) > 0 {
			console.ErrorGroup("Missing required secrets", missingRequired)
			os.Exit(1)
		}

		runWith(args, resolved)
	},
}

func runWith(args []string, secrets map[string]string) {
	c := exec.Command(args[0], args[1:]...)
	c.Env = os.Environ()
	for k, v := range secrets {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		console.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	execCmd.Flags().SortFlags = false
	execCmd.Flags().StringVarP(&flagExecNamespace, "namespace", "n", "", "Execute command with secrets from a specific namespace")
}
