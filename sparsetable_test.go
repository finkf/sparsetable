package sparsetable

import "testing"

func TestAddEmpty(t *testing.T) {
	var st SparseTable
	for i, tc := range []struct {
		final     bool
		data, pos uint32
	}{
		{true, 42, 0},
		{false, 42, 1},
		{false, 42, 2},
		{true, 42, 3},
	} {
		x := st.Add(TmpState{Data: tc.data, Final: tc.final})
		if x != tc.pos {
			t.Errorf("[%d] expected pos = %d; got pos = %d\n",
				i, tc.pos, x)
		}
		if _, final := st.Cells[tc.pos].Final(); final != tc.final {
			t.Errorf("[%d] expected cell.final = %t; got cell.final = %t\n",
				i, tc.final, final)
		}
		// only final state cells have data
		var f uint32
		if tc.final {
			f = 1
		}
		if st.Cells[tc.pos].data != (f * tc.data) {
			t.Errorf("[%d] expected cell.data = %d; got cell.data = %d\n",
				i, tc.data, st.Cells[tc.pos].data)
		}
	}
}

func TestAdd(t *testing.T) {
	var st SparseTable
	for i, tc := range []struct {
		pos uint32
		ts  []TmpStateTransition
	}{
		{0, []TmpStateTransition{
			TmpStateTransition{'a', 0},
			TmpStateTransition{'c', 1},
		}},
		{1, []TmpStateTransition{
			TmpStateTransition{'a', 0},
			TmpStateTransition{'z', 2},
		}},
		{2, []TmpStateTransition{
			TmpStateTransition{'A', 0},
			TmpStateTransition{'Z', 0},
		}},
		{3, []TmpStateTransition{
			TmpStateTransition{'Z', 0},
		}},
	} {
		x := st.Add(TmpState{Transitions: tc.ts})
		if x != tc.pos {
			t.Errorf("[%d] expected pos = %d; got pos = %d\n",
				i, tc.pos, x)
		}
		for j, tt := range tc.ts {
			cell := st.Cells[tc.pos+uint32(tt.char)]
			if !cell.Transition() {
				t.Errorf("[%d:%d] expected transition cell\n", i, j)
			}
			if cell.Char() != tt.char {
				t.Errorf("[%d:%d] expected char = %c; got char = %c\n",
					i, j, tt.char, cell.char)
			}
			if cell.Target() != tt.target {
				t.Errorf("[%d:%d] expected data = %d; got data = %d\n",
					i, j, tt.char, cell.data)
			}
		}
	}
}
