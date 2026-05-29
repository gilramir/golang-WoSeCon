// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	. "github.com/gilramir/gocheck-extra"
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestDirectionsNaturalLTR(c *C) {
	var d Direction = NaturalLTRDirections

	directions := DirectionSlice(d)
	c.Assert(len(directions), Equals, 4)

	// Has these
	c.Check(Down, InSlice, directions)
	c.Check(LTRHorizontal, InSlice, directions)
	c.Check(LTRDescending, InSlice, directions)
	c.Check(LTRAscending, InSlice, directions)

	// Does not have these
	c.Check(Up, Not(InSlice), directions)
	c.Check(RTLHorizontal, Not(InSlice), directions)
	c.Check(RTLDescending, Not(InSlice), directions)
	c.Check(RTLAscending, Not(InSlice), directions)
}

func (s *MySuite) TestDirectionsAll(c *C) {
	var d Direction = AllDirections

	directions := DirectionSlice(d)
	c.Assert(len(directions), Equals, 8)

	// Has these
	c.Check(Down, InSlice, directions)
	c.Check(LTRHorizontal, InSlice, directions)
	c.Check(LTRDescending, InSlice, directions)
	c.Check(LTRAscending, InSlice, directions)

	// And these too
	c.Check(Up, InSlice, directions)
	c.Check(RTLHorizontal, InSlice, directions)
	c.Check(RTLDescending, InSlice, directions)
	c.Check(RTLAscending, InSlice, directions)
}

// Regression test for a validPlacement boundary off-by-one. With Up,
// LTRAscending, RTLAscending, RTLHorizontal, RTLDescending — anything
// whose endRow/endCol is one *before* the first cell rather than one
// past the last — the old check rejected placements whose last cell
// landed on row 0 or column 0. In a 3x3 grid with only LTRAscending
// enabled, the single valid placement for a 3-letter word is the
// diagonal (col=0, row=2) -> (col=2, row=0). The buggy check rejected
// it because endRow = -1 < 0; with the fix the comparison is < -1,
// which lets endRow = -1 (last cell at row 0) through.
func (s *MySuite) TestAscendingPlacementAtRowZero(c *C) {
	words := []string{"abc"}
	ws, err := Construct(3, 3, words,
		AddDirection(LTRAscending),
		RandomSeed(0),
	)
	c.Assert(err, IsNil)
	c.Check(ws.WordPlacements["abc"], Equals, WordPlacement{
		Col: 0, Row: 2, Direction: LTRAscending,
	})
}
