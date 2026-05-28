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

// Cap on the number of memoized dead-end states. At ~200 bytes/entry
// (fingerprint string + map overhead) this bounds the dead-end cache to
// roughly 20 MB. Once the cap is reached, further failures stop being
// recorded; existing entries keep serving lookups. The cache mainly pays
// off when distinct backtrack paths reach the same cell state at the
// same word index — for word lists with no placement-equivalent entries
// that rarely happens, so an aggressive cap costs little in practice.
const deadEndCacheCap = 100000

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

	// Memoized dead-end states. Key is (cell-state fingerprint, current
	// word index): a state from which no completion of the remaining
	// words is possible. Sound for word lists with no palindromic /
	// otherwise placement-equivalent entries — when two paths reach the
	// same cell state at the same word index, they have the same
	// available locator positions and hence the same set of reachable
	// completions.
	deadEnds map[deadEndKey]bool
}

type deadEndKey struct {
	fingerprint string
	wordIndex   int
}

func (s *wordSearchConstructor) init(numCols int, numRows int, sequences []Sequence, opts ...WordSearchOption) error {

	s.numCols = numCols
	s.numRows = numRows
	s.solutionMatrix.initialize(numCols, numRows)
	s.deadEnds = make(map[deadEndKey]bool)

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

		if s.tryPlace(currentWord, currentWordIndex) {
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

// tryPlace wraps locateOne with the dead-end memoization. The cell-state
// fingerprint at this point reflects words 0..currentWordIndex-1 placed and
// the current word not yet placed (in backward mode, clearPlacement on the
// previous iteration restored this state before we returned to the loop
// top). If a prior locateOne already proved this (state, index) yields no
// completion, return immediately; otherwise record the failure for later.
func (s *wordSearchConstructor) tryPlace(currentWord *wordInfo, currentWordIndex int) bool {
	key := deadEndKey{s.solutionMatrix.fingerprint(), currentWordIndex}
	if s.deadEnds[key] {
		return false
	}
	if s.locateOne(currentWord, currentWordIndex) {
		return true
	}
	if len(s.deadEnds) < deadEndCacheCap {
		s.deadEnds[key] = true
	}
	return false
}

// Place one word into the puzzle. If it cannot, returns false.
// currentWordIndex is the position of currentWord in s.wordInfos; it is used
// by the forward-check that ensures every still-unplaced word has at least
// one fitting location before the search commits to this placement.
func (s *wordSearchConstructor) locateOne(currentWord *wordInfo, currentWordIndex int) bool {
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
			// Forward-check: if any still-unplaced word now has zero
			// fitting locations, this placement is a dead end. Undo it,
			// record the location as tested so we don't re-evaluate it
			// on a subsequent backward-mode visit to this word, and
			// continue searching for another placement.
			if !s.remainingHaveFit(currentWordIndex) {
				s.clearPlacement(currentWord)
				currentWord.testedLocations = append(currentWord.testedLocations, suitableLocation)
				currentWord.placement = nilDirectedLocation()
				continue
			}
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

// canFit reports whether wordInfo would fit at location given the current
// solution matrix. It is a pure read-only feasibility check — no state is
// mutated. Used by the forward-check in locateOne.
func (s *wordSearchConstructor) canFit(wordInfo *wordInfo, location directedLocation) bool {
	var endCol, endRow, colAdj, rowAdj int

	if location.direction&GoesDownward > 0 {
		endRow = location.row + wordInfo.seqLen
		rowAdj = 1
	} else if location.direction&GoesUpward > 0 {
		endRow = location.row - wordInfo.seqLen
		rowAdj = -1
	}

	if location.direction&GoesLTR != 0 {
		endCol = location.col + wordInfo.seqLen
		colAdj = 1
	} else if location.direction&GoesRTL != 0 {
		endCol = location.col - wordInfo.seqLen
		colAdj = -1
	}

	if endCol < -1 || endCol > s.numCols || endRow < -1 || endRow > s.numRows {
		return false
	}

	col := location.col
	row := location.row
	for _, seqString := range wordInfo.seq.Items() {
		if !s.solutionMatrix.isCellAvailableFor(col, row, seqString) {
			return false
		}
		col += colAdj
		row += rowAdj
	}
	return true
}

// remainingHaveFit reports whether every word with index > currentWordIndex
// still has at least one location where it could be placed given the current
// solution matrix. Because future placements can only consume more cells,
// a word with no fit now will never have a fit, so this is a sound prune.
func (s *wordSearchConstructor) remainingHaveFit(currentWordIndex int) bool {
nextWord:
	for i := currentWordIndex + 1; i < len(s.wordInfos); i++ {
		wi := s.wordInfos[i]
		for _, loc := range s.locator.availableLocations {
			if s.canFit(wi, loc) {
				continue nextWord
			}
		}
		return false
	}
	return true
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
	// For DOWN/LTR, endRow/endCol is "one past the last cell" (startRow+seqLen),
	// so out-of-grid is > numRows / > numCols.
	// For UP/RTL, endRow/endCol is "one before the last cell" (startRow-seqLen).
	// The last cell ends up at endRow+1 / endCol+1, so the lowest legal value
	// of endRow/endCol is -1 (last cell at index 0). Reject only when < -1.
	if endCol < -1 || endCol > s.numCols || endRow < -1 || endRow > s.numRows {
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
