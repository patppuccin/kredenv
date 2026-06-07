# Secrets Commands

Commands for managing secrets in the encrypted vault.

---

## `kredenv set`

Stores a secret in the vault. If no value is provided, kredenv prompts for it interactively with masked input.

```bash
kredenv set <key> [value]
```

**Flags**

| Flag              | Description                                 |
| ----------------- | ------------------------------------------- |
| `-n, --namespace` | Store the secret under a specific namespace |

**Examples**

```bash
# prompt for value interactively
kredenv set AWS_ACCESS_KEY_ID

# provide value directly (use with caution — may appear in shell history)
kredenv set AWS_ACCESS_KEY_ID AKIAIOSFODNN7EXAMPLE

# store under a namespace
kredenv set DATABASE_PASSWORD -n staging
kredenv set DATABASE_PASSWORD -n production

# using colon syntax
kredenv set staging:DATABASE_PASSWORD
```

::: warning
Passing the value directly as an argument may expose it in your shell history. Use the interactive prompt when possible.
:::

---

## `kredenv get`

Retrieves a secret from the vault and prints it to stdout.

```bash
kredenv get <key>
```

**Flags**

| Flag              | Description                              |
| ----------------- | ---------------------------------------- |
| `-n, --namespace` | Get the secret from a specific namespace |

**Examples**

```bash
kredenv get AWS_ACCESS_KEY_ID
kredenv get DATABASE_PASSWORD -n staging

# using colon syntax
kredenv get staging:DATABASE_PASSWORD
```

---

## `kredenv delete`

Deletes one or more secrets from the vault.

```bash
kredenv delete <key> [keys...]
```

**Flags**

| Flag              | Description                           |
| ----------------- | ------------------------------------- |
| `-n, --namespace` | Delete keys from a specific namespace |

**Examples**

```bash
# delete a single key
kredenv delete AWS_ACCESS_KEY_ID

# delete multiple keys
kredenv delete AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY

# delete a namespaced key
kredenv delete DATABASE_PASSWORD -n staging

# using colon syntax
kredenv delete staging:DATABASE_PASSWORD
```

If a key is not found, kredenv logs the error and continues to the next key.

---

## `kredenv list`

Lists secrets declared in the `kredsfile.yaml` and checks which ones are set in the vault.

```bash
kredenv list
```

**Flags**

| Flag              | Description                                           |
| ----------------- | ----------------------------------------------------- |
| `-a, --all`       | List all secrets in the vault, ignoring the kredsfile |
| `--show-values`   | Show secret values (use with caution)                 |
| `-n, --namespace` | Filter by namespace                                   |

**Examples**

```bash
# list secrets from kredsfile.yaml and their vault status
kredenv list

# list all secrets in the vault
kredenv list --all

# filter by namespace
kredenv list -n staging

# show values
kredenv list --show-values
```

---

## `kredenv export`

Exports secrets from the vault to stdout or a file. Supports `env`, `json`, `yaml`, and `toml` formats.

```bash
kredenv export
```

**Flags**

| Flag               | Description                                                   |
| ------------------ | ------------------------------------------------------------- |
| `-f, --format`     | Export format: `env`, `json`, `yaml`, `toml` (default: `env`) |
| `-o, --output`     | Output path (default: stdout)                                 |
| `--all`            | Export all secrets in the vault                               |
| `--encrypt`        | Encrypt secret values with a password                         |
| `-n, --namespaces` | Export specific namespaces (repeatable)                       |

**Examples**

```bash
# export to stdout
kredenv export

# export to a file
kredenv export -o backup.env

# export as yaml
kredenv export -f yaml -o backup.yaml

# export a specific namespace
kredenv export -n staging

# export multiple namespaces
kredenv export -n staging -n production

# export with value-level encryption
kredenv export --encrypt -o backup.yaml
```

When exporting multiple namespaces as `env`, kredenv writes one file per namespace (`.env.staging`, `.env.production`). Structured formats write a single file with namespaces as top-level keys.

---

## `kredenv import`

Imports secrets from a file into the vault. Supports `env`, `json`, `yaml`, and `toml` formats.

```bash
kredenv import <file>
```

**Flags**

| Flag               | Description                                           |
| ------------------ | ----------------------------------------------------- |
| `--overwrite`      | Overwrite existing keys in the vault                  |
| `-n, --namespaces` | Import specific namespaces from the file (repeatable) |

**Examples**

```bash
# import from an env file
kredenv import .env

# import from a namespaced env file
kredenv import .env.staging

# import from yaml
kredenv import backup.yaml

# overwrite existing keys
kredenv import backup.yaml --overwrite

# import a specific namespace
kredenv import backup.yaml -n staging
```

After importing, kredenv prints hints for any secrets not yet declared in your `kredsfile.yaml`.
