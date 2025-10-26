// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"sort"
)

func (s *WordSearch) findAllPossibleSolutions(possibleDirections Direction) {

	// Populate AllPossibleWordPlacements from WordPlacements
	for word, placement := range s.WordPlacements {
		placements := make([]WordPlacement, 1)
		placements[0] = placement
		s.AllPossibleWordPlacements[word] = placements
	}

	// Sort the words from largest to smallest. We must search
	// in this order to avoid matching a shorter word as a substring
	// of a larger word.
	// Make other slices of the same size that will help us during
	// the search.
	orderedWords := make([]string, len(s.WordPlacements))
	firstRunes := make([]rune, len(s.WordPlacements))
	officialPlacements := make([]WordPlacement, len(s.WordPlacements))
	wordRuneLength := make([]int, len(s.WordPlacements))

	i := 0
	for word, placement := range s.WordPlacements {
		orderedWords[i] = word
		firstRunes[i] = []rune(word)[0]
		officialPlacements[i] = placement
		wordRuneLength[i] = len([]rune(word))
		i++
	}

	// Stable-sort them by size, descending
	sort.Slice(orderedWords, func(i, j int) bool {
		wi := orderedWords[i]
		wj := orderedWords[j]
		if len(wi) == len(wj) {
			// No reason, but, we do a stable sort
			return wi < wj
		} else {
			return len(wi) > len(wj)
		}
	})

	// Find all possible solutions
	possibleDirectionsSlice := DirectionSlice(possibleDirections)
	for row, columns := range s.PuzzleRows {
	next_column:
		for col, cellRune := range columns {
		next_word:
			for wi, word := range orderedWords {
				// Can this word even start here?
				wordFirstRune := firstRunes[wi]
				if cellRune != wordFirstRune {
					continue next_word
				}

				// Yes, it can start here.
				// It can fit here. Is it a single-rune word?
				officialPlacement := officialPlacements[wi]
				if wordRuneLength[wi] == 1 {
					if officialPlacement.Row != row || officialPlacement.Col != col {
						placements := s.AllPossibleWordPlacements[word]
						placements = append(placements, WordPlacement{
							Col:       col,
							Row:       row,
							Direction: NilDirection,
						})
						s.AllPossibleWordPlacements[word] = placements
					}
					continue next_word
				}

				// See if we can place this word at this point,
				// in any direction
				manyD := s.directionsForWordAtLocation(word, row, col,
					possibleDirectionsSlice, officialPlacement)

				// Here, NilDirection means "not found"
				// because we took care of single-rune
				// words above
				if manyD != NilDirection {
					for _, d := range DirectionSlice(manyD) {
						placements := s.AllPossibleWordPlacements[word]
						placements = append(placements, WordPlacement{
							Col:       col,
							Row:       row,
							Direction: d,
						})
						s.AllPossibleWordPlacements[word] = placements
					}
					continue next_column
				}
			}
		}
	}
}

func (s *WordSearch) directionsForWordAtLocation(word string,
	startingRow int, startingCol int, possibleDirectionsSlice []Direction, officialPlacement WordPlacement) Direction {

	var solutionDirections Direction

	for _, d := range possibleDirectionsSlice {
		var colAdj int
		var rowAdj int

		if d&GoesDownward > 0 {
			rowAdj = 1
		} else if d&GoesUpward > 0 {
			rowAdj = -1
		} // else horizontal

		if d&GoesLTR > 0 {
			colAdj = 1
		} else if d&GoesRTL > 0 {
			colAdj = -1
		} // else vertical

		row := startingRow
		col := startingCol

		// Is this cell this word's official placement? Don't record it twice
		if officialPlacement.Row == row && officialPlacement.Col == col && d == officialPlacement.Direction {
			continue
		}

		// Iterate rune by rune
		matched := true
		for _, wordRune := range word {
			// Out of bounds?
			if row < 0 || row >= s.NumRows || col < 0 || col >= s.NumCols {
				matched = false
				break
			}
			// We already know the first rune matched,
			// since our single caller checked that already
			// But we'll check again in case we make a
			// mistake in the future.
			cellRune := s.PuzzleRows[row][col]
			if cellRune != wordRune {
				// Not a match
				matched = false
				break
			}
			row += rowAdj
			col += colAdj
		}
		// Found a solution in this direction
		if matched {
			solutionDirections |= d
		}
	}

	return solutionDirections
}
