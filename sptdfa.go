package fsa

import (
	"fmt"
	"io"
	"unicode"
)

// SparseTableDFA is a DFA implementation using a sparse table.
type SparseTableDFA struct {
	table   []Cell
	initial uint32
}

// Initial returns the initial state of the DFA.
func (d *SparseTableDFA) Initial() uint32 {
	return d.initial
}

// Delta makes on transition from the given state s with the given byte c.
func (d *SparseTableDFA) Delta(s uint32, c byte) uint32 {
	n := uint32(len(d.table))
	o := uint32(c)
	if s <= 0 || s > n || s+o > n {
		return 0
	}
	s--
	if !d.table[s].State() ||
		d.table[s+o].typ != transitionCellType ||
		d.table[s+o].char != c {
		return 0
	}
	return d.table[s+o].data + 1
}

// Final returns the (data, true) if the given state is final.
// If the given state is not final, (0, false) is returned.
func (d *SparseTableDFA) Final(s uint32) (uint32, bool) {
	n := uint32(len(d.table))
	if s <= 0 || n <= s || d.table[s-1].typ != finalCellType {
		return 0, false
	}
	return d.table[s-1].data, true
}

// EachTransition iterates over all transitions of the given state.
func (d *SparseTableDFA) EachTransition(s uint32, f func(Cell)) {
	n := uint32(len(d.table))
	if s <= 0 || s > n || !d.table[s-1].State() {
		panic("invalid cell type in EachTransition: not a state cell")
	}
	for i := d.table[s-1].next; i > 0; i = d.table[s-1+uint32(i)].next {
		cell := d.table[s+uint32(i)-1]
		if cell.typ != transitionCellType {
			panic("invalid cell type in EachTransition: not a transition cell")
		}
		f(cell)
	}
}

// Dot prints out the dotcode for the DFA.
func (d *SparseTableDFA) Dot(out io.Writer) {
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
		case emptyCellType:
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
