// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	. "gopkg.in/check.v1"
)

// Exercise the dead-end memoization code path with the cache enabled.
// The default cap is 0 (disabled), so without this test the lookup,
// fingerprint, and cap-check branches in tryPlace are never executed
// by go test. The result should be a valid puzzle — turning the cache
// on must not change correctness, only performance characteristics.
func (s *MySuite) TestMemoizationLimitEnabledSolves(c *C) {
	words := []string{"cat", "dog", "bird", "fish"}
	ws, err := Construct(8, 8, words,
		FillUniformlyFromString("ABCDEFGH"),
		AddNaturalLTRDirections(),
		RandomSeed(0),
		WithMemoizationLimit(100),
	)
	c.Assert(err, IsNil)
	c.Check(len(ws.WordPlacements), Equals, 4)
}

// Enabling the cache must produce the same official solution as not
// enabling it for the same seed. The cache only short-circuits states
// the algorithm would have failed at anyway, so for a deterministic
// search seed the chosen placements should be identical.
func (s *MySuite) TestMemoizationDoesNotChangeResult(c *C) {
	words := []string{"cat", "dog", "bird", "fish"}
	base := []WordSearchOption{
		FillUniformlyFromString("ABCDEFGH"),
		AddNaturalLTRDirections(),
		RandomSeed(0),
	}

	noCache, err := Construct(8, 8, words, base...)
	c.Assert(err, IsNil)

	withCache, err := Construct(8, 8, words,
		append(append([]WordSearchOption{}, base...), WithMemoizationLimit(100))...,
	)
	c.Assert(err, IsNil)

	c.Check(withCache.WordPlacements, DeepEquals, noCache.WordPlacements)
}
