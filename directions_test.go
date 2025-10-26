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
