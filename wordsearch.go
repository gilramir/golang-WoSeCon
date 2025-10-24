// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"math/rand"
	"slices"
)

// The constructed word search puzzle
type WordSearch struct {
	NumCols int
	NumRows int

	// The solution, as words and placments (location+direction)
	WordPlacements map[string]WordPlacement

	// The solution, as runes
	// First dimension is y (rows)
	// Second dimension is x (column)
	SolutionRows [][]rune

	// The puzzle, as runes
	// Every cell is filled. It has the solution
	// an the filler, together.
	// First dimension is y (rows)
	// Second dimension is x (column)
	PuzzleRows [][]rune
}

type WordPlacement struct {
	Col       int
	Row       int
	Direction Direction
}

func (s WordPlacement) DirectionString() string {
	return DirectionString(s.Direction)
}

// Take a successful wordSearchConstructor,
// and return a WordSearch result which will be sent back to the caller
func (s *wordSearchConstructor) translateToWordSearch() *WordSearch {

	ws := &WordSearch{
		NumCols:        s.numCols,
		NumRows:        s.numRows,
		WordPlacements: make(map[string]WordPlacement),
		// Allocate the rows
		SolutionRows: make([][]rune, s.numRows),
		PuzzleRows:   make([][]rune, s.numRows),
	}
	// Allocate the columns, filled with SPACE characteers
	for row := 0; row < s.numRows; row++ {
		ws.SolutionRows[row] = slices.Repeat([]rune{' '}, s.numCols)
		ws.PuzzleRows[row] = slices.Repeat([]rune{' '}, s.numCols)
	}

	// Fill in the solution and puzzle
	for _, wordInfo := range s.wordInfos {
		dl := wordInfo.placement
		ws.WordPlacements[wordInfo.text] = WordPlacement{
			Col:       dl.col,
			Row:       dl.row,
			Direction: dl.direction,
		}

		var colAdj int
		var rowAdj int

		// Find endRow
		if dl.direction&GoesDownward > 0 {
			rowAdj = 1
		} else if dl.direction&GoesUpward > 0 {
			rowAdj = -1
		} // else horizontal

		// Find endCol
		if dl.direction&GoesLTR != 0 {
			colAdj = 1
		} else if dl.direction&GoesRTL != 0 {
			colAdj = -1
		} // else vetical

		var col = dl.col
		var row = dl.row
		// Iterate rune by rune in the string
		for _, r := range wordInfo.text {
			ws.SolutionRows[row][col] = r
			ws.PuzzleRows[row][col] = r
			col += colAdj
			row += rowAdj
		}
	}

	// Add the filler
	if len(s.fillerWeights) > 0 {
		s.applyWeightedFiller(ws)
	} else if len(s.fillerRunes) > 0 {
		s.applyUniformFiller(ws)
	}

	return ws
}

func (s *wordSearchConstructor) applyUniformFiller(ws *WordSearch) {
	for col := 0; col < ws.NumCols; col++ {
		for row := 0; row < ws.NumRows; row++ {
			if s.isCellAvailable(col, row) {
				// Pick a random number, thus picking a rune
				n := rand.Intn(len(s.fillerRunes))
				ws.PuzzleRows[row][col] = s.fillerRunes[n]
			}
		}
	}
}

func (s *wordSearchConstructor) applyWeightedFiller(ws *WordSearch) {
	for col := 0; col < ws.NumCols; col++ {
		for row := 0; row < ws.NumRows; row++ {
			if s.isCellAvailable(col, row) {
				// Pick a random number
				rn := rand.Int63n(s.fillerWeightsSum)
				// Find where it would fit in the slice of
				// weights
				n, _ := slices.BinarySearch(s.fillerWeights, rn)
				// Use that rune
				ws.PuzzleRows[row][col] = s.fillerRunes[n]
			}
		}
	}
}
