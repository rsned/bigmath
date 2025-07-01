package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestCos(t *testing.T) {
	testCases := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"Cos(0)", 0, 1, 1e-15},
		{"Cos(π/6)", math.Pi / 6, math.Sqrt(3) / 2, 1e-15},
		{"Cos(π/4)", math.Pi / 4, math.Sqrt2 / 2, 1e-15},
		{"Cos(π/3)", math.Pi / 3, 0.5, 1e-15},
		{"Cos(π/2)", math.Pi / 2, 0, 1e-14},
		{"Cos(π)", math.Pi, -1, 1e-15},
		{"Cos(3π/2)", 3 * math.Pi / 2, 0, 1e-13},
		{"Cos(2π)", 2 * math.Pi, 1, 1e-13},
		{"Cos(-π/4)", -math.Pi / 4, math.Sqrt2 / 2, 1e-15},
		{"Cos(1)", 1, math.Cos(1), 1e-15},
		{"Cos(0.1)", 0.1, math.Cos(0.1), 1e-15},
		{"Cos(0.01)", 0.01, math.Cos(0.01), 1e-15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(tc.input)

			result := Cos(x)
			actual, _ := result.Float64()

			diff := math.Abs(actual - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("Cos(%v) = %v, want %v (diff: %v, tolerance: %v)",
					tc.input, actual, tc.expected, diff, tc.tolerance)
			}
		})
	}
}

func TestCosVsMathCos(t *testing.T) {
	// Test against math.Cos for a range of values
	testValues := []float64{
		0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
		1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 2.0,
		math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2,
		2 * math.Pi / 3, 3 * math.Pi / 4, 5 * math.Pi / 6, math.Pi,
		-0.1, -0.5, -1.0, -math.Pi / 4, -math.Pi / 2, -math.Pi,
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("Cos(%0.4f)", val), func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(val)

			bigResult := Cos(x)
			bigFloat, _ := bigResult.Float64()
			mathResult := math.Cos(val)

			tolerance := 1e-14
			diff := math.Abs(bigFloat - mathResult)

			if diff > tolerance {
				t.Errorf("Cos(%0.4f): bigmath=%v, math=%v, diff=%v",
					val, bigFloat, mathResult, diff)
			}
		})
	}
}

func BenchmarkCos(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	for b.Loop() {
		Cos(x)
	}
}

func BenchmarkCosVsMathCos(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.Run("BigMath", func(b *testing.B) {
		for b.Loop() {
			Cos(x)
		}
	})

	b.Run("MathLib", func(b *testing.B) {
		for b.Loop() {
			math.Cos(math.Pi / 3.0)
		}
	})
}

func BenchmarkCosPrecision(b *testing.B) {
	for _, prec := range precisions {
		b.Run(fmt.Sprintf("prec=%d", prec), func(b *testing.B) {
			x := new(big.Float).SetPrec(prec)
			x.SetFloat64(math.Pi / 3.0)

			b.ResetTimer()
			for b.Loop() {
				Cos(x)
			}
		})
	}
}

func BenchmarkAcos(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Acos(x)
	}
}

func BenchmarkAcosVsMathAcos(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.Run("BigMath", func(b *testing.B) {
		for b.Loop() {
			Acos(x)
		}
	})

	b.Run("MathLib", func(b *testing.B) {
		for b.Loop() {
			math.Acos(math.Pi / 3.0)
		}
	})
}

func BenchmarkCosh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Cosh(x)
	}
}

func BenchmarkAcosh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Acosh(x)
	}
}

// Basic timing benchmark.
func BenchmarkCosInternalFunctions(b *testing.B) {
	runTrigBenchmark(b, cosMethods, precisions)
}
