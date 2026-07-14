package Core

import "fmt"

//no comments in this file are ai, the documentation is just assanine to i had to make my own notes

type GameGenieEngine struct {
	cheatTable map[uint16]cheat
	Enabled    bool
}

// A = 0000
// P = 0001
// Z = 0010
// L = 0011
// G = 0100
// I = 0101
// T = 0110
// Y = 0111
// E = 1000
// O = 1001
// X = 1010
// U = 1011
// K = 1100
// S = 1101
// V = 1110
// N = 1111

var charBin = map[rune]uint16{
	'A': 0x0,
	'P': 0x1,
	'Z': 0x2,
	'L': 0x3,
	'G': 0x4,
	'I': 0x5,
	'T': 0x6,
	'Y': 0x7,
	'E': 0x8,
	'O': 0x9,
	'X': 0xA,
	'U': 0xB,
	'K': 0xC,
	'S': 0xD,
	'V': 0xE,
	'N': 0xF,
}

// Char # |   1   |   2   |   3   |   4   |   5   |   6   |
// Bit  # |3|2|1|0|3|2|1|0|3|2|1|0|3|2|1|0|3|2|1|0|3|2|1|0|
// maps to|1|6|7|8|H|2|3|4|-|I|J|K|L|A|B|C|D|M|N|O|5|E|F|G|

// note char 3 bit 3 is used by the game genie to determine the length
// of the code.

// The value is made of 12345678 of the maps to line.
// The address is made of ABCDEFGHIJKLMNO of the maps to line.

// if I take code SXIOPO(infinite lives in smb1) I get this

// |   S   |   X   |   I   |   O   |   P   |   O   |
// |1|1|0|1|1|0|1|0|0|1|0|1|1|0|0|1|0|0|0|1|1|0|0|1|

var charBitMap = [][]rune{
	{
		'8', '7', '6', '1', //char 1
	}, {
		'4', '3', '2', 'H', //char 2
	}, {
		'K', 'J', 'I', '_', //char 3
	}, {
		'C', 'B', 'A', 'L', //char 4
	}, {
		'O', 'N', 'M', 'D', //chat 5
	}, {
		'G', 'F', 'E', '5', //char 6
	},
} // charBitMap[char][bit] = mapped

var AddrString6 string = "ABCDEFGHIJKLMNO"
var ValStrin6 string = "12345678"

type cheat struct {
	enabled    bool
	addr       uint16
	val        uint8
	compare    bool
	compareVal uint8
}

func (g *GameGenieEngine) AddCheat(cheatCode string) {

	if g.cheatTable == nil {
		g.cheatTable = make(map[uint16]cheat)
	}

	cheat := DecodeCheat(cheatCode)
	g.cheatTable[cheat.addr] = cheat
}

func DecodeCheat(code string) cheat {
	c := cheat{}
	switch len(code) {
	case 6:
		c.compare = false
		c.enabled = true

		decoded := make(map[rune]int)
		for i, char := range code {
			bits := charBin[char]
			for j := range 4 {
				fmt.Println(charBitMap[i][j], ":", getBit16LSB(bits, j))
				decoded[charBitMap[i][j]] = getBit16LSB(bits, j)
			}
		}

		var addr uint16
		n := len(AddrString6)
		for i, r := range AddrString6 {
			addr |= uint16(decoded[r]) << (n - 1 - i)
		}
		c.addr = 0x8000 + addr

		var val uint8
		m := len(ValStrin6)
		for i, r := range ValStrin6 {
			val |= uint8(decoded[r]) << (m - 1 - i)
		}
		c.val = val

		c.val = val

	}
	return c
}

func getBit16LSB(val uint16, pos int) int {
	return int((val >> pos) & 1)
}
