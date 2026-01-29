# ============================================================
# OS Detection
# ============================================================
# Sets OS_TYPE environment variable for OS-specific configuration
#
# Usage in other scripts:
#   if [[ "$OS_TYPE" == "macos" ]]; then
#       # macOS-specific code
#   elif [[ "$OS_TYPE" == "linux" || "$OS_TYPE" == "ubuntu" || "$OS_TYPE" == "wsl" ]]; then
#       # Linux-specific code
#   fi
# ============================================================

if [[ "$(uname)" == "Darwin" ]]; then
    export OS_TYPE="macos"
elif [[ "$(uname -s)" == Linux* ]]; then
    if [[ -n "$WSL_DISTRO_NAME" ]]; then
        export OS_TYPE="wsl"
    elif [[ -e /etc/os-release ]] && grep -q '^ID=ubuntu' /etc/os-release; then
        export OS_TYPE="ubuntu"
    else
        export OS_TYPE="linux"
    fi
else
    export OS_TYPE="unknown"
fi
