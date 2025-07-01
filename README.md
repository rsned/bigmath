# bigmath

Bigmath is a Go package that provides high-precision mathematical functions using Go's `math/big` library. It implements arbitrary-precision arithmetic methods that parallel much of what is found in the standard `math` package, enabling computations with hundreds or thousands of decimal places of precision.

## Overview

The `bigmath` package extends Go's `math/big` capabilities by providing essential mathematical functions that operate on `*big.Float` and `*big.Int` types. This allows for computations that far exceed the precision limitations of standard floating-point types.

## Features

- **High-precision arithmetic**: Computations with configurable precision (typically 100-1000+ decimal places)
- **Mathematical constants**: Pre-computed high-precision values of π and e
- **Well-tested**: Extensive test coverage with validation against known mathematical constants and standard library methods.

## Installation

```bash
go get github.com/rsned/bigmath
```

## Quick Start

```go
package main

import (
    "fmt"
    "math/big"
    "github.com/rsned/bigmath"
)

func main() {
    // Use pre-computed high-precision constants
    fmt.Printf("π = %s\n", bigmath.Pi().Text('f', 50))
    fmt.Printf("e = %s\n", bigmath.E().Text('f', 50))
    
    // Compute e^π with high precision
    x := bigmath.Pi()
    result := bigmath.Exp(x)
    fmt.Printf("e^π = %s\n", result.Text('f', 50))
    
    // Compute sqrt(2) with custom precision
    two := big.NewFloat(2).SetPrec(256)
    sqrt2 := bigmath.Sqrt(two)
    fmt.Printf("√2 = %s\n", sqrt2.Text('f', 50))
}
```

## Available Functions

### Mathematical Constants
- **`Pi`** - Pre-computed π to 1000 decimal places
- **`E`** - Pre-computed e to 1000 decimal places

### Exponential and Logarithmic Functions
- **`Exp(x *big.Float) *big.Float`** - Computes e^x using Taylor series expansion
- **`Ln(x *big.Float) *big.Float`** - Natural logarithm using high-precision algorithms
- **`Log(x *big.Float) (*big.Float, error)`** - Natural logarithm with error handling

### Power Functions
- **`Pow(x, y *big.Float) *big.Float`** - Computes x^y for arbitrary precision big.Floats
- **`PowInt(x *big.Float, n int64) *big.Float`** - Optimized for integer exponentiation
- **`PowFloat64(x, y float64) *big.Float`** - Convenience function for float64 inputs which would exceed math.MaxFloat64
- **`Sqrt(x *big.Float) *big.Float`** - Square root using combination of methods.

### Trigonometric Functions
- **`Sin(x *big.Float) *big.Float`** - Sine 
- **`Cos(x *big.Float) *big.Float`** - Cosine 
- **`Tan(x *big.Float) *big.Float`** - Tangent 
- **`Secant(x *big.Float) *big.Float`** - Sine 
- **`Cosecant(x *big.Float) *big.Float`** - Cosine
- **`Cotangent(x *big.Float) *big.Float`** - Tangent 

- **`Arcsin(x *big.Float) *big.Float`** - Sine⁻¹
- **`Arccos(x *big.Float) *big.Float`** - Cosine⁻¹
- **`Arctan(x *big.Float) *big.Float`** - Tangent⁻¹
- **`Arcsec(x *big.Float) *big.Float`** - Secant⁻¹
- **`Arccsc(x *big.Float) *big.Float`** - Cosecant⁻¹ 
- **`Arccot(x *big.Float) *big.Float`** - Cotangent⁻¹

- **`Sinh(x *big.Float) *big.Float`** - Hyperbolic Sine 
- **`Cosh(x *big.Float) *big.Float`** - Hyperbolic Cosine 
- **`Tanh(x *big.Float) *big.Float`** - Hyperbolic Tangent 
- **`Secanth(x *big.Float) *big.Float`** - Hyperbolic Secant
- **`Cosecanth(x *big.Float) *big.Float`** - Cosine using Taylor series
- **`Cotangenth(x *big.Float) *big.Float`** - Tangent using Taylor series

### Gamma and Factorial Functions
- **`Gamma(x *big.Float) *big.Float`** - Gamma function using Lanczos approximation
- **`GammaFloat64(x float64) *big.Float`** - Convenience function for float64 input
- **`Factorial(n int64) *big.Int`** - Integer factorial for large numbers
- **`FactorialFloat(x *big.Float) *big.Float`** - Factorial for non-integers using Gamma function
- **`FactorialInt(x int) *big.Float`** - Factorial for integer > 170 which would overflow normal math.
- **`StirlingApproximation(x *big.Float) *big.Float`** - Stirling's approximation

### High-Precision Constant Computation
- **`ComputePi(precision uint) *big.Float`** - Compute π using Machin's formula with the given bits of precision.
- **`ComputeE(precision uint) *big.Float`** - Compute e using series expansion with the given bits of precision. 
- **`ComputeLn2(precision uint) *big.Float`** - Compute ln(2) with high precision with the given bits of precision.

## Precision and Performance

The package is designed to handle computations with:
- **100+ decimal places**: Typical use cases with good performance
- **1000+ decimal places**: Advanced applications with reasonable performance
- **Custom precision**: Configurable precision based on application needs

### Performance Characteristics

Some example timings on a moderate system (i7700)

- **ComputeE**: ~28μs (100 digits) to ~1ms (2000 digits)
- **ComputePi**: ~22μs (100 digits) to ~2ms (2000 digits)
- **Exp**: ~36μs (small values) to ~480μs (large values)

## Accuracy

The package provides excellent accuracy:
- **e computation**: Accurate to at least 350 decimal places
- **π computation**: Accurate to at least 700 decimal places  
- **Exponential function**: Relative error < 1e-70 for typical inputs
- **Other functions**: Generally accurate to hundreds of decimal places

## Examples

### Computing Mathematical Constants
```go
// Compute π to 500 decimal places
pi500 := bigmath.ComputePi(2000) // ~2000 bits ≈ 600 decimal places
fmt.Printf("π = %s\n", pi500.Text('f', 500))

// Compute e to 200 decimal places  
e200 := bigmath.ComputeE(800)
fmt.Printf("e = %s\n", e200.Text('f', 200))
```

### High-Precision Calculations
```go
// Compute e^(π*i) + 1 ≈ 0 (Euler's identity, imaginary part omitted)
pi := new(big.Float).Copy(bigmath.Pi)
ePi := bigmath.Exp(pi)
fmt.Printf("e^π = %s\n", ePi.Text('f', 50))

// Compute 2^100 exactly
base := big.NewFloat(2)
exp := int64(100)
result := bigmath.PowInt(base, exp)
fmt.Printf("2^100 = %s\n", result.Text('f', 0))
```

### Gamma Function Applications
```go
// Compute Γ(0.5) = √π
x := big.NewFloat(0.5).SetPrec(256)
gamma := bigmath.Gamma(x)
fmt.Printf("Γ(0.5) = %s\n", gamma.Text('f', 50))

// Compute 100!
factorial100 := bigmath.Factorial(100)
fmt.Printf("100! = %s\n", factorial100.String())
```

## Testing and Benchmarking

The package includes comprehensive tests validating accuracy against known mathematical constants and comparing results with reference implementations.

```bash
go test ./...                    # Run all tests
go test --timeout=10m -bench=.   # Run all benchmarks (may need to adjust the timeout)
go test -run TestComputeE -v     # Test e computation specifically
```
### Tests

Most methods have basic tests along with related unit tests comparing the results to other systems or sources on common sets of inputs.

### Benchmarks

In addition to basic Benchmarks  (e.g., BenchmarkSin), for most methods, there are additional benchmarks that test over a few common ranges of precision bits (53, 64, 128, 256, 500, 1000, 2000) to better measure the impact of increasing precision.  

For many of the methods, I've included more than one common implementation method to better gauge which algorithm is the better choice.  For example, in **log.go**, there are three implementations, **logNewton**, **logHalley**, and **logTaylor**.  


## License

Apache 2.0

## See Also

- [Go math/big documentation](https://pkg.go.dev/math/big)
