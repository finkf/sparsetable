package sparsetable

import "fmt"

const (
	emptyCellType = iota
	finalCellType
	nonFinalCellType
	transitionCellType
)

// Cell represents either a final or non-final state cell
// or transition cell or an empty cell. Cells are used to repesent
// either transitions or states in a DFA.
type Cell struct {
	data            uint32
	char, typ, next byte
}

// FinalCell creates a final cell.
func FinalCell(data uint32, next byte) Cell {
	return Cell{data: data, next: next, typ: finalCellType}
}

// NonFinalCell creates a non-final cell.
func NonFinalCell(next byte) Cell {
	return Cell{next: next, typ: nonFinalCellType}
}

// TransitionCell creates a transtion cell
func TransitionCell(target uint32, char byte, next byte) Cell {
	return Cell{data: target, char: char, next: next, typ: transitionCellType}
}

// State retruns true iff the cell is either a final or a non final state cell.
func (c Cell) State() bool {
	return c.typ == finalCellType || c.typ == nonFinalCellType
}

// Final returns the asociated data of the cell and true if the
// cell represent a final state, Otherwise it returns 0, false.
func (c Cell) Final() (uint32, bool) {
	return c.data, c.typ == finalCellType
}

// Transition return true iff the cell represents a transtion.
func (c Cell) Transition() bool {
	return c.typ == transitionCellType
}

// Empty returns true iff the cell is empty.
func (c Cell) Empty() bool {
	return c.typ == emptyCellType
}

// Target returns the target position of transition cell. It panics
// if it is called on a cell that does not represent a transtion.
func (c Cell) Target() uint32 {
	if !c.Transition() {
		panic("called Target() on a cell that is not a transition")
	}
	return c.data
}

// Char returns the character (byte) that the transition cell represents.
// It panics if it is called on a cell that does not represent a transtion.
func (c Cell) Char() byte {
	if !c.Transition() {
		panic("called Char() on a cell that is not a transition")
	}
	return c.char
}

// Next returns the next transition cell in a states outgoing transitions.
func (c Cell) Next() uint32 {
	return uint32(c.next)
}

// String returns a string representation of the cell.
func (c Cell) String() string {
	switch c.typ {
	case emptyCellType:
		return "{}"
	case finalCellType:
		return fmt.Sprintf("{data: %d, next: %d}", c.data, c.next)
	case nonFinalCellType:
		return fmt.Sprintf("{next: %d}", c.next)
	case transitionCellType:
		return fmt.Sprintf("{char: %c, target: %d, next: %d}",
			c.char, c.data, c.next)
	default:
		panic("invalid cell type")
	}
}
