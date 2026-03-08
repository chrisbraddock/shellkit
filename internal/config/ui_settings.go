package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UISettings stores persisted shellkit TUI preferences.
type UISettings struct {
	Header HeaderSettings `json:"header"`
}

// HeaderSettings stores persisted header and animation preferences.
type HeaderSettings struct {
	StartCollapsed       bool     `json:"start_collapsed"`
	ShowVersion          bool     `json:"show_version"`
	ShowPlatform         bool     `json:"show_platform"`
	ShowCompactTabAccent bool     `json:"show_compact_tab_accent"`
	EnabledAnimations    []string `json:"enabled_animations,omitempty"`
}

// DefaultUISettings returns the baseline TUI settings.
func DefaultUISettings() UISettings {
	return UISettings{
		Header: HeaderSettings{
			ShowVersion:          true,
			ShowPlatform:         true,
			ShowCompactTabAccent: false,
		},
	}
}

// LoadUISettings reads persisted UI settings, returning defaults when absent.
func LoadUISettings(metricsDir string) (UISettings, error) {
	settings := DefaultUISettings()
	path := UISettingsPath(metricsDir)
	if path == "" {
		return settings, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}
		return settings, err
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		return DefaultUISettings(), err
	}
	return settings, nil
}

// SaveUISettings persists UI settings to disk.
func SaveUISettings(metricsDir string, settings UISettings) error {
	path := UISettingsPath(metricsDir)
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// UISettingsPath returns the persisted TUI settings path.
func UISettingsPath(metricsDir string) string {
	if metricsDir == "" {
		return ""
	}
	return filepath.Join(metricsDir, "ui-config.json")
}
