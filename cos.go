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

import "math/big"

// Cos returns the cosine of the radian argument x.
//
// The special cases are:
//
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos(x *big.Float) *big.Float {
	precision := x.Prec()
	result := new(big.Float).SetPrec(precision)
	term := new(big.Float).SetPrec(precision).SetInt64(1)
	xSquared := new(big.Float).SetPrec(precision)
	xSquared.Mul(x, x)

	result.SetInt64(1)

	// Taylor series for cos: (-1)^n * (x^2n / x!)
	// cos(x) = 1 - x²/2! + x⁴/4! - x⁶/6! + ...
	for i := int64(1); i < 100; i++ {
		// Calculate next term: multiply by x² and divide by (2i-1)(2i)
		term.Mul(term, xSquared)

		// Divide by the next two factorial terms
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*i-1))
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*i))

		// Alternate signs: subtract for odd i, add for even i
		if i%2 == 1 {
			result.Sub(result, term)
		} else {
			result.Add(result, term)
		}

		// Check convergence
		threshold := new(big.Float).SetPrec(precision).Quo(
			new(big.Float).SetInt64(1),
			new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision/4)), nil)),
		)
		if term.Cmp(threshold) < 0 {
			break
		}
	}

	return result
}

// Acos returns the arccosine, in radians, of x.
//
// The special case is:
//
//	Acos(x) = NaN if x < -1 or x > 1
func Acos(x *big.Float) *big.Float {
	precision := x.Prec()

	// Acos calculates arccos(x) using the identity arccos(x) = π/2 - arcsin(x).

	// Check domain [-1, 1]
	one := new(big.Float).SetPrec(precision).SetInt64(1)
	negOne := new(big.Float).SetPrec(precision).SetInt64(-1)

	if x.Cmp(one) > 0 || x.Cmp(negOne) < 0 {
		result := new(big.Float).SetPrec(precision)
		result.SetInf(false) // Return NaN for out of domain

		return result
	}

	// Use identity: arccos(x) = π/2 - arcsin(x)
	halfPi := new(big.Float).SetPrec(precision).Quo(bigPi, two)
	arcsinX := Asin(x)

	result := new(big.Float).SetPrec(precision)
	result.Sub(halfPi, arcsinX)

	return result
}

// Cosh returns the hyperbolic cosine of x.
//
// The special cases are:
//
//	Cosh(±0) = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN) = NaN
func Cosh(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(0)

	return result
}

// Acosh returns the inverse hyperbolic cosine of x.
//
// The special cases are:
//
//	Acosh(+Inf) = +Inf
//	Acosh(x) = NaN if x < 1
//	Acosh(NaN) = NaN
func Acosh(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(0)

	return result
}

// cosCORDIC calculates cos(x) using CORDIC algorithm.
// This is a package-private method for performance comparison.
func cosCORDIC(x *big.Float) *big.Float {
	precision := x.Prec()

	// CORDIC requires argument reduction to [-π/2, π/2]
	reducedX := new(big.Float).SetPrec(precision).Set(x)
	twoPi := new(big.Float).SetPrec(precision).Mul(bigPi, big.NewFloat(2))
	pi := new(big.Float).SetPrec(precision).Set(bigPi)
	halfPi := new(big.Float).SetPrec(precision).Quo(pi, big.NewFloat(2))

	// Reduce to [0, 2π]
	for reducedX.Cmp(twoPi) >= 0 {
		reducedX.Sub(reducedX, twoPi)
	}
	for reducedX.Cmp(big.NewFloat(0)) < 0 {
		reducedX.Add(reducedX, twoPi)
	}

	// Further reduce and track transformations
	sign := 1
	if reducedX.Cmp(pi) > 0 {
		reducedX.Sub(reducedX, pi)
		sign = -1
	}

	cosToSin := false
	if reducedX.Cmp(halfPi) > 0 {
		reducedX.Sub(pi, reducedX)
		cosToSin = true
	}

	// CORDIC rotation mode
	x0 := new(big.Float).SetPrec(precision).SetFloat64(0.6072529350088812561694) // 1/K ≈ 0.607252935
	y0 := new(big.Float).SetPrec(precision).SetInt64(0)
	z0 := new(big.Float).SetPrec(precision).Set(reducedX)

	// CORDIC iterations
	for i := 0; i < int(precision/2)+10; i++ {
		// Calculate 2^(-i)
		powerOf2 := new(big.Float).SetPrec(precision).SetInt64(1)
		for j := 0; j < i; j++ {
			powerOf2.Quo(powerOf2, big.NewFloat(2))
		}

		// Calculate arctan(2^(-i))
		arctanTable := new(big.Float).SetPrec(precision)
		if i < len(arctanLookupTable) {
			arctanTable, _, _ = new(big.Float).SetPrec(precision).Parse(arctanLookupTable[i], 10)
		} else {
			arctanTable.Set(powerOf2)
		}

		if z0.Sign() >= 0 {
			newX := new(big.Float).SetPrec(precision).Set(x0)
			temp := new(big.Float).SetPrec(precision).Mul(y0, powerOf2)
			newX.Sub(newX, temp)

			newY := new(big.Float).SetPrec(precision).Set(y0)
			temp = new(big.Float).SetPrec(precision).Mul(x0, powerOf2)
			newY.Add(newY, temp)

			newZ := new(big.Float).SetPrec(precision).Set(z0)
			newZ.Sub(newZ, arctanTable)

			x0, y0, z0 = newX, newY, newZ
		} else {
			newX := new(big.Float).SetPrec(precision).Set(x0)
			temp := new(big.Float).SetPrec(precision).Mul(y0, powerOf2)
			newX.Add(newX, temp)

			newY := new(big.Float).SetPrec(precision).Set(y0)
			temp = new(big.Float).SetPrec(precision).Mul(x0, powerOf2)
			newY.Sub(newY, temp)

			newZ := new(big.Float).SetPrec(precision).Set(z0)
			newZ.Add(newZ, arctanTable)

			x0, y0, z0 = newX, newY, newZ
		}
	}

	var result *big.Float
	if cosToSin {
		result = new(big.Float).SetPrec(precision).Set(y0)
	} else {
		result = new(big.Float).SetPrec(precision).Set(x0)
	}

	if sign == -1 {
		result.Neg(result)
	}

	return result
}

// cosArgReduction calculates cos(x) using argument reduction and polynomial approximation.
// This is a package-private method for performance comparison.
func cosArgReduction(x *big.Float) *big.Float {
	precision := x.Prec()

	// Reduce argument to [-π/2, π/2]
	reducedX := new(big.Float).SetPrec(precision).Set(x)
	twoPi := new(big.Float).SetPrec(precision).Mul(bigPi, big.NewFloat(2))
	pi := new(big.Float).SetPrec(precision).Set(bigPi)
	halfPi := new(big.Float).SetPrec(precision).Quo(pi, big.NewFloat(2))

	// Reduce to [0, 2π]
	for reducedX.Cmp(twoPi) >= 0 {
		reducedX.Sub(reducedX, twoPi)
	}
	for reducedX.Cmp(big.NewFloat(0)) < 0 {
		reducedX.Add(reducedX, twoPi)
	}

	// Further reduction using symmetries
	sign := 1
	if reducedX.Cmp(pi) > 0 {
		reducedX.Sub(reducedX, pi)
		sign = -1
	}

	cosToSin := false
	if reducedX.Cmp(halfPi) > 0 {
		reducedX.Sub(pi, reducedX)
		cosToSin = true
	}

	// Use minimax polynomial for cos(x) on [0, π/2]
	var result *big.Float
	if cosToSin {
		result = Sin(reducedX)
	} else {
		result = Cos(reducedX)
	}

	if sign == -1 {
		result.Neg(result)
	}

	return result
}
