#!/usr/bin/env bash
# KREDENV Installation Convenience Script for Linux & macOS
# Usage: bash <(curl -fsSL https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.sh)

set -euo pipefail

REPO="patppuccin/kredenv"
BIN_NAME="kredenv"
INSTALL_DIR="$HOME/.local/bin"

log_inf() { printf "\033[34mINF\033[0m %s\n" "$*"; }
log_wrn() { printf "\033[33mWRN\033[0m %s\n" "$*"; }
log_err() { printf "\033[31mERR\033[0m %s\n" "$*" >&2; }

get_latest_version() {
    curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
        grep '"tag_name"' |
        sed 's/.*"tag_name": *"\([^"]*\)".*/\1/'
}

get_os() {
    case "$(uname -s)" in
        Linux)  echo "linux" ;;
        Darwin) echo "darwin" ;;
        *)
            log_err "Unsupported OS: $(uname -s)"
            exit 1
            ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64)  echo "x86_64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)
            log_err "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

verify_checksum() {
    local file="$1"
    local checksum_file="$2"
    local archive_name="$3"

    local expected
    expected=$(grep "$archive_name" "$checksum_file" | awk '{print $1}' | tr '[:lower:]' '[:upper:]')

    if [[ -z "$expected" ]]; then
        log_err "Checksum entry not found for $archive_name"
        exit 1
    fi

    local actual
    if command -v sha256sum &>/dev/null; then
        actual=$(sha256sum "$file" | awk '{print $1}' | tr '[:lower:]' '[:upper:]')
    elif command -v shasum &>/dev/null; then
        actual=$(shasum -a 256 "$file" | awk '{print $1}' | tr '[:lower:]' '[:upper:]')
    else
        log_err "No SHA256 utility found (sha256sum or shasum required)"
        exit 1
    fi

    if [[ "$actual" != "$expected" ]]; then
        log_err "Checksum verification failed"
        log_err "Expected: $expected"
        log_err "Actual:   $actual"
        exit 1
    fi

    log_inf "Checksum verified"
}

install_kredenv() {
    local version
    version=$(get_latest_version)
    local clean_version="${version#v}"
    local os
    os=$(get_os)
    local arch
    arch=$(get_arch)

    local archive_name="kredenv-$os-$arch.tar.gz"
    local checksum_file="kredenv_${clean_version}_checksums.txt"
    local download_url="https://github.com/$REPO/releases/download/$version/$archive_name"
    local checksum_url="https://github.com/$REPO/releases/download/$version/$checksum_file"

    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap '[[ -n "${tmp_dir:-}" ]] && rm -rf "$tmp_dir"' EXIT

    local archive_path="$tmp_dir/$archive_name"
    local checksum_path="$tmp_dir/$checksum_file"

    log_inf "Installing kredenv $version for $os/$arch"
    log_inf "Downloading from $download_url"
    log_inf "Installing to $INSTALL_DIR"

    # check existing installation
    if command -v "$BIN_NAME" &>/dev/null; then
        local existing_version
        existing_version=$("$BIN_NAME" --version 2>/dev/null || true)
        log_wrn "Found existing $existing_version (will be upgraded to $version)"
    fi

    # download
    log_inf "Downloading artifacts (archive, checksums)"
    curl -fsSL "$download_url" -o "$archive_path"
    curl -fsSL "$checksum_url" -o "$checksum_path"

    # verify
    verify_checksum "$archive_path" "$checksum_path" "$archive_name"

    # extract
    tar -xzf "$archive_path" -C "$tmp_dir"

    # install
    mkdir -p "$INSTALL_DIR"

    local extracted_bin="$tmp_dir/$BIN_NAME"
    if [[ ! -f "$extracted_bin" ]]; then
        log_err "Extracted binary not found at expected location"
        exit 1
    fi

    mv "$extracted_bin" "$INSTALL_DIR/$BIN_NAME"
    chmod +x "$INSTALL_DIR/$BIN_NAME"

    log_inf "Installed kredenv $version to $INSTALL_DIR/$BIN_NAME"

    # PATH check
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_wrn "$INSTALL_DIR is not in PATH"
        log_wrn "Add the following to your shell profile:"
        log_wrn "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi

    echo ""
    printf "\033[32mNext steps:\033[0m\n"
    echo "  1. Ensure $INSTALL_DIR is in your PATH"
    echo "  2. Run: 'kredenv hook --help' for help with configuring the shell of your choice"
}

# Entrypoint
echo ""
printf "\033[32mKREDENV Linux/macOS Installer\033[0m\n"
printf "Source Repo: https://github.com/$REPO\n"
echo ""

install_kredenv