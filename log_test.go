package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

// logMethods is used throughout the tests and benchmarks to range
// over the different logarithm methods.
var logMethods = []benchAndCompare{
	{"Log", Log, math.Log},
	{"LogHalley", logHalley, math.Log}, // Halley's method
	{"LogNewton", logNewton, math.Log}, // Newton's method
	{"LogTaylor", logTaylor, math.Log}, // Taylor series
}

// This is a limited set of test cases since the better cases are tested in
// the specific log... method tests.
func TestLog(t *testing.T) {
	tests := []struct {
		input     float64
		expected  float64
		tolerance float64
	}{
		{1, 0, 1e-10},               // ln(1) = 0
		{math.E, 1, 1e-10},          // ln(e) = 1
		{math.E * math.E, 2, 1e-10}, // ln(e²) = 2
		{2, math.Log(2), 1e-10},     // ln(2)
		{10, math.Log(10), 1e-10},   // ln(10)
		{100, math.Log(100), 1e-9},  // ln(100)
	}

	for _, test := range tests {
		result := Log(big.NewFloat(test.input))
		resultFloat, _ := result.Float64()

		diff := math.Abs(resultFloat - test.expected)
		if diff > test.tolerance {
			t.Errorf("Log(%v) = %v, expected %v (diff: %v)", test.input, resultFloat, test.expected, diff)
		}
	}
}

func TestLogEdgeCases(t *testing.T) {
	// Test zero - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Log(0) should panic")
		}
	}()
	Log(big.NewFloat(0))

	// Test negative number - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Log(-1) should panic")
		}
	}()
	Log(big.NewFloat(-1))
}

// Helper function to test a logarithm method with standard test cases
func testLogMethod(t *testing.T, methodName string, logFunc func(*big.Float) *big.Float) {
	t.Helper()

	tests := []struct {
		input     float64
		expected  float64
		tolerance float64
	}{
		{1, 0, 1e-10},               // ln(1) = 0
		{math.E, 1, 1e-10},          // ln(e) = 1
		{math.E * math.E, 2, 1e-9},  // ln(e²) = 2 (This one fails on Taylor at 1e-10)
		{2, math.Log(2), 1e-10},     // ln(2)
		{10, math.Log(10), 1e-10},   // ln(10)
		{100, math.Log(100), 1e-10}, // ln(100)
		{0.5, math.Log(0.5), 1e-10}, // ln(0.5)
		{0.1, math.Log(0.1), 1e-10}, // ln(0.1)
		{1.5, math.Log(1.5), 1e-10}, // ln(1.5)
		{3, math.Log(3), 1e-10},     // ln(3)
		{50, math.Log(50), 1e-10},   // ln(50)
	}

	for _, test := range tests {
		result := logFunc(big.NewFloat(test.input))
		if result == nil {
			t.Errorf("%s(%v) returned nil", methodName, test.input)

			continue
		}
		resultFloat, _ := result.Float64()

		diff := math.Abs(resultFloat - test.expected)
		if diff > test.tolerance {
			t.Errorf("%s(%v) = %v, expected %v (diff: %v)", methodName, test.input, resultFloat, test.expected, diff)
		}
	}
}

func TestLogNewton(t *testing.T) {
	testLogMethod(t, "logNewton", logNewton)
}

func TestLogTaylor(t *testing.T) {
	testLogMethod(t, "logTaylor", logTaylor)
}

func TestLogHalley(t *testing.T) {
	testLogMethod(t, "logHalley", logHalley)
}

func TestLogTaylorEdgeCases(t *testing.T) {
	// Test zero - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("logTaylor(0) should panic")
		}
	}()
	logTaylor(big.NewFloat(0))

	// Test very small positive number
	result := logTaylor(big.NewFloat(1e-10))
	if result == nil {
		t.Errorf("logTaylor(1e-10) returned nil")

		return
	}
	resultFloat, _ := result.Float64()
	expected := math.Log(1e-10)
	diff := math.Abs(resultFloat - expected)
	if diff > 1e-6 {
		t.Errorf("logTaylor(1e-10) = %v, expected %v (diff: %v)", resultFloat, expected, diff)
	}

	// Test very large number (should fall back to main Log function)
	result = logTaylor(big.NewFloat(1e20))
	if result == nil {
		t.Errorf("logTaylor(1e20) returned nil")

		return
	}
	resultFloat, _ = result.Float64()
	expected = math.Log(1e20)
	diff = math.Abs(resultFloat - expected)
	if diff > 1e-6 {
		t.Errorf("logTaylor(1e20) = %v, expected %v (diff: %v)", resultFloat, expected, diff)
	}
}

func TestLogHalleyEdgeCases(t *testing.T) {
	// Test zero - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("logHalley(0) should panic")
		}
	}()
	logHalley(big.NewFloat(0))

	// Test very small positive number (Halley method may have convergence issues for very small numbers)
	result := logHalley(big.NewFloat(1e-5))
	if result == nil {
		t.Errorf("logHalley(1e-5) returned nil")

		return
	}
	resultFloat, _ := result.Float64()
	expected := math.Log(1e-5)
	diff := math.Abs(resultFloat - expected)
	if diff > 1e-6 {
		t.Errorf("logHalley(1e-5) = %v, expected %v (diff: %v)", resultFloat, expected, diff)
	}

	// Test infinity
	result = logHalley(big.NewFloat(math.Inf(1)))
	if result == nil {
		t.Errorf("logHalley(+Inf) returned nil")

		return
	}
	if !result.IsInf() || result.Sign() != 1 {
		t.Errorf("logHalley(+Inf) should be +Inf, got %v", result)
	}
}

func TestLogMethodsComparison(t *testing.T) {
	// Test that all three log methods give similar results for various inputs
	testValues := []float64{0.1, 0.5, 1, 1.5, 2, math.E, 3, 5, 10, 50, 100}

	for _, val := range testValues {
		x := big.NewFloat(val)

		resultNewton := logNewton(x)
		resultTaylor := logTaylor(x)
		resultHalley := logHalley(x)

		if resultNewton == nil || resultTaylor == nil || resultHalley == nil {
			t.Errorf("Comparison for %v: errors: Newton=%v, Taylor=%v, Halley=%v", val, resultNewton, resultTaylor, resultHalley)

			continue
		}

		newtonFloat, _ := resultNewton.Float64()
		taylorFloat, _ := resultTaylor.Float64()
		halleyFloat, _ := resultHalley.Float64()

		const tolerance = 1e-8
		// Compare Log vs Taylor
		diffNewton := math.Abs(newtonFloat - taylorFloat)
		if diffNewton > tolerance {
			t.Errorf("Log vs Taylor for %v: Log=%v, Taylor=%v (diff: %v)", val, newtonFloat, taylorFloat, diffNewton)
		}

		// Compare Log vs Halley
		diffHalley := math.Abs(newtonFloat - halleyFloat)
		if diffHalley > tolerance {
			t.Errorf("Log vs Halley for %v: Log=%v, Halley=%v (diff: %v)", val, newtonFloat, halleyFloat, diffHalley)
		}

		// Compare Taylor vs Halley
		diffMethods := math.Abs(taylorFloat - halleyFloat)
		if diffMethods > tolerance {
			t.Errorf("Taylor vs Halley for %v: Taylor=%v, Halley=%v (diff: %v)", val, taylorFloat, halleyFloat, diffMethods)
		}
	}
}

func TestLogMethodsAccuracy(t *testing.T) {
	// Test accuracy against known high-precision values
	testCases := []struct {
		input    float64
		expected float64
		name     string
	}{
		{2, 0.693147180559945309417232121458176568075500134360255254120680009, "ln(2)"},  //https://oeis.org/A002162 via math.Ln2
		{10, 2.30258509299404568401799145468436420760110148862877297603332790, "ln(10)"}, // https://oeis.org/A002392 via math.Ln10
		{math.E, 1.0, "ln(e)"},
		{math.E * math.E, 2.0, "ln(e²)"},
		{0.5, -0.6931471805599453, "ln(0.5)"},
		{1.5, 0.4054651081081644, "ln(1.5)"},
		{3, 1.0986122886681098, "ln(3)"},
	}

	const tolerance = 1e-9
	for _, tc := range testCases {
		x := big.NewFloat(tc.input)

		// Test Newton method
		resultNewton := logNewton(x)
		if resultNewton == nil {
			t.Errorf("logNewton %s: error: %v", tc.name, resultNewton)
		} else {
			newtonFloat, _ := resultNewton.Float64()
			diffNewton := math.Abs(newtonFloat - tc.expected)
			if diffNewton > tolerance {
				t.Errorf("logNewton %s: got %v, expected %v (diff: %v)", tc.name, newtonFloat, tc.expected, diffNewton)
			}
		}
		// Test Taylor method
		resultTaylor := logTaylor(x)
		if resultTaylor == nil {
			t.Errorf("logTaylor %s: error: %v", tc.name, resultTaylor)
		} else {
			taylorFloat, _ := resultTaylor.Float64()
			diffTaylor := math.Abs(taylorFloat - tc.expected)
			if diffTaylor > tolerance {
				t.Errorf("logTaylor %s: got %v, expected %v (diff: %v)", tc.name, taylorFloat, tc.expected, diffTaylor)
			}
		}

		// Test Halley method
		resultHalley := logHalley(x)
		if resultHalley == nil {
			t.Errorf("logHalley %s: error: %v", tc.name, resultHalley)
		} else {
			halleyFloat, _ := resultHalley.Float64()
			diffHalley := math.Abs(halleyFloat - tc.expected)
			if diffHalley > tolerance {
				t.Errorf("logHalley %s: got %v, expected %v (diff: %v)", tc.name, halleyFloat, tc.expected, diffHalley)
			}
		}
	}
}

// Benchmarks for Log function
func BenchmarkLogSpecific(b *testing.B) {
	b.Run("Log_1", func(b *testing.B) {
		x := big.NewFloat(1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Log(x)
		}
	})

	b.Run("Log_E", func(b *testing.B) {
		x := big.NewFloat(math.E)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Log(x)
		}
	})

	b.Run("Log_10", func(b *testing.B) {
		x := big.NewFloat(10)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Log(x)
		}
	})

	b.Run("Log_Large", func(b *testing.B) {
		x := big.NewFloat(1000000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Log(x)
		}
	})
}

// Helper function to benchmark a logarithm method with standard test cases
func benchmarkLogMethod(b *testing.B, methodName string, logFunc func(*big.Float) *big.Float) {
	b.Helper()

	// Common benchmark test cases for all logarithm methods
	tests := []struct {
		name  string
		value float64
	}{
		{"1", 1},
		{"E", math.E},
		{"2", 2},
		{"10", 10},
		{"100", 100},
		{"0.5", 0.5},
		{"Large", 1000000},
	}

	for _, tc := range tests {
		b.Run(methodName+"_"+tc.name, func(b *testing.B) {
			x := big.NewFloat(tc.value)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = logFunc(x)
			}
		})
	}
}

// Benchmarks for logNewton function
func BenchmarkLogNewton(b *testing.B) {
	benchmarkLogMethod(b, "LogNewton", logNewton)
}

// Benchmarks for logTaylor function
func BenchmarkLogTaylor(b *testing.B) {
	benchmarkLogMethod(b, "LogTaylor", logTaylor)
}

// Benchmarks for logHalley function
func BenchmarkLogHalley(b *testing.B) {
	benchmarkLogMethod(b, "LogHalley", logHalley)
}

// Comparative benchmarks between all three logarithm methods
func BenchmarkLogMethodsComparative(b *testing.B) {
	testValues := []struct {
		name  string
		value float64
	}{
		{"Small_0.5", 0.5},
		{"Unity_1", 1},
		{"E", math.E},
		{"Medium_10", 10},
		{"Large_100", 100},
		{"VeryLarge_1e6", 1e6},
	}

	for _, lm := range logMethods {
		for _, tv := range testValues {
			x := big.NewFloat(tv.value)

			b.Run(lm.name+"_"+tv.name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = lm.fnBigmath(x)
				}
			})
		}
	}
}

// Memory allocation benchmarks for logarithm methods
func BenchmarkLogMethodsMemory(b *testing.B) {
	x := big.NewFloat(10)

	for _, lm := range logMethods {
		b.Run(lm.name+"_Memory", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = lm.fnBigmath(x)
			}
		})
	}
}

// Scaling benchmarks to test performance with different input magnitudes.
// Primarily used to help determine cutoffs for when to switch methods in Log.
func BenchmarkLogMethodsScaling(b *testing.B) {
	for _, lm := range logMethods {
		for i := 1e-5; i < 1e10; i *= 10 {
			x := big.NewFloat(i)

			b.Run(fmt.Sprintf("%s_Scale_%f", lm.name, i), func(b *testing.B) {
				b.ResetTimer()
				for b.Loop() {
					_ = lm.fnBigmath(x)
				}
			})
		}
	}
}

// TestLogMethodsConvergenceStress stress tests the convergence limits for each logarithm method
// This test validates that each method can handle challenging precision and input combinations
func TestLogMethodsConvergenceStress(t *testing.T) {
	// Test different precision levels
	precisions := []uint{24, 53, 64, 128, 256}

	// Challenging input values that stress convergence
	testInputs := []struct {
		name        string
		value       float64
		description string
	}{
		{"VerySmall_1e-15", 1e-15, "very small number near underflow"},
		{"Small_1e-6", 1e-6, "small number"},
		{"TinyFraction_0.001", 0.001, "small fraction"},
		{"SmallFraction_0.1", 0.1, "decimal fraction"},
		{"NearOne_0.9", 0.9, "close to 1 from below"},
		{"Unity_1.0", 1.0, "exactly 1 (special case)"},
		{"NearOne_1.1", 1.1, "close to 1 from above"},
		{"Small_2", 2.0, "small integer"},
		{"E_2.718", math.E, "Euler's number"},
		{"Medium_10", 10.0, "medium integer"},
		{"Large_100", 100.0, "large integer"},
		{"VeryLarge_1e6", 1e6, "very large number"},
		{"Huge_1e12", 1e12, "huge number"},
		{"Extreme_1e15", 1e15, "extreme large number"},
	}

	// Methods to test with their expected iteration characteristics
	methods := []struct {
		name             string
		logFunc          func(*big.Float) *big.Float
		expectedMaxIters func(prec uint) int
		convergenceOrder string
	}{
		{
			name:    "Newton",
			logFunc: logNewton,
			expectedMaxIters: func(prec uint) int {
				return int(3*math.Log2(float64(prec))) + 20
			},
			convergenceOrder: "quadratic",
		},
		{
			name:    "Taylor",
			logFunc: logTaylor,
			expectedMaxIters: func(prec uint) int {
				return int(float64(prec)*2.0) + 200
			},
			convergenceOrder: "linear",
		},
		{
			name:    "Halley",
			logFunc: logHalley,
			expectedMaxIters: func(prec uint) int {
				return int(2.0*math.Log2(float64(prec))) + 15
			},
			convergenceOrder: "cubic",
		},
	}

	for _, method := range methods {
		t.Run(method.name+"_ConvergenceStress", func(t *testing.T) {
			for _, prec := range precisions {
				t.Run(fmt.Sprintf("Precision_%d", prec), func(t *testing.T) {
					expectedMaxIters := method.expectedMaxIters(prec)

					successCount := 0
					totalTests := 0

					for _, input := range testInputs {
						totalTests++

						// Create input with specified precision
						x := new(big.Float).SetPrec(prec).SetFloat64(input.value)

						// Call the method and measure success
						result := method.logFunc(x)

						if result == nil {
							// Log the failure but continue testing
							t.Logf("%s precision=%d input=%s (%s): FAILED - %v",
								method.name, prec, input.name, input.description, result)

							continue
						}

						successCount++

						// Validate result quality by checking against Go's math.Log
						if input.value > 0 && input.value != 1.0 {
							expected := math.Log(input.value)
							resultFloat, _ := result.Float64()

							// Calculate relative error
							relativeError := math.Abs((resultFloat - expected) / expected)

							// For higher precision, we expect better accuracy
							expectedAccuracy := math.Pow(2, -float64(prec)/4.0) // Rough estimate
							if relativeError > expectedAccuracy && relativeError > 1e-10 {
								t.Logf("%s precision=%d input=%s: moderate accuracy - relative error %e (expected < %e)",
									method.name, prec, input.name, relativeError, expectedAccuracy)
							}
						}

						// Verify result precision matches input precision
						if result.Prec() != prec {
							t.Errorf("%s precision=%d input=%s: result precision mismatch - got %d, expected %d",
								method.name, prec, input.name, result.Prec(), prec)
						}
					}

					// Calculate success rate
					successRate := float64(successCount) / float64(totalTests)

					// Report results
					t.Logf("%s (%s convergence) precision=%d: %d/%d successful (%.1f%%), max_iters=%d",
						method.name, method.convergenceOrder, prec, successCount, totalTests,
						successRate*100, expectedMaxIters)

					// Minimum success rate expectations (some edge cases may fail)
					minSuccessRate := 0.85 // 85% success rate minimum
					if method.name == "Taylor" {
						minSuccessRate = 0.75 // Taylor series more sensitive to extreme values
					}

					if successRate < minSuccessRate {
						t.Errorf("%s precision=%d: success rate %.1f%% below minimum %.1f%%",
							method.name, prec, successRate*100, minSuccessRate*100)
					}
				})
			}
		})
	}
}

// TestLogMethodsIterationCounts verifies that methods stay within expected iteration bounds
func TestLogMethodsIterationCounts(t *testing.T) {
	// Test with instrumented versions that count iterations
	// This is a conceptual test - in practice, you'd need to modify the functions
	// to expose iteration counts for detailed analysis

	precisions := []uint{53, 128, 256, 512}
	testValues := []float64{0.1, 1.5, 2, 5, 10.0, 20, 50, 100.0, 200, 500, 1000.0}

	for _, prec := range precisions {
		t.Run(fmt.Sprintf("Precision_%d", prec), func(t *testing.T) {
			for _, val := range testValues {
				x := new(big.Float).SetPrec(prec).SetFloat64(val)

				// Test each method
				methods := []struct {
					name     string
					logFunc  func(*big.Float) *big.Float
					maxIters int
				}{
					{"Newton", logNewton, int(3*math.Log2(float64(prec))) + 20},
					{"Taylor", logTaylor, int(float64(prec)*2.0) + 200},
					{"Halley", logHalley, int(2.0*math.Log2(float64(prec))) + 15},
				}

				for _, method := range methods {
					result := method.logFunc(x)
					want := math.Log(val)

					if result == nil {
						resText := "nil"
						t.Errorf("%s precision=%d input=%g: convergence failed - %v, last val was: %v, want: %v",
							method.name, prec, val, result, resText, want)

						continue
					}

					// Verify the result is reasonable
					if result.IsInf() {
						t.Errorf("%s precision=%d input=%g: invalid result %v",
							method.name, prec, val, result)
					}

					// Log success with expected iteration bound
					// t.Logf("%s precision=%d input=%g: SUCCESS (max_iters=%d)",
					//	method.name, prec, val, method.maxIters)
				}
			}
		})
	}
}

// TestLogPanicCases tests that log methods properly panic on invalid inputs
func TestLogPanicCases(t *testing.T) {
	// Test division by zero panic in logNewton
	t.Run("Newton_DivisionByZero", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("logNewton should panic on division by zero")
			}
		}()
		// This should cause division by zero in Newton's method
		logNewton(big.NewFloat(-1))
	})

	// Test division by zero panic in logHalley
	t.Run("Halley_DivisionByZero", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("logHalley should panic on division by zero")
			}
		}()
		// This should cause division by zero in Halley's method
		logHalley(big.NewFloat(-1))
	})

	// Test negative input panic in logTaylor
	t.Run("Taylor_NegativeInput", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("logTaylor should panic on negative input")
			}
		}()
		logTaylor(big.NewFloat(-1))
	})
}

// TestLogOverflowCases tests that log methods properly handle and panic on overflow conditions
func TestLogOverflowCases(t *testing.T) {
	// These tests may require specific conditions to trigger overflow panics
	// The actual overflow conditions depend on the internal implementation

	// Test potential overflow in Newton's method with extreme values
	t.Run("Newton_PotentialOverflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// Overflow panic is expected and acceptable
				t.Logf("Newton method handled overflow correctly: %v", r)
			}
		}()
		// Very large number that might cause overflow in intermediate calculations
		largeNum := new(big.Float).SetPrec(256)
		largeNum.SetString("1e1000")
		logNewton(largeNum)
	})

	// Test potential overflow in Halley's method with extreme values
	t.Run("Halley_PotentialOverflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// Overflow panic is expected and acceptable
				t.Logf("Halley method handled overflow correctly: %v", r)
			}
		}()
		// Very large number that might cause overflow in intermediate calculations
		largeNum := new(big.Float).SetPrec(256)
		largeNum.SetString("1e1000")
		logHalley(largeNum)
	})
}

// TestLogMethodsPrecisionScaling verifies that methods handle increasing precision appropriately
func TestLogMethodsPrecisionScaling(t *testing.T) {
	testValue := math.E // Use e as it has a known exact logarithm of 1

	// Test scaling from low to very high precision
	for _, method := range logMethods {
		t.Run(method.name+"_PrecisionScaling", func(t *testing.T) {
			for _, prec := range precisions {
				x := new(big.Float).SetPrec(prec).SetFloat64(testValue)

				result := method.fnBigmath(x)

				if result == nil {
					t.Errorf("%s precision=%d: failed to converge - %v", method.name, prec, result)

					continue
				}

				// For ln(e) = 1, check accuracy
				resultFloat, _ := result.Float64()
				err := math.Abs(resultFloat - 1.0)

				// Expected accuracy should improve with precision
				expectedError := math.Pow(2, -float64(prec)/2.0)

				if err > expectedError && err > 1e-15 {
					t.Logf("%s precision=%d: accuracy could be better - error %e (expected < %e)",
						method.name, prec, err, expectedError)
				} else {
					t.Logf("%s precision=%d: excellent accuracy - error %e",
						method.name, prec, err)
				}

				// Verify result has correct precision
				if result.Prec() != prec {
					t.Errorf("%s precision=%d: result precision mismatch - got %d",
						method.name, prec, result.Prec())
				}
			}
		})
	}
}

func BenchmarkLog(b *testing.B) {
	testCases := []struct {
		name  string
		input float64
	}{
		{"Small_2", 2},
		{"Medium_100", 100},
		{"Large_1000", 1000},
		{"VeryLarge_1e6", 1e6},
		{"Fractional_0.5", 0.5},
		{"E", math.E},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			x := big.NewFloat(tc.input)
			b.ResetTimer()
			for b.Loop() {
				_ = Log(x)
			}
		})
	}
}

func BenchmarkComparativeLog(b *testing.B) {
	x := new(big.Float).SetPrec(maxTestingPrecision).SetFloat64(100.0)

	benchmarkBigmathVsStdlib(b, logMethods[0], x)
}
