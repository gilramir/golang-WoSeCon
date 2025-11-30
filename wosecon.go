// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"log"
	"math/rand"
	"sort"
	"time"
)

// The algorithm works in 2 modoes, forward or backward
type algoMode byte

const (
	forwardMode  algoMode = 'f'
	backwardMode algoMode = 'b'
)

type wordSearchConstructor struct {
	numCols            int
	numRows            int
	wordInfos          []*wordInfo
	possibleDirections Direction
	mode               algoMode
	locator            randomLocator
	solutionMatrix     solutionMatrix

	randomSeed      int64
	randomSeedGiven bool
	rng             *rand.Rand

	// Used for uniformly filling in the puzzle, and
	// also as one of the data structures for filling in
	// by weights
	fillerStrings []string

	// Used for filling in the puzzle by weights
	fillerWeights    []int64
	fillerWeightsSum int64

	// Optional time limit
	timeLimit time.Duration
}

func (s *wordSearchConstructor) init(numCols int, numRows int, sequences []Sequence, opts ...WordSearchOption) error {

	s.numCols = numCols
	s.numRows = numRows
	s.solutionMatrix.initialize(numCols, numRows)

	// Use the list of words (strings) from the user
	// to create our wordInfo objects, but unique-ify the word lsit
	s.wordInfos = make([]*wordInfo, 0, len(sequences))

	seen := make(map[string]bool)
	for _, seq := range sequences {
		if seen[seq.String()] {
			continue
		}
		s.wordInfos = append(s.wordInfos, newWordInfo(seq))
		seen[seq.String()] = true
	}

	// Apply the options
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return err
		}
	}

	// If no possibleDirections were given, use the default
	if s.possibleDirections == 0 {
		s.possibleDirections = NaturalLTRDirections
	}

	// Is any one word bigger than the grid?
	var maxWordSize int
	if s.possibleDirections&(GoesUpward|GoesDownward) != 0 {
		maxWordSize = numRows
	}
	if s.possibleDirections&(GoesLTR|GoesRTL) != 0 {
		maxWordSize = max(maxWordSize, numCols)
	}

	for _, wordInfo := range s.wordInfos {
		if wordInfo.seqLen > maxWordSize {
			return ErrWordIsTooLong
		}
	}

	// Initialize the random number generator if given a seed
	if s.randomSeedGiven {
		s.rng = rand.New(rand.NewSource(s.randomSeed))
	} else {
		s.rng = rand.New(rand.NewSource(rand.Int63()))
	}

	// Stable sort words, longest to shortest
	// If the same length, then in alphabetical order
	// It should be easiest to fit the longest words first.
	sort.Slice(s.wordInfos, func(i, j int) bool {
		wi := s.wordInfos[i]
		wj := s.wordInfos[j]
		if wi.seqLen == wj.seqLen {
			// is wi < wj?
			return wi.seq.Cmp(wj.seq) == -1
		} else {
			// Descending order
			return wj.seqLen < wi.seqLen
		}
	})
	return nil
}

// The main loop of the algorithm
func (s *wordSearchConstructor) construct() error {

	s.locator.initialize(s.numCols, s.numRows, s.possibleDirections, s.rng)

	s.mode = forwardMode
	currentWordIndex := 0
	currentWord := s.wordInfos[currentWordIndex]

	// Configure the timer, if we have a time limit
	var timer *time.Timer
	hasTimeLimit := s.timeLimit.Milliseconds() > 0
	if hasTimeLimit {
		timer = time.NewTimer(s.timeLimit)
		defer timer.Stop()
	}

	for {
		//		fmt.Printf("currentWordIndex=%d\n", currentWordIndex)
		//		s.solutionMatrix.dump()
		// Reached our time limit?
		if hasTimeLimit {
			select {
			case _ = <-timer.C:
				return ErrReachedTimeLimit
			default:
				// keep going
			}
		}

		if s.locateOne(currentWord) {
			if currentWordIndex == len(s.wordInfos)-1 {
				// Finished
				break
			}
			currentWordIndex++
			currentWord = s.wordInfos[currentWordIndex]
			s.mode = forwardMode
		} else {
			//			fmt.Printf("Need to backtrack, currentWordIndex=%d\n", currentWordIndex)
			// Couldn't place a word. Go backwards
			if currentWordIndex == 0 {
				return ErrCannotFitWords
			} else {
				currentWord.deleteTested()
				currentWordIndex--
				currentWord = s.wordInfos[currentWordIndex]
				s.clearPlacement(currentWord)
				s.mode = backwardMode
			}
		}
	}
	return nil
}

// Place one word into the puzzle. If it cannot, returns false.
func (s *wordSearchConstructor) locateOne(currentWord *wordInfo) bool {
	//	log.Printf("locateOne %s mode=%s", currentWord.text, string(s.mode))

	var localLocator *randomLocator

	if s.mode == backwardMode {
		dl := currentWord.getPlacement()
		s.locator.add(dl)
		currentWord.moveLocationToTested()
		localLocator = s.locator.minus(currentWord.getTested())
	} else {
		localLocator = &s.locator
	}

	/*
		if s.mode == backwardMode {
			fmt.Printf("trying %d suitableLocations\n", localLocator.size())
		}
	*/
	for locationIndex := 0; locationIndex < localLocator.size(); locationIndex++ {
		suitableLocation := localLocator.get(locationIndex)
		if s.validPlacement(currentWord, suitableLocation) {
			if localLocator == &s.locator {
				s.locator.removeN(locationIndex)
			} else {
				s.locator.remove(suitableLocation)
			}
			return true
		}
	}
	//	log.Printf("could not locateOne %s mode=%s", currentWord.text, string(s.mode))
	return false
}

// Will this word fit in the puzzle at this directedLocation?
// If so, the wordInfo's placement is updated
func (s *wordSearchConstructor) validPlacement(wordInfo *wordInfo, location directedLocation) bool {

	/*
		log.Printf("validPlacement %s trying %d,%d,%s", wordInfo.text,
			location.col, location.row, DirectionString(location.direction))
	*/

	// Perchance did we already test this location?
	for _, testedLocation := range wordInfo.testedLocations {
		if testedLocation.equals(location) {
			log.Printf("Tested already")
			return false
		}
	}

	// Starting from the start of the directedLocation, check if the word will
	// fit in the puzzle.
	var startCol = location.col
	var startRow = location.row
	var endCol int
	var endRow int
	var colAdj int
	var rowAdj int

	// Find endRow
	if location.direction&GoesDownward > 0 {
		endRow = startRow + wordInfo.seqLen
		rowAdj = 1
	} else if location.direction&GoesUpward > 0 {
		endRow = startRow - wordInfo.seqLen
		rowAdj = -1
	} // else horizontal

	// Find endCol
	if location.direction&GoesLTR != 0 {
		endCol = startCol + wordInfo.seqLen
		colAdj = 1
	} else if location.direction&GoesRTL != 0 {
		endCol = startCol - wordInfo.seqLen
		colAdj = -1
	} // else vertical

	// Does it fit?
	if endCol < 0 || endCol > s.numCols || endRow < 0 || endRow > s.numRows {
		/*
			log.Printf("validPlacement %s doesn't fit; %s start=%d,%d end=%d,%d", wordInfo.text,
				DirectionString(location.direction), startCol, startRow, endCol, endRow)
		*/
		return false
	}

	// Check if each coordinate is unused, or, if used, has the same
	// content that we need
	col := startCol
	row := startRow
	//	for i := 0; i < wordInfo.seqLen; i++ {
	//if !s.isCellAvailable(col, row) {
	for _, seqString := range wordInfo.seq.Items() {
		if !s.solutionMatrix.isCellAvailableFor(col, row, seqString) {
			// log.Printf("validPlacement %s overlap at col=%d row=%d", wordInfo.text, col, row)
			return false
		}
		col += colAdj
		row += rowAdj
	}

	// All good! Place it.
	wordInfo.placement = location

	// Mark the cells as used
	col = startCol
	row = startRow
	for _, seqString := range wordInfo.seq.Items() {
		s.solutionMatrix.setCellUsedFor(col, row, seqString)
		col += colAdj
		row += rowAdj
	}

	//	log.Printf("validPlacement %s chosen, col=%d row=%d", wordInfo.text, startCol, startRow)
	return true
}

func (s *wordSearchConstructor) clearPlacement(wordInfo *wordInfo) {
	location := wordInfo.placement
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

	// Mark the cells as avaialble
	col := location.col
	row := location.row
	//	log.Printf("clearPlacement %s at col=%d row=%d", wordInfo.text, col, row)
	for i := 0; i < wordInfo.seqLen; i++ {
		s.solutionMatrix.setCellAvailable(col, row)
		col += colAdj
		row += rowAdj
	}
}
