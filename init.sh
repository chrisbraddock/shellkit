#!/usr/bin/env bash
set -euo pipefail

SHELLKIT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHEZMOI_SOURCE="$HOME/.local/share/chezmoi"

echo "==> Initializing shellkit..."

# Install chezmoi if not present
if ! command -v chezmoi &>/dev/null; then
    echo "==> Installing chezmoi..."
    sh -c "$(curl -fsLS get.chezmoi.io)" -- -b "$HOME/.local/bin"
    export PATH="$HOME/.local/bin:$PATH"
fi

# Handle existing chezmoi source
if [[ -e "$CHEZMOI_SOURCE" ]]; then
    if [[ -L "$CHEZMOI_SOURCE" ]]; then
        current_target="$(readlink "$CHEZMOI_SOURCE")"
        if [[ "$current_target" == "$SHELLKIT_DIR" ]]; then
            echo "Symlink already correct."
        else
            echo "ERROR: $CHEZMOI_SOURCE points to $current_target"
            echo "Expected: $SHELLKIT_DIR"
            exit 1
        fi
    else
        echo "ERROR: $CHEZMOI_SOURCE exists as directory."
        echo "Remove it first: rm -rf $CHEZMOI_SOURCE"
        exit 1
    fi
else
    mkdir -p "$(dirname "$CHEZMOI_SOURCE")"
    ln -s "$SHELLKIT_DIR" "$CHEZMOI_SOURCE"
    echo "Created symlink: $CHEZMOI_SOURCE -> $SHELLKIT_DIR"
fi

# Initialize chezmoi config if needed
if [[ ! -f "$HOME/.config/chezmoi/chezmoi.toml" ]]; then
    echo "==> Running chezmoi init..."
    chezmoi init
else
    echo "==> Existing chezmoi config found."
    echo "    To reconfigure shellkit options, run: chezmoi init --prompt"
fi

# Apply
echo "==> Applying configuration..."
chezmoi apply

echo "==> Done! Start a new shell to use your updated config."
