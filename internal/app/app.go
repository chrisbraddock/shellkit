package app

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/chrisbraddock/shellkit/internal/config"
	"github.com/chrisbraddock/shellkit/internal/data"
	"github.com/chrisbraddock/shellkit/internal/ui"
)

type tabID int

const (
	tabDashboard tabID = iota
	tabAliases
	tabFunctions
	tabPackages
	tabTmux
	tabSearch
	tabDoctor
)

// headerLines is the number of lines consumed by the header (logo + blank).
const headerLines = 5

// tabBarLines is the number of lines consumed by the tab bar + gradient line.
const tabBarLines = 2

// statusBarLines is the number of lines consumed by the status bar.
const statusBarLines = 1

// chromeLines is total lines consumed by header + tab bar + status bar.
const chromeLines = headerLines + tabBarLines + statusBarLines

// Messages for async data loading
type packagesLoadedMsg struct{ pkgs []data.Package }
type sysInfoLoadedMsg struct{ info data.SystemInfo }
type metricsLoadedMsg struct{ entries []data.MetricEntry }

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
	dashboard ui.DashboardTab
	doctor    ui.DoctorTab

	// Track lazy-loaded data
	allAliases     []data.Alias
	allFuncs       []data.Function
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
		packages:       ui.NewPackageTab(nil, styles),
		tmux:           ui.NewTmuxTab(styles),
		search:         ui.NewSearchTab(aliases, funcs, nil, keybindings, styles),
		dashboard:      ui.NewDashboardTab(nil, data.SystemInfo{}, cfg.Version, styles),
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
		m.loadMetricsAsync(),
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

func (m Model) loadMetricsAsync() tea.Cmd {
	cfg := m.cfg
	return func() tea.Msg {
		entries, _ := data.LoadMetrics(cfg.HomeDir)
		return metricsLoadedMsg{entries: entries}
	}
}

func (m Model) contentSize() (int, int) {
	h, v := m.styles.Doc.GetFrameSize()
	contentW := m.width - h
	contentH := m.height - v - chromeLines
	if contentH < 1 {
		contentH = 1
	}
	return contentW, contentH
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
		m.dashboard.SetStyles(m.styles)
		m.doctor.SetStyles(m.styles)
		return m, nil

	case packagesLoadedMsg:
		m.packages = ui.NewPackageTab(msg.pkgs, m.styles)
		m.search = ui.NewSearchTab(m.allAliases, m.allFuncs, msg.pkgs, m.allKeybindings, m.styles)
		if m.width > 0 {
			contentW, contentH := m.contentSize()
			m.packages.SetSize(contentW, contentH)
			m.search.SetSize(contentW, contentH)
		}
		return m, nil

	case sysInfoLoadedMsg:
		m.dashboard.SetSysInfo(msg.info)
		if m.width > 0 {
			contentW, contentH := m.contentSize()
			m.dashboard.SetSize(contentW, contentH)
		}
		return m, nil

	case metricsLoadedMsg:
		m.dashboard.SetMetrics(msg.entries)
		if m.width > 0 {
			contentW, contentH := m.contentSize()
			m.dashboard.SetSize(contentW, contentH)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		contentW, contentH := m.contentSize()

		m.aliases.SetSize(contentW, contentH)
		m.functions.SetSize(contentW, contentH)
		m.packages.SetSize(contentW, contentH)
		m.tmux.SetSize(contentW, contentH)
		m.search.SetSize(contentW, contentH)
		m.dashboard.SetSize(contentW, contentH)
		m.doctor.SetSize(contentW, contentH)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1":
			m.activeTab = tabDashboard
			return m, nil
		case "2":
			m.activeTab = tabAliases
			return m, nil
		case "3":
			m.activeTab = tabFunctions
			return m, nil
		case "4":
			m.activeTab = tabPackages
			return m, nil
		case "5":
			m.activeTab = tabTmux
			return m, nil
		case "6":
			m.activeTab = tabSearch
			return m, nil
		case "7":
			m.activeTab = tabDoctor
			return m, nil
		}

		if msg.String() == "q" && !m.isFiltering() {
			return m, tea.Quit
		}

		if !m.isFiltering() {
			switch msg.String() {
			case "tab":
				m.activeTab = tabID((int(m.activeTab) + 1) % len(ui.TabNames))
				return m, nil
			case "shift+tab":
				m.activeTab = tabID((int(m.activeTab) - 1 + len(ui.TabNames)) % len(ui.TabNames))
				return m, nil
			}
		}
	}

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
	case tabDashboard:
		cmd = m.dashboard.Update(msg)
	case tabDoctor:
		cmd = m.doctor.Update(msg)
	}

	return m, cmd
}

func (m Model) isFiltering() bool {
	return false
}

func (m Model) stats() map[string]string {
	s := make(map[string]string)
	s["Aliases"] = fmt.Sprintf("%d aliases", len(m.allAliases))
	s["Functions"] = fmt.Sprintf("%d functions", len(m.allFuncs))
	s["Packages"] = fmt.Sprintf("%d packages", m.packages.Count())
	s["Tmux"] = ""
	total := len(m.allAliases) + len(m.allFuncs) + m.packages.Count() + len(m.allKeybindings)
	s["Search"] = fmt.Sprintf("%d items", total)
	s["Dashboard"] = m.dashboard.Summary()
	s["Doctor"] = m.doctor.Summary()
	return s
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("  Loading shellkit...")
	}

	var doc strings.Builder

	header := ui.RenderHeader(m.cfg.Version, m.cfg.OS, m.cfg.Arch, m.width, m.styles)
	doc.WriteString(header)

	tabBar := ui.RenderTabBar(int(m.activeTab), m.width, m.styles)
	doc.WriteString(tabBar)

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
	case tabDashboard:
		content = m.dashboard.View()
	case tabDoctor:
		content = m.doctor.View()
	}
	doc.WriteString(content)

	doc.WriteString("\n")
	statusBar := ui.RenderStatusBar(int(m.activeTab), m.stats(), m.width, m.styles)
	doc.WriteString(statusBar)

	v := tea.NewView(m.styles.Doc.Render(doc.String()))
	v.AltScreen = true
	return v
}
