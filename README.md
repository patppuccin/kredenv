![banner](assets/banner.svg)

[![Latest Release](https://img.shields.io/github/v/release/patppuccin/kredenv)](https://github.com/patppuccin/kredenv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/patppuccin/kredenv)](https://goreportcard.com/report/github.com/patppuccin/kredenv)
[![License](https://img.shields.io/github/license/patppuccin/kredenv)](LICENSE)

Inject secrets and environment variables from a locally encrypted vault into your shell session.  
No plaintext files. No accidental commits. No secret leaks.

## Quick Start

Get `kredenv` running in under a minute.

### 1. Initialize your vault

```sh
kredenv setup
```

This walks you through creating a master password and setting up your encrypted vault.

### 2. Initialize a project

Inside your project directory:

```sh
kredenv init
```

This creates a `.kredsfile` that defines which secrets the project expects.

### 3. Store a secret

```sh
kredenv set API_KEY
```

You will be prompted to enter the secret value, which is stored in your encrypted vault.

Add the secret to your `.kredsfile`:

```txt
# .kredsfile
needs API_KEY
```

### 4. Use the secret

Once a `.kredsfile` is in scope, `kredenv` automatically injects secrets into your shell environment.

```sh
echo $API_KEY
```

If the secret was just created and the directory has not changed, you can manually load the secret:

```sh
kredenv load
```

Applications running in that directory can now read the variable as a normal environment variable.

### 5. Run commands with secrets

You can also run commands without loading secrets into your interactive shell:

```sh
kredenv exec -- npm run dev
kredenv exec -- terraform apply
```

Secrets are injected only for the lifetime of the executed command.

---

## Up and running

### Installation

`kredenv` is distributed as a single standalone binary and runs on **Linux, macOS, and Windows**.
No runtime dependencies. No package managers required.

You can install it using the convenience installer or by downloading a prebuilt release.

#### Convenience Script (Recommended)

The install script detects your platform and architecture, downloads the **latest release**, and installs the binary into your system `PATH`.

**Linux & macOS**

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.sh)
```

**Windows (PowerShell 5.1+)**

```powershell
powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.ps1 | iex"
```

#### Install via Prebuilt Binary

You can also install `kredenv` manually.

1. Download the appropriate binary for your **OS and architecture** from the
   [GitHub Releases](https://github.com/patppuccin/kredenv/releases) page.
2. Extract the archive.
3. Move the binary to a directory in your `PATH`.

Example (Linux / macOS):

```bash
chmod +x kredenv
sudo mv kredenv /usr/local/bin/
```

Example (Windows):

Move `kredenv.exe` to a directory included in your `PATH`.

> [!TIP]
> Run the following command to confirm the installation:
> ```bash
> kredenv --version
> ```
> If the command prints the version, the installation succeeded.

### Setting up Shell Integration

`kredenv` installs a shell hook that watches for directory changes and loads or unloads secrets based on the `.kredsfile` currently in scope.

**Supported shells:**

* Bash
* Zsh
* Fish
* PowerShell
* Nushell

Add the appropriate hook to your shell configuration file, then restart your shell session.

#### Bash

```sh
echo 'eval "$(kredenv hook bash)"' >> ~/.bashrc
```

#### Zsh

```sh
echo 'eval "$(kredenv hook zsh)"' >> ~/.zshrc
```

#### Fish

```sh
echo 'kredenv hook fish | source' >> $__fish_config_dir/config.fish
```

#### PowerShell

```powershell
Add-Content $PROFILE 'Invoke-Expression (& { (kredenv hook powershell | Out-String) })'
```

#### Nushell

> [!NOTE]
> Nushell loads the hook from a saved script. After upgrading `kredenv`, run the command below again to regenerate the hook.

```sh
kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")
```

### Setting up the kredenv Vault

`kredenv` stores secrets in a **locally encrypted vault**.

During setup you create a **master password** used to encrypt and decrypt the vault. The password is stored securely on the machine and reused automatically, so you are not prompted on every command.

Setup is typically performed **once per machine**.

```sh
# Initialize kredenv on this machine
kredenv setup
```

The command will:

* Prompt you to create a master password
* Create the encrypted vault
* Store the password securely on the machine

#### Recover Local Credentials

If the vault exists but the locally stored password is missing, you can restore the configuration without deleting secrets.

```sh
kredenv setup --overwrite
```

You will be prompted for the **existing master password** and then asked to set a **new one**. All secrets are re-encrypted using the new password.

#### Start Fresh

To delete all configuration and secrets and start from scratch:

```sh
kredenv setup --nuke
```

This removes:

* The encrypted vault
* All stored credentials
* Local kredenv configuration

You will be prompted for confirmation before deletion.

### Usage

`kredenv` provides a small set of commands grouped by purpose.  
Use the `-h` or `--help` flag at any level to see available commands and flags.

**For example:**

```sh
kredenv --help
kredenv set --help
kredenv exec --help
```
**Primary commands:**

```txt
Setup Commands
  kredenv setup          : Initialize kredenv on this machine and create the encrypted vault
  kredenv init           : Initialize a .kredsfile and optionally populate missing secrets
  kredenv hook <shell>   : Emit the shell integration script for the specified shell

Environment Commands
  kredenv load           : Load secrets from the .kredsfile in scope into the environment
  kredenv unload         : Remove kredenv secrets from the current shell session
  kredenv exec <command> : Execute a command with secrets injected into its environment
  kredenv which          : Print the path to the .kredsfile currently in scope
  kredenv validate       : Validate .kredsfile syntax

Secrets Commands
  kredenv set <key>      : Store a secret in the encrypted vault
  kredenv get <key>      : Retrieve a secret from the encrypted vault
  kredenv list           : List secrets defined in the .kredsfile or stored in the vault
  kredenv delete <key>   : Delete one or more secrets from the encrypted vault
  kredenv export         : Export secrets from the vault to stdout or a file
  kredenv import <file>  : Import secrets from a file into the encrypted vault
```

---

## How kredenv Works

`kredenv` follows a **local‑first secret management model**. Secrets live in an encrypted vault on your machine and are injected into your environment only when required.

### Encrypted Local Vault

Running `kredenv setup` initializes the local vault.

* A master password is created
* Argon2 derives a 256‑bit encryption key from that password
* Secrets are encrypted using **AES‑256‑GCM**

The vault is stored locally on your machine. Secret values are never written to disk in plaintext.

### Project Secret Requirements

Each project defines required secrets using a manifest file called a `.kredsfile`. This file contains **only secret names and rules**, never secret values.

Example:

```txt
needs AWS_ACCESS_KEY_ID
needs AWS_SECRET_ACCESS_KEY
maybe ANALYTICS_ID
```

The `.kredsfile` is safe to commit to version control. It acts as a contract describing which secrets developers must populate in their local vault.

### Per‑Developer Vault

Each developer maintains their own encrypted vault on their machine.

The repository contains only the `.kredsfile`, which defines which secrets are required. Each developer populates their vault independently.

This design avoids:

* committing secrets to version control
* distributing plaintext `.env` files
* sharing credentials across machines

Secrets remain local to each developer while the project still defines what is required.

### Directory Discovery

When inside a project directory, `kredenv` searches for the nearest `.kredsfile` by walking up the directory tree.

The recursion depth can be controlled inside the file:

```txt
recurse to 3
```

This allows nested project structures while keeping lookup predictable.

### Namespaces

Namespaces allow different environments (such as staging or production) to store separate values for the same secret name.

Example:

```sh
kredenv set AWS_ACCESS_KEY_ID -n staging
kredenv set AWS_ACCESS_KEY_ID -n production
```

Inside a `.kredsfile`, namespaced secrets must declare an alias:

```txt
needs staging:AWS_ACCESS_KEY_ID as AWS_ACCESS_KEY_ID
needs staging:AWS_SECRET_ACCESS_KEY as AWS_SECRET_ACCESS_KEY

needs production:AWS_ACCESS_KEY_ID as AWS_ACCESS_KEY_ID
needs production:AWS_SECRET_ACCESS_KEY as AWS_SECRET_ACCESS_KEY
```

Namespaces allow teams to manage secrets for environments such as:

* staging
* production
* development

within the same encrypted vault.

### Secret Injection

Whenever your working directory changes, `kredenv` checks whether a `.kredsfile` is in scope and:

1. Locates the `.kredsfile`
2. Decrypts the vault in memory
3. Resolves required secrets
4. Injects them into the shell environment

Applications then read them as normal environment variables.

### Automatic Unloading

When leaving the project scope, `kredenv` unloads the injected variables from the shell session. This prevents secrets from leaking into unrelated commands or projects.

Secrets exist only in memory while the project scope is active.

### Direct Command Execution

Sometimes you may want to run a command with secrets **without loading them into the interactive shell**.

`kredenv exec` runs a command in a temporary environment populated with secrets from the `.kredsfile`.

Example:

```sh
kredenv exec -- terraform apply
kredenv exec -- npm run dev
kredenv exec -n staging -- rails db:migrate
```

This approach is useful for:

* running build tools
* executing deployment commands
* scripting workflows

Secrets are injected only for the lifetime of the executed process.

### Autoloading

By default, `kredenv` automatically loads secrets when a `.kredsfile` is in scope. This behavior can be controlled using the autoload directive inside the `.kredsfile`.

```txt
# enable autoloading (either can be used - default behavior)
autoload
autoload on

# disable autoloading
autoload off
```

The directive can also be used to set the default namespace to load:

```txt
# load secrets from the staging namespace by default
autoload for staging
```

> [!NOTE]
> The autoload directive uses the default namespace when running `kredenv exec`. But a specific namespace can be specified using the `-n` or the `--namespace` flag.

---

## Export & Import

kredenv can export secrets to a file for backup or migration and re-import them on another machine. By default, it prints to stdout, but by using the `-o` flag, you can export to a file. YAML, JSON, TOML, and standard `.env` formats are supported.

```sh
# export to stdout (env format, default)
kredenv export

# export to a file
kredenv export -o .env

# export a specific namespace
kredenv export -n staging

# export multiple namespaces
kredenv export -n staging -n production

# export as json, yaml, or toml
kredenv export -f yaml

# export with value-level encryption
kredenv export --encrypt
```

When exporting multiple namespaces, env files are written per namespace (`.env.staging`, `.env.production`). Structured formats (json, yaml, toml) write a single file with namespaces as top-level keys.

Importing restores secrets to the vault and optionally updates or creates a `.kredsfile`.

```sh
# import from a file
kredenv import .env

# import only a specific namespace
kredenv import creds.yaml -n staging

# overwrite existing keys
kredenv import .env --overwrite

# import to keyring only, skip .kredsfile update
kredenv import .env --no-kredsfile
```

---

## Integrations

### Terminal Prompts

`kredenv` exposes environment variables that prompt frameworks can use to display when secrets are loaded.

Available variables:

* `KREDENV_LOADED_COUNT` – number of secrets loaded in the current scope
* `KREDENV_LOADED_VARS` – comma-separated list of secrets loaded in the current scope

> [!NOTE]
> Any prompt framework that supports environment variables can integrate with `kredenv`.  
> Use `KREDENV_LOADED_COUNT` or `KREDENV_LOADED_VARS` to display indicators, counts, or secret names in your prompt.

The following example uses the `env_var` module of [Starship](https://starship.rs) to display a prompt indicator when secrets are loaded. More information about the `env_var` module can be found [here](https://starship.rs/config/#environment-variable).

```toml
format = """\
$username \
on $hostname \
at $directory \
(with $git_branch )\
(having ${env_var.kredenv})\
$line_break\
$character\
"""

# ... other config ...

[env_var.kredenv]
variable = "KREDENV_LOADED_COUNT"
format = "[$symbol $env_value]($style)"
symbol = "🔑"
style = "bold yellow"
disabled = false
```

This displays an indicator showing the number of secrets loaded whenever a `.kredsfile` is in scope.

## Caveats

**Local-first:** kredenv stores secrets in a locally encrypted vault. It is designed for developer machines and does not provide centralized secret management. Losing your master password means losing access to the vault unless you have exported or backed up your secrets.

**CI / Production:** kredenv is not a remote secrets manager. For CI pipelines, use your platform’s native secret injection (GitHub Actions secrets, GitLab CI variables, etc.). For production systems, use a dedicated secrets manager such as HashiCorp Vault or AWS SSM.

**Per-machine:** The encrypted vault is per-user, per-machine. Each developer must run `kredenv setup` to create their vault and populate required secrets.

## Acknowledgements

kredenv is built on the shoulders of these open source projects:

- [**fatih/color**](https://github.com/fatih/color) - For dead-simple terminal styling
- [**spf13/cobra**](https://github.com/spf13/cobra) - For an intuitive CLI framework
- [**zalando/go-keyring**](https://github.com/zalando/go-keyring) - The keyring API with cross-platform compatibility
- [**mattn/go-isatty**](https://github.com/mattn/go-isatty) - For checking if a terminal is interactive

And to the Go team, for designing a language where a tool like this can go from idea to working binary in a weekend.

## A Note on LLMs

_kredenv was not vibe coded_. LLMs were used during development as a thinking partner to explore architecture ideas, critique design decisions, and research shell edge cases.

The code, design decisions, and trade-offs are the author's own.

## Roadmap

- [ ] Expand automated test coverage
- [ ] IDE support for `.kredsfile` (syntax highlighting, linting, autocomplete)
- [ ] Explore secure backup and sync mechanisms for the encrypted vault

## Contributing

Contributions are not being accepted at this time while the project structure and roadmap are still evolving.

Issues are welcome for:

- bug reports
- feature ideas
- feedback

Pull requests may be opened once the project stabilizes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) for details.