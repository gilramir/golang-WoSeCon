// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	. "gopkg.in/check.v1"
)

// TotalLetters and LetterDensity on the returned WordSearch should
// reflect the (deduplicated) input across a normal grid.
func (s *MySuite) TestTotalLettersAndLetterDensity(c *C) {
	words := []string{"cat", "dog", "bird", "fish"} // 3+3+4+4 = 14 cells
	ws, err := Construct(10, 10, words,
		FillUniformlyFromString("ABCDEFGH"),
		RandomSeed(0),
	)
	c.Assert(err, IsNil)
	c.Check(ws.TotalLetters, Equals, 14)
	c.Check(ws.LetterDensity, Equals, 14.0/100.0)
}

// TotalLetters must count the deduplicated word list — duplicates are
// removed in init(), and the result struct should report the post-dedup
// total.
func (s *MySuite) TestTotalLettersAfterDeduplication(c *C) {
	words := []string{"cat", "cat", "dog"} // dedup -> {cat, dog} = 6 cells
	ws, err := Construct(8, 8, words,
		FillUniformlyFromString("ABC"),
		RandomSeed(0),
	)
	c.Assert(err, IsNil)
	c.Check(ws.TotalLetters, Equals, 6)
	c.Check(ws.LetterDensity, Equals, 6.0/64.0)
}
