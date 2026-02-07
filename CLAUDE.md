# CLAUDE.md - Shellkit

## Project Overview

Shellkit is a cross-platform Zsh configuration system managed with [chezmoi](https://www.chezmoi.io/). It provides a fast, declarative dotfiles setup that works on macOS, Linux, and WSL.

## Key Directories

```
shellkit/
├── dot_zsh/              # Zsh configuration
│   ├── aliases/          # Git, system, and tool aliases
│   ├── functions/        # Shell utilities (shellkit, mon, secrets, etc.)
│   ├── os-detect.zsh     # Sets $OS_TYPE variable
│   └── paths.zsh         # PATH setup
├── dot_config/           # Tool configs (atuin, nvim)
│   └── nvim/lua/         # Neovim Lua config (lazy.nvim, nvim-tree)
├── private_dot_ssh/      # SSH configuration templates
├── iterm2/               # iTerm2 Python API scripts (deployed via run script)
├── .chezmoidata/         # Package definitions (packages.yaml)
└── run_*.sh.tmpl         # Bootstrap and install scripts
```

## Conventions

- **Templates**: Files ending in `.tmpl` are chezmoi templates
- **OS detection**: Use `$OS_TYPE` variable (`macos`, `linux`, `wsl`, `ubuntu`)
- **Profile check**: Use `{{ if eq .profile "full" }}` for full-profile-only features
- **Functions**: Add new functions to `dot_zsh/functions/` and autoload in `dot_zshrc.tmpl`
- **Aliases**: Add to appropriate file in `dot_zsh/aliases/`
- **Packages**: Define in `.chezmoidata/packages.yaml`
- **iTerm2 scripts**: Add to `iterm2/`, deployed to AutoLaunch via `run_onchange_after_deploy-iterm2-scripts.sh.tmpl`
- **Neovim plugins**: Add plugin specs to `dot_config/nvim/lua/plugins.lua`, config to `dot_config/nvim/lua/plugin/`
- **Local overrides**: Users can create `~/.zshrc.local` for machine-specific config (not managed by chezmoi)

## Development Commands

```bash
chezmoi apply              # Apply dotfile changes
chezmoi diff               # Preview pending changes
chezmoi edit <file>        # Edit a managed file
shellkit update            # Pull latest and apply
shellkit diff              # Alias for chezmoi diff
```

## Testing Changes

1. Edit files in the shellkit directory
2. Run `chezmoi diff` to preview changes
3. Run `chezmoi apply` to apply
4. Open a new terminal or run `reload` to test

## Package Management

Packages are defined in `.chezmoidata/packages.yaml`:
- `shell`: Shell runtime essentials (zsh, fzf, zoxide, direnv)
- `editor`: Editor tools (neovim, ripgrep, fd)
- `tools`: CLI utilities (jq, bat, tmux, etc.)
- `monitoring`: System monitors (btop, htop)
- `darwin`/`linux`: OS-specific packages

## Adding a New Feature

1. **Function**: Create `dot_zsh/functions/<name>`, add autoload to `dot_zshrc.tmpl`
2. **Alias**: Add to appropriate file in `dot_zsh/aliases/`
3. **Package**: Add to `.chezmoidata/packages.yaml`
4. **Documentation**: Update README.md and this file

## Profiles

- **full**: All packages, plugins, biometric sudo (workstations)
- **minimal**: Basic shell config only (servers, containers)

## Multi-Identity Git

Git identity is selected based on repo directory:
- `~/repos/personal/` → personal identity
- `~/repos/work/` → work identity

SSH host aliases: `github.com-personal`, `github.com-work`
