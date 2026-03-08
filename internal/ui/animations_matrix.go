package ui

import "math"

var matrixGlyphs = []rune("01[]{}<>/\\|+-=*#$%&?;:ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func renderMatrixRain(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	cycleFrame := frame
	if framesPerAnim > 0 {
		cycleFrame = frame % framesPerAnim
	}
	phase := float64(cycleFrame) / float64(maxInt(framesPerAnim, 1))

	spacing := 2
	if w < 26 {
		spacing = 1
	}

	for base := 0; base < w; base += spacing {
		offsetSpan := minInt(spacing, w-base)
		col := base
		if offsetSpan > 1 {
			col += pseudoRand(base*41+17) % offsetSpan
		}

		seed := float64(pseudoRand(base*97+31)%1000) / 1000.0
		speed := 1 + pseudoRand(base*67+11)%3
		trail := maxInt(3, rows/2+1+pseudoRand(base*53+7)%maxInt(1, rows))
		gap := 2 + pseudoRand(base*89+13)%4
		span := rows + trail + gap
		headY := fract(seed+phase*float64(speed))*float64(span) - float64(trail)

		for row := 0; row < rows; row++ {
			dist := headY - float64(row)
			if dist < 0 || dist > float64(trail) {
				continue
			}

			intensity := 1.0 - dist/float64(trail+1)
			if intensity <= 0.03 {
				continue
			}

			glyphPhase := pseudoRand(base*131+row*61+cycleFrame*29+int(math.Round(dist*11))) % len(matrixGlyphs)
			ch := string(matrixGlyphs[glyphPhase])

			cell := animCell{ch: ch, customColor: demoMatrixTail, dim: true}
			switch {
			case dist < 0.55:
				cell.bold = true
				cell.customColor = demoMatrixHead
			case intensity > 0.74:
				cell.customColor = demoMatrixGlow
			case intensity > 0.48:
				cell.customColor = demoMatrixCore
				cell.dim = false
			case intensity > 0.20:
				cell.customColor = demoMatrixTrail
			}
			setCell(grid, row, col, cell)

			if col+1 < w && dist < 0.4 && (base/spacing+cycleFrame)%5 == 0 {
				setCellIfEmpty(grid, row, col+1, animCell{ch: string(matrixGlyphs[(glyphPhase+7)%len(matrixGlyphs)]), customColor: demoMatrixCore, dim: true})
			}
		}

		glintRow := clampInt(int(math.Round(headY))-trail/2, 0, rows-1)
		if rows > 1 && glintRow >= 0 && glintRow < rows && (base/spacing+cycleFrame)%7 == 0 {
			setCellIfEmpty(grid, glintRow, col, animCell{ch: ".", customColor: demoMatrixTrail, dim: true})
		}
	}

	for col := 0; col < w; col++ {
		if (col*5+cycleFrame)%37 == 0 {
			row := pseudoRand(col*17+cycleFrame*13) % rows
			setCellIfEmpty(grid, row, col, animCell{ch: ".", customColor: demoMatrixTail, dim: true})
		}
	}
}
