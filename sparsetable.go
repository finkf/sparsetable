package fsa

import (
	"fmt"
)

// TmpStateTransition represent an outgoing transition from a temporary
// state.
type TmpStateTransition struct {
	char   byte
	target uint32
}

// String returns a strin representation for a temporary state
// transition. It is used basically for hashing.
func (t TmpStateTransition) String() string {
	return fmt.Sprintf("%c %d", t.char, t.target)
}

// TmpState represents a state that should be inerter into a sparse table.
// It contains a sorted list of
type TmpState struct {
	Transitions []TmpStateTransition
	Data        uint32
	Final       bool
}

// String returns a strin representation for a temporary state.
// It is used basically for hashing.
func (t TmpState) String() string {
	return fmt.Sprintf("%t %d %v", t.Final, t.Data, t.Transitions)
}

// SparseTable is a sparse table of cells.
type SparseTable struct {
	Cells []Cell
	free  uint32
}

// Add adds a temporary state into the sparse table. It returns the
// absolute position where the state was inserted.
// The transitions of the temorary state must b sorted.
func (t *SparseTable) Add(tmp TmpState) uint32 {
	start := t.findFreeTableCell(tmp)
	t.doInsert(start, tmp)
	t.nextFreeCell()
	return start
}

func (t *SparseTable) doInsert(i uint32, tmp TmpState) {
	var next byte
	if len(tmp.Transitions) > 0 {
		next = tmp.Transitions[0].char
	}
	if tmp.Final {
		t.Cells[i] = FinalCell(tmp.Data, next)
	} else {
		t.Cells[i] = NonFinalCell(next)
	}
	for j, trans := range tmp.Transitions {
		next = 0
		if (j + 1) < len(tmp.Transitions) {
			next = tmp.Transitions[j+1].char
		}
		pos := i + uint32(trans.char)
		t.Cells[pos] = TransitionCell(trans.target, trans.char, next)
	}
}

func (t *SparseTable) findFreeTableCell(tmp TmpState) uint32 {
	for i := t.free; ; i++ {
		t.resize(i, tmp)
		if t.fits(i, tmp) {
			return i
		}
	}
}

func (t *SparseTable) resize(i uint32, tmp TmpState) {
	if len(tmp.Transitions) > 0 {
		i += uint32(tmp.Transitions[len(tmp.Transitions)-1].char)
	}
	for uint32(len(t.Cells)) < (i + 1) {
		t.Cells = append(t.Cells, Cell{})
	}
}

func (t *SparseTable) fits(i uint32, tmp TmpState) bool {
	if !t.Cells[i].Empty() {
		return false
	}
	for _, trans := range tmp.Transitions {
		if t.Cells[i+uint32(trans.char)].typ != emptyCellType {
			return false
		}
	}
	return true
}

func (t *SparseTable) nextFreeCell() {
	for {
		if uint32(len(t.Cells)) <= t.free {
			t.Cells = append(t.Cells, Cell{})
		}
		if t.Cells[t.free].Empty() {
			break
		}
		t.free += 1
	}
}
