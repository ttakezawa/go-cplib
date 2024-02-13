package cplib

import "testing"

func TestWaveletMatrix(t *testing.T) {
	v := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9}
	wm := NewWaveletMatrix(v, 32)
	if wm.Access(4) != 5 {
		t.Errorf("wm.Access(4) = %v, want %v", wm.Access(4), 5)
	}
	if wm.Rank(8, 9, 5) != 1 {
		t.Errorf("wm.Rank(8, 9, 5) = %v, want %v", wm.Rank(8, 9, 5), 1)
	}
	if wm.Rank(0, len(v), 9) != 3 {
		t.Errorf("wm.Rank(0, len(v), 9) = %v, want %v", wm.Rank(0, len(v), 9), 3)
	}
	if wm.KthSmallest(1, 4, 2) != 4 {
		t.Errorf("wm.KthSmallest(1, 4, 2) = %v, want %v", wm.KthSmallest(1, 4, 2), 4)
	}
	if wm.KthLargest(1, 5, 3) != 1 {
		t.Errorf("wm.KthLargest(1, 5, 3) = %v, want %v", wm.KthLargest(1, 5, 3), 1)
	}
	if wm.RangeFreq(3, 10, 5, 7) != 3 {
		t.Errorf("wm.RangeFreq(3, 10, 5, 7) = %v, want %v", wm.RangeFreq(3, 10, 5, 7), 3)
	}
	if wm.PrevVal(4, 9, 5) != 2 {
		t.Errorf("wm.PrevVal(4, 9, 5) = %v, want %v", wm.PrevVal(4, 9, 5), 2)
	}
	if wm.NextVal(4, 9, 7) != 9 {
		t.Errorf("wm.NextVal(4, 9, 7) = %v, want %v", wm.NextVal(4, 9, 7), 9)
	}
}

func TestCompressedWaveletMatrix(t *testing.T) {
	v := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9}
	cwm := NewCompressedWaveletMatrix(v)
	if cwm.Access(4) != 5 {
		t.Errorf("cwm.Access(4) = %v, want %v", cwm.Access(4), 5)
	}
	if cwm.Rank(8, 9, 5) != 1 {
		t.Errorf("cwm.Rank(8, 9, 5) = %v, want %v", cwm.Rank(8, 9, 5), 1)
	}
	if cwm.Rank(0, len(v), 9) != 3 {
		t.Errorf("cwm.Rank(0, len(v), 9) = %v, want %v", cwm.Rank(0, len(v), 9), 3)
	}
	if cwm.KthSmallest(1, 4, 2) != 4 {
		t.Errorf("cwm.KthSmallest(1, 4, 2) = %v, want %v", cwm.KthSmallest(1, 4, 2), 4)
	}
	if cwm.KthLargest(1, 5, 3) != 1 {
		t.Errorf("cwm.KthLargest(1, 5, 3) = %v, want %v", cwm.KthLargest(1, 5, 3), 1)
	}
	if cwm.RangeFreq(3, 10, 5, 7) != 3 {
		t.Errorf("cwm.RangeFreq(3, 10, 5, 7) = %v, want %v", cwm.RangeFreq(3, 10, 5, 7), 3)
	}
	if cwm.PrevVal(4, 9, 5) != 2 {
		t.Errorf("cwm.PrevVal(4, 9, 5) = %v, want %v", cwm.PrevVal(4, 9, 5), 2)
	}
	if cwm.NextVal(4, 9, 7) != 9 {
		t.Errorf("cwm.NextVal(4, 9, 7) = %v, want %v", cwm.NextVal(4, 9, 7), 9)
	}
}
