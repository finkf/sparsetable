package sparsetable

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"unicode"
)

// DFA is a DFA implementation using a sparse table.
type DFA struct {
	table   []Cell
	initial int
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
func (d *DFA) Initial() int {
	return d.initial
}

// Delta makes on transition from the given state s with the given byte c.
func (d *DFA) Delta(s int, c byte) int {
	n := len(d.table)
	o := int(c)
	if s < 0 || s >= n || s+o >= n {
		return -1
	}
	if !d.table[s].State() ||
		!d.table[s+o].Transition() ||
		d.table[s+o].Char() != c {
		return -1
	}
	return int(d.table[s+o].Target())
}

// Final returns the (data, true) if the given state is final.
// If the given state is not final, (0, false) is returned.
func (d *DFA) Final(s int) (int32, bool) {
	n := len(d.table)
	if s < 0 || s >= n || d.table[s].typ != finalCellType {
		return 0, false
	}
	return d.table[s].Final()
}

// EachTransition iterates over all transitions of the given state calling
// the callback function f for each transition cell.
func (d *DFA) EachTransition(s int, f func(Cell)) {
	n := int(len(d.table))
	if s < 0 || s >= n {
		return
	}
	if !d.table[s].State() {
		panic(fmt.Sprintf("invalid cell type in EachTransition: %s", d.table[s]))
	}
	for i := int(d.table[s].Next()); i > 0; i = int(d.table[s+i].Next()) {
		cell := d.table[s+i]
		if !cell.Transition() {
			panic(fmt.Sprintf("invalid cell type in EachTransition: %s", cell))
		}
		f(cell)
	}
}

// CellAt returns the the cell of the given state.
func (d *DFA) CellAt(s int) Cell {
	if s < 0 || s >= len(d.table) {
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
