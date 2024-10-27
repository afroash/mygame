package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/afroash/mygame/logic"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	diffFontSize   = 24
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

type StatusMessage struct {
	text      string
	color     color.RGBA
	timer     int
	isVisible bool
}

// Game struct
type Game struct {
	cursorX        int // X position of the game box
	cursorY        int // Y position of the game box
	Puzzle         *logic.GameLogic
	logic          *logic.GameLogic
	state          GameState
	difficulty     DifficultyLevel
	selected       int
	drawer         *DrawHandler
	shoudlExit     bool
	showWinMessage bool
	messageTimer   int
	statusMessage  StatusMessage
}

func NewGame() *Game {
	// Initialize font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		cursorX:    gridSize / 2,
		cursorY:    gridSize / 2,
		state:      MainMenu,
		shoudlExit: false,
		statusMessage: StatusMessage{
			timer:     0,
			isVisible: false,
		},
	}

	// Initialize the drawer
	game.drawer = NewDrawHandler(game, s)

	return game
}

func (g *Game) showStatus(text string, color color.RGBA) {
	g.statusMessage = StatusMessage{
		text:      text,
		color:     color,
		timer:     120,
		isVisible: true,
	}
}

func (g *Game) updateStatusMessage() {
	if g.statusMessage.isVisible {
		g.statusMessage.timer--
		if g.statusMessage.timer <= 0 {
			g.statusMessage.isVisible = false
		}
	}
}

func (g *Game) Update() error {
	// check if the game should exit
	if g.shoudlExit {
		return ebiten.Termination
	}

	switch g.state {
	case MainMenu:
		g.handleMainMenu()
	case DifficultyMenu:
		g.handleDifficultyMenu()
	case Playing:
		if g.logic != nil {
			g.handlePlayingInput()
		}
	}

	// Global Exit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch g.state {
		case Playing:
			// If in playing state, go back to main menu
			g.state = MainMenu
			g.selected = 0
		case DifficultyMenu:
			// If in difficulty menu, go back to main menu
			g.state = MainMenu
			g.selected = 1
		case MainMenu:
			// If in main menu, exit the game
			g.shoudlExit = true
		}
	}

	return nil
}

func (g *Game) handleMainMenu() {
	// Only process one key press per frame for smoother navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.selected = (g.selected + 1) % 3
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.selected = (g.selected - 1)
		if g.selected < 0 {
			g.selected = 2
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch g.selected {
		case 0: // New Game
			if g.difficulty == 0 {
				// If difficulty hasn't been set, go to difficulty menu first
				g.state = DifficultyMenu
				g.selected = 0
			} else {
				// If difficulty is already set, start the game
				g.startGame()
			}
		case 1: // Difficulty
			g.state = DifficultyMenu
			g.selected = 0
		case 2: // Exit
			g.shoudlExit = true // Exit the game
		}
	}
}

func (g *Game) handleDifficultyMenu() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.selected = (g.selected + 1) % 3
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.selected = (g.selected - 1)
		if g.selected < 0 {
			g.selected = 2
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.difficulty = DifficultyLevel(g.selected)
		g.startGame()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = MainMenu
		g.selected = 1 // Select "Difficulty" option when returning
	}
}

// handlePlayingInput will handle the input when the game is in the Playing state
func (g *Game) handlePlayingInput() {
	//update the status message timer
	g.updateStatusMessage()

	// Handle  progress check [P] key
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.CheckProgress()

	}
	// Move the cursor
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

	// Handle number input
	if g.logic.Puzzle[g.cursorY][g.cursorX] == 0 {
		for i := ebiten.Key0; i <= ebiten.Key9; i++ {
			if inpututil.IsKeyJustPressed(i) {
				num := int(i - ebiten.Key0)
				if g.isNumValid(g.cursorY, g.cursorX, num) {
					g.logic.Puzzle[g.cursorY][g.cursorX] = num
					g.logic.MoveStack = append(g.logic.MoveStack, logic.Action{
						Row:      g.cursorY,
						Col:      g.cursorX,
						OldValue: 0,
						NewValue: num,
					})
				} else {
					// Number is invalid
					fmt.Println("Invalid number")
				}
				// Check game status after each move
				if g.logic.IsGridFull() {
					if g.logic.IsGridValid() {
						g.showWinMessage = true
						g.messageTimer = 180
					}
				}
			}
		}
	}

	// Handle undo
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.logic.UndoMove()
	}

	// Handle win message timer
	if g.showWinMessage {
		g.messageTimer--
		if g.messageTimer <= 0 {
			g.showWinMessage = false
		}
	}

	// Add a check progress button (Space key)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch g.logic.GetGameStatus() {
		case logic.InProgress:
			// show a "Keep going!" message

		case logic.Completed:
			g.showWinMessage = true
			g.messageTimer = 180
		case logic.Invalid:
			// Optionally show an "Something's not right" message
		}
	}
}

// CheckProgress will check the progress of the game
func (g *Game) CheckProgress() {
	if g.logic == nil {
		return
	}
	emptyCount := 0
	invalidCount := 0

	// Count empty cells and check for invalid entries
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if g.logic.Puzzle[i][j] == 0 {
				emptyCount++
			} else if !g.isNumValid(i, j, g.logic.Puzzle[i][j]) {
				invalidCount++
			}
		}
	}

	// Prepare status message
	if invalidCount > 0 {
		g.showStatus(
			fmt.Sprintf("Found %d incorrect numbers", invalidCount),
			color.RGBA{255, 0, 0, 255}, // Red
		)
	} else if emptyCount > 0 {
		g.showStatus(
			fmt.Sprintf("%d cells left to fill", emptyCount),
			color.RGBA{0, 128, 255, 255}, // Blue
		)
	} else if g.logic.IsGridValid() {
		g.showStatus(
			"Puzzle completed correctly!",
			color.RGBA{0, 255, 0, 255}, // Green
		)
		g.showWinMessage = true
		g.messageTimer = 180
	}
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
	if g.logic == nil {
		return false
	}

	// Check row
	for i := 0; i < gridSize; i++ {
		if g.logic.Puzzle[row][i] == num {
			return false
		}
	}

	// Check column
	for i := 0; i < gridSize; i++ {
		if g.logic.Puzzle[i][col] == num {
			return false
		}
	}

	// Check 3x3 subgrid
	subGridRowStart := (row / 3) * 3
	subGridColStart := (col / 3) * 3
	for r := subGridRowStart; r < subGridRowStart+3; r++ {
		for c := subGridColStart; c < subGridColStart+3; c++ {
			if g.logic.Puzzle[r][c] == num {
				return false
			}
		}
	}

	return true
}

// Draw will draw a 9x9 grid.
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawer.Draw(screen)

}

// Layout sets the logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the constant screen width and height
	return screenWidth, screenHeight
}

func main() {

	// Create a new game instance
	game := NewGame()
	// Run the game (this will open a window and start rendering)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sudoku BY Ash!")

	if err := ebiten.RunGame(game); err != nil {
		if err == ebiten.Termination {
			// Clean Exit
			os.Exit(0)
		}

		log.Fatal(err)
	}
}
