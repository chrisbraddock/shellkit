# ============================================================
# Path and Environment Configuration
# ============================================================

# ------------------------------------------------------------
# Helper Functions
# ------------------------------------------------------------

# Add directory to PATH if it exists and isn't already present
add_to_path() {
    if [[ -d "$1" ]] && [[ ":$PATH:" != *":$1:"* ]]; then
        PATH="$1:$PATH"
    fi
}

# Add directory to LD_LIBRARY_PATH if it exists
add_to_ld_library_path() {
    if [[ -d "$1" ]] && [[ ":$LD_LIBRARY_PATH:" != *":$1:"* ]]; then
        if [[ -z "$LD_LIBRARY_PATH" ]]; then
            LD_LIBRARY_PATH="$1"
        else
            LD_LIBRARY_PATH="$1:$LD_LIBRARY_PATH"
        fi
    fi
}

# ------------------------------------------------------------
# Homebrew
# ------------------------------------------------------------

if [[ -f "/opt/homebrew/bin/brew" ]]; then
    # Apple Silicon macOS
    eval "$(/opt/homebrew/bin/brew shellenv)"
elif [[ -f "/usr/local/bin/brew" ]]; then
    # Intel macOS
    eval "$(/usr/local/bin/brew shellenv)"
elif [[ -f "/home/linuxbrew/.linuxbrew/bin/brew" ]]; then
    # Linux (Linuxbrew)
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
fi

# ------------------------------------------------------------
# Base PATH
# ------------------------------------------------------------

add_to_path "$HOME/bin"
add_to_path "$HOME/.local/bin"
add_to_path "/usr/local/share/npm/bin"

# macOS-specific paths
if [[ "$OS_TYPE" == "macos" ]]; then
    add_to_path "/Applications/Postgres.app/Contents/MacOS/bin"
fi

# ------------------------------------------------------------
# Rust / Cargo
# ------------------------------------------------------------

if [[ -d "$HOME/.cargo/bin" ]]; then
    add_to_path "$HOME/.cargo/bin"
fi

# ------------------------------------------------------------
# Bun
# ------------------------------------------------------------

if [[ -d "$HOME/.bun" ]]; then
    export BUN_INSTALL="$HOME/.bun"
    add_to_path "$BUN_INSTALL/bin"
fi

# ------------------------------------------------------------
# NVM (Node Version Manager) — lazy loaded for fast startup
# ------------------------------------------------------------

if [[ -d "$HOME/.nvm" ]]; then
    export NVM_DIR="$HOME/.nvm"
    # Lazy-load: nvm init runs on first call to nvm, node, npm, npx, or corepack
    _nvm_lazy_load() {
        unfunction nvm node npm npx corepack 2>/dev/null
        [[ -s "$NVM_DIR/nvm.sh" ]] && \. "$NVM_DIR/nvm.sh"
        [[ -s "$NVM_DIR/bash_completion" ]] && \. "$NVM_DIR/bash_completion"
    }
    nvm()      { _nvm_lazy_load; nvm "$@"; }
    node()     { _nvm_lazy_load; node "$@"; }
    npm()      { _nvm_lazy_load; npm "$@"; }
    npx()      { _nvm_lazy_load; npx "$@"; }
    corepack() { _nvm_lazy_load; corepack "$@"; }
    # Add default node to PATH immediately (no nvm overhead)
    # Resolves alias chains like: default → lts/* → lts/jod → v22.21.0
    [[ -d "$NVM_DIR/versions/node" ]] && {
        local _ver="" _alias_file="$NVM_DIR/alias/default"
        while [[ -f "$_alias_file" ]]; do
            _ver=$(cat "$_alias_file")
            if [[ "$_ver" == v* ]]; then
                break
            elif [[ -f "$NVM_DIR/alias/$_ver" ]]; then
                _alias_file="$NVM_DIR/alias/$_ver"
            else
                _ver=""
                break
            fi
        done
        if [[ -n "$_ver" ]]; then
            local _node_path="$NVM_DIR/versions/node/${_ver}/bin"
            [[ -d "$_node_path" ]] && add_to_path "$_node_path"
        fi
    }
fi

# ------------------------------------------------------------
# Conda — lazy loaded for fast startup
# ------------------------------------------------------------

# Detect conda installation (miniforge preferred for Apple Silicon)
_conda_root=""
if [[ -d "$HOME/miniforge3" ]]; then
    _conda_root="$HOME/miniforge3"
elif [[ -d "$HOME/miniconda3" ]]; then
    _conda_root="$HOME/miniconda3"
elif [[ -d "$HOME/anaconda3" ]]; then
    _conda_root="$HOME/anaconda3"
fi

if [[ -n "$_conda_root" ]]; then
    # Add conda to PATH immediately (no shell hook overhead)
    add_to_path "${_conda_root}/bin"
    # Lazy-load: full conda init runs on first call to conda or mamba
    _conda_lazy_load() {
        unfunction conda mamba 2>/dev/null
        local _root="${_conda_root}"
        __conda_setup="$("${_root}/bin/conda" 'shell.zsh' 'hook' 2> /dev/null)"
        if [ $? -eq 0 ]; then
            eval "$__conda_setup"
        else
            [[ -f "${_root}/etc/profile.d/conda.sh" ]] && \. "${_root}/etc/profile.d/conda.sh"
        fi
        unset __conda_setup
    }
    conda() { _conda_lazy_load; conda "$@"; }
    mamba() { _conda_lazy_load; mamba "$@"; }
fi
unset _conda_root

# ------------------------------------------------------------
# CUDA (Linux)
# ------------------------------------------------------------

if [[ -d "/usr/local/cuda" ]]; then
    add_to_path "/usr/local/cuda/bin"
    add_to_ld_library_path "/usr/local/cuda/lib64"
fi

# WSL NVIDIA GPU libraries
if [[ "$OS_TYPE" == "wsl" ]]; then
    add_to_ld_library_path "/usr/lib/wsl/lib"
fi

# ------------------------------------------------------------
# Environment Variables
# ------------------------------------------------------------

export EDITOR="${EDITOR:-vi}"
export LSCOLORS="exfxcxdxbxegedabagacad"
# GNU LS_COLORS equivalent for Linux and zsh completion colors
export LS_COLORS="di=34:ln=35:so=32:pi=33:ex=31:bd=34;46:cd=34;43:su=30;41:sg=30;46:tw=30;42:ow=30;43"
export CLICOLOR=true

# iTerm2 detection
if [[ "$TERM_PROGRAM" == "iTerm.app" ]]; then
    export ITERM2_CUSTOM="true"
fi

# ------------------------------------------------------------
# Export Modified Paths
# ------------------------------------------------------------

export PATH
export LD_LIBRARY_PATH
