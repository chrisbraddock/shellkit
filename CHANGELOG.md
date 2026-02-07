# Changelog

All notable changes to shellkit will be documented in this file.

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

