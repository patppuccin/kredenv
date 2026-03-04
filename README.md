![banner](assets/banner.svg)

[![Latest Release](https://img.shields.io/github/v/release/patppuccin/kredenv)](https://github.com/patppuccin/kredenv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/patppuccin/kredenv)](https://goreportcard.com/report/github.com/patppuccin/kredenv)
[![License](https://img.shields.io/github/license/patppuccin/kredenv)](LICENSE)

Inject secrets and environment variables from a locally encrypted vault into your shell session.  
No plaintext files. No accidental commits. No secret leaks.

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

Move `kredenv.exe` to a directory already included in your `PATH`.


> [!TIP]
> Run the following command to confirm the installation:
> ```bash
> kredenv --version
> ```
> If the command prints the version, the installation succeeded.

### Setting up Shell Integration

`kredenv` integrates with your shell to automatically **load and unload environment secrets when you enter or leave project directories**.

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

## How it works

*kredenv* reads a `.kredsfile` in your project directory and maps each entry to a secret stored in your locally encrypted vault.

When you run `kredenv setup`, you create a master password. That password derives a 256-bit key using Argon2 and decrypts the vault using AES-256-GCM.

When you `cd` into a directory containing a `.kredsfile`, kredenv:

- Decrypts the vault in memory
- Resolves required secrets
- Injects them into your shell session

When you leave the directory or project scope, the variables are unloaded. Secrets are never written to disk in plaintext. They never leave your machine unless you explicitly export them.

Each developer maintains their own encrypted vault. There is no shared secrets file, no `.env` to `.gitignore`, and no risk of committing credentials. The `.kredsfile` must be committed so collaborators know which secrets are required and can populate their own vault securely.

## The `.kredsfile`

Place a `.kredsfile` in your project root. kredenv walks up the directory tree to locate the nearest one.

```txt
# .kredsfile
# safe to commit - contains no secrets
# kredenv errors on missing 'needs', warns on missing 'maybe'

# levels to recurse when looking for .kredsfile
recurse to 3

# autoload controls what gets injected on cd
autoload on            # load flat keys (default)
autoload off           # disable autoloading
autoload for staging   # load only staging namespace keys

# mandatory secrets
needs AWS_ACCESS_KEY_ID
needs AWS_SECRET_ACCESS_KEY

# optional secrets
maybe ANALYTICS_ID

# namespaced secrets (must have an alias)
needs staging:AWS_ACCESS_KEY_ID as AWS_ACCESS_KEY_ID
needs staging:AWS_SECRET_ACCESS_KEY as AWS_SECRET_ACCESS_KEY
```

To initialize a `.kredsfile` in the current directory:

```sh
kredenv init
```

To store secrets in the keyring:

```sh
kredenv set AWS_ACCESS_KEY_ID
kredenv set AWS_SECRET_ACCESS_KEY super-secret-value
```

To validate your `.kredsfile`:

```sh
kredenv validate
```

## Namespaces

Namespaces let you manage secrets for multiple environments inside the same encrypted vault without collisions.

```sh
kredenv set AWS_ACCESS_KEY_ID -n staging
kredenv set AWS_ACCESS_KEY_ID -n production

# load a specific namespace
kredenv load -n staging

# run a command with a specific namespace
kredenv exec -n production -- terraform apply
```

In the `.kredsfile`, namespaced keys must declare an alias:

```txt
needs staging:AWS_ACCESS_KEY_ID as AWS_ACCESS_KEY_ID
needs production:AWS_ACCESS_KEY_ID as AWS_ACCESS_KEY_ID
```

The `autoload for <namespace>` directive controls which namespace gets loaded automatically on `cd`.

## Usage

### Setup

| Command                | Description                                                   |
| ---------------------- | ------------------------------------------------------------- |
| `kredenv setup`        | Initialize kredenv on this machine and create the vault       |
| `kredenv init`         | Initialize a `.kredsfile` and optionally fill missing secrets |
| `kredenv hook <shell>` | Emit the shell integration script                             |

### Environment

| Command              | Description                                                      |
| -------------------- | ---------------------------------------------------------------- |
| `kredenv load`       | Load secrets from the `.kredsfile` in scope into the environment |
| `kredenv unload`     | Unload kredenv secrets from the current session                  |
| `kredenv exec <cmd>` | Run a command with secrets injected into its environment         |
| `kredenv which`      | Print the path to the `.kredsfile` in scope                      |
| `kredenv validate`   | Validate `.kredsfile` syntax                                     |

### Secrets

| Command                 | Description                                         |
| ----------------------- | --------------------------------------------------- |
| `kredenv set <key>`     | Store a secret in the encrypted vault               |
| `kredenv get <key>`     | Retrieve a secret from the encrypted vault          |
| `kredenv list`          | List secrets defined in the `.kredsfile` or vault   |
| `kredenv delete <key>`  | Delete one or more secrets from the encrypted vault |
| `kredenv export`        | Export secrets to stdout or a file                  |
| `kredenv import <file>` | Import secrets from a file into the encrypted vault |

## Export & Import

kredenv can export secrets to a file for backup or migration and re-import them on another machine. By default, it prints to stdout, but by using the `-o` flag, you can export to a file. YAML, JSON, and TOML formats along with the standard env format are supported.

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

Importing restores secrets to the keyring and updates or creates a `.kredsfile`:

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

## Running a command with secrets

To run a single command with secrets injected without loading them into your session:

```sh
kredenv exec -- terraform apply
kredenv exec -- npm run dev
kredenv exec -n staging -- rails db:migrate
```

## Integrations

### Prompt

## Caveats

**Local-first:** kredenv stores secrets in a locally encrypted vault. It is designed for developer machines. It does not provide centralized secret management. Losing your master password means losing access to your secrets.

**CI / Production:** kredenv is not a remote secrets manager. For CI, use your platform’s native secret injection (GitHub Actions secrets, GitLab CI variables, etc.). For production systems, use a dedicated secrets manager such as HashiCorp Vault or AWS SSM.

**Per-machine:** The encrypted vault is per-user, per-machine. Each developer must run `kredenv setup` to create their vault and populate required secrets.

## Acknowledgements

kredenv is built on the shoulders of these open source projects:

- [**fatih/color**](https://github.com/fatih/color) - For dead-simple terminal styling
- [**spf13/cobra**](https://github.com/spf13/cobra) - For an intuitive CLI framework
- [**zalando/go-keyring**](https://github.com/zalando/go-keyring) - The keyring API with cross-platform compatibility
- [**mattn/go-isatty**](https://github.com/mattn/go-isatty) - For checking if a terminal is interactive

And to the Go team, for designing a language where a tool like this can go from idea to working binary in a weekend.

## A note on LLMs

_kredenv was not vibe coded_. However, LLMs were used throughout development as a thinking partner to ideate on architecture, critique design decisions, and research edge cases across shell environments. Every decision was reasoned through, understood, and deliberately chosen.

The code, the design, the opinions and most importantly noob'ish coding practices in this tool are the author's own.

## Roadmap

- [ ] Add a reasonable test suite
- [ ] Guides on integrating with terminal prompt libraries (e.g. Starship, Oh My Posh)
- [ ] IDE integrations for syntax highlighting and autocompletion
- [ ] Explore a remote sync mechanism for `.kredsfile` and keyring secrets

## Contributing

Not accepting contributions at this time while the project structure and roadmap are being figured out. Feel free to open issues for bugs or ideas. Feedback is welcome even if PRs are not yet.

## License

[MIT](LICENSE)
