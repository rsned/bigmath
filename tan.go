package bigmath

import "math/big"

// Tan returns the tangent of the radian argument x.
//
// The special cases are:
//
//	Tan(±0) = ±0
//	Tan(±Inf) = NaN
//	Tan(NaN) = NaN
func Tan(x *big.Float) *big.Float {
	return tanNaive(x)
}

// Atan returns the arctangent, in radians, of x.
//
// The special cases are:
//
//	Atan(±0) = ±0
//	Atan(±Inf) = ±Pi/2
//	Atan(NaN) = NaN
func Atan(x *big.Float) *big.Float {
	prec := x.Prec()

	// For |x| > 1, use the identity:
	//    arctan(x) = π/2 - arctan(1/x) for x > 0
	// or arctan(x) = -π/2 - arctan(1/x) for x < 0
	abs := new(big.Float).SetPrec(prec).Abs(x)

	if abs.Cmp(one) > 0 {
		reciprocal := new(big.Float).SetPrec(prec).Quo(one, x)
		arctanRecip := Atan(reciprocal)

		if x.Sign() > 0 {
			return new(big.Float).SetPrec(prec).Sub(bigHalfPi, arctanRecip)
		}

		result := new(big.Float).SetPrec(prec).Add(bigHalfPi, arctanRecip)

		return result.Neg(result)
	}

	tolerance := new(big.Float).SetPrec(prec).Quo(one, new(big.Float).SetInt(new(big.Int).Exp(intTen, big.NewInt(int64(prec/4)), nil)))

	result := new(big.Float).SetPrec(prec)
	term := new(big.Float).SetPrec(prec).Set(x)
	xSquared := new(big.Float).SetPrec(prec)
	xSquared.Mul(x, x)

	result.Set(x)

	// arctan(x) = x - x³/3 + x⁵/5 - x⁷/7 + ...
	for i := int64(1); i < 500; i++ {
		term.Mul(term, xSquared)
		denominator := new(big.Float).SetPrec(prec).SetInt64(2*i + 1)
		termValue := new(big.Float).SetPrec(prec).Quo(term, denominator)

		if i%2 == 0 {
			result.Add(result, termValue)
		} else {
			result.Sub(result, termValue)
		}

		// Check convergence
		if termValue.Cmp(tolerance) < 0 {
			break
		}
	}

	return result
}

// Tanh returns the hyperbolic tangent of x.
//
// The special cases are:
//
//	Tanh(±0) = ±0
//	Tanh(±Inf) = ±1
//	Tanh(NaN) = NaN
func Tanh(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO: Implement efficient hyperbolic tangent calculation
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(0)

	return result
}

// Atanh returns the inverse hyperbolic arc tangent of x.
//
// The special cases are:
//
//	Atanh(1) = +Inf
//	Atanh(±0) = ±0
//	Atanh(-1) = -Inf
//	Atanh(x) = NaN if x < -1 or x > 1
//	Atanh(NaN) = NaN
func Atanh(x *big.Float) *big.Float {
	precision := x.Prec()

	// TODO: Implement efficient hyperbolic arc tangent calculation
	// For now, return a placeholder result
	result := new(big.Float).SetPrec(precision)
	result.SetInt64(0)

	return result
}

func tanNaive(x *big.Float) *big.Float {
	precision := x.Prec()
	sinX := Sin(x)
	cosX := Cos(x)

	result := new(big.Float).SetPrec(precision)
	result.Quo(sinX, cosX)

	return result
}

// tanContinuedFraction calculates tan(x) using continued fraction expansion.
// Uses: tan(x) = x/(1 - x²/(3 - x²/(5 - x²/(7 - ...))))
// This is a package-private method for performance comparison.
func tanContinuedFraction(x *big.Float) *big.Float {
	prec := x.Prec()

	// Track transformations needed
	negateResult := false
	reciprocal := false

	// Simple argument reduction to [-π/4, π/4] for better convergence
	reducedX := new(big.Float).SetPrec(prec).Set(x)

	// Handle negative values: tan(-x) = -tan(x)
	if reducedX.Sign() < 0 {
		negateResult = !negateResult
		reducedX.Abs(reducedX)
	}

	// Reduce argument using tan(x) periodicity and symmetries
	for reducedX.Cmp(bigPi) >= 0 {
		reducedX.Sub(reducedX, bigPi)
	}

	// Use tan(π - x) = -tan(x) for [π/2, π]
	if reducedX.Cmp(bigHalfPi) > 0 {
		reducedX.Sub(bigPi, reducedX)
		negateResult = !negateResult
	}

	// Use tan(π/2 - x) = cot(x) = 1/tan(x) for [π/4, π/2]
	if reducedX.Cmp(bigQuarterPi) > 0 {
		reducedX.Sub(bigHalfPi, reducedX)
		reciprocal = !reciprocal
	}

	// Now reducedX is in [0, π/4] where continued fraction converges well

	// For very small values, use Taylor expansion instead
	threshold := new(big.Float).SetPrec(prec).SetFloat64(0.1)
	if reducedX.Cmp(threshold) < 0 {
		result := tanTaylor(reducedX)
		// Apply transformations in reverse order
		if reciprocal {
			result.Quo(one, result)
		}

		if negateResult {
			result.Neg(result)
		}
	}

	xSquared := new(big.Float).SetPrec(prec).Mul(reducedX, reducedX)

	// Evaluate continued fraction from the bottom up
	// Use more conservative number of terms
	maxTerms := 30 + int(prec/16)
	if maxTerms > 100 {
		maxTerms = 100
	}

	cf := new(big.Float).SetPrec(prec).SetInt64(2*int64(maxTerms) + 1)

	for i := maxTerms; i >= 1; i-- {
		// cf = (2i+1) - x²/cf
		temp := new(big.Float).SetPrec(prec).Quo(xSquared, cf)
		cf.SetInt64(2*int64(i) + 1)
		cf.Sub(cf, temp)
	}

	// Final step: tan(x) = x / (1 - x²/cf)
	temp := new(big.Float).SetPrec(prec).Quo(xSquared, cf)
	denominator := one
	denominator.Sub(denominator, temp)

	result := new(big.Float).SetPrec(prec).Quo(reducedX, denominator)

	// Apply transformations in reverse order
	if reciprocal {
		result.Quo(one, result)
	}

	if negateResult {
		result.Neg(result)
	}

	return result
}

// tanCORDIC calculates tan(x) using CORDIC rotation mode (sin/cos approach).
// This is a package-private method for performance comparison.
func tanCORDIC(x *big.Float) *big.Float {
	// Use existing CORDIC sin/cos functions for better reliability
	precision := x.Prec()
	sinX := sinCORDIC(x)
	cosX := cosCORDIC(x)

	// Check for division by zero (cos(x) ≈ 0 near π/2 + nπ)
	zero := new(big.Float).SetPrec(precision).SetInt64(0)
	threshold := new(big.Float).SetPrec(precision).Quo(
		one,
		new(big.Float).SetInt(new(big.Int).Exp(intTen, big.NewInt(int64(precision/8)), nil)),
	)

	if cosX.Cmp(threshold) < 0 && cosX.Cmp(zero.Neg(threshold)) > 0 {
		// Return a large value indicating approach to infinity
		result := new(big.Float).SetPrec(precision)
		if sinX.Sign() > 0 {
			result.SetInf(false) // +Inf
		} else {
			result.SetInf(true) // -Inf
		}

		return result
	}

	result := new(big.Float).SetPrec(precision)
	result.Quo(sinX, cosX)

	return result
}

// First 15 coefficients for the tan taylor series.
var taylorTanSeriesCoeffs = []*big.Float{
	big.NewFloat(1.0 / 3),                                               // x^3
	big.NewFloat(2.0 / 15),                                              // x^5
	big.NewFloat(17.0 / 315),                                            // x^7
	big.NewFloat(62.0 / 2835),                                           // x^9
	big.NewFloat(1382.0 / 155925),                                       // x^11
	big.NewFloat(21844.0 / 6081075),                                     // x^13
	big.NewFloat(929569.0 / 638512875),                                  // x^15
	big.NewFloat(6404582.0 / 10854718875),                               // x^17
	big.NewFloat(443861162.0 / 1856156927625),                           // x^19
	big.NewFloat(18888466084.0 / 194896477400625),                       // x^21
	big.NewFloat(113927491862.0 / 2900518163668125),                     // x^23
	big.NewFloat(58870668456604.0 / 3698160658676859375),                // x^25
	big.NewFloat(8374643517010684.0 / 1298054391195577640625),           // x^27
	big.NewFloat(689005380505609448.0 / 263747951750360639921875),       // x^29
	big.NewFloat(129848163681107301953.0 / 122529844256906551386796875), // x^31
}

// tanTaylor calculates tan(x) using direct Taylor series.
// Uses the series: tan(x) = x + x³/3 + 2x⁵/15 + 17x⁷/315 + ...
//
// This is a package-private method for performance comparison.
func tanTaylor(x *big.Float) *big.Float {
	prec := x.Prec()

	// Track transformations needed
	negateResult := false
	reciprocal := false

	// Reduce argument to [-π/4, π/4] for better convergence
	reducedX := new(big.Float).SetPrec(prec).Set(x)

	// Handle negative values: tan(-x) = -tan(x)
	if reducedX.Sign() < 0 {
		negateResult = !negateResult
		reducedX.Abs(reducedX)
	}

	// Reduce argument using tan(x) periodicity and symmetries
	for reducedX.Cmp(bigPi) >= 0 {
		reducedX.Sub(reducedX, bigPi)
	}

	// Use tan(π - x) = -tan(x) for [π/2, π]
	if reducedX.Cmp(bigHalfPi) > 0 {
		reducedX.Sub(bigPi, reducedX)
		negateResult = !negateResult
	}

	// Use tan(π/2 - x) = cot(x) = 1/tan(x) for [π/4, π/2]
	if reducedX.Cmp(bigQuarterPi) > 0 {
		reducedX.Sub(bigHalfPi, reducedX)
		reciprocal = !reciprocal
	}

	// Now reducedX is in [0, π/4]
	result := new(big.Float).SetPrec(prec).Set(reducedX)
	term := new(big.Float).SetPrec(prec).Set(reducedX)
	xSquared := new(big.Float).SetPrec(prec).Mul(reducedX, reducedX)

	convergenceThreshold := new(big.Float).SetPrec(prec).Quo(
		one,
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec/4)), nil)),
	)

	for _, coeff := range taylorTanSeriesCoeffs {
		term.Mul(term, xSquared)
		termValue := new(big.Float).SetPrec(prec).Mul(term, coeff)

		result.Add(result, termValue)

		// Check convergence
		if termValue.Cmp(convergenceThreshold) <= 0 {
			// TODO(rsned): Add an indicator to signal when convergence
			// was not reached and we ran out of iterations.
			break
		}
	}
	// Apply transformations
	if reciprocal {
		result.Quo(one, result)
	}
	if negateResult {
		result.Neg(result)
	}

	return result
}
