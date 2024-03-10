package cplib

func LCM(xs ...int) int {
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
