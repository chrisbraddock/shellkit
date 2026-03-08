package ui

import (
	"fmt"
	"math"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	configToggleCompactAccent
	configAdjustAnimationSpeed
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
	cfg          config.Config
	settings     config.UISettings
	width        int
	height       int
	cursor       int
	status       string
	styles       *Styles
	previewMode  animMode
	previewFrame int
}

// NewConfigTab creates the config tab.
func NewConfigTab(cfg config.Config, settings config.UISettings, styles *Styles) ConfigTab {
	previewMode := animWaveDots
	if modes := allAnimModes(); len(modes) > 0 {
		previewMode = modes[0]
	}

	return ConfigTab{
		cfg:         cfg,
		settings:    settings,
		styles:      styles,
		previewMode: previewMode,
	}
}

// AtTop returns true when the selection is at the top.
func (t *ConfigTab) AtTop() bool { return t.cursor == 0 }

func (t *ConfigTab) SetStyles(s *Styles) { t.styles = s }

func (t *ConfigTab) SetPreviewFrame(frame int) {
	t.previewFrame = frame
}

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
	case "left", "h":
		t.syncPreview(items[t.cursor])
		return t.adjustCurrent(items[t.cursor], -config.AnimationSpeedStep)
	case "right", "l":
		t.syncPreview(items[t.cursor])
		return t.adjustCurrent(items[t.cursor], config.AnimationSpeedStep)
	case " ", "enter":
		t.syncPreview(items[t.cursor])
		return t.toggleCurrent(items[t.cursor])
	}
	t.syncPreview(items[t.cursor])
	return nil
}

func (t *ConfigTab) View() string {
	var b strings.Builder

	b.WriteString(t.styles.Title.Render("  Header Config"))
	b.WriteString("\n")
	b.WriteString(t.styles.Subtle.Render("  Toggle header elements, startup state, animation speed, and rotation."))
	b.WriteString("\n\n")

	items := t.items()
	if len(items) == 0 {
		b.WriteString(t.styles.Subtle.Render("  No settings available."))
		return b.String()
	}

	if t.useSplitLayout() {
		b.WriteString(t.renderSplitBody(items))
		b.WriteString("\n\n")
		b.WriteString(t.styles.HelpBar.Render("  up/down move · enter toggle/reset · left/right adjust speed · settings saved immediately"))
		if t.status != "" {
			b.WriteString("\n")
			b.WriteString(t.styles.Info.Render("  " + t.status))
		}
		return b.String()
	}

	preview := t.renderPreview()
	if preview != "" {
		b.WriteString(preview)
		b.WriteString("\n\n")
	}

	b.WriteString(t.renderConfigList(items, 0, true))

	b.WriteString("\n")
	b.WriteString(t.styles.HelpBar.Render("  up/down move · enter toggle/reset · left/right adjust speed · settings saved immediately"))
	if t.status != "" {
		b.WriteString("\n")
		b.WriteString(t.styles.Info.Render("  " + t.status))
	}

	return b.String()
}

func (t *ConfigTab) renderSplitBody(items []configTabItem) string {
	const gap = 4

	leftWidth := clampInt(t.width/3, 28, 46)
	rightWidth := maxInt(36, t.width-leftWidth-gap-4)

	left := t.renderPreviewAtWidth(leftWidth)
	right := t.renderConfigList(items, rightWidth, false)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		strings.Repeat(" ", gap),
		right,
	)
}

func (t *ConfigTab) renderConfigList(items []configTabItem, contentWidth int, previewConsumesHeight bool) string {
	var b strings.Builder

	visible := len(items)
	if t.height > 0 {
		reserved := 8
		if previewConsumesHeight {
			preview := t.renderPreview()
			if preview != "" {
				reserved += strings.Count(preview, "\n") + 2
			}
		}
		visible = clampInt(t.height-reserved, 6, len(items))
	}
	start := 0
	if t.cursor >= visible {
		start = t.cursor - visible + 1
	}
	end := minInt(len(items), start+visible)

	section := ""
	for i := start; i < end; {
		item := items[i]
		if item.section != section {
			if i > start {
				b.WriteString("\n")
			}
			b.WriteString(t.styles.SectionTitle.Render("  " + item.section))
			b.WriteString("\n")
			section = item.section
		}

		sectionEnd := i + 1
		for sectionEnd < end && items[sectionEnd].section == item.section {
			sectionEnd++
		}

		if t.useTwoColumnAnimations(contentWidth) && item.kind == configToggleAnimation && sectionEnd-i > 1 {
			b.WriteString(t.renderAnimationColumns(items[i:sectionEnd], i, contentWidth))
			b.WriteString("\n")
		} else {
			for idx := i; idx < sectionEnd; idx++ {
				b.WriteString(t.renderItemBlock(items[idx], idx, contentWidth))
				b.WriteString("\n")
			}
		}

		i = sectionEnd
	}

	return b.String()
}

// Summary returns a short status string for the status bar.
func (t *ConfigTab) Summary() string {
	state := "open"
	if t.settings.Header.StartCollapsed {
		state = "collapsed"
	}
	return fmt.Sprintf("%s · %d fx · %d%%", state, len(t.enabledAnimationIDs()), t.animationSpeed())
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
		{
			section:     "Compact Mode",
			kind:        configToggleCompactAccent,
			label:       "Show animated tab accent",
			description: "Collapsed mode keeps the extra animated line below the navigation tabs.",
		},
		{
			section:     "Animation",
			kind:        configAdjustAnimationSpeed,
			label:       "Header animation speed",
			description: "Control how fast the header, tab chrome, and preview animate.",
		},
	}

	for _, mode := range allAnimModes() {
		items = append(items, configTabItem{
			section:     animationSection(mode),
			kind:        configToggleAnimation,
			mode:        mode,
			label:       modeName(mode),
			description: modeDescription(mode),
		})
	}

	return items
}

func (t *ConfigTab) syncPreview(item configTabItem) {
	if item.kind == configToggleAnimation {
		t.previewMode = item.mode
	}
}

func (t *ConfigTab) renderPreview() string {
	panelWidth := 52
	if t.width > 0 {
		panelWidth = clampInt(t.width-4, 28, 58)
	}
	return t.renderPreviewAtWidth(panelWidth)
}

func (t *ConfigTab) renderPreviewAtWidth(panelWidth int) string {
	innerWidth := clampInt(panelWidth-4, 20, 54)

	body := renderPreviewTitle(t.previewMode, t.styles) + "\n" +
		renderAnimationPreview(t.previewMode, innerWidth, t.previewFrame+11, t.styles)

	return t.styles.Preview.Width(panelWidth).Render(body)
}

func (t *ConfigTab) useSplitLayout() bool {
	return t.width >= 110
}

func (t *ConfigTab) useTwoColumnAnimations(contentWidth int) bool {
	return contentWidth >= 72
}

func (t *ConfigTab) renderAnimationColumns(items []configTabItem, startIndex, contentWidth int) string {
	const gap = 4

	availableWidth := maxInt(60, contentWidth-2)
	colWidth := maxInt(28, (availableWidth-gap)/2)
	mid := (len(items) + 1) / 2

	var b strings.Builder
	for row := 0; row < mid; row++ {
		left := t.renderItemBlock(items[row], startIndex+row, colWidth)
		right := blankConfigItemBlock(colWidth)
		if row+mid < len(items) {
			right = t.renderItemBlock(items[row+mid], startIndex+row+mid, colWidth)
		}

		b.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			left,
			strings.Repeat(" ", gap),
			right,
		))
		if row < mid-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *ConfigTab) renderItemBlock(item configTabItem, index, width int) string {
	if item.kind == configAdjustAnimationSpeed {
		return t.renderSpeedItemBlock(item, index, width)
	}

	cursor := "  "
	if index == t.cursor {
		cursor = "> "
	}
	lineStyle := t.styles.Subtle
	if index == t.cursor {
		lineStyle = t.styles.Highlight
	}

	checked := "[ ]"
	if t.itemEnabled(item) {
		checked = "[x]"
	}

	label := "  " + cursor + checked + " " + item.label
	desc := "      " + truncateConfigText(item.description, maxInt(20, width-6))

	if width <= 0 {
		return lineStyle.Render(label) + "\n" + t.styles.Subtle.Render(desc)
	}

	return lineStyle.Width(width).Render(label) + "\n" +
		t.styles.Subtle.Width(width).Render(desc)
}

func (t *ConfigTab) renderSpeedItemBlock(item configTabItem, index, width int) string {
	cursor := "  "
	if index == t.cursor {
		cursor = "> "
	}
	lineStyle := t.styles.Subtle
	if index == t.cursor {
		lineStyle = t.styles.Highlight
	}

	speed := t.animationSpeed()
	label := fmt.Sprintf("  %s[%3d%%] %s", cursor, speed, item.label)
	slider := fmt.Sprintf("      %s", renderConfigSlider(speed, 16))
	if width > 0 {
		slider = truncateConfigText(slider+"  "+item.description, maxInt(20, width-1))
		return lineStyle.Width(width).Render(label) + "\n" +
			t.styles.Subtle.Width(width).Render(slider)
	}
	return lineStyle.Render(label) + "\n" + t.styles.Subtle.Render(slider)
}

func blankConfigItemBlock(width int) string {
	if width <= 0 {
		return " \n "
	}
	style := lipgloss.NewStyle().Width(width)
	return style.Render(" ") + "\n" + style.Render(" ")
}

func truncateConfigText(text string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

func (t *ConfigTab) itemEnabled(item configTabItem) bool {
	switch item.kind {
	case configToggleCollapsed:
		return t.settings.Header.StartCollapsed
	case configToggleVersion:
		return t.settings.Header.ShowVersion
	case configTogglePlatform:
		return t.settings.Header.ShowPlatform
	case configToggleCompactAccent:
		return t.settings.Header.ShowCompactTabAccent
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

func (t *ConfigTab) animationSpeed() int {
	return config.NormalizeAnimationSpeed(t.settings.Header.AnimationSpeed)
}

func (t *ConfigTab) toggleCurrent(item configTabItem) tea.Cmd {
	switch item.kind {
	case configToggleCollapsed:
		t.settings.Header.StartCollapsed = !t.settings.Header.StartCollapsed
	case configToggleVersion:
		t.settings.Header.ShowVersion = !t.settings.Header.ShowVersion
	case configTogglePlatform:
		t.settings.Header.ShowPlatform = !t.settings.Header.ShowPlatform
	case configToggleCompactAccent:
		t.settings.Header.ShowCompactTabAccent = !t.settings.Header.ShowCompactTabAccent
	case configAdjustAnimationSpeed:
		if t.animationSpeed() == config.DefaultAnimationSpeed {
			t.status = "Animation speed already at default."
			return nil
		}
		t.settings.Header.AnimationSpeed = config.DefaultAnimationSpeed
	case configToggleAnimation:
		if !t.toggleAnimation(item.mode) {
			return nil
		}
	}

	return t.saveSettings()
}

func (t *ConfigTab) adjustCurrent(item configTabItem, delta int) tea.Cmd {
	if item.kind != configAdjustAnimationSpeed || delta == 0 {
		return nil
	}

	next := config.NormalizeAnimationSpeed(t.animationSpeed() + delta)
	if next == t.animationSpeed() {
		return nil
	}
	t.settings.Header.AnimationSpeed = next
	return t.saveSettings()
}

func (t *ConfigTab) saveSettings() tea.Cmd {
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

func renderConfigSlider(speed, width int) string {
	width = maxInt(8, width)
	pos := 0
	if config.MaxAnimationSpeed > config.MinAnimationSpeed {
		pos = int(math.Round(float64(speed-config.MinAnimationSpeed) / float64(config.MaxAnimationSpeed-config.MinAnimationSpeed) * float64(width-1)))
	}
	pos = clampInt(pos, 0, width-1)

	var b strings.Builder
	b.WriteRune('[')
	for i := 0; i < width; i++ {
		switch {
		case i < pos:
			b.WriteRune('=')
		case i == pos:
			b.WriteRune('|')
		default:
			b.WriteRune('-')
		}
	}
	b.WriteRune(']')
	return b.String()
}
