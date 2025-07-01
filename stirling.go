package bigmath

import (
	"math"
	"math/big"
)

// stirlingApproximation computes n! using Stirling's approximation
// n! ≈ sqrt(2πn) * (n/e)^n
func stirlingApproximation(x *big.Float) *big.Float {
	xFloat, _ := x.Float64()

	if xFloat <= 0 {
		return big.NewFloat(1)
	}

	// sqrt(2πn)
	twoPiN := 2 * math.Pi * xFloat
	sqrtTwoPiN := big.NewFloat(math.Sqrt(twoPiN))

	// (n/e)^n = n^n / e^n
	nToN := Pow(x, x)
	eToN := Pow(big.NewFloat(math.E), x)
	nOverEToN := new(big.Float).Quo(nToN, eToN)

	// Combine: sqrt(2πn) * (n/e)^n
	result := new(big.Float).Mul(sqrtTwoPiN, nOverEToN)

	return result
}
