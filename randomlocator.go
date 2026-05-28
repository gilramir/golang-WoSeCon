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
	// Shared RNG, kept so that minus() can reshuffle on each call.
	rng *rand.Rand
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
		rng:                s.rng,
	}
	for _, target := range targets {
		newLocator.remove(target)
	}
	// Freshly shuffle so each backward-mode entry probes a different
	// order. This matches the C++ reference's RandomLocator::minus(),
	// which builds a fresh, freshly-shuffled RandomLocator on every
	// call. The global locator's ordering becomes increasingly biased
	// as items are added and removed during the search, so re-shuffling
	// here injects randomness that helps the algorithm escape patterns
	// where the same constrained locations keep being tried first.
	s.rng.Shuffle(len(newLocator.availableLocations), func(i, j int) {
		newLocator.availableLocations[i], newLocator.availableLocations[j] =
			newLocator.availableLocations[j], newLocator.availableLocations[i]
	})
	return newLocator
}

// Create all the directedLocations that fit within the possible directions
// the caller set.
func (s *randomLocator) initialize(numCols, numRows int, possibleDirections Direction, rng *rand.Rand) {

	s.rng = rng

	var dlmatrix directedLocationMatrix
	dlmatrix.initialize(numCols, numRows, possibleDirections)

	s.availableLocations = dlmatrix.getAllDirectedLocations()

	// Shuffle
	rng.Shuffle(len(s.availableLocations), func(i, j int) {
		s.availableLocations[i], s.availableLocations[j] = s.availableLocations[j], s.availableLocations[i]
	})
	//	fmt.Printf("locator init: %d directedlocations\n", len(s.availableLocations))

	/*
		for i, d := range s.availableLocations {
			fmt.Printf("#%d. %d,%d %s\n", i+1, d.col, d.row, DirectionString(d.direction))
		}
	*/
}
