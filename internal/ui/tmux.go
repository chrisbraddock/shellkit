package ui

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/charmbracelet/glamour"
	"github.com/chrisbraddock/shellkit/internal/data"
)

// TmuxTab shows the tmux reference page.
type TmuxTab struct {
	viewport viewport.Model
	width    int
	height   int
	styles   *Styles
}

// NewTmuxTab creates the tmux reference tab.
func NewTmuxTab(styles *Styles) TmuxTab {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))

	style := "dark"
	if !styles.IsDark {
		style = "light"
	}
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(78),
	)

	md := data.TmuxReference()
	if rendered, err := r.Render(md); err == nil {
		vp.SetContent(rendered)
	} else {
		vp.SetContent(md)
	}

	return TmuxTab{
		viewport: vp,
		styles:   styles,
	}
}

func (t *TmuxTab) SetStyles(s *Styles) {
	t.styles = s
	style := "dark"
	if !s.IsDark {
		style = "light"
	}
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(t.width-4),
	)
	md := data.TmuxReference()
	if rendered, err := r.Render(md); err == nil {
		t.viewport.SetContent(rendered)
	}
}

func (t *TmuxTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.viewport.SetWidth(w)
	t.viewport.SetHeight(h)
}

func (t *TmuxTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *TmuxTab) View() string {
	return t.viewport.View()
}
