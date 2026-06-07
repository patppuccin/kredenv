![banner](assets/banner.svg)

[![Latest Release](https://img.shields.io/github/v/release/patppuccin/kredenv)](https://github.com/patppuccin/kredenv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/patppuccin/kredenv)](https://goreportcard.com/report/github.com/patppuccin/kredenv)
[![License](https://img.shields.io/github/license/patppuccin/kredenv)](LICENSE)

Keep your secrets encrypted locally and inject them into your shell as you move across projects.  
No plaintext files. No accidental commits. No secret leaks.

## Installation

**Linux & macOS**

```bash
bash <(curl -fsSL https://kredenv.patppuccin.com/install.sh)
```

**Windows (PowerShell 5.1+)**

```powershell
powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://kredenv.patppuccin.com/install.ps1 | iex"
```

## Shell Hook

Add the hook to your shell configuration file:

```bash
# Bash
echo 'eval "$(kredenv hook bash)"' >> ~/.bashrc

# Zsh
echo 'eval "$(kredenv hook zsh)"' >> ~/.zshrc

# Fish
echo 'kredenv hook fish | source' >> $__fish_config_dir/config.fish

# PowerShell
Add-Content $PROFILE 'Invoke-Expression (& { (kredenv hook powershell | Out-String) })'

# Nushell
kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")
```

## Documentation

Full documentation, guides, and command reference at **[kredenv.patppuccin.com](https://kredenv.patppuccin.com)**.

## License

Released under the [Apache License 2.0](LICENSE).
