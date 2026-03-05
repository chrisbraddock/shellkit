package data

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ToolStatus represents the install status of a tool.
type ToolStatus struct {
	Name      string
	Installed bool
	Version   string
}

// SystemInfo holds system-level information.
type SystemInfo struct {
	OS       string
	Arch     string
	Shell    string
	Terminal string
	Tools    []ToolStatus
}

// LoadSystemInfo gathers system information.
func LoadSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Shell:    os.Getenv("SHELL"),
		Terminal: os.Getenv("TERM_PROGRAM"),
	}

	tools := []struct {
		name    string
		cmd     string
		verFlag string
	}{
		{"chezmoi", "chezmoi", "--version"},
		{"zsh", "zsh", "--version"},
		{"fzf", "fzf", "--version"},
		{"neovim", "nvim", "--version"},
		{"tmux", "tmux", "-V"},
		{"git", "git", "--version"},
		{"brew", "brew", "--version"},
		{"atuin", "atuin", "--version"},
		{"zoxide", "zoxide", "--version"},
		{"direnv", "direnv", "version"},
		{"bat", "bat", "--version"},
		{"ripgrep", "rg", "--version"},
		{"fd", "fd", "--version"},
		{"eza", "eza", "--version"},
		{"delta", "delta", "--version"},
		{"pet", "pet", "version"},
	}

	for _, t := range tools {
		ts := ToolStatus{Name: t.name}
		if _, err := exec.LookPath(t.cmd); err == nil {
			ts.Installed = true
			if out, err := exec.Command(t.cmd, t.verFlag).Output(); err == nil {
				ver := strings.TrimSpace(string(out))
				// Take first line only
				if idx := strings.IndexByte(ver, '\n'); idx > 0 {
					ver = ver[:idx]
				}
				ts.Version = ver
			}
		}
		info.Tools = append(info.Tools, ts)
	}

	return info
}

// FormatSystemInfo returns a markdown-formatted system info string.
func FormatSystemInfo(info SystemInfo, version string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Shellkit Info\n\n")
	if version != "" {
		fmt.Fprintf(&b, "**Version:** %s\n\n", version)
	}
	fmt.Fprintf(&b, "**OS:** %s/%s\n", info.OS, info.Arch)
	fmt.Fprintf(&b, "**Shell:** %s\n", info.Shell)
	if info.Terminal != "" {
		fmt.Fprintf(&b, "**Terminal:** %s\n", info.Terminal)
	}
	fmt.Fprintf(&b, "\n## Installed Tools\n\n")
	for _, t := range info.Tools {
		if t.Installed {
			if t.Version != "" {
				fmt.Fprintf(&b, "- **%s** — %s\n", t.Name, t.Version)
			} else {
				fmt.Fprintf(&b, "- **%s** — installed\n", t.Name)
			}
		} else {
			fmt.Fprintf(&b, "- ~~%s~~ — not installed\n", t.Name)
		}
	}
	return b.String()
}
