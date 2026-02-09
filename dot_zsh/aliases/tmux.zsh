# ============================================================
# Tmux Aliases
# ============================================================
# Session, window, and pane management
# ============================================================

if command -v tmux &>/dev/null; then

    # ------------------------------------------------------------
    # Session Management
    # ------------------------------------------------------------

    alias ta='tmux attach'                  # Attach to most recent session
    alias tad='tmux attach -d'              # Attach and detach other clients
    alias tas='tmux attach -t'              # Attach to named session
    alias tns='tmux new-session -s'         # New named session
    alias tls='tmux list-sessions'          # List sessions
    alias tks='tmux kill-session -t'        # Kill named session
    alias tkill='tmux kill-server'          # Kill entire tmux server

    # ------------------------------------------------------------
    # Window & Pane
    # ------------------------------------------------------------

    alias tlw='tmux list-windows'           # List windows
    alias tlp='tmux list-panes'             # List panes

    # ------------------------------------------------------------
    # iTerm2 Control Mode (tmux -CC)
    # ------------------------------------------------------------
    # Maps tmux windows/panes to native iTerm2 tabs/splits

    if [[ "$TERM_PROGRAM" == "iTerm.app" ]]; then
        alias tcc='tmux -CC'                    # New session in control mode
        alias tccs='tmux -CC new-session -s'    # New named session in control mode
        alias tcca='tmux -CC attach'            # Attach in control mode
        alias tccas='tmux -CC attach -t'        # Attach to named session in control mode
    fi

fi
