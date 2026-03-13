<div align="center">

# 🐚 shellkit

**A fast, modern shell configuration managed with [chezmoi](https://www.chezmoi.io/)**

[![Version](https://img.shields.io/badge/version-1.17.0-blue.svg)](https://github.com/chrisbraddock/shellkit/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20WSL-lightgrey.svg)]()

[Features](#-features) • [Installation](#-installation) • [Profiles](#-profiles) • [Customization](#-customization) • [Reference](#-aliases-reference)

</div>

---

## 📖 Overview

Shellkit is a complete Zsh dotfiles system built on chezmoi and Antidote. It replaces heavy frameworks like Oh-My-Zsh with a lean, declarative configuration that starts fast and works across macOS, Linux, and WSL.

```mermaid
flowchart LR
    A["⚡ P10k Instant Prompt<br/><i>~50ms</i>"] --> B["🔌 Antidote<br/><i>lazy load</i>"] --> C["🛠️ Tool Inits<br/><i>fzf/zoxide</i>"]
```

**Key benefits:**

| | |
|---|---|
| ⚡ **Fast startup** | Deferred plugin loading with Antidote + Powerlevel10k instant prompt |
| 🌍 **Cross-platform** | Works on macOS, Ubuntu/Debian (APT), and WSL |
| 📝 **Declarative** | All configuration managed via chezmoi templates |
| 📦 **Portable** | Single `chezmoi init` bootstraps everything |

## ✨ Features

| Feature | Description |
|:--------|:------------|
| 🔌 **Antidote** | Plugin manager with deferred loading for fast shell startup |
| 🎨 **Powerlevel10k** | Feature-rich prompt with instant prompt support |
| 🔍 **Atuin** | Shell history with fuzzy search (local-only, no cloud sync) |
| 📁 **Zoxide** | Smart `cd` replacement — jump anywhere with `z` |
| 🔎 **FZF** | Fuzzy finder for files, history, and completions |
| 🔀 **40+ Git aliases** | Shortcuts for status, branches, commits, diffs, stash |
| 🛠️ **Git helper functions** | `git-undo`, `git-amend`, `git-up`, and more |
| 💻 **OS-aware aliases** | Platform-specific commands for macOS/Linux/WSL |
| 👆 **Touch ID for sudo** | Biometric sudo on macOS (auto-enabled) |
| 👥 **Multi-identity Git** | Separate personal and work Git/SSH identities |
| 🔐 **Secret management** | Per-directory env vars with direnv + 1Password CLI |
| 📂 **File Explorer** | nvim-tree sidebar with auto-open, git status, and icons (neovim) |
| 📋 **Copy Claude Plan** | iTerm2 hotkey to copy Claude Code plan blocks to clipboard (macOS, full profile) |
| 🖥️ **tmux auto-start** | Auto-launches tmux on new shells with per-repo workspace isolation in cmux |
| 📝 **Markdown preview** | Live preview with Mermaid.js diagram rendering in Neovim |

## 📥 Installation

### Prerequisites

- 🐚 **Zsh** — installed automatically if missing
- 📦 **Git** — required
- 🌐 **curl** — required for Homebrew bootstrap
- 🔤 **[Nerd Font](https://www.nerdfonts.com/)** — configured in your terminal for Powerlevel10k icons

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/chrisbraddock/shellkit/main/install.sh | bash
```

### Manual Install

```bash
git clone https://github.com/chrisbraddock/shellkit ~/repos/shellkit
cd ~/repos/shellkit && ./init.sh
```

<details>
<summary><strong>🔧 What happens on first run?</strong></summary>

1. Prompt for Git identity (name, email, SSH key)
2. Prompt for optional work identity (can skip for personal-only setups)
3. Prompt for repo directories and optional tool toggles
3. Install Homebrew (if not present)
4. Install all packages from `.chezmoidata/packages.yaml`
5. Apply all dotfiles and enable biometric sudo (if selected)

Restart your terminal to load the new configuration.

</details>

## 🎭 Profiles

Shellkit supports two profiles for different use cases:

| Profile | Description | Use Case |
|:--------|:------------|:---------|
| 🖥️ **full** | All packages, plugins, and tools | Workstations, development machines |
| 📟 **minimal** | Basic shell config only | Servers, SSH hosts, containers |

Select your profile during `chezmoi init`.

**Both profiles include:** zsh, Powerlevel10k prompt, git config, shell aliases

**Full profile adds:**
- Package installation (Homebrew, APT)
- Antidote plugin manager + zsh plugins
- Tool integrations (atuin, zoxide, fzf)
- Biometric sudo setup

<details>
<summary><strong>📂 Directory Structure</strong></summary>

```
shellkit/
├── dot_zshrc.tmpl              # Main Zsh config (chezmoi template)
├── dot_tmux.conf.tmpl          # tmux config (splits, vim-nav, persistence)
├── dot_zsh_plugins.txt         # Antidote plugin list
├── dot_p10k.zsh                # Powerlevel10k theme config
├── dot_gitconfig.tmpl          # Git configuration
├── dot_gitconfig-personal.tmpl # Personal identity (name/email)
├── dot_gitconfig-work.tmpl     # Work identity (name/email)
├── dot_vimrc                   # Vim configuration
├── .chezmoi.toml.tmpl          # Chezmoi variables (name, email, etc.)
│
├── dot_zsh/
│   ├── aliases/
│   │   ├── git.zsh             # Git shortcuts (gs, gc, gd, etc.)
│   │   ├── system.zsh.tmpl     # File/network helpers (ll, etc.)
│   │   ├── tmux.zsh            # tmux session/window/pane aliases
│   │   └── tools.zsh           # Editor and tool aliases
│   ├── functions/              # Shell utilities and git helpers
│   │   ├── shellkit            # Shellkit management (update, edit, diff)
│   │   ├── mon                 # System monitoring launcher
│   │   ├── reload              # Reload shell configuration
│   │   ├── sshx                # SSH with automatic tmux attach
│   │   ├── tw                  # Per-repo tmux workspace (isolated socket)
│   │   ├── sudo-biometrics     # Biometric sudo setup
│   │   └── git-*               # Git helper functions
│   ├── completion.zsh          # Completion styles
│   ├── os-detect.zsh           # OS_TYPE detection
│   └── paths.zsh               # PATH setup for tools
│
├── dot_config/
│   ├── atuin/config.toml       # Local-only history config
│   └── nvim/
│       ├── init.vim            # Sources vimrc + loads Lua config
│       └── lua/                # Neovim Lua config (plugins, nvim-tree)
│
├── private_dot_ssh/
│   └── config.tmpl             # SSH config with host aliases
│
├── iterm2/
│   └── copy_claude_plan.py     # Copy Claude plan block to clipboard (iTerm2 Python API)
│
├── .chezmoidata/
│   └── packages.yaml           # Declarative package list
│
└── run_*.sh.tmpl               # Bootstrap and package scripts
```

</details>

## ⚙️ Customization

### Modify Aliases

Edit files in `dot_zsh/aliases/`:

| File | Contents |
|:-----|:---------|
| `git.zsh` | Git shortcuts |
| `system.zsh.tmpl` | OS commands, navigation |
| `tools.zsh` | Development tool aliases |

After editing, run `chezmoi apply` to update your live config.

### Machine-Specific Config

For customizations that shouldn't sync across machines, create `~/.zshrc.local`:

```bash
# ~/.zshrc.local - not managed by chezmoi
export MY_LOCAL_VAR="value"
alias myalias="command"
source ~/work-specific-tools.sh
```

This file is sourced at the end of `.zshrc` and won't cause conflicts during `shellkit update`.

### Add/Remove Plugins

Edit `dot_zsh_plugins.txt`:

```txt
# Standard plugin
username/plugin-name

# Deferred loading (faster startup)
username/plugin-name kind:defer

# Completions (loads before compinit)
username/plugin-name kind:fpath
```

Then: `chezmoi apply && source ~/.zshrc`

<details>
<summary><strong>📋 Template Variables</strong></summary>

Configure personal settings in `.chezmoi.toml.tmpl`. On first run, chezmoi prompts for these values:

| Variable | Purpose |
|:---------|:--------|
| `.identities.personal.name` | Personal full name (for git commits) |
| `.identities.personal.email` | Personal email (for git commits) |
| `.identities.personal.sshKey` | Personal SSH key filename |
| `.identities.work.enabled` | Enable separate work identity (optional) |
| `.identities.work.label` | Work identity label (e.g., `work`, `company`) |
| `.identities.work.name` | Work full name |
| `.identities.work.email` | Work email |
| `.identities.work.sshKey` | Work SSH key filename |
| `.gitDirs.personal` | Personal repos directory (e.g., `~/repos/personal/`) |
| `.gitDirs.work` | Work repos directory (e.g., `~/repos/work/`) |
| `.install_atuin` | Enable atuin shell init |
| `.install_zoxide` | Enable zoxide shell init |
| `.install_fzf` | Enable fzf shell init |
| `.enable_sudo_biometrics` | Enable biometric sudo (Touch ID on macOS) |

Re-run `chezmoi init --prompt` to update these values.

</details>

### 🔑 Multi-Identity SSH Setup

Shellkit configures separate SSH keys for personal and work GitHub accounts:

| Host Alias | SSH Key Used |
|:-----------|:-------------|
| `github.com` | Personal SSH key |
| `github.com-personal` | Personal SSH key (explicit) |
| `github.com-work` | Work SSH key (uses your configured label) |

```bash
# Clone with work identity
git clone git@github.com-work:org/repo.git

# Or update existing remote
git remote set-url origin git@github.com-work:org/repo.git
```

Git identity (name/email) is selected automatically based on the repository directory.

<details>
<summary><strong>🔒 Security Considerations</strong></summary>

This configuration makes intentional security tradeoffs for convenience:

| Setting | Platform | Impact |
|:--------|:---------|:-------|
| `credential.helper = store` | Linux/WSL | Stores Git credentials in plaintext (`~/.git-credentials`). Consider using a credential manager for sensitive repos. |
| `protocol.file.allow = always` | All | Allows `file://` URLs in Git. Safe for local use but be cautious with untrusted repositories. |

These are acceptable for personal development machines but should be reviewed for shared or production environments.

</details>

### 🔐 Secret Management

Shellkit includes [direnv](https://direnv.net/) for per-directory environment variables and optional [1Password CLI](https://developer.1password.com/docs/cli/) integration.

```bash
# Create .envrc in your project
secrets init

# Edit and add secrets
secrets edit

# Trust and load the file
secrets allow
```

<details>
<summary><strong>Using 1Password CLI</strong></summary>

Instead of hardcoding secrets, reference them from 1Password:

```bash
# In your .envrc:
export OPENAI_API_KEY=$(op read "op://Development/OpenAI/api-key")
export GITHUB_TOKEN=$(op read "op://Development/GitHub CLI/token")
```

Benefits:
- Secrets are fetched at load time, never written to disk
- Uses biometric/password authentication
- Easy to share patterns (`.envrc.example`) without exposing values

Setup:
```bash
# Check 1Password status
secrets op-status

# Sign in (requires 1Password app integration enabled)
secrets op-login
```

</details>

---

## 📚 Aliases Reference

### 🔀 Git Aliases

| Alias | Command | Description |
|:------|:--------|:------------|
| `gs` | `git status` | Show working tree status |
| `ga` | `git add` | Stage files |
| `gaa` | `git add --all` | Stage all changes |
| `gc` | `git commit` | Commit changes |
| `gcm` | `git commit -m` | Commit with message |
| `gca` | `git commit --amend` | Amend last commit |
| `gd` | `git diff` | Show unstaged changes |
| `gds` | `git diff --staged` | Show staged changes |
| `gco` | `git checkout` | Switch branches |
| `gcob` | `git checkout -b` | Create and switch branch |
| `gb` | `git branch` | List branches |
| `gl` | `git log` | Show commit log |
| `glog` | `git log --oneline --graph --decorate` | Pretty one-line log |
| `glg` | `git log --graph` | Commit graph |
| `gf` | `git fetch` | Fetch remote |
| `gpl` | `git pull` | Pull changes |
| `gst` | `git stash` | Stash changes |
| `gstp` | `git stash pop` | Apply stash |

> See `dot_zsh/aliases/git.zsh` for the complete list.

### 🛠️ Git Helper Functions

| Function | Description |
|:---------|:------------|
| `git-all` | Stage all changes (`git add -A`) |
| `git-amend` | Amend last commit, keep message |
| `git-undo` | Undo last commit, keep changes staged |
| `git-up` | Pull and show what changed |
| `git-promote` | Push branch and set up tracking |
| `git-credit` | Rewrite commit author: `git-credit "Name" email` |
| `git-delete-local-merged` | Delete local branches already merged |

### 🖥️ Shell Utilities

| Function | Description |
|:---------|:------------|
| `shellkit update` | Pull latest config and apply changes |
| `shellkit edit` | Edit config files with chezmoi |
| `shellkit diff` | Show pending changes |
| `shellkit cd` | Jump to shellkit source directory |
| `reload` | Reload shell configuration |
| `mon` | System monitoring (`mon cpu`, `mon gpu`, `mon net`\*, `mon disk`\*) |
| `sudo-biometrics` | Biometric sudo auth (`status`, `enable`, `disable`, `test`) |
| `secrets` | Secret management (`status`, `init`, `edit`, `allow`, `op-login`) |
| `sshx` | SSH with automatic tmux attach on remote |
| `tw` | Per-repo tmux workspace with isolated socket |
| `wimip` | Show machine's IP address |

\* Requires sudo

### 🖥️ tmux

tmux auto-starts on new shells (full profile) with these key features:

| Feature | Description |
|:--------|:------------|
| **Auto-start** | New shells launch into tmux automatically |
| **cmux isolation** | Each cmux tab gets an isolated tmux socket — no session collisions |
| **SSH auto-attach** | SSH to `.home.lan` hosts auto-attaches tmux |
| **Per-repo workspaces** | `tw <path>` creates an isolated tmux workspace per repository |
| **Nested tmux** | `Ctrl-b Ctrl-b` sends prefix to inner tmux (SSH sessions) |
| **Session persistence** | tmux-resurrect + tmux-continuum auto-saves every 10 minutes |
| **Mouse support** | Click panes, drag borders to resize, scroll with mouse wheel |

<details>
<summary><strong>tmux Aliases</strong></summary>

| Alias | Description |
|:------|:------------|
| `ta` | Attach to most recent session |
| `tas <name>` | Attach to named session |
| `tad` | Attach and detach other clients |
| `tns <name>` | New named session |
| `tls` | List sessions |
| `tks <name>` | Kill named session |
| `tkill` | Kill entire tmux server |
| `tlw` | List windows |
| `tlp` | List panes |

</details>

### ⌨️ System Aliases

| Alias | Description |
|:------|:------------|
| `ll` | List files with details (colored) |
| `e` | Open in editor (VS Code) |
| `o` | Open current directory in file manager |
| `r` | `pushd ~/repos` |
| `c` | Clear terminal |
| `hosts` | Edit /etc/hosts |
| `flushdns` | Clear DNS cache |

> **Note:** Clone repos into `~/repos/personal/` or `~/repos/work/` to automatically use the correct Git identity.

---

## 📋 iTerm2: Copy Claude Plan

**macOS, full profile only.** An iTerm2 Python API script that copies the most recent Claude Code plan block from terminal scrollback to the clipboard via a hotkey.

<details>
<summary><strong>Setup</strong></summary>

The script is auto-deployed to `~/Library/Application Support/iTerm2/Scripts/AutoLaunch/` by `chezmoi apply`. To activate:

1. **Enable Python API**: iTerm2 > Settings > General > Magic > Enable Python API
2. **Install Python Runtime** if prompted: Scripts > Manage > Install Python Runtime
3. **Restart iTerm2** (AutoLaunch picks up the script automatically)
4. **Bind a hotkey**: Settings > Keys > Key Bindings > **+**
   - Action: **Invoke Script Function**
   - Function: `copy_claude_plan()`
   - Assign your preferred key combination

</details>

<details>
<summary><strong>How it works</strong></summary>

When invoked, the script:
1. Scans the last 20,000 lines of scrollback (bottom-up)
2. Finds the most recent `Ready to code?` marker
3. Includes any decorative border lines above it
4. Captures everything up to (but not including) the end divider
5. Copies the block to the system clipboard via `pbcopy`
6. Highlights the copied region in the terminal

Configuration constants (in `iterm2/copy_claude_plan.py`): `MAX_LINES`, `SELECT_IN_TERMINAL`, `DEBUG`.

</details>

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Use [conventional commits](https://www.conventionalcommits.org/) for commit messages
4. Submit a pull request

## 📄 License

MIT
