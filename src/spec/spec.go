package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

var MinimalTemplate = `# kredsfile.yaml
# safe to commit - contains no secrets

# recurse: 3                       # walk up N levels looking for a kredsfile.yaml
# autoload: true                   # inject secrets into shell on cd (default: false)
# autoload_namespace: development  # namespace to autoload (default: secrets without namespace)

secrets: []
  # - key: MY_SECRET
  #
  # - key: DATABASE_PASSWORD
  #   alias: DB_PASSWORD
  #   namespace: production
  #
  # - key: GOOGLE_ANALYTICS_ID
  #   alias: ANALYTICS_ID
`

type Secret struct {
	Key       string `yaml:"key"`
	Namespace string `yaml:"namespace,omitempty"`
	Alias     string `yaml:"alias,omitempty"`
}

func (s Secret) VaultKey() string {
	if s.Namespace != "" {
		return s.Namespace + ":" + s.Key
	}
	return s.Key
}

func (s Secret) EnvKey() string {
	if s.Alias != "" {
		return s.Alias
	}
	return s.Key
}

type Kredsfile struct {
	RecurseDepth      int      `yaml:"recurse,omitempty"`
	Autoload          bool     `yaml:"autoload,omitempty"`
	AutoloadNamespace string   `yaml:"autoload_namespace,omitempty"`
	Secrets           []Secret `yaml:"secrets,omitempty"`
}

func Locate() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine current directory: %w", err)
	}

	levelsUp := 0

	for {
		candidate := filepath.Join(current, "kredsfile.yaml")

		if _, err := os.Stat(candidate); err == nil {
			recurseDepth, err := peekRecurseDepth(candidate)
			if err != nil {
				return "", err
			}

			if recurseDepth == 0 {
				if levelsUp == 0 {
					return candidate, nil
				}
				return "", nil
			}

			if levelsUp <= recurseDepth {
				return candidate, nil
			}

			return "", nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("could not access %s: %w", candidate, err)
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", nil
		}

		current = parent
		levelsUp++
	}
}

func Parse(path string) (*Kredsfile, []error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, []error{fmt.Errorf("could not read %s: %w", path, err)}
	}

	var kf Kredsfile
	if err := yaml.Unmarshal(data, &kf); err != nil {
		return nil, []error{fmt.Errorf("could not parse %s: %w", path, err)}
	}

	var parseErrs []error

	for i, s := range kf.Secrets {
		if s.Key == "" {
			parseErrs = append(parseErrs, fmt.Errorf("secret at index %d is missing required field 'key'", i))
			continue
		}
		if s.Namespace != "" && s.Alias == "" {
			parseErrs = append(parseErrs, fmt.Errorf("secret %q has a namespace but no alias", s.Key))
		}
	}

	if kf.AutoloadNamespace != "" {
		found := false
		for _, s := range kf.Secrets {
			if s.Namespace == kf.AutoloadNamespace {
				found = true
				break
			}
		}
		if !found {
			parseErrs = append(parseErrs, fmt.Errorf("autoload_namespace %q has no matching secrets in this file", kf.AutoloadNamespace))
		}
	}

	if len(parseErrs) > 0 {
		return nil, parseErrs
	}

	return &kf, nil
}

func peekRecurseDepth(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("could not read %s: %w", path, err)
	}

	var kf Kredsfile
	if err := yaml.Unmarshal(data, &kf); err != nil {
		return 0, fmt.Errorf("could not parse %s: %w", path, err)
	}

	return kf.RecurseDepth, nil
}
