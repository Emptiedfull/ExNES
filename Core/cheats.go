package Core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//no comments in this file are ai, the documentation is just assanine to i had to make my own notes

type GameGenieEngine struct {
	cheatTable map[uint16]cheatPart
	Enabled    bool

	cheatDir       string
	cheatDirExists bool
}

func InitCheat() GameGenieEngine {
	gg := GameGenieEngine{
		cheatTable: make(map[uint16]cheatPart),
	}

	gg.cheatDir = "./cheats"
	_, err := os.Stat(gg.cheatDir)
	if err != nil {
		gg.cheatDirExists = false
	} else {
		gg.cheatDirExists = true
	}

	return gg
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

var char6BitMap = [][]rune{
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

var char8BitMap = [][]rune{
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
		'G', 'F', 'E', '%', //char 6
	}, {
		'*', '&', '^', '!', //char 7
	}, {
		'$', '#', '@', '5', // char 8
	},
} // charBitMap[char][bit] = mapped

var AddrString string = "ABCDEFGHIJKLMNO"
var ValString string = "12345678"
var CompareVal string = "!@#$%^&*"

type cheat struct {
	enabled     bool
	Description string

	part cheatPart

	multipart bool
	multi     []cheatPart
}

func (g *GameGenieEngine) decodeCheat(cheatCode string, description string) (cheat, error) {

	if g.cheatTable == nil {
		g.cheatTable = make(map[uint16]cheatPart)
	}

	part, err := parseCode(cheatCode)
	if err != nil {
		return cheat{}, fmt.Errorf("unable to add cheat: %v", err)
	}

	fmt.Println(g.cheatTable)
	return cheat{
		Description: description,
		part:        part,
	}, nil
}

func (gg *GameGenieEngine) addCheat(c cheat) {
	if c.multipart {
		for _, part := range c.multi {
			gg.cheatTable[part.addr] = part
		}
	} else {
		gg.cheatTable[c.part.addr] = c.part
	}
}

func (gg *GameGenieEngine) findCheat(name string) (*[]cheat, bool) {

	if !gg.cheatDirExists {
		return nil, false
	}

	base := strings.TrimPrefix(name, filepath.Ext(name))

	f := filepath.Join(gg.cheatDir, base[0:len(base)-4]+".cht")

	fmt.Println("file:", f)

	if _, err := os.Stat(f); err != nil {
		fmt.Println("candiate not found")
		return nil, false
	}

	fmt.Println("canidate found")

	parseCheatFile(f)

	return nil, true
}

func parseCheatFile(filepath string) (*[]cheat, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	raw := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Println("bad line", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)

		raw[key] = val

	}
	fmt.Println(raw)

	count, err := strconv.Atoi(raw["cheats"])

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	res := make([]cheat, count)
	for i := range count {
		c := cheat{
			enabled: false,
		}
		pre := fmt.Sprintf("cheat%d_", i)

		c.Description = raw[pre+"desc"]

		code := raw[pre+"code"]

		if strings.Contains(code, "+") {
			c.multipart = true
			codeParts := strings.SplitN(code, "+", -1)

			multiArr := make([]cheatPart, len(codeParts))
			for i, part := range codeParts {
				fmt.Println(part)

				cheatpart, err := parseCode(part)
				if err != nil {
					return nil, fmt.Errorf("bad cheat part:%v for fucking up %w", part, err)
				}

				multiArr[i] = cheatpart

			}

			c.multi = multiArr
		} else {
			c.multipart = false

			part, err := parseCode(code)

			if err != nil {
				return nil, fmt.Errorf("bad code: %v for %w", code, err)
			}

			c.part = part

		}

		res = append(res, c)

	}

	return &res, nil

}

type cheatPart struct {
	addr       uint16
	val        uint8
	compare    bool
	compareVal uint8
}

func parseCode(code string) (cheatPart, error) {
	if strings.Contains(code, ":") {
		parts := strings.Split(code, ":")

		if len(parts[0]) != 4 || len(parts[1]) != 2 {
			return cheatPart{}, fmt.Errorf("bad code:", code)
		}

		addrInt, err := strconv.Atoi(parts[0])
		if err != nil {
			return cheatPart{}, fmt.Errorf("invalid code addr:", parts[0])
		}

		valInt, err := strconv.Atoi(parts[1])
		if err != nil {
			return cheatPart{}, fmt.Errorf("invalid val addr:", parts[1])
		}

		return cheatPart{
			addr:    uint16(addrInt),
			val:     uint8(valInt),
			compare: false,
		}, nil
	} else {

		switch len(code) {
		case 6:

			decoded := make(map[rune]int)
			for i, char := range code {
				bits := charBin[char]
				for j := range 4 {

					decoded[char6BitMap[i][j]] = getBit16LSB(bits, j)
				}
			}

			var addr uint16
			n := len(AddrString)
			for i, r := range AddrString {
				addr |= uint16(decoded[r]) << (n - 1 - i)
			}

			var val uint8
			m := len(ValString)
			for i, r := range ValString {
				val |= uint8(decoded[r]) << (m - 1 - i)
			}

			return cheatPart{
				addr:    0x8000 + addr,
				val:     val,
				compare: false,
			}, nil

		case 8:

			decoded := make(map[rune]int)
			for i, char := range code {
				bits := charBin[char]
				for j := range 4 {
					decoded[char8BitMap[i][j]] = getBit16LSB(bits, j)
				}
			}

			var addr uint16
			n := len(AddrString)
			for i, r := range AddrString {
				addr |= uint16(decoded[r]) << (n - 1 - i)
			}

			var val uint8
			m := len(ValString)
			for i, r := range ValString {
				val |= uint8(decoded[r]) << (m - 1 - i)
			}

			var compare uint8
			o := len(CompareVal)
			for i, r := range CompareVal {
				compare |= uint8(decoded[r]) << (o - 1 - i)
			}

			return cheatPart{
				addr:       addr,
				val:        val,
				compare:    true,
				compareVal: compare,
			}, nil
		default:
			return cheatPart{}, fmt.Errorf("bad code:", code)
		}

	}
}
