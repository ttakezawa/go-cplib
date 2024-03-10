// Verify https://atcoder.jp/contests/abc177/tasks/abc177_e
package cplib

import (
	"sort"
)

type Sieve struct {
	_minFactor []int
	_primes    []int
}

// 初期化 O(n log log n)
func NewSieve(max int) *Sieve {
	s := new(Sieve)
	s._minFactor = make([]int, max+1)
	s._primes = make([]int, 0)
	for i := 0; i <= max; i++ {
		s._minFactor[i] = i
	}
	for i := 2; i <= max; i++ {
		if s._minFactor[i] != i {
			continue
		}
		s._primes = append(s._primes, i)
		for j := i * i; j <= max; j += i {
			if s._minFactor[j] == j {
				s._minFactor[j] = i
			}
		}
	}
	return s
}

func (s *Sieve) GetPrimes() []int {
	return s._primes
}

func (s *Sieve) IsPrime(n int) bool {
	return n > 1 && s._minFactor[n] == n
}

// SPFによる素因数分解 O(log n)
func (s *Sieve) Factorize(n int) map[int]int {
	ret := make(map[int]int)
	for n > 1 {
		ret[s._minFactor[n]]++
		n /= s._minFactor[n]
	}
	return ret
}

// SPFによる約数列挙 O(√n)
func (s *Sieve) Divisors(n int) []int {
	ret := []int{1}
	for factor, cnt := range s.Factorize(n) {
		for j := 0; j < len(ret); j++ {
			p := 1
			for k := 0; k < cnt; k++ {
				p *= factor
				ret = append(ret, p*ret[j])
			}
		}
	}
	sort.Ints(ret)
	return ret
}
