# Caveats

Things to be aware of before relying on kredenv.

## Local-first

kredenv stores secrets in a locally encrypted vault. It is designed for developer machines and does not provide centralized secret management.

**Losing your master password means losing access to your vault.** There is no recovery mechanism. Export your secrets regularly:

```bash
kredenv export -f yaml --encrypt -o backup.yaml
```

## Per-machine

The encrypted vault is per-user, per-machine. Each developer must run `kredenv setup` on their own machine and populate required secrets independently.

There is no built-in sync between machines. Use `kredenv export` and `kredenv import` to migrate secrets to a new machine.

## CI / Production

kredenv is not a remote secrets manager. It is not designed for CI pipelines or production systems.

For CI, use your platform's native secret injection:

- **GitHub Actions** — encrypted repository secrets
- **GitLab CI** — CI/CD variables
- **Bitbucket Pipelines** — repository variables

For production, use a dedicated secrets manager such as HashiCorp Vault, AWS SSM, or Azure Key Vault.

## Shell history

When using `kredenv set <key> <value>` with the value as a direct argument, the value may appear in your shell history. Use the interactive prompt instead:

```bash
kredenv set MY_SECRET
# kredenv prompts for the value with masked input
```

## autoload and performance

The shell hook runs on every directory change. On most systems this is imperceptible. On very slow filesystems or with very large vaults it may introduce a slight delay. Use `autoload: false` and `kredenv exec` if this is a concern.

## Comments in kredsfile.yaml

Running `kredenv import` with kredsfile updates enabled will rewrite your `kredsfile.yaml` without preserving comments. Edit the file manually to maintain comments and formatting.
