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

import "math/big"

// Sec calculates sec(x) as 1/cos(x).
func Sec(x *big.Float) *big.Float {
	precision := x.Prec()
	cosX := Cos(x)

	result := new(big.Float).SetPrec(precision)
	one := new(big.Float).SetPrec(precision).SetInt64(1)
	result.Quo(one, cosX)

	return result
}

// Asec calculates inverse secant using the identity asec(x) = acos(1/x).
func Asec(x *big.Float) *big.Float {
	precision := x.Prec()

	// Check domain |x| >= 1
	one := new(big.Float).SetPrec(precision).SetInt64(1)
	abs := new(big.Float).SetPrec(precision).Abs(x)

	if abs.Cmp(one) < 0 {
		result := new(big.Float).SetPrec(precision)
		result.SetInf(false) // Return NaN for out of domain

		return result
	}

	// Use identity: arcsec(x) = arccos(1/x)
	reciprocal := new(big.Float).SetPrec(precision).Quo(one, x)

	return Acos(reciprocal)
}

// Sech calculates hyperbolic secant using the formula: sech(x) = 1/cosh(x)
// This is a placeholder implementation.
func Sech(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(1) // Placeholder (sech(0) = 1)

	return result
}

// Asech calculates inverse hyperbolic secant using the formula: asech(x) = acosh(1/x)
// This is a placeholder implementation.
func Asech(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(1) // Placeholder (arsech(0) = 1)

	return result
}

// secSeries calculates sec(x) using direct series expansion.
// Uses the series: sec(x) = 1 + x²/2 + 5x⁴/24 + 61x⁶/720 + ...
// This is a package-private method for performance comparison.
func secSeries(x *big.Float) *big.Float {
	precision := x.Prec()

	// Euler numbers for sec(x) series (simplified approximation)
	result := new(big.Float).SetPrec(precision).SetInt64(1)
	term := new(big.Float).SetPrec(precision).SetInt64(1)
	xSquared := new(big.Float).SetPrec(precision).Mul(x, x)

	// Coefficients for sec(x) = 1 + x²/2! + 5x⁴/4! + 61x⁶/6! + ...
	eulerCoeffs := []float64{1.0 / 2, 5.0 / 24, 61.0 / 720, 1385.0 / 40320}

	for i := 0; i < len(eulerCoeffs) && i < 20; i++ {
		term.Mul(term, xSquared)
		coeff := new(big.Float).SetPrec(precision).SetFloat64(eulerCoeffs[i])
		termValue := new(big.Float).SetPrec(precision).Mul(term, coeff)

		result.Add(result, termValue)

		// Check convergence
		convergenceThreshold := new(big.Float).SetPrec(precision).Quo(
			new(big.Float).SetInt64(1),
			new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision/4)), nil)),
		)
		if termValue.Cmp(convergenceThreshold) < 0 {
			break
		}
	}

	return result
}
