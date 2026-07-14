package Core

import "fmt"

//no comments in this file are ai, the documentation is just assanine to i had to make my own notes

type GameGenieEngine struct {
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
	'A': 0x0000,
	'P': 0x1000,
	'Z': 0x0010,
	'L': 0x0011,
	'G': 0x0100,
	'I': 0x0101,
	'T': 0x0110,
	'Y': 0x0111,
	'E': 0x1000,
	'O': 0x1001,
	'X': 0x1010,
	'U': 0x1011,
	'K': 0x1100,
	'S': 0x1101,
	'V': 0x1110,
	'N': 0x1111,
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

var charBitMap = [][]string{
	{
		"8", "7", "6", "1", //char 1
	}, {
		"4", "3", "2", "H", //char 2
	}, {
		"K", "J", "I", "_", //char 3
	}, {
		"C", "B", "A", "L", //char 4
	}, {
		"O", "N", "M", "D", //chat 5
	}, {
		"G", "F", "E", "5", //char 6
	},
} // charBitMap[char][bit] = mapped

type cheat struct {
	enabled    bool
	addr       uint16
	val        uint8
	compare    bool
	compareVal uint8
}

func DecodeCheat(code string) cheat {
	c := cheat{}
	switch len(code) {
	case 6:
		for i, char := range code {
			bits := charBin[char]
			for j := range 4 {
				fmt.Println(charBitMap[i][j], ":", getBit16LSB(bits, j))
			}
		}
	}
	return c
}

func getBit16LSB(val uint16, pos int) int {
	return int((val >> pos) & 1)
}
