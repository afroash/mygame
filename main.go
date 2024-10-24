package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/afroash/mygame/logic"

	"github.com/afroash/ashlog"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 450 // Width of the screen
	screenHeight = 450 // Height of the screen
	gridSize     = 9   // Size of the grid
	cellSize     = 50  // Size of each cell
)

const normalFontSize = 12

// // Sample grid data, 0 means empty cell
// var sampleGrid = [9][9]int{
// 	{5, 3, 0, 0, 7, 0, 0, 0, 0},
// 	{6, 0, 0, 1, 9, 5, 0, 0, 0},
// 	{0, 9, 8, 0, 0, 0, 0, 6, 0},
// 	{8, 0, 0, 0, 6, 0, 0, 0, 3},
// 	{4, 0, 0, 8, 0, 3, 0, 0, 1},
// 	{7, 0, 0, 0, 2, 0, 0, 0, 6},
// 	{0, 6, 0, 0, 0, 0, 2, 8, 0},
// 	{0, 0, 0, 4, 1, 9, 0, 0, 5},
// 	{0, 0, 0, 0, 8, 0, 0, 7, 9},
// }

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

type Game struct {
	cursorX int // X position of the game box
	cursorY int // Y position of the game box
	Puzzle  *logic.GameLogic
	logic   *logic.GameLogic
}

func (g *Game) Update() error {
	// Move the cursor based on the arrow keys or WASD keys
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.cursorY > 0 {
			g.cursorY--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.cursorY < gridSize-1 {
			g.cursorY++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if g.cursorX > 0 {
			g.cursorX--
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.cursorX < gridSize-1 {
			g.cursorX++
		}

	}

	// handle number input for empty cells
	if g.logic.Puzzle[g.cursorY][g.cursorX] == 0 {
		for i := ebiten.Key0; i <= ebiten.Key9; i++ {
			if inpututil.IsKeyJustPressed(i) {
				num := int(i - ebiten.Key0)
				//validate the number
				if g.isNumValid(g.cursorY, g.cursorX, num) {
					g.logic.Puzzle[g.cursorY][g.cursorX] = num
				} else {
					log.Println("Invalid number try again")
				}
			}
		}
	}

	//handle undo input via z or backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.logic.UndoMove()
	}
	return nil
}

// Lets check if the entered number is valid as per Sudoku rules.
func (g *Game) isNumValid(row, col, num int) bool {
	// Check if the number is already present in the row
	for i := 0; i < gridSize; i++ {
		if g.logic.Puzzle[row][i] == num {
			return false
		}
	}

	// Check if the number is already present in the column
	for i := 0; i < gridSize; i++ {
		if g.logic.Puzzle[i][col] == num {
			return false
		}
	}

	// Check if the number is already present in the 3x3 subgrid
	subGridRowSrt := row - row%3
	subGridColSrt := col - col%3
	for r := subGridRowSrt; r < subGridRowSrt; r++ {
		for c := subGridColSrt; c < subGridColSrt+3; c++ {
			if g.logic.Puzzle[r][c] == num {
				return false
			}
		}
	}

	// If the number is not present in the row, column, or subgrid, then it is valid
	return true
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

// Draw will draw a 9x9 grid.
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with white background
	screen.Fill(color.RGBA{255, 255, 255, 255})
	g.DrawGrid(screen)
	g.DrawNumbers(screen)
}

// Layout sets the logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the constant screen width and height
	return screenWidth, screenHeight
}

func main() {
	//Load puzzle from samples.
	puzzles, err := logic.LoadPuzzles("sample.txt")
	if err != nil {
		ashlog.LogFatal("Error loading puzzles", err)
	}

	//Select a random puzzle
	randomPuzzle := logic.GetRandomPuzzle(puzzles)

	//shuffle the puzzle
	//logic.ShuffleAsh(&randomPuzzle)

	//remove some numbers from the puzzle
	logic.RemoveNumbersFromGrid(&randomPuzzle, 1)

	gamelogic := &logic.GameLogic{
		Puzzle:    randomPuzzle,
		MoveStack: []logic.Action{},
	}
	// Create a new game instance
	game := &Game{
		//Puzzle:  randomPuzzle,
		cursorX: gridSize / 2,
		cursorY: gridSize / 2,
		logic:   gamelogic,
	}

	// Run the game (this will open a window and start rendering)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sudoku BY Ash!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
