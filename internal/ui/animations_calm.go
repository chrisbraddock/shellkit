package ui

import "math"

func renderAuroraDrift(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	maxRow := float64(maxInt(rows-1, 1))
	phase := float64(frame) * 0.032
	for col := 0; col < w; col++ {
		bandA := (0.22 + 0.42*(math.Sin(float64(col)*0.17+phase)+1)/2) * maxRow
		bandB := (0.38 + 0.36*(math.Sin(float64(col)*0.11-phase*0.7)+1)/2) * maxRow

		for row := 0; row < rows; row++ {
			distA := math.Abs(float64(row) - bandA)
			distB := math.Abs(float64(row) - bandB)

			intensity := math.Max(0, 1.10-distA*0.80) + math.Max(0, 0.95-distB*0.75)
			if intensity < 0.22 {
				continue
			}

			cell := animCell{ch: "░", dim: intensity < 0.65, customColor: demoBlue}
			if intensity > 0.55 {
				cell.ch = "▒"
				cell.customColor = demoCyan
			}
			if intensity > 1.15 {
				cell.ch = "▓"
				cell.bold = true
				cell.customColor = demoViolet
			}
			setCellIfEmpty(grid, row, col, cell)
		}

		if rows > 1 && (col+frame/9)%23 == 0 {
			setCellIfEmpty(grid, 0, col, animCell{ch: "·", dim: true, customColor: demoWhite})
		}
	}
}

func renderTidalLines(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	amp := math.Max(0.8, float64(rows-1)/3.0)
	center := float64(rows-1) / 2

	for col := 0; col < w; col++ {
		x := float64(col)
		top := int(math.Round(center - 0.7 + math.Sin(x*0.18+float64(frame)*0.034)*amp*0.85))
		mid := int(math.Round(center + math.Sin(x*0.13-float64(frame)*0.028)*amp*0.55))
		low := int(math.Round(center + 0.9 + math.Sin(x*0.09+float64(frame)*0.021+1.6)*amp*0.65))

		setCell(grid, clampInt(top, 0, rows-1), col, animCell{ch: "-", dim: true, customColor: demoBlue})
		setCellIfEmpty(grid, clampInt(mid, 0, rows-1), col, animCell{ch: "~", customColor: demoCyan})
		setCellIfEmpty(grid, clampInt(low, 0, rows-1), col, animCell{ch: "_", dim: true, customColor: demoWhite})

		if rows > 2 && (col+frame/5)%19 == 0 {
			setCellIfEmpty(grid, clampInt(low-1, 0, rows-1), col, animCell{ch: "·", dim: true, customColor: demoWhite})
		}
	}
}

func renderMoonDrift(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	moonRow := clampInt(rows/3-1, 0, rows-1)
	moonX := clampInt(w/2+int(math.Sin(float64(frame)*0.018)*float64(maxInt(2, w/5))), 1, maxInt(1, w-2))

	setCell(grid, moonRow, moonX, animCell{ch: "O", bold: true, customColor: demoWhite})
	setCellIfEmpty(grid, moonRow, moonX-1, animCell{ch: "(", dim: true, customColor: demoBlue})
	setCellIfEmpty(grid, moonRow+1, moonX, animCell{ch: "·", dim: true, customColor: demoCyan})

	starRows := maxInt(1, rows/2)
	for col := 0; col < w; col++ {
		if (col*7+frame/5)%29 == 0 {
			row := (col + frame/17) % starRows
			setCellIfEmpty(grid, row, col, animCell{ch: "·", dim: true, customColor: demoWhite})
		}
	}

	cloudBase := clampInt(rows-1, 0, rows-1)
	for band := 0; band < 2; band++ {
		row := clampInt(cloudBase-(band+1), 0, rows-1)
		span := w/3 + 6
		drift := fract(float64(frame)/(52.0+float64(band)*26.0) + float64(band)*0.17)
		start := int(math.Round((1.0-drift)*float64(w+span))) - span
		for step := 0; step < w/3+6; step++ {
			col := start + step
			if col < 0 || col >= w {
				continue
			}
			ch := "~"
			if step%4 == 1 {
				ch = "_"
			}
			setCellIfEmpty(grid, row, col, animCell{ch: ch, dim: band == 1, customColor: demoBlue})
		}
	}
}

func renderSoftRain(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	drops := maxInt(10, w/3)
	for i := 0; i < drops; i++ {
		baseX := pseudoRand(i*137+19) % w
		seed := float64(pseudoRand(i*83+7)%100) * 0.01
		progress := fract(float64(frame)*0.03 + seed)
		rowFloat := progress*float64(rows+4) - 2
		row := int(math.Floor(rowFloat))
		drift := 2 + pseudoRand(i*61+13)%4
		x := (baseX + int(progress*float64(drift))) % w

		if row >= 0 && row < rows {
			setCell(grid, row, x, animCell{ch: "╲", dim: true, customColor: demoCyan})
			setCellIfEmpty(grid, row-1, x, animCell{ch: "·", dim: true, customColor: demoWhite})
			setCellIfEmpty(grid, row+1, x+1, animCell{ch: "·", dim: true, customColor: demoBlue})
		}

		if rows > 1 && progress > 0.82 {
			rippleRow := rows - 1
			rippleX := clampInt(x, 0, w-1)
			ripple := "_"
			if (i+frame/3)%2 == 0 {
				ripple = "."
			}
			setCellIfEmpty(grid, rippleRow, rippleX, animCell{ch: ripple, dim: true, customColor: demoBlue})
			setCellIfEmpty(grid, rippleRow, rippleX-1, animCell{ch: "·", dim: true, customColor: demoWhite})
			setCellIfEmpty(grid, rippleRow, rippleX+1, animCell{ch: "·", dim: true, customColor: demoWhite})
		}
	}

	renderRainLightning(grid, w, frame)
}

func renderLanternDrift(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows == 0 || w <= 0 {
		return
	}

	count := maxInt(3, w/16)
	palette := []animCell{
		{ch: "◌", customColor: demoAmber},
		{ch: "○", customColor: demoPink},
		{ch: "◍", customColor: demoCyan},
	}

	for i := 0; i < count; i++ {
		cycle := rows*12 + 24
		progress := float64((frame*2+i*19)%cycle) / float64(cycle)
		x := (i*w)/count + int(math.Sin(progress*math.Pi*2+float64(i))*2)
		x = clampInt(x, 0, w-1)
		row := rows - 1 - int(progress*float64(rows+3))
		if row >= 0 && row < rows {
			cell := palette[i%len(palette)]
			cell.bold = progress > 0.45
			setCell(grid, row, x, cell)
		}

		for trail := 1; trail <= 2; trail++ {
			trailRow := row + trail
			if trailRow >= 0 && trailRow < rows {
				setCellIfEmpty(grid, trailRow, x, animCell{ch: "·", dim: true, customColor: demoWhite})
			}
		}
	}
}

func renderRainLightning(grid [][]animCell, w, frame int) {
	rows := len(grid)
	if rows < 3 || w < 6 {
		return
	}

	segmentLen := 240
	segment := frame / segmentLen
	phase := frame % segmentLen
	strikeStart := 18 + pseudoRand(segment*149+w*17+rows*13)%maxInt(24, segmentLen-54)
	if phase < strikeStart || phase >= strikeStart+12 {
		return
	}

	flashPhase := phase - strikeStart
	flash := math.Sin(float64(flashPhase) / 11.0 * math.Pi)
	if flash <= 0 {
		return
	}

	boltX := 2 + pseudoRand(segment*211+37)%maxInt(1, w-4)
	boltHeight := minInt(rows-1, 2+rows/2+pseudoRand(segment*53+11)%maxInt(1, rows/3+1))

	for step := 0; step < boltHeight; step++ {
		x := clampInt(boltX+int(math.Round(math.Sin(float64(step)*1.2+float64(segment))*1.4)), 0, w-1)
		y := step
		ch := "|"
		if step%2 == 0 {
			ch = "\\"
		} else if step%3 == 0 {
			ch = "/"
		}
		cell := animCell{ch: ch, bold: flash > 0.45, customColor: demoWhite}
		if flash < 0.45 {
			cell.customColor = demoCyan
		}
		setCell(grid, y, x, cell)
		setCellIfEmpty(grid, y, x-1, animCell{ch: ".", dim: true, customColor: demoWhite})
		setCellIfEmpty(grid, y, x+1, animCell{ch: ".", dim: true, customColor: demoCyan})
	}

	if flash > 0.28 {
		for col := 0; col < w; col++ {
			if (col+segment)%4 == 0 {
				setCell(grid, 0, col, animCell{ch: ".", bold: flash > 0.55, customColor: demoWhite})
			}
			if rows > 1 && (col+segment)%7 == 0 {
				setCellIfEmpty(grid, 1, col, animCell{ch: ":", dim: true, customColor: demoBlue})
			}
		}
	}
}
