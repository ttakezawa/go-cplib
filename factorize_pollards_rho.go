package cplib

// Originated from https://qiita.com/Kiri8128/items/eca965fe86ea5f4cbb98
// Verify https://algo-method.com/tasks/553

import (
	"math/bits"
	"sort"
)

// ポラード・ロー素因数分解法 O(n⁽¹/⁴⁾)
func FactorizePollardsRho(n int) map[int]int {
	ret := make(map[int]int)
	i := 2
	rhoFlg := 0
	for i*i <= n {
		k := 0
		for n%i == 0 {
			n /= i
			k++
		}
		if k > 0 {
			ret[i] = k
		}
		if i%3 == 1 {
			i += i%2 + 3
		} else {
			i += i%2 + 1
		}
		if i == 101 && n >= 1<<20 {
			for n > 1 {
				if IsPrime(n) {
					ret[n], n = 1, 1
				} else {
					rhoFlg = 1
					j := _findFactorRho(n)
					k := 0
					for n%j == 0 {
						n /= j
						k++
					}
					ret[j] = k
				}
			}
		}
	}
	if n > 1 {
		ret[n] = 1
	}
	if rhoFlg > 0 {
		keys := make([]int, 0, len(ret))
		for k := range ret {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		tmp := make(map[int]int)
		for _, k := range keys {
			tmp[k] = ret[k]
		}
		ret = tmp
	}
	return ret
}

// ミラーラビン素数判定 O(1)
func IsPrime(n int) bool {
	if n == 2 {
		return true
	}
	if n < 2 || n&1 == 0 {
		return false
	}
	n1 := n - 1
	d, s := n1, 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	for _, a := range []int{2, 325, 9375, 28178, 450775, 9780504, 1795265022} {
		if a%n == 0 {
			continue
		}
		t := _powmod(a, d, n)
		if t == 1 || t == n1 {
			continue
		}
		for i := 0; i < s-1; i++ {
			t = _powmod(t, 2, n)
			if t == n1 {
				break
			}
		}
		if t != n1 {
			return false
		}
	}
	return true
}

// b^e mod m O(log e)
func _powmod(b, e, m int) int {
	if e == 0 {
		return 1
	}
	t := _powmod(b, e/2, m)
	t = int(NewUint128(uint64(t)).Mul64(uint64(t)).Mod64(uint64(m)))
	if e%2 == 1 {
		t = int(NewUint128(uint64(t)).Mul64(uint64(b)).Mod64(uint64(m)))
	}
	return t
}

func _findFactorRho(n int) int {
	m := 1 << (bits.Len(uint(n)) / 8)
	for c := 1; c < 99; c++ {
		f := func(x int) int {
			// x*x+c mod n
			return int(NewUint128(uint64(x)).Mul64(uint64(x)).Add64(uint64(c)).Mod64(uint64(n)))
		}
		y, r, q, g := 2, 1, 1, 1
		x, ys := 0, 0
		for g == 1 {
			x = y
			for i := 0; i < r; i++ {
				y = f(y)
			}
			k := 0
			for k < r && g == 1 {
				ys = y
				for i := 0; i < min(m, r-k); i++ {
					y = f(y)
					// q = q * abs(x-y) % n
					ret := NewUint128(uint64(q))
					ret = ret.Mul64(uint64(abs(x - y)))
					q = int(ret.Mod64(uint64(n)))
				}
				g = gcd(q, n)
				k += m
			}
			r <<= 1
		}
		if g == n {
			g = 1
			for g == 1 {
				ys = f(ys)
				g = gcd(abs(x-ys), n)
			}
		}
		if g < n {
			if IsPrime(g) {
				return g
			} else if IsPrime(n / g) {
				return n / g
			}
			return _findFactorRho(g)
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(xs ...int) int {
	if len(xs) == 2 {
		if xs[1] == 0 {
			return xs[0]
		}
		return gcd(xs[1], xs[0]%xs[1])
	}
	g := xs[0]
	for _, x := range xs[1:] {
		g = gcd(g, x)
	}
	return g
}

type Uint128 struct {
	Lo, Hi uint64
}

func NewUint128(v uint64) Uint128 {
	return Uint128{v, 0}
}

func (u Uint128) Add64(v uint64) Uint128 {
	lo, carry := bits.Add64(u.Lo, v, 0)
	hi, carry := bits.Add64(u.Hi, 0, carry)
	if carry != 0 {
		panic("overflow")
	}
	return Uint128{lo, hi}
}

func (u Uint128) Mul64(v uint64) Uint128 {
	hi, lo := bits.Mul64(u.Lo, v)
	p0, p1 := bits.Mul64(u.Hi, v)
	hi, c0 := bits.Add64(hi, p1, 0)
	if p0 != 0 || c0 != 0 {
		panic("overflow")
	}
	return Uint128{lo, hi}
}

func (u Uint128) QuoRem64(v uint64) (q Uint128, r uint64) {
	if u.Hi < v {
		q.Lo, r = bits.Div64(u.Hi, u.Lo, v)
	} else {
		q.Hi, r = bits.Div64(0, u.Hi, v)
		q.Lo, r = bits.Div64(r, u.Lo, v)
	}
	return
}

func (u Uint128) Mod64(v uint64) (r uint64) {
	_, r = u.QuoRem64(v)
	return
}
