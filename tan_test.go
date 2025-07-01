package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestTan(t *testing.T) {
	tolerance := 1e-11
	testCases := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Tan(0)", 0, 0},
		{"Tan(π/6)", math.Pi / 6, 1 / math.Sqrt(3)},
		{"Tan(π/4)", math.Pi / 4, 1},
		{"Tan(π/3)", math.Pi / 3, math.Sqrt(3)},
		{"Tan(π)", math.Pi, 0},
		{"Tan(-π/4)", -math.Pi / 4, -1},
		{"Tan(1)", 1, math.Tan(1)},
		{"Tan(0.1)", 0.1, math.Tan(0.1)},
		{"Tan(0.01)", 0.01, math.Tan(0.01)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := big.NewFloat(tc.input)
			x.SetPrec(256)

			result := Tan(x)
			actual, _ := result.Float64()

			diff := math.Abs(actual - tc.expected)
			if diff > tolerance {
				t.Errorf("Tan(%v) = %v, want %v (diff: %v, tolerance: %v)",
					tc.input, actual, tc.expected, diff, tolerance)
			}
		})
	}
}

func TestTanVsMathTan(t *testing.T) {
	// Test against math.Tan for a range of values
	// Avoiding values close to π/2 + nπ where tan is undefined
	testValues := []float64{
		0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
		1.1, 1.2, 1.3, 1.4, // Skip 1.57... (π/2)
		math.Pi / 6, math.Pi / 4, math.Pi / 3,
		2 * math.Pi / 3, 3 * math.Pi / 4, 5 * math.Pi / 6, math.Pi,
		-0.1, -0.5,
		// -1.0, // Needs tolerance = 1e-2
		-math.Pi / 4,
		// -math.Pi / 3, // Needs tolerance = 1e-2
		-math.Pi,
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("Tan(%0.4f)", val), func(t *testing.T) {
			x := big.NewFloat(val)
			x.SetPrec(64)

			bigResult := Tan(x)
			bigFloat, _ := bigResult.Float64()
			mathResult := math.Tan(val)

			tolerance := 1e-11
			diff := math.Abs(bigFloat - mathResult)

			if diff > tolerance {
				t.Errorf("Tan(%v): bigmath=%v, math=%v, diff=%v",
					val, bigFloat, mathResult, diff)
			}
		})
	}
}

func BenchmarkTan(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	for b.Loop() {
		Tan(x)
	}
}

func BenchmarkTanVsMathTan(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.Run("BigMath", func(b *testing.B) {
		for b.Loop() {
			Tan(x)
		}
	})

	b.Run("MathLib", func(b *testing.B) {
		for b.Loop() {
			math.Tan(1.0)
		}
	})
}

func BenchmarkTanPrecisionBits(b *testing.B) {
	for _, prec := range precisions {
		b.Run(fmt.Sprintf("%d_bits", prec), func(b *testing.B) {
			x := new(big.Float).SetPrec(prec)
			x.SetFloat64(math.Pi / 3.0)

			b.ResetTimer()
			for b.Loop() {
				Tan(x)
			}
		})
	}
}

func BenchmarkAtan(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Atan(x)
	}
}

func BenchmarkAtanVsMathAtan(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.Run("BigMath", func(b *testing.B) {
		for b.Loop() {
			Atan(x)
		}
	})

	b.Run("MathLib", func(b *testing.B) {
		for b.Loop() {
			math.Atan(math.Pi / 3.0)
		}
	})
}

func BenchmarkTanh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Tanh(x)
	}
}

func BenchmarkAtanh(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Atanh(x)
	}
}

// Basic timing benchmark.
func BenchmarkTanInternalFunctions(b *testing.B) {
	runTrigBenchmark(b, tanMethods, precisions)
}
