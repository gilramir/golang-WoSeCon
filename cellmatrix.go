// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

// A matrix indicating if the cell is used or not
type cellMatrix []uint8

func newCellMatrix(numCols, numRows int) cellMatrix {
	numBytesPerRow := numCols / 8
	if numCols%8 > 0 {
		numBytesPerRow += 1
	}
	totalBytes := numBytesPerRow * numRows
	/*
		fmt.Printf("newCellMatrix c=%d r=%d bpr=%d total=%d\n",
			numCols, numRows, numBytesPerRow, totalBytes)
	*/
	return cellMatrix(make([]byte, totalBytes))
}

func (c cellMatrix) getIndexBitMask(col, row, numCols int) (int, uint8) {
	numBytesPerRow := numCols / 8
	if numCols%8 > 0 {
		numBytesPerRow += 1
	}

	colByte := col / 8
	colBit := col % 8

	index := (numBytesPerRow * row) + colByte

	bitMask := uint8(1) << colBit

	/*
		fmt.Printf("getIndexBitMask c=%d r=%d i=%d bitMask=%b\n",
			col, row, index, bitMask)
	*/
	return index, bitMask
}

func (s cellMatrix) isAvailable(col, row, numCols int) bool {
	index, bitMask := s.getIndexBitMask(col, row, numCols)
	targetByte := s[index]
	return targetByte&bitMask == 0
}

func (s cellMatrix) isUsed(col, row, numCols int) bool {
	return !s.isAvailable(col, row, numCols)
}

func (s cellMatrix) setAvailable(col, row, numCols int) {
	index, bitMask := s.getIndexBitMask(col, row, numCols)
	s[index] &= ^bitMask
}

func (s cellMatrix) setUsed(col, row, numCols int) {
	index, bitMask := s.getIndexBitMask(col, row, numCols)
	s[index] |= bitMask
}
