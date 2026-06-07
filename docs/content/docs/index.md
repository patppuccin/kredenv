# Introduction

**kredenv** keeps your secrets encrypted locally and injects them into your shell as you move between projects.

No plaintext `.env` files. No accidental commits. No secret leaks.

## What it does

When you `cd` into a project, kredenv reads the `kredsfile.yaml` manifest, looks up the declared secrets in your encrypted local vault, and injects them as environment variables into your shell session. When you leave the project directory, they're unloaded automatically.

For commands that need secrets without polluting your interactive shell, `kredenv exec` runs any command in a temporary environment populated with the right secrets.

## Why kredenv

Most developers manage secrets one of two ways — `.env` files committed by accident, or complex remote secrets managers that are overkill for local development. kredenv sits in between:

- **Local-first** — secrets live on your machine, encrypted. Nothing goes to a server.
- **Shell-native** — secrets load and unload as you move between projects, just like `direnv` but for secrets.
- **Declarative** — `kredsfile.yaml` is a manifest of what a project needs. Safe to commit. Easy to read.
- **Single binary** — no daemon, no runtime dependencies, no installation ceremony.

## How it compares

|                         | kredenv | `.env` files | Remote secrets manager |
| ----------------------- | ------- | ------------ | ---------------------- |
| Encrypted at rest       | ✓       | ✗            | ✓                      |
| Safe to commit manifest | ✓       | ✗            | ✓                      |
| Works offline           | ✓       | ✓            | ✗                      |
| Shell integration       | ✓       | via direnv   | ✗                      |
| Per-developer vault     | ✓       | ✗            | ✗                      |
| Zero infrastructure     | ✓       | ✓            | ✗                      |

## Getting started

If you're new to kredenv, start with [Installation](/docs/installation) then follow the [Quick Start](/docs/quick-start) guide to get up and running in under a minute.
