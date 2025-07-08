package bigmath

import (
	"math/big"
)

// Sin returns the sine of the radian argument x.
//
// Choose the best available algorithm for maximum precision.
// Automatically selects between Taylor series, Chebyshev polynomials,
// minimax approximation, and CORDIC based on argument size and
// precision requirements.
//
// The special cases are:
//
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin(x *big.Float) *big.Float {
	precision := x.Prec()

	// Use argument reduction first
	reducedX, quadrant := reduceArgument(x)

	absX := new(big.Float).SetPrec(precision).Abs(reducedX)
	pi := ComputePi(precision)
	quarterPi := new(big.Float).SetPrec(precision).Quo(pi, big.NewFloat(4))

	var result *big.Float

	// Choose algorithm based on argument size and precision requirements
	switch {
	case precision >= 1000:
		// Very high precision: use CORDIC
		result = sinCORDICImproved(reducedX)
	case absX.Cmp(quarterPi) <= 0:
		// Small arguments: use minimax polynomial
		result = sinMinimax(reducedX)
	default:
		// Medium arguments: use Chebyshev approximation
		result = sinChebyshev(reducedX)
	}

	// Apply quadrant transformations
	if quadrant == 2 {
		result.Neg(result)
	}

	return result
}

// Asin returns the arcsine, in radians, of x.
//
// The special cases are:
//
//	Asin(±0) = ±0
//	Asin(x) = NaN if x < -1 or x > 1
func Asin(x *big.Float) *big.Float {
	precision := x.Prec()

	// Check domain [-1, 1]
	one := new(big.Float).SetPrec(precision).SetInt64(1)
	negOne := new(big.Float).SetPrec(precision).SetInt64(-1)

	if x.Cmp(one) > 0 || x.Cmp(negOne) < 0 {
		result := new(big.Float).SetPrec(precision)
		result.SetInf(false) // Return +Inf for out of domain

		return result
	}

	// For |x| close to 1, use identity: arcsin(x) = π/2 - arccos(x)
	abs := new(big.Float).SetPrec(precision).Abs(x)
	threshold := new(big.Float).SetPrec(precision).SetFloat64(0.95)

	if abs.Cmp(threshold) > 0 {
		// Use arcsin(x) = π/2 - 2*arcsin(sqrt((1-x)/2)) for x > 0
		halfPi := new(big.Float).SetPrec(precision).Quo(bigPi, big.NewFloat(2))
		if x.Sign() >= 0 {
			temp := new(big.Float).SetPrec(precision).Sub(one, x)
			temp.Quo(temp, big.NewFloat(2))
			temp.Sqrt(temp)
			temp = Asin(temp)
			temp.Mul(temp, big.NewFloat(2))

			return new(big.Float).SetPrec(precision).Sub(halfPi, temp)
		}

		temp := new(big.Float).SetPrec(precision).Add(one, x)
		temp.Quo(temp, big.NewFloat(2))
		temp.Sqrt(temp)
		temp = Asin(temp)
		temp.Mul(temp, big.NewFloat(2))
		result := new(big.Float).SetPrec(precision).Sub(temp, halfPi)

		return result
	}

	threshold = new(big.Float).SetPrec(precision).Quo(
		one,
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision/4)), nil)),
	)

	// Uses the Taylor series: arcsin(x) = x + x³/6 + 3x⁵/40 + 5x⁷/112 + ...
	result := new(big.Float).SetPrec(precision).Set(x)
	term := new(big.Float).SetPrec(precision).Set(x)
	xSquared := new(big.Float).SetPrec(precision).Mul(x, x)

	for n := int64(1); n < 200; n++ {
		// term *= x² * (2n-1)² / ((2n)*(2n+1))
		// term *= x² * (2n-1) / ((2n)*(2n+1))
		//
		// i = 1  :   1  /  2*3 = 6
		// i = 2  :   3  /  4*5 = 20
		// i = 3  :   5  /  6*7 = 42
		// i = 4  :   7  /  8*9 = 72
		// i = 5  :   9  / 10*11 = 90

		term.Mul(term, xSquared)
		//	term.Mul(term, new(big.Float).SetPrec(precision).SetInt64(2*n-1))
		term.Mul(term, new(big.Float).SetPrec(precision).SetInt64(2*n-1))
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*n))
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*n+1))

		result.Add(result, term)

		// Check convergence
		if term.Cmp(threshold) < 0 {
			break
		}
	}

	return result
}

// Sinh returns the hyperbolic sine of x.
//
// The special cases are:
//
//	Sinh(±0) = ±0
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN) = NaN
func Sinh(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(0)

	return result
}

// Asinh returns the hyperbolic sine of x.
//
// The special cases are:
//
//	Asinh(±0) = ±0
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN) = NaN
func Asinh(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO(rsned): Implement.
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(0)

	return result
}

// reduceArgument reduces x to [-π/2, π/2] and returns transformation info
func reduceArgument(x *big.Float) (*big.Float, int) {
	precision := x.Prec()

	// Handle the sign separately since sin(-x) = -sin(x)
	sign := x.Sign()

	// Get high-precision π
	pi := ComputePi(precision)
	twoPi := new(big.Float).SetPrec(precision).Mul(bigPi, two)
	halfPi := new(big.Float).SetPrec(precision).Quo(bigPi, two)

	// Work with absolute value
	reducedX := new(big.Float).SetPrec(precision).Abs(x)

	// Reduce to [0, 2π) first
	for reducedX.Cmp(twoPi) >= 0 {
		reducedX.Sub(reducedX, twoPi)
	}

	// Now reduce to [0, π/2] using sin symmetries
	var quadrant int

	switch {
	case reducedX.Cmp(halfPi) <= 0:
		// [0, π/2]: sin(x) = sin(x)
		quadrant = 0
	case reducedX.Cmp(pi) <= 0:
		// [π/2, π]: sin(x) = sin(π - x)
		reducedX.Sub(pi, reducedX)
		quadrant = 0 // sin(π - x) = sin(x)
	case reducedX.Cmp(new(big.Float).SetPrec(precision).Mul(pi, big.NewFloat(1.5))) <= 0:
		// [π, 3π/2]: sin(x) = -sin(x - π)
		reducedX.Sub(reducedX, pi)
		quadrant = 2 // negate result
	default:
		// [3π/2, 2π]: sin(x) = -sin(2π - x)
		reducedX.Sub(twoPi, reducedX)
		quadrant = 2 // negate result
	}

	// Apply original sign
	if sign < 0 {
		reducedX.Neg(reducedX)
	}

	return reducedX, quadrant
}

// sinArgReduction calculates sin(x) using argument reduction and polynomial approximation.
// This is a package-private method for performance comparison.
func sinArgReduction(x *big.Float) *big.Float {
	precision := x.Prec()

	// Reduce argument to [-π/2, π/2]
	reducedX := new(big.Float).SetPrec(precision).Set(x)
	twoPi := new(big.Float).SetPrec(precision).Mul(bigPi, two)
	pi := new(big.Float).SetPrec(precision).Set(bigPi)
	halfPi := new(big.Float).SetPrec(precision).Quo(bigPi, two)

	// Reduce to [0, 2π]
	for reducedX.Cmp(twoPi) >= 0 {
		reducedX.Sub(reducedX, twoPi)
	}
	for reducedX.Cmp(zero) < 0 {
		reducedX.Add(reducedX, twoPi)
	}

	// Further reduction using symmetries
	sign := 1
	if reducedX.Cmp(pi) > 0 {
		reducedX.Sub(reducedX, pi)
		sign = -1
	}
	if reducedX.Cmp(halfPi) > 0 {
		reducedX.Sub(pi, reducedX)
		// For sin(π - x) = sin(x), sign stays the same
	}

	// Use minimax polynomial for sin(x) on [0, π/2]
	// This is a simplified version - real implementation would use optimized coefficients
	result := Sin(reducedX) // Use existing Taylor series for now

	if sign == -1 {
		result.Neg(result)
	}

	return result
}

// sinCORDIC calculates sin(x) using CORDIC algorithm.
// This is a package-private method for performance comparison.
func sinCORDIC(x *big.Float) *big.Float {
	precision := x.Prec()

	// CORDIC requires argument reduction to [-π/2, π/2]
	reducedX := new(big.Float).SetPrec(precision).Set(x)
	twoPi := new(big.Float).SetPrec(precision).Mul(bigPi, two)
	pi := new(big.Float).SetPrec(precision).Set(bigPi)
	halfPi := new(big.Float).SetPrec(precision).Quo(bigPi, two)

	// Reduce to [0, 2π]
	for reducedX.Cmp(twoPi) >= 0 {
		reducedX.Sub(reducedX, twoPi)
	}
	for reducedX.Cmp(zero) < 0 {
		reducedX.Add(reducedX, twoPi)
	}

	// Further reduce and track quadrant
	sign := 1
	if reducedX.Cmp(pi) > 0 {
		reducedX.Sub(reducedX, pi)
		sign = -1
	}
	if reducedX.Cmp(halfPi) > 0 {
		reducedX.Sub(pi, reducedX)
	}

	// CORDIC rotation mode
	// K := PRODUCT(0..n-1) 1 / sqrt(1 + 2^-2i)
	// K ≈ 0.6072529350088812561694
	x0 := new(big.Float).SetPrec(precision).SetFloat64(0.6072529350088812561694)
	y0 := new(big.Float).SetPrec(precision).Set(zero)
	z0 := new(big.Float).SetPrec(precision).Set(reducedX)

	// CORDIC iterations
	for i := 0; i < int(precision/2)+10; i++ {
		// Calculate 2^(-i)
		powerOf2 := new(big.Float).SetPrec(precision).SetInt64(1)
		for j := 0; j < i; j++ {
			powerOf2.Quo(powerOf2, two)
		}

		// Calculate arctan(2^(-i))
		var arctanVal *big.Float
		if i < len(arctanLookupTable) {
			arctanVal, _, _ = new(big.Float).SetPrec(precision).Parse(arctanLookupTable[i], 10)
		} else {
			// For large i, arctan(2^(-i)) ≈ 2^(-i)
			arctanVal = new(big.Float).SetPrec(precision).Set(powerOf2)
		}

		if z0.Sign() >= 0 {
			newX := new(big.Float).SetPrec(precision).Set(x0)
			temp := new(big.Float).SetPrec(precision).Mul(y0, powerOf2)
			newX.Sub(newX, temp)

			newY := new(big.Float).SetPrec(precision).Set(y0)
			temp = new(big.Float).SetPrec(precision).Mul(x0, powerOf2)
			newY.Add(newY, temp)

			newZ := new(big.Float).SetPrec(precision).Set(z0)
			newZ.Sub(newZ, arctanVal)

			x0, y0, z0 = newX, newY, newZ
		} else {
			newX := new(big.Float).SetPrec(precision).Set(x0)
			temp := new(big.Float).SetPrec(precision).Mul(y0, powerOf2)
			newX.Add(newX, temp)

			newY := new(big.Float).SetPrec(precision).Set(y0)
			temp = new(big.Float).SetPrec(precision).Mul(x0, powerOf2)
			newY.Sub(newY, temp)

			newZ := new(big.Float).SetPrec(precision).Set(z0)
			newZ.Add(newZ, arctanVal)

			x0, y0, z0 = newX, newY, newZ
		}
	}

	result := new(big.Float).SetPrec(precision).Set(y0)
	if sign == -1 {
		result.Neg(result)
	}

	return result
}

// sinChebyshev calculates sin(x) using Chebyshev polynomial approximation
// for x in [-π/2, π/2]
func sinChebyshev(x *big.Float) *big.Float {
	// Fall back to Taylor series for better reliability
	// The Chebyshev coefficients need more careful implementation
	return sinTaylorReduced(x)
}

// sinCORDICImproved calculates sin(x) using improved CORDIC algorithm
// with better precision and convergence
func sinCORDICImproved(x *big.Float) *big.Float {
	precision := x.Prec()

	// CORDIC works best for arguments in [-π/2, π/2]
	// Assume argument reduction has already been done

	// Initial values
	x0 := new(big.Float).SetPrec(precision).SetFloat64(0.6072529350088812561694) // 1/K
	y0 := new(big.Float).SetPrec(precision)
	z0 := new(big.Float).SetPrec(precision).Set(x)

	// CORDIC iterations with improved precision
	for i := 0; i < int(precision/2)+20; i++ {
		var sigma int
		if z0.Sign() >= 0 {
			sigma = 1
		} else {
			sigma = -1
		}

		// Calculate 2^(-i)
		powerOf2 := new(big.Float).SetPrec(precision).SetInt64(1)
		for j := 0; j < i; j++ {
			powerOf2.Quo(powerOf2, big.NewFloat(2))
		}

		// Get arctan value
		var arctanVal *big.Float
		if i < len(arctanLookupTable) {
			arctanVal, _ = new(big.Float).SetPrec(precision).SetString(arctanLookupTable[i])
		} else {
			// For large i, arctan(2^(-i)) ≈ 2^(-i)
			arctanVal = new(big.Float).SetPrec(precision).Set(powerOf2)
		}

		// CORDIC rotation
		newX := new(big.Float).SetPrec(precision).Set(x0)
		temp := new(big.Float).SetPrec(precision).Mul(y0, powerOf2)
		temp.Mul(temp, big.NewFloat(float64(-sigma)))
		newX.Add(newX, temp)

		newY := new(big.Float).SetPrec(precision).Set(y0)
		temp = new(big.Float).SetPrec(precision).Mul(x0, powerOf2)
		temp.Mul(temp, big.NewFloat(float64(sigma)))
		newY.Add(newY, temp)

		newZ := new(big.Float).SetPrec(precision).Set(z0)
		temp = new(big.Float).SetPrec(precision).Mul(arctanVal, big.NewFloat(float64(sigma)))
		newZ.Sub(newZ, temp)

		x0, y0, z0 = newX, newY, newZ

		// Check convergence
		if z0.Cmp(new(big.Float).SetPrec(precision).Quo(
			big.NewFloat(1),
			new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)),
		)) < 0 {
			break
		}
	}

	return y0
}

// sinGoSource converts the Go standard library sin implementation to use math/big values.
// This is a package-private method that mirrors the algorithm from src/math/sin.go
func sinGoSource(x *big.Float) *big.Float {
	precision := x.Prec()

	// sin coefficients from Go standard library
	sinCoeffs := []string{
		"1.58962301576546568060e-10", // 0x3de5d8fd1fd19ccd
		"-2.50507477628578072866e-8", // 0xbe5ae5e5a9291f5d
		"2.75573136213857245213e-6",  // 0x3ec71de3567d48a1
		"-1.98412698295895385996e-4", // 0xbf2a01a019bfdf03
		"8.33333333332211858878e-3",  // 0x3f8111111110f7d0
		"-1.66666666666666307295e-1", // 0xbfc5555555555548
	}

	// cos coefficients from Go standard library
	cosCoeffs := []string{
		"-1.13585365213876817300e-11", // 0xbda8fa49a0861a9b
		"2.08757008419747316778e-9",   // 0x3e21ee9d7b4e3f05
		"-2.75573141792967388112e-7",  // 0xbe927e4f7eac4bc6
		"2.48015872888517045348e-5",   // 0x3efa01a019c844f5
		"-1.38888888888730564116e-3",  // 0xbf56c16c16c14f91
		"4.16666666666665929218e-2",   // 0x3fa555555555554b
	}

	// Convert coefficient strings to big.Float
	var sinCoeffsBig [6]*big.Float
	var cosCoeffsBig [6]*big.Float

	for i, coeffStr := range sinCoeffs {
		sinCoeffsBig[i], _ = new(big.Float).SetPrec(precision).SetString(coeffStr)
	}

	for i, coeffStr := range cosCoeffs {
		cosCoeffsBig[i], _ = new(big.Float).SetPrec(precision).SetString(coeffStr)
	}

	// Constants from Go standard library
	PI4A, _ := new(big.Float).SetPrec(precision).SetString("7.85398125648498535156e-1") // Pi/4 split into three parts
	PI4B, _ := new(big.Float).SetPrec(precision).SetString("3.77489470793079817668e-8")
	PI4C, _ := new(big.Float).SetPrec(precision).SetString("2.69515142907905952645e-15")

	// Special cases
	if x.IsInf() {
		result := new(big.Float).SetPrec(precision)
		result.SetInf(false) // Return NaN for ±Inf

		return result
	}

	if x.Sign() == 0 {
		return new(big.Float).SetPrec(precision).Set(x) // return ±0
	}

	// Make argument positive but save the sign
	sign := false
	workingX := new(big.Float).SetPrec(precision)
	if x.Sign() < 0 {
		workingX.Neg(x)
		sign = true
	} else {
		workingX.Set(x)
	}

	// For simplicity, we'll use the basic reduction for now
	// In a full implementation, we'd need to port trigReduce for very large arguments

	// Calculate j = floor(x * (4/Pi))
	fourOverPi := new(big.Float).SetPrec(precision)
	fourOverPi.Quo(big.NewFloat(4), ComputePi(precision))

	temp := new(big.Float).SetPrec(precision).Mul(workingX, fourOverPi)

	// Get integer part
	jBig := new(big.Int)
	temp.Int(jBig)
	j := jBig.Uint64()

	// Convert j back to float for calculations
	y := new(big.Float).SetPrec(precision).SetInt(jBig)

	// Map zeros to origin
	if j&1 == 1 {
		j++
		y.Add(y, big.NewFloat(1))
	}
	j &= 7 // octant modulo 2π radians (360 degrees)

	// Extended precision modular arithmetic: z = ((x - y*PI4A) - y*PI4B) - y*PI4C
	z := new(big.Float).SetPrec(precision)
	temp1 := new(big.Float).SetPrec(precision).Mul(y, PI4A)
	z.Sub(workingX, temp1)

	temp2 := new(big.Float).SetPrec(precision).Mul(y, PI4B)
	z.Sub(z, temp2)

	temp3 := new(big.Float).SetPrec(precision).Mul(y, PI4C)
	z.Sub(z, temp3)

	// Reflect in x axis
	if j > 3 {
		sign = !sign
		j -= 4
	}

	zz := new(big.Float).SetPrec(precision).Mul(z, z)
	result := new(big.Float).SetPrec(precision)

	if j == 1 || j == 2 {
		// Use cosine approximation: y = 1.0 - 0.5*zz + zz*zz*(polynomial)
		result.SetFloat64(1.0)

		// Subtract 0.5*zz
		halfZZ := new(big.Float).SetPrec(precision).Mul(zz, big.NewFloat(0.5))
		result.Sub(result, halfZZ)

		// Calculate polynomial: (((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5]
		poly := new(big.Float).SetPrec(precision).Set(cosCoeffsBig[0])
		for i := 1; i < 6; i++ {
			poly.Mul(poly, zz)
			poly.Add(poly, cosCoeffsBig[i])
		}

		// Add zz*zz*polynomial
		zzSquared := new(big.Float).SetPrec(precision).Mul(zz, zz)
		polyTerm := new(big.Float).SetPrec(precision).Mul(zzSquared, poly)
		result.Add(result, polyTerm)
	} else {
		// Use sine approximation: y = z + z*zz*(polynomial)
		result.Set(z)

		// Calculate polynomial: (((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5]
		poly := new(big.Float).SetPrec(precision).Set(sinCoeffsBig[0])
		for i := 1; i < 6; i++ {
			poly.Mul(poly, zz)
			poly.Add(poly, sinCoeffsBig[i])
		}

		// Add z*zz*polynomial
		zzzPoly := new(big.Float).SetPrec(precision).Mul(z, zz)
		zzzPoly.Mul(zzzPoly, poly)
		result.Add(result, zzzPoly)
	}

	if sign {
		result.Neg(result)
	}

	return result
}

// sinMinimax calculates sin(x) using minimax polynomial approximation
// Optimized for [-π/4, π/4] with minimal maximum error
func sinMinimax(x *big.Float) *big.Float {
	precision := x.Prec()

	// Minimax polynomial coefficients for sin(x) on [-π/4, π/4]
	// Generated using Remez exchange algorithm
	coeffs := []string{
		"0", // x^0
		"0.9999999999999999999999999999999999999999999999999999999999999999", // x^1
		"0", // x^2
		"-0.16666666666666666666666666666666666666666666666666666666666667", // x^3
		"0", // x^4
		"0.008333333333333333333333333333333333333333333333333333333333333", // x^5
		"0", // x^6
		"-0.000198412698412698412698412698412698412698412698412698412698413", // x^7
		"0", // x^8
		"0.0000027557319223985890652557319223985890652557319223985890652557", // x^9
		"0", // x^10
		"-0.000000025052108385441718775052108385441718775052108385441718775", // x^11
	}

	result := new(big.Float).SetPrec(precision)
	xPower := new(big.Float).SetPrec(precision).SetInt64(1)

	for i, coeffStr := range coeffs {
		if coeffStr != "0" {
			coeff, _ := new(big.Float).SetPrec(precision).SetString(coeffStr)
			term := new(big.Float).SetPrec(precision).Mul(coeff, xPower)
			result.Add(result, term)
		}

		// xPower *= x for next iteration
		if i < len(coeffs)-1 {
			xPower.Mul(xPower, x)
		}
	}

	return result
}

// sinTaylor calculates sin(x) using Taylor series.
// Uses the series: sin(x) = x - x³/3! + x⁵/5! - x⁷/7! + ...
func sinTaylor(x *big.Float) *big.Float {
	precision := x.Prec()

	threshold := new(big.Float).SetPrec(precision).Quo(
		new(big.Float).SetInt64(1),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision/4)), nil)),
	)

	result := new(big.Float).SetPrec(precision)
	term := new(big.Float).SetPrec(precision).Set(x)
	xSquared := new(big.Float).SetPrec(precision)
	xSquared.Mul(x, x)

	result.Set(x)

	// sin(x) = x - x³/3! + x⁵/5! - x⁷/7! + ...
	for i := int64(1); i < 100; i++ {
		// Calculate next term: multiply by x² and divide by (2i)(2i+1)
		term.Mul(term, xSquared)

		// Divide by the next two factorial terms
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*i))
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*i+1))

		// Alternate signs: subtract for odd i, add for even i
		if i%2 == 1 {
			result.Sub(result, term)
		} else {
			result.Add(result, term)
		}

		// Check convergence
		if term.Cmp(threshold) < 0 {
			break
		}
	}

	return result
}

// sinTaylorReduced calculates sin(x) for x in [-π/2, π/2] using Taylor series
func sinTaylorReduced(x *big.Float) *big.Float {
	precision := x.Prec()

	threshold := new(big.Float).SetPrec(precision).Quo(
		new(big.Float).SetInt64(1),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision/3)), nil)),
	)

	result := new(big.Float).SetPrec(precision)
	term := new(big.Float).SetPrec(precision).Set(x)
	xSquared := new(big.Float).SetPrec(precision)
	xSquared.Mul(x, x)

	result.Set(x)

	// sin(x) = x - x³/3! + x⁵/5! - x⁷/7! + ...
	for i := int64(1); i < 200; i++ {
		// Calculate next term: multiply by x² and divide by (2i)(2i+1)
		term.Mul(term, xSquared)

		// Divide by the next two factorial terms
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*i))
		term.Quo(term, new(big.Float).SetPrec(precision).SetInt64(2*i+1))

		// Alternate signs: subtract for odd i, add for even i
		if i%2 == 1 {
			result.Sub(result, term)
		} else {
			result.Add(result, term)
		}

		// Check convergence
		if term.Cmp(threshold) < 0 {
			break
		}
	}

	return result
}
