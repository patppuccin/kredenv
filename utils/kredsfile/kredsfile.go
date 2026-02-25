package kredsfile

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

# needs <key> as <env_var>
# maybe <key> as <env_var>
`

type Secret struct {
	Key      string // the `needs` or `maybe` key
	Alias    string // the `as` alias
	Optional bool
}

type Kredsfile struct {
	RecurseDepth int
	Secrets      []Secret
}

func Locate() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine current directory: %w", err)
	}

	for {
		candidate := filepath.Join(current, ".kredsfile")

		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
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

	if len(parseErrs) > 0 {
		return nil, parseErrs
	}

	return kf, nil
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

func parseSecretsDirective(fields []string, lineNum int, optional bool) (Secret, error) {
	if len(fields) < 2 {
		return Secret{}, fmt.Errorf("line %d: invalid syntax, expected: needs <key> [as <alias>]", lineNum)
	}

	key := fields[1]

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
