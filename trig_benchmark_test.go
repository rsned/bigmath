package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

// This file has a collection of comparison benchmarks for the trigonometric functions.

var (
	trigPublicMethods = []benchAndCompare{
		{"Sin", Sin, math.Sin},
		{"Cos", Cos, math.Cos},
		{"Tan", Tan, math.Tan},
		{"Sec", Sec, func(x float64) float64 { return 1 / math.Cos(x) }},
		{"Csc", Csc, func(x float64) float64 { return 1 / math.Sin(x) }},
		{"Cot", Cot, func(x float64) float64 { return math.Cos(x) / math.Sin(x) }},
		// Inverse trigonometric versions
		{"Asin", Asin, math.Asin},
		{"Acos", Acos, math.Acos},
		{"Atan", Atan, math.Atan},
		{"Asec", Asec, func(x float64) float64 { return math.Acos(1 / x) }},
		{"Acsc", Acsc, func(x float64) float64 { return math.Asin(1 / x) }},
		{"Acot", Acot, func(x float64) float64 { return math.Atan(1 / x) }},
		// Hyperbolic versions
		{"Sinh", Sinh, math.Sinh},
		{"Cosh", Cosh, math.Cosh},
		{"Tanh", Tanh, math.Tanh},
		{"Sech", Sech, func(x float64) float64 { return 1 / math.Cosh(x) }},
		{"Csch", Csch, func(x float64) float64 { return 1 / math.Sinh(x) }},
		{"Coth", Coth, func(x float64) float64 { return math.Cosh(x) / math.Sinh(x) }},
		// Inverse hyperbolic versions
		{"Asinh", Asinh, math.Asinh},
		{"Acosh", Acosh, math.Acosh},
		{"Atanh", Atanh, math.Atanh},
		{"Asech", Asech, func(x float64) float64 { return math.Acosh(1 / x) }},
		{"Acsch", Acsch, func(x float64) float64 { return math.Asinh(1 / x) }},
		{"Acoth", Acoth, func(x float64) float64 { return math.Atanh(1 / x) }},
	}

	// TODO(rsned): Move these to their respective files.
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

	cosMethods = []benchAndCompare{
		{"Cos", Cos, math.Cos},
		{"CosArgReduction", cosArgReduction, math.Cos},
		{"CosCORDIC", cosCORDIC, math.Cos},
	}

	tanMethods = []benchAndCompare{
		{"Tan", Tan, math.Tan},
		{"TanTaylor", tanTaylor, math.Tan},
		{"TanCORDIC", tanCORDIC, math.Tan},
		{"TanContinuedFraction", tanContinuedFraction, math.Tan},
	}

	cscMethods = []benchAndCompare{
		{"Csc", Csc, func(x float64) float64 { return 1 / math.Sin(x) }},
	}

	secMethods = []benchAndCompare{
		{"Sec", Sec, func(x float64) float64 { return 1 / math.Cos(x) }},
		{"SecSeries", secSeries, func(x float64) float64 { return 1 / math.Cos(x) }},
	}

	cotMethods = []benchAndCompare{
		{"Cot", Cot, func(x float64) float64 { return math.Cos(x) / math.Sin(x) }},
	}

	sinhMethods = []benchAndCompare{
		{"Sinh", Sinh, math.Sinh},
	}

	coshMethods = []benchAndCompare{
		{"Cosh", Cosh, math.Cosh},
	}

	tanhMethods = []benchAndCompare{
		{"Tanh", Tanh, math.Tanh},
	}

	sechMethods = []benchAndCompare{
		{"Sech", Sech, func(x float64) float64 { return 1 / math.Cosh(x) }},
	}

	cschMethods = []benchAndCompare{
		{"Csch", Csch, func(x float64) float64 { return 1 / math.Sinh(x) }},
	}

	cothMethods = []benchAndCompare{
		{"Coth", Coth, func(x float64) float64 { return math.Cosh(x) / math.Sinh(x) }},
	}

	// Throughout the trig tests, there are a number of common cases related
	// to Pi that are used to test the accuracy of the different algorithms.
	specialTrigValues = []float64{
		-math.Pi, -3 * math.Pi / 4, -2 * math.Pi / 3, -math.Pi / 2, -math.Pi / 3, -math.Pi / 4, -math.Pi / 6,
		0,
		math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, 2 * math.Pi / 3, 3 * math.Pi / 4, math.Pi,
	}
)

func runTrigBenchmark(b *testing.B, methods []benchAndCompare, precisions []uint) {
	b.Helper()

	for _, method := range methods {
		for _, prec := range precisions {
			x := big.NewFloat(0.987)
			x.SetPrec(prec)
			b.Run(fmt.Sprintf("%s_%d_prec", method.name, prec), func(b *testing.B) {
				for b.Loop() {
					method.fnBigmath(x)
				}
			})
		}
	}
}

// Basic timing benchmark.
func BenchmarkTrigFunctions(b *testing.B) {
	runTrigBenchmark(b, trigPublicMethods, precisions)
}

// Performance comparison with standard math library.
func BenchmarkTrigVsStdLib(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(0.9876)

	for _, method := range trigPublicMethods {
		b.Run(method.name, func(b *testing.B) {
			for b.Loop() {
				method.fnBigmath(x)
			}
		})

		b.Run(method.name+"_StdLib", func(b *testing.B) {
			for b.Loop() {
				method.fnStdlib(0.9876)
			}
		})
	}
}

// Precision scaling benchmarks.
func BenchmarkTrigPrecisionScaling(b *testing.B) {
	for _, prec := range precisions {
		x := new(big.Float).SetPrec(prec)
		x.SetFloat64(0.9876)

		for _, method := range trigPublicMethods {
			b.Run(fmt.Sprintf("%s_%d_precs", method.name, prec), func(b *testing.B) {
				for b.Loop() {
					method.fnBigmath(x)
				}
			})
		}
	}
}

// Input value scaling benchmarks.
func BenchmarkTrigCommonInputValues(b *testing.B) {
	testValues := []float64{0.1, 1.0, math.Pi / 4, math.Pi / 3, math.Pi / 2, math.Pi, 2 * math.Pi}

	for _, val := range testValues {
		x := big.NewFloat(val)
		x.SetPrec(256)

		for _, method := range trigPublicMethods {
			b.Run(fmt.Sprintf("%s_%0.4f_rad", method.name, val), func(b *testing.B) {
				for b.Loop() {
					method.fnBigmath(x)
				}
			})
		}
	}
}

// Memory allocation benchmarks.
func BenchmarkTrigMemory(b *testing.B) {
	x := big.NewFloat(.789)
	x.SetPrec(256)

	for _, method := range trigPublicMethods {
		b.Run(method.name+"_Allocs", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				method.fnBigmath(x)
			}
		})
	}
}
