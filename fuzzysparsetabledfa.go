package fsa

type fuzzyState struct {
	str  string
	k, i int
	s    uint32
}

type fuzzyStack []fuzzyState

func (f fuzzyStack) push(max int, dfa *SparseTableDFA, s fuzzyState) fuzzyStack {
	if s.k < max {
		dfa.EachTransition(s.s, func(cell Cell) {
			f = f.push(max, dfa, fuzzyState{
				k:   s.k + 1,
				s:   cell.data,
				i:   s.i,
				str: s.str,
			})
		})
	}
	if s.k <= max {
		f = append(f, s)
	}
	return f
}

type FuzzySparseTableDFA struct {
	dfa *SparseTableDFA
	k   int
}

func NewFuzzySparseTableDFA(k int, dfa *SparseTableDFA) *FuzzySparseTableDFA {
	return &FuzzySparseTableDFA{k: k, dfa: dfa}
}

func (d *FuzzySparseTableDFA) Initial(str string) fuzzyStack {
	var s fuzzyStack
	return s.push(d.k, d.dfa, fuzzyState{
		k:   0,
		s:   d.dfa.Initial(),
		i:   0,
		str: str,
	})
}

func (d *FuzzySparseTableDFA) Delta(f fuzzyStack, str string, cb func(int, string, uint32)) fuzzyStack {
	n := len(f)
	if n == 0 {
		return nil
	}
	top := f[n-1]
	f = f[:n-1]
	f = d.deltaDiagonal(f, top)
	f = d.deltaHorizontal(f, top)
	f = d.deltaVertical(f, top)
	if data, final := d.dfa.Final(top.s); final {
		cb(top.i, top.str, data)
	}
	return f
}

func (d *FuzzySparseTableDFA) deltaDiagonal(f fuzzyStack, s fuzzyState) fuzzyStack {
	if d.k <= s.k || len(s.str) <= s.i {
		return f
	}
	d.dfa.EachTransition(s.s, func(cell Cell) {
		f = f.push(d.k, d.dfa, fuzzyState{
			k:   s.k + 1,
			s:   cell.data,
			i:   s.i + 1,
			str: s.str,
		})
	})
	return f
}

func (d *FuzzySparseTableDFA) deltaVertical(f fuzzyStack, s fuzzyState) fuzzyStack {
	if d.k <= s.k || len(s.str) <= s.i {
		return f
	}
	return f.push(d.k, d.dfa, fuzzyState{
		k:   s.k + 1,
		s:   s.s,
		i:   s.i + 1,
		str: s.str,
	})
}

func (d *FuzzySparseTableDFA) deltaHorizontal(f fuzzyStack, s fuzzyState) fuzzyStack {
	if len(s.str) <= s.i {
		return f
	}
	x := d.dfa.Delta(s.s, s.str[s.i])
	if x == 0 {
		return f
	}
	return f.push(d.k, d.dfa, fuzzyState{
		k:   s.k,
		s:   x,
		i:   s.i + 1,
		str: s.str,
	})
}
