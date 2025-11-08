package wosecon

import (
	"iter"
)

type Sequence interface {
	String() string
	Len() int
	Cmp(Sequence) int
	Index(int) string
	Items() iter.Seq2[int, string]
}

// RuneSequence is used for the most
// common type of word search; it's for strings
type RuneSequence struct {
	runes []rune
}

func NewRuneSequence(text string) RuneSequence {
	return RuneSequence{
		runes: []rune(text),
	}
}

func (s RuneSequence) String() string {
	return string(s.runes)
}

func (s RuneSequence) Len() int {
	return len(s.runes)
}

func (s RuneSequence) Cmp(other Sequence) int {
	m := string(s.runes)
	o := other.String()
	if m < o {
		return -1
	} else if m > o {
		return 1
	} else {
		return 0
	}
}

func (s RuneSequence) Index(n int) string {
	if n < len(s.runes) {
		return string(s.runes[n])
	} else {
		return ""
	}
}
func (s RuneSequence) Items() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i, r := range s.runes {
			if !yield(i, string(r)) {
				return
			}
		}
	}
}
