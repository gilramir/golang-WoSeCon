// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	wosecon "github.com/gilramir/golang-WoSeCon/v2"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func (s *Program) run() error {

	words, err := s.read_input_words()
	if err != nil {
		return err
	}

	wsOptions := make([]wosecon.WordSearchOption, 1)

	if s.WeightsFilename == "" {
		wsOptions[0] = wosecon.FillUniformlyFromString(alphabet)
	} else {
		// Read the weights file
		fh, err := os.Open(s.WeightsFilename)
		if err != nil {
			return err
		}
		fillerWeights := make([]wosecon.FillerWeight, 0)

		scanner := bufio.NewScanner(fh)
		lineno := 0
		for scanner.Scan() {
			lineno++
			text := scanner.Text()
			fields := strings.Fields(text)
			if len(fields) != 2 {
				continue
			}
			character, weightText := fields[0], fields[1]
			weight, err := strconv.ParseInt(weightText, 10, 64)
			if err != nil {
				fh.Close()
				return fmt.Errorf("In %s on line %d: %w", s.WeightsFilename, lineno, err)
			}
			fillerWeights = append(fillerWeights, wosecon.FillerWeight{
				String: character,
				Weight: weight,
			})
		}
		err = scanner.Err()
		if err != nil {
			fh.Close()
			return err
		}
		fh.Close()

		wsOptions[0] = wosecon.FillWeighted(fillerWeights)
	}

	ws, err := wosecon.Construct(s.NumCols, s.NumRows, words, wsOptions...)
	if err != nil {
		return err
	}

	for _, word := range words {
		placement, has := ws.WordPlacements[word]
		if !has {
			return fmt.Errorf("%s is not in WordPlacements", word)
		}
		fmt.Printf("%s - col:%d row:%d direction:%s\n", word,
			placement.Col, placement.Row, placement.DirectionString())
	}
	fmt.Println()

	s.printRows(ws.SolutionRows)
	s.printRows(ws.PuzzleRows)

	return nil
}

func (s *Program) printRows(rows [][]string) {
	// Header
	fmt.Print("   ")
	for col := 0; col < len(rows[0]); col++ {
		fmt.Printf("%2d ", col)
	}
	fmt.Println()
	fmt.Println()

	// Body
	for row, rowSlice := range rows {
		fmt.Printf("%2d ", row)
		for _, cellString := range rowSlice {
			if cellString == "" {
				cellString = " "
			}
			fmt.Printf(" %s ", cellString)
		}
		fmt.Println()
		fmt.Println()
	}

	fmt.Println()
}

func (s *Program) read_input_words() ([]string, error) {
	var reader io.Reader
	var err error

	if s.WordListFilename == "-" {
		reader = os.Stdin
	} else {

		fh, err := os.Open(s.WordListFilename)
		if err != nil {
			return nil, err
		}
		defer fh.Close()
		reader = fh
	}

	words := make([]string, 0, 20)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		words = append(words, text)
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return words, nil
}
