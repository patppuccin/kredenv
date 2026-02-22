package shells

import (
	_ "embed"
	"slices"
	"strings"
)

//go:embed integrations/bash-hook.sh
var hookBash string

//go:embed integrations/zsh-hook.sh
var hookZsh string

//go:embed integrations/fish-hook.fish
var hookFish string

//go:embed integrations/pwsh-hook.ps1
var hookPowershell string

//go:embed integrations/nushell-hook.nu
var hookNushell string

type Shell struct {
	Name     string
	Aliases  []string
	SetupCmd string
	Hook     string
}

var Supported = []Shell{
	{
		Name:     "bash",
		Aliases:  []string{},
		SetupCmd: `echo 'eval "$(kredenv hook bash)"' >> ~/.bashrc`,
		Hook:     hookBash,
	},
	{
		Name:     "zsh",
		Aliases:  []string{},
		SetupCmd: `echo 'eval "$(kredenv hook zsh)"' >> ~/.zshrc`,
		Hook:     hookZsh,
	},
	{
		Name:     "fish",
		Aliases:  []string{},
		SetupCmd: `echo 'kredenv hook fish | source' >> $__fish_config_dir/config.fish`,
		Hook:     hookFish,
	},
	{
		Name:     "powershell",
		Aliases:  []string{"pwsh"},
		SetupCmd: `Add-Content $PROFILE 'Invoke-Expression (& { (kredenv hook powershell | Out-String) })'`,
		Hook:     hookPowershell,
	},
	{
		Name:     "nushell",
		Aliases:  []string{"nu"},
		SetupCmd: `kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")`,
		Hook:     hookNushell,
	},
}

func Get(name string) (*Shell, bool) {
	for i, s := range Supported {
		if s.Name == name {
			return &Supported[i], true
		}
		if slices.Contains(s.Aliases, name) {
			return &Supported[i], true
		}
	}
	return nil, false
}

func Names() string {
	names := make([]string, 0, len(Supported))
	for _, s := range Supported {
		curr := s.Name
		if len(s.Aliases) > 0 {
			curr += "/" + strings.Join(s.Aliases, "/")
		}
		names = append(names, curr)
	}
	return strings.Join(names, ", ")
}

func SetupCmds() []string {
	cmds := make([]string, 0, len(Supported))
	for _, s := range Supported {
		cmds = append(cmds, s.SetupCmd)
	}
	return cmds
}
