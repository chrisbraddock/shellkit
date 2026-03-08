package ui

import "math"

func renderStarBattle(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	for col := 0; col < w; col++ {
		seed := col*11 + frame/3
		if seed%19 == 0 {
			row := (col + frame/9) % rows
			cell := animCell{ch: "·", dim: true, customColor: demoWhite}
			if (col+frame)%7 == 0 {
				cell.ch = "✦"
				cell.bold = true
				cell.customColor = demoCyan
			}
			setCellIfEmpty(grid, row, col, cell)
		}
	}

	centerY := clampInt(rows/2, 1, maxInt(1, rows-2))

	xWingX := (frame*2)%(w+18) - 8
	xWingY := clampInt(centerY+int(math.Round(math.Sin(float64(frame)*0.08))), 1, maxInt(1, rows-2))
	drawXWing(grid, xWingX, xWingY)

	tieAX := w - 1 - ((frame*2 + 8) % (w + 18)) + 6
	tieAY := clampInt(centerY-2, 0, rows-1)
	drawTieFighter(grid, tieAX, tieAY)

	tieBX := w - 1 - ((frame*2 + 28) % (w + 20)) + 10
	tieBY := clampInt(centerY+2, 0, rows-1)
	drawTieFighter(grid, tieBX, tieBY)

	if xWingX > 2 && xWingX < w-4 {
		for shot := 0; shot < 2; shot++ {
			laserX := xWingX + 3 + shot*5 + (frame/2)%6
			laserY := clampInt(xWingY-shot, 0, rows-1)
			if laserX < w {
				setCellIfEmpty(grid, laserY, laserX, animCell{ch: "═", bold: true, customColor: demoLime})
				setCellIfEmpty(grid, laserY, laserX+1, animCell{ch: "•", customColor: demoWhite})
			}
		}
	}

	if tieAX > 0 && tieAX < w {
		enemyLaserX := tieAX - 2 - (frame / 2 % 8)
		enemyLaserY := clampInt(tieAY+1, 0, rows-1)
		if enemyLaserX >= 0 {
			setCellIfEmpty(grid, enemyLaserY, enemyLaserX, animCell{ch: "═", bold: true, customColor: demoPink})
			setCellIfEmpty(grid, enemyLaserY, enemyLaserX-1, animCell{ch: "•", customColor: demoViolet})
		}
	}

	deathStarPhase := fract(float64(frame)/220.0 + 0.35)
	if deathStarPhase > 0.48 {
		entry := math.Sin((deathStarPhase - 0.48) / 0.52 * math.Pi)
		if entry > 0 {
			travel := float64(minInt(10, maxInt(4, w/5)))
			deathStarX := int(math.Round(float64(w+3) - travel*entry))
			deathStarY := clampInt(centerY, 1, maxInt(1, rows-2))
			drawDeathStar(grid, deathStarX, deathStarY)
		}
	}
}

func drawXWing(grid [][]animCell, x, y int) {
	setCellIfEmpty(grid, y-1, x-1, animCell{ch: "╲", customColor: demoCyan})
	setCellIfEmpty(grid, y-1, x+1, animCell{ch: "╱", customColor: demoCyan})
	setCell(grid, y, x-2, animCell{ch: "═", customColor: demoLime})
	setCell(grid, y, x-1, animCell{ch: "╳", customColor: demoWhite})
	setCell(grid, y, x, animCell{ch: "█", bold: true, customColor: demoWhite})
	setCell(grid, y, x+1, animCell{ch: ">", bold: true, customColor: demoAmber})
	setCellIfEmpty(grid, y+1, x-1, animCell{ch: "╱", customColor: demoPink})
	setCellIfEmpty(grid, y+1, x+1, animCell{ch: "╲", customColor: demoPink})
	setCellIfEmpty(grid, y, x-3, animCell{ch: "·", dim: true, customColor: demoBlue})
}

func drawTieFighter(grid [][]animCell, x, y int) {
	setCellIfEmpty(grid, y, x-1, animCell{ch: "│", customColor: demoViolet})
	setCell(grid, y, x, animCell{ch: "◉", bold: true, customColor: demoPink})
	setCellIfEmpty(grid, y, x+1, animCell{ch: "│", customColor: demoViolet})
	setCellIfEmpty(grid, y-1, x-1, animCell{ch: "│", dim: true, customColor: demoBlue})
	setCellIfEmpty(grid, y+1, x-1, animCell{ch: "│", dim: true, customColor: demoBlue})
	setCellIfEmpty(grid, y-1, x+1, animCell{ch: "│", dim: true, customColor: demoBlue})
	setCellIfEmpty(grid, y+1, x+1, animCell{ch: "│", dim: true, customColor: demoBlue})
}

func drawDeathStar(grid [][]animCell, x, y int) {
	setCellIfEmpty(grid, y-1, x-1, animCell{ch: "◜", customColor: demoBlue})
	setCell(grid, y-1, x, animCell{ch: "◠", dim: true, customColor: demoWhite})
	setCellIfEmpty(grid, y-1, x+1, animCell{ch: "◝", customColor: demoBlue})
	setCell(grid, y, x-1, animCell{ch: "◖", customColor: demoWhite})
	setCell(grid, y, x, animCell{ch: "◉", bold: true, customColor: demoWhite})
	setCellIfEmpty(grid, y, x+1, animCell{ch: "◗", customColor: demoBlue})
	setCellIfEmpty(grid, y, x+2, animCell{ch: "•", dim: true, customColor: demoCyan})
	setCellIfEmpty(grid, y+1, x-1, animCell{ch: "◟", customColor: demoBlue})
	setCell(grid, y+1, x, animCell{ch: "◡", dim: true, customColor: demoWhite})
	setCellIfEmpty(grid, y+1, x+1, animCell{ch: "◞", customColor: demoBlue})
}
