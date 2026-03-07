package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/charmbracelet/glamour"
	"github.com/chrisbraddock/shellkit/internal/data"
)

type funcItem struct {
	fn data.Function
}

func (i funcItem) Title() string       { return i.fn.Name }
func (i funcItem) Description() string { return i.fn.Description }
func (i funcItem) FilterValue() string { return i.fn.Name + " " + i.fn.Description }

// FunctionTab is the functions tab model.
type FunctionTab struct {
	list     list.Model
	preview  viewport.Model
	funcs    []data.Function
	width    int
	height   int
	styles   *Styles
	renderer *glamour.TermRenderer
}

// NewFunctionTab creates the functions tab.
func NewFunctionTab(funcs []data.Function, styles *Styles) FunctionTab {
	items := make([]list.Item, len(funcs))
	for i, f := range funcs {
		items[i] = funcItem{fn: f}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Functions"
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	vp := viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))

	style := "dark"
	if !styles.IsDark {
		style = "light"
	}
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(60),
	)

	return FunctionTab{
		list:     l,
		preview:  vp,
		funcs:    funcs,
		styles:   styles,
		renderer: r,
	}
}

// AtTop returns true when the list cursor is at the first item.
func (t *FunctionTab) AtTop() bool { return t.list.Index() == 0 }

func (t *FunctionTab) SetStyles(s *Styles) {
	t.styles = s
	style := "dark"
	if !s.IsDark {
		style = "light"
	}
	t.renderer, _ = glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(60),
	)
}

func (t *FunctionTab) SetSize(w, h int) {
	t.width = w
	t.height = h

	listWidth := w * 45 / 100
	previewWidth := w - listWidth - 3

	t.list.SetSize(listWidth, h)
	t.preview.SetWidth(previewWidth - 4)
	t.preview.SetHeight(h - 4)

	if t.renderer != nil {
		t.renderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle(func() string {
				if t.styles.IsDark {
					return "dark"
				}
				return "light"
			}()),
			glamour.WithWordWrap(previewWidth-6),
		)
	}
}

func (t *FunctionTab) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update preview with source code
	if item, ok := t.list.SelectedItem().(funcItem); ok {
		md := "```zsh\n" + item.fn.Source + "\n```"
		if t.renderer != nil {
			if rendered, err := t.renderer.Render(md); err == nil {
				t.preview.SetContent(rendered)
			} else {
				t.preview.SetContent(item.fn.Source)
			}
		} else {
			t.preview.SetContent(item.fn.Source)
		}
	}

	t.preview, cmd = t.preview.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (t *FunctionTab) View() string {
	if t.width == 0 {
		return ""
	}

	listView := t.list.View()

	previewWidth := t.width - lipgloss.Width(listView) - 5
	previewHeight := t.height - 2

	// Bordered preview section
	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.styles.SubtleColor).
		Padding(0, 1).
		Width(previewWidth).
		Height(previewHeight)

	previewContent := t.preview.View()
	previewView := previewStyle.Render(previewContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, " ", previewView)
}
