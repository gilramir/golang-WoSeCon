// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"sort"
)

type wordDetailsType struct {
	word              string
	firstRune         rune
	officialPlacement WordPlacement
	runeLength        int
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
		orderedWordDetails[i] = wordDetailsType{
			word:              word,
			firstRune:         []rune(word)[0],
			officialPlacement: placement,
			runeLength:        len([]rune(word)),
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
		if wdi.runeLength == wdi.runeLength {
			// No reason, but, we do a stable sort
			return wdi.word < wdj.word
		} else {
			return wdi.runeLength > wdj.runeLength
		}
	})

	// Find all possible solutions
	possibleDirectionsSlice := DirectionSlice(possibleDirections)
	for row, columns := range s.PuzzleRows {
	next_column:
		for col, cellRune := range columns {
		next_word:
			for _, wordDetails := range orderedWordDetails {
				word := wordDetails.word
				// Can this word even start here?
				if cellRune != wordDetails.firstRune {
					//					fmt.Printf("(%d, %d) cell=%s impossible, skip word=%s\n",
					//						row, col, string(cellRune), word)
					continue next_word
				}
				//				fmt.Printf("(%d, %d) cell=%s checking rune %s word=%s\n",
				//					row, col, string(cellRune), string(wordDetails.firstRune), word)

				// Yes, it can start here.
				// It can fit here. Is it a single-rune word?
				officialPlacement := wordDetails.officialPlacement
				//				fmt.Printf("Official placement: %+v\n", officialPlacement)
				if wordDetails.runeLength == 1 {
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
				//				fmt.Printf("manyD: %v\n", manyD)

				// Here, NilDirection means "not found"
				// because we took care of single-rune
				// words above
				if manyD != NilDirection {
					for _, d := range DirectionSlice(manyD) {
						//						fmt.Printf("Adding placement d=%s\n", DirectionString(d))
						placements := s.AllPossibleWordPlacements[word]
						placements = append(placements, WordPlacement{
							Col:       col,
							Row:       row,
							Direction: d,
						})
						s.AllPossibleWordPlacements[word] = placements
						//						fmt.Printf("now %s => %+v\n", word, placements)
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
		// Is this cell this word's official placement? Don't record it twice
		if officialPlacement.Row == startingRow && officialPlacement.Col == startingCol && d == officialPlacement.Direction {
			//			fmt.Printf("Skipping official placement")
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
		//		fmt.Printf("Checking %s rowAdj=%d colAdj=%d\n", DirectionString(d), rowAdj, colAdj)

		// Iterate rune by rune
		matched := true
		for _, wordRune := range word {
			//			fmt.Printf("Checking (%d, %d)\n", row, col)
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
			cellRune := s.PuzzleRows[row][col]
			//			fmt.Printf("Checking cellRune %s ; wordRune %s\n", string(cellRune), string(wordRune))
			if cellRune != wordRune {
				//				fmt.Printf("not a match\n")
				// Not a match
				matched = false
				break
			}
			row += rowAdj
			col += colAdj
		}
		// Found a solution in this direction
		if matched {
			//			fmt.Printf("Matched: %s\n", DirectionString(d))
			solutionDirections |= d
		}
	}

	return solutionDirections
}
