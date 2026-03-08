package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

type animMeta struct {
	mode        animMode
	id          string
	name        string
	category    string
	description string
}

var animationCatalog = []animMeta{
	{mode: animWaveDots, id: "wave-dots", name: "Wave Dots", category: "Calm", description: "Gentle dots rolling in shallow waves."},
	{mode: animParticles, id: "particles", name: "Particles", category: "Calm", description: "Slow particle drift with soft depth."},
	{mode: animSineLine, id: "sine-line", name: "Sine Line", category: "Calm", description: "A clean oscilloscope ribbon."},
	{mode: animGradientRain, id: "gradient-rain", name: "Gradient Rain", category: "Calm", description: "Muted rainfall of drifting color."},
	{mode: animConstellations, id: "constellations", name: "Constellations", category: "Calm", description: "Sparse stars and quiet connections."},
	{mode: animAuroraDrift, id: "aurora-drift", name: "Aurora Drift", category: "Calm", description: "Layered aurora curtains folding overhead."},
	{mode: animTidalLines, id: "tidal-lines", name: "Tidal Lines", category: "Calm", description: "Overlapping wave bands with slow phase shifts."},
	{mode: animMoonDrift, id: "moon-drift", name: "Moon Drift", category: "Calm", description: "A moonlit sky with drifting cloud bands."},
	{mode: animSoftRain, id: "soft-rain", name: "Soft Rain", category: "Calm", description: "Light drizzle with tiny ripples below."},
	{mode: animLanternDrift, id: "lantern-drift", name: "Lantern Drift", category: "Calm", description: "Floating lights rising through the dark."},
	{mode: animWormhole, id: "wormhole", name: "Wormhole", category: "Spectacle", description: "A tunnel pull toward the center."},
	{mode: animLaserShow, id: "laser-show", name: "Laser Show", category: "Spectacle", description: "Neon sweeps and mirrored beams."},
	{mode: animFireworks, id: "fireworks", name: "Fireworks", category: "Spectacle", description: "Skybursts and falling embers."},
	{mode: animPlasma, id: "plasma", name: "Plasma", category: "Spectacle", description: "Classic demoscene plasma turbulence."},
	{mode: animHyperspace, id: "hyperspace", name: "Hyperspace", category: "Sci-Fi", description: "Warp streaks rushing from a vanishing point."},
	{mode: animTimeRift, id: "time-rift", name: "Time Rift", category: "Sci-Fi", description: "A glitching seam in spacetime."},
	{mode: animReactorPulse, id: "reactor-pulse", name: "Reactor Pulse", category: "Sci-Fi", description: "An orbital core cycling with energy."},
	{mode: animStarBattle, id: "star-battle", name: "Star Battle", category: "Sci-Fi", description: "TIE fighters and an X-wing trading fire as the Death Star rolls in."},
	{mode: animAlienRun, id: "alien-run", name: "Alien Run", category: "Arcade", description: "A side-scrolling alien platform dash."},
}

func animationMeta(mode animMode) animMeta {
	for _, meta := range animationCatalog {
		if meta.mode == mode {
			return meta
		}
	}
	return animMeta{
		mode:        mode,
		name:        "Unknown",
		category:    "Other",
		description: "Unknown animation mode.",
	}
}

func modeCategory(mode animMode) string {
	return animationMeta(mode).category
}

func modeDescription(mode animMode) string {
	return animationMeta(mode).description
}

func animationSection(mode animMode) string {
	return modeCategory(mode) + " Animations"
}

func renderAnimationPreview(mode animMode, width, frame int, styles *Styles) string {
	if width <= 0 {
		return ""
	}

	localF := 0
	if framesPerAnim > 0 {
		localF = frame % framesPerAnim
	}

	grid := buildAnimationGrid(4, width, frame, mode, localF, false)
	colors := lipgloss.Blend1D(width, styles.GradientStart, styles.GradientEnd)

	var b strings.Builder
	for rowIdx, row := range grid {
		for col, cell := range row {
			if cell.ch == "" || cell.ch == " " {
				b.WriteRune(' ')
				continue
			}
			base := colors[col]
			if cell.customColor != nil {
				base = cell.customColor
			}
			b.WriteString(renderAnimCell(cell, base))
		}
		if rowIdx < len(grid)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func renderPreviewTitle(mode animMode, styles *Styles) string {
	parts := []string{
		styles.Highlight.Render(modeName(mode)),
		styles.Subtle.Render("[" + modeCategory(mode) + "]"),
	}
	return strings.Join(parts, " ")
}
