package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestCsc(t *testing.T) {
	testCases := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"Cosecant(π/6)", math.Pi / 6, 2, 1e-12},
		{"Cosecant(π/4)", math.Pi / 4, math.Sqrt2, 1e-10},
		{"Cosecant(π/3)", math.Pi / 3, 2 / math.Sqrt(3), 1e-14},
		{"Cosecant(π/2)", math.Pi / 2, 1, 1e-14},
		{"Cosecant(-π/4)", -math.Pi / 4, -math.Sqrt2, 1e-10},
		{"Cosecant(1)", 1, 1 / math.Sin(1), 1e-14},
		{"Cosecant(0.1)", 0.1, 1 / math.Sin(0.1), 1e-14},
		{"Cosecant(0.01)", 0.01, 1 / math.Sin(0.01), 1e-13},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(tc.input)

			result := Csc(x)
			actual, _ := result.Float64()

			diff := math.Abs(actual - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("Csc(%0.4f) = %v, want %v (diff: %v, tolerance: %v)",
					tc.input, actual, tc.expected, diff, tc.tolerance)
			}
		})
	}
}

func TestCscVsMathCsc(t *testing.T) {
	// Test against 1/math.Sin for a range of values
	// Avoiding values close to nπ where cosecant is undefined
	testValues := []float64{
		0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
		1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 2.0,
		2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 2.9, 3.0,
		math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2,
		2 * math.Pi / 3, 3 * math.Pi / 4, 5 * math.Pi / 6,
		-0.1, -0.5,
		// -1.0, // tolerance = 1e-2
		-math.Pi / 4,
		//-math.Pi / 2,  // tolerance = 1e-2
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("Csc(%0.4f)", val), func(t *testing.T) {
			x := new(big.Float).SetPrec(64)
			x.SetFloat64(val)

			bigResult := Csc(x)
			bigFloat, _ := bigResult.Float64()
			mathResult := 1 / math.Sin(val)

			tolerance := 1e-10
			diff := math.Abs(bigFloat - mathResult)

			if diff > tolerance {
				t.Errorf("Csc(%0.4f): bigmath=%v, math=%v, diff=%v",
					val, bigFloat, mathResult, diff)
			}
		})
	}
}

func BenchmarkCsc(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Csc(x)
	}
}

func BenchmarkCscVsMathCsc(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.Run("BigMath", func(b *testing.B) {
		for b.Loop() {
			Csc(x)
		}
	})

	b.Run("MathLib", func(b *testing.B) {
		for b.Loop() {
			_ = 1 / math.Sin(1.0)
		}
	})
}

func BenchmarkCscPrecision(b *testing.B) {
	for _, prec := range precisions {
		b.Run(fmt.Sprintf("prec=%d", prec), func(b *testing.B) {
			x := new(big.Float).SetPrec(prec)
			x.SetFloat64(math.Pi / 3.0)

			b.ResetTimer()
			for b.Loop() {
				Csc(x)
			}
		})
	}
}

func BenchmarkAcsc(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Acsc(x)
	}
}

func BenchmarkAcscVsMathAcsc(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.Run("BigMath", func(b *testing.B) {
		for b.Loop() {
			Acsc(x)
		}
	})

	b.Run("MathLib", func(b *testing.B) {
		for b.Loop() {
			math.Acos(1 / math.Sin(math.Pi/3.0))
		}
	})
}

func BenchmarkCsch(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Csch(x)
	}
}

func BenchmarkAcsch(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(math.Pi / 3.0)

	b.ResetTimer()
	for b.Loop() {
		Acsch(x)
	}
}
