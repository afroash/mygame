package main

import (
	"testing"

	"github.com/afroash/mygame/logic"
)

// Helper function to create a test game instance
func setupTestGame(t *testing.T) *Game {
	game := &Game{
		cursorX: gridSize / 2,
		cursorY: gridSize / 2,
		state:   MainMenu,
		statusMessage: StatusMessage{
			timer:     0,
			isVisible: false,
		},
	}

	// Load a test puzzle
	puzzles, err := logic.LoadPuzzles("sample.txt")
	if err != nil {
		t.Fatalf("Failed to load test puzzles: %v", err)
	}

	puzzle := puzzles[0] // Use first puzzle for consistent testing
	game.logic = &logic.GameLogic{
		Puzzle:    puzzle,
		MoveStack: []logic.Action{},
	}

	return game
}

// Test number validation
func TestIsNumValid(t *testing.T) {
	tests := []struct {
		name     string
		row      int
		col      int
		num      int
		expected bool
	}{
		{"Valid number placement", 0, 0, 1, true},
		{"Invalid row", 0, 1, 6, false},     // 6 already exists in row 0
		{"Invalid column", 1, 0, 6, false},  // 6 already exists in column 0
		{"Invalid subgrid", 1, 1, 6, false}, // 6 already exists in 3x3 grid
		{"Out of range number", 0, 0, 10, false},
		{"Zero is invalid", 0, 0, 0, false},
		{"Negative number", 0, 0, -1, false},
	}

	game := setupTestGame(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := game.isNumValid(tt.row, tt.col, tt.num)
			if result != tt.expected {
				t.Errorf("isNumValid(%d, %d, %d) = %v; want %v",
					tt.row, tt.col, tt.num, result, tt.expected)
			}
		})
	}
}

// Test game state transitions
func TestGameStateTransitions(t *testing.T) {
	game := setupTestGame(t)

	// Test initial state
	if game.state != MainMenu {
		t.Errorf("Initial state = %v; want MainMenu", game.state)
	}

	// Test transition to difficulty menu
	game.state = DifficultyMenu
	game.handleDifficultyMenu()
	if game.state != DifficultyMenu {
		t.Errorf("State after difficulty menu = %v; want DifficultyMenu", game.state)
	}

	// Test difficulty selection
	game.selected = int(Easy)
	game.difficulty = Easy
	game.startGame()
	if game.state != Playing {
		t.Errorf("State after starting game = %v; want Playing", game.state)
	}
}

// Test status message system
func TestStatusMessages(t *testing.T) {
	game := setupTestGame(t)

	// Test showing status message
	testMsg := "Test Message"
	testColor := errorMessage
	testDuration := normalMessageDuration

	game.showStatus(testMsg, testColor, testDuration)

	if !game.statusMessage.isVisible {
		t.Error("Status message not visible after showStatus")
	}
	if game.statusMessage.text != testMsg {
		t.Errorf("Status message text = %v; want %v", game.statusMessage.text, testMsg)
	}
	if game.statusMessage.timer != testDuration {
		t.Errorf("Status message timer = %v; want %v", game.statusMessage.timer, testDuration)
	}

	// Test message timer
	game.updateStatusMessage()
	if game.statusMessage.timer != testDuration-1 {
		t.Error("Status message timer not decremented properly")
	}
}

// Test progress checking
func TestCheckProgress(t *testing.T) {
	game := setupTestGame(t)
	game.state = Playing

	// Test with empty cells
	game.CheckProgress()
	if !game.statusMessage.isVisible {
		t.Error("Progress check should show status message")
	}

	// Test with invalid numbers
	game.logic.Puzzle[0][0] = 1
	game.logic.Puzzle[0][1] = 1 // Same number in row - invalid
	game.CheckProgress()
	if game.statusMessage.color != errorMessage {
		t.Error("Invalid numbers should show error message")
	}
}

// Test win condition
func TestWinCondition(t *testing.T) {
	game := setupTestGame(t)

	// Load a completed valid puzzle
	puzzles, _ := logic.LoadPuzzles("sample.txt")
	game.logic.Puzzle = puzzles[0] // First puzzle is complete

	if !game.logic.IsGridValid() {
		t.Error("Completed puzzle should be valid")
	}
	if !game.logic.IsGridFull() {
		t.Error("Completed puzzle should be full")
	}
}

// Test undo functionality
func TestUndoMove(t *testing.T) {
	game := setupTestGame(t)

	// Make a move
	oldValue := game.logic.Puzzle[0][0]
	newValue := 5
	game.logic.AddMove(0, 0, oldValue, newValue)

	// Test undo
	game.logic.UndoMove()
	if game.logic.Puzzle[0][0] != oldValue {
		t.Errorf("After undo, value = %v; want %v", game.logic.Puzzle[0][0], oldValue)
	}

	// Test undo with empty stack
	game.logic.UndoMove() // Should not panic
}

// Test puzzle difficulty settings
func TestDifficultySettings(t *testing.T) {
	difficulties := []struct {
		level         DifficultyLevel
		expectedEmpty int
	}{
		{Easy, 30},   // 20 + 1*10
		{Medium, 50}, // 20 + 3*10
		{Hard, 70},   // 20 + 5*10
	}

	for _, diff := range difficulties {
		t.Run(diff.level.String(), func(t *testing.T) {
			game := setupTestGame(t)
			game.difficulty = diff.level
			game.startGame()

			emptyCount := 0
			for i := 0; i < gridSize; i++ {
				for j := 0; j < gridSize; j++ {
					if game.logic.Puzzle[i][j] == 0 {
						emptyCount++
					}
				}
			}

			// Allow some variance due to random removal
			if emptyCount < diff.expectedEmpty-5 || emptyCount > diff.expectedEmpty+5 {
				t.Errorf("Difficulty %v: empty cells = %v; want approximately %v",
					diff.level, emptyCount, diff.expectedEmpty)
			}
		})
	}
}

// Helper method to convert DifficultyLevel to string
func (d DifficultyLevel) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	default:
		return "Unknown"
	}
}

// Test game initialization
func TestGameInitialization(t *testing.T) {
	game := NewGame()

	if game.state != MainMenu {
		t.Error("New game should start in MainMenu state")
	}
	if game.cursorX != gridSize/2 || game.cursorY != gridSize/2 {
		t.Error("Cursor should be initialized to center of grid")
	}
	if game.drawer == nil {
		t.Error("DrawHandler should be initialized")
	}
	if game.shoudlExit {
		t.Error("New game should not be in exit state")
	}
}
