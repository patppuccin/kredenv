---
title: About
---

# About kredenv

kredenv started as a personal frustration. Secrets scattered across `.env` files, shell profiles, and password managers but none of it project-scoped and all of it tedious to manage across machines.

The goal was simple: one encrypted vault, shell hooks that load the right secrets as you move between projects, and a declarative manifest that makes it clear what a project needs.

## Philosophy

**Encrypted by default.** Secrets are encrypted with AES-256-GCM and argon2id key derivation. Nothing sensitive ever touches disk unencrypted.

**Shell-native.** kredenv integrates with your shell rather than replacing it. Hooks load and unload secrets automatically as you `cd` between projects.

**Single binary.** No runtime dependencies, no daemon, no installation wizard. Download, run, done.

**Declarative manifest.** `kredsfile.yaml` is the contract for what secrets a project needs. It is safe to commit, easy to read, and designed to be understood at a glance.

**Boring tech.** The stack is deliberately simple — proven tools that do their job without drama.

## Built With

kredenv is built with the following open source tools:

- [Go](https://go.dev) — the binary
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [termactions](https://github.com/patppuccin/termactions) — terminal UI primitives
- [go-keyring](https://github.com/zalando/go-keyring) — OS keyring integration
- [VitePress](https://vitepress.dev) — this documentation site
