# Shell Hooks

The shell hook watches for directory changes and automatically loads or unloads secrets based on the `kredsfile.yaml` in scope. Install it once per shell and it runs transparently on every `cd`.

## How it works

The hook installs a shell function that intercepts `kredenv` commands and wraps directory change detection. When you move into a directory with a `kredsfile.yaml` that has `autoload: true`, secrets are injected into your environment. When you leave, they are unloaded.

The hook also intercepts `kredenv load` and `kredenv unload` so they work correctly in your interactive shell.

## Installation

### Bash

Add to `~/.bashrc`:

```bash
echo 'eval "$(kredenv hook bash)"' >> ~/.bashrc
```

Then restart your shell or run:

```bash
source ~/.bashrc
```

### Zsh

Add to `~/.zshrc`:

```bash
echo 'eval "$(kredenv hook zsh)"' >> ~/.zshrc
```

Then restart your shell or run:

```bash
source ~/.zshrc
```

### Fish

Add to `~/.config/fish/config.fish`:

```bash
echo 'kredenv hook fish | source' >> $__fish_config_dir/config.fish
```

Then restart your shell.

### PowerShell

Add to your PowerShell profile:

```powershell
Add-Content $PROFILE 'Invoke-Expression (& { (kredenv hook powershell | Out-String) })'
```

Then restart PowerShell or run:

```powershell
. $PROFILE
```

### Nushell

Nushell loads hooks from saved scripts in the autoload directory:

```bash
kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")
```

::: warning
After upgrading kredenv, re-run this command to regenerate the hook script with the latest version.
:::

## Verifying the Hook

Check that the hook is active:

```bash
echo $__KREDENV_BIN
```

If this prints the path to the kredenv binary, the hook is installed and active.

## Starship Integration

kredenv exposes two environment variables for use in prompt frameworks:

- `__KREDENV_LOADED_NS` — namespace currently loaded
- `__KREDENV_LOADED_COUNT` — number of secrets currently loaded
- `__KREDENV_LOADED_VARS` — comma-separated list of loaded secret names

See [Prompt Frameworks](/docs/prompt-frameworks) for setup examples.
