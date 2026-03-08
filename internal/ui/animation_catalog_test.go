package ui

import (
	"strings"
	"testing"
)

func TestAnimationCatalogIDsAreUnique(t *testing.T) {
	seen := make(map[string]struct{}, len(animationCatalog))
	for _, meta := range animationCatalog {
		if meta.id == "" {
			t.Fatalf("animation %v has empty id", meta.mode)
		}
		if _, ok := seen[meta.id]; ok {
			t.Fatalf("duplicate animation id %q", meta.id)
		}
		seen[meta.id] = struct{}{}
	}
}

func TestAnimationSectionUsesCategory(t *testing.T) {
	if got := animationSection(animAuroraDrift); got != "Calm Animations" {
		t.Fatalf("animationSection(animAuroraDrift) = %q, want %q", got, "Calm Animations")
	}
}

func TestRenderAnimationPreviewReturnsFourRows(t *testing.T) {
	preview := renderAnimationPreview(animSoftRain, 24, 18, NewStyles(true))
	if strings.Count(preview, "\n") != 3 {
		t.Fatalf("renderAnimationPreview() row count mismatch: got %d newlines", strings.Count(preview, "\n"))
	}
}
