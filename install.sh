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

# Check if zsh is the default login shell
if [[ "$SHELL" != */zsh ]]; then
    echo ""
    echo "NOTE: shellkit requires zsh as your default shell."
    echo ""

    if command -v zsh &>/dev/null; then
        ZSH_PATH=$(command -v zsh)
        echo "zsh is installed at: $ZSH_PATH"
        read -p "Change default shell to zsh now? [y/N] " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            if chsh -s "$ZSH_PATH"; then
                echo "Default shell changed. Log out and back in to use zsh."
            else
                echo "Failed to change shell. Run manually: chsh -s $ZSH_PATH"
            fi
        fi
    else
        echo "zsh will be installed by shellkit. After install completes, run:"
        echo "  chsh -s \$(which zsh)"
    fi
    echo ""
fi

# Install chezmoi if needed
if ! command -v chezmoi &>/dev/null; then
    echo "==> Installing chezmoi..."
    CHEZMOI_BIN="$HOME/.local/bin"
    if ! mkdir -p "$CHEZMOI_BIN" 2>/dev/null; then
        CHEZMOI_BIN="$HOME/bin"
        mkdir -p "$CHEZMOI_BIN"
    fi
    sh -c "$(curl -fsLS get.chezmoi.io)" -- -b "$CHEZMOI_BIN"
    export PATH="$CHEZMOI_BIN:$PATH"
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
