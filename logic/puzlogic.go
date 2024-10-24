package logic

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Action struct {
	Row, Col int
	OldValue int
	NewValue int
}

type Puzzle [9][9]int

type GameLogic struct {
	Puzzle    Puzzle
	MoveStack []Action
}

// Function to load puzzles from the text file
func LoadPuzzles(filename string) ([]Puzzle, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var puzzles []Puzzle
	var currentPuzzle [9][9]int
	scanner := bufio.NewScanner(file)
	row := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" { // Blank line indicates end of a puzzle
			if row == 9 { // Only add if we've completed reading 9 rows
				puzzles = append(puzzles, currentPuzzle)
				row = 0
			}
			continue
		}

		if len(line) != 9 {
			return nil, fmt.Errorf("invalid row length in puzzle: %v. The row has %v", line, len(line))
		}

		for col, char := range line {
			if char < '0' || char > '9' {
				return nil, fmt.Errorf("invalid character in puzzle: %v", char)
			}
			currentPuzzle[row][col] = int(char - '0')
		}

		row++
	}

	// Add the last puzzle if no trailing blank line
	if row == 9 {
		puzzles = append(puzzles, currentPuzzle)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return puzzles, nil
}

// Function to select a random puzzle from the loaded puzzles
func GetRandomPuzzle(puzzles []Puzzle) [9][9]int {
	//rand.Seed(time.Now().UnixNano())
	return puzzles[rand.Intn(len(puzzles))]
}

// Function to shuffle the puzzle
func ShuffleAsh(grid *[9][9]int) {
	numbers := rand.Perm(9)
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] != 0 {
				grid[row][col] = numbers[grid[row][col]-1] + 1
			}
		}
	}

	//Shuffle the rows withing each block.
	for i := 0; i < 9; i += 3 {
		shuffleRowsInBlock(grid, i)
	}

	//Shuffle the columns withing each block.
	for i := 0; i < 9; i += 3 {
		shuffleColsInBlock(grid, i)
	}

}

// Function to shuffle the rows within a block
func shuffleRowsInBlock(grid *[9][9]int, blockStart int) {
	rows := rand.Perm(3)
	for i := 0; i < 3; i++ {
		temp := grid[blockStart+i]
		grid[blockStart+i] = grid[blockStart+rows[i]]
		grid[blockStart+rows[i]] = temp
	}
}

// Function to shuffle the columns within a block
func shuffleColsInBlock(grid *[9][9]int, blockStart int) {
	for i := 0; i < 3; i++ {
		cols := rand.Perm(3)
		for j := 0; j < 9; j++ {
			temp := grid[j][blockStart+i]
			grid[j][blockStart+i] = grid[j][blockStart+cols[i]]
			grid[j][blockStart+cols[i]] = temp
		}
	}
}

// Remove numbers to make the puzzle playable
func RemoveNumbersFromGrid(grid *[9][9]int, difficulty int) {
	blanks := 20 + difficulty*10 // Control how many numbers to remove based on difficulty
	for blanks > 0 {
		row := rand.Intn(9)
		col := rand.Intn(9)
		if grid[row][col] != 0 {
			grid[row][col] = 0
			blanks--
		}
	}
}

// Add moves to the stack
func (g *GameLogic) AddMove(row, col, oldValue, newValue int) {
	oldmove := g.Puzzle[row][col]
	action := Action{
		Row:      row,
		Col:      col,
		OldValue: oldmove,
		NewValue: newValue,
	}
	g.MoveStack = append(g.MoveStack, action)
	g.Puzzle[row][col] = newValue
}

// Undo the last move
func (g *GameLogic) UndoMove() {
	if len(g.MoveStack) == 0 {
		return
	}
	lastMove := g.MoveStack[len(g.MoveStack)-1]
	g.MoveStack = g.MoveStack[:len(g.MoveStack)-1]
	g.Puzzle[lastMove.Row][lastMove.Col] = lastMove.OldValue
}
