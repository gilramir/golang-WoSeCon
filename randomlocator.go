// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"fmt"
	"math/rand"
)

type randomLocator struct {
	// One directedLocation for every cell and possible direction
	availableLocations []directedLocation
}

func (s *randomLocator) size() int {
	return len(s.availableLocations)
}

func (s *randomLocator) get(n int) directedLocation {
	return s.availableLocations[n]
}

func (s *randomLocator) add(d directedLocation) {
	s.availableLocations = append(s.availableLocations, d)
}

func (s *randomLocator) remove(target directedLocation) {
	for i := 0; i < len(s.availableLocations); i++ {

		if s.availableLocations[i].equals(target) {
			s.removeN(i)
			return
		}
	}
	panic(fmt.Sprintf("Should not reach, target=%+v", target))
}

func (s *randomLocator) removeN(index int) {
	s.availableLocations = append(s.availableLocations[:index], s.availableLocations[index+1:]...)
}

func (s *randomLocator) minus(targets []directedLocation) *randomLocator {
	newLocations := make([]directedLocation, len(s.availableLocations))
	copy(newLocations, s.availableLocations)
	newLocator := &randomLocator{
		availableLocations: newLocations,
	}
	for _, target := range targets {
		newLocator.remove(target)
	}
	return newLocator
}

// Create all the directedLocations that fit within the possible directions
// the caller set.
func (s *randomLocator) initialize(numCols, numRows int, possibleDirections Direction, rng *rand.Rand) {

	// The cells inside (non-border) of the grid can go in all directions
	// Initialize those first
	possibleDirectionsSlice := getIndividualDirections(possibleDirections)
	numInnerCells := (numCols - 1) * (numRows - 1) * len(possibleDirectionsSlice)
	s.availableLocations = make([]directedLocation, numInnerCells)
	i := 0
	for col := 1; col < numCols-1; col++ {
		for row := 1; row < numRows-1; row++ {
			for _, direction := range possibleDirectionsSlice {
				d := directedLocation{
					col:       col,
					row:       row,
					direction: direction,
				}
				s.availableLocations[i] = d
				i++
			}
		}
	}

	// Make masks to mask down certain directions
	goingUpMask := ^(Up | LTRAscending | RTLAscending)
	goingRightMask := ^(LTRHorizontal | LTRAscending | LTRDescending)
	goingLeftMask := ^(RTLHorizontal | RTLAscending | RTLDescending)
	goingDownMask := ^(Down | LTRDescending | RTLDescending)

	// Top row; cannot ascend or go up
	topRowDirections := possibleDirections & goingUpMask
	topRowDirectionsSlice := getIndividualDirections(topRowDirections)
	for col := 0; col < numCols; col++ {
		var directionsSlice []Direction
		// Corners have more restrictions
		if col == 0 {
			directionsSlice = getIndividualDirections(topRowDirections & goingLeftMask)
		} else if col == numCols-1 {
			directionsSlice = getIndividualDirections(topRowDirections & goingRightMask)
		} else {
			directionsSlice = topRowDirectionsSlice
		}
		for _, direction := range directionsSlice {
			d := directedLocation{
				col:       col,
				row:       0,
				direction: direction,
			}
			s.availableLocations = append(s.availableLocations, d)
		}
	}

	// Bottomm row; cannot descend or go down
	bottomRowDirections := possibleDirections & goingDownMask
	bottomRowDirectionsSlice := getIndividualDirections(bottomRowDirections)
	for col := 0; col < numCols; col++ {
		var directionsSlice []Direction
		// Corners have more restrictions
		if col == 0 {
			directionsSlice = getIndividualDirections(bottomRowDirections & goingLeftMask)
		} else if col == numCols-1 {
			directionsSlice = getIndividualDirections(bottomRowDirections & goingRightMask)
		} else {
			directionsSlice = bottomRowDirectionsSlice
		}
		for _, direction := range directionsSlice {
			d := directedLocation{
				col:       col,
				row:       numRows - 1,
				direction: direction,
			}
			s.availableLocations = append(s.availableLocations, d)
		}
	}

	// Left column; cannot go left. We already did the top and bottom
	// corners, so we can skip them
	leftColDirectionsSlice := getIndividualDirections(possibleDirections & goingLeftMask)
	for row := 1; row < numRows-1; row++ {
		for _, direction := range leftColDirectionsSlice {
			d := directedLocation{
				col:       0,
				row:       row,
				direction: direction,
			}
			s.availableLocations = append(s.availableLocations, d)
		}
	}

	// Right column; cannot go right. We already did the top and bottom
	// corners, so we can skip them
	rightColDirectionsSlice := getIndividualDirections(possibleDirections & goingRightMask)
	for row := 1; row < numRows-1; row++ {
		for _, direction := range rightColDirectionsSlice {
			d := directedLocation{
				col:       numCols - 1,
				row:       row,
				direction: direction,
			}
			s.availableLocations = append(s.availableLocations, d)
		}
	}

	// Shuffle
	rng.Shuffle(len(s.availableLocations), func(i, j int) {
		s.availableLocations[i], s.availableLocations[j] = s.availableLocations[j], s.availableLocations[i]
	})
}
