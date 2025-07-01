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

// Bernoulli B2n numbers are used in some Taylor series expansions.
var (
	BernoulliNumbers = []big.Float{
		*big.NewFloat(1),                                                         // B(0)
		*big.NewFloat(1.0 / 2),                                                   // B(1)
		*big.NewFloat(1.0 / 6),                                                   // B(2)
		*big.NewFloat(-1.0 / 30),                                                 // B(4)
		*big.NewFloat(1.0 / 42),                                                  // B(6)
		*big.NewFloat(-1.0 / 30),                                                 // B(8)
		*big.NewFloat(5.0 / 66),                                                  // B(10)
		*big.NewFloat(-691.0 / 2730),                                             // B(12)
		*big.NewFloat(7.0 / 6),                                                   // B(14)
		*big.NewFloat(-3617.0 / 510),                                             // B(16)
		*big.NewFloat(43867.0 / 798),                                             // B(18)
		*big.NewFloat(-174611.0 / 330),                                           // B(20)
		*big.NewFloat(854513.0 / 138),                                            // B(22)
		*big.NewFloat(-236364091.0 / 2730),                                       // B(24)
		*big.NewFloat(8553103.0 / 6),                                             // B(26)
		*big.NewFloat(-23749461029.0 / 870),                                      // B(28)
		*big.NewFloat(8615841276005.0 / 14322),                                   // B(30)
		*big.NewFloat(-7709321041217.0 / 510),                                    // B(32)
		*big.NewFloat(2577687858367.0 / 6),                                       // B(34)
		*big.NewFloat(-26315271553053477373.0 / 1919190),                         // B(36)
		*big.NewFloat(2929993913841559.0 / 6),                                    // B(38)
		*big.NewFloat(-261082718496449122051.0 / 13530),                          // B(40)
		*big.NewFloat(1520097643918070802691.0 / 1806),                           // B(42)
		*big.NewFloat(-27833269579301024235023.0 / 690),                          // B(44)
		*big.NewFloat(596451111593912163277961.0 / 282),                          // B(46)
		*big.NewFloat(-5609403368997817686249127547.0 / 46410),                   // B(48)
		*big.NewFloat(495057205241079648212477525.0 / 66),                        // B(50)
		*big.NewFloat(-801165718135489957347924991853.0 / 1590),                  // B(52)
		*big.NewFloat(29149963634884862421418123812691.0 / 798),                  // B(54)
		*big.NewFloat(-2479392929313226753685415739663229.0 / 870),               // B(56)
		*big.NewFloat(84483613348880041862046775994036021.0 / 354),               // B(58)
		*big.NewFloat(-1215233140483755572040304994079820246041491.0 / 56786730), // B(60)
	}
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
		term.Quo(new(big.Float).SetInt64(1), factorial)
		e.Add(e, term)

		// Check if term is small enough to stop
		if term.Cmp(new(big.Float).SetPrec(precision).Quo(
			new(big.Float).SetInt64(1),
			new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision/4)), nil)),
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
