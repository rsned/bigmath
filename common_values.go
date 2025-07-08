package bigmath

import "math/big"

// Predefine some example values for pi and e to use.
var (
	bigPi        *big.Float
	bigHalfPi    *big.Float
	bigQuarterPi *big.Float

	bigE *big.Float
)

// Create some of the big.Floats we use a lot.
var (
	zero  = new(big.Float).SetInt64(0)
	one   = new(big.Float).SetInt64(1)
	two   = new(big.Float).SetInt64(2)
	four  = new(big.Float).SetInt64(4)
	five  = new(big.Float).SetInt64(5)
	six   = new(big.Float).SetInt64(6)
	seven = new(big.Float).SetInt64(7)
	eight = new(big.Float).SetInt64(8)
	nine  = new(big.Float).SetInt64(9)
	ten   = new(big.Float).SetInt64(10)

	intZero  = new(big.Int).SetInt64(0)
	intOne   = new(big.Int).SetInt64(1)
	intTwo   = new(big.Int).SetInt64(2)
	intFour  = new(big.Int).SetInt64(4)
	intFive  = new(big.Int).SetInt64(5)
	intSix   = new(big.Int).SetInt64(6)
	intSeven = new(big.Int).SetInt64(7)
	intEight = new(big.Int).SetInt64(8)
	intNine  = new(big.Int).SetInt64(9)
	intTen   = new(big.Int).SetInt64(10)
)

func init() {
	bigPi = ComputePi(1000)
	bigE = ComputeE(1000)

	bigQuarterPi = new(big.Float).SetPrec(1000).Quo(bigPi, four)
	bigHalfPi = new(big.Float).SetPrec(1000).Quo(bigPi, two)
}

// Pi returns a copy of the precomputed high-precision π value.
//
// Because we can't make const *big.Floats, and don't want someone
// downstream altering the value, we return a copy.
func Pi() *big.Float {
	return new(big.Float).Copy(bigPi)
}

// E returns a copy of the precomputed high-precision e value.
//
// Because we can't make const *big.Floats, and don't want someone
// downstream altering the value, we return a copy.
func E() *big.Float {
	return new(big.Float).Copy(bigE)
}

// ComputeE calculates e with the given precision using series expansion
func ComputeE(precision uint) *big.Float {
	e := new(big.Float).SetPrec(precision)
	term := new(big.Float).SetPrec(precision).SetInt64(1)
	factorial := new(big.Float).SetPrec(precision).SetInt64(1)

	e.SetInt64(1) // Start with 1

	// e = 1 + 1/1! + 1/2! + 1/3! + ...
	for i := int64(1); i < 200; i++ { // 200 terms should be enough for 1000 digits
		factorial.Mul(factorial, new(big.Float).SetInt64(i))
		term.Quo(one, factorial)
		e.Add(e, term)

		// Check if term is small enough to stop
		if term.Cmp(new(big.Float).SetPrec(precision).Quo(one,
			new(big.Float).SetInt(new(big.Int).Exp(intTen, big.NewInt(int64(precision/4)), nil)),
		)) < 0 {
			break
		}
	}

	return e
}

// ComputePi calculates π with the given precision using Machin's formula
func ComputePi(precision uint) *big.Float {
	pi := new(big.Float).SetPrec(precision)

	// Use Machin's formula: π/4 = 4*arctan(1/5) - arctan(1/239)
	one := new(big.Float).SetPrec(precision).SetInt64(1)
	four := new(big.Float).SetPrec(precision).SetInt64(4)
	five := new(big.Float).SetPrec(precision).SetInt64(5)
	val239 := new(big.Float).SetPrec(precision).SetInt64(239)

	// arctan(1/5)
	arctan1_5 := Atan(new(big.Float).SetPrec(precision).Quo(one, five))

	// arctan(1/239)
	arctan1_239 := Atan(new(big.Float).SetPrec(precision).Quo(one, val239))

	// π/4 = 4*arctan(1/5) - arctan(1/239)
	pi.Mul(four, arctan1_5)
	pi.Sub(pi, arctan1_239)
	pi.Mul(pi, four) // Convert π/4 to π

	return pi
}
