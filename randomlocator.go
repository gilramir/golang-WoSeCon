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

/*
func (s *rancomLocator) addLocationsForLength(startLoc directedLocation, seqLen int, possibleDirections Direction) {

	location := startLoc
	var colAdj int
	var rowAdj int

	// Find endRow
	if location.direction&GoesDownward > 0 {
		rowAdj = 1
	} else if location.direction&GoesUpward > 0 {
		rowAdj = -1
	} // else horizontal

	// Find endCol
	if location.direction&GoesLTR != 0 {
		colAdj = 1
	} else if location.direction&GoesRTL != 0 {
		colAdj = -1
	} // else vertical

	dlcons = DirectedLocationConstructor{
		numRows:            numRows,
		numCols:            numCols,
		possibleDirections: possibleDirections,
	}

	col := location.col
	row := location.row
	for i := 0; i < wordInfo.seqLen; i++ {
		newDirectedLocations = dlcons.constructPossibleDirectedLocations(col, row)

		s.availableLocations = append(s.availableLocations, newDirectedLocations)

		col += colAdj
		row += rowAdj
	}
}

func (s *rancomLocator) removeLocationsForLength(startLoc directedLocation, seqLen int) {

	toRemove := make([]directedLocation, seqLen)

	location := startLoc
	var colAdj int
	var rowAdj int

	// Find endRow
	if location.direction&GoesDownward > 0 {
		rowAdj = 1
	} else if location.direction&GoesUpward > 0 {
		rowAdj = -1
	} // else horizontal

	// Find endCol
	if location.direction&GoesLTR != 0 {
		colAdj = 1
	} else if location.direction&GoesRTL != 0 {
		colAdj = -1
	} // else vertical

	col := location.col
	row := location.row
	for i := 0; i < wordInfo.seqLen; i++ {
		toRemove = append(toRemove, directedLocation{
			col: col,
			row: row,
		})

		col += colAdj
		row += rowAdj
	}

	newLocations := make([]directedLocation, 0, len(s.availableLocations))
locs:
	for _, loc := range newLocations {
		for _, toRem := range toRemove {
			if loc.col == toRem.col && loc.row == toRem.row {
				// ignore it
				continue locs
			} else {
				newLocations = append(newLocations, loc)
			}
		}
	}
	s.availableLocations = newLocations

}
*/

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
func (s *randomLocator) initialize(dlmatrix *directedLocationMatrix, rng *rand.Rand) {

	s.availableLocations = dlmatrix.getAllDirectedLocations()

	// Shuffle
	rng.Shuffle(len(s.availableLocations), func(i, j int) {
		s.availableLocations[i], s.availableLocations[j] = s.availableLocations[j], s.availableLocations[i]
	})

	/*
		for i, d := range s.availableLocations {
			fmt.Printf("#%d. %d,%d %s\n", i+1, d.col, d.row, DirectionString(d.direction))
		}
	*/
}
