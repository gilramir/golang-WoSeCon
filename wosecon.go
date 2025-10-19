// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"math/rand"
	"sort"
	"unicode/utf8"
)

// The algorithm works in 2 modoes, forward or backward
type algoMode byte

const (
	forwardMode  algoMode = 'f'
	backwardMode algoMode = 'b'
)

type wordInfo struct {
	text            string
	runeLen         int
	placement       directedLocation
	testedLocations []directedLocation
}

func newWordInfo(text string) *wordInfo {
	return &wordInfo{
		text:            text,
		runeLen:         utf8.RuneCountInString(text),
		placement:       nilDirectedLocation(),
		testedLocations: make([]directedLocation, 0),
	}
}

func (s *wordInfo) getPlacement() directedLocation {
	return s.placement
}

func (s *wordInfo) getTested() []directedLocation {
	return s.testedLocations
}

func (s *wordInfo) moveLocationToTested() {
	s.testedLocations = append(s.testedLocations, s.placement)
	s.placement = nilDirectedLocation()
}

func (s *wordInfo) deleteTested() {
	s.testedLocations = nil
}

type directedLocation struct {
	col       int
	row       int
	direction Direction
}

func (s directedLocation) equals(target directedLocation) bool {
	return s.col == target.col && s.row == target.row && s.direction == target.direction
}

// Return a nil DirectedLocation, which has -1 col/row
func nilDirectedLocation() directedLocation {
	return directedLocation{-1, -1, 0}
}
func IsNilDirectedLocation(d directedLocation) bool {
	return d.col == -1 || d.row == -1
}

// XXX - can I make this private?
type WordSearchConstructor struct {
	numCols            int
	numRows            int
	wordInfos          []*wordInfo
	possibleDirections Direction
	mode               algoMode
	locator            randomLocator
	cellMatrix         cellMatrix
	//	badWords           []string

	randomSeed      int64
	randomSeedGiven bool
	rng             *rand.Rand

	fillerRunes []rune
}

func (s *WordSearchConstructor) init(numCols int, numRows int, words []string, opts ...WordSearchOption) error {

	s.numCols = numCols
	s.numRows = numRows
	s.cellMatrix = newCellMatrix(numCols, numRows)

	// Use the list of words (strings) from the user
	// to create our wordInfo objects
	s.wordInfos = make([]*wordInfo, len(words))

	for i, wordString := range words {
		s.wordInfos[i] = newWordInfo(wordString)
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
	if s.possibleDirections&(goesUp|goesDown) != 0 {
		maxWordSize = numRows
	}
	if s.possibleDirections&(goesLTR|goesRTL) != 0 {
		maxWordSize = max(maxWordSize, numCols)
	}

	for _, wordInfo := range s.wordInfos {
		if wordInfo.runeLen > maxWordSize {
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
	// If the same length, then in alphbetical order
	sort.Slice(s.wordInfos, func(i, j int) bool {
		wi := s.wordInfos[i]
		wj := s.wordInfos[j]
		if wi.runeLen == wj.runeLen {
			return wi.text < wj.text
		} else {
			// Descending order
			return wj.runeLen < wi.runeLen
		}
	})
	return nil
}

// The main loop of the algorithm
func (s *WordSearchConstructor) construct() error {

	s.locator.initialize(s.numCols, s.numRows, s.possibleDirections, s.rng)

	s.mode = forwardMode
	currentWordIndex := 0
	currentWord := s.wordInfos[currentWordIndex]
	for {
		if s.locateOne(currentWord) {
			if currentWordIndex == len(s.wordInfos)-1 {
				// Finished
				break
			}
			currentWordIndex++
			currentWord = s.wordInfos[currentWordIndex]
			s.mode = forwardMode
		} else {
			// Couldn't place a word. Go backwards
			if currentWordIndex == 0 {
				return ErrCannotFitWords
			} else {
				currentWord.deleteTested()
				currentWordIndex--
				currentWord = s.wordInfos[currentWordIndex]
				s.mode = backwardMode
			}
		}
	}
	return nil
}

// Place one word into the puzzle. If it cannot, returns false.
func (s *WordSearchConstructor) locateOne(currentWord *wordInfo) bool {
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
	return false
}

// Will this word fit in the puzzle at this directedLocation?
// If so, the wordInfo's placement is updated
func (s *WordSearchConstructor) validPlacement(wordInfo *wordInfo, location directedLocation) bool {

	/*
		log.Printf("validPlacement %s trying %d,%d,%s", wordInfo.text,
			location.col, location.row, DirectionString(location.direction))
	*/

	// Perchance did we already test this location?
	for _, testedLocation := range wordInfo.testedLocations {
		if testedLocation.equals(location) {
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
	if location.direction&goesDown > 0 {
		endRow = startRow + wordInfo.runeLen
		rowAdj = 1
	} else if location.direction&goesUp > 0 {
		endRow = startRow - wordInfo.runeLen
		rowAdj = -1
	} // else horizontal

	// Find endCol
	if location.direction&goesLTR != 0 {
		endCol = startCol + wordInfo.runeLen
		colAdj = 1
	} else if location.direction&goesRTL != 0 {
		endCol = startCol - wordInfo.runeLen
		colAdj = -1
	} // else vetical

	// Does it fit?
	if endCol < 0 || endCol > s.numCols || endRow < 0 || endRow > s.numRows {
		/*
			log.Printf("validPlacement %s doesn't fit; %s start=%d,%d end=%d,%d", wordInfo.text,
				DirectionString(location.direction), startCol, startRow, endCol, endRow)
		*/
		return false
	}

	// Check if each coordinate is unused
	col := startCol
	row := startRow
	for i := 0; i < wordInfo.runeLen; i++ {
		if !s.isCellAvailable(col, row) {
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
	for i := 0; i < wordInfo.runeLen; i++ {
		s.setCellUsed(col, row)
		col += colAdj
		row += rowAdj
	}

	//log.Printf("validPlacement %s chosen, col=%d row=%d", wordInfo.text, startCol, startRow)
	return true
}

func (s *WordSearchConstructor) isCellAvailable(col, row int) bool {
	return s.cellMatrix.isAvailable(col, row, s.numCols)
}

func (s *WordSearchConstructor) setCellUsed(col, row int) {
	s.cellMatrix.setUsed(col, row, s.numCols)
}
