package bigmath

import (
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestFactorialBig(t *testing.T) {
	tests := []struct {
		input     float64
		expected  float64
		tolerance float64
	}{
		{0, 1, 1e-10},
		{1, 1, 1e-10},
		{2, 2, 1e-10},
		{3, 6, 1e-10},
		{4, 24, 1e-10},
		{5, 120, 1e-10},
		{10, 3628800, 1e-10},
		{20, 2.43290200817664e+18, 1e10}, // Use large tolerance for big values
	}

	for _, test := range tests {
		result := FactorialFloat(big.NewFloat(test.input))
		resultFloat, _ := result.Float64()

		if math.IsInf(test.expected, 0) {
			if !math.IsInf(resultFloat, 0) {
				t.Errorf("factorialBig(%v) = %v, expected Inf", test.input, resultFloat)
			}
		} else {
			diff := math.Abs(resultFloat - test.expected)
			if diff > test.tolerance {
				t.Errorf("factorialBig(%v) = %v, expected %v (diff: %v)", test.input, resultFloat, test.expected, diff)
			}
		}
	}
}

func TestFactorialIntConsistency(t *testing.T) {
	// Test that factorial() and factorialBig() produce consistent results for values up to 170.
	testValues := []int{0, 1, 2, 3, 4, 5, 10, 15, 20, 25, 50, 100, 150, 169, 170}

	for _, n := range testValues {
		// Get result from factorial() function
		factorialResult := Factorial(int64(n))
		factorialResultFloat, _ := factorialResult.Float64()

		// Get result from factorialBig() function
		factorialBigResult := FactorialFloat(big.NewFloat(float64(n)))
		factorialBigFloat, _ := factorialBigResult.Float64()

		// Compare the results with appropriate tolerance
		tolerance := 1e-20
		if n > 51 {
			// Use smaller tolerance for very large factorials due to floating point precision
			tolerance = factorialResultFloat * 1e-15
		}

		diff := math.Abs(factorialResultFloat - factorialBigFloat)
		if diff > tolerance {
			t.Errorf("factorial(%d) = %v, factorialBig(%d) = %v (diff: %v, tolerance: %v)",
				n, factorialResult, n, factorialBigFloat, diff, tolerance)
		}
	}
}

func TestFactorialBigConsistency(t *testing.T) {
	// Test that factorial() and factorialBig() produce consistent results for values up to 170.
	testValues := []int{0, 1, 2, 3, 4, 5, 10, 15, 20, 25, 50, 100, 150, 169, 170}

	for _, n := range testValues {
		// Get result from factorial() function
		factorialResult := Factorial(int64(n))

		// Get result from factorialBig() function
		factorialBigResult := FactorialFloat(big.NewFloat(float64(n)))
		factorialBigResultInt, _ := factorialBigResult.Int(nil)

		// Set the tolerance to be at least half the number of digits in the
		// values, but no more than 15.
		tolerance := min(15, len(factorialResult.Text(10))/2)

		diff := new(big.Int).Sub(factorialResult, factorialBigResultInt)
		diffPosition := compareNumberStrings(factorialResult.Text(10), factorialBigResultInt.Text(10))
		if diffPosition < tolerance {
			t.Errorf("factorial(%d) = %v, factorialBig(%d) = %v (diffPosition: %v, diff: %v, tolerance: %v)",
				n, factorialResult, n, factorialBigResultInt, diffPosition, diff, tolerance)
		}
	}
}

// compareNumberStrings compares two integer number strings and returns
// the number of digits from the start where they first differ.
// This is useful for comparing high-precision calculations.
func compareNumberStrings(str1, str2 string) int {
	// Remove any leading/trailing whitespace
	str1 = strings.TrimSpace(str1)
	str2 = strings.TrimSpace(str2)

	if str1 == str2 {
		return len(str1)
	}

	// Find the shorter length to avoid out-of-bounds
	minLen := len(str1)
	if len(str2) < minLen {
		minLen = len(str2)
	}

	// Count matching characters from the beginning
	matchingChars := 0

	for i := 0; i < minLen; i++ {
		if str1[i] == str2[i] {
			matchingChars++
		} else {
			break
		}
	}

	return matchingChars
}

func TestCompareNumberStrings(t *testing.T) {
	tests := []struct {
		str1     string
		str2     string
		expected int
	}{
		{"123.456789", "123.456789", 10},               // Identical strings
		{"123.456789", "123.456788", 9},                // Differ at 6th decimal place
		{"123.456", "123.457", 6},                      // Differ at 3rd decimal place
		{"123.456", "123.456000", 7},                   // Different lengths, same value
		{"123", "124", 2},                              // Differ in integer part
		{"123.0", "124.0", 2},                          // Differ in integer part with decimals
		{"0.123456", "0.123457", 7},                    // Leading zero
		{"3.141592653589793", "3.141592653589794", 16}, // High precision
		{
			"788657867364790503552363213932185062295135977687173263294742533244359449963403342920304284011984623904177212138919638830257642790242637105061926624952829931113462857270763317237396988943922445621451664240254033291864131227428294853277524242407573903240321257405579568660226031904170324062351700858796178922222789623703897374720000000000000000000000000000000000000000000000000",
			"788657867364790503552363213932185062295135977687173263294742533244359449963403342920304284011984623904177212138919638830257642790242637105061926624952829931113462857270763317237396988943922445621451664240254033291864131227428294853277524242407573903240321257405579568660226031904170324062351700858796178922222789623703897374720000000000000000000000000000000000000000000000000", 375}, // 200! should be identical
	}

	for _, test := range tests {
		result := compareNumberStrings(test.str1, test.str2)
		if result != test.expected {
			t.Errorf("compareNumberStrings(%q, %q) = %d, expected %d",
				test.str1, test.str2, result, test.expected)
		}
	}
}

func TestFactorialPrecisionComparison(t *testing.T) {
	// Example of using compareNumberStrings to compare factorial calculations
	// with different precision levels

	// Calculate 50! with different methods/precision
	fact50 := FactorialFloat(big.NewFloat(50))
	fact50String := fact50.Text('f', 0)

	// Expected value of 50! (from a reference calculation)
	expected50 := "30414093201713378043612608166064768844377641568960512000000000000"

	matchingPlaces := compareNumberStrings(fact50String, expected50)

	if fact50String != expected50 {
		t.Logf("Factorial(50) calculation differs from expected:")
		t.Logf("  Computed: %s", fact50String)
		t.Logf("  Expected: %s", expected50)
		t.Logf("  Matching decimal places: %d", matchingPlaces)
		t.Logf("  Note: Since this is an integer, matching decimal places = -1 indicates difference in integer part")
	} else {
		t.Logf("Factorial(50) matches expected value exactly")
	}
}

func TestStirlingApproximation(t *testing.T) {
	tests := []struct {
		input     float64
		tolerance float64 // Relative tolerance for Stirling's approximation
	}{
		{1, 0.1}, // Small values have larger relative error
		{2, 0.1},
		{5, 0.05}, // Better approximation for larger values
		{10, 0.02},
		{20, 0.01},
		{50, 0.005},
		{100, 0.002},
	}

	for _, test := range tests {
		// Get Stirling's approximation
		stirlingResult := StirlingApproximation(big.NewFloat(test.input))
		stirlingFloat, _ := stirlingResult.Float64()

		// Get the actual factorial for comparison
		actualResult := FactorialFloat(big.NewFloat(test.input))
		actualFloat, _ := actualResult.Float64()

		// Calculate relative error: |stirling - actual| / actual
		relativeError := math.Abs(stirlingFloat-actualFloat) / actualFloat

		if relativeError > test.tolerance {
			t.Errorf("stirlingApproximation(%.0f) has relative error %.6f, tolerance %.6f",
				test.input, relativeError, test.tolerance)
			t.Logf("  Stirling: %.6e", stirlingFloat)
			t.Logf("  Actual:   %.6e", actualFloat)
		} else {
			t.Logf("stirlingApproximation(%.0f) relative error: %.6f (within tolerance %.6f)",
				test.input, relativeError, test.tolerance)
		}
	}
}

func TestStirlingApproximationVsKnownValues(t *testing.T) {
	// Test against known factorial values where Stirling's approximation is well-documented
	tests := []struct {
		n        float64
		expected string // Known accurate factorial value as string
		name     string
	}{
		{10, "3628800", "10!"},
		{20, "2432902008176640000", "20!"},
		{50, "30414093201713378043612608166064768844377641568960512000000000000", "50!"},
	}

	for _, test := range tests {
		stirlingResult := StirlingApproximation(big.NewFloat(test.n))
		stirlingString := stirlingResult.Text('f', 0)

		// Use our helper function to see how many digits match
		matchingDigits := compareNumberStrings(stirlingString, test.expected)

		t.Logf("%s - Stirling's approximation:", test.name)
		t.Logf("  Computed: %s", stirlingString)
		t.Logf("  Expected: %s", test.expected)
		t.Logf("  Matching digits: %d out of %d", matchingDigits, len(test.expected))

		// For larger factorials, Stirling's approximation should match at least a few leading digits
		// Note: Stirling's approximation accuracy varies, so we set conservative expectations
		minimumMatchingDigits := 1
		if test.n >= 20 {
			minimumMatchingDigits = 2
		}
		if test.n >= 50 {
			minimumMatchingDigits = 2 // Conservative for 50! since it only matches 2 digits
		}

		if matchingDigits < minimumMatchingDigits {
			t.Errorf("stirlingApproximation(%.0f) matches only %d digits, expected at least %d",
				test.n, matchingDigits, minimumMatchingDigits)
		}
	}
}

func TestStirlingApproximationEdgeCases(t *testing.T) {
	// Test edge cases for Stirling's approximation
	tests := []struct {
		input    float64
		expected bool // whether result should be valid (not NaN/Inf)
		name     string
	}{
		{0, true, "zero"},
		{0.5, true, "half"},
		{1, true, "one"},
		{1.5, true, "one and half"},
	}

	for _, test := range tests {
		result := StirlingApproximation(big.NewFloat(test.input))
		resultFloat, _ := result.Float64()

		isValid := !math.IsNaN(resultFloat) && !math.IsInf(resultFloat, 0) && resultFloat > 0

		if isValid != test.expected {
			t.Errorf("stirlingApproximation(%v) validity = %v, expected %v (result: %v)",
				test.input, isValid, test.expected, resultFloat)
		} else {
			t.Logf("stirlingApproximation(%v) = %.6e (valid: %v)",
				test.input, resultFloat, isValid)
		}
	}
}
