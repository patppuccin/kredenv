---
title: Roadmap
---

# Roadmap

A rough look at what's been built and what's coming.

## Done

- Encrypted local vault with AES-256-GCM and argon2id key derivation
- Shell hooks for bash, zsh, fish, nushell, and PowerShell
- Automatic secret loading and unloading on directory change
- `kredsfile.yaml` declarative manifest
- Namespace support for multi-environment secrets
- `kredenv exec` for scoped command execution
- Export and import in env, json, yaml, and toml formats
- Value-level encryption on export
- Prompt integration with Starship
- OS keyring integration with file-based fallback
- Single binary, cross-platform — Linux, macOS, Windows
- VitePress documentation site

## Planned

- JSON schema for `kredsfile.yaml` with IDE intellisense
- External secrets store support (AWS SSM, HashiCorp Vault)
- Automated test coverage
- Package manager distribution (Homebrew, Scoop, Winget, AUR)
- Secure vault backup and sync
