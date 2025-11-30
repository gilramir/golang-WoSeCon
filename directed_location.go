package wosecon

type directedLocation struct {
	col       int
	row       int
	direction Direction
}

// XXX don't need this; Go's == will work.
func (s directedLocation) equals(target directedLocation) bool {
	return s.col == target.col && s.row == target.row && s.direction == target.direction
}

// Return a nil DirectedLocation, which has -1 col/row
func nilDirectedLocation() directedLocation {
	return directedLocation{-1, -1, 0}
}
func isNilDirectedLocation(d directedLocation) bool {
	return d.col == -1 || d.row == -1
}

type directedLocationMatrix struct {
	numCols int
	numRows int
	// Dimension 1: Column (x)
	// Dimension 2: Row (y)
	// Dimension 3: the possible directedLocations
	dlmatrix [][][]directedLocation

	numDirectedLocations int
}

func (s *directedLocationMatrix) getDirectionsAt(col int, row int) []directedLocation {
	return s.dlmatrix[col][row]
}

func (s *directedLocationMatrix) initialize(numCols int, numRows int, possibleDirections Direction) {
	s.numCols = numCols
	s.numRows = numRows

	possibleDirectionsSlice := DirectionSlice(possibleDirections)

	s.dlmatrix = make([][][]directedLocation, s.numCols)

	for col := 0; col < s.numCols; col++ {
		s.dlmatrix[col] = make([][]directedLocation, s.numRows)
	}

	// Allocate possibleDirections for the the non-border cells
	for col := 1; col < numCols-1; col++ {
		for row := 1; row < numRows-1; row++ {
			directedLocations := make([]directedLocation, len(possibleDirectionsSlice))
			for di, direction := range possibleDirectionsSlice {
				d := directedLocation{
					col:       col,
					row:       row,
					direction: direction,
				}
				directedLocations[di] = d
			}
			s.dlmatrix[col][row] = directedLocations
			s.numDirectedLocations += len(directedLocations)
		}
	}

	// Make masks to mask down certain directions
	goingUpMask := ^(Up | LTRAscending | RTLAscending)
	goingRightMask := ^(LTRHorizontal | LTRAscending | LTRDescending)
	goingLeftMask := ^(RTLHorizontal | RTLAscending | RTLDescending)
	goingDownMask := ^(Down | LTRDescending | RTLDescending)

	// Top row; cannot ascend or go up
	topRowDirections := possibleDirections & goingUpMask
	topRowDirectionsSlice := DirectionSlice(topRowDirections)
	for col := 0; col < numCols; col++ {
		var directionsSlice []Direction
		// Corners have more restrictions
		if col == 0 {
			directionsSlice = DirectionSlice(topRowDirections & goingLeftMask)
		} else if col == numCols-1 {
			directionsSlice = DirectionSlice(topRowDirections & goingRightMask)
		} else {
			directionsSlice = topRowDirectionsSlice
		}
		directedLocations := make([]directedLocation, 0)
		for _, direction := range directionsSlice {
			d := directedLocation{
				col:       col,
				row:       0,
				direction: direction,
			}
			directedLocations = append(directedLocations, d)
		}
		s.dlmatrix[col][0] = directedLocations
		s.numDirectedLocations += len(directedLocations)
	}

	// Bottom row; cannot descend or go down
	bottomRowDirections := possibleDirections & goingDownMask
	bottomRowDirectionsSlice := DirectionSlice(bottomRowDirections)
	for col := 0; col < numCols; col++ {
		var directionsSlice []Direction
		// Corners have more restrictions
		if col == 0 {
			directionsSlice = DirectionSlice(bottomRowDirections & goingLeftMask)
		} else if col == numCols-1 {
			directionsSlice = DirectionSlice(bottomRowDirections & goingRightMask)
		} else {
			directionsSlice = bottomRowDirectionsSlice
		}
		directedLocations := make([]directedLocation, 0)
		for _, direction := range directionsSlice {
			d := directedLocation{
				col:       col,
				row:       numRows - 1,
				direction: direction,
			}
			directedLocations = append(directedLocations, d)
		}
		s.dlmatrix[col][s.numRows-1] = directedLocations
		s.numDirectedLocations += len(directedLocations)
	}

	// Left column; cannot go left. We already did the top and bottom
	// corners, so we can skip them
	leftColDirectionsSlice := DirectionSlice(possibleDirections & goingLeftMask)
	for row := 1; row < numRows-1; row++ {
		directedLocations := make([]directedLocation, 0)
		for _, direction := range leftColDirectionsSlice {
			d := directedLocation{
				col:       0,
				row:       row,
				direction: direction,
			}
			directedLocations = append(directedLocations, d)
		}
		s.dlmatrix[0][row] = directedLocations
		s.numDirectedLocations += len(directedLocations)
	}

	// Right column; cannot go right. We already did the top and bottom
	// corners, so we can skip them
	rightColDirectionsSlice := DirectionSlice(possibleDirections & goingRightMask)
	for row := 1; row < numRows-1; row++ {
		directedLocations := make([]directedLocation, 0)
		for _, direction := range rightColDirectionsSlice {
			d := directedLocation{
				col:       numCols - 1,
				row:       row,
				direction: direction,
			}
			directedLocations = append(directedLocations, d)
		}
		s.dlmatrix[s.numCols-1][row] = directedLocations
		s.numDirectedLocations += len(directedLocations)
	}
	/*
		fmt.Printf("initialized dlmatrix; #dl=%d lendlm=%d\n", s.numDirectedLocations,
			len(s.dlmatrix))
	*/
}

func (s *directedLocationMatrix) getAllDirectedLocations() []directedLocation {

	allDirectedLocations := make([]directedLocation, 0, s.numDirectedLocations)

	//	fmt.Printf("len dlmatrix=%d\n", len(s.dlmatrix))
	// We could do it easily like this...
	/*
		for _, row := range s.dlmatrix {
			for _, cell := range row {
				//			fmt.Printf("col=%d row=%d -> %v", coli, rowi, cell)
				allDirectedLocations = append(allDirectedLocations, cell...)
			}
		}
	*/

	// but we do it like this since this was the original order, and some
	// unit tests have RNG seeds and depend on this order
	for col := 1; col < s.numCols-1; col++ {
		for row := 1; row < s.numRows-1; row++ {
			allDirectedLocations = append(allDirectedLocations, s.dlmatrix[col][row]...)
		}
	}

	for col := 0; col < s.numCols; col++ {
		allDirectedLocations = append(allDirectedLocations, s.dlmatrix[col][0]...)
	}

	for col := 0; col < s.numCols; col++ {
		allDirectedLocations = append(allDirectedLocations, s.dlmatrix[col][s.numRows-1]...)
	}

	for row := 1; row < s.numRows-1; row++ {
		allDirectedLocations = append(allDirectedLocations, s.dlmatrix[0][row]...)
	}
	for row := 1; row < s.numRows-1; row++ {
		allDirectedLocations = append(allDirectedLocations, s.dlmatrix[s.numCols-1][row]...)
	}

	//	fmt.Printf("returned alldl=%d\n", len(allDirectedLocations))
	return allDirectedLocations

}
