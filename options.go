// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

type WordSearchOption func(*wordSearchConstructor) error

func AddNaturalLTRDirections() WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections |= NaturalLTRDirections
		return nil
	}
}

func AddUnnaturalDirections() WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections |= UnnaturalLTRDirections
		return nil
	}
}

func AddDirection(direction Direction) WordSearchOption {
	return func(constructor *wordSearchConstructor) error {
		constructor.possibleDirections |= direction
		return nil
	}
}

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
		constructor.fillerRunes = make([]rune, len(filler))
		copy(constructor.fillerRunes, filler)
		return nil
	}
}
