package sparsetable

import (
	"unicode/utf8"
)

type fuzzyState struct {
	lev, next int
	state     State
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
	f.dfa.EachTransition(s.state, func(cell Cell) {
		f.push(fuzzyState{
			lev:   incrError(s.lev, cell.Char()),
			state: State(cell.Target()),
			next:  s.next,
		})
	})
	if s.lev <= f.max && s.next <= len(f.str) && s.state.Valid() {
		if s.lev == 0 {
			// log.Printf("pushing %v (%s)", s, f.str[s.next:])
		}
		f.stack = append(f.stack, s)
	}
}

func (f *FuzzyStack) deltaDiagonal(s fuzzyState) {
	f.dfa.EachTransition(s.state, func(cell Cell) {
		f.push(fuzzyState{
			lev:   incrError(s.lev, cell.Char()),
			state: State(cell.Target()),
			next:  s.next + 1,
		})
	})
}

func (f *FuzzyStack) deltaVertical(s fuzzyState) {
	if s.next < len(f.str) {
		f.push(fuzzyState{
			lev:   incrError(s.lev, f.str[s.next]),
			state: s.state,
			next:  s.next + 1,
		})
	}
}

func (f *FuzzyStack) deltaHorizontal(s fuzzyState) {
	if s.next >= len(f.str) {
		return
	}
	x := f.dfa.Delta(s.state, f.str[s.next])
	if !x.Valid() {
		return
	}
	f.push(fuzzyState{
		lev:   s.lev,
		state: x,
		next:  s.next + 1,
	})
}

func (f *FuzzyStack) delta(top fuzzyState) {
	f.deltaDiagonal(top)
	f.deltaHorizontal(top)
	f.deltaVertical(top)
}

func incrError(k int, b byte) int {
	if utf8.RuneStart(b) {
		return k + 1
	}
	return k
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
	f.delta(top)
	if data, final := d.dfa.Final(top.state); final {
		cb(top.lev, top.next, data)
	}
	return true
}
