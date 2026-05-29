// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"runtime"
	"slices"
	"sort"
	"time"
)

// This type used by options which can be given to Construct()
// Construct accepts zero or more options.
type WordSearchOption func(*wordSearchConstructor) error

// Add all left-to-right directions (and "down")
func AddNaturalLTRDirections() WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections |= NaturalLTRDirections
		return nil
	}
}

// Add all right-to-left directions (and "up")
func AddUnnaturalLTRDirections() WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections |= UnnaturalLTRDirections
		return nil
	}
}

// Add a direction. The direction argument may be one direction, or it
// may represent multiple directions.
func AddDirection(direction Direction) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections |= direction
		return nil
	}
}

// Use all directions
func UseAllDirections() WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections = AllDirections
		return nil
	}
}

// Initialize the Constructor with this random seed
func RandomSeed(seed int64) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.randomSeed = seed
		constructor.randomSeedGiven = true
		return nil
	}
}

// Use the runes in this string for the filler. Every rune has
// an equal chance of being used
func FillUniformlyFromString(filler string) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		// Convert the string into a slice of strings, one for each rune
		fillerStrings := make([]string, 0, len(filler))
		// Iterate rune by rune
		for _, r := range filler {
			fillerStrings = append(fillerStrings, string(r))
		}
		constructor.fillerStrings = fillerStrings
		return nil
	}
}

// Use the strings in this slice for the filler. Every string has
// an equal chance of being used. This routine makes a private
// copy of the slice.
func FillUniformlyFromStringSlice(filler []string) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		// Copy the user's slice
		constructor.fillerStrings = slices.Clone(filler)
		return nil
	}
}

// Use the runes in this slice for the filler. Every rune has
// an equal chance of being used
func FillUniformlyFromRuneSlice(filler []rune) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		// Convert the rune slice into a slice of strings, one for each rune
		fillerStrings := make([]string, len(filler))
		for i, r := range filler {
			fillerStrings[i] = string(r)
		}
		constructor.fillerStrings = fillerStrings
		return nil
	}
}

// Combines a string for a cell's item and a relative weight. This is used
// with the FillWeighted option to the Construct function.
type FillerWeight struct {
	String string
	Weight int64
}

// Use the weighted runes in this slice for the filler. The weights
// do not need to add up to any specific value. The Constructor will
// sum the weights and calculate the percentages.
func FillWeighted(filler []FillerWeight) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		// Make our own temporary copy of the slice, so we can
		// sort it without destroying the caller's slice.
		fillerWeights := make([]FillerWeight, len(filler))
		copy(fillerWeights, filler)

		// Sort the slice, ascending, by weight
		sort.Slice(fillerWeights, func(i, j int) bool {
			rwi := fillerWeights[i]
			rwj := fillerWeights[j]
			if rwi.Weight == rwj.Weight {
				return rwi.String < rwj.String
			} else {
				return rwi.Weight < rwj.Weight
			}
		})

		// Allocate our own slices
		constructor.fillerStrings = make([]string, len(filler))
		constructor.fillerWeights = make([]int64, len(filler))

		// Calculate the sum of all the weights, and adjust
		// the running weight sum in each FillerWeight
		for i, rw := range fillerWeights {
			if rw.Weight <= 0 {
				return ErrFillerWeightNotPositive
			}
			constructor.fillerStrings[i] = rw.String
			constructor.fillerWeightsSum += rw.Weight
			constructor.fillerWeights[i] = constructor.fillerWeightsSum
		}

		return nil
	}
}

// Set the time limit for constructing the puzzle. If the time limit is
// reached before finding a solution, ErrReachedTimeLimit is returned.
func WithTimeLimit(duration time.Duration) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.timeLimit = duration
		return nil
	}
}

// Set the maximum number of entries kept in the dead-end memoization
// cache. The cache short-circuits the search when a (cell-state, word
// index) pair is revisited; once the cap is reached, additional
// failures stop being recorded but existing entries keep serving
// lookups. Each entry takes roughly 150-250 bytes of memory depending
// on grid dimensions and cell-content byte length, so the default of
// 10,000 holds the cache around 2 MB. Pass 0 to disable memoization
// entirely.
func WithMemoizationLimit(maxEntries int) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.deadEndCacheCap = maxEntries
		return nil
	}
}

// Run the search in parallel across n goroutines, each starting from a
// distinct RNG seed. The first goroutine to find a solution wins; the
// others are cancelled. This helps in the "borderline density" regime
// where a single seed may wander for a long time but a different seed
// finds the answer quickly (see the "Density and search time" section
// in the README).
//
// Semantics:
//   - n == 0: use runtime.NumCPU() workers.
//   - n == 1: equivalent to leaving the option off (single-threaded).
//   - n  > 1: race that many workers.
//   - n  < 0: treated as 1.
//
// Combining WithParallelism(n > 1) with RandomSeed returns
// ErrSeedWithParallelism, because parallel workers each require an
// independent seed and silently overriding the caller's seed would be
// surprising.
func WithParallelism(n int) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		if n == 0 {
			n = runtime.NumCPU()
		}
		if n < 1 {
			n = 1
		}
		constructor.parallelism = n
		return nil
	}
}
