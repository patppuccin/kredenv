---
layout: home

hero:
  name: KREDENV
  text: Secrets Manager
  tagline: Keep your secrets encrypted locally and inject them into your shell as you move across projects.
  image:
    light: /hero-light.svg
    dark: /hero-dark.svg
    alt: kredenv hero image
  actions:
    - theme: brand
      text: Get Started
      link: /docs/
    - theme: alt
      text: View on GitHub
      link: https://github.com/patppuccin/kredenv

features:
  - icon: 🔐
    title: Encrypted vault
    details: AES-256-GCM encryption with argon2id key derivation. Secrets never touch disk unencrypted.
  - icon: 🐚
    title: Shell-native
    details: Shell hooks for bash, zsh, fish, nushell, and PowerShell that load and unload secrets automatically as you move between projects.
  - icon: 📦
    title: Single binary
    details: One binary, no runtime dependencies, no daemon. Drop it on any machine and go.
  - icon: 📋
    title: kredsfile.yaml
    details: A declarative manifest that defines what secrets a project needs. Safe to commit, easy to read.
---
