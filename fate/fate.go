// Copyright © 2017 Walter Scheper <walter.scheper@gmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fate

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Roll is a single roll of a fate dice.
type Roll rune

const (
	// Plus is a + roll
	Plus Roll = '+'
	// Minus is a - roll
	Minus Roll = '-'
	// Zero is a o roll
	Zero Roll = 'o'
	// Unknown is a roll that we shouldn't get
	Unknown Roll = 'x'
)

const (
	cardTop           string = "+-----------+"
	cardBottom        string = "+-----------+"
	topMarker         string = "| %c         |"
	topValueMarker    string = "| %c      %+d |"
	bottomMarker      string = "|         %c |"
	bottomValueMarker string = "| %+d      %c |"
)

// RenderFunc is type of function that takes a slice of Rolls and converts them
// into a string for printing.
type RenderFunc func([]Roll) string

// RenderCard takes a slice of rolls and returns a string representation of a Fate card.
func RenderCard(rolls []Roll) string {
	var value = SumRolls(rolls)
	var rv = []string{
		cardTop,
		fmt.Sprintf(topValueMarker, rolls[0], value),
		fmt.Sprintf(topMarker, rolls[1]),
		fmt.Sprintf(topMarker, rolls[2]),
		fmt.Sprintf(topMarker, rolls[3]),
		fmt.Sprintf(bottomMarker, rolls[3]),
		fmt.Sprintf(bottomMarker, rolls[2]),
		fmt.Sprintf(bottomMarker, rolls[1]),
		fmt.Sprintf(bottomValueMarker, value, rolls[0]),
		cardBottom,
	}
	return strings.Join(rv, "\n")
}

// RenderDice takes a slice of Rolls and returns a string representation of an
// equivalent number of Fate dice.
func RenderDice(rolls []Roll) string {
	var value = SumRolls(rolls)
	return fmt.Sprintf("%+d: %c %c %c %c", value, rolls[0], rolls[1], rolls[2], rolls[3])
}

// RollDice returns a slice of count Rolls.
func RollDice(count int) []Roll {
	rv := make([]Roll, count)
	for i := range rv {
		rv[i] = RollDie()
	}
	return rv
}

// RollDie returns a Roll: -, +, or o
func RollDie() Roll {
	V, _ := rand.Int(rand.Reader, big.NewInt(6))
	v := V.Int64()
	switch {
	case v < 2:
		return Minus
	case v < 4:
		return Zero
	case v < 6:
		return Plus
	default:
		return Unknown
	}
}

// SumRolls returns the value of a slice of Rolls
func SumRolls(rolls []Roll) int {
	var value int
	for _, roll := range rolls {
		switch roll {
		case Minus:
			value = value - 1
		case Plus:
			value = value + 1
		}
	}
	return value
}
