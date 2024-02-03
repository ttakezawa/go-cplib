package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"golang.org/x/exp/constraints"
)

func main() {
	defer _w.Flush()
	N := ReadInt()
	A := ReadInts(N)
	prints(A)
}

var _r, _w = bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)

func ReadString() (s string) { fmt.Fscan(_r, &s); return }
func ReadInt() (i int)       { fmt.Fscan(_r, &i); return }
func ReadInts(n int) (s []int) {
	for i := 0; i < n; i++ {
		s = append(s, ReadInt())
	}
	return
}
func print(a ...any)       { fmt.Fprint(_w, a...) }
func println(a ...any)     { fmt.Fprintln(_w, a...) }
func prints[T any](xs []T) { v := fmt.Sprint(xs); println(v[1 : len(v)-1]) }

// Max, Min
func max[T constraints.Ordered](x T, ys ...T) T {
	for _, y := range ys {
		if x < y || isNaN(x) {
			x = y
		}
	}
	return x
}

func min[T constraints.Ordered](x T, ys ...T) T {
	for _, y := range ys {
		if x > y || isNaN(x) {
			x = y
		}
	}
	return x
}
func Chmax[T constraints.Ordered](p *T, xs ...T)     { *p = max(*p, xs...) }
func Chmin[T constraints.Ordered](p *T, xs ...T)     { *p = min(*p, xs...) }
func isNaN[T constraints.Ordered](x T) bool          { return x != x }
func MaxSlice[S ~[]E, E constraints.Ordered](xs S) E { return max(xs[0], xs[1:]...) }
func MinSlice[S ~[]E, E constraints.Ordered](xs S) E { return min(xs[0], xs[1:]...) }

// SortBy sorts slice xs by f(x) in ascending order.
func SortBy[T any, U constraints.Ordered](xs []T, f func(x T) U) {
	SortBySlice(xs, func(x T) []U { return []U{f(x)} })
}

func SortBySlice[T any, U constraints.Ordered](xs []T, f func(x T) []U) {
	sort.Slice(xs, func(i, j int) bool {
		s1, s2 := f(xs[i]), f(xs[j])
		for i := 0; i < len(s1) && i < len(s2); i++ {
			if s1[i] != s2[i] {
				return s1[i] < s2[i]
			}
		}
		return len(s1) < len(s2)
	})
}

func SortBy2[T any, U constraints.Ordered](xs []T, f func(x T) (U, U)) {
	SortBySlice(xs, func(x T) []U {
		a, b := f(x)
		return []U{a, b}
	})
}
