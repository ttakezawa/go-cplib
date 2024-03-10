package cplib

func GCD(xs ...int) int {
	if len(xs) == 2 {
		if xs[1] == 0 {
			return xs[0]
		}
		return GCD(xs[1], xs[0]%xs[1])
	}
	g := xs[0]
	for _, x := range xs[1:] {
		g = GCD(g, x)
	}
	return g
}
