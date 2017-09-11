package fsa

import (
	"testing"
)

func fuzzyAccepts(dfa *FuzzySparseTableDFA, str string) (bool, int) {
	s := dfa.Initial(str)
	mink := dfa.MaxError() + 1
	var final bool
	for len(s) > 0 {
		s = dfa.Delta(s, str, func(k, pos int, data uint32) {
			if k < mink {
				mink = k
			}
			final = true
		})
	}
	return final, mink
}

func TestEmptyFuzzyDFA(t *testing.T) {
	tests := []struct {
		name, test string
	}{
		{"empty", ""},
		{"non-empty", "non-empty-string"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := NewFuzzySparseTableDFA(3, new(SparseTableDFA))
			final, _ := fuzzyAccepts(dfa, tc.test)
			if final {
				t.Fatalf("empty DFA should not accept %q", tc.test)
			}
		})
	}
}
