package cplib

/*
// 区間和
type S = int
seg := NewSegtree(
	make([]S, N),
	func(x, y S) S { return x + y },
	func() S { return 0 },  // 単位元E: op(x, E) = op(E, x) = x
)
*/

type Segtree[S any] struct {
	n, size, log int
	d            []S
	e            func() S
	merger       func(x, y S) S
}

func NewSegtree[S any](v []S, op func(x, y S) S, e func() S) *Segtree[S] {
	seg := &Segtree[S]{}
	seg.n = len(v)
	seg.log = seg._ceilPow2(seg.n)
	seg.size = 1 << uint(seg.log)
	seg.d = make([]S, 2*seg.size)
	seg.e = e
	seg.merger = op
	for i := range seg.d {
		seg.d[i] = seg.e()
	}
	for i := 0; i < seg.n; i++ {
		seg.d[seg.size+i] = v[i]
	}
	for i := seg.size - 1; i >= 1; i-- {
		seg._update(i)
	}
	return seg
}

func (seg *Segtree[S]) _update(k int) {
	seg.d[k] = seg.merger(seg.d[2*k], seg.d[2*k+1])
}

func (seg *Segtree[S]) Set(p int, x S) {
	p += seg.size
	seg.d[p] = x
	for i := 1; i <= seg.log; i++ {
		seg._update(p >> uint(i))
	}
}

func (seg *Segtree[S]) Get(p int) S {
	return seg.d[p+seg.size]
}

func (seg *Segtree[S]) Prod(l, r int) S {
	if l < 0 {
		l = 0
	}
	if r > seg.n {
		r = seg.n
	}
	sml, smr := seg.e(), seg.e()
	l += seg.size
	r += seg.size
	for l < r {
		if (l & 1) == 1 {
			sml = seg.merger(sml, seg.d[l])
			l++
		}
		if (r & 1) == 1 {
			r--
			smr = seg.merger(seg.d[r], smr)
		}
		l >>= 1
		r >>= 1
	}
	return seg.merger(sml, smr)
}

func (seg *Segtree[S]) AllProd() S {
	return seg.d[1]
}

func (seg *Segtree[S]) MaxRight(l int, cmp func(x S) bool) int {
	if l == seg.n {
		return seg.n
	}
	l += seg.size
	sm := seg.e()
	for {
		for l%2 == 0 {
			l >>= 1
		}
		if !cmp(seg.merger(sm, seg.d[l])) {
			for l < seg.size {
				l = 2 * l
				if cmp(seg.merger(sm, seg.d[l])) {
					sm = seg.merger(sm, seg.d[l])
					l++
				}
			}
			return l - seg.size
		}
		sm = seg.merger(sm, seg.d[l])
		l++
		if l&-l == l {
			break
		}
	}
	return seg.n
}

func (seg *Segtree[S]) MinLeft(r int, cmp func(x S) bool) int {
	if r == 0 {
		return 0
	}
	r += seg.size
	sm := seg.e()
	for {
		r--
		for r > 1 && r%2 != 0 {
			r >>= 1
		}
		if !cmp(seg.merger(seg.d[r], sm)) {
			for r < seg.size {
				r = 2*r + 1
				if cmp(seg.merger(seg.d[r], sm)) {
					sm = seg.merger(seg.d[r], sm)
					r--
				}
			}
			return r + 1 - seg.size
		}
		sm = seg.merger(seg.d[r], sm)
		if r&-r == r {
			break
		}
	}
	return 0
}

func (seg *Segtree[S]) _ceilPow2(n int) (x int) {
	for (1 << uint(x)) < n {
		x++
	}
	return
}
func (seg *Segtree[S]) Add(p int, s S) { seg.Set(p, seg.merger(seg.Get(p), s)) }
