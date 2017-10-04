package sparsetable

type fuzzyState struct {
	lev, next, state int
}

// FuzzyStack keeps track of the active states during the apporimxate search.
type FuzzyStack struct {
	stack []fuzzyState
	dfa   *DFA
	str   string
	max   int
}

func (f *FuzzyStack) empty() bool {
	return len(f.stack) == 0
}

func (f *FuzzyStack) pop() fuzzyState {
	n := len(f.stack)
	if n == 0 {
		panic("called pop() on empty stack")
	}
	top := f.stack[n-1]
	f.stack = f.stack[0 : n-1]
	return top
}

func (f *FuzzyStack) push(s fuzzyState) {
	if s.lev < f.max {
		// log.Printf("State(%v): %v", s.state, f.dfa.CellAt(s.state))
		f.dfa.EachTransition(s.state, func(cell Cell) {
			// log.Printf("Cell: %v", cell)
			// log.Printf("Epsilon from %d %c -> %d", s.state, cell.Char(), cell.Target())
			// log.Printf("Next state(%v): %v", cell.Target(), f.dfa.CellAt(int(cell.Target())))
			f.push(fuzzyState{
				lev:   s.lev + 1,
				state: int(cell.Target()),
				next:  s.next,
			})
		})
	}
	if s.lev <= f.max && s.next <= len(f.str) && s.state >= 0 {
		// log.Printf("pusing %v", s)
		f.stack = append(f.stack, s)
	}
}

func (f *FuzzyStack) deltaDiagonal(s fuzzyState) {
	f.dfa.EachTransition(s.state, func(cell Cell) {
		// log.Printf("deltaDiagonal(%v): at %d with %c -> %d",
		// 	s, s.next, cell.Char(), cell.Target())
		f.push(fuzzyState{
			lev:   s.lev + 1,
			state: int(cell.Target()),
			next:  s.next + 1,
		})
	})
}

func (f *FuzzyStack) deltaVertical(s fuzzyState) {
	// log.Printf("deltaVertical(%v): at %d", s, s.next)
	f.push(fuzzyState{
		lev:   s.lev + 1,
		state: s.state,
		next:  s.next + 1,
	})
}

func (f *FuzzyStack) deltaHorizontal(s fuzzyState) {
	if s.next >= len(f.str) {
		return
	}
	x := f.dfa.Delta(s.state, f.str[s.next])
	if x < 0 {
		return
	}
	// log.Printf("deltaHorizontal(%v): (at pos = %d) from %d with %c -> %d",
	// 	s, s.next, s.state, f.str[s.next], x)
	f.push(fuzzyState{
		lev:   s.lev,
		state: x,
		next:  s.next + 1,
	})
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
func (d *FuzzyDFA) Initial(str string) *FuzzyStack {
	s := &FuzzyStack{
		str: str,
		dfa: d.dfa,
		max: d.k,
	}
	// log.Printf("Initial: %d", d.dfa.Initial())
	s.push(fuzzyState{
		lev:   0,
		state: d.dfa.Initial(),
		next:  0,
	})
	return s
}

// FinalStateCallback is a callback function that is called on final states.
// It is called using the active error, the next position and the data.
type FinalStateCallback func(int, int, int32)

// Delta make one transtion on the top of the stack. If a final state is encountered,
// the callback function is called. It returns false if no more transitions
// can be done with the active stack.
func (d *FuzzyDFA) Delta(f *FuzzyStack, cb FinalStateCallback) bool {
	if f.empty() {
		return false
	}
	top := f.pop()
	// log.Printf("top: %v", top)
	// log.Printf("stack: %v", f.stack)
	f.deltaDiagonal(top)
	f.deltaHorizontal(top)
	f.deltaVertical(top)
	if data, final := d.dfa.Final(top.state); final {
		cb(top.lev, top.next, data)
	}
	// log.Printf("stack: %v", f.stack)
	// log.Printf("####")
	return true
}
