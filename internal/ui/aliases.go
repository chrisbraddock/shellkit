package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/chrisbraddock/shellkit/internal/data"
)

// aliasItem implements list.DefaultItem for the alias list.
type aliasItem struct {
	alias data.Alias
}

func (i aliasItem) Title() string       { return i.alias.Name }
func (i aliasItem) FilterValue() string { return i.alias.Name + " " + i.alias.Command }
func (i aliasItem) Description() string {
	desc := i.alias.Command
	if i.alias.Comment != "" {
		desc = i.alias.Comment
	}
	return desc
}

// AliasTab is the aliases tab model.
type AliasTab struct {
	list    list.Model
	preview viewport.Model
	aliases []data.Alias
	width   int
	height  int
	styles  *Styles
}

// NewAliasTab creates the aliases tab.
func NewAliasTab(aliases []data.Alias, styles *Styles) AliasTab {
	items := make([]list.Item, len(aliases))
	for i, a := range aliases {
		items[i] = aliasItem{alias: a}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Aliases"
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	vp := viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))

	return AliasTab{
		list:    l,
		preview: vp,
		aliases: aliases,
		styles:  styles,
	}
}

func (t *AliasTab) SetStyles(s *Styles) { t.styles = s }

func (t *AliasTab) SetSize(w, h int) {
	t.width = w
	t.height = h

	listWidth := w * 55 / 100
	previewWidth := w - listWidth - 3 // gap

	t.list.SetSize(listWidth, h)
	t.preview.SetWidth(previewWidth - 4) // border padding
	t.preview.SetHeight(h - 4)
}

func (t *AliasTab) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update preview content based on selection
	if item, ok := t.list.SelectedItem().(aliasItem); ok {
		t.preview.SetContent(t.renderPreview(item.alias))
	}

	t.preview, cmd = t.preview.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (t *AliasTab) View() string {
	if t.width == 0 {
		return ""
	}

	listView := t.list.View()

	previewStyle := t.styles.Preview.
		Width(t.width - lipgloss.Width(listView) - 5).
		Height(t.height - 2)

	previewView := previewStyle.Render(t.preview.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, " ", previewView)
}

func (t *AliasTab) renderPreview(a data.Alias) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  Alias:     %s\n", a.Name)
	fmt.Fprintf(&b, "  Command:   %s\n", a.Command)
	fmt.Fprintf(&b, "  Category:  %s\n", a.Category)
	if a.Comment != "" {
		fmt.Fprintf(&b, "  Note:      %s\n", a.Comment)
	}
	return b.String()
}
