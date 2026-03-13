package data

// Keybinding represents a keyboard shortcut reference.
type Keybinding struct {
	Key         string
	Description string
	Category    string // "tmux", "shell", "shellkit"
}

// LoadKeybindings returns all known keybinding references.
func LoadKeybindings() []Keybinding {
	return []Keybinding{
		// Tmux prefix bindings (Ctrl-b + key)
		{Key: "Ctrl-b |", Description: "Split pane horizontally", Category: "tmux"},
		{Key: "Ctrl-b -", Description: "Split pane vertically", Category: "tmux"},
		{Key: "Ctrl-b h/j/k/l", Description: "Navigate panes (vim-style)", Category: "tmux"},
		{Key: "Ctrl-b H/J/K/L", Description: "Resize panes", Category: "tmux"},
		{Key: "Ctrl-b c", Description: "New window", Category: "tmux"},
		{Key: "Ctrl-b n/p", Description: "Next/previous window", Category: "tmux"},
		{Key: "Ctrl-b 1-9", Description: "Jump to window by number", Category: "tmux"},
		{Key: "Ctrl-b ,", Description: "Rename current window", Category: "tmux"},
		{Key: "Ctrl-b w", Description: "Window/session picker", Category: "tmux"},
		{Key: "Ctrl-b d", Description: "Detach from session", Category: "tmux"},
		{Key: "Ctrl-b [", Description: "Enter copy mode (vi keys)", Category: "tmux"},
		{Key: "Ctrl-b I", Description: "Install TPM plugins", Category: "tmux"},
		{Key: "Ctrl-b Ctrl-s", Description: "Save session (resurrect)", Category: "tmux"},
		{Key: "Ctrl-b Ctrl-r", Description: "Restore session (resurrect)", Category: "tmux"},

		// Shell keybindings
		{Key: "Alt-A", Description: "AI command prompt (claude CLI)", Category: "shell"},
		{Key: "Ctrl-S", Description: "Pet snippet search/insert", Category: "shell"},
		{Key: "Ctrl-R", Description: "Atuin history search", Category: "shell"},
		{Key: "Ctrl-T", Description: "FZF file search", Category: "shell"},
		{Key: "Alt-C", Description: "FZF cd to directory", Category: "shell"},
	}
}

// TmuxReference returns markdown-formatted tmux reference.
func TmuxReference() string {
	return `# Tmux Quick Reference

## Concept: Session > Window > Pane

tmux is hierarchical — think of it like a workspace manager:

- **Session** = workspace (e.g. "main", "project-x")
- **Window** = tab within a session (shown in status bar)
- **Pane** = split within a window

## Sessions
| Command | Description |
|---------|-------------|
| ` + "`ta`" + ` | Attach to most recent session |
| ` + "`tas <name>`" + ` | Attach to named session |
| ` + "`tad`" + ` | Attach and detach other clients |
| ` + "`tns <name>`" + ` | New named session |
| ` + "`tls`" + ` | List sessions |
| ` + "`tks <name>`" + ` | Kill named session |
| ` + "`tkill`" + ` | Kill entire tmux server |
| ` + "`sshx <host>`" + ` | SSH + auto-attach tmux |
| ` + "`tw <path>`" + ` | Open repo workspace (isolated socket) |
| ` + "`Ctrl-b d`" + ` | Detach from current session |

## Windows (Tabs)
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b c`" + ` | New window |
| ` + "`Ctrl-b n / p`" + ` | Next / previous window |
| ` + "`Ctrl-b 1-9`" + ` | Jump to window by number |
| ` + "`Ctrl-b ,`" + ` | Rename current window |
| ` + "`Ctrl-b w`" + ` | Window/session picker (interactive) |
| ` + "`Ctrl-b &`" + ` | Close current window |

## Panes (Splits)
| Key | Description |
|-----|-------------|
| ` + "` Ctrl-b \\| `" + ` | Split horizontally |
| ` + "`Ctrl-b -`" + ` | Split vertically |
| ` + "`Ctrl-b h/j/k/l`" + ` | Navigate panes (vim-style) |
| ` + "`Alt-h/j/k/l`" + ` | Resize panes |
| ` + "`Ctrl-b z`" + ` | Toggle pane zoom (fullscreen) |
| ` + "`Ctrl-b x`" + ` | Close current pane |
| ` + "`Ctrl-b {`" + ` | Swap pane up |
| ` + "`Ctrl-b }`" + ` | Swap pane down |
| ` + "`Ctrl-b Space`" + ` | Cycle pane layouts |
| Mouse drag | Resize pane borders |

## Copy Mode (vi)
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b [`" + ` | Enter copy mode |
| ` + "`v`" + ` | Start selection |
| ` + "`y`" + ` | Copy selection to clipboard |
| ` + "`q`" + ` | Exit copy mode |
| Mouse scroll | Scroll buffer (enters copy mode) |

## Persistence
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b Ctrl-s`" + ` | Save session |
| ` + "`Ctrl-b Ctrl-r`" + ` | Restore session |
| Auto-save | Every 10 minutes (continuum) |

## Nested tmux (SSH)
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b Ctrl-b`" + ` | Send prefix to inner tmux |

## First-Time Setup
1. Open tmux: ` + "`tmux`" + `
2. Install plugins: ` + "`Ctrl-b I`" + ` (capital I)
3. Wait for "TMUX environment reloaded"
4. Done! Resurrect + continuum are now active
`
}
