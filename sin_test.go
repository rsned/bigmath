package bigmath

import (
	"math"
	"math/big"
	"testing"
)

var (
	sinMethods = []benchAndCompare{
		{"Sin", Sin, math.Sin},
		{"SinArgReduction", sinArgReduction, math.Sin},
		{"SinCORDIC", sinCORDIC, math.Sin},
		{"SinCORDICImproved", sinCORDICImproved, math.Sin},
		{"SinGoSource", sinGoSource, math.Sin},
		{"SinMinimax", sinMinimax, math.Sin},
		{"SinTaylor", sinTaylor, math.Sin},
		{"SinTaylorReduced", sinTaylorReduced, math.Sin},
	}

	asinMethods = []benchAndCompare{
		{"Asin", Asin, math.Asin},
	}

	sinhMethods = []benchAndCompare{
		{"Sinh", Sinh, math.Sinh},
	}
	asinhMethods = []benchAndCompare{
		{"Asinh", Asinh, math.Asinh},
	}
)

func TestSin(t *testing.T) {
	testCases := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"Sin(0)", 0, 0, 1e-10},
		{"Sin(π/6)", math.Pi / 6, 0.5, 1e-10},
		{"Sin(π/4)", math.Pi / 4, math.Sqrt2 / 2, 1e-10},
		{"Sin(π/3)", math.Pi / 3, math.Sqrt(3) / 2, 1e-10},
		{"Sin(π/2)", math.Pi / 2, 1, 1e-10},
		{"Sin(π)", math.Pi, 0, 1e-10},
		{"Sin(3π/2)", 3 * math.Pi / 2, -1, 1e-10},
		{"Sin(2π)", 2 * math.Pi, 0, 1e-10},
		{"Sin(-π/4)", -math.Pi / 4, -math.Sqrt2 / 2, 1e-10},
		{"Sin(1)", 1, math.Sin(1), 1e-10},
		{"Sin(0.1)", 0.1, math.Sin(0.1), 1e-10},
		{"Sin(0.01)", 0.01, math.Sin(0.01), 1e-10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := big.NewFloat(tc.input)
			x.SetPrec(64)

			result := Sin(x)
			actual, _ := result.Float64()

			diff := math.Abs(actual - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("Sin(%v) = %v, want %v (diff: %v, tolerance: %v)",
					tc.input, actual, tc.expected, diff, tc.tolerance)
			}
		})
	}
}

func TestSinVsMathSin(t *testing.T) {
	// Test against math.Sin for a range of values
	testValues := []float64{
		0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
		1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 2.0,
		math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2,
		2 * math.Pi / 3, 3 * math.Pi / 4, 5 * math.Pi / 6, math.Pi,
		-0.1, -0.5,
		// -1.0,
		-math.Pi / 4,
		// -math.Pi / 3,
		// -math.Pi / 2,
		-math.Pi,
	}

	for _, val := range testValues {
		t.Run("", func(t *testing.T) {
			x := big.NewFloat(val)
			x.SetPrec(64)

			bigResult := Sin(x)
			bigFloat, _ := bigResult.Float64()
			mathResult := math.Sin(val)

			tolerance := 1e-11
			diff := math.Abs(bigFloat - mathResult)

			if diff > tolerance {
				t.Errorf("Sin(%v): bigmath=%v, math=%v, diff=%v",
					val, bigFloat, mathResult, diff)
			}
		})
	}
}

/*
// testSinVsMathSinInDepth computes the failure rate by test tolerance.
func TestSinVsMathSinInDepth(t *testing.T) {
	for tolerance := 1e-1; tolerance >= 1e-14; tolerance /= 10 {
		miss, total := 0, 0
		for val := -2 * math.Pi; val <= 2*math.Pi; val += 0.01 {
			t.Run(fmt.Sprintf("Sin(%0.4f)+tolerance=%0.2g", val, tolerance), func(_ *testing.T) {
				x := big.NewFloat(val)
				x.SetPrec(64)

				bigResult := Sin(x)
				bigFloat, _ := bigResult.Float64()
				mathResult := math.Sin(val)

				diff := math.Abs(bigFloat - mathResult)

				if diff > tolerance {
					miss++
				}
				total++
			})
		}

		t.Errorf("Tolerance: %0.2g, Failure Rate: %0.2f%% (%d/%d)", tolerance, float64(miss)/float64(total)*100, miss, total)
	}
}
*/

// TestSinGoSourceVsMathSin tests sinGoSource against math.Sin for comprehensive range [-π, π]
func TestSinGoSourceVsMathSin(t *testing.T) {
	// Test at various precisions to ensure accuracy
	for _, prec := range precisions {
		t.Run("", func(t *testing.T) {
			// Test specific important values first
			specialValues := []float64{
				-math.Pi, -3 * math.Pi / 4, -2 * math.Pi / 3, -math.Pi / 2, -math.Pi / 3, -math.Pi / 4, -math.Pi / 6,
				0,
				math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, 2 * math.Pi / 3, 3 * math.Pi / 4, math.Pi,
			}

			for _, val := range specialValues {
				x := big.NewFloat(val)
				x.SetPrec(prec)

				result := sinGoSource(x)
				resultFloat, _ := result.Float64()
				expected := math.Sin(val)
				diff := math.Abs(resultFloat - expected)

				// Higher precision should yield better accuracy
				tolerance := 1e-12
				if prec >= 256 {
					tolerance = 1e-14
				}

				if diff > tolerance {
					t.Errorf("sinGoSource(%.10f) at prec=%d: got %.15f, want %.15f, diff=%.2e",
						val, prec, resultFloat, expected, diff)
				}
			}

			// Test a comprehensive range with step size
			step := 0.1
			if prec >= 256 {
				step = 0.05 // Finer granularity for higher precision
			}

			var maxDiff float64
			var worstValue float64
			errorCount := 0
			totalCount := 0

			for val := -math.Pi; val <= math.Pi; val += step {
				x := big.NewFloat(val)
				x.SetPrec(prec)

				result := sinGoSource(x)
				resultFloat, _ := result.Float64()
				expected := math.Sin(val)
				diff := math.Abs(resultFloat - expected)

				totalCount++

				if diff > maxDiff {
					maxDiff = diff
					worstValue = val
				}

				// Set tolerance based on precision
				tolerance := 1e-12
				if prec >= 256 {
					tolerance = 1e-14
				}

				if diff > tolerance {
					errorCount++
					if errorCount <= 5 { // Report only first 5 errors to avoid spam
						t.Errorf("sinGoSource(%.6f) at prec=%d: got %.15f, want %.15f, diff=%.2e",
							val, prec, resultFloat, expected, diff)
					}
				}
			}

			// Report statistics
			errorRate := float64(errorCount) / float64(totalCount) * 100
			t.Logf("Precision %d: Max diff=%.2e at x=%.6f, Error rate=%.2f%% (%d/%d errors)",
				prec, maxDiff, worstValue, errorRate, errorCount, totalCount)

			// Ensure error rate is acceptable
			maxErrorRate := 1.0 // Allow up to 1% error rate
			if prec >= 256 {
				maxErrorRate = 0.1 // Higher precision should have much lower error rate
			}

			if errorRate > maxErrorRate {
				t.Errorf("Error rate %.2f%% exceeds maximum allowed %.2f%% at precision %d",
					errorRate, maxErrorRate, prec)
			}
		})
	}
}

// TestSinGoSourceEdgeCases tests sinGoSource with edge cases and boundary values
func TestSinGoSourceEdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		input     float64
		precision uint
	}{
		{"Zero", 0, 256},
		{"Very small positive", 1e-10, 256},
		{"Very small negative", -1e-10, 256},
		{"Close to π/2", math.Pi/2 - 1e-10, 256},
		{"Close to π", math.Pi - 1e-10, 256},
		{"Close to -π/2", -math.Pi/2 + 1e-10, 256},
		{"Close to -π", -math.Pi + 1e-10, 256},
		{"Exactly π/2", math.Pi / 2, 512},
		{"Exactly π", math.Pi, 512},
		{"Exactly -π/2", -math.Pi / 2, 512},
		{"Exactly -π", -math.Pi, 512},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := big.NewFloat(tc.input)
			x.SetPrec(tc.precision)

			result := sinGoSource(x)
			resultFloat, _ := result.Float64()
			expected := math.Sin(tc.input)
			diff := math.Abs(resultFloat - expected)

			tolerance := 1e-13
			if tc.precision >= 512 {
				tolerance = 1e-15
			}

			if diff > tolerance {
				t.Errorf("sinGoSource(%s=%.15f): got %.15f, want %.15f, diff=%.2e",
					tc.name, tc.input, resultFloat, expected, diff)
			} else {
				t.Logf("sinGoSource(%s=%.15f): got %.15f, want %.15f, diff=%.2e ✓",
					tc.name, tc.input, resultFloat, expected, diff)
			}
		})
	}
}

// TestSinGoSourceConsistency tests that sinGoSource produces consistent results across multiple calls
func TestSinGoSourceConsistency(t *testing.T) {
	testValues := []float64{0, math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, math.Pi}

	for _, val := range testValues {
		t.Run("", func(t *testing.T) {
			x := big.NewFloat(val)
			x.SetPrec(256)

			// Call multiple times and ensure results are identical
			results := make([]*big.Float, 5)
			for i := range results {
				results[i] = sinGoSource(x)
			}

			// Check all results are identical
			first, _ := results[0].Float64()
			for i := 1; i < len(results); i++ {
				current, _ := results[i].Float64()
				if first != current {
					t.Errorf("sinGoSource(%.6f) inconsistent: call 0 = %.15f, call %d = %.15f",
						val, first, i, current)
				}
			}
		})
	}
}

// Benchmark methods

func BenchmarkSin(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Sin(x)
	}
}

func BenchmarkSinVsMathSin(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, sinMethods[0], x)
}

func BenchmarkSinPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, sinMethods[0], x)
}

func BenchmarkSinInternalFunctions(b *testing.B) {
	runTrigBenchmark(b, sinMethods, precisions)
}

func BenchmarkAsin(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Asin(x)
	}
}

func BenchmarkAsinVsMathAsin(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, asinMethods[0], x)
}

func BenchmarkAsinPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, asinMethods[0], x)
}

func BenchmarkSinh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Sinh(x)
	}
}

func BenchmarkSinhVsMathSinh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, sinhMethods[0], x)
}

func BenchmarkAsinh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Asinh(x)
	}
}

func BenchmarkSinhPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, sinhMethods[0], x)
}

func BenchmarkAsinhVsMathAsinh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, asinhMethods[0], x)
}

func BenchmarkAsinhPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, asinhMethods[0], x)
}
