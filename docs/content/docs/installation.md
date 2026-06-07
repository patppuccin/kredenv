# Installation

kredenv is distributed as a single standalone binary. No runtime dependencies, no package managers required. It runs on **Linux**, **macOS**, and **Windows**.

## Convenience Script (Recommended)

The install script detects your platform and architecture, downloads the latest release, and installs the binary into your system `PATH`.

**Linux & macOS**

```bash
bash <(curl -fsSL https://kredenv.patppuccin.com/install.sh)
```

**Windows (PowerShell 5.1+)**

```powershell
powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://kredenv.patppuccin.com/install.ps1 | iex"
```

## Manual Installation

1. Download the appropriate binary for your OS and architecture from the [GitHub Releases](https://github.com/patppuccin/kredenv/releases) page.
2. Extract the archive.
3. Move the binary to a directory in your `PATH`.

**Linux & macOS**

```bash
chmod +x kredenv
sudo mv kredenv /usr/local/bin/
```

**Windows**

Move `kredenv.exe` to a directory included in your `PATH`.

## Verify Installation

```bash
kredenv --version
```

If the command prints the version, the installation succeeded.

## Next Steps

Once installed, set up your vault and shell integration by following the [Quick Start](/docs/quick-start) guide.
