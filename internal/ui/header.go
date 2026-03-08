package ui

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/chrisbraddock/shellkit/internal/config"
)

var asciiLogo = []string{
	` ╔═╗ ╦ ╦ ╔═╗ ╦   ╦   ╦╔═ ╦ ╔╦╗`,
	` ╚═╗ ╠═╣ ║╣  ║   ║   ╠╩╗ ║  ║ `,
	` ╚═╝ ╩ ╩ ╚═╝ ╩═╝ ╩═╝ ╩ ╩ ╩  ╩ `,
}

const logoRows = 3
const headerSceneRows = 7
const compactHeaderRows = 1

// HeaderLineCount is the expanded rendered header height including the badge row.
const HeaderLineCount = headerSceneRows + 1

// Animation modes
type animMode int

const (
	animWaveDots       animMode = iota // sine wave dot field
	animParticles                      // parallax particle drift
	animSineLine                       // continuous sine wave line
	animGradientRain                   // gradient block rain/parallax
	animConstellations                 // twinkling star field
	animWormhole                       // psychedelic tunnel / portal
	animLaserShow                      // neon laser grid sweep
	animFireworks                      // exploding skybursts
	animPlasma                         // demoscene plasma field
	animHyperspace                     // warp-speed star streaks
	animTimeRift                       // VHS glitch / spacetime tear
	animReactorPulse                   // orbital reactor core pulse
	animAlienRun                       // alien side-scroller
	animModeCount                      // sentinel for rotation
)

// Frames per animation before rotating (~8 seconds at 15fps)
const framesPerAnim = 120

// AnimTickMsg signals an animation frame advance.
type AnimTickMsg struct{}

func animTick() tea.Cmd {
	return tea.Tick(66*time.Millisecond, func(time.Time) tea.Msg {
		return AnimTickMsg{}
	})
}

// HeaderState holds animation state for the header.
type HeaderState struct {
	frame        int
	width        int
	styles       *Styles
	locked       bool
	lockedMode   animMode
	compact      bool
	showVersion  bool
	showPlatform bool
	enabledModes []animMode
}

func NewHeaderState(styles *Styles, settings config.UISettings) HeaderState {
	h := HeaderState{
		styles:       styles,
		showVersion:  true,
		showPlatform: true,
	}
	h.ApplySettings(settings, true)
	h.frame = randomHeaderFrame(len(h.availableModes()))
	return h
}

func (h *HeaderState) Init() tea.Cmd {
	return animTick()
}

func (h *HeaderState) SetSize(w int) {
	h.width = w
}

func (h *HeaderState) SetStyles(s *Styles) {
	h.styles = s
}

func (h *HeaderState) Update(msg tea.Msg) tea.Cmd {
	if _, ok := msg.(AnimTickMsg); ok {
		h.frame++
		return animTick()
	}
	return nil
}

func (h *HeaderState) autoMode() animMode {
	modes := h.availableModes()
	return modes[(h.frame/framesPerAnim)%len(modes)]
}

func (h *HeaderState) currentMode() animMode {
	if h.locked {
		if h.modeEnabled(h.lockedMode) {
			return h.lockedMode
		}
		h.locked = false
	}
	return h.autoMode()
}

// localFrame returns the frame number within the current animation cycle.
func (h *HeaderState) localFrame() int {
	return h.frame % framesPerAnim
}

func (h *HeaderState) ToggleLock() {
	if h.locked {
		h.locked = false
		return
	}
	h.locked = true
	h.lockedMode = h.currentMode()
}

func (h *HeaderState) CycleMode(delta int) {
	modes := h.availableModes()
	if len(modes) == 0 {
		return
	}

	base := h.modeIndex(h.currentMode())
	base = (base + delta) % len(modes)
	if base < 0 {
		base += len(modes)
	}
	h.locked = true
	h.lockedMode = modes[base]
}

func (h *HeaderState) ModeName() string {
	return modeName(h.currentMode())
}

func (h *HeaderState) ModeStatus() string {
	state := "auto"
	if h.locked {
		state = "lock"
	}
	return "FX " + h.ModeName() + " [" + state + "]"
}

func (h *HeaderState) Collapse() {
	h.compact = true
}

func (h *HeaderState) Expand() {
	h.compact = false
}

func (h *HeaderState) IsCompact() bool {
	return h.compact
}

func (h *HeaderState) LineCount() int {
	if h.compact {
		return compactHeaderRows
	}
	return HeaderLineCount
}

func (h *HeaderState) ApplySettings(settings config.UISettings, syncCompact bool) {
	h.showVersion = settings.Header.ShowVersion
	h.showPlatform = settings.Header.ShowPlatform
	h.enabledModes = enabledModesFromIDs(settings.Header.EnabledAnimations)
	if syncCompact {
		h.compact = settings.Header.StartCollapsed
	}
	if h.locked && !h.modeEnabled(h.lockedMode) {
		h.locked = false
	}
}

func (h *HeaderState) availableModes() []animMode {
	if len(h.enabledModes) == 0 {
		return allAnimModes()
	}
	return h.enabledModes
}

func (h *HeaderState) modeEnabled(mode animMode) bool {
	for _, candidate := range h.availableModes() {
		if candidate == mode {
			return true
		}
	}
	return false
}

func (h *HeaderState) modeIndex(mode animMode) int {
	modes := h.availableModes()
	for i, candidate := range modes {
		if candidate == mode {
			return i
		}
	}
	return 0
}

func randomHeaderFrame(modeCount int) int {
	if modeCount <= 1 {
		return 0
	}

	mode := rand.IntN(modeCount)
	local := 0
	if framesPerAnim > 30 {
		local = 15 + rand.IntN(framesPerAnim-30)
	} else if framesPerAnim > 0 {
		local = rand.IntN(framesPerAnim)
	}
	return mode*framesPerAnim + local
}

// RenderHeader renders the gradient ASCII logo with animated effects and badges.
func RenderHeader(version, os, arch string, width int, styles *Styles, state *HeaderState) string {
	if state != nil {
		if !state.showVersion {
			version = ""
		}
		if !state.showPlatform {
			os, arch = "", ""
		}
	}

	if state != nil && state.compact {
		return renderCompactHeader(version, os, arch, width, styles, state)
	}

	var doc strings.Builder

	logoWidth := asciiLogoWidth()

	// Generate gradient for the logo
	logoColors := lipgloss.Blend1D(logoWidth, styles.GradientStart, styles.GradientEnd)

	// The scene spans the full content width, with the logo composited on top.
	sceneWidth := width - 4
	if sceneWidth < 0 {
		sceneWidth = 0
	}

	// Generate gradient colors for the full scene.
	var sceneColors []color.Color
	if sceneWidth > 0 {
		sceneColors = lipgloss.Blend1D(sceneWidth, styles.GradientStart, styles.GradientEnd)
	}

	frame := 0
	mode := animWaveDots
	localF := 0
	if state != nil {
		frame = state.frame
		mode = state.currentMode()
		localF = state.localFrame()
	}

	logoTop := clampInt((headerSceneRows-logoRows)/2-1, 0, headerSceneRows-logoRows)
	versionLabel := formatVersionLabel(version)
	versionRunes := []rune(versionLabel)
	versionRow := clampInt(logoTop+logoRows, 0, headerSceneRows-1)
	versionLeft := 0
	var versionColors []color.Color
	showVersionLine := false
	if sceneWidth > 0 && len(versionRunes) > 0 && len(versionRunes) <= sceneWidth {
		versionLeft = clampInt((logoWidth-len(versionRunes))/2, 0, maxInt(sceneWidth-len(versionRunes), 0))
		versionColors = lipgloss.Blend1D(len(versionRunes), styles.SubtleColor, styles.AccentColor)
		showVersionLine = true
	}

	grid := buildAnimationGrid(headerSceneRows, sceneWidth, frame, mode, localF, true)

	// Render each scene row with the logo composited over the animated field.
	for row := 0; row < headerSceneRows; row++ {
		var lineBuilder strings.Builder

		for col := 0; col < sceneWidth; col++ {
			cell := grid[row][col]

			if row >= logoTop && row < logoTop+logoRows && col < logoWidth {
				r := logoRuneAt(row-logoTop, col)
				if r != ' ' {
					lineBuilder.WriteString(renderLogoGlyph(
						r,
						colorAt(logoColors, col, demoWhite),
						cell,
						colorAt(sceneColors, col, styles.AccentColor),
					))
					continue
				}
			}

			if showVersionLine && row == versionRow && col >= versionLeft && col < versionLeft+len(versionRunes) {
				r := versionRunes[col-versionLeft]
				if r == ' ' {
					lineBuilder.WriteRune(' ')
				} else {
					lineBuilder.WriteString(renderVersionGlyph(
						r,
						colorAt(versionColors, col-versionLeft, styles.SubtleColor),
						colorAt(sceneColors, col, styles.AccentColor),
					))
				}
				continue
			}

			if cell.ch == "" || cell.ch == " " {
				lineBuilder.WriteRune(' ')
				continue
			}

			base := colorAt(sceneColors, col, styles.AccentColor)
			if cell.customColor != nil {
				base = cell.customColor
			}
			lineBuilder.WriteString(renderAnimCell(cell, base))
		}

		doc.WriteString(lineBuilder.String())
		doc.WriteString("\n")
	}

	// Version + platform badge on the right
	if (!showVersionLine && version != "") || os != "" {
		var badge strings.Builder
		if !showVersionLine && version != "" {
			badge.WriteString(styles.VersionBadge.Render(version))
		}
		if os != "" && arch != "" {
			if badge.Len() > 0 {
				badge.WriteString(" ")
			}
			badge.WriteString(styles.PlatformBadge.Render(os + "/" + arch))
		}

		badgeStr := badge.String()
		gap := width - lipgloss.Width(badgeStr) - 4
		if gap < 2 {
			gap = 2
		}
		doc.WriteString(strings.Repeat(" ", gap))
		doc.WriteString(badgeStr)
		doc.WriteString("\n")
	}

	return doc.String()
}

func renderCompactHeader(version, os, arch string, width int, styles *Styles, state *HeaderState) string {
	sceneWidth := width - 4
	if sceneWidth <= 0 {
		return "\n"
	}

	frame := 0
	mode := animWaveDots
	localF := 0
	if state != nil {
		frame = state.frame + 9
		mode = state.currentMode()
		localF = state.localFrame()
	}

	grid := buildAnimationGrid(compactHeaderRows, sceneWidth, frame, mode, localF, true)
	seedChromeBleed(grid[0], frame+23)

	logoText := "shellkit"
	pos := overlayPaletteText(grid[0], 1, logoText, []color.Color{
		demoPink,
		demoViolet,
		demoBlue,
		demoCyan,
		demoWhite,
	}, animCell{
		bold: true,
	})

	if version != "" {
		versionText := version
		if !strings.HasPrefix(versionText, "v") {
			versionText = "v" + versionText
		}
		pos = overlayText(grid[0], pos, " // ", animCell{
			customColor: demoViolet,
			dim:         true,
		})
		overlayText(grid[0], pos, versionText+" ", animCell{
			customColor: demoCyan,
		})
	} else {
		overlayText(grid[0], pos, " ", animCell{})
	}

	if os != "" && arch != "" {
		rightLabel := " " + os + "/" + arch + " "
		rightWidth := len([]rune(rightLabel))
		if rightWidth < sceneWidth {
			overlayText(grid[0], sceneWidth-rightWidth, rightLabel, animCell{
				customColor: styles.SubtleColor,
				dim:         true,
			})
		}
	}

	baseColors := lipgloss.Blend1D(sceneWidth, styles.GradientStart, styles.GradientEnd)
	return renderTabRow(grid[0], baseColors) + "\n"
}

// animCell represents one character in the animation grid.
type animCell struct {
	ch          string
	bold        bool
	dim         bool
	underline   bool
	customColor color.Color // nil = use gradient
}

var (
	demoWhite  color.Color = lipgloss.Color("#FFF7FB")
	demoPink   color.Color = lipgloss.Color("#FF4FD8")
	demoViolet color.Color = lipgloss.Color("#9A6BFF")
	demoBlue   color.Color = lipgloss.Color("#5B8CFF")
	demoCyan   color.Color = lipgloss.Color("#35F3FF")
	demoLime   color.Color = lipgloss.Color("#C8FF4D")
	demoAmber  color.Color = lipgloss.Color("#FFB347")
)

// fadeGrid reduces visibility of the grid by the given factor [0, 1].
func fadeGrid(grid [][]animCell, w int, factor float64) {
	for r := range grid {
		for c := 0; c < w; c++ {
			if grid[r][c].ch == "" {
				continue
			}
			if factor < 0.3 {
				grid[r][c].ch = ""
			} else if factor < 0.6 {
				grid[r][c].dim = true
				grid[r][c].bold = false
			}
		}
	}
}

func setCell(grid [][]animCell, row, col int, cell animCell) {
	if row < 0 || row >= len(grid) {
		return
	}
	if col < 0 || col >= len(grid[row]) {
		return
	}
	grid[row][col] = cell
}

func setCellIfEmpty(grid [][]animCell, row, col int, cell animCell) {
	if row < 0 || row >= len(grid) {
		return
	}
	if col < 0 || col >= len(grid[row]) {
		return
	}
	if grid[row][col].ch == "" {
		grid[row][col] = cell
	}
}

func colorFromPalette(palette []color.Color, phase float64) color.Color {
	if len(palette) == 0 {
		return nil
	}
	idx := int(fract(phase) * float64(len(palette)))
	idx = clampInt(idx, 0, len(palette)-1)
	return palette[idx]
}

func fract(v float64) float64 {
	return v - math.Floor(v)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func asciiLogoWidth() int {
	width := 0
	for _, line := range asciiLogo {
		if w := len([]rune(line)); w > width {
			width = w
		}
	}
	return width
}

func formatVersionLabel(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return ""
	}
	if strings.HasPrefix(version, "v") {
		return "shellkit " + version
	}
	return "shellkit v" + version
}

func allAnimModes() []animMode {
	return []animMode{
		animWaveDots,
		animParticles,
		animSineLine,
		animGradientRain,
		animConstellations,
		animWormhole,
		animLaserShow,
		animFireworks,
		animPlasma,
		animHyperspace,
		animTimeRift,
		animReactorPulse,
		animAlienRun,
	}
}

func animModeID(mode animMode) string {
	switch mode {
	case animWaveDots:
		return "wave-dots"
	case animParticles:
		return "particles"
	case animSineLine:
		return "sine-line"
	case animGradientRain:
		return "gradient-rain"
	case animConstellations:
		return "constellations"
	case animWormhole:
		return "wormhole"
	case animLaserShow:
		return "laser-show"
	case animFireworks:
		return "fireworks"
	case animPlasma:
		return "plasma"
	case animHyperspace:
		return "hyperspace"
	case animTimeRift:
		return "time-rift"
	case animReactorPulse:
		return "reactor-pulse"
	case animAlienRun:
		return "alien-run"
	default:
		return ""
	}
}

func animModeFromID(id string) (animMode, bool) {
	for _, mode := range allAnimModes() {
		if animModeID(mode) == id {
			return mode, true
		}
	}
	return 0, false
}

func enabledModesFromIDs(ids []string) []animMode {
	if len(ids) == 0 {
		return allAnimModes()
	}

	var modes []animMode
	seen := make(map[animMode]struct{})
	for _, id := range ids {
		mode, ok := animModeFromID(id)
		if !ok {
			continue
		}
		if _, ok := seen[mode]; ok {
			continue
		}
		seen[mode] = struct{}{}
		modes = append(modes, mode)
	}
	if len(modes) == 0 {
		return allAnimModes()
	}
	return modes
}

func modeName(mode animMode) string {
	switch mode {
	case animWaveDots:
		return "Wave Dots"
	case animParticles:
		return "Particles"
	case animSineLine:
		return "Sine Line"
	case animGradientRain:
		return "Gradient Rain"
	case animConstellations:
		return "Constellations"
	case animWormhole:
		return "Wormhole"
	case animLaserShow:
		return "Laser Show"
	case animFireworks:
		return "Fireworks"
	case animPlasma:
		return "Plasma"
	case animHyperspace:
		return "Hyperspace"
	case animTimeRift:
		return "Time Rift"
	case animReactorPulse:
		return "Reactor Pulse"
	case animAlienRun:
		return "Alien Run"
	default:
		return "Unknown"
	}
}

func buildAnimationGrid(rows, w, frame int, mode animMode, localF int, applyFade bool) [][]animCell {
	grid := make([][]animCell, rows)
	for r := range grid {
		grid[r] = make([]animCell, w)
	}

	if w <= 0 || rows <= 0 {
		return grid
	}

	switch mode {
	case animWaveDots:
		renderWaveDots(grid, w, frame)
	case animParticles:
		renderParticles(grid, w, frame)
	case animSineLine:
		renderSineLine(grid, w, frame)
	case animGradientRain:
		renderGradientRain(grid, w, frame)
	case animConstellations:
		renderConstellations(grid, w, frame, localF)
	case animWormhole:
		renderWormhole(grid, w, frame)
	case animLaserShow:
		renderLaserShow(grid, w, frame)
	case animFireworks:
		renderFireworks(grid, w, frame)
	case animPlasma:
		renderPlasma(grid, w, frame)
	case animHyperspace:
		renderHyperspace(grid, w, frame)
	case animTimeRift:
		renderTimeRift(grid, w, frame)
	case animReactorPulse:
		renderReactorPulse(grid, w, frame)
	case animAlienRun:
		renderAlienRun(grid, w, frame)
	}

	if applyFade {
		if localF < 15 {
			fadeGrid(grid, w, float64(localF)/15.0)
		}
		if localF > framesPerAnim-15 {
			fadeGrid(grid, w, float64(framesPerAnim-localF)/15.0)
		}
	}

	return grid
}

func renderAnimCell(cell animCell, base color.Color) string {
	style := lipgloss.NewStyle().Foreground(base)
	if cell.bold {
		style = style.Bold(true)
	}
	if cell.dim {
		style = style.Faint(true)
	}
	if cell.underline {
		style = style.Underline(true)
	}
	return style.Render(cell.ch)
}

func logoRuneAt(row, col int) rune {
	if row < 0 || row >= len(asciiLogo) || col < 0 {
		return ' '
	}
	runes := []rune(asciiLogo[row])
	if col >= len(runes) {
		return ' '
	}
	return runes[col]
}

func renderLogoGlyph(r rune, baseFg color.Color, cell animCell, fallback color.Color) string {
	glow := fallback
	if cell.customColor != nil {
		glow = cell.customColor
	}

	strength := 0.0
	switch {
	case cell.bold:
		strength = 0.14
	case cell.dim:
		strength = 0.04
	case cell.ch != "" && cell.ch != " ":
		strength = 0.08
	}

	style := lipgloss.NewStyle().
		Foreground(mixColors(baseFg, glow, 0.06+strength)).
		Bold(true)
	return style.Render(string(r))
}

func renderVersionGlyph(r rune, baseFg, glow color.Color) string {
	return lipgloss.NewStyle().
		Foreground(mixColors(baseFg, glow, 0.18)).
		Faint(true).
		Render(string(r))
}

func colorAt(colors []color.Color, idx int, fallback color.Color) color.Color {
	if len(colors) == 0 {
		return fallback
	}
	return colors[clampInt(idx, 0, len(colors)-1)]
}

func mixColors(a, b color.Color, t float64) color.Color {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	r := uint8((1.0-t)*float64(ar>>8) + t*float64(br>>8))
	g := uint8((1.0-t)*float64(ag>>8) + t*float64(bg>>8))
	bl := uint8((1.0-t)*float64(ab>>8) + t*float64(bb>>8))
	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", r, g, bl))
}

// ──────────────────────────────────────────────────
// Animation 1: Wave Dots (original)
// ──────────────────────────────────────────────────

var dotChars = []string{"·", "•", "⦁"}

func renderWaveDots(grid [][]animCell, w, frame int) {
	rows := len(grid)
	for col := 0; col < w; col++ {
		x := float64(col)
		f := float64(frame)

		y := math.Sin(x*0.15+f*0.08) + 0.5*math.Sin(x*0.08-f*0.12)
		rowFloat := (y + 1.5) / 3.0 * float64(rows-1)
		targetRow := clampInt(int(math.Round(rowFloat)), 0, rows-1)

		absY := math.Abs(y)
		dotIdx := 0
		if absY > 0.8 {
			dotIdx = 2
		} else if absY > 0.3 {
			dotIdx = 1
		}

		// Fade right edge
		if col >= w-6 {
			dotIdx = 0
		}

		grid[targetRow][col] = animCell{ch: dotChars[dotIdx], bold: dotIdx == 2}
	}
}

// ──────────────────────────────────────────────────
// Animation 2: Particle Field (parallax layers)
// ──────────────────────────────────────────────────

type particleLayer struct {
	count int
	speed float64
	chars []string
	dim   bool
	bold  bool
}

var particleLayers = []particleLayer{
	{count: 30, speed: 0.3, chars: []string{"·", "⋅", "∘"}, dim: true},              // far (slow, dim)
	{count: 20, speed: 0.7, chars: []string{"·", "•", "◦"}, dim: false},             // mid
	{count: 12, speed: 1.4, chars: []string{"•", "⦁", "✦"}, dim: false, bold: true}, // near (fast, bold)
}

func renderParticles(grid [][]animCell, w, frame int) {
	rows := len(grid)
	for _, layer := range particleLayers {
		for i := 0; i < layer.count; i++ {
			// Deterministic position from particle index
			baseX := pseudoRand(i*7919+3) % w
			row := pseudoRand(i*6271+7) % rows
			charIdx := pseudoRand(i*4177+13) % len(layer.chars)

			// Move particle based on speed (rightward drift, wrapping)
			x := (baseX + int(float64(frame)*layer.speed)) % w

			// Add slight vertical wobble
			wobble := math.Sin(float64(frame)*0.05 + float64(i)*1.7)
			if math.Abs(wobble) > 0.8 {
				row = (row + 1) % rows
			}

			if grid[row][x].ch == "" {
				grid[row][x] = animCell{
					ch:   layer.chars[charIdx],
					bold: layer.bold,
					dim:  layer.dim,
				}
			}
		}
	}
}

// ──────────────────────────────────────────────────
// Animation 3: Sine Wave Line
// ──────────────────────────────────────────────────

// Wave characters for smooth line rendering
var waveLineChars = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▂", "▁"}

func renderSineLine(grid [][]animCell, w, frame int) {
	rows := len(grid)
	f := float64(frame)

	for col := 0; col < w; col++ {
		x := float64(col)

		// Primary wave + harmonic
		y := math.Sin(x*0.12+f*0.06) + 0.3*math.Sin(x*0.25-f*0.1) + 0.2*math.Cos(x*0.07+f*0.04)

		// Map y [-1.5, 1.5] to row [0, 2]
		rowFloat := (y + 1.5) / 3.0 * float64(rows-1)
		targetRow := clampInt(int(math.Round(rowFloat)), 0, rows-1)

		// Map y to wave character (vertical position within cell)
		subIdx := int((rowFloat - float64(targetRow) + 0.5) * float64(len(waveLineChars)-1))
		subIdx = clampInt(subIdx, 0, len(waveLineChars)-1)

		// Trailing particles behind the wave
		for dr := -1; dr <= 1; dr++ {
			r := targetRow + dr
			if r < 0 || r >= rows {
				continue
			}
			if r == targetRow {
				grid[r][col] = animCell{ch: waveLineChars[subIdx], bold: true}
			} else if math.Abs(y) > 0.5 && col%3 == 0 {
				grid[r][col] = animCell{ch: "·", dim: true}
			}
		}
	}
}

// ──────────────────────────────────────────────────
// Animation 4: Gradient Rain (parallax block scrolling)
// ──────────────────────────────────────────────────

var blockChars = []string{" ", "░", "▒", "▓", "█", "▓", "▒", "░"}

func renderGradientRain(grid [][]animCell, w, frame int) {
	rows := len(grid)
	// Three layers of scrolling gradient blocks at different speeds
	type rainLayer struct {
		speed    float64
		period   float64
		rowShift int
		dim      bool
	}
	layers := []rainLayer{
		{speed: 0.3, period: 12.0, rowShift: 0, dim: true}, // slow background
		{speed: 0.8, period: 8.0, rowShift: 1, dim: false}, // mid
		{speed: 1.5, period: 5.0, rowShift: 2, dim: false}, // fast foreground
	}

	for _, layer := range layers {
		for col := 0; col < w; col++ {
			f := float64(frame)
			x := float64(col)

			// Scrolling wave pattern
			phase := x/layer.period - f*layer.speed*0.1
			idx := int(math.Floor(math.Mod(phase+100, float64(len(blockChars))))) % len(blockChars)
			if idx < 0 {
				idx += len(blockChars)
			}

			ch := blockChars[idx]
			if ch == " " {
				continue
			}

			// Vertical placement: wave-based
			rowWave := math.Sin(x*0.1 + f*0.05 + float64(layer.rowShift)*2.0)
			row := clampInt(int(math.Round((rowWave+1.0)*0.5*float64(rows-1))), 0, rows-1)

			// Only overwrite empty cells (layer ordering)
			if grid[row][col].ch == "" {
				grid[row][col] = animCell{ch: ch, dim: layer.dim}
			}
		}
	}
}

// ──────────────────────────────────────────────────
// Animation 5: Constellations (twinkling starfield)
// ──────────────────────────────────────────────────

var starChars = []string{"·", "∘", "•", "✦", "⊹", "✧"}

func renderConstellations(grid [][]animCell, w, frame, localF int) {
	rows := len(grid)
	numStars := w / 3
	if numStars > 40 {
		numStars = 40
	}

	for i := 0; i < numStars; i++ {
		// Fixed star positions (deterministic)
		x := pseudoRand(i*8761+17) % w
		row := pseudoRand(i*3571+31) % rows
		charBase := pseudoRand(i*2909+41) % len(starChars)

		// Twinkle: each star has its own phase
		twinklePhase := float64(pseudoRand(i*1433+53)) * 0.01
		brightness := math.Sin(float64(frame)*0.08+twinklePhase*6.28)*0.5 + 0.5

		if brightness < 0.2 {
			continue // star is dim/invisible this frame
		}

		// Brighter stars use bolder characters
		charIdx := charBase
		if brightness > 0.8 {
			charIdx = clampInt(charBase+1, 0, len(starChars)-1)
		}

		bold := brightness > 0.7
		dim := brightness < 0.4

		if grid[row][x].ch == "" {
			grid[row][x] = animCell{ch: starChars[charIdx], bold: bold, dim: dim}
		}
	}

	// Occasional connecting lines between nearby stars (constellation effect)
	for i := 0; i < numStars-1; i++ {
		x1 := pseudoRand(i*8761+17) % w
		r1 := pseudoRand(i*3571+31) % rows
		x2 := pseudoRand((i+1)*8761+17) % w
		r2 := pseudoRand((i+1)*3571+31) % rows

		// Only connect close stars
		dist := math.Abs(float64(x2 - x1))
		if dist > 3 && dist < 10 && r1 == r2 {
			// Draw connection with dim dots
			pulse := math.Sin(float64(frame)*0.04+float64(i)*0.5)*0.5 + 0.5
			if pulse > 0.6 {
				lo, hi := x1, x2
				if lo > hi {
					lo, hi = hi, lo
				}
				for cx := lo + 1; cx < hi; cx++ {
					if cx >= 0 && cx < w && grid[r1][cx].ch == "" {
						grid[r1][cx] = animCell{ch: "·", dim: true}
					}
				}
			}
		}
	}
}

// ──────────────────────────────────────────────────
// Animation 6: Wormhole (psychedelic tunnel)
// ──────────────────────────────────────────────────

var tunnelChars = []string{"·", "∘", "◌", "◎", "◉"}

func renderWormhole(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	f := float64(frame)
	centerX := float64(w-1)*0.5 + math.Sin(f*0.031)*float64(w)*0.18
	centerY := float64(rows-1)*0.5 + math.Sin(f*0.049)*0.35

	for row := 0; row < rows; row++ {
		for col := 0; col < w; col++ {
			dx := float64(col) - centerX
			dy := (float64(row) - centerY) * 1.8
			dist := math.Sqrt(dx*dx + dy*dy)
			ring := math.Sin(dist*0.52 - f*0.28)
			swirl := math.Cos(float64(col)*0.09 - float64(row)*0.7 + f*0.11)
			intensity := ring*0.72 + swirl*0.28

			if intensity < 0.08 {
				continue
			}

			cell := animCell{
				ch:          tunnelChars[0],
				customColor: colorFromPalette([]color.Color{demoViolet, demoPink, demoCyan, demoWhite}, dist*0.03+f*0.012),
			}

			switch {
			case intensity > 0.88:
				cell.ch = tunnelChars[4]
				cell.bold = true
			case intensity > 0.62:
				cell.ch = tunnelChars[3]
				cell.bold = true
			case intensity > 0.38:
				cell.ch = tunnelChars[2]
			case intensity > 0.18:
				cell.ch = tunnelChars[1]
				cell.dim = true
			default:
				cell.ch = tunnelChars[0]
				cell.dim = true
			}

			grid[row][col] = cell
		}
	}

	coreX := clampInt(int(math.Round(centerX)), 0, w-1)
	mid := rows / 2
	setCell(grid, mid, coreX, animCell{ch: "◉", bold: true, customColor: demoWhite})
	setCellIfEmpty(grid, mid, coreX-1, animCell{ch: "◎", customColor: demoPink})
	setCellIfEmpty(grid, mid, coreX+1, animCell{ch: "◎", customColor: demoCyan})
}

// ──────────────────────────────────────────────────
// Animation 7: Laser Show (neon sweep grid)
// ──────────────────────────────────────────────────

func renderLaserShow(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	mid := rows / 2
	stepSpan := minInt(3, rows/2)
	f := float64(frame)
	sweep := int((math.Sin(f*0.09) + 1.0) * 0.5 * float64(w-1))
	diagA := int((math.Sin(f*0.07+1.4) + 1.0) * 0.5 * float64(w-1))
	diagB := int((math.Sin(f*0.11+3.2) + 1.0) * 0.5 * float64(w-1))
	scanner := int(math.Mod(f*1.8, float64(w+12))) - 6

	for dx := -6; dx <= 6; dx++ {
		x := sweep + dx
		if x < 0 || x >= w {
			continue
		}
		ch := "─"
		if absInt(dx) <= 1 {
			ch = "═"
		}
		col := demoPink
		if dx%2 == 0 {
			col = demoCyan
		}
		setCell(grid, mid, x, animCell{ch: ch, bold: absInt(dx) <= 2, customColor: col})
	}

	for step := -stepSpan; step <= stepSpan; step++ {
		setCell(grid, clampInt(mid+step, 0, rows-1), diagA+step, animCell{ch: "╲", bold: true, customColor: demoCyan})
		setCell(grid, clampInt(mid-step, 0, rows-1), diagB+step, animCell{ch: "╱", bold: true, customColor: demoLime})
	}

	for row := 0; row < rows; row++ {
		for dx := -1; dx <= 1; dx++ {
			x := scanner + dx
			if x < 0 || x >= w {
				continue
			}
			ch := "│"
			if dx == 0 {
				ch = "█"
			}
			setCell(grid, row, x, animCell{ch: ch, bold: dx == 0, customColor: demoAmber})
		}
	}

	for _, x := range []int{sweep, diagA, diagB, scanner} {
		if x < 0 || x >= w {
			continue
		}
		row := clampInt(int(math.Round(float64(mid)+math.Sin(float64(x)*0.12+f*0.09)*float64(stepSpan))), 0, rows-1)
		setCell(grid, row, x, animCell{ch: "✦", bold: true, customColor: demoWhite})
	}

	for i := 0; i < w/7; i++ {
		x := (pseudoRand(i*971+frame*7) + frame*3) % w
		row := pseudoRand(i*557+frame*5) % rows
		setCellIfEmpty(grid, row, x, animCell{ch: "·", dim: true, customColor: demoViolet})
	}
}

// ──────────────────────────────────────────────────
// Animation 8: Fireworks (skyburst finale)
// ──────────────────────────────────────────────────

func renderFireworks(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	type burst struct {
		xFrac float64
		yFrac float64
		delay int
		col   color.Color
	}

	bursts := []burst{
		{xFrac: 0.16, yFrac: 0.12, delay: 0, col: demoPink},
		{xFrac: 0.36, yFrac: 0.32, delay: 17, col: demoAmber},
		{xFrac: 0.58, yFrac: 0.18, delay: 34, col: demoCyan},
		{xFrac: 0.78, yFrac: 0.35, delay: 51, col: demoViolet},
	}

	rows := len(grid)
	for i, burst := range bursts {
		cx := clampInt(int(math.Round(float64(w-1)*burst.xFrac)), 0, w-1)
		targetRow := clampInt(int(math.Round(float64(rows-1)*burst.yFrac)), 0, rows-1)
		age := (frame + burst.delay) % 72

		if age < 14 {
			progress := float64(age) / 14.0
			y := clampInt(int(math.Round(float64(rows-1)-progress*float64((rows-1)-targetRow))), 0, rows-1)
			x := clampInt(cx-int((1.0-progress)*2.0*math.Sin(float64(age)+float64(i)*0.7)), 0, w-1)
			setCell(grid, y, x, animCell{ch: "│", bold: true, customColor: burst.col})
			setCellIfEmpty(grid, y+1, x, animCell{ch: "·", dim: true, customColor: burst.col})
			continue
		}

		explode := age - 14
		if explode > 42 {
			continue
		}

		radiusX := float64(explode) * 0.55
		radiusY := math.Min(float64(rows)/2.5, float64(explode)*0.16)
		headChar := "·"
		switch {
		case explode < 6:
			headChar = "✺"
		case explode < 14:
			headChar = "✦"
		case explode < 26:
			headChar = "•"
		}

		if explode < 8 {
			setCell(grid, targetRow, cx, animCell{ch: "✹", bold: true, customColor: demoWhite})
		}

		for spark := 0; spark < 12; spark++ {
			angle := float64(spark)*math.Pi/6.0 + float64(i)*0.3
			x := cx + int(math.Round(math.Cos(angle)*radiusX))
			row := targetRow + int(math.Round(math.Sin(angle)*radiusY))
			setCell(grid, row, x, animCell{ch: headChar, bold: explode < 14, customColor: burst.col})
		}
	}
}

// ──────────────────────────────────────────────────
// Animation 9: Plasma (demoscene interference field)
// ──────────────────────────────────────────────────

var plasmaChars = []string{" ", "░", "▒", "▓", "█"}

func renderPlasma(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	f := float64(frame)
	centerX := float64(w)*0.5 + math.Sin(f*0.02)*float64(w)*0.15
	centerY := float64(rows-1) * 0.5

	for row := 0; row < rows; row++ {
		for col := 0; col < w; col++ {
			x := float64(col)
			y := float64(row) * 1.4
			v := math.Sin(x*0.11+f*0.09) +
				math.Sin((x+y)*0.07-f*0.05) +
				math.Sin(math.Hypot(x-centerX, y-centerY*1.4)*0.16-f*0.12)

			intensity := (v + 3.0) / 6.0
			if intensity < 0.28 {
				continue
			}

			idx := clampInt(int(intensity*float64(len(plasmaChars)-1)), 1, len(plasmaChars)-1)
			grid[row][col] = animCell{
				ch:          plasmaChars[idx],
				bold:        intensity > 0.72,
				dim:         intensity < 0.42,
				customColor: colorFromPalette([]color.Color{demoPink, demoViolet, demoBlue, demoCyan, demoWhite, demoAmber}, intensity+f*0.01+x*0.004),
			}
		}
	}
}

// ──────────────────────────────────────────────────
// Animation 10: Hyperspace (warp-speed star streaks)
// ──────────────────────────────────────────────────

func renderHyperspace(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	centerX := float64(w-1) * 0.5
	centerY := float64(rows-1) * 0.5
	stars := minInt(w/2, 28)
	if stars < 10 {
		stars = w
	}

	for i := 0; i < stars; i++ {
		seed := pseudoRand(i*811 + 97)
		edgeX := 0
		if seed%2 == 0 {
			edgeX = w - 1
		}
		edgeRow := pseudoRand(i*1733+43) % rows
		speed := 0.006 * float64((seed%6)+8)
		progress := fract(float64(frame)*speed + float64(seed%100)*0.01)
		head := math.Pow(progress, 1.9)
		tail := math.Max(0, progress-0.12)
		tail = math.Pow(tail, 1.9)

		x0 := centerX + (float64(edgeX)-centerX)*tail
		x1 := centerX + (float64(edgeX)-centerX)*head
		r0 := centerY + (float64(edgeRow)-centerY)*tail
		r1 := centerY + (float64(edgeRow)-centerY)*head

		steps := absInt(int(math.Round(x1-x0))) + 1
		if steps < 1 {
			steps = 1
		}

		trailChar := "─"
		if math.Abs(r1-r0) > 0.35 {
			if r1 < r0 {
				trailChar = "╱"
			} else {
				trailChar = "╲"
			}
		}

		for step := 0; step < steps; step++ {
			t := 0.0
			if steps > 1 {
				t = float64(step) / float64(steps-1)
			}
			x := int(math.Round(x0 + (x1-x0)*t))
			row := int(math.Round(r0 + (r1-r0)*t))
			setCellIfEmpty(grid, row, x, animCell{ch: trailChar, dim: progress < 0.35, customColor: demoBlue})
		}

		headChar := "•"
		switch {
		case progress > 0.72:
			headChar = "✦"
		case progress > 0.46:
			headChar = "◦"
		}
		setCell(grid, int(math.Round(r1)), int(math.Round(x1)), animCell{ch: headChar, bold: progress > 0.55, customColor: demoWhite})
	}

	setCellIfEmpty(grid, int(math.Round(centerY)), int(math.Round(centerX)), animCell{ch: "✧", bold: true, customColor: demoCyan})
}

// ──────────────────────────────────────────────────
// Animation 11: Time Rift (glitching spacetime tear)
// ──────────────────────────────────────────────────

var glitchChars = []string{"·", "░", "▒", "▓"}

func renderTimeRift(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	mid := rows / 2
	f := float64(frame)
	crackX := int((math.Sin(f*0.057) + 1.0) * 0.5 * float64(w-1))

	for row := 0; row < rows; row++ {
		offset := int(math.Round(math.Sin(f*0.21+float64(row)*1.7) * 2.0))
		x := crackX + offset + (row - mid)
		ch := "│"
		if row < mid {
			ch = "╲"
		} else if row > mid {
			ch = "╱"
		}
		setCell(grid, row, x, animCell{ch: ch, bold: true, customColor: demoWhite})
		setCellIfEmpty(grid, row, x-1, animCell{ch: "╳", customColor: demoPink})
		setCellIfEmpty(grid, row, x+1, animCell{ch: "╳", customColor: demoCyan})
	}

	for row := 0; row < rows; row++ {
		tear := (frame*2 + row*11) % (w + 10)
		tear -= 5
		for dx := 0; dx < 5; dx++ {
			x := tear + dx
			ch := "▀"
			if row > mid {
				ch = "▄"
			} else if row == mid {
				ch = "▌"
			}
			col := demoAmber
			if dx%2 == 0 {
				col = demoViolet
			}
			if dx == 2 {
				col = demoCyan
			}
			setCellIfEmpty(grid, row, x, animCell{ch: ch, dim: dx != 2, customColor: col})
		}
	}

	for i := 0; i < w/6; i++ {
		x := (pseudoRand(i*139+frame*17) + frame*5) % w
		row := pseudoRand(i*733+frame*13) % rows
		glyph := glitchChars[(i+frame/3)%len(glitchChars)]
		col := demoCyan
		if (i+frame/5)%2 == 0 {
			col = demoPink
		}
		setCellIfEmpty(grid, row, x, animCell{ch: glyph, dim: glyph == "·", customColor: col})
	}
}

// ──────────────────────────────────────────────────
// Animation 12: Reactor Pulse (orbital core shockwave)
// ──────────────────────────────────────────────────

func renderReactorPulse(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	centerY := rows / 2
	f := float64(frame)
	maxRadius := minInt(w/2, 12)
	if maxRadius < 2 {
		maxRadius = 2
	}

	centerX := clampInt(int(math.Round(float64(w-1)*0.5+math.Sin(f*0.041)*float64(w)*0.12)), 0, w-1)
	radius := clampInt(int(math.Round(fract(f*0.032)*float64(maxRadius))), 0, maxRadius)

	setCell(grid, centerY, centerX, animCell{ch: "◉", bold: true, customColor: demoWhite})
	setCellIfEmpty(grid, centerY-1, centerX, animCell{ch: "│", customColor: demoCyan})
	setCellIfEmpty(grid, centerY+1, centerX, animCell{ch: "│", customColor: demoPink})

	for dx := -2; dx <= 2; dx++ {
		if dx == 0 {
			continue
		}
		cell := animCell{ch: "•", customColor: demoPink}
		if absInt(dx) == 2 {
			cell.ch = "·"
			cell.dim = true
		}
		setCellIfEmpty(grid, centerY, centerX+dx, cell)
	}

	if radius > 0 {
		left := centerX - radius
		right := centerX + radius
		setCell(grid, centerY, left, animCell{ch: "◌", bold: true, customColor: demoAmber})
		setCell(grid, centerY, right, animCell{ch: "◌", bold: true, customColor: demoAmber})

		for x := left + 1; x < right; x++ {
			if absInt(x-centerX)%2 == (frame/2)%2 {
				setCellIfEmpty(grid, centerY, x, animCell{ch: "═", dim: true, customColor: demoViolet})
			}
		}

		if radius > 2 {
			setCellIfEmpty(grid, centerY-2, centerX-radius/2, animCell{ch: "•", customColor: demoCyan})
			setCellIfEmpty(grid, centerY+2, centerX+radius/2, animCell{ch: "•", customColor: demoPink})
		}
	}

	orbit := int(math.Round(math.Sin(f*0.18) * float64(minInt(8, maxRadius))))
	setCell(grid, centerY-2, centerX+orbit, animCell{ch: "◦", bold: true, customColor: demoLime})
	setCell(grid, centerY+2, centerX-orbit, animCell{ch: "◦", bold: true, customColor: demoCyan})
}

// ──────────────────────────────────────────────────
// Animation 13: Alien Run (retro side-scroller)
// ──────────────────────────────────────────────────

func renderAlienRun(grid [][]animCell, w, frame int) {
	if w <= 0 {
		return
	}

	rows := len(grid)
	if rows < 4 {
		renderParticles(grid, w, frame)
		return
	}

	scroll := frame / 2
	skyRows := maxInt(1, rows-4)
	logoW := asciiLogoWidth()
	runnerX := clampInt(maxInt(logoW+6, w/2), 6, maxInt(6, w-4))

	surfaceRowAt := func(worldX int) int {
		segment := worldX / 7
		surface := rows - 2
		switch pseudoRand(segment*61+11) % 6 {
		case 0, 1:
			surface = rows - 3
		case 2:
			if segment%4 == 0 {
				surface = rows - 4
			}
		}
		return clampInt(surface, maxInt(1, rows-4), rows-2)
	}

	for i := 0; i < minInt(w/3, 28); i++ {
		baseX := pseudoRand(i*947+13) % w
		x := (baseX - frame/3) % w
		if x < 0 {
			x += w
		}
		row := pseudoRand(i*557+17) % skyRows
		twinkle := (frame/4 + i) % 6
		cell := animCell{ch: "·", dim: true, customColor: demoBlue}
		switch {
		case twinkle == 0:
			cell.ch = "✦"
			cell.bold = true
			cell.dim = false
			cell.customColor = demoWhite
		case twinkle < 3:
			cell.ch = "•"
			cell.customColor = demoCyan
		case twinkle == 5:
			cell.ch = "∘"
			cell.customColor = demoPink
		}
		setCellIfEmpty(grid, row, x, cell)
	}

	planetX := w - 1 - ((frame / 4) % (w + 12)) + 6
	planetRow := clampInt(1, 0, rows-1)
	setCellIfEmpty(grid, planetRow, planetX, animCell{ch: "◉", bold: true, customColor: demoAmber})
	setCellIfEmpty(grid, planetRow, planetX-1, animCell{ch: "◌", customColor: demoPink})
	setCellIfEmpty(grid, planetRow, planetX+1, animCell{ch: "◌", customColor: demoCyan})
	setCellIfEmpty(grid, planetRow+1, planetX, animCell{ch: "═", dim: true, customColor: demoViolet})

	ufoX := w - 1 - ((frame * 2) % (w + 18)) + 8
	ufoRow := clampInt(1+int(math.Round(math.Sin(float64(frame)*0.08))), 0, maxInt(0, rows-4))
	setCellIfEmpty(grid, ufoRow, ufoX-1, animCell{ch: "╭", customColor: demoCyan})
	setCellIfEmpty(grid, ufoRow, ufoX, animCell{ch: "▔", bold: true, customColor: demoWhite})
	setCellIfEmpty(grid, ufoRow, ufoX+1, animCell{ch: "╮", customColor: demoPink})
	setCellIfEmpty(grid, ufoRow+1, ufoX, animCell{ch: "◡", customColor: demoAmber})

	for col := 0; col < w; col++ {
		worldX := col + scroll

		ridgeWave := math.Sin(float64(worldX)*0.065) + 0.55*math.Sin(float64(worldX)*0.16)
		ridgeTop := clampInt(rows-4-int(math.Round(ridgeWave)), 1, rows-3)
		for row := ridgeTop; row < rows-3; row++ {
			cell := animCell{ch: "▒", dim: true, customColor: demoViolet}
			if row == ridgeTop {
				cell.ch = "▔"
				cell.customColor = demoBlue
			}
			setCellIfEmpty(grid, row, col, cell)
		}

		platformSeg := worldX / 11
		if pseudoRand(platformSeg*109+27)%5 == 0 {
			platformRow := clampInt(rows-4-pseudoRand(platformSeg*131+17)%2, 1, rows-4)
			if worldX%11 < 4 {
				setCell(grid, platformRow, col, animCell{ch: "═", bold: true, customColor: demoCyan})
				if worldX%11 == 1 {
					setCellIfEmpty(grid, platformRow-1, col, animCell{ch: "✦", bold: true, customColor: demoAmber})
				}
			}
		}

		surface := surfaceRowAt(worldX)
		for row := surface; row < rows; row++ {
			cell := animCell{ch: "█", customColor: demoBlue}
			switch {
			case row == surface:
				cell.ch = "▀"
				cell.customColor = demoLime
			case row == surface+1:
				cell.ch = "▓"
				cell.customColor = demoViolet
			}
			setCell(grid, row, col, cell)
		}

		if surface > 0 {
			floraSeg := worldX / 9
			switch pseudoRand(floraSeg*89+23) % 8 {
			case 0:
				if worldX%9 == 2 {
					setCellIfEmpty(grid, surface-1, col, animCell{ch: "▲", customColor: demoPink})
				}
			case 1:
				if worldX%9 == 4 {
					setCellIfEmpty(grid, surface-1, col, animCell{ch: "◇", bold: true, customColor: demoCyan})
				}
			case 2:
				if worldX%9 == 6 {
					setCellIfEmpty(grid, surface-1, col, animCell{ch: "⌇", dim: true, customColor: demoAmber})
				}
			}
		}
	}

	jumpPhase := frame % 44
	jumpOffset := 0
	switch {
	case jumpPhase >= 10 && jumpPhase < 14:
		jumpOffset = 1
	case jumpPhase >= 14 && jumpPhase < 18:
		jumpOffset = 2
	case jumpPhase >= 18 && jumpPhase < 24:
		jumpOffset = 1
	}

	runnerSurface := surfaceRowAt(runnerX + scroll)
	footRow := clampInt(runnerSurface-jumpOffset, 1, rows-1)
	bodyRow := clampInt(footRow-1, 0, rows-1)
	headRow := clampInt(footRow-2, 0, rows-1)

	legLeft, legRight := "╱", "╲"
	if (frame/3)%2 == 1 {
		legLeft, legRight = "╲", "╱"
	}

	setCell(grid, headRow, runnerX, animCell{ch: "◉", bold: true, customColor: demoLime})
	setCellIfEmpty(grid, headRow, runnerX-1, animCell{ch: "·", dim: true, customColor: demoPink})
	setCell(grid, bodyRow, runnerX, animCell{ch: "█", bold: true, customColor: demoLime})
	setCell(grid, bodyRow, runnerX-1, animCell{ch: "╱", customColor: demoCyan})
	setCell(grid, bodyRow, runnerX+1, animCell{ch: "╲", customColor: demoPink})
	setCell(grid, footRow, runnerX-1, animCell{ch: legLeft, bold: true, customColor: demoAmber})
	setCell(grid, footRow, runnerX+1, animCell{ch: legRight, bold: true, customColor: demoAmber})
	setCellIfEmpty(grid, footRow, runnerX, animCell{ch: "┴", customColor: demoWhite})
}

// ──────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// pseudoRand returns a deterministic "random" number from a seed.
// Simple hash for reproducible particle/star positions.
func pseudoRand(seed int) int {
	seed = (seed ^ 61) ^ (seed >> 16)
	seed = seed + (seed << 3)
	seed = seed ^ (seed >> 4)
	seed = seed * 0x27d4eb2d
	seed = seed ^ (seed >> 15)
	if seed < 0 {
		seed = -seed
	}
	return seed
}
