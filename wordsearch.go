// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"slices"
)

// The constructed word search puzzle, with the solution too.
type WordSearch struct {
	// The number of columns in the puzzle
	NumCols int
	// The number of rows in the puzzle
	NumRows int

	// Strings and how they are split into sequences
	// This was used when constructing the puzzle
	Sequences map[string]Sequence

	// The "official" solution, as words and placements (location+direction)
	// This solution was create by the algorithm
	WordPlacements map[string]WordPlacement

	// All possible solutions. After filler is added, multiple
	// solutions might appear for a single word. If a single word
	// is a sub-string of a larger word, the solution is counted
	// for the larger word, not the shorter word.
	// If there are multiple solutions for the same word, starting
	// at the same placement (going in different directions), there
	// will be mulitple WordPlacements for that word. In other words,
	// the Direction in any single WordPlacement indicates just a single
	// direction.
	AllPossibleWordPlacements map[string][]WordPlacement

	// How many words have more than one solution
	NumWordsWithMultipleSolutions int

	// The "official" solution.
	// This solution was create by the algorithm and does not
	// take into account other possible solutions after the filler
	// was added.
	// Places where filler would be in the grid are indicated by
	// an empty string ("").
	// First dimension is y (rows)
	// Second dimension is x (column)
	SolutionRows [][]string

	// The puzzle.
	// Every cell is filled. It has the algorithm's "official" solution
	// and the filler, combined.
	// First dimension is y (rows)
	// Second dimension is x (column)
	PuzzleRows [][]string
}

// Indicate the placement of a word in the puzzle
// by the (column, row) coordinate of its starting point,
// and the direction in which the word goes.
type WordPlacement struct {
	Col       int
	Row       int
	Direction Direction
}

// The string representation of the Direction in the WordPlacement
func (s WordPlacement) DirectionString() string {
	return DirectionString(s.Direction)
}

// Take a successful wordSearchConstructor,
// and return a WordSearch result which will be sent back to the caller
func (s *wordSearchConstructor) translateToWordSearch() *WordSearch {

	ws := &WordSearch{
		NumCols:                   s.numCols,
		NumRows:                   s.numRows,
		Sequences:                 make(map[string]Sequence),
		WordPlacements:            make(map[string]WordPlacement),
		AllPossibleWordPlacements: make(map[string][]WordPlacement),
		SolutionRows:              make([][]string, s.numRows),
		PuzzleRows:                make([][]string, s.numRows),
	}

	// Allocate the columns, filled with empty strings
	for row := 0; row < s.numRows; row++ {
		ws.SolutionRows[row] = slices.Repeat([]string{""}, s.numCols)
		ws.PuzzleRows[row] = slices.Repeat([]string{""}, s.numCols)
	}

	// Fill in the solution and puzzle, and copy the sequences
	for _, wordInfo := range s.wordInfos {
		ws.Sequences[wordInfo.seq.String()] = wordInfo.seq

		dl := wordInfo.placement
		direction := dl.direction
		// If the word is only 1-item long, it's directionless
		if wordInfo.seqLen == 1 {
			direction = NilDirection
		}
		ws.WordPlacements[wordInfo.seq.String()] = WordPlacement{
			Col:       dl.col,
			Row:       dl.row,
			Direction: direction,
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
		} // else vertical

		var col = dl.col
		var row = dl.row
		// Iterate item by item in the string
		for _, cellItem := range wordInfo.seq.Items() {
			ws.SolutionRows[row][col] = cellItem
			ws.PuzzleRows[row][col] = cellItem
			col += colAdj
			row += rowAdj
		}
	}

	// Add the filler
	if len(s.fillerWeights) > 0 {
		s.applyWeightedFiller(ws)
	} else if len(s.fillerStrings) > 0 {
		s.applyUniformFiller(ws)
	}

	// The filler may have created more solutions.
	// Find them all
	ws.findAllPossibleSolutions(s.possibleDirections)

	// Count the words with multiple worde placements
	for _, wordPlacements := range ws.AllPossibleWordPlacements {
		if len(wordPlacements) > 1 {
			ws.NumWordsWithMultipleSolutions += 1
		}
	}

	return ws
}

func (s *wordSearchConstructor) applyUniformFiller(ws *WordSearch) {
	for col := 0; col < ws.NumCols; col++ {
		for row := 0; row < ws.NumRows; row++ {
			if s.isCellAvailable(col, row) {
				// Pick a random number, thus picking a rune
				n := s.rng.Intn(len(s.fillerStrings))
				ws.PuzzleRows[row][col] = s.fillerStrings[n]
			}
		}
	}
}

func (s *wordSearchConstructor) applyWeightedFiller(ws *WordSearch) {
	for col := 0; col < ws.NumCols; col++ {
		for row := 0; row < ws.NumRows; row++ {
			if s.isCellAvailable(col, row) {
				// Pick a random number
				rn := s.rng.Int63n(s.fillerWeightsSum)
				// Find where it would fit in the slice of
				// weights
				n, _ := slices.BinarySearch(s.fillerWeights, rn)
				// Use that rune
				ws.PuzzleRows[row][col] = s.fillerStrings[n]
			}
		}
	}
}
