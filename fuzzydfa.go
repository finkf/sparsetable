package sparsetable

type fuzzyState struct {
	str  string
	k, i int
	s    uint32
}

// FuzzyStack keeps track of the active states during the apporimxate search.
type FuzzyStack []fuzzyState

func (f FuzzyStack) push(max int, dfa *DFA, s fuzzyState) FuzzyStack {
	if s.k < max {
		dfa.EachTransition(s.s, func(cell Cell) {
			f = f.push(max, dfa, fuzzyState{
				k:   s.k + 1,
				s:   cell.data,
				i:   s.i,
				str: s.str,
			})
		})
	}
	if s.k <= max {
		f = append(f, s)
	}
	return f
}

// FuzzyDFA is the basic struct for approximate matching on a DFA.
type FuzzyDFA struct {
	dfa *DFA
	k   int
}

// NewFuzzySparseTableDFA create a new FuzzyDFA with a given
// error limit k and a given DFA
func NewFuzzySparseTableDFA(k int, dfa *DFA) *FuzzyDFA {
	return &FuzzyDFA{k: k, dfa: dfa}
}

// MaxError returns the maximum allowed error for the fuzzy DFA.
func (d *FuzzyDFA) MaxError() int {
	return d.k
}

// Initial returns the initial active states of the approximate match for str.
func (d *FuzzyDFA) Initial(str string) FuzzyStack {
	var s FuzzyStack
	return s.push(d.k, d.dfa, fuzzyState{
		k:   0,
		s:   d.dfa.Initial(),
		i:   0,
		str: str,
	})
}

// FinalStateCallback is a callback function that is called on final states.
// It is called using the active error, the next position and the data.
type FinalStateCallback func(int, int, uint32)

// Delta make one transtion on the top of the stack. If a final state is encountered,
// the callback function is called.
func (d *FuzzyDFA) Delta(f FuzzyStack, str string, cb FinalStateCallback) FuzzyStack {
	n := len(f)
	if n == 0 {
		return nil
	}
	top := f[n-1]
	f = f[:n-1]
	f = d.deltaDiagonal(f, top)
	f = d.deltaHorizontal(f, top)
	f = d.deltaVertical(f, top)
	if data, final := d.dfa.Final(top.s); final {
		cb(top.k, top.i, data)
	}
	return f
}

func (d *FuzzyDFA) deltaDiagonal(f FuzzyStack, s fuzzyState) FuzzyStack {
	if d.k <= s.k || len(s.str) <= s.i {
		return f
	}
	d.dfa.EachTransition(s.s, func(cell Cell) {
		f = f.push(d.k, d.dfa, fuzzyState{
			k:   s.k + 1,
			s:   cell.data,
			i:   s.i + 1,
			str: s.str,
		})
	})
	return f
}

func (d *FuzzyDFA) deltaVertical(f FuzzyStack, s fuzzyState) FuzzyStack {
	if d.k <= s.k || len(s.str) <= s.i {
		return f
	}
	return f.push(d.k, d.dfa, fuzzyState{
		k:   s.k + 1,
		s:   s.s,
		i:   s.i + 1,
		str: s.str,
	})
}

func (d *FuzzyDFA) deltaHorizontal(f FuzzyStack, s fuzzyState) FuzzyStack {
	if len(s.str) <= s.i {
		return f
	}
	x := d.dfa.Delta(s.s, s.str[s.i])
	if x == 0 {
		return f
	}
	return f.push(d.k, d.dfa, fuzzyState{
		k:   s.k,
		s:   x,
		i:   s.i + 1,
		str: s.str,
	})
}
