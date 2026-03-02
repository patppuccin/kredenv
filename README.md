![banner](assets/banner.svg)

[![Latest Release](https://img.shields.io/github/v/release/patppuccin/kredenv)](https://github.com/patppuccin/kredenv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/patppuccin/kredenv)](https://goreportcard.com/report/github.com/patppuccin/kredenv)
[![License](https://img.shields.io/github/license/patppuccin/kredenv)](LICENSE)

Inject secrets and environment variables from a locally encrypted vault into your shell session. No plaintext files. No accidental commits. No secrets leakage.

## How it works

*kredenv* reads a `.kredsfile` in your project directory and maps each entry to a secret stored in your locally encrypted vault.

When you run `kredenv setup`, you create a master password. That password derives a 256-bit key using Argon2 and decrypts the vault using AES-256-GCM.

When you `cd` into a directory containing a `.kredsfile`, kredenv:

- Decrypts the vault in memory
- Resolves required secrets
- Injects them into your shell session

When you leave the directory or project scope, the variables are unloaded. Secrets are never written to disk in plaintext. They never leave your machine unless you explicitly export them.

Each developer maintains their own encrypted vault. There is no shared secrets file, no `.env` to `.gitignore`, and no risk of committing credentials. The `.kredsfile` must be committed so collaborators know which secrets are required and can populate their own vault securely.

## Installation

**Convenience script (recommended):**

_Linux & macOS:_

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.sh)
```

_Windows (PowerShell 5.1+):_

```powershell
powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.ps1 | iex"
```
The script downloads the latest release and installs it into your PATH.

**Via prebuilt binary:**

Download the appropriate binary for your platform and architecture from the [releases page](https://github.com/patppuccin/kredenv/releases).

## Setup

kredenv uses a shell hook to automatically load and unload secrets as you move between directories.

Add the appropriate hook to your shell configuration and restart your session.

**Bash**

```sh
echo 'eval "$(kredenv hook bash)"' >> ~/.bashrc
```

Zsh

```sh
echo 'eval "$(kredenv hook zsh)"' >> ~/.zshrc
```

Fish

```sh
echo 'kredenv hook fish | source' >> $__fish_config_dir/config.fish
```

PowerShell

```powershell
Add-Content $PROFILE 'Invoke-Expression (& { (kredenv hook powershell | Out-String) })'
```

Nushell

```sh
kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")
```

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
