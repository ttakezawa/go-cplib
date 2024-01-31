package main

/*
// SEE: https://betrue12.hateblo.jp/entry/2020/09/23/005940
type S = int
type F = int
ID := -1 << 60

xs := make([]S, N)
seg := NewLazySegtree(
	xs,
	// op(x, y) -> S
	func(x, y S) S {},
	// 単位元E: op(x, E) = op(E, x) = x
	func() S {},
	// mapping: f(x) -> S
	func(f F, x S) S {
		if f == ID {
			return x
		}
		return // TODO
	},
	// composition: f∘g(x) = f(g(x)) -> F(x)
	func(f, g F) F {
		if f == ID {
			return g
		}
		if g == ID {
			return f
		}
		return // TODO
	},
	// 恒等写像identityを返す関数
	func() F { return ID },
)
*/

type LazySegtree[S, F any] struct {
	n, size, log int
	d            []S
	lz           []F
	e            func() S
	merger       func(x, y S) S
	mapper       func(f F, x S) S
	comp         func(f, g F) F
	id           func() F
}

func NewLazySegtree[S, F any](v []S, op func(x, y S) S, e func() S, mapping func(f F, x S) S, composition func(f, g F) F, id func() F) *LazySegtree[S, F] {
	lseg := new(LazySegtree[S, F])
	lseg.n = len(v)
	lseg.log = lseg._ceilPow2(lseg.n)
	lseg.size = 1 << uint(lseg.log)
	lseg.d = make([]S, 2*lseg.size)
	lseg.e = e
	lseg.lz = make([]F, lseg.size)
	lseg.merger = op
	lseg.mapper = mapping
	lseg.comp = composition
	lseg.id = id
	for i := range lseg.d {
		lseg.d[i] = lseg.e()
	}
	for i := range lseg.lz {
		lseg.lz[i] = lseg.id()
	}
	for i := 0; i < lseg.n; i++ {
		lseg.d[lseg.size+i] = v[i]
	}
	for i := lseg.size - 1; i >= 1; i-- {
		lseg._update(i)
	}
	return lseg
}

func (lseg *LazySegtree[S, F]) _update(k int) {
	lseg.d[k] = lseg.merger(lseg.d[2*k], lseg.d[2*k+1])
}

func (lseg *LazySegtree[S, F]) _allApply(k int, f F) {
	lseg.d[k] = lseg.mapper(f, lseg.d[k])
	if k < lseg.size {
		lseg.lz[k] = lseg.comp(f, lseg.lz[k])
	}
}

func (lseg *LazySegtree[S, F]) _push(k int) {
	lseg._allApply(2*k, lseg.lz[k])
	lseg._allApply(2*k+1, lseg.lz[k])
	lseg.lz[k] = lseg.id()
}

func (lseg *LazySegtree[S, F]) Set(p int, x S) {
	p += lseg.size
	for i := lseg.log; i >= 1; i-- {
		lseg._push(p >> uint(i))
	}
	lseg.d[p] = x
	for i := 1; i <= lseg.log; i++ {
		lseg._update(p >> uint(i))
	}
}

func (lseg *LazySegtree[S, F]) Get(p int) S {
	p += lseg.size
	for i := lseg.log; i >= 1; i-- {
		lseg._push(p >> uint(i))
	}
	return lseg.d[p]
}
func (lseg *LazySegtree[S, F]) RangeProd(l, r int) S { return lseg.Prod(l, r) }
func (lseg *LazySegtree[S, F]) Prod(l, r int) S {
	if l < 0 {
		l = 0
	}
	if r > lseg.n {
		r = lseg.n
	}
	if l == r {
		return lseg.e()
	}
	l += lseg.size
	r += lseg.size
	for i := lseg.log; i >= 1; i-- {
		if (l>>uint(i))<<uint(i) != l {
			lseg._push(l >> uint(i))
		}
		if (r>>uint(i))<<uint(i) != r {
			lseg._push(r >> uint(i))
		}
	}
	sml, smr := lseg.e(), lseg.e()
	for l < r {
		if (l & 1) == 1 {
			sml = lseg.merger(sml, lseg.d[l])
			l++
		}
		if (r & 1) == 1 {
			r--
			smr = lseg.merger(lseg.d[r], smr)
		}
		l >>= 1
		r >>= 1
	}
	return lseg.merger(sml, smr)
}

func (lseg *LazySegtree[S, F]) AllProd() S {
	return lseg.d[1]
}

func (lseg *LazySegtree[S, F]) Apply(p int, f F) {
	p += lseg.size
	for i := lseg.log; i >= 1; i-- {
		lseg._push(p >> uint(i))
	}
	lseg.d[p] = lseg.mapper(f, lseg.d[p])
	for i := 1; i <= lseg.log; i++ {
		lseg._update(p >> uint(i))
	}
}

func (lseg *LazySegtree[S, F]) RangeApply(l int, r int, f F) {
	if l < 0 {
		l = 0
	}
	if r > lseg.n {
		r = lseg.n
	}
	if l == r {
		return
	}
	l += lseg.size
	r += lseg.size
	for i := lseg.log; i >= 1; i-- {
		if (l>>uint(i))<<uint(i) != l {
			lseg._push(l >> uint(i))
		}
		if (r>>uint(i))<<uint(i) != r {
			lseg._push((r - 1) >> uint(i))
		}
	}
	l2, r2 := l, r
	for l < r {
		if l&1 == 1 {
			lseg._allApply(l, f)
			l++
		}
		if r&1 == 1 {
			r--
			lseg._allApply(r, f)
		}
		l >>= 1
		r >>= 1
	}
	l, r = l2, r2
	for i := 1; i <= lseg.log; i++ {
		if (l>>uint(i))<<uint(i) != l {
			lseg._update(l >> uint(i))
		}
		if (r>>uint(i))<<uint(i) != r {
			lseg._update((r - 1) >> uint(i))
		}
	}
}

func (lseg *LazySegtree[S, F]) MaxRight(l int, cmp func(v S) bool) int {
	if l == lseg.n {
		return lseg.n
	}
	l += lseg.size
	for i := lseg.log; i >= 1; i-- {
		lseg._push(l >> uint(i))
	}
	sm := lseg.e()
	for {
		for l%2 == 0 {
			l >>= 1
		}
		if !cmp(lseg.merger(sm, lseg.d[l])) {
			for l < lseg.size {
				lseg._push(l)
				l = 2 * l
				if cmp(lseg.merger(sm, lseg.d[l])) {
					sm = lseg.merger(sm, lseg.d[l])
					l++
				}
			}
			return l - lseg.size
		}
		sm = lseg.merger(sm, lseg.d[l])
		l++
		if l&-l == l {
			break
		}
	}
	return lseg.n
}

func (lseg *LazySegtree[S, F]) MinLeft(r int, cmp func(v S) bool) int {
	if r == 0 {
		return 0
	}
	r += lseg.size
	for i := lseg.log; i >= 1; i-- {
		lseg._push(r - 1>>uint(i))
	}
	sm := lseg.e()
	for {
		r--
		for r > 1 && r%2 != 0 {
			r >>= 1
		}
		if !cmp(lseg.merger(lseg.d[r], sm)) {
			for r < lseg.size {
				lseg._push(r)
				r = 2*r + 1
				if cmp(lseg.merger(lseg.d[r], sm)) {
					sm = lseg.merger(lseg.d[r], sm)
					r--
				}
			}
			return r + 1 - lseg.size
		}
		sm = lseg.merger(lseg.d[r], sm)
		if r&-r == r {
			break
		}
	}
	return 0
}

func (lseg *LazySegtree[S, F]) _ceilPow2(n int) (x int) {
	for (1 << uint(x)) < n {
		x++
	}
	return
}
func (lseg *LazySegtree[S, F]) Add(p int, s S) { lseg.Set(p, lseg.merger(lseg.Get(p), s)) }
