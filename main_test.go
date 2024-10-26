package main

import (
	"testing"

	"github.com/afroash/ashlog"
	"github.com/afroash/mygame/logic"
)

// func to test row validation
func TestIsNumberValidRow(t *testing.T) {
	g := &Game{}

	//test case 1 - number is already present in the row, lets use 6
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
	g.logic.Puzzle = randomPuzzle
	if g.isNumValid(0, 5, 6) != false {
		t.Errorf("Expected false, got true. Number 6 is already present in the row")
	}
}

// func to test column validation
func TestIsNumberValidColumn(t *testing.T) {
	g := &Game{}

	//test case 1 - number is already present in the column, lets use 5
	//Load puzzle from samples.
	puzzles, err := logic.LoadPuzzles("sample.txt")
	if err != nil {
		t.Errorf("Error loading puzzles: %v", err)
	}

	//Select a random puzzle
	randomPuzzle := logic.GetRandomPuzzle(puzzles)

	//shuffle the puzzle
	logic.ShuffleAsh(&randomPuzzle)

	//remove some numbers from the puzzle
	logic.RemoveNumbersFromGrid(&randomPuzzle, 1)
	g.logic.Puzzle = randomPuzzle
	if g.isNumValid(5, 5, 5) != false {
		t.Errorf("Expected false, got true. Number 5 is already present in the column")
	}
}

// func to test subgrid validation
func TestIsNumberValidSubgrid(t *testing.T) {
	g := &Game{}

	//test case 1 - number is already present in the subgrid, lets use 9
	//Load puzzle from samples.
	puzzles, err := logic.LoadPuzzles("sample.txt")
	if err != nil {
		t.Errorf("Error loading puzzles: %v", err)
	}

	//Select a random puzzle
	randomPuzzle := logic.GetRandomPuzzle(puzzles)

	//shuffle the puzzle
	logic.ShuffleAsh(&randomPuzzle)

	//remove some numbers from the puzzle
	logic.RemoveNumbersFromGrid(&randomPuzzle, 1)
	g.logic.Puzzle = randomPuzzle
	if g.isNumValid(1, 3, 4) != false {
		t.Errorf("Expected false, got true. Number 9 is already present in the subgrid")
	}
}
