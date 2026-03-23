# Changelog

All notable changes to shellkit will be documented in this file.

## [1.22.0] - 2026-03-22


### Bug Fixes

- Fix tdev right-column horizontal split


### Features

- Enable passthrough for terminal notifications


### Miscellaneous

- Bump version to 1.22.0
## [1.21.0] - 2026-03-22


### Features

- Add ai function for claude + codex tmux workspace


### Miscellaneous

- Bump version to 1.21.0


### Refactoring

- Rename ai function to tdev
## [1.20.0] - 2026-03-21


### Features

- Add pane content capture, nvim session restore, 5min saves
- Add nvim auto-session save and tmux boot LaunchAgent


### Miscellaneous

- Bump version to 1.20.0
## [1.19.2] - 2026-03-21


### Documentation

- Record nvm/conda lazy-load as a known pitfall


### Miscellaneous

- Bump version to 1.19.2


### Revert

- Remove lazy-load for nvm and conda
## [1.19.1] - 2026-03-21


### Bug Fixes

- Resolve nvm alias chains for lazy-load PATH setup


### Miscellaneous

- Bump version to 1.19.1
## [1.19.0] - 2026-03-21


### Features

- Add 3-pane tiled layout and per-repo layout scripts to tw


### Miscellaneous

- Bump version to 1.19.0
## [1.18.1] - 2026-03-21


### Miscellaneous

- Bump version to 1.18.1


### Performance

- Lazy-load nvm and conda for 10x faster startup
## [1.18.0] - 2026-03-20


### Features

- Add tmux snippet manager with Ctrl-b S paste binding
- Add Ctrl-b ? help popup with quick reference


### Miscellaneous

- Bump version to 1.18.0
## [1.17.1] - 2026-03-13


### Bug Fixes

- Dispatch mouse events regardless of tab focus state


### Miscellaneous

- Bump version to 1.17.1
## [1.17.0] - 2026-03-13


### Bug Fixes

- Prevent exec tmux on source ~/.zshrc reload
- Skip tmux auto-start inside cmux
- Also check CMUX_BUNDLE_ID for tmux auto-start guard
- Use $HOME instead of ZDOTDIR for dotfile paths
- Fallback to xterm-256color for unknown TERM types


### Documentation

- Add tmux integration and workspace isolation documentation


### Features

- Auto-start tmux inside cmux
- Enable mouse scroll support for viewport tabs
- Add tw function for per-repo tmux workspaces
- Enable automatic window renaming
- Restructure tmux reference with hierarchy and full keybindings
- Auto-isolate tmux sessions per cmux workspace


### Miscellaneous

- Bump version to 1.17.0


### Styling

- Add heavy pane borders with active highlight
## [1.16.0] - 2026-03-13


### Features

- Make tmux the default terminal experience
## [1.15.0] - 2026-03-11


### Bug Fixes

- Use npm install for markdown-preview build step


### Features

- Calm animations, star battle, config preview, tab styling
- Matrix animation, config tab enhancements, settings sync
## [1.14.0] - 2026-03-08


### Features

- Visual overhaul with gradient header, rich tab bar, and styled components
- Replace Info tab with Dashboard showing startup metrics and system info
- Add startup time recording to zshrc
- Move Dashboard tab to first position
- Animated dot wave in header and fix metrics timestamp parsing
- Add arrow key navigation between tab bar and content
- Add UI settings persistence
- Expanded animation system, config tab, and animated chrome
- Compact tab accent toggle and config tab sizing
- Add markdown preview with Mermaid.js support


### Miscellaneous

- Add bluera-knowledge config and update gitignore


### Refactoring

- Robust version detection with multi-root search
## [1.13.0] - 2026-03-05


### Bug Fixes

- Clarify update progress messages
- Fix tab styling and speed up startup
- Remove double line under tab bar


### Features

- Auto-download TUI binary on shellkit update
## [1.12.0] - 2026-03-05


### Bug Fixes

- Use brew tap for pet (not in core homebrew)
- Include cmd/shellkit-tui entry point and fix gitignore pattern


### Features

- Add pet CLI snippet manager with Ctrl-S widget
- Add Go TUI binary built on Charm stack
## [1.11.0] - 2026-03-03


### Features

- Add terminalizer for terminal-to-GIF recording
- Add AI command prompt widget (Alt-A)
## [1.10.0] - 2026-03-03


### Documentation

- Update tmux TUI page with setup and shellkit bindings


### Features

- Add pbcopy/pbpaste shims for Linux and WSL
- Add managed tmux.conf with TPM, resurrect, and sshx
## [1.9.0] - 2026-03-01


### Bug Fixes

- Create ~/.local/bin before chezmoi install
- Fall back to ~/bin when ~/.local/bin is not writable


### Features

- Add unified search across all shellkit content
- Isolate shell history per session with merge on exit
- Add local GitLab server (gitlab.home.lan) support
- Support local overrides via ~/.ssh/config.local
- Colorize glog alias with hash, refs, date, and author
## [1.8.1] - 2026-02-11


### Bug Fixes

- Auto-close nvim-tree when it's the last window
## [1.8.0] - 2026-02-10


### Bug Fixes

- Anchor plan selection on divider lines


### Features

- Add clipshot package and TUI quick reference page
## [1.7.0] - 2026-02-10


### Bug Fixes

- Select plan content instead of header block


### Features

- Add tmux quick reference page to TUI
## [1.6.0] - 2026-02-09


### Features

- Promote to core dependency with aliases and iTerm2 support
## [1.5.1] - 2026-02-08


### Bug Fixes

- Strip all trailing spaces from package names
- Add binary mappings for git-delta and 1password-cli
- Suppress p10k instant prompt warning from direnv output
## [1.5.0] - 2026-02-07


### Bug Fixes

- Auto-reload functions and aliases after update/apply
- Add navigation loops to fzf TUI menus
- Make left arrow no-op on main TUI menu
- Render colors correctly in all fzf screens


### Features

- Add fzf drill-down for alias categories
## [1.4.1] - 2026-02-07


### Bug Fixes

- Use top-down layout for fzf TUI menus
- Bind arrow keys for TUI navigation
## [1.4.0] - 2026-02-07


### Bug Fixes

- Use lightweight shell for fzf preview commands


### Features

- Add treesitter for syntax highlighting
- Rewrite as full-fledged CLI with discovery and diagnostics
- Add fzf-powered interactive TUI navigation
## [1.3.0] - 2026-02-07


### Features

- Add glog alias for pretty one-line git log
- Add iTerm2 script to copy Claude Code plan blocks to clipboard
- Isolate shell history per terminal session
- Add file explorer sidebar with nvim-tree
- Redesign help output with colors and grouped commands
- Add .zshrc protection with migrate-to-local option


### Miscellaneous

- Update gitignore for plugin working files and local overrides
## [1.2.0] - 2026-02-04


### Features

- Add secret management with direnv and 1Password CLI
- Add .zshrc.local for machine-specific customizations
## [1.1.0] - 2026-02-02


### Documentation

- Replace ASCII diagram with mermaid flowchart


### Features

- Add vim to bootstrap tier for all profiles
- Detect conditions that disable Touch ID on macOS
## [1.0.1] - 2026-01-29


### Bug Fixes

- Guard work gitconfig template when work identity disabled
- Clarify identity prompts are for Git config
- Remove shellkit alias that conflicts with function
- Add shellkit to autoload list


### Documentation

- Document optional work identity


### Features

- Make work identity optional
- Check for zsh and offer to change default shell
- Bootstrap zsh install for all profiles
- Make powerlevel10k common to all profiles
- Add shellkit command for config management
- Auto-create GitHub Release in release script
## [1.0.0] - 2026-01-29


### Features

- Initial release of shellkit v1.0.0

