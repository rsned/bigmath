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

	// Throughout the trig tests, there are a number of common cases related
	// to Pi that are used to test the accuracy of the different algorithms.
	specialTrigValues = []float64{
		-math.Pi, -3 * math.Pi / 4, -2 * math.Pi / 3, -math.Pi / 2, -math.Pi / 3, -math.Pi / 4, -math.Pi / 6,
		0,
		math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, 2 * math.Pi / 3, 3 * math.Pi / 4, math.Pi,
	}
)

// runTrigBenchmark boilerplates the work of running the trig benchmarks
// for the given methods and precisions.
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

// BenchmarkTrigFunctions runa all the public trig methods benchmarks.
func BenchmarkTrigFunctions(b *testing.B) {
	runTrigBenchmark(b, trigPublicMethods, precisions)
}

// BenchmarkTrigVsStdLib runs the trig methods against the standard math library.
func BenchmarkTrigVsStdLib(b *testing.B) {
	x := new(big.Float).SetPrec(64)
	x.SetFloat64(0.9876)

	for _, method := range trigPublicMethods {
		benchmarkBigmathVsStdlib(b, method, x)
	}
}

// BenchmarkTrigPrecisionScaling runs the trig methods against a range of precisions.
func BenchmarkTrigPrecisionScaling(b *testing.B) {
	for _, prec := range precisions {
		x := new(big.Float).SetPrec(prec)
		x.SetFloat64(0.9876)

		for _, method := range trigPublicMethods {
			benchmarkBigmathFunctionVsPrecision(b, method, x)
		}
	}
}

// BenchmarkTrigCommonInputValues runs the trig methods against a range of common input values.
func BenchmarkTrigCommonInputValues(b *testing.B) {
	//testValues := []float64{0.1, 1.0, math.Pi / 4, math.Pi / 3, math.Pi / 2, math.Pi, 2 * math.Pi}

	for _, val := range specialTrigValues {
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
