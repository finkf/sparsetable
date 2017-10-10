package sparsetable

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"unicode"
	"unicode/utf8"
)

// State represents a the state of a DFA.
// It is a simple integer that points to the active state of the DFA's
// cell table.
type State int

// Valid returns true if the state is still valid.
func (s State) Valid() bool {
	return s >= 0
}

// DFA is a DFA implementation using a sparse table.
type DFA struct {
	table   []Cell
	initial State
}

// NewDictionary builds a minimized sparse table DFA from a list of strings.
// NewDictionary panics if the build process fails.
func NewDictionary(strs ...string) *DFA {
	b := NewBuilder()
	sort.Slice(strs, func(i, j int) bool {
		return bytes.Compare([]byte(strs[i]), []byte(strs[j])) < 0
	})
	for _, str := range strs {
		if err := b.Add(str, 1); err != nil {
			panic(err)
		}
	}
	return b.Build()
}

// Initial returns the initial state of the DFA.
// The state of the DFA is a simple integer that give the position
// of the active cell in the DFA's cell table.
// Values less than 0 mark invalid states.
func (d *DFA) Initial() State {
	return d.initial
}

// Delta makes on transition from the given state s with the given byte c.
func (d DFA) Delta(s State, c byte) State {
	if !d.valid(s, validAnyState) {
		return -1
	}
	o := State(c)
	if int(s+o) >= len(d.table) ||
		!d.table[s+o].Transition() ||
		d.table[s+o].Char() != c {
		return -1
	}
	return State(d.table[s+o].Target())
}

// Final returns the (data, true) if the given state is final.
// If the given state is not final, (0, false) is returned.
func (d *DFA) Final(s State) (int32, bool) {
	if !d.valid(s, validFinalState) {
		return 0, false
	}
	return d.table[s].Final()
}

// EachTransition iterates over all transitions of the given state calling
// the callback function f for each transition cell.
func (d *DFA) EachTransition(s State, f func(Cell)) {
	if !d.valid(s, validAnyState) {
		return
	}
	d.forEachTransition(s, f)
}

var (
	ulen = [...]int{
		1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0,
		2, 2,
		3,
		4,
	}
)

// EachUTF8Transition iterates over all transition of the given state
// calling the callback function f for each transition.
// EachUTF8Transition follows UTF8 mutlibyte sequences to ensure
// that the callback is called for each valid unicode transition.
func (d *DFA) EachUTF8Transition(s State, f func(rune, State)) {
	if !d.valid(s, validAnyState) {
		return
	}
	d.forEachTransition(s, func(cell Cell) {
		buf := [utf8.UTFMax]byte{cell.Char()}
		switch ulen[cell.Char()>>4] {
		case 0:
			f(0, State(cell.Target()))
		case 1:
			f(rune(cell.Char()), State(cell.Target()))
		case 2: // two bytes
			d.forEachUTF8Transition(buf[:], 1, 1, State(cell.Target()), f)
		case 3: // three bytes
			d.forEachUTF8Transition(buf[:], 1, 2, State(cell.Target()), f)
		case 4: // four bytes
			d.forEachUTF8Transition(buf[:], 1, 3, State(cell.Target()), f)
		default: // something else
			panic(fmt.Sprintf("invalid utf8 byte %b encountered", cell.Char()))
		}
	})
}

func (d DFA) forEachUTF8Transition(buf []byte, i, end int, s State, f func(rune, State)) {
	if !d.valid(s, validAnyState) {
		return
	}
	d.forEachTransition(s, func(cell Cell) {
		if !utf8.RuneStart(cell.Char()) {
			buf[i] = cell.Char()
			if i == end {
				r, _ := utf8.DecodeRune(buf)
				f(r, State(cell.Target()))
			} else {
				d.forEachUTF8Transition(buf, i+1, end, State(cell.Target()), f)
			}
		}
	})
}

func (d DFA) forEachTransition(s State, f func(Cell)) {
	if !d.valid(s, validAnyState) {
		return
	}
	for i := State(d.table[s].Next()); i > 0; i = State(d.table[s+i].Next()) {
		cell := d.table[s+i]
		if !cell.Transition() {
			panic(fmt.Sprintf("invalid cell type in EachTransition: %s", cell))
		}
		f(cell)
	}
}

const (
	validTransition = iota
	validAnyState
	validFinalState
	validAny
)

func (d DFA) valid(s State, typ int) bool {
	if s < 0 || int(s) >= len(d.table) {
		return false
	}
	switch typ {
	case validAny:
		return true
	case validAnyState:
		return d.table[s].State()
	case validFinalState:
		_, final := d.table[s].Final()
		return final
	case validTransition:
		return d.table[s].Transition()
	}
	return false
}

// CellAt returns the the cell of the given state.
func (d *DFA) CellAt(s State) Cell {
	if !d.valid(s, validAny) {
		return Cell{}
	}
	return d.table[s]
}

// Dot prints out the dotcode for the DFA.
func (d *DFA) Dot(out io.Writer) {
	dot := "// dotcode\n"
	fmt.Fprintf(out, "digraph dfa { %s", dot)
	fmt.Fprintf(out, " rankdir=LR; %s", dot)
	fmt.Fprintf(out, " -1 [style=invisible,label=\"\",width=0,height=0] %s",
		dot)
	fmt.Fprintf(out, " -1 -> %d %s", d.initial-1, dot)
	for i, cell := range d.table {
		switch cell.typ {
		case finalCellType:
			fmt.Fprintf(out, " %d[peripheries=2] %s", i, dot)
		case nonFinalCellType:
			fmt.Fprintf(out, " %d[peripheries=1] %s", i, dot)
		case transitionCellType:
			fmt.Fprintf(out, " %d -> %d [label=%q] %s",
				i-int(cell.char), cell.data, byte2string(cell.char), dot)
		default:
		}
	}
	fmt.Fprintf(out, "} %s", dot)
}

func byte2string(c byte) string {
	if c < 0x80 && unicode.IsPrint(rune(c)) {
		return fmt.Sprintf("%c", c)
	}
	return fmt.Sprintf("0x%x", c)
}
