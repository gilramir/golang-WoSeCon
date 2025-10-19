package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	wosecon "github.com/gilramir/golang-WoSeCon"
)

func (s *Program) run() error {

	words, err := s.read_input_words()
	if err != nil {
		return err
	}

	ws, err := wosecon.Construct(s.NumCols, s.NumRows, words)
	if err != nil {
		return err
	}

	// Header
	fmt.Print("   ")
	for col := 0; col < ws.NumCols; col++ {
		fmt.Printf("%2d ", col)
	}
	fmt.Println()
	fmt.Println()

	// Body
	for row, rowSlice := range ws.SolutionRows {
		fmt.Printf("%2d ", row)
		for _, cellRune := range rowSlice {
			fmt.Printf(" %s ", string(cellRune))
		}
		fmt.Println()
		fmt.Println()
	}

	fmt.Println()

	hasError := false
	for _, word := range words {
		placement, has := ws.WordPlacements[word]
		if !has {
			fmt.Printf("ERROR: %s is not in WordPlacements\n", word)
			hasError = true
		}
		fmt.Printf("%s - col:%d row:%d direction:%s\n", word,
			placement.Col, placement.Row, placement.DirectionString())
	}

	if hasError {
		return errors.New("Invalid wordSearchConstructor data")
	}

	return nil
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
