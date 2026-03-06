package ui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// Styles holds all shared styles for the TUI.
type Styles struct {
	Doc         lipgloss.Style
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style
	TabContent  lipgloss.Style
	Title       lipgloss.Style
	Subtle      lipgloss.Style
	Highlight   lipgloss.Style
	StatusOK    lipgloss.Style
	StatusFail  lipgloss.Style
	HelpBar     lipgloss.Style
	Category    lipgloss.Style
	Preview     lipgloss.Style

	// New styles for visual overhaul
	Logo         lipgloss.Style
	VersionBadge lipgloss.Style
	PlatformBadge lipgloss.Style
	TabDot       lipgloss.Style
	TabNumber    lipgloss.Style
	SectionTitle lipgloss.Style
	StatsBadge   lipgloss.Style
	TableHeader  lipgloss.Style
	TableRowAlt  lipgloss.Style
	Warning      lipgloss.Style
	Info         lipgloss.Style
	HelpKey      lipgloss.Style
	HelpSep      lipgloss.Style

	IsDark      bool
	AccentColor color.Color
	SubtleColor color.Color

	// Gradient stops
	GradientStart color.Color
	GradientEnd   color.Color

	// Metric health colors
	MetricFast   color.Color
	MetricMedium color.Color
	MetricSlow   color.Color

	// Category colors for search badges
	CategoryAlias    lipgloss.Style
	CategoryFunction lipgloss.Style
	CategoryPackage  lipgloss.Style
	CategoryKeybind  lipgloss.Style
}

// NewStyles creates a style set based on terminal background.
func NewStyles(isDark bool) *Styles {
	lightDark := lipgloss.LightDark(isDark)

	accent := lightDark(lipgloss.Color("#6C3FC5"), lipgloss.Color("#7D56F4"))
	subtle := lightDark(lipgloss.Color("#555555"), lipgloss.Color("#666666"))
	text := lightDark(lipgloss.Color("#1a1a1a"), lipgloss.Color("#FAFAFA"))
	dimText := lightDark(lipgloss.Color("#888888"), lipgloss.Color("#555555"))

	gradientStart := lightDark(lipgloss.Color("#6C3FC5"), lipgloss.Color("#7D56F4"))
	gradientEnd := lightDark(lipgloss.Color("#00B090"), lipgloss.Color("#00D4AA"))

	okColor := lightDark(lipgloss.Color("#0a6e0a"), lipgloss.Color("#73F59F"))
	failColor := lightDark(lipgloss.Color("#a00"), lipgloss.Color("#FF6B6B"))
	warnColor := lightDark(lipgloss.Color("#8B6914"), lipgloss.Color("#FFD93D"))
	infoColor := lightDark(lipgloss.Color("#2B6CB0"), lipgloss.Color("#6CB4FF"))

	// Category badge colors
	aliasColor := lightDark(lipgloss.Color("#C0392B"), lipgloss.Color("#F97583"))
	funcColor := lightDark(lipgloss.Color("#0E8A6E"), lipgloss.Color("#56D4DD"))
	pkgColor := lightDark(lipgloss.Color("#7D56F4"), lipgloss.Color("#D2A8FF"))
	keybindColor := lightDark(lipgloss.Color("#2B6CB0"), lipgloss.Color("#79C0FF"))

	s := &Styles{
		IsDark:        isDark,
		AccentColor:   accent,
		SubtleColor:   subtle,
		GradientStart: gradientStart,
		GradientEnd:   gradientEnd,
		MetricFast:    okColor,
		MetricMedium:  warnColor,
		MetricSlow:    failColor,
	}

	s.Doc = lipgloss.NewStyle().Padding(0, 1)

	// Active tab: bold accent text
	s.ActiveTab = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Padding(0, 2)

	// Inactive tab: subtle text
	s.InactiveTab = lipgloss.NewStyle().
		Foreground(subtle).
		Padding(0, 2)

	s.TabContent = lipgloss.NewStyle().
		Padding(1, 0)

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(accent)

	s.Subtle = lipgloss.NewStyle().
		Foreground(subtle)

	s.Highlight = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	s.StatusOK = lipgloss.NewStyle().
		Foreground(okColor)

	s.StatusFail = lipgloss.NewStyle().
		Foreground(failColor)

	s.HelpBar = lipgloss.NewStyle().
		Foreground(subtle)

	s.Category = lipgloss.NewStyle().
		Foreground(okColor).
		Bold(true)

	s.Preview = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(0, 1)

	// --- New styles ---

	s.Logo = lipgloss.NewStyle().
		Bold(true)

	s.VersionBadge = lipgloss.NewStyle().
		Foreground(text).
		Background(accent).
		Padding(0, 1).
		Bold(true)

	s.PlatformBadge = lipgloss.NewStyle().
		Foreground(dimText)

	s.TabDot = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	s.TabNumber = lipgloss.NewStyle().
		Foreground(dimText)

	s.SectionTitle = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	s.StatsBadge = lipgloss.NewStyle().
		Foreground(dimText)

	s.TableHeader = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Padding(0, 1)

	s.TableRowAlt = lipgloss.NewStyle().
		Padding(0, 1)

	s.Warning = lipgloss.NewStyle().
		Foreground(warnColor)

	s.Info = lipgloss.NewStyle().
		Foreground(infoColor)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(accent)

	s.HelpSep = lipgloss.NewStyle().
		Foreground(dimText)

	// Category badge styles
	s.CategoryAlias = lipgloss.NewStyle().
		Foreground(aliasColor).
		Bold(true)

	s.CategoryFunction = lipgloss.NewStyle().
		Foreground(funcColor).
		Bold(true)

	s.CategoryPackage = lipgloss.NewStyle().
		Foreground(pkgColor).
		Bold(true)

	s.CategoryKeybind = lipgloss.NewStyle().
		Foreground(keybindColor).
		Bold(true)

	return s
}

// RenderGradientLine renders a full-width gradient line using Blend1D.
func RenderGradientLine(width int, start, end color.Color) string {
	if width <= 0 {
		return ""
	}
	colors := lipgloss.Blend1D(width, start, end)
	var b strings.Builder
	for _, c := range colors {
		b.WriteString(lipgloss.NewStyle().Foreground(c).Render("━"))
	}
	return b.String()
}
