package fsa

import "fmt"

const (
	EmptyCellType = iota
	FinalCellType
	NonFinalCellType
	TransitionCellType
)

type Cell struct {
	data            uint32
	char, typ, next byte
}

func FinalCell(data uint32, next byte) Cell {
	return Cell{data: data, next: next, typ: FinalCellType}
}

func NonFinalCell(next byte) Cell {
	return Cell{next: next, typ: NonFinalCellType}
}

func TransitionCell(target uint32, char byte, next byte) Cell {
	return Cell{data: target, char: char, next: next, typ: TransitionCellType}
}

func (c Cell) State() bool {
	return c.typ == FinalCellType || c.typ == NonFinalCellType
}

func (c Cell) Final() bool {
	return c.typ == FinalCellType
}

func (c Cell) String() string {
	switch c.typ {
	case EmptyCellType:
		return "{}"
	case FinalCellType:
		return fmt.Sprintf("{data: %d, next: %d}", c.data, c.next)
	case NonFinalCellType:
		return fmt.Sprintf("{next: %d}", c.next)
	case TransitionCellType:
		return fmt.Sprintf("{char: %c, target: %d, next: %d}",
			c.char, c.data, c.next)
	default:
		panic("invalid cell type")
	}
}
