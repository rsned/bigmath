// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigmath

import (
	"math"
	"math/big"
)

// stirlingApproximation computes n! using Stirling's approximation
// n! ≈ sqrt(2πn) * (n/e)^n
func stirlingApproximation(x *big.Float) *big.Float {
	xFloat, _ := x.Float64()

	if xFloat <= 0 {
		return big.NewFloat(1)
	}

	// sqrt(2πn)
	twoPiN := 2 * math.Pi * xFloat
	sqrtTwoPiN := big.NewFloat(math.Sqrt(twoPiN))

	// (n/e)^n = n^n / e^n
	nToN := Pow(x, x)
	eToN := Pow(big.NewFloat(math.E), x)
	nOverEToN := new(big.Float).Quo(nToN, eToN)

	// Combine: sqrt(2πn) * (n/e)^n
	result := new(big.Float).Mul(sqrtTwoPiN, nOverEToN)

	return result
}
