package cplib

// const Mod = int(1e9) + 7
const Mod = 998244353

// ModInv の引数aとmodは互いに素(GCD(a,mod)=1)であることが必要。modが素数なら常に成立する。
func ModInverse(a, mod int) int {
	x, _, _ := ExtGCD(a, mod)
	return (x%mod + mod) % mod
}

// ExtGCD は ax + by = gcd(a, b) を満たす (x,y,aとbの最大公約数) を返す。
func ExtGCD(a, b int) (x, y, gcd int) {
	if b == 0 {
		return 1, 0, a
	}
	y, x, gcd = ExtGCD(b, a%b)
	return x, y - a/b*x, gcd
}
func ModDiv(a, b int) int { return a * ModInverse(b, Mod) % Mod }
func ModPow(a, b, mod int) (x int) {
	a = ((a % mod) + mod) % mod
	for x = 1; b > 0; b >>= 1 {
		if b&1 == 1 {
			x = x * a % mod
		}
		a = a * a % mod
	}
	return
}
func ModFactorial(n int) int { return fac[n] }
func ModPermutation(n, r int) int {
	if n < 0 || n-r < 0 {
		return 0
	}
	return fac[n] * finv[n-r] % Mod
}

func ModCombination(n, k int) (ret int) {
	if n < k {
		return 0
	}
	for ret = 1; n > 0 || k > 0; n, k = n/Mod, k/Mod {
		ret = ret * fac[n%Mod] % Mod * finv[k%Mod] % Mod * finv[(n-k)%Mod] % Mod
	}
	return
}

const maxArg = 5e6

var fac, finv, inv [maxArg + 1]int

func init() {
	fac[0], fac[1], finv[0], finv[1], inv[1] = 1, 1, 1, 1, 1
	for i := 2; i <= maxArg; i++ {
		fac[i], inv[i] = fac[i-1]*i%Mod, Mod-inv[Mod%i]*(Mod/i)%Mod
		finv[i] = finv[i-1] * inv[i] % Mod
	}
}
