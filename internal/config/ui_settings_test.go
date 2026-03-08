package config

import "testing"

func TestDefaultUISettingsCompactAccentDisabled(t *testing.T) {
	settings := DefaultUISettings()
	if settings.Header.ShowCompactTabAccent {
		t.Fatal("DefaultUISettings() enables compact tab accent, want disabled")
	}
}

func TestUISaveLoadRoundTripPreservesCompactAccent(t *testing.T) {
	metricsDir := t.TempDir()
	settings := DefaultUISettings()
	settings.Header.StartCollapsed = true
	settings.Header.ShowCompactTabAccent = true

	if err := SaveUISettings(metricsDir, settings); err != nil {
		t.Fatalf("SaveUISettings() error = %v", err)
	}

	got, err := LoadUISettings(metricsDir)
	if err != nil {
		t.Fatalf("LoadUISettings() error = %v", err)
	}
	if !got.Header.ShowCompactTabAccent {
		t.Fatal("LoadUISettings() lost compact tab accent setting")
	}
	if !got.Header.StartCollapsed {
		t.Fatal("LoadUISettings() lost start-collapsed setting")
	}
}
