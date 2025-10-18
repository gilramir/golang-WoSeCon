// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

// Define a new type that implements the error interface
type ErrorString string

func (e ErrorString) Error() string {
	return string(e)
}

// The errors that Construct() can return
const (
	ErrCannotFitWords         ErrorString = "the words cannot be placed into this size of grid"
	ErrSmallerThanMinimumSize ErrorString = "the grid must be at least 8x8"
)

// The constructed word search puzzle
type WordSearch struct {
	NumCols int
	NumRows int

	// The solution, as words and placments (location+direction)
	WordPlacements map[string]WordPlacement

	// The soluation, as runes
	// First dimension is x (columns)
	// Second dimension is y (rows)
	SolutionRuneMatrix [][]rune
}

type WordPlacement struct {
	Col       int
	Row       int
	Direction Direction
}

// Construct a WordSearch
func Construct(numCols int, numRows int, words []string, opts ...WordSearchOption) (*WordSearch, error) {
	var constructor WordSearchConstructor

	if numRows < 8 || numCols < 8 {
		return nil, ErrSmallerThanMinimumSize
	}

	err := constructor.init(numCols, numRows, words, opts...)
	if err != nil {
		return nil, err
	}
	err = constructor.construct()
	if err != nil {
		return nil, err
	}

	return constructor.translateToWordSearch(), nil
}
