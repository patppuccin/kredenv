package spec

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var MinimalTemplate = `# .kredsfile
# safe to commit - contains no secrets
# kredenv errors on missing 'needs', warns on missing 'maybe'

# recurse to <depth>

# autoload                  (load flat keys on cd - default behaviour)
# autoload on               (same as above, explicit form)
# autoload off              (disable autoloading entirely)
# autoload for <namespace>  (set default namespace to load)

# needs <key> as <env_var>
# maybe <key> as <env_var>
`

type Secret struct {
	Key      string // the `needs` or `maybe` key
	Alias    string // the `as` alias
	Optional bool
}

type Kredsfile struct {
	RecurseDepth      int
	AutoloadOff       bool
	AutoloadNamespace string
	Secrets           []Secret
}

func Locate() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine current directory: %w", err)
	}

	start := current

	for {
		candidate := filepath.Join(current, ".kredsfile")

		if _, err := os.Stat(candidate); err == nil {
			recurseDepth, err := peekRecurseDepth(candidate)
			if err != nil {
				return "", err
			}

			rel, err := filepath.Rel(current, start)
			if err != nil {
				return "", fmt.Errorf("could not calculate relative path: %w", err)
			}
			levelsDeep := len(strings.Split(rel, string(filepath.Separator))) - 1

			if recurseDepth == 0 || levelsDeep <= recurseDepth {
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
	}
}

func Parse(path string) (*Kredsfile, []error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, []error{fmt.Errorf("could not open %s: %w", path, err)}
	}
	defer f.Close()

	kf := &Kredsfile{}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	sawAutoload := false
	var parseErrs []error

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)

		switch fields[0] {
		case "recurse":
			depth, err := parseRecurseDirective(fields, lineNum)
			if err != nil {
				parseErrs = append(parseErrs, err)
				continue
			}
			kf.RecurseDepth = depth

		case "autoload":
			if sawAutoload {
				parseErrs = append(parseErrs, fmt.Errorf("line %d: duplicate autoload directive", lineNum))
				continue
			}
			sawAutoload = true

			isOff, namespace, err := parseAutoloadDirective(fields, lineNum)
			if err != nil {
				parseErrs = append(parseErrs, err)
				continue
			}
			kf.AutoloadOff = isOff
			kf.AutoloadNamespace = namespace

		case "needs":
			secret, err := parseSecretsDirective(fields, lineNum, false)
			if err != nil {
				parseErrs = append(parseErrs, err)
				continue
			}
			kf.Secrets = append(kf.Secrets, secret)

		case "maybe":
			secret, err := parseSecretsDirective(fields, lineNum, true)
			if err != nil {
				parseErrs = append(parseErrs, err)
				continue
			}
			kf.Secrets = append(kf.Secrets, secret)

		default:
			parseErrs = append(parseErrs, fmt.Errorf("line %d: found unknown directive %q", lineNum, fields[0]))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, []error{fmt.Errorf("error reading file: %w", err)}
	}

	// Validate autoload namespace exists in secrets.
	if len(parseErrs) == 0 && kf.AutoloadNamespace != "" {
		found := false
		for _, s := range kf.Secrets {
			if strings.HasPrefix(s.Key, kf.AutoloadNamespace+":") {
				found = true
				break
			}
		}
		if !found {
			parseErrs = append(parseErrs, fmt.Errorf("autoload namespace %q has no matching keys in this file", kf.AutoloadNamespace))
		}
	}

	if len(parseErrs) > 0 {
		return nil, parseErrs
	}

	return kf, nil
}

func peekRecurseDepth(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("unable to open file %s", path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 3 && fields[0] == "recurse" && fields[1] == "to" {
			depth, err := strconv.Atoi(fields[2])
			if err != nil || depth < 1 {
				return 0, fmt.Errorf("invalid recurse depth in file %s", path)
			}
			return depth, nil
		}
	}

	return 0, nil
}

func parseRecurseDirective(fields []string, lineNum int) (int, error) {
	if len(fields) != 3 || fields[1] != "to" {
		return 0, fmt.Errorf("line %d: invalid syntax, expected: recurse to <depth-as-integer>", lineNum)
	}

	depth, err := strconv.Atoi(fields[2])
	if err != nil || depth < 1 {
		return 0, fmt.Errorf("line %d: depth must be +ve integer, got %q", lineNum, fields[2])
	}

	return depth, nil
}

func parseAutoloadDirective(fields []string, lineNum int) (isOff bool, namespace string, err error) {
	// bare autoload (or) autoload on == autoload flat keys (without namespace)
	if len(fields) == 1 || (len(fields) == 2 && fields[1] == "on") {
		return false, "", nil
	}

	// autoload off == disable autoloading
	if len(fields) == 2 && fields[1] == "off" {
		return true, "", nil
	}

	// autoload for <namespace> == set default namespace to load
	if len(fields) == 3 && fields[1] == "for" {
		return false, fields[2], nil
	}

	return false, "", fmt.Errorf("line %d: invalid autoload syntax, expected: autoload [on|off|for <namespace>]", lineNum)
}

func parseSecretsDirective(fields []string, lineNum int, optional bool) (Secret, error) {
	if len(fields) < 2 {
		return Secret{}, fmt.Errorf("line %d: invalid syntax, expected: needs <key> [as <alias>]", lineNum)
	}

	key := fields[1]

	if strings.Count(key, ":") > 1 {
		return Secret{}, fmt.Errorf("line %d: key %q has more than one namespace separator ':'", lineNum, key)
	}

	asIdx := -1
	for i, f := range fields {
		if f == "as" {
			asIdx = i
			break
		}
	}

	var alias string

	if asIdx == -1 {
		if strings.Contains(key, ":") {
			return Secret{}, fmt.Errorf("line %d: namespaced key %q requires an 'as' alias", lineNum, key)
		}
		alias = key
	} else {
		if asIdx == len(fields)-1 {
			return Secret{}, fmt.Errorf("line %d: missing alias after 'as'", lineNum)
		}
		alias = fields[asIdx+1]
	}

	return Secret{
		Key:      key,
		Alias:    alias,
		Optional: optional,
	}, nil
}

func SplitNamespacedKey(key string) (ns string, name string) {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", parts[0]
}
