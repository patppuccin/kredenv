package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/patppuccin/kredenv/src/auth"
	"github.com/patppuccin/kredenv/src/spec"
	"github.com/patppuccin/kredenv/src/store"
	"github.com/patppuccin/termactions"
	"github.com/spf13/cobra"
)

const helpExecCmd = "Execute a command with secrets injected into its environment"

var execExamples = `  # Execute a command with secrets from the default namespace
  kredenv exec -- node server.js

  # Execute a command with secrets from a specific namespace
  kredenv exec -n production -- kubectl apply -f deploy.yaml
  kredenv exec --namespace staging -- terraform apply plan.tfplan`

var flagExecNamespace string

var execCmd = &cobra.Command{
	Use:           "exec -- <command> [args...]",
	Short:         helpExecCmd,
	Long:          banner(helpExecCmd),
	GroupID:       "env",
	SilenceUsage:  true,
	SilenceErrors: true,
	Example:       execExamples,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			termactions.Log().Error("No command provided to execute")
			os.Exit(1)
		}

		path, err := spec.Locate()
		if err != nil || path == "" {
			termactions.Log().Error("No kredsfile.yaml found")
			os.Exit(1)
		}

		kf, errs := spec.Parse(path)
		if len(errs) > 0 {
			errMsgs := make([]string, len(errs))
			for i, e := range errs {
				errMsgs[i] = e.Error()
			}
			termactions.LogGroup().Error("Failed to parse kredsfile.yaml", errMsgs...)
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

		ns := flagExecNamespace
		if ns == "" {
			ns = kf.AutoloadNamespace
		}

		resolved := map[string]string{}
		var missingRequired []string

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
				missingRequired = append(missingRequired, secret.VaultKey())
				continue
			}

			resolved[secret.Key] = value
		}

		if len(missingRequired) > 0 {
			termactions.LogGroup().Error("Missing required secrets", missingRequired...)
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
		termactions.Log().Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	execCmd.Flags().SortFlags = false
	execCmd.Flags().StringVarP(&flagExecNamespace, "namespace", "n", "", "Execute command with secrets from a specific namespace")
}
