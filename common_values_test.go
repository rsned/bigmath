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
	"fmt"
	"math/big"
	"strings"
	"testing"
)

const (
	// Known value of e to ~600 decimal places
	eKnown1000 = "2.7182818284590452353602874713526624977572470936999595749669676277240766303535475945713821785251664274274663919320030599218174135966290435729003342952605956307381323286279434907632338298807531952510190115738341879307021540891499348841675092447614606680822648001684774118537423454424371075390777449920695517027618386062613313845830007520449338265602976067371132007093287091274437470472306969772093101416928368190255151086574637721112523897844250569536967707854499699679468644549059879316368892300987931277361782154249992295763514822082698951936680331825288693984964651058209392398294887933203625094649524825692302"

	// Known value of π to ~1950 decimal places
	piKnown1000 = "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989380952572010654858632788659361533818279682303019520353018529689957736225994138912497217752834791315155748572424541506959508295331168617278558890750983817546374649393192550604009277016711390098488240128583616035637076601047101819429555961989467678374494482553797747268471040475346462080466842590694912933136770289891521047521620569660240580381501935112533824300355876402474964732639141992726042699227967823547816360093417216412199245863150302861829745557067498385054945885869269956909272107975093029553211653449872027559602364806654991198818347977535663698074265425278625518184175746728909777727938000816470600161452491921732172147723501414419735685481613611573525521334757418494684385233239073941433345477624168625189835694855620992192221842725502542568876717904946016746097659798123254675172727753314097692073032012871893975476646821197846862177359068200808439951846244775893033264701926159139988889100011752300816000684151655644072461608423317090816080"
)

func TestComputeE(t *testing.T) {
	tests := []struct {
		name      string
		precision uint
		digits    int
		tolerance float64
	}{
		{"100 digits", 400, 100, 1e-95},
		{"500 digits", 2000, 500, 1e-495},
		{"1000 digits", 4000, min(len(eKnown1000)-2, 1000), 1e-995}, // Available digits
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Compute e using our function
			computed := ComputeE(test.precision)

			// Parse known value with same precision - ensure we don't exceed available digits
			availableDigits := min(test.digits, len(eKnown1000)-2) // -2 for "2."
			expected := new(big.Float).SetPrec(test.precision)
			expected.SetString(eKnown1000[:availableDigits+2]) // +2 for "2."

			// Calculate relative error
			diff := new(big.Float).SetPrec(test.precision)
			diff.Sub(computed, expected)
			diff.Abs(diff)

			relativeError := new(big.Float).SetPrec(test.precision)
			relativeError.Quo(diff, expected)
			relativeErrorFloat, _ := relativeError.Float64()

			if relativeErrorFloat > test.tolerance {
				t.Errorf("ComputeE(%d) relative error %.2e exceeds tolerance %.2e",
					test.precision, relativeErrorFloat, test.tolerance)
				t.Errorf("Expected: %s", expected.Text('f', min(availableDigits, 50)))
				t.Errorf("Got:      %s", computed.Text('f', min(availableDigits, 50)))
			}
		})
	}
}

func TestComputePi(t *testing.T) {
	tests := []struct {
		name      string
		precision uint
		digits    int
		tolerance float64
	}{
		{"100 digits", 400, 100, 1e-95},
		{"500 digits", 2000, 500, 1e-495},
		{"1000 digits", 4000, min(len(piKnown1000)-2, 1000), 1e-995}, // Available digits
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Compute π using our function
			computed := ComputePi(test.precision)

			// Parse known value with same precision - ensure we don't exceed available digits
			availableDigits := min(test.digits, len(piKnown1000)-2) // -2 for "3."
			expected := new(big.Float).SetPrec(test.precision)
			expected.SetString(piKnown1000[:availableDigits+2]) // +2 for "3."

			// Calculate relative error
			diff := new(big.Float).SetPrec(test.precision)
			diff.Sub(computed, expected)
			diff.Abs(diff)

			relativeError := new(big.Float).SetPrec(test.precision)
			relativeError.Quo(diff, expected)
			relativeErrorFloat, _ := relativeError.Float64()

			if relativeErrorFloat > test.tolerance {
				t.Errorf("ComputePi(%d) relative error %.2e exceeds tolerance %.2e",
					test.precision, relativeErrorFloat, test.tolerance)
				t.Errorf("Expected: %s", expected.Text('f', min(availableDigits, 50)))
				t.Errorf("Got:      %s", computed.Text('f', min(availableDigits, 50)))
			}
		})
	}
}

func TestComputeEDigitByDigit(t *testing.T) {
	// Test individual digits of e for accuracy
	computed := ComputeE(4000) // High precision
	computedStr := computed.Text('f', 1000)

	// Remove the "2." prefix and check digit by digit
	computedDigits := strings.Replace(computedStr, "2.", "", 1)
	knownDigits := strings.Replace(eKnown1000, "2.", "", 1)

	// Find first differing digit
	minLen := min(len(computedDigits), len(knownDigits))
	for i := range minLen {
		if computedDigits[i] != knownDigits[i] {
			if i < 300 { // Should be accurate to at least 300 digits
				t.Errorf("ComputeE differs from known value at digit %d: got %c, expected %c",
					i+1, computedDigits[i], knownDigits[i])
				t.Errorf("Computed: 2.%s...", computedDigits[:min(i+10, len(computedDigits))])
				t.Errorf("Expected: 2.%s...", knownDigits[:min(i+10, len(knownDigits))])
			} else {
				t.Logf("ComputeE matches known value for first %d digits", i)
			}

			break
		}
	}
}

func TestComputePiDigitByDigit(t *testing.T) {
	// Test individual digits of π for accuracy
	computed := ComputePi(4000) // High precision
	computedStr := computed.Text('f', 1000)

	// Remove the "3." prefix and check digit by digit
	computedDigits := strings.Replace(computedStr, "3.", "", 1)
	knownDigits := strings.Replace(piKnown1000, "3.", "", 1)

	// Find first differing digit
	minLen := min(len(computedDigits), len(knownDigits))
	for i := range minLen {
		if computedDigits[i] != knownDigits[i] {
			if i < 300 { // Should be accurate to at least 300 digits
				t.Errorf("ComputePi differs from known value at digit %d: got %c, expected %c",
					i+1, computedDigits[i], knownDigits[i])
				t.Errorf("Computed: 3.%s...", computedDigits[:min(i+10, len(computedDigits))])
				t.Errorf("Expected: 3.%s...", knownDigits[:min(i+10, len(knownDigits))])
			} else {
				t.Logf("ComputePi matches known value for first %d digits", i)
			}

			break
		}
	}
}

func TestPrecomputedConstants(t *testing.T) {
	// Test that the precomputed Pi and E constants are accurate

	// Test Pi
	expected := new(big.Float).SetPrec(4000)
	expected.SetString(piKnown1000[:1002]) // First 1000 decimal places + "3."

	diff := new(big.Float).SetPrec(4000)
	diff.Sub(bigPi, expected)
	diff.Abs(diff)

	relativeError := new(big.Float).SetPrec(4000)
	relativeError.Quo(diff, expected)
	relativeErrorFloat, _ := relativeError.Float64()

	if relativeErrorFloat > 1e-95 {
		t.Errorf("Precomputed Pi relative error %.2e exceeds tolerance %.2e",
			relativeErrorFloat, 1e-95)
	}

	// Test E
	expected.SetString(eKnown1000[:502]) // First 500 decimal places + "2."

	diff.Sub(bigE, expected)
	diff.Abs(diff)

	relativeError.Quo(diff, expected)
	relativeErrorFloat, _ = relativeError.Float64()

	if relativeErrorFloat > 1e-95 {
		t.Errorf("Precomputed E relative error %.2e exceeds tolerance %.2e",
			relativeErrorFloat, 1e-95)
	}
}

// Benchmarks for performance testing
func BenchmarkComputeE(b *testing.B) {
	for _, prec := range precisions {
		b.Run(fmt.Sprintf("precision_%d", prec), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = ComputeE(prec)
			}
		})
	}
}

func BenchmarkComputePi(b *testing.B) {
	for _, prec := range precisions {
		b.Run(fmt.Sprintf("precision_%d", prec), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = ComputePi(prec)
			}
		})
	}
}
