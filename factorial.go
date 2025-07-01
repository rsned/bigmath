package bigmath

import (
	"math/big"
)

// Factorial calculates n! using big.Int math to ensure no overflow.
//
// big.Int does not have a concept of Inf or NaN, so the best we can
// do for negatives is return 0.
func Factorial(n int64) *big.Int {
	if n < 0 {
		// Factorial is undefined for negative numbers
		return big.NewInt(0)
	}

	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}

	return result
}

// FactorialFloat is a function that returns the factorial of a given big.Float.
// For integer values, computes n! = n * (n-1) * ... * 2 * 1
// For non-integer values, uses the Gamma function property: n! = Gamma(n+1)
//
// Negative values will return +Inf.
func FactorialFloat(x *big.Float) *big.Float {
	// Handle special cases
	if x.Sign() < 0 {
		// Factorial is undefined for negative numbers
		result := big.NewFloat(0)
		result.SetInf(false) // +Inf to indicate undefined

		return result
	}

	// For integer values we can just calculate directly.
	if x.IsInt() {
		i, _ := x.Int64()

		return new(big.Float).SetInt(Factorial(i))
	}

	// For non-integer values, use the relation n! = Gamma(n+1)
	// TODO(rsned): Need to block circular infinite loop cycle because Gamma calls
	// Factorial in some cases.
	return Gamma(x.Add(x, big.NewFloat(1)))
}
