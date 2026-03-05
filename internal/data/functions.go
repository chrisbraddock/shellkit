package data

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Function represents a shell function.
type Function struct {
	Name        string
	Description string
	Source      string // full source code
}

// LoadFunctions reads all function files from the given directory.
func LoadFunctions(dir string) ([]Function, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var funcs []Function
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), "_") || strings.HasPrefix(e.Name(), ".") {
			continue
		}

		path := filepath.Join(dir, e.Name())
		f, err := parseFunction(path, e.Name())
		if err != nil {
			continue
		}
		funcs = append(funcs, f)
	}
	return funcs, nil
}

func parseFunction(path, name string) (Function, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Function{}, err
	}

	source := string(data)
	desc := extractDescription(source)

	return Function{
		Name:        name,
		Description: desc,
		Source:      source,
	}, nil
}

// extractDescription pulls description from the first comment lines.
func extractDescription(source string) string {
	scanner := bufio.NewScanner(strings.NewReader(source))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line == "#!/usr/bin/env zsh" || line == "#!/usr/bin/env bash" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			desc := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			// Skip lines that are just separators
			if strings.Trim(desc, "=-") == "" {
				continue
			}
			return desc
		}
		break
	}
	return ""
}
