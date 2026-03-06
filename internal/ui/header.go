package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var asciiLogo = []string{
	` ╔═╗ ╦ ╦ ╔═╗ ╦   ╦   ╦╔═ ╦ ╔╦╗`,
	` ╚═╗ ╠═╣ ║╣  ║   ║   ╠╩╗ ║  ║ `,
	` ╚═╝ ╩ ╩ ╚═╝ ╩═╝ ╩═╝ ╩ ╩ ╩  ╩ `,
}

// RenderHeader renders the gradient ASCII logo with version and platform info.
func RenderHeader(version, os, arch string, width int, styles *Styles) string {
	var doc strings.Builder

	// Render each line of the ASCII logo with gradient colors
	logoWidth := 0
	for _, line := range asciiLogo {
		if len([]rune(line)) > logoWidth {
			logoWidth = len([]rune(line))
		}
	}

	colors := lipgloss.Blend1D(logoWidth, styles.GradientStart, styles.GradientEnd)

	for _, line := range asciiLogo {
		runes := []rune(line)
		var lineBuilder strings.Builder
		for i, r := range runes {
			if r == ' ' {
				lineBuilder.WriteRune(' ')
				continue
			}
			ci := i
			if ci >= len(colors) {
				ci = len(colors) - 1
			}
			lineBuilder.WriteString(
				lipgloss.NewStyle().
					Foreground(colors[ci]).
					Bold(true).
					Render(string(r)),
			)
		}
		doc.WriteString(lineBuilder.String())
		doc.WriteString("\n")
	}

	// Version + platform badge on the right of the last logo line
	if version != "" || os != "" {
		var badge strings.Builder
		if version != "" {
			badge.WriteString(styles.VersionBadge.Render(version))
		}
		if os != "" && arch != "" {
			badge.WriteString(" ")
			badge.WriteString(styles.PlatformBadge.Render(os + "/" + arch))
		}

		// Position badge at the right side
		badgeStr := badge.String()
		logoRenderedWidth := lipgloss.Width(asciiLogo[0]) // approximate
		gap := width - logoRenderedWidth - lipgloss.Width(badgeStr) - 4
		if gap < 2 {
			gap = 2
		}

		doc.WriteString(strings.Repeat(" ", gap))
		doc.WriteString(badgeStr)
		doc.WriteString("\n")
	}

	return doc.String()
}
