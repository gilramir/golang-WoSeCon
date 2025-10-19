// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"math/rand"

	. "github.com/gilramir/gocheck-extra"
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestRandomLocator01(c *C) {
	rng := rand.New(rand.NewSource(42))

	var locator randomLocator
	locator.initialize(2, 3, Down|Up, rng)

	// # of locations = # row * # cols * # directions
	c.Assert(len(locator.availableLocations), Equals, 12)
	c.Assert(locator.size(), Equals, 12)

	newdl := directedLocation{
		col:       10,
		row:       10,
		direction: LTRHorizontal,
	}
	locator.add(newdl)

	c.Assert(locator.size(), Equals, 13)

	// Remove it
	locator.remove(newdl)
	c.Assert(locator.size(), Equals, 12)

	// It's not there to remove
	c.Assert(func() { locator.remove(newdl) }, PanicMatches, "Should not reach.*")
}

func (s *MySuite) TestRandomLocator02(c *C) {
	rng := rand.New(rand.NewSource(42))

	var locator randomLocator
	locator.initialize(2, 3, Down|Up, rng)

	// # of locations = # row * # cols * # directions
	c.Assert(locator.size(), Equals, 12)

	targetdl := directedLocation{
		col:       1,
		row:       1,
		direction: Down,
	}

	// Now you see it
	c.Assert(targetdl, InSlice, locator.availableLocations)

	newLocator := locator.minus([]directedLocation{targetdl})
	c.Assert(newLocator.size(), Equals, 11)

	// Now you don't
	c.Assert(targetdl, Not(InSlice), newLocator.availableLocations)

}
