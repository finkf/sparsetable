package sparsetable

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// CellType represents the type of a cell.
type CellType byte

// There are four types of cells:
// empty (unused) cells, final state cells, non final state cells
// and transition cells.
// Final and non final state cells represent states in the automaton.
// Transition cells represent transtions between the different states
// in the automaton.
const (
	EmptyCell CellType = iota
	FinalCell
	NonFinalCell
	TransitionCell
)

// Cell represents either a final or non-final state cell
// or transition cell or an empty cell. Cells are used to repesent
// either transitions or states in a DFA.
type Cell struct {
	data       int32
	char, next byte
	typ        CellType
}

// NewFinalCell creates a final cell.
func NewFinalCell(data int32, next byte) Cell {
	return Cell{data: data, next: next, typ: FinalCell}
}

// NewNonFinalCell creates a non-final cell.
func NewNonFinalCell(next byte) Cell {
	return Cell{next: next, typ: NonFinalCell}
}

// NewTransitionCell creates a transtion cell
func NewTransitionCell(target uint32, char byte, next byte) Cell {
	return Cell{data: int32(target), char: char, next: next, typ: TransitionCell}
}

// State retruns true iff the cell is either a final or a non final state cell.
func (c Cell) State() bool {
	return c.typ == FinalCell || c.typ == NonFinalCell
}

// Final returns the asociated data of the cell and true if the
// cell represent a final state, Otherwise it returns 0, false.
func (c Cell) Final() (int32, bool) {
	if c.typ != FinalCell {
		return 0, false
	}
	return c.data, true
}

// Transition return true iff the cell represents a transtion.
func (c Cell) Transition() bool {
	return c.typ == TransitionCell
}

// Empty returns true iff the cell is empty.
func (c Cell) Empty() bool {
	return c.typ == EmptyCell
}

// Target returns the target position of transition cell.
func (c Cell) Target() uint32 {
	if !c.Transition() {
		return 0
	}
	return uint32(c.data)
}

// Char returns the character (byte) that the transition cell represents.
func (c Cell) Char() byte {
	if !c.Transition() {
		return 0
	}
	return c.char
}

// Next returns the next transition cell in a states outgoing transitions.
func (c Cell) Next() uint32 {
	return uint32(c.next)
}

// Type returns the type of this cell.
func (c Cell) Type() CellType {
	return c.typ
}

// String returns a string representation of the cell.
func (c Cell) String() string {
	switch c.typ {
	case EmptyCell:
		return "EmptyCell{}"
	case FinalCell:
		return fmt.Sprintf("FinalCell{data:%d,next:%d}", c.data, c.next)
	case NonFinalCell:
		return fmt.Sprintf("NonFinalCell{next:%d}", c.next)
	case TransitionCell:
		return fmt.Sprintf("TransitionCell{target:%d,char:%c,next:%d}",
			c.data, c.char, c.next)
	default:
		panic("invalid cell type")
	}
}

// GobDecode decodes a cells from gob.
func (c *Cell) GobDecode(bs []byte) error {
	buffer := bytes.NewBuffer(bs)
	decoder := gob.NewDecoder(buffer)
	var data int32
	var char, next byte
	var typ CellType
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if err := decoder.Decode(&char); err != nil {
		return err
	}
	if err := decoder.Decode(&next); err != nil {
		return err
	}
	if err := decoder.Decode(&typ); err != nil {
		return err
	}
	c.data = data
	c.char = char
	c.next = next
	c.typ = typ
	return nil
}

// GobEncode encodes a cell to gob.
func (c Cell) GobEncode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(c.data); err != nil {
		return nil, err
	}
	if err := encoder.Encode(c.char); err != nil {
		return nil, err
	}
	if err := encoder.Encode(c.next); err != nil {
		return nil, err
	}
	if err := encoder.Encode(c.typ); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
