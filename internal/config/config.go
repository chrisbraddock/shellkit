package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const modulePath = "github.com/chrisbraddock/shellkit"

// BuildVersion can be injected by the main package for release builds.
var BuildVersion string

// Config holds detected paths and platform info.
type Config struct {
	ZshDir      string // ~/.zsh
	AliasDir    string // ~/.zsh/aliases
	FunctionDir string // ~/.zsh/functions
	ChezmoiSrc  string // chezmoi source-path
	OS          string // "darwin", "linux"
	Arch        string // "amd64", "arm64"
	HomeDir     string
	MetricsDir  string
	Version     string
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

	c.ChezmoiSrc = detectChezmoiSrc(home)
	c.Version = detectVersion(c.ChezmoiSrc)
	return c
}

func detectChezmoiSrc(home string) string {
	for _, candidate := range []string{
		filepath.Join(home, ".local", "share", "chezmoi"),
		filepath.Join(home, ".config", "share", "chezmoi"),
	} {
		if looksLikeShellkitRoot(candidate) {
			return candidate
		}
	}
	return ""
}

func detectVersion(chezmoiSrc string) string {
	for _, root := range versionRoots(chezmoiSrc) {
		if version := readVersionFromRoot(root); version != "" {
			return version
		}
	}
	return sanitizeVersion(BuildVersion)
}

func versionRoots(chezmoiSrc string) []string {
	var roots []string
	seen := make(map[string]struct{})

	add := func(root string) {
		root = strings.TrimSpace(root)
		if root == "" {
			return
		}
		root = filepath.Clean(root)
		if _, ok := seen[root]; ok {
			return
		}
		seen[root] = struct{}{}
		roots = append(roots, root)
	}

	add(sourceRoot())
	add(findShellkitRoot(chezmoiSrc))

	if exe, err := os.Executable(); err == nil {
		add(findShellkitRoot(filepath.Dir(exe)))
	}

	if wd, err := os.Getwd(); err == nil {
		add(findShellkitRoot(wd))
	}

	return roots
}

func sourceRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok || file == "" {
		return ""
	}

	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	if looksLikeShellkitRoot(root) {
		return root
	}
	return ""
}

func findShellkitRoot(start string) string {
	start = strings.TrimSpace(start)
	if start == "" {
		return ""
	}

	dir := filepath.Clean(start)
	for {
		if looksLikeShellkitRoot(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func looksLikeShellkitRoot(dir string) bool {
	if dir == "" {
		return false
	}

	if matchesModule(filepath.Join(dir, "go.mod")) {
		return true
	}

	name, _, err := readPackageJSON(filepath.Join(dir, "package.json"))
	return err == nil && name == "shellkit"
}

func matchesModule(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")) == modulePath
		}
	}
	return false
}

func readVersionFromRoot(root string) string {
	root = strings.TrimSpace(root)
	if root == "" {
		return ""
	}

	if data, err := os.ReadFile(filepath.Join(root, "VERSION")); err == nil {
		if version := sanitizeVersion(string(data)); version != "" {
			return version
		}
	}

	_, version, err := readPackageJSON(filepath.Join(root, "package.json"))
	if err != nil {
		return ""
	}
	return sanitizeVersion(version)
}

func readPackageJSON(path string) (string, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	var pkg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", "", err
	}
	return strings.TrimSpace(pkg.Name), strings.TrimSpace(pkg.Version), nil
}

func sanitizeVersion(version string) string {
	version = strings.TrimSpace(version)
	switch version {
	case "", "dev", "(devel)", "unknown":
		return ""
	default:
		return version
	}
}
