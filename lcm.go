package cplib

func LCM(xs ...int) int {
	// if len(xs) == 2 {
	// 	return xs[0] / _GCD2(xs[0], xs[1]) * xs[1]
	// }
	l := xs[0]
	for _, x := range xs[1:] {
		l = _LCM2(l, x)
	}
	return l
}

func _LCM2(a, b int) int { return a / _GCD2(a, b) * b }
func _GCD2(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
