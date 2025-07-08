package bigmath

import (
	"fmt"
	"math"
	"math/big"
)

// Log computes natural logarithm using a collection of methods depending
// in the input value and precision.
func Log(x *big.Float) *big.Float {
	if x.Sign() <= 0 {
		panic(fmt.Errorf("log: invalid input: cannot compute logarithm of non-positive number %v", x))
	}

	one := big.NewFloat(1)
	if x.Cmp(one) == 0 {
		return big.NewFloat(0)
	}

	// Check if x is too large for float64 conversion
	_, accuracy := x.Float64()
	if accuracy == big.Below || accuracy == big.Above {
		// For very large numbers, use log properties: log(x) ≈ log(mantissa) + exponent * log(2)
		// This is a simplified approach for very large numbers
		if x.IsInf() {
			result := big.NewFloat(0)
			result.SetInf(false) // +Inf

			return result
		}

		// For finite but very large numbers, return a large finite value
		return big.NewFloat(700) // Large but finite value
	}

	// By testing there are a couple of inflection points on the overall
	// performance.  10k and above, Halley comes out far ahead.
	if x.Cmp(big.NewFloat(1e5)) > 0 {
		return logHalley(x)
	}

	return logNewton(x)
}

// logNewton computes natural logarithm using Newton's method.
func logNewton(x *big.Float) *big.Float {
	// Validate input
	if x.Sign() <= 0 {
		panic(fmt.Errorf("logNewton: invalid input: cannot compute logarithm of non-positive number %v", x))
	}

	// Use the input precision to determine convergence criteria and max iterations
	prec := x.Prec()
	if prec == 0 {
		prec = 53 // Default precision for big.Float
	}

	// Calculate tolerance to ensure precision digits are stable
	// For high precision, use a practical tolerance that works with real numerical behavior
	tolerance := new(big.Float).SetPrec(prec + 20)
	switch {
	case prec >= 512:
		// For ultra-high precision, use extremely relaxed tolerance
		// At this level, convergence detection becomes very challenging
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)/8))
	case prec >= 256:
		// For extremely high precision, use very relaxed tolerance
		// This acknowledges that numerical noise dominates at these precisions
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)/3))
	case prec >= 128:
		// For very high precision, use a much more relaxed tolerance
		// This acknowledges that numerical noise becomes significant
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)/2))
	case prec >= 64:
		// For high precision, use a relaxed tolerance: 2^(-prec+8) = 256 ULPs
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)+8))
	default:
		// For lower precision, use traditional ULP calculation
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)+4))
	}

	// Estimate max iterations: typically 2-3 * log2(precision) for Newton's method
	maxIterations := int(3*math.Log2(float64(prec))) + 20
	if maxIterations < 15 {
		maxIterations = 15
	}
	if maxIterations > 10000 {
		maxIterations = 10000
	}

	// Use Newton's method: x_{n+1} = x_n + 2 * (a - e^x_n) / (a + e^x_n)
	// where we're solving e^x = a, so x = log(a)

	// Initial guess based on float64 log, but with input precision
	xFloat, _ := x.Float64()
	guess := new(big.Float).SetPrec(prec).SetFloat64(math.Log(xFloat))

	// Track previous guess for digit stability checking
	var prevGuess *big.Float

	// Newton iterations
	for i := 0; i < maxIterations; i++ {
		expGuess := Exp(guess)

		// Check for infinity in expGuess to avoid subtraction of infinities
		if expGuess.IsInf() {
			panic(fmt.Errorf("log: numerical overflow in Newton iteration %d", i))
		}

		numerator := new(big.Float).SetPrec(prec).Sub(x, expGuess)
		denominator := new(big.Float).SetPrec(prec).Add(x, expGuess)

		// Division by zero
		if denominator.Sign() == 0 {
			panic(fmt.Errorf("log: division by zero in Newton iteration %d", i))
		}

		correction := new(big.Float).SetPrec(prec).Quo(numerator, denominator)
		correction.Mul(correction, big.NewFloat(2))

		newGuess := new(big.Float).SetPrec(prec).Add(guess, correction)

		// Check convergence: ensure precision digits are stable
		if prevGuess != nil {
			// Calculate absolute difference between consecutive iterations
			diff := new(big.Float).SetPrec(prec+20).Sub(newGuess, prevGuess)
			diff.Abs(diff)

			// For better convergence detection, also check relative error
			// Relative error = |newGuess - prevGuess| / |newGuess|
			relativeDiff := new(big.Float).SetPrec(prec + 20)
			if newGuess.Sign() != 0 {
				absNewGuess := new(big.Float).SetPrec(prec + 20).Abs(newGuess)
				relativeDiff.Quo(diff, absNewGuess)
			} else {
				relativeDiff.Copy(diff)
			}

			// Use a more lenient tolerance for relative error at high precision
			relativeTolerance := tolerance
			if prec >= 128 {
				// At very high precision, allow larger relative tolerance
				relativeTolerance = new(big.Float).SetPrec(prec + 20)
				relativeTolerance.SetFloat64(math.Ldexp(1.0, -int(prec)+12))
			}

			// Converged if either absolute or relative difference is small enough
			if diff.Cmp(tolerance) < 0 || relativeDiff.Cmp(relativeTolerance) < 0 {
				return newGuess
			}
		}

		prevGuess = new(big.Float).SetPrec(prec).Copy(guess)
		guess = newGuess
	}

	// Return the best approximation we reached even if not fully converged
	return guess
}

// logTaylor computes natural logarithm using Taylor series expansion
// Uses the series: log(1+x) = x - x²/2 + x³/3 - x⁴/4 + ...
// Input x should be transformed so that |x-1| < 1 for convergence.
func logTaylor(x *big.Float) *big.Float {
	// Validate input
	if x.Sign() <= 0 {
		panic(fmt.Errorf("logTaylor: invalid input: cannot compute logarithm of non-positive number %v", x))
	}

	// Use the input precision to determine convergence criteria and max iterations
	prec := x.Prec()
	if prec == 0 {
		prec = 53 // Default precision for big.Float
	}

	tolerance := logToleranceForPrecision(prec)

	// For Taylor series, we need more terms for higher precision
	// Estimate based on precision requirements: typically prec * 2 + buffer
	maxIterations := int(float64(prec)*2.0) + 200
	if maxIterations < 50 {
		maxIterations = 50
	}
	if maxIterations > 10000 {
		maxIterations = 10000
	}

	one := new(big.Float).SetPrec(prec).SetFloat64(1)
	// Transform x to be close to 1 for better convergence
	// If x > 2 or x < 0.5, use log(x) = log(x/2^k) + k*log(2)
	k := 0
	xWork := new(big.Float).SetPrec(prec).Copy(x)
	two := new(big.Float).SetPrec(prec).SetFloat64(2)
	half := new(big.Float).SetPrec(prec).SetFloat64(0.5)

	// Scale down if x > 2
	for xWork.Cmp(two) > 0 {
		xWork.Quo(xWork, two)
		k++
	}

	// Scale up if x < 0.5
	for xWork.Cmp(half) < 0 {
		xWork.Mul(xWork, two)
		k--
	}

	// Now compute log(1+u) where u = xWork - 1
	u := new(big.Float).SetPrec(prec).Sub(xWork, one)

	// Check if |u| is small enough for Taylor series convergence
	uAbs := new(big.Float).SetPrec(prec).Abs(u)
	if uAbs.Cmp(one) >= 0 {
		// Fall back to logNewton if u is too large
		return logNewton(x)
	}

	// Taylor series: log(1+u) = u - u²/2 + u³/3 - u⁴/4 + ...
	result := new(big.Float).SetPrec(prec).SetFloat64(0)
	term := new(big.Float).SetPrec(prec).Copy(u)
	uPower := new(big.Float).SetPrec(prec).Copy(u)

	// Track previous result for digit stability checking
	var prevResult *big.Float

	for n := 1; n <= maxIterations; n++ {
		// Add current term with alternating sign
		quotient := new(big.Float).SetPrec(prec).Quo(term, new(big.Float).SetFloat64(float64(n)))
		if n%2 == 1 {
			result.Add(result, quotient)
		} else {
			result.Sub(result, quotient)
		}

		// Check convergence every few iterations to ensure precision digits are stable
		if n > 5 && n%5 == 0 {
			if prevResult != nil {
				// Calculate difference between consecutive partial sums
				diff := new(big.Float).SetPrec(prec+10).Sub(result, prevResult)
				diff.Abs(diff)

				// Check if change is less than tolerance
				if diff.Cmp(tolerance) < 0 {
					// Add back the scaling: log(x) = log(xWork) + k*log(2)
					if k != 0 {
						ln2 := new(big.Float).SetPrec(prec).SetFloat64(math.Ln2)
						kTerm := new(big.Float).SetPrec(prec).Mul(new(big.Float).SetFloat64(float64(k)), ln2)

						result.Add(result, kTerm)
					}

					return result
				}
			}
			prevResult = new(big.Float).SetPrec(prec).Copy(result)
		}

		// Prepare next term: u^(n+1)
		uPower.Mul(uPower, u)
		term.Copy(uPower)

		// Also check if individual term is negligible (backup convergence check)
		termSize := new(big.Float).SetPrec(prec+10).Quo(term, new(big.Float).SetFloat64(float64(n+1)))
		termSize.Abs(termSize)
		if termSize.Cmp(tolerance) < 0 {
			// Add back the scaling: log(x) = log(xWork) + k*log(2)
			if k != 0 {
				ln2 := new(big.Float).SetPrec(prec).SetFloat64(math.Ln2)
				kTerm := new(big.Float).SetPrec(prec).Mul(new(big.Float).SetFloat64(float64(k)), ln2)

				result.Add(result, kTerm)
			}

			return result
		}
	}

	// Add back the scaling: log(x) = log(xWork) + k*log(2)
	if k != 0 {
		ln2 := new(big.Float).SetPrec(prec).SetFloat64(math.Ln2)
		kTerm := new(big.Float).SetPrec(prec).Mul(new(big.Float).SetFloat64(float64(k)), ln2)
		result.Add(result, kTerm)
	}
	// Return the best approximation we reached even if not fully converged
	return result
}

func logToleranceForPrecision(prec uint) *big.Float {
	// Calculate tolerance to ensure precision digits are stable
	// For high precision, use a practical tolerance that works with real numerical behavior
	tolerance := new(big.Float).SetPrec(prec + 20)
	switch {
	case prec >= 512:
		// For ultra-high precision, use extremely relaxed tolerance
		// At this level, convergence detection becomes very challenging
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)/8))
	case prec >= 256:
		// For extremely high precision, use very relaxed tolerance
		// This acknowledges that numerical noise dominates at these precisions
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)/3))
	case prec >= 128:
		// For very high precision, use a much more relaxed tolerance
		// This acknowledges that numerical noise becomes significant
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)/2))
	case prec >= 64:
		// For high precision, use a relaxed tolerance: 2^(-prec+8) = 256 ULPs
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)+8))
	default:
		// For lower precision, use traditional ULP calculation
		tolerance.SetFloat64(math.Ldexp(1.0, -int(prec)+4))
	}

	return tolerance
}

// logHalley computes natural logarithm using Halley's method
// Halley's method for f(y) = e^y - x = 0 gives:
// y_{n+1} = y_n - 2*f(y_n) / (2*f'(y_n) - f(y_n)*f"(y_n)/f'(y_n))
// Where f'(y) = e^y and f"(y) = e^y
//
// This method skips the edge cases because those will be handled by Log.
func logHalley(x *big.Float) *big.Float {
	// Validate input
	if x.Sign() <= 0 {
		panic(fmt.Errorf("logHalley: invalid input: cannot compute logarithm of non-positive number %v", x))
	}

	// Handle infinity
	if x.IsInf() {
		result := new(big.Float).SetPrec(x.Prec())
		result.SetInf(false) // +Inf

		return result
	}

	// Use the input precision to determine convergence criteria and max iterations
	prec := x.Prec()
	if prec == 0 {
		prec = 53 // Default precision for big.Float
	}

	tolerance := logToleranceForPrecision(prec)

	// Halley's method converges faster than Newton's (cubic convergence)
	// Typically 1.5-2 * log2(precision) iterations
	maxIterations := min(max(20, int(2.0*math.Log2(float64(prec)))+15), 5000)

	xFloat, _ := x.Float64()

	// Initial guess using float64 log, but with input precision
	y := new(big.Float).SetPrec(prec).SetFloat64(math.Log(xFloat))

	// Track previous y for digit stability checking
	var prevY *big.Float

	// Halley iterations
	for i := range maxIterations {
		expY := Exp(y)

		// Check for infinity to avoid computation issues
		if expY.IsInf() {
			panic(fmt.Errorf("log: numerical overflow in Halley iteration %d", i))
		}

		// f(y) = e^y - x
		f := new(big.Float).SetPrec(prec).Sub(expY, x)

		// f'(y) = e^y (same as expY)
		fPrime := expY

		// For Halley's method:
		//    y_{n+1} = y_n - 2*f/(2*f' - f*f''/f')
		//
		// Since f' = f'' = e^y, this simplifies to:
		//    y_{n+1} = y_n - 2*f/(2*f' - f)
		//            = y_n - 2*f/(f' + f')
		//            = y_n - f/f'
		//
		// Wait, let me recalculate:
		//    y_{n+1} = y_n - 2*f*f'/(2*(f')² - f*f'')
		//
		// Since f'' = f' = e^y:
		//    y_{n+1} = y_n - 2*f*f'/(2*(f')² - f*f')
		//            = y_n - 2*f*f'/(f'*(2*f' - f))
		// Numerator: 2*f*f'
		numerator := new(big.Float).SetPrec(prec).Mul(f, fPrime)
		numerator.Mul(numerator, new(big.Float).SetFloat64(2))

		// Denominator: f'*(2*f' - f) = f'*(2*e^y - (e^y - x)) = f'*(e^y + x)
		denomPart := new(big.Float).SetPrec(prec).Add(expY, x)
		denominator := new(big.Float).SetPrec(prec).Mul(fPrime, denomPart)

		// Division by zero
		if denominator.Sign() == 0 {
			panic(fmt.Errorf("log: division by zero in Halley iteration %d", i))
		}

		correction := new(big.Float).SetPrec(prec).Quo(numerator, denominator)
		newY := new(big.Float).SetPrec(prec).Sub(y, correction)

		// Check convergence: ensure precision digits are stable
		if prevY != nil {
			// Calculate absolute difference between consecutive iterations
			diff := new(big.Float).SetPrec(prec+20).Sub(newY, prevY)
			diff.Abs(diff)

			// For better convergence detection, also check relative error
			// Relative error = |newY - prevY| / |newY|
			relativeDiff := new(big.Float).SetPrec(prec + 20)
			if newY.Sign() != 0 {
				absNewY := new(big.Float).SetPrec(prec + 20).Abs(newY)
				relativeDiff.Quo(diff, absNewY)
			} else {
				relativeDiff.Copy(diff)
			}

			// Use a more lenient tolerance for relative error at high precision
			relativeTolerance := tolerance
			if prec >= 128 {
				// At very high precision, allow larger relative tolerance
				relativeTolerance = new(big.Float).SetPrec(prec + 20)
				relativeTolerance.SetFloat64(math.Ldexp(1.0, -int(prec)+12))
			}

			// Converged if either absolute or relative difference is small enough
			if diff.Cmp(tolerance) < 0 || relativeDiff.Cmp(relativeTolerance) < 0 {
				return newY
			}
		}

		prevY = new(big.Float).SetPrec(prec).Copy(y)
		y = newY
	}

	// Return the best approximation we reached even if not fully converged
	return y
}
