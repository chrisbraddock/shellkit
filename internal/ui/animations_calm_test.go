package ui

import "testing"

func TestRenderMoonDriftKeepsCloudBandsBelowLogoZone(t *testing.T) {
	const rows = 7
	const cols = 48

	grid := make([][]animCell, rows)
	for row := range grid {
		grid[row] = make([]animCell, cols)
	}

	renderMoonDrift(grid, cols, 37)

	for row := 1; row <= 3; row++ {
		for col := 0; col < cols; col++ {
			switch grid[row][col].ch {
			case "~", "_":
				t.Fatalf("renderMoonDrift() placed cloud glyph %q in logo row %d col %d", grid[row][col].ch, row, col)
			}
		}
	}
}
