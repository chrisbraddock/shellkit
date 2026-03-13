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
	tabConfig
)

// statusBarLines is the number of lines consumed by the status bar.
const statusBarLines = 1

// Messages for async data loading
type packagesLoadedMsg struct{ pkgs []data.Package }
type sysInfoLoadedMsg struct{ info data.SystemInfo }
type metricsLoadedMsg struct{ entries []data.MetricEntry }

// Model is the root Bubble Tea model.
type Model struct {
	cfg        config.Config
	activeTab  tabID
	tabFocused bool // true = tab bar focused, false = content focused
	width      int
	height     int
	isDark     bool
	styles     *ui.Styles
	ready      bool
	settings   config.UISettings

	headerState ui.HeaderState
	aliases     ui.AliasTab
	functions   ui.FunctionTab
	packages    ui.PackageTab
	tmux        ui.TmuxTab
	search      ui.SearchTab
	dashboard   ui.DashboardTab
	doctor      ui.DoctorTab
	configTab   ui.ConfigTab

	// Track lazy-loaded data
	allAliases     []data.Alias
	allFuncs       []data.Function
	allKeybindings []data.Keybinding
}

// New creates the initial application model with fast-loading data only.
func New() Model {
	cfg := config.Detect()
	styles := ui.NewStyles(true) // assume dark until detected
	settings, _ := config.LoadUISettings(cfg.MetricsDir)

	// Fast: file parsing only, no subprocesses
	aliases, _ := data.LoadAliases(cfg.AliasDir)
	funcs, _ := data.LoadFunctions(cfg.FunctionDir)
	keybindings := data.LoadKeybindings()

	return Model{
		cfg:            cfg,
		styles:         styles,
		isDark:         true,
		settings:       settings,
		tabFocused:     true,
		headerState:    ui.NewHeaderState(styles, settings),
		aliases:        ui.NewAliasTab(aliases, styles),
		functions:      ui.NewFunctionTab(funcs, styles),
		packages:       ui.NewPackageTab(nil, styles),
		tmux:           ui.NewTmuxTab(styles),
		search:         ui.NewSearchTab(aliases, funcs, nil, keybindings, styles),
		dashboard:      ui.NewDashboardTab(nil, data.SystemInfo{}, cfg.Version, styles),
		doctor:         ui.NewDoctorTab(cfg, styles),
		configTab:      ui.NewConfigTab(cfg, settings, styles),
		allAliases:     aliases,
		allFuncs:       funcs,
		allKeybindings: keybindings,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		m.headerState.Init(),
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
	contentH := m.height - v - m.chromeLines()
	if contentH < 1 {
		contentH = 1
	}
	return contentW, contentH
}

func (m Model) chromeLines() int {
	return m.headerState.LineCount() + ui.TabBarLineCount(&m.headerState) + statusBarLines
}

func (m *Model) syncContentSize() {
	if m.width <= 0 {
		return
	}

	contentW, contentH := m.contentSize()
	m.aliases.SetSize(contentW, contentH)
	m.functions.SetSize(contentW, contentH)
	m.packages.SetSize(contentW, contentH)
	m.tmux.SetSize(contentW, contentH)
	m.search.SetSize(contentW, contentH)
	m.dashboard.SetSize(contentW, contentH)
	m.doctor.SetSize(contentW, contentH)
	m.configTab.SetSize(contentW, contentH)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ui.AnimTickMsg:
		cmd := m.headerState.Update(msg)
		m.configTab.SetPreviewFrame(m.headerState.Frame())
		return m, cmd

	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		m.styles = ui.NewStyles(m.isDark)
		m.headerState.SetStyles(m.styles)
		m.aliases.SetStyles(m.styles)
		m.functions.SetStyles(m.styles)
		m.packages.SetStyles(m.styles)
		m.tmux.SetStyles(m.styles)
		m.search.SetStyles(m.styles)
		m.dashboard.SetStyles(m.styles)
		m.doctor.SetStyles(m.styles)
		m.configTab.SetStyles(m.styles)
		return m, nil

	case ui.HeaderSettingsChangedMsg:
		m.settings = msg.Settings
		m.headerState.ApplySettings(msg.Settings, false)
		m.configTab.SetSettings(msg.Settings)
		m.syncContentSize()
		return m, nil

	case packagesLoadedMsg:
		m.packages = ui.NewPackageTab(msg.pkgs, m.styles)
		m.search = ui.NewSearchTab(m.allAliases, m.allFuncs, msg.pkgs, m.allKeybindings, m.styles)
		m.syncContentSize()
		return m, nil

	case sysInfoLoadedMsg:
		m.dashboard.SetSysInfo(msg.info)
		m.syncContentSize()
		return m, nil

	case metricsLoadedMsg:
		m.dashboard.SetMetrics(msg.entries)
		m.syncContentSize()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.headerState.SetSize(m.width)
		m.syncContentSize()
		return m, nil

	case tea.KeyPressMsg:
		key := msg.String()

		// Global keys always work
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "shift+up":
			m.headerState.Collapse()
			m.syncContentSize()
			return m, nil
		case "shift+down":
			m.headerState.Expand()
			m.syncContentSize()
			return m, nil
		case "[":
			m.headerState.CycleMode(-1)
			return m, nil
		case "]":
			m.headerState.CycleMode(1)
			return m, nil
		case "m":
			m.headerState.ToggleLock()
			return m, nil
		case "1":
			m.activeTab = tabDashboard
			m.tabFocused = true
			return m, nil
		case "2":
			m.activeTab = tabAliases
			m.tabFocused = true
			return m, nil
		case "3":
			m.activeTab = tabFunctions
			m.tabFocused = true
			return m, nil
		case "4":
			m.activeTab = tabPackages
			m.tabFocused = true
			return m, nil
		case "5":
			m.activeTab = tabTmux
			m.tabFocused = true
			return m, nil
		case "6":
			m.activeTab = tabSearch
			m.tabFocused = true
			return m, nil
		case "7":
			m.activeTab = tabDoctor
			m.tabFocused = true
			return m, nil
		case "8":
			m.activeTab = tabConfig
			m.tabFocused = true
			return m, nil
		}

		if key == "q" && !m.isFiltering() {
			return m, tea.Quit
		}

		if m.tabFocused {
			// Tab bar is focused: left/right change tabs, down enters content
			switch key {
			case "left", "shift+tab":
				m.activeTab = tabID((int(m.activeTab) - 1 + len(ui.TabNames)) % len(ui.TabNames))
				return m, nil
			case "right", "tab":
				m.activeTab = tabID((int(m.activeTab) + 1) % len(ui.TabNames))
				return m, nil
			case "down", "enter":
				m.tabFocused = false
				return m, nil
			}
			return m, nil
		}

		// Content is focused
		if !m.isFiltering() {
			switch key {
			case "tab":
				m.activeTab = tabID((int(m.activeTab) + 1) % len(ui.TabNames))
				m.tabFocused = true
				return m, nil
			case "shift+tab":
				m.activeTab = tabID((int(m.activeTab) - 1 + len(ui.TabNames)) % len(ui.TabNames))
				m.tabFocused = true
				return m, nil
			}

			// Up at the top of content → return to tab bar
			if key == "up" && m.activeTabAtTop() {
				m.tabFocused = true
				return m, nil
			}
		}
	}

	// Dispatch to active tab content (only when content focused)
	if !m.tabFocused {
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
		case tabConfig:
			cmd = m.configTab.Update(msg)
		}
		return m, cmd
	}

	return m, nil
}

func (m Model) isFiltering() bool {
	return false
}

// activeTabAtTop returns true if the active tab's content is scrolled/selected at the top.
func (m Model) activeTabAtTop() bool {
	switch m.activeTab {
	case tabAliases:
		return m.aliases.AtTop()
	case tabFunctions:
		return m.functions.AtTop()
	case tabPackages:
		return m.packages.AtTop()
	case tabTmux:
		return m.tmux.AtTop()
	case tabSearch:
		return m.search.AtTop()
	case tabDashboard:
		return m.dashboard.AtTop()
	case tabDoctor:
		return m.doctor.AtTop()
	case tabConfig:
		return m.configTab.AtTop()
	}
	return true
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
	s["Config"] = m.configTab.Summary()
	return s
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("  Loading shellkit...")
	}

	var doc strings.Builder

	header := ui.RenderHeader(m.cfg.Version, m.cfg.OS, m.cfg.Arch, m.width, m.styles, &m.headerState)
	doc.WriteString(header)

	tabBar := ui.RenderTabBar(int(m.activeTab), m.tabFocused, m.width, m.styles, &m.headerState)
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
	case tabConfig:
		content = m.configTab.View()
	}
	doc.WriteString(content)

	doc.WriteString("\n")
	statusBar := ui.RenderStatusBar(int(m.activeTab), m.stats(), m.headerState.ModeStatus(), m.width, m.styles)
	doc.WriteString(statusBar)

	v := tea.NewView(m.styles.Doc.Render(doc.String()))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
