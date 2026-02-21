package shells

import (
	"strings"
)

type Shell struct {
	Name    string
	Aliases []string
	Hook    func() string
}

var Supported = []Shell{
	{
		Name:    "bash",
		Aliases: []string{},
		Hook: func() string {
			return `_kredenv_hook() {
  local previous="$OLDPWD"
  local current="$PWD"
  if [ "$previous" != "$current" ]; then
    if output=$(kredenv load 2>/dev/null); then
      eval "$output"
    fi
  fi
}

PROMPT_COMMAND="_kredenv_hook;${PROMPT_COMMAND}"`
		},
	},
	{
		Name:    "zsh",
		Aliases: []string{},
		Hook: func() string {
			return `_kredenv_hook() {
  if output=$(kredenv load 2>/dev/null); then
    eval "$output"
  fi
}

autoload -Uz add-zsh-hook
add-zsh-hook chpwd _kredenv_hook`
		},
	},
	{
		Name:    "fish",
		Aliases: []string{},
		Hook: func() string {
			return `function _kredenv_hook --on-variable PWD
  if set output (kredenv load 2>/dev/null)
    eval $output
  end
end`
		},
	},
	{
		Name:    "powershell",
		Aliases: []string{"pwsh"},
		Hook: func() string {
			return `function global:prompt {
  $output = kredenv load 2>$null
  if ($output) {
    $output | Invoke-Expression
  }
  "PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) "
}`
		},
	},
	{
		Name:    "nushell",
		Aliases: []string{"nu"},
		Hook: func() string {
			return `$env.config.hooks.env_change.PWD = ($env.config.hooks.env_change.PWD | append {|before, after|
  let output = (kredenv load | complete)
  if $output.exit_code == 0 {
    $output.stdout | lines | each {|line|
      let parts = ($line | str replace "export " "" | split row "=")
      load-env {($parts.0): ($parts.1 | str trim -c '"')}
    }
  }
})`
		},
	},
}

func Get(name string) (*Shell, bool) {
	for i, s := range Supported {
		if s.Name == name {
			return &Supported[i], true
		}
		for _, a := range s.Aliases {
			if a == name {
				return &Supported[i], true
			}
		}
	}
	return nil, false
}

func Names() string {
	names := make([]string, len(Supported))
	for i, s := range Supported {
		names[i] = s.Name
		if len(s.Aliases) > 0 {
			names[i] += " (" + strings.Join(s.Aliases, ", ") + ")"
		}
	}
	return strings.Join(names, ", ")
}
