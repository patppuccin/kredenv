# Kredsfile Manifest

The `kredsfile.yaml` is a declarative manifest that defines what secrets a project needs. It contains only secret names — never values. It is safe to commit to version control.

## Location

kredenv looks for `kredsfile.yaml` in the current directory and walks up the directory tree based on the `recurse` setting.

## Full Reference

```yaml
# kredsfile.yaml
# safe to commit - contains no secrets

recurse: 3 # walk up N levels looking for a kredsfile.yaml
autoload: true # inject secrets into shell on cd (default: false)
autoload_namespace: staging # namespace to autoload (default: secrets without namespace)

secrets:
  - key: AWS_ACCESS_KEY_ID

  - key: DATABASE_PASSWORD
    namespace: production

  - key: GOOGLE_ANALYTICS_ID
```

## Fields

### `recurse`

Controls how many directory levels kredenv walks up when searching for a `kredsfile.yaml`.

```yaml
recurse: 3
```

Default is `0` — only the current directory is checked. Set to a positive integer to allow kredenv to find the manifest from subdirectories.

### `autoload`

When `true`, the shell hook automatically injects secrets when you `cd` into the project directory.

```yaml
autoload: true
```

Default is `false` — secrets must be loaded manually with `kredenv load` or injected via `kredenv exec`.

### `autoload_namespace`

Sets which namespace to load automatically. When set, only secrets from that namespace are injected. When absent, only flat secrets (no namespace) are injected.

```yaml
autoload_namespace: staging
```

### `secrets`

A list of secrets the project requires. Each entry has:

| Field       | Required | Description                         |
| ----------- | -------- | ----------------------------------- |
| `key`       | Yes      | The secret name in the vault        |
| `namespace` | No       | The vault namespace for this secret |

**Flat secret** — stored and injected as `AWS_ACCESS_KEY_ID`:

```yaml
secrets:
  - key: AWS_ACCESS_KEY_ID
```

**Namespaced secret** — stored as `staging:DATABASE_PASSWORD`, injected as `DATABASE_PASSWORD`:

```yaml
secrets:
  - key: DATABASE_PASSWORD
    namespace: staging
```

## Validation

Run `kredenv validate` to check your `kredsfile.yaml` for errors:

```bash
kredenv validate
```

kredenv errors on:

- Missing `key` field on any secret entry
- Duplicate secrets (same key and namespace declared twice)
- `autoload_namespace` set to a namespace with no matching secrets

## Initialization

Run `kredenv init` to create a `kredsfile.yaml` in the current directory and interactively populate the declared secrets:

```bash
kredenv init
```
