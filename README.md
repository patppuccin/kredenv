# KREDENV

Inject secrets from your OS keyring into your shell environment via a `.kredsfile`.

No plaintext credentials. No `.env` files. Commit your `.kredsfile` safely.

> ⚠️ Work in progress — not ready for use yet.

## How it works

- Declare what secrets your project needs in a `.kredsfile`
- `kredenv` fetches them from your OS keyring and injects them as environment variables
- Will support major shells (bash, zsh, fish, powershell, nushell)

## Status

Very early development. API and file format subject to change.