package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

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

func (t *PackageTab) SetStyles(s *Styles) {
	t.styles = s
	t.viewport.SetContent(t.renderContent())
}

func (t *PackageTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.viewport.SetWidth(w)
	t.viewport.SetHeight(h)
}

func (t *PackageTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *PackageTab) View() string {
	return t.viewport.View()
}

func (t *PackageTab) renderContent() string {
	var b strings.Builder

	currentTier := ""
	for _, pkg := range t.packages {
		if pkg.Tier != currentTier {
			if currentTier != "" {
				b.WriteString("\n")
			}
			currentTier = pkg.Tier
			b.WriteString(t.styles.Title.Render(fmt.Sprintf("  %s", strings.ToUpper(currentTier))))
			b.WriteString("\n\n")
		}

		status := t.styles.StatusOK.Render("  ✓")
		if !pkg.Installed {
			status = t.styles.StatusFail.Render("  ✗")
		}

		name := fmt.Sprintf(" %-20s", pkg.Name)
		comment := ""
		if pkg.Comment != "" {
			comment = t.styles.Subtle.Render("  " + pkg.Comment)
		}

		b.WriteString(fmt.Sprintf("%s %s%s\n", status, name, comment))
	}

	return b.String()
}
