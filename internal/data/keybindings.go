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
		// Tmux prefix bindings (Oh my tmux! + shellkit)
		{Key: "Ctrl-b |", Description: "Split pane horizontally", Category: "tmux"},
		{Key: "Ctrl-b -", Description: "Split pane vertically", Category: "tmux"},
		{Key: "Ctrl-b h/j/k/l", Description: "Navigate panes (vim-style)", Category: "tmux"},
		{Key: "Ctrl-b H/J/K/L", Description: "Resize panes", Category: "tmux"},
		{Key: "Ctrl-b >/<", Description: "Swap pane next/prev", Category: "tmux"},
		{Key: "Ctrl-b +", Description: "Maximize pane", Category: "tmux"},
		{Key: "Ctrl-b c", Description: "New window", Category: "tmux"},
		{Key: "Ctrl-b C-h/C-l", Description: "Previous/next window", Category: "tmux"},
		{Key: "Ctrl-b Tab", Description: "Last active window", Category: "tmux"},
		{Key: "Ctrl-b 1-9", Description: "Jump to window by number", Category: "tmux"},
		{Key: "Ctrl-b ,", Description: "Rename current window", Category: "tmux"},
		{Key: "Ctrl-b w", Description: "Window/session picker", Category: "tmux"},
		{Key: "Ctrl-b C-c", Description: "New session", Category: "tmux"},
		{Key: "Ctrl-b C-f", Description: "Find session", Category: "tmux"},
		{Key: "Ctrl-b BTab", Description: "Last session", Category: "tmux"},
		{Key: "Ctrl-b d", Description: "Detach from session", Category: "tmux"},
		{Key: "Ctrl-b Enter", Description: "Enter copy mode", Category: "tmux"},
		{Key: "Ctrl-b b", Description: "List paste buffers", Category: "tmux"},
		{Key: "Ctrl-b P", Description: "Choose buffer to paste", Category: "tmux"},
		{Key: "Ctrl-b S", Description: "Paste a tmux snippet (fzf)", Category: "tmux"},
		{Key: "Ctrl-b ?", Description: "Shellkit quick reference popup", Category: "tmux"},
		{Key: "Ctrl-b m", Description: "Toggle mouse", Category: "tmux"},
		{Key: "Ctrl-b e", Description: "Edit tmux config", Category: "tmux"},
		{Key: "Ctrl-b r", Description: "Reload tmux config", Category: "tmux"},
		{Key: "Ctrl-b I", Description: "Install/update TPM plugins", Category: "tmux"},
		{Key: "Ctrl-b Ctrl-s", Description: "Save session (resurrect)", Category: "tmux"},
		{Key: "Ctrl-b Ctrl-r", Description: "Restore session (resurrect)", Category: "tmux"},
		{Key: "Ctrl-b Ctrl-b", Description: "Send prefix to nested tmux", Category: "tmux"},
		{Key: "C-l", Description: "Clear screen and scrollback", Category: "tmux"},

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
	return `# Tmux Quick Reference (Oh my tmux! + shellkit)

## Concept: Session > Window > Pane

- **Session** = workspace (e.g. "main", "project-x")
- **Window** = tab within a session (shown in status bar)
- **Pane** = split within a window

## Sessions
| Command / Key | Description |
|---------------|-------------|
| ` + "`ta`" + ` | Attach to most recent session |
| ` + "`tas <name>`" + ` | Attach to named session |
| ` + "`tns <name>`" + ` | New named session |
| ` + "`tls`" + ` | List sessions |
| ` + "`tks <name>`" + ` | Kill named session |
| ` + "`sshx <host>`" + ` | SSH + auto-attach tmux |
| ` + "`tw <path>`" + ` | Open repo workspace |
| ` + "`tdev <path>`" + ` | Claude + Codex + shell workspace |
| ` + "`tsys <path>`" + ` | 4-pane claude workspace |
| ` + "`Ctrl-b C-c`" + ` | New session |
| ` + "`Ctrl-b C-f`" + ` | Find session |
| ` + "`Ctrl-b BTab`" + ` | Switch to last session |
| ` + "`Ctrl-b d`" + ` | Detach |

## Windows (Tabs)
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b c`" + ` | New window |
| ` + "`Ctrl-b C-h / C-l`" + ` | Previous / next window |
| ` + "`Ctrl-b Tab`" + ` | Last active window |
| ` + "`Ctrl-b 1-9`" + ` | Jump to window by number |
| ` + "`Ctrl-b ,`" + ` | Rename current window |
| ` + "`Ctrl-b w`" + ` | Window/session picker |
| ` + "`Ctrl-b &`" + ` | Close current window |

## Panes (Splits)
| Key | Description |
|-----|-------------|
| ` + "` Ctrl-b \\| `" + ` | Split horizontally |
| ` + "`Ctrl-b -`" + ` | Split vertically |
| ` + "`Ctrl-b h/j/k/l`" + ` | Navigate panes (vim-style) |
| ` + "`Ctrl-b H/J/K/L`" + ` | Resize panes |
| ` + "`Ctrl-b z`" + ` | Toggle pane zoom (fullscreen) |
| ` + "`Ctrl-b +`" + ` | Maximize pane |
| ` + "`Ctrl-b x`" + ` | Close current pane |
| ` + "`Ctrl-b > / <`" + ` | Swap pane next / prev |
| ` + "`Ctrl-b Space`" + ` | Cycle pane layouts |
| Mouse drag | Resize pane borders |

## Snippets
| Key / Command | Description |
|---------------|-------------|
| ` + "`Ctrl-b S`" + ` | fzf pick snippet → paste into pane |
| ` + "`snip new <name>`" + ` | Create a new snippet |
| ` + "`snip edit <name>`" + ` | Edit an existing snippet |
| ` + "`snip rm <name>`" + ` | Delete a snippet |
| ` + "`snip cp <name>`" + ` | Copy snippet to clipboard |
| ` + "`snip`" + ` | Browse all snippets (fzf) |

## Copy Mode (vi)
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b Enter`" + ` | Enter copy mode |
| ` + "`Ctrl-b [`" + ` | Enter copy mode (alternate) |
| ` + "`v`" + ` | Start selection |
| ` + "`C-v`" + ` | Rectangle (block) selection |
| ` + "`y`" + ` | Copy to clipboard |
| ` + "`Escape`" + ` | Cancel |
| Mouse scroll | Scroll buffer (enters copy mode) |

## Buffers
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b b`" + ` | List paste buffers |
| ` + "`Ctrl-b p`" + ` | Paste from top buffer |
| ` + "`Ctrl-b P`" + ` | Choose buffer to paste |

## Persistence
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b Ctrl-s`" + ` | Save session |
| ` + "`Ctrl-b Ctrl-r`" + ` | Restore session |
| Auto-save | Every 5 minutes (continuum) |
| Pane capture | Saved and restored with resurrect |
| Neovim restore | Sessions restored through tmux-resurrect |
| macOS boot restore | LaunchAgent starts tmux server at login |

## Other
| Key | Description |
|-----|-------------|
| ` + "`Ctrl-b Ctrl-b`" + ` | Send prefix to nested tmux |
| ` + "`Ctrl-b ?`" + ` | Quick reference popup |
| ` + "`Ctrl-b e`" + ` | Edit tmux config |
| ` + "`Ctrl-b r`" + ` | Reload tmux config |
| ` + "`Ctrl-b m`" + ` | Toggle mouse |
| ` + "`C-l`" + ` | Clear screen and scrollback |
| ` + "`Ctrl-b I`" + ` | Install/update TPM plugins |

## First-Time Setup
1. Run ` + "`chezmoi apply`" + ` (installs Oh my tmux! framework)
2. Open tmux: ` + "`tmux`" + `
3. Plugins auto-install on first launch
4. Done! Resurrect + continuum are active
`
}
