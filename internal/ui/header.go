package ui

import (
	"image/color"
	"math"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var asciiLogo = []string{
	` ╔═╗ ╦ ╦ ╔═╗ ╦   ╦   ╦╔═ ╦ ╔╦╗`,
	` ╚═╗ ╠═╣ ║╣  ║   ║   ╠╩╗ ║  ║ `,
	` ╚═╝ ╩ ╩ ╚═╝ ╩═╝ ╩═╝ ╩ ╩ ╩  ╩ `,
}

const logoRows = 3

// Animation modes
type animMode int

const (
	animWaveDots      animMode = iota // sine wave dot field
	animParticles                     // parallax particle drift
	animSineLine                      // continuous sine wave line
	animGradientRain                  // gradient block rain/parallax
	animConstellations                // twinkling star field
	animModeCount                     // sentinel for rotation
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
	frame  int
	width  int
	styles *Styles
}

func NewHeaderState(styles *Styles) HeaderState {
	return HeaderState{styles: styles}
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

func (h *HeaderState) currentMode() animMode {
	return animMode((h.frame / framesPerAnim) % int(animModeCount))
}

// localFrame returns the frame number within the current animation cycle.
func (h *HeaderState) localFrame() int {
	return h.frame % framesPerAnim
}

// RenderHeader renders the gradient ASCII logo with animated effects and badges.
func RenderHeader(version, os, arch string, width int, styles *Styles, state *HeaderState) string {
	var doc strings.Builder

	// Compute logo rune width
	logoWidth := 0
	for _, line := range asciiLogo {
		if w := len([]rune(line)); w > logoWidth {
			logoWidth = w
		}
	}

	// Generate gradient for the logo
	logoColors := lipgloss.Blend1D(logoWidth, styles.GradientStart, styles.GradientEnd)

	// Compute animation field dimensions
	fieldStart := logoWidth + 4
	fieldEnd := width - 4
	fieldWidth := fieldEnd - fieldStart
	if fieldWidth < 0 {
		fieldWidth = 0
	}

	// Generate gradient colors for the field
	var fieldColors []color.Color
	if fieldWidth > 0 {
		fieldColors = lipgloss.Blend1D(fieldWidth, styles.GradientStart, styles.GradientEnd)
	}

	frame := 0
	mode := animWaveDots
	localF := 0
	if state != nil {
		frame = state.frame
		mode = state.currentMode()
		localF = state.localFrame()
	}

	// Build a 3×fieldWidth grid for the animation
	grid := make([][]animCell, logoRows)
	for r := range grid {
		grid[r] = make([]animCell, fieldWidth)
	}

	if fieldWidth > 0 {
		// Cross-fade: first and last 15 frames blend with neighbors
		switch mode {
		case animWaveDots:
			renderWaveDots(grid, fieldWidth, frame)
		case animParticles:
			renderParticles(grid, fieldWidth, frame)
		case animSineLine:
			renderSineLine(grid, fieldWidth, frame)
		case animGradientRain:
			renderGradientRain(grid, fieldWidth, frame)
		case animConstellations:
			renderConstellations(grid, fieldWidth, frame, localF)
		}

		// Fade-in during first 15 frames of each animation
		if localF < 15 {
			fadeGrid(grid, fieldWidth, float64(localF)/15.0)
		}
		// Fade-out during last 15 frames
		if localF > framesPerAnim-15 {
			fadeGrid(grid, fieldWidth, float64(framesPerAnim-localF)/15.0)
		}
	}

	// Render each logo line + animation field
	for row := 0; row < logoRows; row++ {
		line := asciiLogo[row]
		runes := []rune(line)

		var lineBuilder strings.Builder
		for i, r := range runes {
			if r == ' ' {
				lineBuilder.WriteRune(' ')
				continue
			}
			ci := i
			if ci >= len(logoColors) {
				ci = len(logoColors) - 1
			}
			lineBuilder.WriteString(
				lipgloss.NewStyle().
					Foreground(logoColors[ci]).
					Bold(true).
					Render(string(r)),
			)
		}

		if fieldWidth > 0 {
			lineBuilder.WriteString("    ") // gap
			for col := 0; col < fieldWidth; col++ {
				cell := grid[row][col]
				if cell.ch == "" || cell.ch == " " {
					lineBuilder.WriteRune(' ')
					continue
				}
				ci := col
				if ci >= len(fieldColors) {
					ci = len(fieldColors) - 1
				}
				c := fieldColors[ci]
				if cell.customColor != nil {
					c = cell.customColor
				}
				style := lipgloss.NewStyle().Foreground(c)
				if cell.bold {
					style = style.Bold(true)
				}
				if cell.dim {
					style = style.Faint(true)
				}
				lineBuilder.WriteString(style.Render(cell.ch))
			}
		}

		doc.WriteString(lineBuilder.String())
		doc.WriteString("\n")
	}

	// Version + platform badge on the right
	if version != "" || os != "" {
		var badge strings.Builder
		if version != "" {
			badge.WriteString(styles.VersionBadge.Render(version))
		}
		if os != "" && arch != "" {
			badge.WriteString(" ")
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

// animCell represents one character in the animation grid.
type animCell struct {
	ch          string
	bold        bool
	dim         bool
	customColor color.Color // nil = use gradient
}

// fadeGrid reduces visibility of the grid by the given factor [0, 1].
func fadeGrid(grid [][]animCell, w int, factor float64) {
	for r := 0; r < logoRows; r++ {
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

// ──────────────────────────────────────────────────
// Animation 1: Wave Dots (original)
// ──────────────────────────────────────────────────

var dotChars = []string{"·", "•", "⦁"}

func renderWaveDots(grid [][]animCell, w, frame int) {
	for col := 0; col < w; col++ {
		x := float64(col)
		f := float64(frame)

		y := math.Sin(x*0.15+f*0.08) + 0.5*math.Sin(x*0.08-f*0.12)
		rowFloat := (y + 1.5) / 3.0 * float64(logoRows-1)
		targetRow := clampInt(int(math.Round(rowFloat)), 0, logoRows-1)

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
	{count: 30, speed: 0.3, chars: []string{"·", "⋅", "∘"}, dim: true},             // far (slow, dim)
	{count: 20, speed: 0.7, chars: []string{"·", "•", "◦"}, dim: false},             // mid
	{count: 12, speed: 1.4, chars: []string{"•", "⦁", "✦"}, dim: false, bold: true}, // near (fast, bold)
}

func renderParticles(grid [][]animCell, w, frame int) {
	for _, layer := range particleLayers {
		for i := 0; i < layer.count; i++ {
			// Deterministic position from particle index
			baseX := pseudoRand(i*7919+3) % w
			row := pseudoRand(i*6271+7) % logoRows
			charIdx := pseudoRand(i*4177+13) % len(layer.chars)

			// Move particle based on speed (rightward drift, wrapping)
			x := (baseX + int(float64(frame)*layer.speed)) % w

			// Add slight vertical wobble
			wobble := math.Sin(float64(frame)*0.05 + float64(i)*1.7)
			if math.Abs(wobble) > 0.8 {
				row = (row + 1) % logoRows
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
	f := float64(frame)

	for col := 0; col < w; col++ {
		x := float64(col)

		// Primary wave + harmonic
		y := math.Sin(x*0.12+f*0.06) + 0.3*math.Sin(x*0.25-f*0.1) + 0.2*math.Cos(x*0.07+f*0.04)

		// Map y [-1.5, 1.5] to row [0, 2]
		rowFloat := (y + 1.5) / 3.0 * float64(logoRows-1)
		targetRow := clampInt(int(math.Round(rowFloat)), 0, logoRows-1)

		// Map y to wave character (vertical position within cell)
		subIdx := int((rowFloat - float64(targetRow) + 0.5) * float64(len(waveLineChars)-1))
		subIdx = clampInt(subIdx, 0, len(waveLineChars)-1)

		// Trailing particles behind the wave
		for dr := -1; dr <= 1; dr++ {
			r := targetRow + dr
			if r < 0 || r >= logoRows {
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
	// Three layers of scrolling gradient blocks at different speeds
	type rainLayer struct {
		speed    float64
		period   float64
		rowShift int
		dim      bool
	}
	layers := []rainLayer{
		{speed: 0.3, period: 12.0, rowShift: 0, dim: true},  // slow background
		{speed: 0.8, period: 8.0, rowShift: 1, dim: false},   // mid
		{speed: 1.5, period: 5.0, rowShift: 2, dim: false},   // fast foreground
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
			row := clampInt(int(math.Round((rowWave+1.0)*0.5*float64(logoRows-1))), 0, logoRows-1)

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
	numStars := w / 3
	if numStars > 40 {
		numStars = 40
	}

	for i := 0; i < numStars; i++ {
		// Fixed star positions (deterministic)
		x := pseudoRand(i*8761+17) % w
		row := pseudoRand(i*3571+31) % logoRows
		charBase := pseudoRand(i*2909+41) % len(starChars)

		// Twinkle: each star has its own phase
		twinklePhase := float64(pseudoRand(i*1433+53)) * 0.01
		brightness := math.Sin(float64(frame)*0.08+twinklePhase*6.28) * 0.5 + 0.5

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
		r1 := pseudoRand(i*3571+31) % logoRows
		x2 := pseudoRand((i+1)*8761+17) % w
		r2 := pseudoRand((i+1)*3571+31) % logoRows

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
