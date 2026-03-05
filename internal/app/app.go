package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/chrisbraddock/shellkit/internal/config"
	"github.com/chrisbraddock/shellkit/internal/data"
	"github.com/chrisbraddock/shellkit/internal/ui"
)

type tabID int

const (
	tabAliases tabID = iota
	tabFunctions
	tabPackages
	tabTmux
	tabSearch
	tabInfo
	tabDoctor
)

var tabNames = []string{
	"Aliases",
	"Functions",
	"Packages",
	"Tmux",
	"Search",
	"Info",
	"Doctor",
}

// Messages for async data loading
type packagesLoadedMsg struct{ pkgs []data.Package }
type sysInfoLoadedMsg struct{ info data.SystemInfo }
type doctorLoadedMsg struct{}

// Model is the root Bubble Tea model.
type Model struct {
	cfg       config.Config
	activeTab tabID
	width     int
	height    int
	isDark    bool
	styles    *ui.Styles
	ready     bool

	aliases   ui.AliasTab
	functions ui.FunctionTab
	packages  ui.PackageTab
	tmux      ui.TmuxTab
	search    ui.SearchTab
	info      ui.InfoTab
	doctor    ui.DoctorTab

	// Track lazy-loaded data
	allAliases    []data.Alias
	allFuncs      []data.Function
	allKeybindings []data.Keybinding
}

// New creates the initial application model with fast-loading data only.
func New() Model {
	cfg := config.Detect()
	styles := ui.NewStyles(true) // assume dark until detected

	// Fast: file parsing only, no subprocesses
	aliases, _ := data.LoadAliases(cfg.AliasDir)
	funcs, _ := data.LoadFunctions(cfg.FunctionDir)
	keybindings := data.LoadKeybindings()

	return Model{
		cfg:            cfg,
		styles:         styles,
		isDark:         true,
		aliases:        ui.NewAliasTab(aliases, styles),
		functions:      ui.NewFunctionTab(funcs, styles),
		packages:       ui.NewPackageTab(nil, styles), // loaded async
		tmux:           ui.NewTmuxTab(styles),
		search:         ui.NewSearchTab(aliases, funcs, nil, keybindings, styles),
		info:           ui.NewInfoTab(data.SystemInfo{}, cfg.Version, styles), // loaded async
		doctor:         ui.NewDoctorTab(cfg, styles),
		allAliases:     aliases,
		allFuncs:       funcs,
		allKeybindings: keybindings,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		m.loadPackagesAsync(),
		m.loadSysInfoAsync(),
	)
}

func (m Model) loadPackagesAsync() tea.Cmd {
	cfg := m.cfg
	return func() tea.Msg {
		pkgs, _ := data.LoadPackages(cfg.ChezmoiSrc, cfg.OS)
		return packagesLoadedMsg{pkgs: pkgs}
	}
}

func (m Model) loadSysInfoAsync() tea.Cmd {
	return func() tea.Msg {
		info := data.LoadSystemInfo()
		return sysInfoLoadedMsg{info: info}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		m.styles = ui.NewStyles(m.isDark)
		m.aliases.SetStyles(m.styles)
		m.functions.SetStyles(m.styles)
		m.packages.SetStyles(m.styles)
		m.tmux.SetStyles(m.styles)
		m.search.SetStyles(m.styles)
		m.info.SetStyles(m.styles)
		m.doctor.SetStyles(m.styles)
		return m, nil

	case packagesLoadedMsg:
		m.packages = ui.NewPackageTab(msg.pkgs, m.styles)
		m.search = ui.NewSearchTab(m.allAliases, m.allFuncs, msg.pkgs, m.allKeybindings, m.styles)
		if m.width > 0 {
			h, v := m.styles.Doc.GetFrameSize()
			contentW := m.width - h
			contentH := m.height - v - 4
			m.packages.SetSize(contentW, contentH)
			m.search.SetSize(contentW, contentH)
		}
		return m, nil

	case sysInfoLoadedMsg:
		m.info = ui.NewInfoTab(msg.info, m.cfg.Version, m.styles)
		if m.width > 0 {
			h, v := m.styles.Doc.GetFrameSize()
			contentW := m.width - h
			contentH := m.height - v - 4
			m.info.SetSize(contentW, contentH)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Calculate content area
		h, v := m.styles.Doc.GetFrameSize()
		contentW := msg.Width - h
		contentH := msg.Height - v - 4 // tab bar + help bar

		m.aliases.SetSize(contentW, contentH)
		m.functions.SetSize(contentW, contentH)
		m.packages.SetSize(contentW, contentH)
		m.tmux.SetSize(contentW, contentH)
		m.search.SetSize(contentW, contentH)
		m.info.SetSize(contentW, contentH)
		m.doctor.SetSize(contentW, contentH)
		return m, nil

	case tea.KeyPressMsg:
		// Global keys (not intercepted by tabs)
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1":
			m.activeTab = tabAliases
			return m, nil
		case "2":
			m.activeTab = tabFunctions
			return m, nil
		case "3":
			m.activeTab = tabPackages
			return m, nil
		case "4":
			m.activeTab = tabTmux
			return m, nil
		case "5":
			m.activeTab = tabSearch
			return m, nil
		case "6":
			m.activeTab = tabInfo
			return m, nil
		case "7":
			m.activeTab = tabDoctor
			return m, nil
		}

		// q quits only when not filtering in a list
		if msg.String() == "q" && !m.isFiltering() {
			return m, tea.Quit
		}

		// Tab/shift-tab for tab switching (only when not filtering)
		if !m.isFiltering() {
			switch msg.String() {
			case "tab":
				m.activeTab = tabID((int(m.activeTab) + 1) % len(tabNames))
				return m, nil
			case "shift+tab":
				m.activeTab = tabID((int(m.activeTab) - 1 + len(tabNames)) % len(tabNames))
				return m, nil
			}
		}
	}

	// Delegate to active tab
	var cmd tea.Cmd
	switch m.activeTab {
	case tabAliases:
		cmd = m.aliases.Update(msg)
	case tabFunctions:
		cmd = m.functions.Update(msg)
	case tabPackages:
		cmd = m.packages.Update(msg)
	case tabTmux:
		cmd = m.tmux.Update(msg)
	case tabSearch:
		cmd = m.search.Update(msg)
	case tabInfo:
		cmd = m.info.Update(msg)
	case tabDoctor:
		cmd = m.doctor.Update(msg)
	}

	return m, cmd
}

func (m Model) isFiltering() bool {
	return false // Lists handle their own key routing when filtering
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("  Loading shellkit...")
	}

	var doc strings.Builder

	// Tab bar
	var tabs []string
	for i, name := range tabNames {
		if tabID(i) == m.activeTab {
			tabs = append(tabs, m.styles.ActiveTab.Render(name))
		} else {
			tabs = append(tabs, m.styles.InactiveTab.Render(name))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)

	// Fill remaining width with a subtle bottom border
	tabBarWidth := lipgloss.Width(tabBar)
	if gap := m.width - tabBarWidth - 4; gap > 0 {
		tabBar += m.styles.TabBar.Render(strings.Repeat(" ", gap))
	}
	doc.WriteString(tabBar)
	doc.WriteString("\n")

	// Tab content
	var content string
	switch m.activeTab {
	case tabAliases:
		content = m.aliases.View()
	case tabFunctions:
		content = m.functions.View()
	case tabPackages:
		content = m.packages.View()
	case tabTmux:
		content = m.tmux.View()
	case tabSearch:
		content = m.search.View()
	case tabInfo:
		content = m.info.View()
	case tabDoctor:
		content = m.doctor.View()
	}
	doc.WriteString(content)

	// Help bar
	help := m.styles.HelpBar.Render("  tab: switch  1-7: jump to tab  /: filter  q: quit")
	doc.WriteString("\n")
	doc.WriteString(help)

	v := tea.NewView(m.styles.Doc.Render(doc.String()))
	v.AltScreen = true
	return v
}
