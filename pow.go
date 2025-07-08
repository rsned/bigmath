// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigmath

import (
	"math"
	"math/big"
)

// isOddInteger checks if a big.Float represents an odd integer
func isOddInteger(f *big.Float) bool {
	if !f.IsInt() {
		return false
	}

	// Convert to big.Int to check if odd
	bigInt := new(big.Int)
	f.Int(bigInt)

	// Check if the least significant bit is 1 (odd)
	return bigInt.Bit(0) == 1
}

// Pow returns x**y, the base-x exponential of y using *big.Float arithmetic.
//
// Special cases are (in order):
//
// Pow(x, ±0) = 1 for any x
// Pow(1, y) = 1 for any y
// Pow(x, 1) = x for any x
// Pow(NaN, y) = NaN
// Pow(x, NaN) = NaN
// Pow(±0, y) = ±Inf for y an odd integer < 0
// Pow(±0, -Inf) = +Inf
// Pow(±0, +Inf) = +0
// Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
// Pow(±0, y) = ±0 for y an odd integer > 0
// Pow(±0, y) = +0 for finite y > 0 and not an odd integer
// Pow(-1, ±Inf) = 1
// Pow(x, +Inf) = +Inf for |x| > 1
// Pow(x, -Inf) = +0 for |x| > 1
// Pow(x, +Inf) = +0 for |x| < 1
// Pow(x, -Inf) = +Inf for |x| < 1
// Pow(+Inf, y) = +Inf for y > 0
// Pow(+Inf, y) = +0 for y < 0
// Pow(-Inf, y) = Pow(-0, -y)
// Pow(x, y) = NaN for finite x < 0 and finite non-integer y
func Pow(x, y *big.Float) *big.Float {
	zero := big.NewFloat(0)
	one := big.NewFloat(1)
	negOne := big.NewFloat(-1)

	// Handle special cases in order as documented

	// Pow(x, ±0) = 1 for any x
	if y.Cmp(zero) == 0 {
		return big.NewFloat(1)
	}

	// Pow(1, y) = 1 for any y
	if x.Cmp(one) == 0 {
		return big.NewFloat(1)
	}

	// Pow(x, 1) = x for any x
	if y.Cmp(one) == 0 {
		return new(big.Float).Copy(x)
	}

	// Check for infinity in x or y
	xIsInf := x.IsInf()
	yIsInf := y.IsInf()

	// Handle ±0 base cases
	if x.Cmp(zero) == 0 {
		ySign := y.Sign()
		if ySign < 0 {
			// Check for odd integer
			if y.IsInt() && isOddInteger(y) {
				// Pow(±0, y) = ±Inf for y an odd integer < 0
				result := big.NewFloat(0)
				// Check if x is negative zero by checking both zero and signbit
				if x.Cmp(zero) == 0 && x.Signbit() {
					result.SetInf(true) // -Inf for -0
				} else {
					result.SetInf(false) // +Inf for +0
				}

				return result
			}
			// Pow(±0, -Inf) = +Inf or Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			result := big.NewFloat(0)
			result.SetInf(false) // +Inf

			return result
		}

		// y > 0
		if yIsInf {
			// Pow(±0, +Inf) = +0
			return big.NewFloat(0)
		}
		if y.IsInt() && isOddInteger(y) {
			// Pow(±0, y) = ±0 for y an odd integer > 0
			if x.Cmp(zero) == 0 && x.Signbit() {
				result := big.NewFloat(0)
				result.Neg(result) // This should create -0

				return result
			}

			return big.NewFloat(0)
		}

		// Pow(±0, y) = +0 for finite y > 0 and not an odd integer
		return big.NewFloat(0)
	}

	// Pow(-1, ±Inf) = 1
	if x.Cmp(negOne) == 0 && yIsInf {
		return big.NewFloat(1)
	}

	// Handle infinite exponent cases
	if yIsInf {
		absX := new(big.Float).Abs(x)
		if y.Signbit() {
			// y = -Inf
			if absX.Cmp(one) > 0 {
				// Pow(x, -Inf) = +0 for |x| > 1
				return big.NewFloat(0)
			} else if absX.Cmp(one) < 0 {
				// Pow(x, -Inf) = +Inf for |x| < 1
				result := big.NewFloat(0)
				result.SetInf(false)

				return result
			}
		} else {
			// y = +Inf
			if absX.Cmp(one) > 0 {
				// Pow(x, +Inf) = +Inf for |x| > 1
				result := big.NewFloat(0)
				result.SetInf(false)

				return result
			} else if absX.Cmp(one) < 0 {
				// Pow(x, +Inf) = +0 for |x| < 1
				return big.NewFloat(0)
			}
		}
	}

	// Handle infinite base cases
	if xIsInf {
		if x.Signbit() {
			// Pow(-Inf, y) = Pow(-0, -y) - implement directly to avoid recursion
			negY := new(big.Float).Neg(y)
			negYSign := negY.Sign()
			if negYSign < 0 {
				// -y < 0, so y > 0
				if negY.IsInt() && isOddInteger(negY) {
					// Pow(-0, -y) = -Inf for -y an odd integer < 0 (i.e., y odd positive)
					result := big.NewFloat(0)
					result.SetInf(true) // -Inf

					return result
				}
				// Pow(-0, -y) = +Inf for -y < 0 and not odd integer
				result := big.NewFloat(0)
				result.SetInf(false) // +Inf

				return result
			} else if negYSign > 0 {
				// -y > 0, so y < 0
				if negY.IsInt() && isOddInteger(negY) {
					// Pow(-0, -y) = -0 for -y an odd integer > 0 (i.e., y odd negative)
					result := big.NewFloat(0)
					result.Neg(result)

					return result
				}

				// Pow(-0, -y) = +0 for -y > 0 and not odd integer
				return big.NewFloat(0)
			}

			// -y = 0, so y = 0: Pow(-0, 0) = 1
			return big.NewFloat(1)
		}

		// +Inf
		if y.Sign() > 0 {
			// Pow(+Inf, y) = +Inf for y > 0
			result := big.NewFloat(0)
			result.SetInf(false)

			return result
		}

		// Pow(+Inf, y) = +0 for y < 0
		return big.NewFloat(0)
	}

	// Pow(x, y) = NaN for finite x < 0 and finite non-integer y
	if x.Sign() < 0 && !y.IsInt() {
		// Return +Inf to indicate undefined behavior (big.Float has no NaN)
		result := big.NewFloat(0)
		result.SetInf(false)

		return result
	}

	// For integer exponents, use PowInt
	if y.IsInt() && y.Abs(y).Cmp(big.NewFloat(1000000)) < 0 {
		yInt, _ := y.Int64()

		return PowInt(x, yInt)
	}

	// For non-integer exponents, use exp(y * ln(x))
	// x^y = e^(y * ln(x))
	logX := Log(x)
	if logX == nil {
		// If log fails, return +Inf to indicate undefined behavior
		result := big.NewFloat(0)
		result.SetInf(false)

		return result
	}

	// Compute y * log(x)
	yLogX := new(big.Float).Mul(y, logX)

	// Compute e^(y * log(x))
	return Exp(yLogX)
}

// PowFloat64 returns x**y, the base-x exponential of y from float64 inputs
// returning a *big.Float.  Useful when x**y would overflow a float64 normally.
//
// Special cases are:
//
// PowFloat64(x, ±0) = 1 for any x
// PowFloat64(1, y) = 1 for any y
// PowFloat64(x, 1) = x for any x
// PowFloat64(NaN, y) = ErrNaN
// PowFloat64(x, NaN) = ErrNaN
// Pow(±0, -Inf) = +Inf
// Pow(±0, +Inf) = +0
// Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
// Pow(±0, y) = +0 for finite y > 0 and not an odd integer
// Pow(±0, y) = ±Inf for y an odd integer < 0
// Pow(±0, y) = ±0 for y an odd integer > 0
// Pow(-1, ±Inf) = 1
// Pow(x, +Inf) = +Inf for |x| > 1
// Pow(x, -Inf) = +0 for |x| > 1
// Pow(x, +Inf) = +0 for |x| < 1
// Pow(x, -Inf) = +Inf for |x| < 1
// Pow(+Inf, y) = +Inf for y > 0
// Pow(+Inf, y) = +0 for y < 0
// Pow(-Inf, y) = Pow(-0, -y)
// Pow(x, y) = ErrNan for finite x < 0 and finite non-integer y
func PowFloat64(x, y float64) *big.Float {
	// Handle NaN inputs - return +Inf to indicate undefined (big.Float has no NaN)
	if math.IsNaN(x) || math.IsNaN(y) {
		result := big.NewFloat(0)
		result.SetInf(false)

		return result
	}

	// Handle special cases for float64 inputs before converting to big.Float
	if y == 0 {
		return big.NewFloat(1)
	}
	if x == 1 {
		return big.NewFloat(1)
	}

	// Handle infinity cases
	if math.IsInf(x, 0) || math.IsInf(y, 0) {
		// Convert to big.Float and let Pow handle the infinity cases
		bigX := big.NewFloat(x)
		bigY := big.NewFloat(y)
		if math.IsInf(x, 1) {
			bigX.SetInf(false) // +Inf
		} else if math.IsInf(x, -1) {
			bigX.SetInf(true) // -Inf
		}
		if math.IsInf(y, 1) {
			bigY.SetInf(false) // +Inf
		} else if math.IsInf(y, -1) {
			bigY.SetInf(true) // -Inf
		}

		return Pow(bigX, bigY)
	}

	// Convert inputs to *big.Float and use Pow
	return Pow(big.NewFloat(x), big.NewFloat(y))
}

// PowInt returns x**n for integer exponent n using repeated multiplication.
func PowInt(x *big.Float, n int64) *big.Float {
	if n == 0 {
		return big.NewFloat(1)
	}

	result := big.NewFloat(1)
	base := new(big.Float).Copy(x)

	if n < 0 {
		// For negative exponent, compute 1/x^|n|
		n = -n
		base.Quo(big.NewFloat(1), base)
	}

	// Binary exponentiation
	for n > 0 {
		if n&1 == 1 {
			result.Mul(result, base)
		}
		base.Mul(base, base)
		n >>= 1
	}

	return result
}
