// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"time"

	. "gopkg.in/check.v1"
)

// WithParallelism(n > 1) combined with RandomSeed should fail before any
// search runs.
func (s *MySuite) TestParallelismRejectsRandomSeed(c *C) {
	words := []string{"cat", "dog", "bird"}
	_, err := Construct(8, 8, words,
		FillUniformlyFromString("ABC"),
		AddNaturalLTRDirections(),
		WithParallelism(4),
		RandomSeed(0),
	)
	c.Check(err, Equals, ErrSeedWithParallelism)
}

// Reversed option order should still fail — neither option implicitly
// "wins" over the other.
func (s *MySuite) TestParallelismRejectsRandomSeedReversedOrder(c *C) {
	words := []string{"cat", "dog", "bird"}
	_, err := Construct(8, 8, words,
		FillUniformlyFromString("ABC"),
		AddNaturalLTRDirections(),
		RandomSeed(0),
		WithParallelism(4),
	)
	c.Check(err, Equals, ErrSeedWithParallelism)
}

// WithParallelism(1) is equivalent to leaving the option off; a single
// seed is allowed alongside it.
func (s *MySuite) TestParallelismOneAllowsRandomSeed(c *C) {
	words := []string{"cat", "dog", "bird"}
	ws, err := Construct(8, 8, words,
		FillUniformlyFromString("ABC"),
		AddNaturalLTRDirections(),
		RandomSeed(0),
		WithParallelism(1),
	)
	c.Assert(err, IsNil)
	c.Check(ws.NumCols, Equals, 8)
	c.Check(ws.NumRows, Equals, 8)
}

// Construct should solve when WithParallelism is used on a workable
// puzzle. We don't assert specific placements because each goroutine
// picks a different seed.
func (s *MySuite) TestParallelismSolvesEasyPuzzle(c *C) {
	words := []string{"cat", "dog", "bird", "fish"}
	ws, err := Construct(10, 10, words,
		FillUniformlyFromString("ABCDEFGH"),
		AddNaturalLTRDirections(),
		WithParallelism(4),
		WithTimeLimit(5*time.Second),
	)
	c.Assert(err, IsNil)
	c.Check(len(ws.WordPlacements), Equals, 4)
}
