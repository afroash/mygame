package drawings

import (
	"bytes"
	"image/color"
	"log"

	"github.com/afroash/mygame/logic"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const normalFontSize = 12

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

const (
	screenWidth  = 450 // Width of the screen
	screenHeight = 450 // Height of the screen
	gridSize     = 9   // Size of the grid
	cellSize     = 50  // Size of each cell
)

type Game struct {
	cursorX int // X position of the game box
	cursorY int // Y position of the game box
	logic   *logic.GameLogic
}

// DrawGrid draws a 9x9 Sudoku grid with thicker lines for 3x3 subgrids
func (g *Game) DrawGrid(screen *ebiten.Image) {
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
	x := float32(g.cursorX * cellSize)
	y := float32(g.cursorY * cellSize)
	vector.StrokeRect(screen, x, y, float32(cellSize), float32(cellSize), float32(2), color.RGBA{255, 0, 0, 255}, false)
}

// DrawNumbers will draw the numbers on the grid
func (g *Game) DrawNumbers(screen *ebiten.Image) {
	// Draw the numbers inside the grid
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if g.logic.Puzzle[row][col] != 0 {
				// Set the position of the number (center it in the cell)
				x := col*cellSize + cellSize/3
				y := row*cellSize + 2*cellSize/3

				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(x), float64(y))
				op.ColorScale.ScaleWithColor(color.Black)
				op.PrimaryAlign = text.AlignCenter
				op.SecondaryAlign = text.AlignCenter

				// Draw the number using basic font with nil options
				numStr := string(rune(g.logic.Puzzle[row][col] + '0')) // Convert int to string
				//face := basicfont.Face7x13
				text.Draw(screen, numStr, &text.GoTextFace{
					Source: mplusFaceSource,
					Size:   normalFontSize,
				}, op)
			}
		}
	}
}
