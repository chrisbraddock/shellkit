package ui

import (
	"strings"
	"testing"

	"github.com/chrisbraddock/shellkit/internal/config"
)

func TestRenderHeaderLineCountWithoutExtraBadgeRow(t *testing.T) {
	styles := NewStyles(true)
	state := NewHeaderState(styles, config.DefaultUISettings())

	rendered := RenderHeader("1.13.0", "darwin", "arm64", 100, styles, &state)
	if got := strings.Count(rendered, "\n"); got != HeaderLineCount {
		t.Fatalf("RenderHeader() line count = %d, want %d", got, HeaderLineCount)
	}
}

func TestRenderTabBarPlainUsesBracketedUppercaseLabels(t *testing.T) {
	styles := NewStyles(true)

	rendered := RenderTabBarPlain(0, styles)
	if !strings.Contains(rendered, "[DASHBOARD]") {
		t.Fatalf("RenderTabBarPlain() missing uppercase bracket label: %q", rendered)
	}
	if strings.Contains(rendered, "Dashboard") {
		t.Fatalf("RenderTabBarPlain() still contains mixed-case label: %q", rendered)
	}
}
