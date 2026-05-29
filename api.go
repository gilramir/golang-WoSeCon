// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

// The error type that Construct() returns. It's just a string.
type ErrorString string

func (e ErrorString) Error() string {
	return string(e)
}

// The errors that Construct() can return. These are the only errors that
// Construct() returns.
const (
	ErrCannotFitWords          ErrorString = "the words cannot be placed into this size of grid"
	ErrSmallerThanMinimumSize  ErrorString = "the grid must be at least 4x4"
	ErrWordIsTooLong           ErrorString = "at least one word is larger than the grid dimensions"
	ErrFillerWeightNotPositive ErrorString = "at least one filler string's weight is <= 0"
	ErrReachedTimeLimit        ErrorString = "reached time limit before finding solution"
	ErrSeedWithParallelism     ErrorString = "RandomSeed cannot be combined with WithParallelism(n > 1); parallel workers each need an independent seed"
)

// errCancelled is an internal sentinel returned by construct() when a
// sibling worker has already found a solution and signalled the rest to
// stop. It is never returned to callers — the parallel orchestrator
// filters it out before returning the final error.
const errCancelled ErrorString = "cancelled by sibling worker"

// Construct a WordSearch, where each word is a Go string.
// One rune fits in one puzzle cell.
// The minum size is 4x4, as anything less doesn't make
// sense as a puzzle. Duplicate words are removed.
func Construct(numCols int, numRows int, words []string, opts ...WordSearchOption) (*WordSearch, error) {
	sequences := make([]Sequence, len(words))
	for i, word := range words {
		sequences[i] = NewRuneSequence(word)
	}

	return ConstructFromSequences(numCols, numRows, sequences, opts...)
}

// Construct a WordSerch, where each word is a slice of objects that are
// comparable. One object fits in one puzzle cell.
func ConstructFromSequences(numCols int, numRows int, sequences []Sequence, opts ...WordSearchOption) (*WordSearch, error) {
	if numRows < 3 || numCols < 3 {
		return nil, ErrSmallerThanMinimumSize
	}

	// Apply options once on a probe constructor so we can detect
	// WithParallelism + RandomSeed conflicts before launching workers,
	// and so we can re-use the constructor on the single-threaded path
	// without paying init twice.
	var constructor wordSearchConstructor
	if err := constructor.init(numCols, numRows, sequences, opts...); err != nil {
		return nil, err
	}

	if constructor.parallelism > 1 {
		if constructor.randomSeedGiven {
			return nil, ErrSeedWithParallelism
		}
		return constructParallel(numCols, numRows, sequences, opts, constructor.parallelism)
	}

	if err := constructor.construct(); err != nil {
		return nil, err
	}
	return constructor.translateToWordSearch(), nil
}
