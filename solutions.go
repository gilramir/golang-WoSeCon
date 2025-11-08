// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"fmt"
	"sort"
)

type wordDetailsType struct {
	seq               Sequence
	firstItem         string
	officialPlacement WordPlacement
	seqLen            int
}

func (s *WordSearch) findAllPossibleSolutions(possibleDirections Direction) {

	// Populate AllPossibleWordPlacements from WordPlacements
	for word, placement := range s.WordPlacements {
		placements := make([]WordPlacement, 1)
		placements[0] = placement
		s.AllPossibleWordPlacements[word] = placements
	}

	// Keep track of various pieces of data so we can sort the words
	// in order
	orderedWordDetails := make([]wordDetailsType, len(s.WordPlacements))
	i := 0
	for word, placement := range s.WordPlacements {
		seq, has := s.Sequences[word]
		if !has {
			panic(fmt.Sprintf("Word %s was not seen in Sequences", word))
		}

		orderedWordDetails[i] = wordDetailsType{
			seq:               seq,
			firstItem:         seq.Index(0),
			officialPlacement: placement,
			seqLen:            len([]rune(word)),
		}
		i++
	}

	// Sort the wordDetails from largest to smallest word. We must search
	// in this order to avoid matching a shorter word as a substring
	// of a larger word.

	// Stable-sort them by size, descending
	sort.Slice(orderedWordDetails, func(i, j int) bool {
		wdi := orderedWordDetails[i]
		wdj := orderedWordDetails[j]
		if wdi.seqLen == wdi.seqLen {
			// No reason, but, we do a stable sort
			// Is wdi < wdj?
			return wdi.seq.Cmp(wdj.seq) == -1
		} else {
			return wdi.seqLen > wdj.seqLen
		}
	})

	// Find all possible solutions
	possibleDirectionsSlice := DirectionSlice(possibleDirections)
	for row, columns := range s.PuzzleRows {
	next_column:
		for col, cellItem := range columns {
		next_word:
			for _, wordDetails := range orderedWordDetails {
				seq := wordDetails.seq
				word := seq.String()
				// Can this word even start here?
				if cellItem != wordDetails.firstItem {
					continue next_word
				}

				// Yes, it can start here.
				// It can fit here. Is it a single-rune word?
				officialPlacement := wordDetails.officialPlacement
				if wordDetails.seqLen == 1 {
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
				manyD := s.directionsForSequenceAtLocation(seq, row, col,
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

func (s *WordSearch) directionsForSequenceAtLocation(seq Sequence,
	startingRow int, startingCol int, possibleDirectionsSlice []Direction, officialPlacement WordPlacement) Direction {

	var solutionDirections Direction

	for _, d := range possibleDirectionsSlice {
		// Is this cell this word's official placement? Don't record it twice
		if officialPlacement.Row == startingRow && officialPlacement.Col == startingCol && d == officialPlacement.Direction {
			continue
		}

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

		// Iterate rune by rune
		matched := true
		for _, seqItem := range seq.Items() {
			// Out of bounds?
			if row < 0 || row >= s.NumRows || col < 0 || col >= s.NumCols {
				//				fmt.Printf("out of bounds\n")
				matched = false
				break
			}
			// We already know the first rune matched,
			// since our single caller checked that already
			// But we'll check again in case we make a
			// mistake in the future.
			cellItem := s.PuzzleRows[row][col]
			if cellItem != seqItem {
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
