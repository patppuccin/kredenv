# Environment Commands

Commands for loading, unloading, and executing with secrets in your shell environment.

## `kredenv which`

Shows paths to kredenv-managed files. Useful for debugging or scripting against kredenv's file locations.

```bash
kredenv which manifest   # path to the kredsfile.yaml in scope
kredenv which store      # path to the encrypted secrets store
kredenv which creds      # path to the credentials file or keyring in use
```

---

## `kredenv load`

Shows the secrets currently loaded in your shell session. Requires the shell hook to be active.

```bash
kredenv load
```

::: tip
The actual loading is performed by the shell hook interceptor. `kredenv load` displays what was loaded.
:::

**Flags**

| Flag              | Description                            |
| ----------------- | -------------------------------------- |
| `-n, --namespace` | Load secrets from a specific namespace |

---

## `kredenv unload`

Unloads kredenv secrets from the current shell session. Requires the shell hook to be active.

```bash
kredenv unload
```

---

## `kredenv exec`

Executes a command with secrets injected into its environment. Secrets exist only for the lifetime of the command and are never loaded into your interactive shell.

```bash
kredenv exec -- <command> [args...]
```

**Flags**

| Flag              | Description                                    |
| ----------------- | ---------------------------------------------- |
| `-n, --namespace` | Execute with secrets from a specific namespace |

**Examples**

```bash
# run with flat secrets
kredenv exec -- node server.js

# run with secrets from a specific namespace
kredenv exec -n staging -- terraform plan
kredenv exec -n production -- kubectl apply -f deploy.yaml

# pipe friendly
kredenv exec -- env | grep AWS
```

**Namespace resolution**

If `--namespace` is not passed, kredenv falls back to `autoload_namespace` from the `kredsfile.yaml`. If neither is set, only flat secrets (no namespace) are injected.

---

## `kredenv validate`

Validates the `kredsfile.yaml` in scope for syntax errors and structural issues.

```bash
kredenv validate
```

You can also validate a specific file:

```bash
kredenv validate ./path/to/kredsfile.yaml
```

kredenv checks for:

- Missing `key` fields
- Duplicate secrets (same key and namespace)
- `autoload_namespace` referencing a namespace with no matching secrets
