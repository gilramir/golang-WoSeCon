package wosecon

import "fmt"

type solutionMatrix struct {
	numCols int
	numRows int
	// Dimension 1: Column (x)
	// Dimension 2: Row (y)
	matrix [][]*solutionCell
}

type solutionCell struct {
	seqString string
	count     int
}

func (s *solutionMatrix) initialize(numCols int, numRows int) {
	s.numCols = numCols
	s.numRows = numRows

	s.matrix = make([][]*solutionCell, s.numCols)

	for col := 0; col < s.numCols; col++ {
		s.matrix[col] = make([]*solutionCell, s.numRows)

		for row := 0; row < s.numRows; row++ {
			s.matrix[col][row] = new(solutionCell)
		}
	}
}

func (s *solutionMatrix) isCellEmpty(col, row int) bool {
	cell := s.matrix[col][row]
	return cell.count == 0
}

func (s *solutionMatrix) isCellAvailableFor(col, row int, text string) bool {
	cell := s.matrix[col][row]
	if cell.count == 0 {
		return true
	} else {
		return cell.seqString == text
	}
}

func (s *solutionMatrix) setCellUsedFor(col, row int, text string) {
	cell := s.matrix[col][row]
	if cell.count == 0 {
		cell.seqString = text
		cell.count = 1
	} else {
		if cell.seqString == text {
			cell.count++
		} else {
			panic(fmt.Sprintf("cell (%d, %d) has %s but we tried to set %s",
				col, row, cell.seqString, text))
		}
	}
}

func (s *solutionMatrix) setCellAvailable(col, row int) {
	cell := s.matrix[col][row]
	if cell.count < 1 {
		panic(fmt.Sprintf("can't reduce cell (%d, %d) any further; at %d",
			col, row, cell.count))
	}
	cell.count--
	if cell.count == 0 {
		cell.seqString = ""
	}
}

func (s *solutionMatrix) dump() {
	// Header
	fmt.Print("   ")
	for col := 0; col < s.numCols; col++ {
		fmt.Printf("%2d ", col)
	}
	fmt.Println()
	fmt.Println()

	// Body
	for row := 0; row < s.numRows; row++ {
		fmt.Printf("%2d ", row)
		for col := 0; col < s.numCols; col++ {
			cell := s.matrix[col][row]
			fmt.Printf(" %1s ", cell.seqString)
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Println()
}
