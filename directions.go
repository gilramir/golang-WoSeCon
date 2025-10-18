// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

type Direction uint8

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
