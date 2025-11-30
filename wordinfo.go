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
