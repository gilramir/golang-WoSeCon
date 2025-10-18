// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

type WordSearchOption func(*WordSearchConstructor) error

func AddNaturalLTRDirections() WordSearchOption {
	return func(maker *WordSearchConstructor) error {
		maker.possibleDirections |= NaturalLTRDirections
		return nil
	}
}

func AddUnnaturalDirections() WordSearchOption {
	return func(maker *WordSearchConstructor) error {
		maker.possibleDirections |= UnnaturalLTRDirections
		return nil
	}
}

func AddDirection(direction Direction) WordSearchOption {
	return func(maker *WordSearchConstructor) error {
		maker.possibleDirections |= direction
		return nil
	}
}

func UseAllDirections() WordSearchOption {
	return func(maker *WordSearchConstructor) error {
		maker.possibleDirections = AllDirections
		return nil
	}
}

/*
// Never produce a WordSearch that has these words
// in *any* direction
func PreventBadWords(badWords []string) WordSearchOption {
	return func(maker *WordSearchConstructor) error {
		maker.badWords = badWords
		return nil
	}
}
*/

// Initialize the Constructor with this random seed
func RandomSeed(seed int64) WordSearchOption {
	return func(maker *WordSearchConstructor) error {
		maker.randomSeed = seed
		maker.randomSeedGiven = true
		return nil
	}
}
