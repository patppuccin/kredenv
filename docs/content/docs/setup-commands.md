# Setup Commands

Commands for initializing and configuring kredenv on your machine.

---

## `kredenv setup`

Initializes kredenv on this machine. Creates the encrypted vault and stores the master password in your OS keyring.

Run this once per machine before using any other kredenv commands.

```bash
kredenv setup
```

**Flags**

| Flag          | Description                                     |
| ------------- | ----------------------------------------------- |
| `--overwrite` | Re-encrypt the vault with a new master password |
| `--nuke`      | Wipe all kredenv configuration and secrets      |

**Re-encrypting with a new password**

```bash
kredenv setup --overwrite
```

You will be prompted for the new master password. All existing secrets are re-encrypted automatically.

**Wiping everything**

```bash
kredenv setup --nuke
```

Permanently deletes the vault, all secrets, and the master password. You will be prompted to confirm.

---

## `kredenv init`

Creates a `kredsfile.yaml` in the current directory and interactively prompts for values for each declared secret.

```bash
kredenv init
```

If a `kredsfile.yaml` already exists, kredenv uses it as-is and only prompts for secrets not yet stored in the vault.

**Flags**

| Flag              | Description                                                       |
| ----------------- | ----------------------------------------------------------------- |
| `--overwrite`     | Overwrite the existing `kredsfile.yaml` with the minimal template |
| `--no-setup`      | Skip the interactive secret prompting after creating the file     |
| `-f, --file`      | Path to the manifest file (default: `kredsfile.yaml`)             |
| `-n, --namespace` | Only prompt for secrets in a specific namespace                   |

**Examples**

```bash
# initialize and prompt for all secrets
kredenv init

# create the file only, skip prompting
kredenv init --no-setup

# prompt only for staging namespace secrets
kredenv init -n staging
```

---

## `kredenv hook`

Emits the shell hook script for the specified shell. Pipe it into your shell configuration to enable automatic secret loading on directory change.

```bash
kredenv hook <shell>
```

**Supported shells:** `bash`, `zsh`, `fish`, `powershell`, `nushell`

**Examples**

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

See [Shell Hooks](/docs/shell-hooks) for detailed setup instructions per shell.
