package main

import (
	"fmt"
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
		if d.game.logic != nil {
			d.DrawGrid(screen)
			d.DrawNumbers(screen)
			d.drawStatusBar(screen)
			d.drawGameMessages(screen)
		}
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
	if d.game == nil || d.game.logic == nil {
		fmt.Println("Game or logic is nil")
		return
	}
	for row := 0; row < d.gridSize; row++ {
		for col := 0; col < d.gridSize; col++ {
			if d.game.logic.Puzzle[row][col] != 0 {

				//Center the number in the cell
				x := col*d.cellSize + d.cellSize/2
				y := row*d.cellSize + d.cellSize/2

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
	// Center the menu on screen
	startX := screenWidth / 2
	startY := screenHeight / 3
	lineSpacing := 50 // Increased spacing between options
	options := []string{"New Game", "Difficulty", "Exit"}

	// Draw title
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(startX), float64(startY-85))
	titleOp.ColorScale.ScaleWithColor(color.Black)
	titleOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Sudoku by Ash", &text.GoTextFace{
		Source: d.fontSource,
		Size:   menuFontSize + 8,
	}, titleOp)

	for i, option := range options {
		yPos := startY + i*lineSpacing

		// Calculate text metrics for centering
		textWidth := len(option) * diffFontSize / 2 // Approximate width
		rectWidth := float32(textWidth + 40)        // Add padding
		rectHeight := float32(40)                   // Fixed height for selection rectangle

		// Draw selection highlight if this option is selected
		if i == d.game.selected {
			vector.DrawFilledRect(
				screen,
				float32(startX)-rectWidth/2, // Center the rectangle
				float32(yPos)-rectHeight/2,  // Center vertically around text
				rectWidth,
				rectHeight,
				color.RGBA{0, 0, 255, 100}, // Lighter blue for better visibility
				false,
			)

			// Draw border for selected option
			vector.StrokeRect(
				screen,
				float32(startX)-rectWidth/2,
				float32(yPos)-rectHeight/2,
				rectWidth,
				rectHeight,
				2,
				color.RGBA{0, 0, 255, 255},
				false,
			)
		}

		// Draw the menu option text
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(yPos))
		op.ColorScale.ScaleWithColor(color.Black)
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter

		text.Draw(screen, option, &text.GoTextFace{
			Source: d.fontSource,
			Size:   menuFontSize,
		}, op)
	}

	// Draw instructions at the bottom
	instructOp := &text.DrawOptions{}
	instructOp.GeoM.Translate(float64(startX), float64(startY+lineSpacing*4))
	instructOp.ColorScale.ScaleWithColor(color.RGBA{100, 100, 100, 255})
	instructOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Use ↑↓ to select, ENTER to confirm", &text.GoTextFace{
		Source: d.fontSource,
		Size:   normalFontSize,
	}, instructOp)
}

// drawDifficultyMenu method in drawing.go
func (d *DrawHandler) drawDifficultyMenu(screen *ebiten.Image) {
	// Center the menu on screen
	startX := screenWidth / 2
	startY := screenHeight / 3
	lineSpacing := 50 // Increased spacing between options
	diffs := []string{"Easy", "Medium", "Hard"}

	// Draw title
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(startX), float64(startY-85))
	titleOp.ColorScale.ScaleWithColor(color.Black)
	titleOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Select Difficulty", &text.GoTextFace{
		Source: d.fontSource,
		Size:   menuFontSize + 4,
	}, titleOp)

	for i, diff := range diffs {
		yPos := startY + i*lineSpacing

		// Draw selection highlight
		if i == d.game.selected {
			vector.DrawFilledRect(
				screen,
				float32(startX-100),
				float32(yPos-20),
				200,
				40,
				color.RGBA{0, 0, 255, 255},
				false)
		}

		// Draw text
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(yPos))
		op.ColorScale.ScaleWithColor(color.Black)
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter

		text.Draw(screen, diff, &text.GoTextFace{
			Source: d.fontSource,
			Size:   diffFontSize,
		}, op)
	}

	// Draw instruction
	instructOp := &text.DrawOptions{}
	instructOp.GeoM.Translate(float64(startX), float64(startY+lineSpacing*4))
	instructOp.ColorScale.ScaleWithColor(color.RGBA{100, 100, 100, 255})
	instructOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Press ESC to return to main menu", &text.GoTextFace{
		Source: d.fontSource,
		Size:   normalFontSize,
	}, instructOp)
}

func (d *DrawHandler) drawGameMessages(screen *ebiten.Image) {
	if !d.game.showWinMessage {
		return
	}

	// Draw semi-transparent overlay
	vector.DrawFilledRect(
		screen,
		0,
		0,
		float32(d.screenWidth),
		float32(d.screenHeight),
		color.RGBA{0, 0, 0, 180},
		false,
	)

	// Draw win message
	message := "Congratulations! Puzzle Solved!"
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(d.screenWidth/2), float64(d.screenHeight/2))
	op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	text.Draw(screen, message, &text.GoTextFace{
		Source: d.fontSource,
		Size:   menuFontSize,
	}, op)

	// Draw sub-message
	subMessage := "Press ESC for menu, ENTER for new game"
	subOp := &text.DrawOptions{}
	subOp.GeoM.Translate(float64(d.screenWidth/2), float64(d.screenHeight/2)+40)
	subOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 200})
	subOp.PrimaryAlign = text.AlignCenter
	subOp.SecondaryAlign = text.AlignCenter

	text.Draw(screen, subMessage, &text.GoTextFace{
		Source: d.fontSource,
		Size:   normalFontSize,
	}, subOp)
}

func (d *DrawHandler) drawStatusBar(screen *ebiten.Image) {
	if d.game.state != Playing {
		return
	}

	// Draw help text for progress check
	helpText := "Press P to check progress"
	helpOp := &text.DrawOptions{}
	helpOp.GeoM.Translate(float64(10), float64(d.screenHeight-20))
	helpOp.ColorScale.ScaleWithColor(color.RGBA{100, 100, 100, 255})

	text.Draw(screen, helpText, &text.GoTextFace{
		Source: d.fontSource,
		Size:   normalFontSize,
	}, helpOp)

	// Draw status message if visible
	if d.game.statusMessage.isVisible {
		msg := d.game.statusMessage

		// Draw background for status message
		vector.DrawFilledRect(
			screen,
			0,
			float32(d.screenHeight-50),
			float32(d.screenWidth),
			40,
			color.RGBA{0, 0, 0, 180},
			false,
		)

		// Draw status message
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(d.screenWidth/2), float64(d.screenHeight-30))
		op.ColorScale.ScaleWithColor(msg.color)
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter

		text.Draw(screen, msg.text, &text.GoTextFace{
			Source: d.fontSource,
			Size:   normalFontSize + 2,
		}, op)
	}
}
