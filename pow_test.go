package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

var (
	powMethods = []benchAndCompare2Args{
		{name: "Pow", fnBigmath: Pow, fnStdlib: math.Pow},
		{name: "PowFloat64", fnBigmath: func(x, y *big.Float) *big.Float {
			fX, _ := x.Float64()
			fY, _ := y.Float64()

			return PowFloat64(fX, fY)
		}, fnStdlib: math.Pow},
	}
)

func TestPow(t *testing.T) {
	tests := []struct {
		x         float64
		y         float64
		expected  float64
		tolerance float64
	}{
		{2, 3, 8, 1e-10},
		{3, 2, 9, 1e-10},
		{10, 0, 1, 1e-10},
		{0, 5, 0, 1e-10},
		{1, 100, 1, 1e-10},
		{2, 0.5, math.Sqrt(2), 1e-10},
		{4, 0.5, 2, 1e-10},
		{math.E, 1, math.E, 1e-10},
		{2, 10, 1024, 1e-10},
	}

	for _, test := range tests {
		result := PowFloat64(test.x, test.y)
		resultFloat, _ := result.Float64()

		switch {
		case math.IsNaN(test.expected):
			if !math.IsNaN(resultFloat) {
				t.Errorf("PowFloat64(%v, %v) = %v, expected NaN", test.x, test.y, resultFloat)
			}
		case math.IsInf(test.expected, 0):
			if !math.IsInf(resultFloat, 0) {
				t.Errorf("PowFloat64(%v, %v) = %v, expected Inf", test.x, test.y, resultFloat)
			}
		default:
			diff := math.Abs(resultFloat - test.expected)
			if diff > test.tolerance {
				t.Errorf("PowFloat64(%v, %v) = %v, expected %v (diff: %v)", test.x, test.y, resultFloat, test.expected, diff)
			}
		}
	}
}

func TestPowEdgeCases(t *testing.T) {
	// Test 0^0 (should be 1)
	result := PowFloat64(0, 0)
	resultFloat, _ := result.Float64()
	if resultFloat != 1 {
		t.Errorf("PowFloat64(0, 0) = %v, expected 1", resultFloat)
	}

	// Test 0^negative (should be +Inf)
	result = PowFloat64(0, -1)
	resultFloat, _ = result.Float64()
	if !math.IsInf(resultFloat, 1) {
		t.Errorf("PowFloat64(0, -1) = %v, expected +Inf", resultFloat)
	}

	// Test negative^non-integer (should be +Inf to indicate undefined)
	result = PowFloat64(-2, 0.5)
	resultFloat, _ = result.Float64()
	if !math.IsInf(resultFloat, 1) {
		t.Errorf("PowFloat64(-2, 0.5) = %v, expected +Inf", resultFloat)
	}
}

func TestPowLargeValues(t *testing.T) {
	// Test a case that would overflow float64 but should work with big.Float
	result := PowFloat64(2, 100)
	expected := math.Pow(2, 100)
	resultFloat, _ := result.Float64()

	// For very large values, check if they're both infinite or compare within tolerance
	if math.IsInf(expected, 0) && math.IsInf(resultFloat, 0) {
		// Both are infinite, that's correct
		return
	}

	// The big.Float should NOT be Inf or a NaN.
	if result.IsInf() {
		t.Errorf("PowFloat64(2, 100) = %v, should not be Infinite.", result)
	}

	if !math.IsInf(expected, 0) {
		diff := math.Abs(resultFloat - expected)
		if diff > 1e10 { // Large tolerance for very large numbers
			t.Errorf("PowFloat64(2, 100) = %v, expected %v (diff: %v)", resultFloat, expected, diff)
		}
	}
}

func TestPowSpecialCases(t *testing.T) {
	zero := big.NewFloat(0)
	one := big.NewFloat(1)
	negOne := big.NewFloat(-1)
	two := big.NewFloat(2)
	half := big.NewFloat(0.5)
	negZero, _, _ := new(big.Float).Parse("-0.0", 10)

	// Create infinity values
	posInf := big.NewFloat(0)
	posInf.SetInf(false)
	negInf := big.NewFloat(0)
	negInf.SetInf(true)

	tests := []struct {
		name     string
		x, y     *big.Float
		expected string // "0", "1", "+Inf", "-Inf", or actual value
	}{
		// Pow(x, ±0) = 1 for any x
		{"any^0", two, zero, "1"},
		{"any^0_neg", negOne, zero, "1"},
		{"inf^0", posInf, zero, "1"},

		// Pow(1, y) = 1 for any y
		{"1^any", one, two, "1"},
		{"1^inf", one, posInf, "1"},
		{"1^-inf", one, negInf, "1"},

		// Pow(x, 1) = x for any x
		{"2^1", two, one, "2"},
		{"(-1)^1", negOne, one, "-1"},

		// Pow(±0, y) cases
		{"0^2", zero, two, "0"},
		{"0^3", zero, big.NewFloat(3), "0"},              // odd positive
		{"(-0)^3", negZero, big.NewFloat(3), "0"},        // odd positive, should be -0 but big.Float may not preserve
		{"0^2.5", zero, big.NewFloat(2.5), "0"},          // non-integer positive
		{"0^(-2)", zero, big.NewFloat(-2), "+Inf"},       // negative even
		{"0^(-3)", zero, big.NewFloat(-3), "+Inf"},       // negative odd
		{"(-0)^(-3)", negZero, big.NewFloat(-3), "-Inf"}, // big.Float may not preserve -0 properly
		{"0^(+Inf)", zero, posInf, "0"},
		{"0^(-Inf)", zero, negInf, "+Inf"},

		// Pow(-1, ±Inf) = 1
		{"(-1)^(+Inf)", negOne, posInf, "1"},
		{"(-1)^(-Inf)", negOne, negInf, "1"},

		// Pow(x, +Inf) and Pow(x, -Inf) for |x| > 1
		{"2^(+Inf)", two, posInf, "+Inf"},
		{"2^(-Inf)", two, negInf, "0"},
		{"(-2)^(+Inf)", big.NewFloat(-2), posInf, "+Inf"},
		{"(-2)^(-Inf)", big.NewFloat(-2), negInf, "0"},

		// Pow(x, +Inf) and Pow(x, -Inf) for |x| < 1
		{"0.5^(+Inf)", half, posInf, "0"},
		{"0.5^(-Inf)", half, negInf, "+Inf"},
		{"(-0.5)^(+Inf)", big.NewFloat(-0.5), posInf, "0"},
		{"(-0.5)^(-Inf)", big.NewFloat(-0.5), negInf, "+Inf"},

		// Pow(+Inf, y) cases
		{"(+Inf)^2", posInf, two, "+Inf"},
		{"(+Inf)^(-2)", posInf, big.NewFloat(-2), "0"},

		// Pow(-Inf, y) = Pow(-0, -y) - test key cases
		{"(-Inf)^2", negInf, two, "+Inf"}, // Pow(-Inf, 2) = Pow(-0, -2) = +Inf

		// Pow(x, y) = undefined for finite x < 0 and finite non-integer y
		{"(-2)^0.5", big.NewFloat(-2), half, "+Inf"},              // undefined, returns +Inf
		{"(-2)^2.5", big.NewFloat(-2), big.NewFloat(2.5), "+Inf"}, // undefined
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Pow(test.x, test.y)

			switch test.expected {
			case "0":
				if result.Cmp(zero) != 0 {
					t.Errorf("Pow(%v, %v) = %v, expected 0", test.x, test.y, result)
				}
			case "1":
				if result.Cmp(one) != 0 {
					t.Errorf("Pow(%v, %v) = %v, expected 1", test.x, test.y, result)
				}
			case "2":
				if result.Cmp(two) != 0 {
					t.Errorf("Pow(%v, %v) = %v, expected 2", test.x, test.y, result)
				}
			case "-1":
				if result.Cmp(negOne) != 0 {
					t.Errorf("Pow(%v, %v) = %v, expected -1", test.x, test.y, result)
				}
			case "+Inf":
				if !result.IsInf() || result.Signbit() {
					t.Errorf("Pow(%v, %v) = %v, expected +Inf", test.x, test.y, result)
				}
			case "-Inf":
				if !result.IsInf() || !result.Signbit() {
					t.Errorf("Pow(%v, %v) = %v, expected -Inf", test.x, test.y, result)
				}
			}
		})
	}
}

func TestPowFloat64SpecialCases(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected string // "0", "1", "+Inf", "-Inf", "NaN"
	}{
		// NaN cases
		{"NaN^2", math.NaN(), 2, "+Inf"}, // big.Float returns +Inf for undefined
		{"2^NaN", 2, math.NaN(), "+Inf"},
		{"NaN^NaN", math.NaN(), math.NaN(), "+Inf"},

		// Infinity input cases
		{"(+Inf)^2", math.Inf(1), 2, "+Inf"},
		{"(+Inf)^(-2)", math.Inf(1), -2, "0"},
		{"(-Inf)^2", math.Inf(-1), 2, "+Inf"}, // Pow(-Inf, 2) = Pow(-0, -2) = +Inf
		{"2^(+Inf)", 2, math.Inf(1), "+Inf"},
		{"2^(-Inf)", 2, math.Inf(-1), "0"},
		{"0.5^(+Inf)", 0.5, math.Inf(1), "0"},
		{"0.5^(-Inf)", 0.5, math.Inf(-1), "+Inf"},

		// Regular special cases
		{"any^0", 5, 0, "1"},
		{"1^any", 1, 123, "1"},
		{"0^positive", 0, 3, "0"},
		{"0^negative", 0, -2, "+Inf"},
		{"negative^non-integer", -2, 0.5, "+Inf"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := PowFloat64(test.x, test.y)
			resultFloat, _ := result.Float64()

			switch test.expected {
			case "0":
				if resultFloat != 0 {
					t.Errorf("PowFloat64(%v, %v) = %v, expected 0", test.x, test.y, resultFloat)
				}
			case "1":
				if resultFloat != 1 {
					t.Errorf("PowFloat64(%v, %v) = %v, expected 1", test.x, test.y, resultFloat)
				}
			case "+Inf":
				if !math.IsInf(resultFloat, 1) && !result.IsInf() {
					t.Errorf("PowFloat64(%v, %v) = %v, expected +Inf", test.x, test.y, resultFloat)
				}
			case "-Inf":
				if !(math.IsInf(resultFloat, -1) || (result.IsInf() && result.Signbit())) {
					t.Errorf("PowFloat64(%v, %v) = %v, expected -Inf", test.x, test.y, resultFloat)
				}
			}
		})
	}
}

// BenchmarkPow benchmarks the Pow function with various input sizes
func BenchmarkPowFloat64(b *testing.B) {
	testCases := []struct {
		name string
		x, y float64
	}{
		{"Small_2^3", 2, 3},
		{"Medium_10^10", 10, 10},
		{"Large_2^100", 2, 100},
		{"Fractional_4^0.5", 4, 0.5},
		{"Negative_Exp_2^-5", 2, -5},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = PowFloat64(tc.x, tc.y)
			}
		})
	}
}

// BenchmarkComparative benchmarks compare bigmath functions with standard library
func BenchmarkComparativePow(b *testing.B) {
	x, y := 2.0, 10.0

	benchmarkBigmathVsStdlib2Args(b, powMethods[0], big.NewFloat(x), big.NewFloat(y))
}

// BenchmarkPrecision benchmarks functions with different precision settings
func BenchmarkPowFloat64Precision(b *testing.B) {
	x, y := 2.0, 50.0

	for _, prec := range precisions {
		b.Run(fmt.Sprintf("PowFloat64_%d_prec", prec), func(b *testing.B) {
			// Note: bigmath doesn't expose precision settings directly
			// This benchmark shows how precision might affect performance
			for b.Loop() {
				result := PowFloat64(x, y)
				result.SetPrec(prec) // Set precision after calculation
				_ = result
			}
		})
	}
}

// Benchmarks for PowFloat64 function
func BenchmarkPowFloat64Specific(b *testing.B) {
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"Pow_2^10", 2, 10},
		{"Pow_E^Pi", math.E, math.Pi},
		{"Pow_1000^5", 1000, 5},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for b.Loop() {
				_ = PowFloat64(test.x, test.y)
			}
		})
	}
}

func BenchmarkPowFloat64Scaling(b *testing.B) {
	// Test how Pow scales with exponent size
	exponents := []float64{10, 50, 100, 200}
	bases := []float64{2.0, 3.0, 4.0, 10.0}

	for _, base := range bases {
		for _, exp := range exponents {
			b.Run(fmt.Sprintf("Pow_%.0f^%.0f", base, exp), func(b *testing.B) {
				for b.Loop() {
					_ = PowFloat64(base, exp)
				}
			})
		}
	}
}
