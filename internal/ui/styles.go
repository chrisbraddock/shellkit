package ui

import (
	"charm.land/lipgloss/v2"
)

// Styles holds all shared styles for the TUI.
type Styles struct {
	Doc         lipgloss.Style
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style
	TabGap      lipgloss.Style
	TabContent  lipgloss.Style
	Title       lipgloss.Style
	Subtle      lipgloss.Style
	Highlight   lipgloss.Style
	StatusOK    lipgloss.Style
	StatusFail  lipgloss.Style
	HelpBar     lipgloss.Style
	Category    lipgloss.Style
	Preview     lipgloss.Style
	IsDark      bool
}

// NewStyles creates a style set based on terminal background.
func NewStyles(isDark bool) *Styles {
	lightDark := lipgloss.LightDark(isDark)

	accent := lightDark(lipgloss.Color("#6C3FC5"), lipgloss.Color("#7D56F4"))
	subtle := lightDark(lipgloss.Color("#555555"), lipgloss.Color("#666666"))
	text := lightDark(lipgloss.Color("#1a1a1a"), lipgloss.Color("#FAFAFA"))

	activeTabBorder := lipgloss.RoundedBorder()
	activeTabBorder.BottomLeft = "│"
	activeTabBorder.Bottom = " "
	activeTabBorder.BottomRight = "│"

	inactiveTabBorder := lipgloss.RoundedBorder()
	inactiveTabBorder.BottomLeft = "┴"
	inactiveTabBorder.Bottom = "─"
	inactiveTabBorder.BottomRight = "┴"

	tabGapBorder := lipgloss.Border{Bottom: "─", BottomLeft: "─", BottomRight: "─"}

	s := &Styles{IsDark: isDark}

	s.Doc = lipgloss.NewStyle().Padding(0, 1)

	s.ActiveTab = lipgloss.NewStyle().
		Border(activeTabBorder, true).
		BorderForeground(accent).
		Foreground(text).
		Bold(true).
		Padding(0, 2)

	s.InactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(subtle).
		Foreground(subtle).
		Padding(0, 2)

	s.TabGap = lipgloss.NewStyle().
		Border(tabGapBorder, false, false, true, false).
		BorderForeground(subtle)

	s.TabContent = lipgloss.NewStyle().
		Padding(1, 2)

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(accent)

	s.Subtle = lipgloss.NewStyle().
		Foreground(subtle)

	s.Highlight = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	s.StatusOK = lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#0a6e0a"), lipgloss.Color("#73F59F")))

	s.StatusFail = lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#a00"), lipgloss.Color("#FF6B6B")))

	s.HelpBar = lipgloss.NewStyle().
		Foreground(subtle).
		Padding(1, 0, 0, 0)

	s.Category = lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#0a6e0a"), lipgloss.Color("#73F59F"))).
		Bold(true)

	s.Preview = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(0, 1)

	return s
}
