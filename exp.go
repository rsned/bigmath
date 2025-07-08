// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigmath

import (
	"math/big"
)

// Exp returns e**x, the base-e exponential of x.
//
// The special cases are:
//
//	Exp(+Inf) = +Inf
//	Exp(NaN) = NaN
//
// Very large values no longer overflow to 0 or +Inf.
// Very small values no longer underflow to 1.
//
// For the time being, there is an explicit upper bound for x of ~700,000
// beyond which we choose to call it Infinite instead of looping excessively.
func Exp(x *big.Float) *big.Float {
	prec := x.Prec()
	result := new(big.Float).SetPrec(prec).SetInt64(1) // Start with 1
	term := new(big.Float).SetPrec(prec).SetInt64(1)   // Current term in series

	const largeX = 700000 // Our chosen limit for when to bail out.

	// Handle special cases
	if x.Sign() == 0 {
		return big.NewFloat(1).SetPrec(prec)
	}

	// For very large positive x, return +Inf
	xFloat, _ := x.Float64()
	if xFloat > largeX {
		result := big.NewFloat(0).SetPrec(prec)
		result.SetInf(false) // +Inf

		return result
	}

	// For very large negative x, return 0
	if xFloat < -largeX {
		return big.NewFloat(0).SetPrec(prec)
	}

	// compute using Taylor series
	// e^x = 1 + x + x^2/2! + x^3/3! + ...
	for i := 1; i < 200; i++ {
		// term = term * x / i
		term.Mul(term, x)
		term.Quo(term, new(big.Float).SetPrec(prec).SetFloat64(float64(i)))
		// Add term to result
		result.Add(result, term)

		// Check for convergence
		termAbs := new(big.Float).SetPrec(prec).Abs(term)

		// Check convergence
		if termAbs.Cmp(new(big.Float).SetPrec(prec).Quo(
			new(big.Float).SetInt64(1),
			new(big.Float).SetPrec(prec).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec/4)), nil)),
		)) < 0 {
			break
		}

		// Prevent infinite loops with very large terms
		absLimit := new(big.Float)
		absLimit, _ = absLimit.SetString("1e1000")
		if termAbs.Cmp(absLimit) > 0 {
			result.SetInf(false) // Return +Inf if terms get too large

			break
		}
	}

	return result
}

// Exp2 returns 2**x, the base-2 exponential of x.
//
// The special cases are:
//
//	Exp2(+Inf) = +Inf
//	Exp2(NaN) = NaN
func Exp2(x *big.Float) *big.Float {
	prec := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(prec).SetInt64(0)

	return result
}
