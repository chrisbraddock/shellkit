package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// RenderStatusBar renders a styled status bar with keybinding hints and context stats.
func RenderStatusBar(activeIdx int, stats map[string]string, width int, styles *Styles) string {
	// Left side: keybinding hints
	sep := styles.HelpSep.Render(" · ")
	hints := []string{
		styles.HelpKey.Render("tab") + styles.HelpBar.Render(" switch"),
		styles.HelpKey.Render("1-7") + styles.HelpBar.Render(" jump"),
		styles.HelpKey.Render("/") + styles.HelpBar.Render(" filter"),
		styles.HelpKey.Render("q") + styles.HelpBar.Render(" quit"),
	}
	left := "  " + strings.Join(hints, sep)

	// Right side: context-aware stats
	var right string
	if stat, ok := stats[TabNames[activeIdx]]; ok && stat != "" {
		right = styles.StatsBadge.Render(stat) + "  "
	}

	// Calculate padding between left and right
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	pad := width - leftWidth - rightWidth - 2 // doc padding
	if pad < 1 {
		pad = 1
	}

	return fmt.Sprintf("%s%s%s", left, strings.Repeat(" ", pad), right)
}
