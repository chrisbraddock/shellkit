package data

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Alias represents a single shell alias.
type Alias struct {
	Name     string
	Command  string
	Category string // filename without extension (e.g. "git", "system", "tmux")
	Comment  string // inline comment if any
}

var aliasRe = regexp.MustCompile(`^\s*alias\s+([a-zA-Z0-9_-]+)=(.+)`)

// LoadAliases reads all alias files from the given directory.
func LoadAliases(dir string) ([]Alias, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var aliases []Alias
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".zsh") && !strings.HasSuffix(name, ".zsh.tmpl") {
			continue
		}

		category := strings.TrimSuffix(strings.TrimSuffix(name, ".tmpl"), ".zsh")
		path := filepath.Join(dir, name)
		parsed, err := parseAliasFile(path, category)
		if err != nil {
			continue // skip files that can't be read
		}
		aliases = append(aliases, parsed...)
	}
	return aliases, nil
}

func parseAliasFile(path, category string) ([]Alias, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var aliases []Alias
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		matches := aliasRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		name := matches[1]
		raw := matches[2]

		// Strip surrounding quotes from command
		cmd := stripQuotes(raw)

		// Check for trailing comment
		var comment string
		if idx := strings.Index(raw, "# "); idx > 0 {
			comment = strings.TrimSpace(raw[idx+2:])
			// Re-strip the command portion
			cmd = stripQuotes(strings.TrimSpace(raw[:idx]))
		}

		aliases = append(aliases, Alias{
			Name:     name,
			Command:  cmd,
			Category: category,
			Comment:  comment,
		})
	}
	return aliases, scanner.Err()
}

func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
