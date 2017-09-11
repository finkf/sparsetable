package sparsetable

import (
	"bytes"
	"fmt"
)

// SparseTableDFABuilder is used to build a DFA.
type SparseTableDFABuilder struct {
	register  map[string]uint32
	curstr    []byte
	curdat    uint32
	tmpStates []TmpState
	table     SparseTable
}

// NewSparseTableDFABuilder return a new instance of a SparseTableDFABuilder.
func NewSparseTableDFABuilder() *SparseTableDFABuilder {
	return &SparseTableDFABuilder{register: make(map[string]uint32)}
}

// Add adds a (string, value) pair to the sparse table. Add returns an error
// iff the added strings are not in byte-wise lexicographical order.
func (b *SparseTableDFABuilder) Add(str string, data uint32) error {
	nextstr := []byte(str)
	if b.curstr == nil {
		b.curstr = nextstr
		b.curdat = data
		return nil
	}
	if !(bytes.Compare(b.curstr, nextstr) < 0) {
		return fmt.Errorf("add: not in lexicographical order: %q >= %q",
			b.curstr, nextstr)
	}
	b.initTmpStates()
	prefix := commonPrefix(b.curstr, nextstr)
	b.insertSuffix(b.curstr, prefix)
	b.curstr = nextstr
	b.curdat = data
	return nil
}

// Build finishes the building of the automaton and returns it.
func (b *SparseTableDFABuilder) Build() *DFA {
	if b.curstr == nil {
		return &DFA{}
	}

	b.initTmpStates()
	b.insertSuffix(b.curstr, 0)
	initial := b.table.Add(b.tmpStates[0])
	return &DFA{
		table:   b.table.Cells,
		initial: initial + 1,
	}
}

func (b *SparseTableDFABuilder) initTmpStates() {
	n := len(b.curstr)
	for len(b.tmpStates) < n+1 {
		b.tmpStates = append(b.tmpStates, TmpState{})
	}
	b.tmpStates[n].Final = true
	b.tmpStates[n].Data = b.curdat
}

func (b *SparseTableDFABuilder) insertSuffix(str []byte, prefix int) {
	for i := len(str); i > prefix; i-- {
		target := b.replaceOrRegister(b.tmpStates[i])
		b.tmpStates[i] = TmpState{Final: false, Data: 0}
		b.tmpStates[i-1].Transitions = append(
			b.tmpStates[i-1].Transitions,
			TmpStateTransition{char: str[i-1], target: target},
		)
	}
}

func (b *SparseTableDFABuilder) replaceOrRegister(tmp TmpState) uint32 {
	str := tmp.String()
	if target, ok := b.register[str]; ok {
		return target
	}
	target := b.table.Add(tmp)
	b.register[str] = target
	return target
}

// a != nil and b != nil
// a < b
func commonPrefix(a, b []byte) int {
	var n int
	for n = 0; n < len(a); n++ {
		if a[n] != b[n] {
			return n
		}
	}
	return n
	// panic(fmt.Sprintf("commonPrefix(%q, %q): not reached", a, b))
}
