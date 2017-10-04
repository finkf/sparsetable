package sparsetable

import (
	"testing"
)

func fuzzyAccepts(dfa *FuzzyDFA, str string) (bool, int) {
	s := dfa.Initial(str)
	mink := dfa.MaxError() + 1
	var final bool
	for dfa.Delta(s, func(k, pos int, data int32) {
		// log.Printf("str=%q", str[:pos])
		// log.Printf(" - k=%d, pos=%d, data=%d", k, pos, data)
		if pos != len(str) {
			return
		}
		if k < mink {
			mink = k
		}
		final = true
		// log.Printf(" - final=%t, mink=%d", final, mink)
	}) {
		// log.Printf("stack: %v", *s)
	}
	// log.Printf("s = %v", *s)
	// log.Printf("accept(%q) = %t, %d", str, final, mink)
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
		{"a with k>3", "a", "xxxax", 0, false},
		{"a with k>3", "a", "xxxxa", 0, false},
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

func TestMatchesFuzzyDFA(t *testing.T) {
	dfa := NewFuzzyDFA(3, NewDictionary("match", "match two"))
	// dfa := NewFuzzyDFA(3, NewDictionary("match"))
	// log.Printf("Initial: %d", dfa.dfa.initial)
	for i := 0; i < len(dfa.dfa.table); i++ {
		if !dfa.dfa.table[i].Empty() {
			// log.Printf("%d %v", i, dfa.dfa.table[i])
		}
	}
	// dfa.dfa.Dot(os.Stdout)
	tests := []struct {
		test   string
		k      int
		accept bool
	}{
		{"match", 0, true},
		{"mxtch", 1, true},
		{"mxxch", 2, true},
		{"mxxxh", 3, true},
		{"ma tch", 1, true},
		{"ma  tch", 2, true},
		{"ma   tch", 3, true},
		{"ma   xch", 0, false},
		{"match two", 0, true},
		{"mxtch two", 1, true},
		{"mxtchtwo", 2, true},
		{"mxtch   two", 3, true},
		{"mxtch to", 2, true},
		{"mxtch tw", 2, true},
		{"mxtc to", 3, true},
		{"mxtc  two", 2, true},
		{"mxtc   two", 3, true},
		{"mxtc    two", 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			final, k := fuzzyAccepts(dfa, tc.test)
			if final != tc.accept {
				t.Fatalf("expected accept(%q)=%t; got %t",
					tc.test, tc.accept, final)
			}
			if final && tc.k != k {
				t.Fatalf("expected accept(%q)=%d; got %d",
					tc.test, tc.k, k)
			}
		})
	}
}
