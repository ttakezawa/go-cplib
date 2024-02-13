package cplib

import (
	"math/bits"
	"sort"
)

type BitVector struct {
	blockNum int
	bit      []uint32
	cnt      []int
}

func NewBitVector(size int) *BitVector {
	blockNum := (size + 31) >> 5
	bit := make([]uint32, blockNum)
	cnt := make([]int, blockNum)
	return &BitVector{blockNum, bit, cnt}
}

func (bv *BitVector) Set(i int) {
	bv.bit[i>>5] |= 1 << (i & 31)
}

func (bv *BitVector) Build() {
	for i := 0; i < bv.blockNum-1; i++ {
		bv.cnt[i+1] = bv.cnt[i] + bits.OnesCount32(bv.bit[i])
	}
}

func (bv *BitVector) Access(i int) int {
	return int(bv.bit[i>>5] >> (i & 31) & 1)
}

func (bv *BitVector) Rank0(r int) int {
	return r - bv.cnt[r>>5] - bits.OnesCount32(bv.bit[r>>5]&((1<<(r&31))-1))
}

func (bv *BitVector) Rank1(r int) int {
	return bv.cnt[r>>5] + bits.OnesCount32(bv.bit[r>>5]&((1<<(r&31))-1))
}

type WaveletMatrix struct {
	maxlog int
	n      int
	mat    []*BitVector
	zs     []int
}

func NewWaveletMatrix(vals []int, maxlog int) *WaveletMatrix {
	n := len(vals)
	mat := make([]*BitVector, 0, maxlog)
	zs := make([]int, maxlog)
	for d := maxlog - 1; d >= 0; d-- {
		vec := NewBitVector(n + 1)
		ls := make([]int, 0, n)
		rs := make([]int, 0, n)
		for i, val := range vals {
			if val>>d&1 == 1 {
				rs = append(rs, val)
				vec.Set(i)
			} else {
				ls = append(ls, val)
			}
		}
		vec.Build()
		mat = append(mat, vec)
		zs[maxlog-d-1] = len(ls)
		vals = append(ls, rs...)
	}
	return &WaveletMatrix{maxlog, n, mat, zs}
}

func (wm *WaveletMatrix) Access(i int) int {
	res := 0
	for d := 0; d < wm.maxlog; d++ {
		res <<= 1
		if wm.mat[d].Access(i) > 0 {
			res |= 1
			i = wm.mat[d].Rank1(i) + wm.zs[d]
		} else {
			i = wm.mat[d].Rank0(i)
		}
	}
	return res
}

func (wm *WaveletMatrix) Rank(l, r, val int) int {
	for d := 0; d < wm.maxlog; d++ {
		if val>>(wm.maxlog-d-1)&1 == 1 {
			l = wm.mat[d].Rank1(l) + wm.zs[d]
			r = wm.mat[d].Rank1(r) + wm.zs[d]
		} else {
			l = wm.mat[d].Rank0(l)
			r = wm.mat[d].Rank0(r)
		}
	}
	return r - l
}

func (wm *WaveletMatrix) Quantile(l, r, k int) int {
	res := 0
	for d := 0; d < wm.maxlog; d++ {
		res <<= 1
		cntl, cntr := wm.mat[d].Rank0(l), wm.mat[d].Rank0(r)
		if k >= cntr-cntl {
			l = wm.mat[d].Rank1(l) + wm.zs[d]
			r = wm.mat[d].Rank1(r) + wm.zs[d]
			res |= 1
			k -= cntr - cntl
		} else {
			l = cntl
			r = cntr
		}
	}
	return res
}

func (wm *WaveletMatrix) KthSmallest(l, r, k int) int {
	return wm.Quantile(l, r, k)
}

func (wm *WaveletMatrix) KthLargest(l, r, k int) int {
	return wm.Quantile(l, r, r-l-k-1)
}

func (wm *WaveletMatrix) rangeFreq(l, r, upper int) int {
	res := 0
	for d := 0; d < wm.maxlog; d++ {
		if upper>>(wm.maxlog-d-1)&1 == 1 {
			res += wm.mat[d].Rank0(r) - wm.mat[d].Rank0(l)
			l = wm.mat[d].Rank1(l) + wm.zs[d]
			r = wm.mat[d].Rank1(r) + wm.zs[d]
		} else {
			l = wm.mat[d].Rank0(l)
			r = wm.mat[d].Rank0(r)
		}
	}
	return res
}

func (wm *WaveletMatrix) RangeFreq(l, r, lower, upper int) int {
	return wm.rangeFreq(l, r, upper) - wm.rangeFreq(l, r, lower)
}

func (wm *WaveletMatrix) PrevVal(l, r, upper int) int {
	cnt := wm.rangeFreq(l, r, upper)
	if cnt == 0 {
		return -1
	}
	return wm.KthSmallest(l, r, cnt-1)
}

func (wm *WaveletMatrix) NextVal(l, r, lower int) int {
	cnt := wm.rangeFreq(l, r, lower)
	if cnt == r-l {
		return -1
	}
	return wm.KthSmallest(l, r, cnt)
}

type CompressedWaveletMatrix struct {
	vals []int
	comp map[int]int
	wm   *WaveletMatrix
}

func NewCompressedWaveletMatrix(vals []int) *CompressedWaveletMatrix {
	cwm := &CompressedWaveletMatrix{}
	cwm.vals = sorted(removeDuplicate(vals))
	cwm.comp = make(map[int]int)
	for idx, val := range cwm.vals {
		cwm.comp[val] = idx
	}
	newVals := make([]int, len(vals))
	for i, val := range vals {
		newVals[i] = cwm.comp[val]
	}
	cwm.wm = NewWaveletMatrix(newVals, bits.Len(uint(len(cwm.vals))))
	return cwm
}

func removeDuplicate(vals []int) []int {
	m := make(map[int]bool)
	res := make([]int, 0, len(vals))
	for _, val := range vals {
		if !m[val] {
			m[val] = true
			res = append(res, val)
		}
	}
	return res
}

func sorted(vals []int) []int {
	res := make([]int, len(vals))
	copy(res, vals)
	sort.Ints(res)
	return res
}

func (cwm *CompressedWaveletMatrix) Access(i int) int {
	return cwm.vals[cwm.wm.Access(i)]
}

func (cwm *CompressedWaveletMatrix) Rank(l, r, val int) int {
	if idx, ok := cwm.comp[val]; ok {
		return cwm.wm.Rank(l, r, idx)
	}
	return 0
}

func (cwm *CompressedWaveletMatrix) KthSmallest(l, r, k int) int {
	return cwm.vals[cwm.wm.KthSmallest(l, r, k)]
}

func (cwm *CompressedWaveletMatrix) KthLargest(l, r, k int) int {
	return cwm.vals[cwm.wm.KthLargest(l, r, k)]
}

func (cwm *CompressedWaveletMatrix) RangeFreq(l, r, lower, upper int) int {
	lower = sort.SearchInts(cwm.vals, lower)
	upper = sort.SearchInts(cwm.vals, upper)
	return cwm.wm.RangeFreq(l, r, lower, upper)
}

func (cwm *CompressedWaveletMatrix) PrevVal(l, r, upper int) int {
	upper = sort.SearchInts(cwm.vals, upper)
	res := cwm.wm.PrevVal(l, r, upper)
	if res == -1 {
		return -1
	}
	return cwm.vals[res]
}

func (cwm *CompressedWaveletMatrix) NextVal(l, r, lower int) int {
	lower = sort.SearchInts(cwm.vals, lower)
	res := cwm.wm.NextVal(l, r, lower)
	if res == -1 {
		return -1
	}
	return cwm.vals[res]
}
