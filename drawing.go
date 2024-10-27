package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DrawHandler struct {
	game         *Game
	fontSource   *text.GoTextFaceSource
	screenWidth  int
	screenHeight int
	gridSize     int
	cellSize     int
}

// NewDrawHandler creates a new DrawHandler instance
func NewDrawHandler(game *Game, fontSource *text.GoTextFaceSource) *DrawHandler {
	return &DrawHandler{
		game:         game,
		fontSource:   fontSource,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		gridSize:     gridSize,
		cellSize:     cellSize,
	}
}

// Draw handles the main drawing logic based on game state
func (d *DrawHandler) Draw(screen *ebiten.Image) {
	// Fill the screen with white background
	screen.Fill(color.RGBA{255, 255, 255, 255})

	switch d.game.state {
	case MainMenu:
		d.drawMainMenu(screen)
	case DifficultyMenu:
		d.drawDifficultyMenu(screen)
	case Playing:
		d.DrawGrid(screen)
		d.DrawNumbers(screen)
	}
}

// DrawGrid draws the 9x9 grid.
func (d *DrawHandler) DrawGrid(screen *ebiten.Image) {
	// Line color (black)

	lineColor := color.RGBA{0, 0, 0, 255} // Black

	// Draw the lines of the grid
	for i := 0; i <= gridSize; i++ {
		thickness := float32(1.0)

		// Thicker lines for the 3x3 subgrids
		if i%3 == 0 {
			thickness = 3.0
		}

		// Vertical lines
		x := float32(i * cellSize)
		vector.StrokeLine(screen, x, 0, x, float32(screenHeight), thickness, lineColor, false)

		// Horizontal lines
		y := float32(i * cellSize)
		vector.StrokeLine(screen, 0, y, float32(screenWidth), y, thickness, lineColor, false)
	}

	//Highlight the active cell
	x := float32(d.game.cursorX * cellSize)
	y := float32(d.game.cursorY * cellSize)
	vector.StrokeRect(screen, x, y, float32(cellSize), float32(cellSize), float32(2), color.RGBA{255, 0, 0, 255}, false)
}

// DrawNumbers draws the numbers on the grid
func (d *DrawHandler) DrawNumbers(screen *ebiten.Image) {
	for row := 0; row < d.gridSize; row++ {
		for col := 0; col < d.gridSize; col++ {
			if d.game.logic.Puzzle[row][col] != 0 {
				x := col*d.cellSize + d.cellSize/3
				y := row*d.cellSize + 2*d.cellSize/3

				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(x), float64(y))
				op.ColorScale.ScaleWithColor(color.Black)
				op.PrimaryAlign = text.AlignCenter
				op.SecondaryAlign = text.AlignCenter

				numStr := string(rune(d.game.logic.Puzzle[row][col] + '0'))
				text.Draw(screen, numStr, &text.GoTextFace{
					Source: d.fontSource,
					Size:   normalFontSize,
				}, op)
			}
		}
	}
}

// drawMainMenu draws the main menu
func (d *DrawHandler) drawMainMenu(screen *ebiten.Image) {
	options := []string{"New Game", "Difficulty", "Exit"}
	startX := 100
	startY := 150
	lineSpacing := 30

	for i, option := range options {
		if i == d.game.selected {
			vector.DrawFilledRect(
				screen,
				float32(startX-10),
				float32(startY+i*lineSpacing-10),
				float32(200),
				float32(30),
				color.RGBA{0, 0, 255, 255},
				false)
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(startY+(i*lineSpacing)))
		op.ColorScale.ScaleWithColor(color.Black)

		text.Draw(screen, option, &text.GoTextFace{
			Source: d.fontSource,
			Size:   menuFontSize,
		}, op)
	}
}

// drawDifficultyMenu draws the difficulty menu
func (d *DrawHandler) drawDifficultyMenu(screen *ebiten.Image) {
	startX := 100
	startY := 150
	lineSpacing := 10
	diffs := []string{"Easy", "medium", "Harder"}

	for i, diff := range diffs {
		if i == d.game.selected {
			vector.DrawFilledRect(
				screen,
				float32(startX-10),
				float32(startY+i*lineSpacing-10),
				float32(200),
				float32(30),
				color.RGBA{0, 0, 255, 255},
				false)
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(startY+(i*lineSpacing)))
		op.ColorScale.ScaleWithColor(color.Black)

		text.Draw(screen, diff, &text.GoTextFace{
			Source: d.fontSource,
			Size:   diffFontSize,
		}, op)
	}
}
