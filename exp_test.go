package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestExp(t *testing.T) {
	tests := []struct {
		input     float64
		expected  float64
		tolerance float64
	}{
		{0, 1, 1e-10},                   // e^0 = 1
		{1, math.E, 1e-10},              // e^1 = e
		{2, math.E * math.E, 1e-9},      // e^2 = e²
		{-1, 1 / math.E, 1e-10},         // e^(-1) = 1/e
		{0.5, math.Sqrt(math.E), 1e-10}, // e^0.5 = √e
		{10, math.Exp(10), 1e-5},        // e^10
	}

	for _, test := range tests {
		result := Exp(big.NewFloat(test.input))
		resultFloat, _ := result.Float64()

		diff := math.Abs(resultFloat - test.expected)
		if diff > test.tolerance {
			t.Errorf("Exp(%v) = %v, expected %v (diff: %v)", test.input, resultFloat, test.expected, diff)
		}
	}
}

func TestExpEdgeCases(t *testing.T) {
	// Test very large positive values
	result := Exp(big.NewFloat(1000000))
	if !result.IsInf() {
		t.Errorf("Exp(1000000) should be +Inf, got %v", result)
	}

	// Test very large negative values
	result = Exp(big.NewFloat(-1000000))
	resultFloat, _ := result.Float64()
	if resultFloat != 0 {
		t.Errorf("Exp(-1000000) should be 0, got %v", resultFloat)
	}
}

func TestExpVsMathExp(t *testing.T) {
	tests := []struct {
		input     float64
		precision uint
		tolerance float64
	}{
		{0, 256, 1e-15},
		{1, 256, 1e-15},
		{2, 256, 1e-15},
		{-1, 256, 1e-15},
		{0.5, 256, 1e-15},
		{-0.5, 256, 1e-15},
		{10, 256, 1e-10},
		{-10, 256, 1e-15},
		{0.1, 256, 1e-15},
		{-0.1, 256, 1e-15},
		{5, 256, 1e-12},
		{-5, 256, 1e-15},
		{0.001, 256, 1e-15},
		{-0.001, 256, 1e-15},
		// Test some values where math.Exp might lose precision
		{20, 256, 1e-8},
		{-20, 256, 1e-15},
		{0.693147180559945, 256, 1e-15}, // ln(2), should give exactly 2
	}

	for _, test := range tests {
		// Create input with specified precision
		input := new(big.Float).SetPrec(test.precision).SetFloat64(test.input)

		// Compute using our Exp function
		ourResult := Exp(input)

		// Compute using math.Exp (standard library)
		mathExpResult := math.Exp(test.input)
		mathExpBigFloat := new(big.Float).SetPrec(test.precision).SetFloat64(mathExpResult)

		// Compare results
		diff := new(big.Float).SetPrec(test.precision)
		diff.Sub(ourResult, mathExpBigFloat)
		diff.Abs(diff)

		// Calculate relative error for better comparison
		relativeError := new(big.Float).SetPrec(test.precision)
		if mathExpBigFloat.Cmp(new(big.Float).SetFloat64(0)) != 0 {
			relativeError.Quo(diff, mathExpBigFloat)
		} else {
			relativeError = diff
		}

		// Convert to float64 for comparison
		relativeErrorFloat, _ := relativeError.Float64()

		if relativeErrorFloat > test.tolerance {
			ourFloat, _ := ourResult.Float64()
			t.Errorf("Exp(%v) mismatch: our=%.15g, math.Exp=%.15g, rel_error=%.2e, tolerance=%.2e",
				test.input, ourFloat, mathExpResult, relativeErrorFloat, test.tolerance)
		}
	}
}

func TestExpHighPrecision(t *testing.T) {
	// Test high precision cases where our implementation should shine
	tests := []struct {
		name      string
		inputStr  string
		precision uint
		// Expected values computed with high-precision tools
		expectedStr string
		tolerance   float64
	}{
		{
			name:        "exp(1) with 100 decimal places",
			inputStr:    "1",
			precision:   400, // ~120 decimal places
			expectedStr: "2.7182818284590452353602874713526624977572470936999595749669676277240766303535475945713821785251664274",
			tolerance:   1e-80,
		},
		{
			name:        "exp(0.5) high precision",
			inputStr:    "0.5",
			precision:   400,
			expectedStr: "1.6487212707001281468486507878141635716537761007101480115750793116406639144845081060401884159251825604",
			tolerance:   1e-60,
		},
		{
			name:        "exp(-1) high precision",
			inputStr:    "-1",
			precision:   400,
			expectedStr: "0.36787944117144232159552377016146086744581113103176783450783680169746149574489980688139251132243124071",
			tolerance:   1e-70,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Parse input with high precision
			input := new(big.Float).SetPrec(test.precision)
			input.SetString(test.inputStr)

			// Compute using our Exp function
			result := Exp(input)

			// Parse expected result
			expected := new(big.Float).SetPrec(test.precision)
			expected.SetString(test.expectedStr)

			// Calculate relative error
			diff := new(big.Float).SetPrec(test.precision)
			diff.Sub(result, expected)
			diff.Abs(diff)

			relativeError := new(big.Float).SetPrec(test.precision)
			relativeError.Quo(diff, expected)
			relativeErrorFloat, _ := relativeError.Float64()

			if relativeErrorFloat > test.tolerance {
				t.Errorf("%s: relative error %.2e exceeds tolerance %.2e",
					test.name, relativeErrorFloat, test.tolerance)
				t.Errorf("Expected: %s", expected.Text('f', 50))
				t.Errorf("Got:      %s", result.Text('f', 50))
			}
		})
	}
}

func BenchmarkExpVsMathExp(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(5.0)

	b.Run("BigMath_Exp", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Exp(x)
		}
	})

	b.Run("StdLib_Exp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = math.Exp(5.0)
		}
	})
}

// Benchmarks for Exp function
func BenchmarkExpSpecific(b *testing.B) {
	tests := []struct {
		val int64
		x   *big.Float
	}{
		{0, big.NewFloat(0)},
		{1, big.NewFloat(1)},
		{10, big.NewFloat(10)},
		{100, big.NewFloat(100)},
		{1000, big.NewFloat(1000)},
		{10000, big.NewFloat(10000)},
		{-5, big.NewFloat(-5)},
	}

	for _, test := range tests {
		b.Run(fmt.Sprintf("Exp(%d)", test.val), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = Exp(test.x)
			}
		})
	}
}

func BenchmarkExp2(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(1025.3)

	b.ResetTimer()
	for b.Loop() {
		Exp2(x)
	}
}
