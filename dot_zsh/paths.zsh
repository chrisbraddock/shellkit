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
# NVM (Node Version Manager)
# ------------------------------------------------------------

if [[ -d "$HOME/.nvm" ]]; then
    export NVM_DIR="$HOME/.nvm"
    [[ -s "$NVM_DIR/nvm.sh" ]] && \. "$NVM_DIR/nvm.sh"
    [[ -s "$NVM_DIR/bash_completion" ]] && \. "$NVM_DIR/bash_completion"
fi

# ------------------------------------------------------------
# Conda (official setup from `conda init zsh`)
# ------------------------------------------------------------

# Detect conda installation (miniforge preferred for Apple Silicon)
if [[ -d "$HOME/miniforge3" ]]; then
    __conda_setup="$("$HOME/miniforge3/bin/conda" 'shell.zsh' 'hook' 2> /dev/null)"
    if [ $? -eq 0 ]; then
        eval "$__conda_setup"
    else
        [[ -f "$HOME/miniforge3/etc/profile.d/conda.sh" ]] && \. "$HOME/miniforge3/etc/profile.d/conda.sh"
    fi
    unset __conda_setup
elif [[ -d "$HOME/miniconda3" ]]; then
    __conda_setup="$("$HOME/miniconda3/bin/conda" 'shell.zsh' 'hook' 2> /dev/null)"
    if [ $? -eq 0 ]; then
        eval "$__conda_setup"
    else
        [[ -f "$HOME/miniconda3/etc/profile.d/conda.sh" ]] && \. "$HOME/miniconda3/etc/profile.d/conda.sh"
    fi
    unset __conda_setup
elif [[ -d "$HOME/anaconda3" ]]; then
    __conda_setup="$("$HOME/anaconda3/bin/conda" 'shell.zsh' 'hook' 2> /dev/null)"
    if [ $? -eq 0 ]; then
        eval "$__conda_setup"
    else
        [[ -f "$HOME/anaconda3/etc/profile.d/conda.sh" ]] && \. "$HOME/anaconda3/etc/profile.d/conda.sh"
    fi
    unset __conda_setup
fi

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
