package sparsetable

import (
	"testing"
)

func fuzzyAccepts(dfa *FuzzyDFA, str string) (bool, int) {
	s := dfa.Initial(str)
	mink := dfa.MaxError() + 1
	var final bool
	for len(s) > 0 {
		s = dfa.Delta(s, func(k, pos int, data int32) {
			if pos != len(str) {
				return
			}
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
			dfa := NewFuzzyDFA(3, new(DFA))
			final, _ := fuzzyAccepts(dfa, tc.test)
			if final {
				t.Fatalf("empty DFA should not accept %q", tc.test)
			}
		})
	}
}

func TestSingleEntryFuzzyDFA(t *testing.T) {
	tests := []struct {
		name, entry, test string
		k                 int
		accept            bool
	}{
		{"empty with k=0", "", "", 0, true},
		{"empty with k=1", "", "a", 1, true},
		{"empty with k=2", "", "aa", 2, true},
		{"empty with k=3", "", "aaa", 3, true},
		{"empty with k>3", "", "aaaa", 0, false},
		{"a with k=0", "a", "a", 0, true},
		{"a with k=1", "a", "xa", 1, true},
		{"a with k=1", "a", "ax", 1, true},
		{"a with k=2", "a", "xxa", 2, true},
		{"a with k=2", "a", "xax", 2, true},
		{"a with k=2", "a", "axx", 2, true},
		{"a with k=3", "a", "axxx", 3, true},
		{"a with k=3", "a", "xaxx", 3, true},
		{"a with k=3", "a", "xxax", 3, true},
		{"a with k=3", "a", "axxx", 3, true},
		{"a with k>3", "a", "axxxx", 0, false},
		{"a with k>3", "a", "xaxxx", 0, false},
		{"a with k>3", "a", "xxaxx", 0, false},
		{"a with k>3", "a", "xaxxx", 0, false},
		{"a with k>3", "a", "axxxx", 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := NewFuzzyDFA(3, NewDictionary(tc.entry))
			final, k := fuzzyAccepts(dfa, tc.test)
			if final != tc.accept {
				t.Fatalf("expected accept(%q) = %t; got %t", tc.test, tc.accept, final)
			}
			if final && tc.k != k { // test only for accepted strings
				t.Fatalf("expected accept(%q) = %d; got %d", tc.test, tc.k, k)
			}
		})
	}
}
