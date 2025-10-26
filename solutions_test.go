// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"fmt"

	//	. "github.com/gilramir/gocheck-extra"
	. "gopkg.in/check.v1"
)

func printRows(rows [][]rune) {
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

// Simple case
func (s *MySuite) TestSolutions01(c *C) {

	words := []string{"cat", "car"}
	filler := "ca"
	ws, apiError := Construct(4, 4, words,
		FillUniformlyFromString(filler),
		UseAllDirections(),
		RandomSeed(42))
	c.Assert(apiError, IsNil)

	/*
		printRows(ws.PuzzleRows)
		printRows(ws.SolutionRows)
	*/

	/*
	      Puzzle
	      0  1  2  3

	   0  a  a  a  c

	   1  c  c  c  c

	   2  c  a  r  a

	   3  c  c  a  t
	*/
	/*
	      Solution
	      0  1  2  3

	   0

	   1

	   2  c  a  r

	   3     c  a  t
	*/

	//	fmt.Printf("AllSolutions: %v\n", ws.AllPossibleWordPlacements)

	carSolutions := ws.AllPossibleWordPlacements["car"]
	c.Check(len(carSolutions), Equals, 1)

	catSolutions := ws.AllPossibleWordPlacements["cat"]
	c.Check(len(catSolutions), Equals, 2)

	c.Check(ws.NumWordsWithMultipleSolutions, Equals, 1)
}

// Extreme case
func (s *MySuite) TestSolutions02(c *C) {

	words := []string{"c"}
	filler := "c"
	ws, apiError := Construct(4, 4, words,
		FillUniformlyFromString(filler),
		UseAllDirections(),
		RandomSeed(42))
	c.Assert(apiError, IsNil)
	//	printRows(ws.SolutionRows)

	solutions := ws.AllPossibleWordPlacements["c"]
	c.Check(len(solutions), Equals, 16)

	for _, placement := range solutions {
		c.Check(placement.Direction, Equals, NilDirection)
	}
	c.Check(ws.NumWordsWithMultipleSolutions, Equals, 1)
}

// Simple case
func (s *MySuite) TestSolutions03(c *C) {

	words := []string{"cat", "car"}
	filler := "cart"
	ws, apiError := Construct(4, 4, words,
		FillUniformlyFromString(filler),
		UseAllDirections(),
		RandomSeed(1))
	c.Assert(apiError, IsNil)

	/*
		printRows(ws.PuzzleRows)
		printRows(ws.SolutionRows)
	*/

	/* Puzzle
	      0  1  2  3

	   0  a  a  t  r

	   1  c  c  a  r

	   2  a  a  t  c

	   3  t  r  t  r
	*/

	/* Solution

	      0  1  2  3

	   0

	   1  c  c

	   2  a  a

	   3  t  r
	*/

	carSolutions := ws.AllPossibleWordPlacements["car"]
	c.Check(len(carSolutions), Equals, 2)
	c.Check(carSolutions[0], Equals, WordPlacement{
		Row:       1,
		Col:       1,
		Direction: Down,
	})
	c.Check(carSolutions[1], Equals, WordPlacement{
		Row:       1,
		Col:       1,
		Direction: LTRHorizontal,
	})

	catSolutions := ws.AllPossibleWordPlacements["cat"]
	c.Check(len(catSolutions), Equals, 2)
	c.Check(catSolutions[0], Equals, WordPlacement{
		Row:       1,
		Col:       0,
		Direction: Down,
	})
	c.Check(catSolutions[1], Equals, WordPlacement{
		Row:       1,
		Col:       0,
		Direction: LTRDescending,
	})

	c.Check(ws.NumWordsWithMultipleSolutions, Equals, 2)
}
