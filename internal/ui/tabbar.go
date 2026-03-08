package ui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// TabNames are the display names for each tab.
var TabNames = []string{
	"Dashboard",
	"Aliases",
	"Functions",
	"Packages",
	"Tmux",
	"Search",
	"Doctor",
	"Config",
}

// TabBarLineCount returns the rendered height of the tab chrome.
func TabBarLineCount(state *HeaderState) int {
	if state != nil && state.IsCompact() {
		if state.showCompactTabAccent {
			return 3
		}
		return 2
	}
	return 2
}

func showCompactTabAccent(state *HeaderState) bool {
	if state == nil || !state.IsCompact() {
		return false
	}
	return state.showCompactTabAccent
}

// RenderTabBar renders a mode-synced animated tab bar with a full-width bleed strip.
func RenderTabBar(activeIdx int, focused bool, width int, styles *Styles, state *HeaderState) string {
	lineWidth := width - 2 // account for doc padding
	if lineWidth <= 0 {
		return strings.Repeat("\n", TabBarLineCount(state))
	}

	frame := 0
	mode := animWaveDots
	localF := 0
	if state != nil {
		frame = state.frame + 17
		mode = state.currentMode()
		localF = (state.localFrame() + 17) % framesPerAnim
	}

	grid := buildAnimationGrid(2, lineWidth, frame, mode, localF, true)
	seedChromeBleed(grid[0], frame)
	seedChromeAccent(grid[1], frame, mode)
	overlayTabs(grid[0], activeIdx, focused, width, styles)

	baseColors := lipgloss.Blend1D(lineWidth, styles.GradientStart, styles.GradientEnd)
	compact := state != nil && state.IsCompact()
	showAccent := !compact || showCompactTabAccent(state)

	var doc strings.Builder
	if compact {
		doc.WriteString(renderStaticDivider(lineWidth, styles))
		doc.WriteString("\n")
	}
	doc.WriteString(renderTabRow(grid[0], baseColors))
	doc.WriteString("\n")
	if showAccent {
		doc.WriteString(renderAccentRow(grid[1], baseColors))
		doc.WriteString("\n")
	}
	return doc.String()
}

// RenderTabBarPlain renders a minimal tab bar for very narrow terminals.
func RenderTabBarPlain(activeIdx int, styles *Styles) string {
	var tabs []string
	for i, name := range TabNames {
		if i == activeIdx {
			tabs = append(tabs, styles.ActiveTab.Render(name))
		} else {
			tabs = append(tabs, styles.InactiveTab.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...) + "\n"
}

func overlayTabs(row []animCell, activeIdx int, focused bool, width int, styles *Styles) {
	pos := 0
	for i, name := range TabNames {
		if pos >= len(row) {
			return
		}

		if i == activeIdx {
			pos = overlayText(row, pos, "●", animCell{
				customColor: styles.AccentColor,
				bold:        true,
			})
			pos = overlayText(row, pos, " "+name+"  ", animCell{
				customColor: demoWhite,
				bold:        true,
				underline:   focused,
			})
			continue
		}

		label := name
		if width >= 100 {
			label = fmt.Sprintf("%d %s", i+1, name)
		}
		pos = overlayText(row, pos, label+"  ", animCell{
			customColor: styles.SubtleColor,
			dim:         true,
		})
	}
}

func overlayText(row []animCell, start int, text string, style animCell) int {
	pos := start
	for _, r := range []rune(text) {
		if pos >= len(row) {
			break
		}
		if r == ' ' {
			row[pos] = animCell{ch: " "}
		} else {
			cell := style
			cell.ch = string(r)
			row[pos] = cell
		}
		pos++
	}
	return pos
}

func overlayPaletteText(row []animCell, start int, text string, palette []color.Color, style animCell) int {
	total := 0
	for _, r := range []rune(text) {
		if r != ' ' {
			total++
		}
	}

	pos := start
	colorIdx := 0
	for _, r := range []rune(text) {
		if pos >= len(row) {
			break
		}
		if r == ' ' {
			row[pos] = animCell{ch: " "}
			pos++
			continue
		}

		cell := style
		cell.ch = string(r)
		if len(palette) > 0 {
			palettePos := 0
			if len(palette) > 1 && total > 1 {
				palettePos = colorIdx * (len(palette) - 1) / (total - 1)
			}
			cell.customColor = palette[palettePos]
		}
		row[pos] = cell
		pos++
		colorIdx++
	}
	return pos
}

func seedChromeBleed(row []animCell, frame int) {
	for col := range row {
		if row[col].ch != "" {
			continue
		}
		switch {
		case (col+frame/2)%17 == 0:
			row[col] = animCell{ch: "·", dim: true, customColor: demoViolet}
		case (col*3+frame)%29 == 0:
			row[col] = animCell{ch: "•", customColor: demoCyan}
		case (col+frame)%37 == 0:
			row[col] = animCell{ch: "∘", dim: true, customColor: demoPink}
		}
	}
}

func seedChromeAccent(row []animCell, frame int, mode animMode) {
	for col := range row {
		if row[col].ch != "" {
			continue
		}

		cell := animCell{ch: "━", dim: true}
		switch mode {
		case animLaserShow, animHyperspace:
			cell.ch = "═"
			cell.customColor = demoCyan
		case animFireworks:
			cell.ch = "┄"
			cell.customColor = demoAmber
		case animTimeRift:
			cell.ch = "▀"
			cell.customColor = demoPink
		case animReactorPulse:
			cell.ch = "━"
			cell.customColor = demoViolet
		default:
			if (col+frame)%9 == 0 {
				cell.ch = "═"
				cell.bold = true
			}
		}

		if (col+frame)%13 == 0 {
			cell.bold = true
			cell.dim = false
		}
		row[col] = cell
	}
}

func renderStaticDivider(width int, styles *Styles) string {
	return RenderGradientLine(width, styles.GradientStart, styles.GradientEnd)
}

func renderTabRow(row []animCell, colors []color.Color) string {
	var b strings.Builder
	for col, cell := range row {
		if cell.ch == "" {
			b.WriteRune(' ')
			continue
		}
		base := colors[col]
		if cell.customColor != nil {
			base = cell.customColor
		}
		b.WriteString(renderAnimCell(cell, base))
	}
	return b.String()
}

func renderAccentRow(row []animCell, colors []color.Color) string {
	var b strings.Builder
	for col, cell := range row {
		base := colors[col]
		if cell.ch == "" {
			cell = animCell{ch: "━", dim: true}
		}
		if cell.customColor != nil {
			base = cell.customColor
		}
		b.WriteString(renderAnimCell(cell, base))
	}
	return b.String()
}
