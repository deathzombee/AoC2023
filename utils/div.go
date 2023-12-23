package utils

import "fmt"

// ExtendedGCD returns the greatest common divisor of a number of integers and the coefficients of Bézout's identity.
func ExtendedGCD(a, b int) (int, int, int) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x, y := ExtendedGCD(b%a, a)
	return gcd, y - (b/a)*x, x
}

// CRT solveCRT solves the system of congruences using the Chinese Remainder Theorem.
func CRT(a, b, mod1, mod2 int) (int, error) {
	invMod1, err := ModInverse(mod1, mod2)
	if err != nil {
		return 0, err
	}
	invMod2, err := ModInverse(mod2, mod1)
	if err != nil {
		return 0, err
	}
	N := mod1 * mod2
	x := (a*mod2*invMod2 + b*mod1*invMod1) % N
	return x, nil
}

// ModInverse returns the modular inverse of a number.
func ModInverse(a, m int) (int, error) {
	gcd, x, _ := ExtendedGCD(a, m)
	if gcd != 1 {
		return 0, fmt.Errorf("modular inverse does not exist")
	}
	return (x + m) % m, nil
}

// LCM returns the least common multiple of a number of integers.
func LCM(nums ...int) int {
	lcm := nums[0]
	for _, num := range nums[1:] {
		lcm = lcm * num / GCD(lcm, num)
	}
	return lcm
}

// GCD returns the greatest common divisor of a number of integers.
func GCD(nums ...int) int {
	gcd := nums[0]
	for _, num := range nums[1:] {
		gcd, _, _ = ExtendedGCD(gcd, num)
	}
	return gcd
}
