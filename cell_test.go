package sparsetable

import (
	"testing"
	"unsafe"
)

func TestCellSize(t *testing.T) {
	tests := []struct {
		name string
		cell Cell
	}{
		{"empty cell", Cell{}},
		{"non final cell", NewNonFinalCell(0)},
		{"final cell", NewFinalCell(1, 0)},
		{"transition cell", NewTransitionCell(1, 'a', 0)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if s := unsafe.Sizeof(tc.cell); s != 8 {
				t.Fatalf("expected cells to have size = 8; got %d", s)
			}
		})
	}
}

func TestCellTraits(t *testing.T) {
	tests := []struct {
		name                            string
		cell                            Cell
		final, transition, state, empty bool
		char                            byte
		next, target                    uint32
		data                            int32
	}{
		{"empty cell", Cell{}, false, false, false, true, 0, 0, 0, 0},
		{"non final cell", NewNonFinalCell(8), false, false, true, false, 0, 8, 0, 0},
		{"final cell", NewFinalCell(1, 10), true, false, true, false, 0, 10, 0, 1},
		{"transition cell", NewTransitionCell(42, 'a', 13), false, true, false, false, 'a', 13, 42, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if res := tc.cell.Transition(); res != tc.transition {
				t.Fatalf("expected transition = %t; got %t", tc.transition, res)
			}
			if _, res := tc.cell.Final(); res != tc.final {
				t.Fatalf("expected final = %t; got %t", tc.final, res)
			}
			if res := tc.cell.State(); res != tc.state {
				t.Fatalf("expected state = %t; got %t", tc.state, res)
			}
			if res := tc.cell.Empty(); res != tc.empty {
				t.Fatalf("expected empty = %t; got %t", tc.empty, res)
			}
			if res := tc.cell.Char(); res != tc.char {
				t.Fatalf("expected char = %c; got %c", tc.char, res)
			}
			if res := tc.cell.Next(); res != tc.next {
				t.Fatalf("expected next = %d; got %d", tc.next, res)
			}
			if res := tc.cell.Target(); res != tc.target {
				t.Fatalf("expected target = %d; got %d", tc.target, res)
			}
			if res, _ := tc.cell.Final(); res != tc.data {
				t.Fatalf("expected data = %d; got %d", tc.data, res)
			}
		})
	}
}

func TestCellInternal(t *testing.T) {
	tests := []struct {
		name string
		cell Cell
		test string
	}{
		{"empty cell", Cell{}, "EmptyCell{}"},
		{"non final cell", NewNonFinalCell(0), "NonFinalCell{next:0}"},
		{"non final cell", NewNonFinalCell(1), "NonFinalCell{next:1}"},
		{"final cell", NewFinalCell(1, 0), "FinalCell{data:1,next:0}"},
		{"final cell", NewFinalCell(2, 1), "FinalCell{data:2,next:1}"},
		{"transition cell", NewTransitionCell(1, 'a', 0), "TransitionCell{target:1,char:a,next:0}"},
		{"transition cell", NewTransitionCell(2, 'b', 1), "TransitionCell{target:2,char:b,next:1}"},
		{"transition cell", NewTransitionCell(3, 'c', 2), "TransitionCell{target:3,char:c,next:2}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if s := tc.cell.String(); s != tc.test {
				t.Fatalf("expected cell = %q; got %q", tc.test, s)
			}
		})
	}
}
