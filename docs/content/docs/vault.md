# The Vault

The kredenv vault is a locally encrypted file that stores all your secret values. It lives on your machine and never leaves it.

## Location

The vault is stored in the kredenv configuration directory:

| Platform | Path                      |
| -------- | ------------------------- |
| Linux    | `~/.kredenv/`             |
| macOS    | `~/.kredenv/`             |
| Windows  | `%LOCALAPPDATA%\kredenv\` |

## Master Password

The vault is protected by a master password you create during `kredenv setup`. This password is used to derive the encryption key via Argon2id — a memory-hard key derivation function designed to resist brute-force attacks.

The master password is stored in your OS keyring after setup so you aren't prompted on every command:

- **macOS** — Keychain
- **Windows** — Credential Manager
- **Linux** — Secret Service (GNOME Keyring / KWallet)

On headless environments where no keyring is available, the password falls back to a file at `~/.kredenv/kredmaster` with `0600` permissions.

::: warning
Losing your master password means losing access to your vault. There is no recovery mechanism. Export your secrets regularly as a backup.
:::

## Encryption

Secrets are encrypted using **AES-256-GCM** with a key derived from your master password via **Argon2id**. The vault file is a single encrypted blob — there is no way to read individual secrets without the master password.

## Setup

```bash
kredenv setup
```

Initializes the vault on this machine. Run once per machine.

## Re-encrypting with a New Password

```bash
kredenv setup --overwrite
```

Re-encrypts all existing secrets with a new master password. The original password is required.

## Wiping the Vault

```bash
kredenv setup --nuke
```

Permanently deletes the vault, all stored secrets, and the master password from the keyring. You will be prompted to confirm before deletion.

::: danger
This action is irreversible. All secrets will be lost unless you have exported them first.
:::

## Backup and Recovery

kredenv does not currently provide built-in backup or sync. To back up your secrets:

```bash
kredenv export -f yaml -o backup.yaml
```

To restore on another machine after running `kredenv setup`:

```bash
kredenv import backup.yaml
```

Keep the backup file secure — it contains your plaintext secrets. Use `--encrypt` to export with value-level encryption.
