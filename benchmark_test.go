package bigmath

import (
	"fmt"
	"math/big"
	"testing"
)

// benchAndCompare defines a common mapping of name and the local bigmath function
// and the standard Go package math method to use to compare accuracy and timings
// against. (or a simple wrapper to make it so)
type benchAndCompare struct {
	name      string
	fnBigmath func(*big.Float) *big.Float
	fnStdlib  func(float64) float64
}

// benchAndCompare2Args defines a common mapping of name and the local bigmath function
// and the standard Go package math method to use to compare accuracy and timings
// against for methods that take 2 params. (or a simple wrapper to make it so)
type benchAndCompare2Args struct {
	name      string
	fnBigmath func(*big.Float, *big.Float) *big.Float
	fnStdlib  func(float64, float64) float64
}

var (
	// precisions is a collection of some reasonable bits of precision to use
	// for testing big.Float values.
	precisions = []uint{53, 64, 128, 256, 500, 1000, 2000}
)

// maxTestingPrecision is the maximum precision to use when generating test vals.
const maxTestingPrecision = 2000

// benchmarkBigmathFunctionVsPrecision benchmarks a bigmath function against various precisions.
// For the purpose of this benchmark, we just use one input value with
// the assumption there is not caching or other changes to the input value.
func benchmarkBigmathFunctionVsPrecision(b *testing.B, data benchAndCompare, val *big.Float) {
	b.Helper()
	for _, prec := range precisions {
		b.Run(fmt.Sprintf("%s_prec_%d", data.name, prec), func(b *testing.B) {
			x := new(big.Float).Copy(val)
			x.SetPrec(prec)

			b.ResetTimer()
			for b.Loop() {
				data.fnBigmath(x)
			}
		})
	}
}

// benchmarkBigmathVsStdlib benchmarks a bigmath function against the standard library
// function using the given value as an input.
func benchmarkBigmathVsStdlib(b *testing.B, data benchAndCompare, val *big.Float) {
	b.Helper()
	x := new(big.Float).Copy(val)

	b.Run(fmt.Sprintf("%s_BigMath", data.name), func(b *testing.B) {
		for b.Loop() {
			data.fnBigmath(x)
		}
	})

	b.Run(fmt.Sprintf("%s_StdLib", data.name), func(b *testing.B) {
		f, _ := val.Float64()
		for b.Loop() {
			data.fnStdlib(f)
		}
	})
}

// benchmarkBigmathVsStdlib2Args benchmarks a bigmath function against the standard library
// function that takes two params using the given values as an input.
func benchmarkBigmathVsStdlib2Args(b *testing.B, data benchAndCompare2Args, val, val2 *big.Float) {
	b.Helper()
	x := new(big.Float).Copy(val)
	y := new(big.Float).Copy(val2)

	b.Run(fmt.Sprintf("%s_BigMath", data.name), func(b *testing.B) {
		for b.Loop() {
			data.fnBigmath(x, y)
		}
	})

	b.Run(fmt.Sprintf("%s_StdLib", data.name), func(b *testing.B) {
		fX, _ := val.Float64()
		fY, _ := val2.Float64()
		for b.Loop() {
			data.fnStdlib(fX, fY)
		}
	})
}
