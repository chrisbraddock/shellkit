package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadVersionFromRootPrefersVersionFile(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "VERSION"), []byte("1.2.3\n"), 0o644); err != nil {
		t.Fatalf("write VERSION: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"name":"shellkit","version":"9.9.9"}`), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}

	if got := readVersionFromRoot(root); got != "1.2.3" {
		t.Fatalf("readVersionFromRoot() = %q, want %q", got, "1.2.3")
	}
}

func TestReadVersionFromRootFallsBackToPackageJSON(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"name":"shellkit","version":"2.4.6"}`), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}

	if got := readVersionFromRoot(root); got != "2.4.6" {
		t.Fatalf("readVersionFromRoot() = %q, want %q", got, "2.4.6")
	}
}

func TestLooksLikeShellkitRootFromGoMod(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module "+modulePath+"\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	if !looksLikeShellkitRoot(root) {
		t.Fatal("looksLikeShellkitRoot() = false, want true")
	}
}

func TestSanitizeVersion(t *testing.T) {
	if got := sanitizeVersion(" dev "); got != "" {
		t.Fatalf("sanitizeVersion(dev) = %q, want empty", got)
	}
	if got := sanitizeVersion("1.13.0"); got != "1.13.0" {
		t.Fatalf("sanitizeVersion(1.13.0) = %q, want %q", got, "1.13.0")
	}
}
