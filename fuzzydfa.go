package sparsetable

type fuzzyState struct {
	str       string
	lev, next int
	state     uint32
}

// FuzzyStack keeps track of the active states during the apporimxate search.
type FuzzyStack []fuzzyState

// Empty returns true iff this stack is empty.
func (f FuzzyStack) Empty() bool {
	return len(f) == 0
}

func (f FuzzyStack) push(max int, dfa *DFA, s fuzzyState) FuzzyStack {
	if s.lev < max {
		dfa.EachTransition(s.state, func(cell Cell) {
			f = f.push(max, dfa, fuzzyState{
				lev:   s.lev + 1,
				state: cell.Target(),
				next:  s.next,
				str:   s.str,
			})
		})
	}
	if s.lev <= max {
		f = append(f, s)
	}
	return f
}

// FuzzyDFA is the basic struct for approximate matching on a DFA.
type FuzzyDFA struct {
	dfa *DFA
	k   int
}

// NewFuzzyDFA create a new FuzzyDFA with a given
// error limit k and a given DFA
func NewFuzzyDFA(k int, dfa *DFA) *FuzzyDFA {
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
		lev:   0,
		state: d.dfa.Initial(),
		next:  0,
		str:   str,
	})
}

// FinalStateCallback is a callback function that is called on final states.
// It is called using the active error, the next position and the data.
type FinalStateCallback func(int, int, int32)

// Delta make one transtion on the top of the stack. If a final state is encountered,
// the callback function is called.
func (d *FuzzyDFA) Delta(f FuzzyStack, cb FinalStateCallback) FuzzyStack {
	n := len(f)
	if n == 0 {
		return nil
	}
	top := f[n-1]
	f = f[:n-1]
	f = d.deltaDiagonal(f, top)
	f = d.deltaHorizontal(f, top)
	f = d.deltaVertical(f, top)
	if data, final := d.dfa.Final(top.state); final {
		cb(top.lev, top.next, data)
	}
	return f
}

func (d *FuzzyDFA) deltaDiagonal(f FuzzyStack, s fuzzyState) FuzzyStack {
	if d.k <= s.lev || len(s.str) <= s.next {
		return f
	}
	d.dfa.EachTransition(s.state, func(cell Cell) {
		f = f.push(d.k, d.dfa, fuzzyState{
			lev:   s.lev + 1,
			state: cell.Target(),
			next:  s.next + 1,
			str:   s.str,
		})
	})
	return f
}

func (d *FuzzyDFA) deltaVertical(f FuzzyStack, s fuzzyState) FuzzyStack {
	if d.k <= s.lev || len(s.str) <= s.next {
		return f
	}
	return f.push(d.k, d.dfa, fuzzyState{
		lev:   s.lev + 1,
		state: s.state,
		next:  s.next + 1,
		str:   s.str,
	})
}

func (d *FuzzyDFA) deltaHorizontal(f FuzzyStack, s fuzzyState) FuzzyStack {
	if len(s.str) <= s.next {
		return f
	}
	x := d.dfa.Delta(s.state, s.str[s.next])
	if x == 0 {
		return f
	}
	return f.push(d.k, d.dfa, fuzzyState{
		lev:   s.lev,
		state: x,
		next:  s.next + 1,
		str:   s.str,
	})
}
