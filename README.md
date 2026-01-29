<div align="center">

# 🐚 shellkit

**A fast, modern shell configuration managed with [chezmoi](https://www.chezmoi.io/)**

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/chrisbraddock/shellkit/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20WSL-lightgrey.svg)]()

[Features](#-features) • [Installation](#-installation) • [Profiles](#-profiles) • [Customization](#-customization) • [Reference](#-aliases-reference)

</div>

---

## 📖 Overview

Shellkit is a complete Zsh dotfiles system built on chezmoi and Antidote. It replaces heavy frameworks like Oh-My-Zsh with a lean, declarative configuration that starts fast and works across macOS, Linux, and WSL.

```
┌─────────────────────────────────────────────────────────────────┐
│                         Shell Startup                           │
├─────────────────────────────────────────────────────────────────┤
│  P10k Instant Prompt  →  Antidote (deferred)  →  Tool Inits    │
│        ⚡ <50ms              🔌 lazy load          🛠️ fzf/zoxide │
└─────────────────────────────────────────────────────────────────┘
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

1. Prompt for identity settings (name, email, SSH key for personal and work)
2. Prompt for repo directories and optional tool toggles
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
│   │   └── tools.zsh           # Editor and tool aliases
│   ├── functions/              # Shell utilities and git helpers
│   │   ├── shellkit            # Shellkit management (update, edit, diff)
│   │   ├── mon                 # System monitoring launcher
│   │   ├── reload              # Reload shell configuration
│   │   ├── sudo-biometrics     # Biometric sudo setup
│   │   └── git-*               # Git helper functions
│   ├── completion.zsh          # Completion styles
│   ├── os-detect.zsh           # OS_TYPE detection
│   └── paths.zsh               # PATH setup for tools
│
├── dot_config/
│   ├── atuin/config.toml       # Local-only history config
│   └── nvim/init.vim           # Neovim configuration
│
├── private_dot_ssh/
│   └── config.tmpl             # SSH config with host aliases
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
| `wimip` | Show machine's IP address |

\* Requires sudo

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

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Use [conventional commits](https://www.conventionalcommits.org/) for commit messages
4. Submit a pull request

## 📄 License

MIT
