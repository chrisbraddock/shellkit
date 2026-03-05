package data

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Package represents a managed package.
type Package struct {
	Name      string
	Tier      string // "shell", "editor", "tools", "monitoring", "npm", "pip", etc.
	Installed bool
	Comment   string
}

type packagesYAML struct {
	Packages struct {
		Shell      []string `yaml:"shell"`
		Editor     []string `yaml:"editor"`
		Tools      []string `yaml:"tools"`
		Monitoring []string `yaml:"monitoring"`
		NPM        []string `yaml:"npm"`
		Pip        []string `yaml:"pip"`
		Darwin     struct {
			Brew []string `yaml:"brew"`
			Cask []string `yaml:"cask"`
		} `yaml:"darwin"`
		Linux struct {
			APT         []string          `yaml:"apt"`
			Brew        []string          `yaml:"brew"`
			APTMappings map[string]string `yaml:"apt_mappings"`
		} `yaml:"linux"`
	} `yaml:"packages"`
}

// LoadPackages reads packages.yaml and checks install status.
func LoadPackages(chezmoiSrc, osType string) ([]Package, error) {
	path := filepath.Join(chezmoiSrc, ".chezmoidata", "packages.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw packagesYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var pkgs []Package

	addTier := func(tier string, names []string) {
		for _, name := range names {
			clean, comment := splitComment(name)
			pkgs = append(pkgs, Package{
				Name:      clean,
				Tier:      tier,
				Installed: isInstalled(clean),
				Comment:   comment,
			})
		}
	}

	addTier("shell", raw.Packages.Shell)
	addTier("editor", raw.Packages.Editor)
	addTier("tools", raw.Packages.Tools)
	addTier("monitoring", raw.Packages.Monitoring)
	addTier("npm", raw.Packages.NPM)
	addTier("pip", raw.Packages.Pip)

	if osType == "darwin" {
		addTier("darwin/brew", raw.Packages.Darwin.Brew)
		addTier("darwin/cask", raw.Packages.Darwin.Cask)
	} else {
		addTier("linux/apt", raw.Packages.Linux.APT)
		addTier("linux/brew", raw.Packages.Linux.Brew)
	}

	return pkgs, nil
}

// splitComment separates "name  # comment" into name and comment.
func splitComment(s string) (string, string) {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "#"); idx > 0 {
		return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+1:])
	}
	return s, ""
}

// commandName maps package names to their actual binary names.
var commandName = map[string]string{
	"neovim":        "nvim",
	"ripgrep":       "rg",
	"git-delta":     "delta",
	"1password-cli": "op",
	"git-lfs":       "git-lfs",
	"git-cliff":     "git-cliff",
}

func isInstalled(name string) bool {
	cmd := name
	if mapped, ok := commandName[name]; ok {
		cmd = mapped
	}
	_, err := exec.LookPath(cmd)
	return err == nil
}
