package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Config holds detected paths and platform info.
type Config struct {
	ZshDir       string // ~/.zsh
	AliasDir     string // ~/.zsh/aliases
	FunctionDir  string // ~/.zsh/functions
	ChezmoiSrc   string // chezmoi source-path
	OS           string // "darwin", "linux"
	Arch         string // "amd64", "arm64"
	HomeDir      string
	Version      string
}

// Detect builds a Config by inspecting the environment.
func Detect() Config {
	home, _ := os.UserHomeDir()

	c := Config{
		ZshDir:      filepath.Join(home, ".zsh"),
		AliasDir:    filepath.Join(home, ".zsh", "aliases"),
		FunctionDir: filepath.Join(home, ".zsh", "functions"),
		HomeDir:     home,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}

	// Detect chezmoi source path
	if out, err := exec.Command("chezmoi", "source-path").Output(); err == nil {
		c.ChezmoiSrc = strings.TrimSpace(string(out))
	}

	// Read version from chezmoi source
	if c.ChezmoiSrc != "" {
		if data, err := os.ReadFile(filepath.Join(c.ChezmoiSrc, "VERSION")); err == nil {
			c.Version = strings.TrimSpace(string(data))
		}
	}

	return c
}
