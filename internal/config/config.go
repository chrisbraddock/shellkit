package config

import (
	"os"
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
	MetricsDir   string
	Version      string
}

// Detect builds a Config by inspecting the environment.
// Avoids shelling out to keep startup fast.
func Detect() Config {
	home, _ := os.UserHomeDir()

	c := Config{
		ZshDir:      filepath.Join(home, ".zsh"),
		AliasDir:    filepath.Join(home, ".zsh", "aliases"),
		FunctionDir: filepath.Join(home, ".zsh", "functions"),
		HomeDir:     home,
		MetricsDir:  filepath.Join(home, ".local", "share", "shellkit"),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}

	// Try standard chezmoi source paths (avoids shelling out)
	for _, candidate := range []string{
		filepath.Join(home, ".local", "share", "chezmoi"),
		filepath.Join(home, ".config", "share", "chezmoi"),
	} {
		if _, err := os.Stat(candidate); err == nil {
			c.ChezmoiSrc = candidate
			break
		}
	}

	// Read version
	if c.ChezmoiSrc != "" {
		if data, err := os.ReadFile(filepath.Join(c.ChezmoiSrc, "VERSION")); err == nil {
			c.Version = strings.TrimSpace(string(data))
		}
	}

	return c
}
