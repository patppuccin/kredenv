![banner](assets/banner.svg)

[![Latest Release](https://img.shields.io/github/v/release/patppuccin/kredenv)](https://github.com/patppuccin/kredenv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/patppuccin/kredenv)](https://goreportcard.com/report/github.com/patppuccin/kredenv)
[![License](https://img.shields.io/github/license/patppuccin/kredenv)](LICENSE)

Inject secrets and environment variables from your OS keyring into your shell session. No plaintext files. No accidental commits.

## How it works

kredenv reads a `.kredsfile` in your project directory. Each entry maps an environment variable name to a key stored in your OS keyring. When you `cd` into the directory, kredenv resolves the keys from the keyring and injects them into your shell session. When you leave, they are unloaded.

Secrets live in your OS keyring, which is _Keychain_ on _macOS_, _Credential Manager_ on _Windows_, _libsecret_ on _Linux_. They are never written to disk in plaintext and never leave your machine unless you explicitly export them.

Each developer manages their own keyring. There is no shared secrets file, no `.env` to `.gitignore`, no risk of accidentally committing credentials. As a matter of fact, `.kredsfile` must be committed to your repository, so that collaborators can use it to setup secrets in their own keyring easily.

## Installation

**Via convenience script (recommended):**

_Linux & macOS:_

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.sh)
```

_Windows (PowerShell 5.1+):_

```powershell
powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.ps1 | iex"
```

**Via prebuilt binary:**

Download the latest release for your platform from the [releases page](https://github.com/patppuccin/kredenv/releases).

## Setup

kredenv requires a shell hook to inject and unload secrets automatically as you navigate directories. Run the appropriate command once and restart your shell session to enable it.

**Bash:**

```sh
echo 'eval "$(kredenv hook bash)"' >> ~/.bashrc
```

**Zsh:**

```sh
echo 'eval "$(kredenv hook zsh)"' >> ~/.zshrc
```

**Fish:**

```sh
echo 'kredenv hook fish | source' >> $__fish_config_dir/config.fish
```

**Nushell:**

```sh
kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")
```

**PowerShell:**

```powershell
Add-Content $PROFILE 'Invoke-Expression (& { (kredenv hook powershell | Out-String) })'
```

## The `.kredsfile`

Place a `.kredsfile` in your project root. kredenv walks up the directory tree to find the nearest one.

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
kredenv set AWS_SECRET_ACCESS_KEY my-super-secret-aws-secret-access-key
```

To validate your `.kredsfile`:

```sh
kredenv validate
```

## Namespaces

Namespaces let you manage secrets for multiple environments in the same keyring without collision.

```sh
# store secrets under a namespace
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

| Command                | Description                                        |
| ---------------------- | -------------------------------------------------- |
| `kredenv init`         | Initialize a `.kredsfile` in the current directory |
| `kredenv setup`        | Interactively store missing secrets in the keyring |
| `kredenv hook <shell>` | Emit the shell integration script                  |

### Environment

| Command              | Description                                              |
| -------------------- | -------------------------------------------------------- |
| `kredenv load`       | Show currently loaded secrets                            |
| `kredenv unload`     | Unload secrets from the current session                  |
| `kredenv exec <cmd>` | Run a command with secrets injected into its environment |
| `kredenv which`      | Print the path to the `.kredsfile` in scope              |
| `kredenv validate`   | Validate `.kredsfile` syntax                             |

### Keyring

| Command                 | Description                                 |
| ----------------------- | ------------------------------------------- |
| `kredenv set <key>`     | Store a secret in the keyring               |
| `kredenv get <key>`     | Retrieve a secret from the keyring          |
| `kredenv delete <key>`  | Delete a secret from the keyring            |
| `kredenv list`          | List keys defined in the `.kredsfile`       |
| `kredenv export`        | Export secrets to stdout or a file          |
| `kredenv import <file>` | Import secrets from a file into the keyring |

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

**Linux:** Requires a running keyring daemon (`gnome-keyring` or `kwallet`). Headless environments and most CI systems do not have one. kredenv is intended for local developer machines, not CI pipelines.

**CI / Production:** kredenv is not a secrets manager for shared or remote environments. For CI, use your platform's native secret injection (GitHub Actions secrets, GitLab CI variables, etc.). For production, use a dedicated secrets manager such as HashiCorp Vault, AWS SSM, or 1Password CLI.

**Per-machine:** Keyring secrets are stored per-user, per-machine. Each developer must run `kredenv setup` or `kredenv set` to populate their own keyring.

## Acknowledgements

kredenv is built on the shoulders of these open source projects:

- [**fatih/color**](https://github.com/fatih/color) - For dead-simple terminal styling
- [**spf13/cobra**](https://github.com/spf13/cobra) - For an intuitive CLI framework
- [**zalando/go-keyring**](https://github.com/zalando/go-keyring) - The keyring API which does all the heavy lifting and cross-platform compatibility
- [**mattn/go-isatty**](https://github.com/mattn/go-isatty) - For checking if a terminal is interactive

And to the Go team, for designing a language where a tool like this can go from idea to working binary in a weekend.

## A note on LLMs

_kredenv was not vibe coded_. However, LLMs were used throughout development as a thinking partner to ideate on architecture, critique design decisions, and research edge cases across shell environments. Every decision was reasoned through, understood, and deliberately chosen.

The code, the design, the opinions and most importantly noob'ish coding practices in this tool are the author's own.

## Roadmap

- [ ] Add a reasonable test suite
- [ ] Guides on integrating with terminal prompt libraries (e.g. Starship, Oh My Posh)
- [ ] IDE integrations for syntax highlighting and autocompletion
- [ ] Add support for alternative backends (e.g. local encryption, secrets managers)
- [ ] Explore a remote sync mechanism for `.kredsfile` and keyring secrets

## Contributing

Not accepting contributions at this time while the project structure and roadmap are being figured out. Feel free to open issues for bugs or ideas. Feedback is welcome even if PRs are not yet.

## License

[MIT](LICENSE)
