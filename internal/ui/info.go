package ui

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/charmbracelet/glamour"
	"github.com/chrisbraddock/shellkit/internal/data"
)

// InfoTab shows system info and tool versions.
type InfoTab struct {
	viewport viewport.Model
	info     data.SystemInfo
	version  string
	width    int
	height   int
	styles   *Styles
}

// NewInfoTab creates the info tab.
func NewInfoTab(info data.SystemInfo, version string, styles *Styles) InfoTab {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))

	t := InfoTab{
		viewport: vp,
		info:     info,
		version:  version,
		styles:   styles,
	}
	t.renderContent()
	return t
}

func (t *InfoTab) SetStyles(s *Styles) {
	t.styles = s
	t.renderContent()
}

func (t *InfoTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.viewport.SetWidth(w)
	t.viewport.SetHeight(h)
}

func (t *InfoTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *InfoTab) View() string {
	return t.viewport.View()
}

func (t *InfoTab) renderContent() {
	md := data.FormatSystemInfo(t.info, t.version)
	style := "dark"
	if !t.styles.IsDark {
		style = "light"
	}
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(t.width-4),
	)
	if rendered, err := r.Render(md); err == nil {
		t.viewport.SetContent(rendered)
	} else {
		t.viewport.SetContent(md)
	}
}
