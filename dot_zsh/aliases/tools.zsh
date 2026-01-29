# ============================================================
# Tool Aliases
# ============================================================
# Aliases for development tools
# ============================================================

# ------------------------------------------------------------
# Editors
# ------------------------------------------------------------

# Use nvim for vi/vim if available
if command -v nvim &>/dev/null; then
    alias vi='nvim'
    alias vim='nvim'
fi

# ------------------------------------------------------------
# Rust
# ------------------------------------------------------------

# Update Rust toolchain
if command -v rustup &>/dev/null; then
    alias update-rust='rustup update'
fi
