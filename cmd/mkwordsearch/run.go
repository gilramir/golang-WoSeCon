package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	wosecon "github.com/gilramir/golang-WoSeCon"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func (s *Program) run() error {

	words, err := s.read_input_words()
	if err != nil {
		return err
	}

	ws, err := wosecon.Construct(s.NumCols, s.NumRows, words,
		wosecon.FillEquallyString(alphabet))
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

func (s *Program) printRows(rows [][]rune) {
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
		for _, cellRune := range rowSlice {
			fmt.Printf(" %s ", string(cellRune))
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
