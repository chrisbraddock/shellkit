package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// TabNames are the display names for each tab.
var TabNames = []string{
	"Dashboard",
	"Aliases",
	"Functions",
	"Packages",
	"Tmux",
	"Search",
	"Doctor",
}

// RenderTabBar renders a rich tab bar with dot indicator, numbers, and gradient accent line.
// When focused is true, the active tab gets an underline to indicate keyboard focus.
func RenderTabBar(activeIdx int, focused bool, width int, styles *Styles) string {
	var doc strings.Builder

	var tabs []string
	for i, name := range TabNames {
		var tab string
		if i == activeIdx {
			// Active tab
			dot := styles.TabDot.Render("●")
			activeStyle := styles.ActiveTab
			if focused {
				activeStyle = activeStyle.Underline(true)
			}
			label := activeStyle.Render(name)
			tab = dot + label
		} else {
			// Inactive: number + name (subtle)
			if width >= 100 {
				num := styles.TabNumber.Render(fmt.Sprintf("%d", i+1))
				label := styles.InactiveTab.Render(name)
				tab = num + label
			} else {
				tab = styles.InactiveTab.Render(name)
			}
		}
		tabs = append(tabs, tab)
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
	doc.WriteString(tabBar)
	doc.WriteString("\n")

	// Gradient accent line
	lineWidth := width - 2 // account for doc padding
	if lineWidth > 0 {
		doc.WriteString(RenderGradientLine(lineWidth, styles.GradientStart, styles.GradientEnd))
		doc.WriteString("\n")
	}

	return doc.String()
}

// RenderTabBarPlain renders a minimal tab bar for very narrow terminals.
func RenderTabBarPlain(activeIdx int, styles *Styles) string {
	var tabs []string
	for i, name := range TabNames {
		if i == activeIdx {
			tabs = append(tabs, styles.ActiveTab.Render(name))
		} else {
			tabs = append(tabs, styles.InactiveTab.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...) + "\n"
}
