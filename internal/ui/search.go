package ui

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/chrisbraddock/shellkit/internal/data"
)

type searchItem struct {
	category string
	name     string
	desc     string
}

func (i searchItem) Title() string       { return fmt.Sprintf("[%s]  %s", i.category, i.name) }
func (i searchItem) Description() string { return i.desc }
func (i searchItem) FilterValue() string { return i.name + " " + i.desc + " " + i.category }

// SearchTab provides unified search across all shellkit content.
type SearchTab struct {
	list   list.Model
	width  int
	height int
	styles *Styles
}

// NewSearchTab creates the search tab with all content indexed.
func NewSearchTab(aliases []data.Alias, funcs []data.Function, pkgs []data.Package, keybindings []data.Keybinding, styles *Styles) SearchTab {
	var items []list.Item

	for _, a := range aliases {
		desc := a.Command
		if a.Comment != "" {
			desc = a.Comment
		}
		items = append(items, searchItem{
			category: "alias",
			name:     a.Name,
			desc:     desc,
		})
	}

	for _, f := range funcs {
		items = append(items, searchItem{
			category: "function",
			name:     f.Name,
			desc:     f.Description,
		})
	}

	for _, p := range pkgs {
		desc := p.Tier
		if p.Comment != "" {
			desc += " — " + p.Comment
		}
		items = append(items, searchItem{
			category: "package",
			name:     p.Name,
			desc:     desc,
		})
	}

	for _, kb := range keybindings {
		items = append(items, searchItem{
			category: "keybind",
			name:     kb.Key,
			desc:     kb.Description,
		})
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Search"
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	return SearchTab{
		list:   l,
		styles: styles,
	}
}

func (t *SearchTab) SetStyles(s *Styles) { t.styles = s }

func (t *SearchTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.list.SetSize(w, h)
}

func (t *SearchTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg)
	return cmd
}

func (t *SearchTab) View() string {
	return t.list.View()
}
