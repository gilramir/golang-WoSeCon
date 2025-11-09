package wosecon

import (
	"iter"
	"slices"
	"strings"
)

type Sequence interface {
	String() string
	Len() int
	Cmp(Sequence) int
	Index(int) string
	Items() iter.Seq2[int, string]
}

// RuneSequence is used for the most
// common type of word search; it's for strings where each
// rune in the string is one cell in the puzzle
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

// StringSequence is a special case. It's for strings
// where more than one rune in the string is one cell in
// the puzzle. It's useful if you're creating a word search
// puzzle for a language like Thai, where a single written "letter"
// might have 1, 2, or 3 runes combined together.
type StringSequence struct {
	parts []string
}

func NewStringSequence(parts []string) StringSequence {
	return StringSequence{
		parts: slices.Clone(parts),
	}
}

func (s StringSequence) String() string {
	return strings.Join(s.parts, "")
}

func (s StringSequence) Len() int {
	return len(s.parts)
}

func (s StringSequence) Cmp(other Sequence) int {
	m := s.String()
	o := other.String()
	if m < o {
		return -1
	} else if m > o {
		return 1
	} else {
		return 0
	}
}

func (s StringSequence) Index(n int) string {
	if n < len(s.parts) {
		return string(s.parts[n])
	} else {
		return ""
	}
}
func (s StringSequence) Items() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i, p := range s.parts {
			if !yield(i, p) {
				return
			}
		}
	}
}
