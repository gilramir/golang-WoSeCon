// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"fmt"
	"strings"
)

// Any of the 8 possible directions is representing by a Direction.
// In some uses, a Direction represents a single direction. In other
// cases, it represents multiple directiosn, because it is a uint8 and
// each bit represents a direction.
type Direction uint8

const (
	// No direction. This is used to represent "not yet initialized",
	// and also the fact that single-rune words need no notion of
	// direction.
	NilDirection Direction = 0

	// Natural directions for LTR languages
	Down          Direction = 1
	LTRHorizontal Direction = 2
	LTRDescending Direction = 4
	LTRAscending  Direction = 8

	// Unnatural directions for LTR languages
	Up            Direction = 16
	RTLHorizontal Direction = 32
	RTLDescending Direction = 64
	RTLAscending  Direction = 128

	// The "natural" directions for a wordsearch puzzle for left-to-right
	// languages. This includes all three left-to-right directions, and
	// "down"
	NaturalLTRDirections = Down | LTRHorizontal | LTRDescending | LTRAscending

	// The opposite of NaturalLTRDirections. It has all three right-to-left
	// diretions, and "up"
	UnnaturalLTRDirections = Up | RTLHorizontal | RTLDescending | RTLAscending

	// All 8 directions.
	AllDirections = NaturalLTRDirections | UnnaturalLTRDirections

	// Any downward direction, including diagonal
	GoesDownward = LTRDescending | RTLDescending | Down

	// Any upward direction, including diagonal
	GoesUpward = LTRAscending | RTLAscending | Up

	// Any left-to-right direction, including diaagonal
	GoesLTR = LTRDescending | LTRAscending | LTRHorizontal

	// Any right-to-left direction, including diagonal
	GoesRTL = RTLDescending | RTLAscending | RTLHorizontal
)

// Return a slice of the directions set in this Direction object.
// It never includes NilDirection
func DirectionSlice(d Direction) []Direction {
	directions := make([]Direction, 0, 8)

	checkAndSet := func(od Direction) {
		if d&od == od {
			directions = append(directions, od)
		}
	}

	checkAndSet(Down)
	checkAndSet(LTRHorizontal)
	checkAndSet(LTRDescending)
	checkAndSet(LTRAscending)

	checkAndSet(Up)
	checkAndSet(RTLHorizontal)
	checkAndSet(RTLDescending)
	checkAndSet(RTLAscending)

	return directions
}

// Gives a string representation of a Direction.
// It never includes "NilDirection"
func DirectionString(d Direction) string {
	directionsSlice := DirectionSlice(d)
	labelsSlice := make([]string, len(directionsSlice))

	for i, idir := range directionsSlice {
		var name string
		switch idir {
		case Down:
			name = "Down"
		case LTRHorizontal:
			name = "LTRHorizontal"
		case LTRDescending:
			name = "LTRDescending"
		case LTRAscending:
			name = "LTRAscending"
		case Up:
			name = "Up"
		case RTLHorizontal:
			name = "RTLHorizontal"
		case RTLDescending:
			name = "RTLDescending"
		case RTLAscending:
			name = "RTLAscending"
		default:
			panic(fmt.Sprintf("Should not reach. value=%b", idir))
		}
		labelsSlice[i] = name
	}

	return strings.Join(labelsSlice, ",")
}
