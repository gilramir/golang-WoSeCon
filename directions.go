// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"fmt"
	"strings"
)

type Direction uint8

// Public
const (
	// Natural directions for LTR languages
	Down          Direction = 1
	LTRHorizontal Direction = 2
	LTRDescending Direction = 4
	LTRAscending  Direction = 8

	// Unnatrual directions for LTR languages
	Up            Direction = 16
	RTLHorizontal Direction = 32
	RTLDescending Direction = 64
	RTLAscending  Direction = 128

	NaturalLTRDirections   = Down | LTRHorizontal | LTRDescending | LTRAscending
	UnnaturalLTRDirections = Up | RTLHorizontal | RTLDescending | RTLAscending

	AllDirections = NaturalLTRDirections | UnnaturalLTRDirections
)

// Private
const (
	goesDown = LTRDescending | RTLDescending | Down
	goesUp   = LTRAscending | RTLAscending | Up
	goesLTR  = LTRDescending | LTRAscending | LTRHorizontal
	goesRTL  = RTLDescending | RTLAscending | RTLHorizontal
)

// Return a slice of the directiosn set in this Direction object
func getIndividualDirections(d Direction) []Direction {
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

func DirectionString(d Direction) string {
	directionsSlice := getIndividualDirections(d)
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
