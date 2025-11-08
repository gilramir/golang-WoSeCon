package wosecon

type wordInfo struct {
	seq             Sequence
	seqLen          int
	placement       directedLocation
	testedLocations []directedLocation
}

func newWordInfo(seq Sequence) *wordInfo {
	return &wordInfo{
		seq:             seq,
		seqLen:          seq.Len(),
		placement:       nilDirectedLocation(),
		testedLocations: make([]directedLocation, 0),
	}
}

func (s *wordInfo) getPlacement() directedLocation {
	return s.placement
}

func (s *wordInfo) getTested() []directedLocation {
	return s.testedLocations
}

func (s *wordInfo) moveLocationToTested() {
	s.testedLocations = append(s.testedLocations, s.placement)
	s.placement = nilDirectedLocation()
}

func (s *wordInfo) deleteTested() {
	s.testedLocations = nil
}

func (s *wordInfo) length() int {
	return s.seqLen
}

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
