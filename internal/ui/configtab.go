package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/chrisbraddock/shellkit/internal/config"
)

// HeaderSettingsChangedMsg notifies the app that persisted header settings changed.
type HeaderSettingsChangedMsg struct {
	Settings config.UISettings
}

type configItemKind int

const (
	configToggleCollapsed configItemKind = iota
	configToggleVersion
	configTogglePlatform
	configToggleAnimation
)

type configTabItem struct {
	section     string
	kind        configItemKind
	mode        animMode
	label       string
	description string
}

// ConfigTab edits persisted TUI and header settings.
type ConfigTab struct {
	cfg      config.Config
	settings config.UISettings
	width    int
	height   int
	cursor   int
	status   string
	styles   *Styles
}

// NewConfigTab creates the config tab.
func NewConfigTab(cfg config.Config, settings config.UISettings, styles *Styles) ConfigTab {
	return ConfigTab{
		cfg:      cfg,
		settings: settings,
		styles:   styles,
	}
}

// AtTop returns true when the selection is at the top.
func (t *ConfigTab) AtTop() bool { return t.cursor == 0 }

func (t *ConfigTab) SetStyles(s *Styles) { t.styles = s }

func (t *ConfigTab) SetSettings(settings config.UISettings) {
	t.settings = settings
	items := t.items()
	if len(items) == 0 {
		t.cursor = 0
		return
	}
	t.cursor = clampInt(t.cursor, 0, len(items)-1)
}

func (t *ConfigTab) SetSize(w, h int) {
	t.width = w
	t.height = h
}

func (t *ConfigTab) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}

	items := t.items()
	if len(items) == 0 {
		return nil
	}

	switch key.String() {
	case "up", "k":
		t.cursor = clampInt(t.cursor-1, 0, len(items)-1)
	case "down", "j":
		t.cursor = clampInt(t.cursor+1, 0, len(items)-1)
	case " ", "enter":
		return t.toggleCurrent(items[t.cursor])
	}
	return nil
}

func (t *ConfigTab) View() string {
	var b strings.Builder

	b.WriteString(t.styles.Title.Render("  Header Config"))
	b.WriteString("\n")
	b.WriteString(t.styles.Subtle.Render("  Toggle header elements, startup state, and animation rotation."))
	b.WriteString("\n\n")

	items := t.items()
	if len(items) == 0 {
		b.WriteString(t.styles.Subtle.Render("  No settings available."))
		return b.String()
	}

	visible := len(items)
	if t.height > 0 {
		visible = clampInt(t.height-8, 6, len(items))
	}
	start := 0
	if t.cursor >= visible {
		start = t.cursor - visible + 1
	}
	end := minInt(len(items), start+visible)

	section := ""
	for i := start; i < end; i++ {
		item := items[i]
		if item.section != section {
			if i > start {
				b.WriteString("\n")
			}
			b.WriteString(t.styles.SectionTitle.Render("  " + item.section))
			b.WriteString("\n")
			section = item.section
		}

		cursor := "  "
		if i == t.cursor {
			cursor = "> "
		}
		lineStyle := t.styles.Subtle
		if i == t.cursor {
			lineStyle = t.styles.Highlight
		}

		checked := "[ ]"
		if t.itemEnabled(item) {
			checked = "[x]"
		}

		b.WriteString(lineStyle.Render("  " + cursor + checked + " " + item.label))
		b.WriteString("\n")
		b.WriteString(t.styles.Subtle.Render("      " + item.description))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(t.styles.HelpBar.Render("  up/down move · enter toggle · settings saved immediately"))
	if t.status != "" {
		b.WriteString("\n")
		b.WriteString(t.styles.Info.Render("  " + t.status))
	}

	return b.String()
}

// Summary returns a short status string for the status bar.
func (t *ConfigTab) Summary() string {
	state := "open"
	if t.settings.Header.StartCollapsed {
		state = "collapsed"
	}
	return fmt.Sprintf("%s · %d fx", state, len(t.enabledAnimationIDs()))
}

func (t *ConfigTab) items() []configTabItem {
	items := []configTabItem{
		{
			section:     "Startup",
			kind:        configToggleCollapsed,
			label:       "Default to collapsed header",
			description: "New shellkit sessions start in the compact header mode.",
		},
		{
			section:     "Header Elements",
			kind:        configToggleVersion,
			label:       "Show version",
			description: "Show the shellkit version in the expanded and compact header.",
		},
		{
			section:     "Header Elements",
			kind:        configTogglePlatform,
			label:       "Show platform",
			description: "Show the current OS and architecture in the header chrome.",
		},
	}

	for _, mode := range allAnimModes() {
		items = append(items, configTabItem{
			section:     "Animations",
			kind:        configToggleAnimation,
			mode:        mode,
			label:       modeName(mode),
			description: "Include this background in auto-rotate and random startup selection.",
		})
	}

	return items
}

func (t *ConfigTab) itemEnabled(item configTabItem) bool {
	switch item.kind {
	case configToggleCollapsed:
		return t.settings.Header.StartCollapsed
	case configToggleVersion:
		return t.settings.Header.ShowVersion
	case configTogglePlatform:
		return t.settings.Header.ShowPlatform
	case configToggleAnimation:
		return t.animationEnabled(item.mode)
	default:
		return false
	}
}

func (t *ConfigTab) enabledAnimationIDs() []string {
	ids := t.settings.Header.EnabledAnimations
	if len(ids) == 0 {
		ids = make([]string, 0, len(allAnimModes()))
		for _, mode := range allAnimModes() {
			ids = append(ids, animModeID(mode))
		}
		return ids
	}

	enabled := enabledModesFromIDs(ids)
	out := make([]string, 0, len(enabled))
	for _, mode := range enabled {
		out = append(out, animModeID(mode))
	}
	return out
}

func (t *ConfigTab) animationEnabled(mode animMode) bool {
	id := animModeID(mode)
	for _, enabled := range t.enabledAnimationIDs() {
		if enabled == id {
			return true
		}
	}
	return false
}

func (t *ConfigTab) toggleCurrent(item configTabItem) tea.Cmd {
	switch item.kind {
	case configToggleCollapsed:
		t.settings.Header.StartCollapsed = !t.settings.Header.StartCollapsed
	case configToggleVersion:
		t.settings.Header.ShowVersion = !t.settings.Header.ShowVersion
	case configTogglePlatform:
		t.settings.Header.ShowPlatform = !t.settings.Header.ShowPlatform
	case configToggleAnimation:
		if !t.toggleAnimation(item.mode) {
			return nil
		}
	}

	if err := config.SaveUISettings(t.cfg.MetricsDir, t.settings); err != nil {
		t.status = fmt.Sprintf("save failed: %v", err)
		return nil
	}

	t.status = "Saved " + config.UISettingsPath(t.cfg.MetricsDir)
	settings := t.settings
	return func() tea.Msg {
		return HeaderSettingsChangedMsg{Settings: settings}
	}
}

func (t *ConfigTab) toggleAnimation(mode animMode) bool {
	id := animModeID(mode)
	enabled := t.enabledAnimationIDs()
	enabledSet := make(map[string]bool, len(enabled))
	for _, current := range enabled {
		enabledSet[current] = true
	}

	if enabledSet[id] && len(enabledSet) == 1 {
		t.status = "Keep at least one animation enabled."
		return false
	}

	enabledSet[id] = !enabledSet[id]

	var next []string
	for _, candidate := range allAnimModes() {
		candidateID := animModeID(candidate)
		if enabledSet[candidateID] {
			next = append(next, candidateID)
		}
	}

	if len(next) == len(allAnimModes()) {
		t.settings.Header.EnabledAnimations = nil
		return true
	}

	t.settings.Header.EnabledAnimations = next
	return true
}
