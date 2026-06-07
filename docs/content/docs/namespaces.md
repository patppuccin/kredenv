# Namespaces

Namespaces let you store multiple values for the same secret key — one per environment — within the same vault.

## The Problem

When working across multiple environments, you often need the same secret name to hold different values:

- `DATABASE_PASSWORD` for staging points to one database
- `DATABASE_PASSWORD` for production points to another

Without namespaces, you'd need to manually swap values or maintain separate vaults.

## How Namespaces Work

A namespace is a prefix on the vault key, separated by a colon:

```
staging:DATABASE_PASSWORD
production:DATABASE_PASSWORD
```

Both are stored in the same vault, independently. You load one at a time.

## Storing Namespaced Secrets

```bash
kredenv set DATABASE_PASSWORD -n staging
kredenv set DATABASE_PASSWORD -n production
```

Or using the colon syntax directly:

```bash
kredenv set staging:DATABASE_PASSWORD
kredenv set production:DATABASE_PASSWORD
```

## Declaring in kredsfile.yaml

```yaml
secrets:
  - key: DATABASE_PASSWORD
    namespace: staging

  - key: DATABASE_PASSWORD
    namespace: production

  - key: API_KEY
    namespace: staging

  - key: API_KEY
    namespace: production
```

## Loading a Namespace

**Via autoload_namespace** — set a default namespace in your kredsfile:

```yaml
autoload: true
autoload_namespace: staging
```

kredenv will inject all `staging:*` secrets automatically on `cd`.

**Via exec** — inject a specific namespace for a single command:

```bash
kredenv exec -n staging -- terraform plan
kredenv exec -n production -- terraform apply
```

**Via load** — manually load a namespace into your shell:

```bash
kredenv load -n staging
```

## Listing Namespaced Secrets

```bash
# list secrets declared in kredsfile.yaml for staging namespace
kredenv list -n staging

# list all secrets in the vault across all namespaces
kredenv list --all
```

## Flat and Namespaced Together

You can mix flat and namespaced secrets in the same kredsfile:

```yaml
secrets:
  - key: GLOBAL_API_KEY # flat, always loaded

  - key: DATABASE_PASSWORD
    namespace: staging # only loaded when namespace is staging

  - key: DATABASE_PASSWORD
    namespace: production # only loaded when namespace is production
```

When `autoload_namespace` is empty, only flat secrets are autoloaded. Namespaced secrets require an explicit namespace to be set.
