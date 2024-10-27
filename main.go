package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/afroash/mygame/logic"

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

const (
	normalFontSize = 12
	menuFontSize   = 24
	diffFontSize   = 18
)

type GameState int

const (
	MainMenu GameState = iota
	DifficultyMenu
	Playing
)

type DifficultyLevel int

const (
	Easy DifficultyLevel = iota
	Medium
	Hard
)

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
	cursorX    int // X position of the game box
	cursorY    int // Y position of the game box
	Puzzle     *logic.GameLogic
	logic      *logic.GameLogic
	state      GameState
	difficulty DifficultyLevel
	selected   int
}

func (g *Game) Update() error {

	switch g.state {
	case MainMenu:
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.selected = (g.selected + 1) % 3 //Allows to cylce through the menu options 3 times.
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.selected = (g.selected + 2) % 3 //Allows to cylce through the menu options 3 times.
		} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			switch g.selected {
			case 0: // Starts a new game
				g.startGame()
			case 1: // Shows the difficulty menu
				g.state = DifficultyMenu
				g.selected = 0
			case 2: // Exits the game
				ebiten.Termination.Error()

			}
		}
	case DifficultyMenu:
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.selected = (g.selected + 1) % 3 //Allows to cylce through the menu options 3 times.
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.selected = (g.selected + 2) % 3 //Allows to cylce through the menu options 3 times.
		} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			switch g.selected {
			case 0: // Easy difficulty
				g.difficulty = Easy
				g.state = Playing
			case 1: // Medium difficulty
				g.difficulty = Medium
				g.state = Playing
			case 2: // Hard difficulty
				g.difficulty = Hard
				g.state = Playing
			}
		}
	case Playing:
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
						g.logic.MoveStack = append(g.logic.MoveStack, logic.Action{
							Row:      g.cursorY,
							Col:      g.cursorX,
							OldValue: 0,
							NewValue: num,
						})

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
	}

	return nil
}

// startGame will start a new game
func (g *Game) startGame() {
	// Load the puzzles from the file
	puzzles, err := logic.LoadPuzzles("sample.txt")
	if err != nil {
		log.Fatalf("Error loading puzzles: %v", err)
	}

	randomPuzzle := logic.GetRandomPuzzle(puzzles)
	logic.ShuffleAsh(&randomPuzzle)
	// Remove numbers from the puzzle based on the difficulty level
	switch g.difficulty {
	case Easy:
		logic.RemoveNumbersFromGrid(&randomPuzzle, 1)
	case Medium:
		logic.RemoveNumbersFromGrid(&randomPuzzle, 3)
	case Hard:
		logic.RemoveNumbersFromGrid(&randomPuzzle, 5)
	}

	// Set the puzzle to the game logic
	g.logic = &logic.GameLogic{
		Puzzle:    randomPuzzle,
		MoveStack: []logic.Action{},
	}
	g.state = Playing
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

// drawMainMenu will draw the main menu
func (g *Game) drawMainMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255}) // White background
	options := []string{"New Game", "Difficulty", "Exit"}
	startX := 100
	startY := 150
	lineSpacing := 30

	for i, option := range options {

		if i == g.selected {
			vector.DrawFilledRect(screen, float32(startX-10), float32(startY+i*lineSpacing-10), float32(200), float32(30), color.RGBA{0, 0, 255, 255}, false)

		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(startY+(i*lineSpacing)))
		//op.LineSpacing = float64(menuFontSize * 20)

		op.ColorScale.ScaleWithColor(color.Black)

		text.Draw(screen, option, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   menuFontSize,
		}, op)
	}
}

// drawDifficultyMenu will draw the difficulty menu
func (g *Game) drawDifficultyMenu(screen *ebiten.Image) {
	startX := 100
	startY := 150
	lineSpacing := 10
	screen.Fill(color.RGBA{255, 255, 255, 255}) // White background
	diffs := []string{"Easy", "medium", "Hard"}
	for i, diff := range diffs {
		if i == g.selected {
			vector.DrawFilledRect(screen, float32(startX-10), float32(startY+i*lineSpacing-10), float32(200), float32(30), color.RGBA{0, 0, 255, 255}, false)
		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(startY+(i*lineSpacing)))
		op.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, diff, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   diffFontSize,
		}, op)
	}
}

// Draw will draw a 9x9 grid.
func (g *Game) Draw(screen *ebiten.Image) {

	// Fill the screen with white background
	screen.Fill(color.RGBA{255, 255, 255, 255})
	switch g.state {
	case MainMenu:
		g.drawMainMenu(screen)
	case DifficultyMenu:
		g.drawDifficultyMenu(screen)
	case Playing:
		g.DrawGrid(screen)
		g.DrawNumbers(screen)
	}
}

// Layout sets the logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the constant screen width and height
	return screenWidth, screenHeight
}

func main() {

	// Create a new game instance
	game := &Game{
		//Puzzle:  randomPuzzle,
		cursorX: gridSize / 2,
		cursorY: gridSize / 2,
		state:   MainMenu, // Start with the main menu
	}

	// Run the game (this will open a window and start rendering)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sudoku BY Ash!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
