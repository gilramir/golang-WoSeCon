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
	ErrSmallerThanMinimumSize ErrorString = "the grid must be at least 4x4"
	ErrWordIsTooLong          ErrorString = "at least one word is larger than the grid dimensions"
)

// Construct a WordSearch
func Construct(numCols int, numRows int, words []string, opts ...WordSearchOption) (*WordSearch, error) {
	var constructor wordSearchConstructor

	if numRows < 4 || numCols < 4 {
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
