// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	. "gopkg.in/check.v1"
)

func (s *MySuite) runCellMatrixTest(c *C, numCols, numRows int) {
	cm := newCellMatrix(numCols, numRows)
	for col := 0; col < numCols; col++ {
		for row := 0; row < numRows; row++ {
			c.Assert(cm.isAvailable(col, row, numCols), Equals, true)
		}
	}
}

// Less than 8 columns
func (s *MySuite) TestCellMatrix5xN(c *C) {
	s.runCellMatrixTest(c, 5, 3)
	s.runCellMatrixTest(c, 5, 8)
	s.runCellMatrixTest(c, 5, 25)
}

// Exactly 8 columns
func (s *MySuite) TestCellMatrix8xN(c *C) {
	s.runCellMatrixTest(c, 8, 3)
	s.runCellMatrixTest(c, 8, 8)
	s.runCellMatrixTest(c, 8, 25)
}

// More than 8, less than 16 columns
func (s *MySuite) TestCellMatrix10lN(c *C) {
	s.runCellMatrixTest(c, 10, 3)
	s.runCellMatrixTest(c, 10, 8)
	s.runCellMatrixTest(c, 10, 25)
}
