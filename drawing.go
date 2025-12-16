package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DrawHandler struct {
	game            *Game
	fontSource      *text.GoTextFaceSource
	screenWidth     int
	screenHeight    int
	gridSize        int
	cellSize        int
	gridTop         int
	statusTop       int
	statusBarHeight int
}

// NewDrawHandler creates a new DrawHandler instance
func NewDrawHandler(game *Game, fontSource *text.GoTextFaceSource) *DrawHandler {
	return &DrawHandler{
		game:            game,
		fontSource:      fontSource,
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		gridSize:        gridSize,
		cellSize:        cellSize,
		gridTop:         gridTop,
		statusTop:       statusTop,
		statusBarHeight: statusBarHeight,
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
			// Add a title at the top
			titleOp := &text.DrawOptions{}
			titleOp.GeoM.Translate(float64(d.screenWidth/2), float64(25))
			titleOp.ColorScale.ScaleWithColor(color.Black)
			titleOp.PrimaryAlign = text.AlignCenter
			titleOp.SecondaryAlign = text.AlignCenter

			text.Draw(screen, "Sudoku", &text.GoTextFace{
				Source: d.fontSource,
				Size:   menuFontSize,
			}, titleOp)

			d.DrawGrid(screen)
			d.DrawNumbers(screen)
			d.drawStatusBar(screen)
			d.drawGameMessages(screen)
		}
	}
}

// DrawGrid draws the 9x9 grid.
func (d *DrawHandler) DrawGrid(screen *ebiten.Image) {
	lineColor := color.RGBA{0, 0, 0, 255} // Black

	// Draw the lines of the grid
	for i := 0; i <= d.gridSize; i++ {
		thickness := float32(1.0)

		// Thicker lines for the 3x3 subgrids
		if i%3 == 0 {
			thickness = 3.0
		}

		// Vertical lines
		x := float32(i * d.cellSize)
		vector.StrokeLine(
			screen,
			x,
			float32(d.gridTop),
			x,
			float32(d.gridTop+(d.gridSize*d.cellSize)), // Fix grid height calculation
			thickness,
			lineColor,
			false,
		)

		// Horizontal lines
		y := float32(d.gridTop + i*d.cellSize)
		vector.StrokeLine(
			screen,
			0,
			y,
			float32(d.gridSize*d.cellSize), // Fix grid width calculation
			y,
			thickness,
			lineColor,
			false,
		)
	}

	// Highlight the active cell
	x := float32(d.game.cursorX * d.cellSize)
	y := float32(d.gridTop + d.game.cursorY*d.cellSize)
	highlightColor := color.RGBA{255, 0, 0, 255}
	if d.game.specialEnterMode {
		highlightColor = color.RGBA{0, 255, 0, 255}
	}
	if !d.game.specialEnterMode {
		highlightColor = color.RGBA{255, 0, 0, 255}
	}
	vector.StrokeRect(
		screen,
		x,
		y,
		float32(d.cellSize),
		float32(d.cellSize),
		float32(2),
		highlightColor,
		false,
	)
}

// DrawPencilMarks draws pencil marks in the corners of a cell
func (d *DrawHandler) DrawPencilMarks(screen *ebiten.Image, row, col int, marks map[int]bool) {
	if len(marks) == 0 {
		return
	}

	// Convert map to sorted slice for consistent corner assignment
	// Iterate 1-9 to get sorted order
	var sortedMarks []int
	for num := 1; num <= 9; num++ {
		if marks[num] {
			sortedMarks = append(sortedMarks, num)
		}
	}

	// Limit to 4 marks (one per corner)
	if len(sortedMarks) > 4 {
		sortedMarks = sortedMarks[:4]
	}

	// Corner positions with padding
	padding := 7
	cellX := col * d.cellSize
	cellY := d.gridTop + row*d.cellSize

	// Define corner positions and alignments
	type cornerInfo struct {
		x, y           int
		primaryAlign   text.Align
		secondaryAlign text.Align
	}

	corners := []cornerInfo{
		{cellX + padding, cellY + padding, text.AlignStart, text.AlignStart},                       // Top-left
		{cellX + d.cellSize - padding, cellY + padding, text.AlignEnd, text.AlignStart},            // Top-right
		{cellX + padding, cellY + d.cellSize - padding, text.AlignStart, text.AlignEnd},            // Bottom-left
		{cellX + d.cellSize - padding, cellY + d.cellSize - padding, text.AlignEnd, text.AlignEnd}, // Bottom-right
	}

	// Draw each mark in its corner
	pencilFontSize := float64(normalFontSize - 4)
	pencilColor := color.RGBA{150, 150, 150, 255} // Gray color for pencil marks

	for i, num := range sortedMarks {
		if i >= len(corners) {
			break
		}

		corner := corners[i]
		numStr := string(rune(num + '0'))

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(corner.x), float64(corner.y))
		op.ColorScale.ScaleWithColor(pencilColor)
		op.PrimaryAlign = corner.primaryAlign
		op.SecondaryAlign = corner.secondaryAlign

		text.Draw(screen, numStr, &text.GoTextFace{
			Source: d.fontSource,
			Size:   pencilFontSize,
		}, op)
	}
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
				y := d.gridTop + row*d.cellSize + d.cellSize/2

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
				
				// Draw pencil marks for filled cells when in help mode
				if d.game.specialEnterMode && len(d.game.pencilMarks[row][col]) > 0 {
					d.DrawPencilMarks(screen, row, col, d.game.pencilMarks[row][col])
				}
			} else {
				// Draw pencil marks for empty cells when in help mode
				if d.game.specialEnterMode && len(d.game.pencilMarks[row][col]) > 0 {
					d.DrawPencilMarks(screen, row, col, d.game.pencilMarks[row][col])
				}
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
	// Draw status area background
	vector.DrawFilledRect(
		screen,
		0,
		float32(d.statusTop),
		float32(d.screenWidth),
		float32(d.screenHeight-d.statusTop),
		color.RGBA{240, 240, 240, 255}, // Light gray background
		false,
	)

	// Draw help mode status
	modeText := "Help Mode: OFF"
	modeColor := color.RGBA{100, 100, 100, 255}
	if d.game.specialEnterMode {
		modeText = "Help Mode: ON"
		modeColor = color.RGBA{0, 150, 0, 255} // Green when active
	}

	modeOp := &text.DrawOptions{}
	modeOp.GeoM.Translate(float64(10), float64(d.statusTop+15))
	modeOp.ColorScale.ScaleWithColor(modeColor)
	modeOp.PrimaryAlign = text.AlignStart
	modeOp.SecondaryAlign = text.AlignStart

	text.Draw(screen, modeText, &text.GoTextFace{
		Source: d.fontSource,
		Size:   normalFontSize,
	}, modeOp)

	// Draw help text
	helpText := "H: Help Mode | N: Normal | P: Check Progress | Z/Backspace: Undo | ESC: Menu"
	helpOp := &text.DrawOptions{}
	helpOp.GeoM.Translate(float64(d.screenWidth/2), float64(d.screenHeight-20))
	helpOp.ColorScale.ScaleWithColor(color.RGBA{100, 100, 100, 255})
	helpOp.PrimaryAlign = text.AlignCenter
	helpOp.SecondaryAlign = text.AlignEnd

	text.Draw(screen, helpText, &text.GoTextFace{
		Source: d.fontSource,
		Size:   normalFontSize,
	}, helpOp)

	// Draw status message if visible
	if d.game.statusMessage.isVisible {
		msg := d.game.statusMessage

		// Draw status message background
		vector.DrawFilledRect(
			screen,
			0,
			float32(d.statusTop),
			float32(d.screenWidth),
			float32(d.statusBarHeight),
			color.RGBA{0, 0, 0, 180},
			false,
		)

		// Draw message text
		op := &text.DrawOptions{}
		op.GeoM.Translate(
			float64(d.screenWidth/2),
			float64(d.statusTop+d.statusBarHeight/2),
		)
		op.ColorScale.ScaleWithColor(msg.color)
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter

		text.Draw(screen, msg.text, &text.GoTextFace{
			Source: d.fontSource,
			Size:   normalFontSize + 2,
		}, op)
	}
}
