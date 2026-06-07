# How it Works

kredenv follows a local-first secret management model. Secrets live in an encrypted vault on your machine and are injected into your shell environment only when required.

## The Flow

```
kredsfile.yaml          local vault
(what's needed)   →   (encrypted secrets)   →   shell environment
```

1. You declare what secrets a project needs in `kredsfile.yaml`
2. Each developer stores their own secret values in their local vault
3. When you enter a project directory, kredenv reads the manifest, decrypts the vault in memory, and injects the declared secrets as environment variables
4. When you leave the directory, the secrets are unloaded from your shell

## Encryption

Running `kredenv setup` initializes the local vault:

- You create a **master password**
- **Argon2id** derives a 256-bit encryption key from that password
- Secrets are encrypted using **AES-256-GCM**

The vault is stored locally on your machine. Secret values are never written to disk in plaintext. The master password is stored in your OS keyring (Keychain on macOS, Credential Manager on Windows, Secret Service on Linux) with a file-based fallback for headless environments.

## Directory Discovery

When your shell detects a directory change, kredenv searches for the nearest `kredsfile.yaml` by walking up the directory tree from the current working directory.

The recursion depth is controlled by the `recurse` field in the manifest:

```yaml
recurse: 3 # walk up to 3 levels looking for a kredsfile.yaml
```

This allows nested project structures while keeping secret lookup predictable.

## Secret Injection

Whenever your working directory changes, the shell hook:

1. Unloads any previously loaded secrets
2. Locates the nearest `kredsfile.yaml`
3. Checks if `autoload: true` is set
4. Decrypts the vault in memory
5. Resolves the declared secrets
6. Injects them into the shell environment

The vault is decrypted in memory only — never written to disk unencrypted.

## Automatic Unloading

When you leave a project directory, kredenv unloads the injected variables from your shell session. This prevents secrets from leaking into unrelated commands or other projects.

Secrets exist in your environment only while the project scope is active.

## Per-Developer Vault

Each developer maintains their own encrypted vault on their machine. The repository contains only the `kredsfile.yaml`, which declares which secrets are required. Each developer populates their vault independently.

This design avoids committing secrets to version control, distributing plaintext `.env` files, or sharing credentials across machines.
