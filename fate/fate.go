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

type FateRoll rune

const (
	PLUS  FateRoll = '+'
	MINUS FateRoll = '-'
	ZERO  FateRoll = 'o'
	NULL  FateRoll = 'x'
)

const (
	CARD_TOP            string = "+-----------+"
	CARD_BOTTOM         string = "+-----------+"
	TOP_MARKER          string = "| %c         |"
	TOP_VALUE_MARKER    string = "| %c      %+d |"
	BOTTOM_MARKER       string = "|         %c |"
	BOTTOM_VALUE_MARKER string = "| %+d      %c |"
)

func RenderCard(rolls []FateRoll) string {
	var value int = SumRolls(rolls)
	var rv []string = []string{
		CARD_TOP,
		fmt.Sprintf(TOP_VALUE_MARKER, rolls[0], value),
		fmt.Sprintf(TOP_MARKER, rolls[1]),
		fmt.Sprintf(TOP_MARKER, rolls[2]),
		fmt.Sprintf(TOP_MARKER, rolls[3]),
		fmt.Sprintf(BOTTOM_MARKER, rolls[3]),
		fmt.Sprintf(BOTTOM_MARKER, rolls[2]),
		fmt.Sprintf(BOTTOM_MARKER, rolls[1]),
		fmt.Sprintf(BOTTOM_VALUE_MARKER, value, rolls[0]),
		CARD_BOTTOM,
	}
	return strings.Join(rv, "\n")
}

func RollDice(count int) []FateRoll {
	rv := make([]FateRoll, count)
	for i, _ := range rv {
		rv[i] = RollDie()
	}
	return rv
}

// RolDie returns a fate dice roll: -, +, or o
func RollDie() FateRoll {
	V, _ := rand.Int(rand.Reader, big.NewInt(6))
	v := V.Int64()
	switch {
	case v < 2:
		return MINUS
	case v < 4:
		return ZERO
	case v < 6:
		return PLUS
	default:
		return NULL
	}
}

// SumRolls returns the value of a list of fate rolls
func SumRolls(rolls []FateRoll) int {
	var value int
	for _, roll := range rolls {
		switch roll {
		case MINUS:
			value += -1
		case PLUS:
			value += 1
		}
	}
	return value
}
