# Quick Start

Get kredenv running in under a minute.

### Step 1: Initialize your vault

```bash
kredenv setup
```

This walks you through creating a master password and setting up your encrypted vault. Run this once per machine.

### Step 2: Set up your shell hook

Add the hook to your shell configuration file so kredenv can automatically load and unload secrets as you move between directories.

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

Restart your shell after adding the hook.

### Step 3: Initialize a project

Inside your project directory:

```bash
kredenv init
```

This creates a `kredsfile.yaml` that declares which secrets the project needs, then prompts you to store values for each one.

### Step 4: Store a secret manually

```bash
kredenv set AWS_ACCESS_KEY_ID
```

You will be prompted to enter the secret value. It is stored encrypted in your local vault.

Then declare it in your `kredsfile.yaml`:

```yaml
secrets:
  - key: AWS_ACCESS_KEY_ID
```

### Step 5: Use the secret

Once a `kredsfile.yaml` is in scope and `autoload: true` is set, kredenv injects secrets automatically when you `cd` into the directory.

```bash
echo $AWS_ACCESS_KEY_ID
```

To load manually without autoload:

```bash
kredenv load
```

### Step 6: Run commands with secrets

To inject secrets into a single command without loading them into your interactive shell:

```bash
kredenv exec -- terraform apply
kredenv exec -- npm run dev
kredenv exec -n staging -- rails db:migrate
```

Secrets exist only for the lifetime of the command.
