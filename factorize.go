// Originated from https://qiita.com/Kiri8128/items/eca965fe86ea5f4cbb98
// Verify https://algo-method.com/tasks/553
package cplib

import (
	"math/big"
	"math/bits"
	"sort"
)

// 素因数分解 O(n⁽¹/⁴⁾)
func factorize(n int) map[int]int {
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
				if is_prime(n) {
					ret[n], n = 1, 1
				} else {
					rhoFlg = 1
					j := find_factor_rho(n)
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
func is_prime(n int) bool {
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
		t := powmod(a, d, n)
		if t == 1 || t == n1 {
			continue
		}
		for i := 0; i < s-1; i++ {
			t = powmod(t, 2, n)
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
func powmod(b, e, m int) int {
	// b**e % m
	bb, ee, mm := big.NewInt(int64(b)), big.NewInt(int64(e)), big.NewInt(int64(m))
	ret := int(new(big.Int).Exp(bb, ee, mm).Int64())
	return ret
}

func find_factor_rho(n int) int {
	m := 1 << (bits.Len(uint(n)) / 8)
	for c := 1; c < 99; c++ {
		f := func(x int) int {
			// return (x*x + c) % n
			ret := new(big.Int)
			ret.Mul(big.NewInt(int64(x)), big.NewInt(int64(x)))
			ret.Add(ret, big.NewInt(int64(c)))
			ret.Mod(ret, big.NewInt(int64(n)))
			return int(ret.Int64())
		}
		var y, r, q, g int = 2, 1, 1, 1
		var x, ys int = 0, 0
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
					ret := new(big.Int)
					ret.Mul(big.NewInt(int64(q)), big.NewInt(int64(abs(x-y))))
					q = int(ret.Mod(ret, big.NewInt(int64(n))).Int64())
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
			if is_prime(g) {
				return g
			} else if is_prime(n / g) {
				return n / g
			}
			return find_factor_rho(g)
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
