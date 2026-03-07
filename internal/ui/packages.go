package ui

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"

	"github.com/chrisbraddock/shellkit/internal/data"
)

// PackageTab displays package install status.
type PackageTab struct {
	viewport viewport.Model
	packages []data.Package
	width    int
	height   int
	styles   *Styles
}

// NewPackageTab creates the packages tab.
func NewPackageTab(pkgs []data.Package, styles *Styles) PackageTab {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	t := PackageTab{
		viewport: vp,
		packages: pkgs,
		styles:   styles,
	}
	t.viewport.SetContent(t.renderContent())
	return t
}

// AtTop returns true when the viewport is scrolled to the top.
func (t *PackageTab) AtTop() bool { return t.viewport.ScrollPercent() <= 0 }

func (t *PackageTab) SetStyles(s *Styles) {
	t.styles = s
	t.viewport.SetContent(t.renderContent())
}

func (t *PackageTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.viewport.SetWidth(w)
	t.viewport.SetHeight(h)
	t.viewport.SetContent(t.renderContent())
}

func (t *PackageTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *PackageTab) View() string {
	return t.viewport.View()
}

// Count returns the number of packages.
func (t *PackageTab) Count() int {
	return len(t.packages)
}

func (t *PackageTab) renderContent() string {
	if len(t.packages) == 0 {
		return t.styles.Subtle.Render("  Loading packages...")
	}

	// Group packages by tier
	tiers := make(map[string][]data.Package)
	var tierOrder []string
	for _, pkg := range t.packages {
		if _, seen := tiers[pkg.Tier]; !seen {
			tierOrder = append(tierOrder, pkg.Tier)
		}
		tiers[pkg.Tier] = append(tiers[pkg.Tier], pkg)
	}

	var b strings.Builder

	accentColor := t.styles.AccentColor
	subtleColor := t.styles.SubtleColor
	okColor := t.styles.StatusOK.GetForeground()
	failColor := t.styles.StatusFail.GetForeground()

	for _, tier := range tierOrder {
		pkgs := tiers[tier]

		// Build table for this tier
		tbl := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(subtleColor)).
			BorderColumn(true).
			BorderRow(false).
			Headers("", "PACKAGE", "DESCRIPTION").
			StyleFunc(func(row, col int) lipgloss.Style {
				s := lipgloss.NewStyle().Padding(0, 1)
				if row == table.HeaderRow {
					return s.Foreground(accentColor).Bold(true)
				}
				// Status column coloring
				if col == 0 && row >= 0 && row < len(pkgs) {
					if pkgs[row].Installed {
						return s.Foreground(okColor)
					}
					return s.Foreground(failColor)
				}
				// Alternating row tint (subtle)
				if row%2 == 1 {
					return s.Foreground(subtleColor)
				}
				return s
			})

		// Set width if available
		if t.width > 4 {
			tbl = tbl.Width(t.width - 2)
		}

		// Add rows
		for _, pkg := range pkgs {
			status := "✓"
			if !pkg.Installed {
				status = "✗"
			}
			comment := pkg.Comment
			if comment == "" {
				comment = pkg.Tier
			}
			tbl = tbl.Row(status, pkg.Name, comment)
		}

		// Tier header
		b.WriteString(t.styles.Title.Render("  " + strings.ToUpper(tier)))
		b.WriteString("\n")
		b.WriteString(tbl.Render())
		b.WriteString("\n\n")
	}

	return b.String()
}
