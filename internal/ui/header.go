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

// AnimTickMsg signals an animation frame advance.
type AnimTickMsg struct{}

func animTick() tea.Cmd {
	return tea.Tick(66*time.Millisecond, func(time.Time) tea.Msg {
		return AnimTickMsg{}
	})
}

// HeaderState holds animation state for the header dot wave.
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

// dot characters by intensity (trough → peak)
var dotChars = []string{"·", "•", "⦁"}

// RenderHeader renders the gradient ASCII logo with animated dot wave and badges.
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

	// Compute wave field dimensions
	waveStart := logoWidth + 4
	waveEnd := width - 4
	waveWidth := waveEnd - waveStart
	if waveWidth < 0 {
		waveWidth = 0
	}

	// Generate gradient for the wave field
	var waveColors []color.Color
	if waveWidth > 0 {
		waveColors = lipgloss.Blend1D(waveWidth, styles.GradientStart, styles.GradientEnd)
	}

	frame := 0
	if state != nil {
		frame = state.frame
	}

	// Render each logo line + wave dots
	for row := 0; row < logoRows; row++ {
		line := asciiLogo[row]
		runes := []rune(line)

		// Render logo characters with gradient
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

		// Gap between logo and wave
		if waveWidth > 0 {
			lineBuilder.WriteString("    ") // 4-char gap

			// Render wave dots for this row
			for col := 0; col < waveWidth; col++ {
				x := float64(col)
				f := float64(frame)

				// Two overlapping sine waves
				y := math.Sin(x*0.15+f*0.08) + 0.5*math.Sin(x*0.08-f*0.12)
				// y ranges roughly [-1.5, 1.5], map to row [0, 2]
				rowFloat := (y + 1.5) / 3.0 * float64(logoRows-1)
				targetRow := int(math.Round(rowFloat))
				if targetRow < 0 {
					targetRow = 0
				}
				if targetRow >= logoRows {
					targetRow = logoRows - 1
				}

				if targetRow == row {
					// Determine dot intensity from wave amplitude
					absY := math.Abs(y)
					dotIdx := 0
					if absY > 0.8 {
						dotIdx = 2 // heavy
					} else if absY > 0.3 {
						dotIdx = 1 // medium
					}

					// Fade out last 6 columns
					fadeStart := waveWidth - 6
					if col >= fadeStart && fadeStart > 0 {
						dotIdx = 0 // light dot only
					}

					ci := col
					if ci >= len(waveColors) {
						ci = len(waveColors) - 1
					}

					lineBuilder.WriteString(
						lipgloss.NewStyle().
							Foreground(waveColors[ci]).
							Render(dotChars[dotIdx]),
					)
				} else {
					lineBuilder.WriteRune(' ')
				}
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
