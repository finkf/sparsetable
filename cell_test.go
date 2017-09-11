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
		{"non final cell", NonFinalCell(0)},
		{"final cell", FinalCell(1, 0)},
		{"transition cell", TransitionCell(1, 'a', 0)},
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
	}{
		{"empty cell", Cell{}, false, false, false, true},
		{"non final cell", NonFinalCell(0), false, false, true, false},
		{"final cell", FinalCell(1, 0), true, false, true, false},
		{"transition cell", TransitionCell(0, 'a', 0), false, true, false, false},
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
		{"non final cell", NonFinalCell(0), "NonFinalCell{next:0}"},
		{"non final cell", NonFinalCell(1), "NonFinalCell{next:1}"},
		{"final cell", FinalCell(1, 0), "FinalCell{data:1,next:0}"},
		{"final cell", FinalCell(2, 1), "FinalCell{data:2,next:1}"},
		{"transition cell", TransitionCell(1, 'a', 0), "TransitionCell{target:1,char:a,next:0}"},
		{"transition cell", TransitionCell(2, 'b', 1), "TransitionCell{target:2,char:b,next:1}"},
		{"transition cell", TransitionCell(3, 'c', 2), "TransitionCell{target:3,char:c,next:2}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if s := tc.cell.String(); s != tc.test {
				t.Fatalf("expected cell = %q; got %q", tc.test, s)
			}
		})
	}
}
