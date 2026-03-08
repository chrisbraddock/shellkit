package config

import "testing"

func TestDefaultUISettingsCompactAccentDisabled(t *testing.T) {
	settings := DefaultUISettings()
	if settings.Header.ShowCompactTabAccent {
		t.Fatal("DefaultUISettings() enables compact tab accent, want disabled")
	}
	if settings.Header.AnimationSpeed != DefaultAnimationSpeed {
		t.Fatalf("DefaultUISettings() animation speed = %d, want %d", settings.Header.AnimationSpeed, DefaultAnimationSpeed)
	}
}

func TestUISaveLoadRoundTripPreservesCompactAccent(t *testing.T) {
	metricsDir := t.TempDir()
	settings := DefaultUISettings()
	settings.Header.StartCollapsed = true
	settings.Header.ShowCompactTabAccent = true
	settings.Header.AnimationSpeed = 140

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
	if got.Header.AnimationSpeed != 140 {
		t.Fatalf("LoadUISettings() animation speed = %d, want %d", got.Header.AnimationSpeed, 140)
	}
}

func TestLoadUISettingsRepairsMissingAnimationSpeed(t *testing.T) {
	metricsDir := t.TempDir()
	settings := DefaultUISettings()
	settings.Header.AnimationSpeed = 0

	if err := SaveUISettings(metricsDir, settings); err != nil {
		t.Fatalf("SaveUISettings() error = %v", err)
	}

	got, err := LoadUISettings(metricsDir)
	if err != nil {
		t.Fatalf("LoadUISettings() error = %v", err)
	}
	if got.Header.AnimationSpeed != DefaultAnimationSpeed {
		t.Fatalf("LoadUISettings() animation speed = %d, want %d", got.Header.AnimationSpeed, DefaultAnimationSpeed)
	}
}
