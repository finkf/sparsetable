package fsa

import (
	"bytes"
	"sort"
	"testing"
)

var teststrs = []string{
	"",
	"abcde",
	"very-long-string",
	"some-non-ascii-chars-ü-äåÅßß",
	"abcabc", // common suffixes
	"ddeabc",
	"floabc",
}

func accepts(dfa *SparseTableDFA, str string) bool {
	s := dfa.Initial()
	for i := 0; i < len(str) && s != 0; i++ {
		s = dfa.Delta(s, str[i])
	}
	_, final := dfa.Final(s)
	return final
}

func TestEmptySparseTableDFA(t *testing.T) {
	dfa := &SparseTableDFA{}
	for i, str := range teststrs {
		if accepts(dfa, str) {
			t.Errorf("[%d] dfa accepts %q", i, str)
		}
	}
}

func TestSingleEntrySparseTableDFA(t *testing.T) {
	for i, str := range teststrs {
		b := NewSparseTableDFABuilder()
		if err := b.Add(str, 1); err != nil {
			t.Errorf("[%d] error: %q", i, err)
		}
		dfa := b.Build()
		for _, test := range teststrs {
			if accepts(dfa, test) && test != str {
				t.Errorf("[%d] dfa accepts %q", i, test)
			}
			if !accepts(dfa, test) && test == str {
				t.Errorf("[%d] dfa does not accept %q", i, test)
				t.Errorf("[%d] %v", i, dfa)
			}
		}
	}
}

func TestSparseTableDFA(t *testing.T) {
	b := NewSparseTableDFABuilder()
	sorted := make([]string, len(teststrs))
	copy(sorted, teststrs)
	sort.Slice(sorted, func(i, j int) bool {
		return bytes.Compare([]byte(sorted[i]), []byte(sorted[j])) < 0
	})
	for i, str := range sorted {
		if err := b.Add(str, 1); err != nil {
			t.Fatalf("[%d] could not add %q: %s", i, str, err)
		}
	}
	dfa := b.Build()
	for i, test := range teststrs {
		if !accepts(dfa, test) {
			t.Errorf("[%d] dfa does not accept %q", i, test)
		}
	}
}

func TestEachTransition(t *testing.T) {
	b := NewSparseTableDFABuilder()
	sorted := make([]string, len(teststrs))
	copy(sorted, teststrs)
	sort.Slice(sorted, func(i, j int) bool {
		return bytes.Compare([]byte(sorted[i]), []byte(sorted[j])) < 0
	})
	for i, str := range sorted {
		if err := b.Add(str, 1); err != nil {
			t.Errorf("[%d] could not add %q: %s", i, str, err)
		}
	}
	dfa := b.Build()
	chars := make(map[byte]bool)
	dfa.EachTransition(dfa.Initial(), func(cell Cell) {
		if cell.typ != transitionCellType {
			t.Errorf("expected transition cell; got %s", cell)
		}
		chars[cell.char] = true
	})
	if len(chars) != 5 {
		t.Errorf("expected 5 transitions; got %d", len(chars))
	}
	for _, c := range []byte{'a', 'v', 'd', 'f', 's'} {
		if !chars[c] {
			t.Errorf("expected chars to contain %c", c)
		}
	}
}
