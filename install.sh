#!/usr/bin/env bash
# shellkit installer
# Usage: curl -fsSL https://raw.githubusercontent.com/chrisbraddock/shellkit/main/install.sh | bash

set -euo pipefail

REPO="chrisbraddock/shellkit"

echo "==> Installing shellkit..."
echo ""

# Check prerequisites
for cmd in curl git; do
    if ! command -v "$cmd" &>/dev/null; then
        echo "ERROR: $cmd is required but not installed." >&2
        exit 1
    fi
done

# Install chezmoi if needed
if ! command -v chezmoi &>/dev/null; then
    echo "==> Installing chezmoi..."
    sh -c "$(curl -fsLS get.chezmoi.io)" -- -b "$HOME/.local/bin"
    export PATH="$HOME/.local/bin:$PATH"
fi

# Check for existing config
if [[ -f "$HOME/.config/chezmoi/chezmoi.toml" ]]; then
    echo "==> Existing chezmoi config found."
    echo "    Running with --prompt to allow reconfiguration..."
    echo ""
    chezmoi init --prompt --apply "$REPO"
else
    # Profile hint for new installs
    echo "You'll be prompted for configuration options."
    echo "Profile options: 'full' (all features) or 'minimal' (basic shell only)"
    echo ""
    chezmoi init --apply "$REPO"
fi

echo ""
echo "==> shellkit installed successfully!"
echo "    Start a new shell to use your updated config."
