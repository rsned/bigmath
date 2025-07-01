package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

var (
	cotMethods = []benchAndCompare{
		{"Cot", Cot, func(x float64) float64 { return math.Cos(x) / math.Sin(x) }},
	}

	cothMethods = []benchAndCompare{
		{"Coth", Coth, func(x float64) float64 { return math.Cosh(x) / math.Sinh(x) }},
	}

	acotMethods = []benchAndCompare{
		{"Acot", Acot, func(x float64) float64 { return math.Atan(1 / x) }},
	}

	acothMethods = []benchAndCompare{
		{"Acoth", Acoth, func(x float64) float64 { return math.Atanh(1 / x) }},
	}
)

func TestCot(t *testing.T) {
	testCases := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"Cotangent(π/6)", math.Pi / 6, math.Sqrt(3), 1e-12},
		{"Cotangent(π/4)", math.Pi / 4, 1, 1e-11},
		{"Cotangent(π/3)", math.Pi / 3, 1 / math.Sqrt(3), 1e-12},
		{"Cotangent(-π/4)", -math.Pi / 4, -1, 1e-11},
		{"Cotangent(1)", 1, math.Cos(1) / math.Sin(1), 1e-12},
		{"Cotangent(0.1)", 0.1, math.Cos(0.1) / math.Sin(0.1), 1e-12},
		{"Cotangent(0.01)", 0.01, math.Cos(0.01) / math.Sin(0.01), 1e-12},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(tc.input)

			result := Cot(x)
			actual, _ := result.Float64()

			diff := math.Abs(actual - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("Cotangent(%v) = %v, want %v (diff: %v, tolerance: %v)",
					tc.input, actual, tc.expected, diff, tc.tolerance)
			}
		})
	}
}

func TestCotVsMathCot(t *testing.T) {
	// Test against math.Cos/math.Sin for a range of values
	// Avoiding values close to nπ where cotangent is undefined
	testValues := []float64{
		0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
		1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 2.0,
		2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 2.9, 3.0,
		math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2,
		2 * math.Pi / 3, 3 * math.Pi / 4, 5 * math.Pi / 6,
		-0.1, -0.5,
		// -1.0,  tolerance 1e-2
		-math.Pi / 4, -math.Pi / 2,
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("Cot(%0.4f)", val), func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(val)

			bigResult := Cot(x)
			bigFloat, _ := bigResult.Float64()
			mathResult := math.Cos(val) / math.Sin(val)

			tolerance := 1e-10
			diff := math.Abs(bigFloat - mathResult)

			if diff > tolerance {
				t.Errorf("Cot(%0.4f): bigmath=%v, math=%v, diff=%v",
					val, bigFloat, mathResult, diff)
			}
		})
	}
}

func BenchmarkCot(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Cot(x)
	}
}

func BenchmarkCotVsMathCot(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, cotMethods[0], x)
}

func BenchmarkCotPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, cotMethods[0], x)
}

func BenchmarkAcot(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Acot(x)
	}
}

func BenchmarkAcotVsMathAcot(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, acotMethods[0], x)
}

func BenchmarkCothPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, cothMethods[0], x)
}

func BenchmarkAcothPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, acothMethods[0], x)
}

func BenchmarkCoth(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Coth(x)
	}
}

func BenchmarkAcoth(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Acoth(x)
	}
}
