// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import "sort"

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

// Add a single direction
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

/*
// Never produce a WordSearch that has these words
// in *any* direction
func PreventBadWords(badWords []string) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.badWords = badWords
		return nil
	}
}
*/

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
		constructor.fillerRunes = []rune(filler)
		return nil
	}
}

// Use the runes in this slice for the filler. Every rune has
// an equal chance of being used
func FillUniformlyFromRuneSlice(filler []rune) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		// Make our own copy of the slice
		constructor.fillerRunes = make([]rune, len(filler))
		copy(constructor.fillerRunes, filler)
		return nil
	}
}

// Combines a rune and a relative weight. This is used
// with the FillWeighted option to the Construct function.
type RuneWeight struct {
	Rune   rune
	Weight int64
}

// Use the weighted runes in this slice for the filler. The weights
// do not need to add up to any specific value. The Constructor will
// sum the weights and calculate the percentages.
func FillWeighted(filler []RuneWeight) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		// Make our own temporary copy of the slice, so we can
		// sort it without destroying the caller's slice.
		runeWeights := make([]RuneWeight, len(filler))
		copy(runeWeights, filler)

		// Sort the slice, ascending, by weight
		sort.Slice(runeWeights, func(i, j int) bool {
			rwi := runeWeights[i]
			rwj := runeWeights[j]
			if rwi.Weight == rwj.Weight {
				return rwi.Rune < rwj.Rune
			} else {
				return rwi.Weight < rwj.Weight
			}
		})

		// Allocate our own slices
		constructor.fillerRunes = make([]rune, len(filler))
		constructor.fillerWeights = make([]int64, len(filler))

		// Calculate the sum of all the weights, and adjust
		// the running weight sum in each RuneWeight
		for i, rw := range runeWeights {
			if rw.Weight <= 0 {
				return ErrFillerWeightNotPositive
			}
			constructor.fillerRunes[i] = rw.Rune
			constructor.fillerWeightsSum += rw.Weight
			constructor.fillerWeights[i] = constructor.fillerWeightsSum
		}

		return nil
	}
}
