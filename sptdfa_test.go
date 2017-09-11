package fsa

import (
	"bytes"
	"fmt"
	"math/rand"
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

var chars = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'ä', 'ö', 'ü', 'ß', 'Ä', 'Ö', 'Ü',
	',', '.', '~', '[', ']', '{', '}', '(', ')', ':', '!', '?', ' ',
}

func makeRandomString(r *rand.Rand) string {
	n := r.Intn(100)
	var str string
	for i := 0; i < n; i++ {
		x := r.Intn(len(chars))
		str = fmt.Sprintf("%s%c", str, chars[x])
	}
	return str
}

func makeRandomStrings(n int, r *rand.Rand) (map[string]bool, []string) {
	m := make(map[string]bool)
	var s []string
	for i := 0; i < n; i++ {
		str := makeRandomString(r)
		if !m[str] {
			m[str] = true
			s = append(s, str)
		}
	}
	sort.Slice(s, func(i, j int) bool {
		return bytes.Compare([]byte(s[i]), []byte(s[j])) < 0
	})
	return m, s
}

func makeRandomSparseTableDFA(n int, seed int64, r *rand.Rand) (*SparseTableDFA, map[string]bool, error) {
	m, s := makeRandomStrings(n, r)
	b := NewSparseTableDFABuilder()
	for _, str := range s {
		if err := b.Add(str, 1); err != nil {
			return nil, nil, err
		}
	}
	dfa := b.Build()
	return dfa, m, nil
}

func makeR() (int64, *rand.Rand) {
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))
	return seed, r
}

func TestFuzzy(t *testing.T) {
	seed, r := makeR()
	dfa, m, err := makeRandomSparseTableDFA(100, seed, r)
	if err != nil {
		t.Fatalf("could not add string: %v (%d)", err, seed)
	}
	for str := range m {
		if !accepts(dfa, str) {
			t.Errorf("dfa does not accept %q (%d)", str, seed)
		}
	}
	for i := 0; i < 10000; i++ {
		str := makeRandomString(r)
		if accepts(dfa, str) && !m[str] {
			t.Errorf("dfa accepts %q; but it shouldn't (%d)", str, seed)
		}
		if !accepts(dfa, str) && m[str] {
			t.Errorf("dfa does not accept %q; but it should (%d)", str, seed)
		}
	}
}
