package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

var (
	secMethods = []benchAndCompare{
		{"Sec", Sec, func(x float64) float64 { return 1 / math.Cos(x) }},
		{"SecSeries", secSeries, func(x float64) float64 { return 1 / math.Cos(x) }},
	}

	sechMethods = []benchAndCompare{
		{"Sech", Sech, func(x float64) float64 { return 1 / math.Cosh(x) }},
	}

	asecMethods = []benchAndCompare{
		{"Asec", Asec, func(x float64) float64 { return math.Acos(1 / x) }},
	}

	asechMethods = []benchAndCompare{
		{"Asech", Asech, func(x float64) float64 { return math.Acosh(1 / x) }},
	}
)

func TestSecant(t *testing.T) {
	testCases := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"Secant(0)", 0, 1, 1e-15},
		{"Secant(π/6)", math.Pi / 6, 2 / math.Sqrt(3), 1e-14},
		{"Secant(π/4)", math.Pi / 4, math.Sqrt2, 1e-14},
		{"Secant(π/3)", math.Pi / 3, 2, 1e-14},
		{"Secant(π)", math.Pi, -1, 1e-14},
		{"Secant(-π/4)", -math.Pi / 4, math.Sqrt2, 1e-14},
		{"Secant(1)", 1, 1 / math.Cos(1), 1e-14},
		{"Secant(0.1)", 0.1, 1 / math.Cos(0.1), 1e-15},
		{"Secant(0.01)", 0.01, 1 / math.Cos(0.01), 1e-15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(tc.input)

			result := Sec(x)
			actual, _ := result.Float64()

			diff := math.Abs(actual - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("Secant(%v) = %v, want %v (diff: %v, tolerance: %v)",
					tc.input, actual, tc.expected, diff, tc.tolerance)
			}
		})
	}
}

func TestSecantVsMathSecant(t *testing.T) {
	// Test against 1/math.Cos for a range of values
	// Avoiding values close to π/2 + nπ where secant is undefined
	testValues := []float64{
		0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
		1.1, 1.2, 1.3, 1.4, // Skip 1.57... (π/2)
		math.Pi / 6, math.Pi / 4, math.Pi / 3,
		2 * math.Pi / 3, 3 * math.Pi / 4, 5 * math.Pi / 6, math.Pi,
		-0.1, -0.5, -1.0, -math.Pi / 4, -math.Pi / 3, -math.Pi,
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("Secant(%0.4f)", val), func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(val)

			bigResult := Sec(x)
			bigFloat, _ := bigResult.Float64()
			mathResult := 1 / math.Cos(val)

			tolerance := 1e-13
			diff := math.Abs(bigFloat - mathResult)

			if diff > tolerance {
				t.Errorf("Secant(%v): bigmath=%v, math=%v, diff=%v",
					val, bigFloat, mathResult, diff)
			}
		})
	}
}

func BenchmarkSec(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Sec(x)
	}
}

func BenchmarkSecVsMathSec(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, secMethods[0], x)
}

func BenchmarkSecPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, secMethods[0], x)
}

func BenchmarkSecInternalFunctions(b *testing.B) {
	runTrigBenchmark(b, secMethods, precisions)
}

func BenchmarkAsec(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Asec(x)
	}
}

func BenchmarkAsecVsMathAsec(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathVsStdlib(b, asecMethods[0], x)
}

func BenchmarkAsecPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, asecMethods[0], x)
}

func BenchmarkSech(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Sech(x)
	}
}

func BenchmarkSechPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, sechMethods[0], x)
}

func BenchmarkAsech(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Asech(x)
	}
}

func BenchmarkAsechPrecision(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision)
	x.SetFloat64(math.Pi / 3.0)

	benchmarkBigmathFunctionVsPrecision(b, asechMethods[0], x)
}
