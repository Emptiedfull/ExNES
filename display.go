package main

type Pallete [][4]byte

var NesPaletteLUT = [64]RGB{
	{0x7C, 0x7C, 0x7C}, {0x00, 0x00, 0xFC}, {0x00, 0x00, 0xBC}, {0x44, 0x28, 0xBC},
	{0x94, 0x00, 0x84}, {0xA8, 0x00, 0x20}, {0xA8, 0x10, 0x00}, {0x88, 0x14, 0x00},
	{0x50, 0x30, 0x00}, {0x00, 0x78, 0x00}, {0x00, 0x68, 0x00}, {0x00, 0x58, 0x00},
	{0x00, 0x40, 0x58}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},

	{0xBC, 0xBC, 0xBC}, {0x00, 0x78, 0xF8}, {0x00, 0x58, 0xF8}, {0x68, 0x44, 0xFC},
	{0xD8, 0x00, 0xCC}, {0xE4, 0x00, 0x58}, {0xF8, 0x38, 0x00}, {0xE4, 0x5C, 0x10},
	{0xAC, 0x7C, 0x00}, {0x00, 0xB8, 0x00}, {0x00, 0xA8, 0x00}, {0x00, 0xA8, 0x44},
	{0x00, 0x88, 0x88}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},

	{0xF8, 0x78, 0xF8}, {0x3C, 0xBC, 0xFC}, {0x68, 0xA0, 0xFC}, {0x98, 0x88, 0xF8},
	{0xF8, 0x78, 0xF8}, {0xFD, 0x54, 0x8C}, {0xF8, 0x78, 0x58}, {0xFC, 0xA0, 0x44},
	{0xFF, 0xBC, 0x3C}, {0xB4, 0xD4, 0x20}, {0x48, 0xCD, 0x48}, {0x58, 0xFA, 0x98},
	{0x58, 0xF8, 0xFC}, {0x78, 0x78, 0x78}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},

	{0xFC, 0xFC, 0xFC}, {0xA4, 0xE4, 0xFC}, {0xB8, 0xB8, 0xF8}, {0xD8, 0xB8, 0xF8},
	{0xF8, 0xB8, 0xF8}, {0xF8, 0xA4, 0xC0}, {0xF0, 0xD0, 0xB0}, {0xFC, 0xE0, 0xA8},
	{0xFD, 0xEA, 0x88}, {0xD8, 0xF8, 0x78}, {0xB8, 0xF8, 0xB8}, {0xB8, 0xF8, 0xD8},
	{0x00, 0xFC, 0xFC}, {0xD8, 0xD8, 0xD8}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},
}

type RGB struct {
	R, G, B uint8
}

func (p *ppu) getColor(colorindex uint8) RGB {
	switch colorindex {
	case 1:
		return RGB{R: 184, G: 56, B: 140}
	case 2:
		return RGB{R: 60, G: 188, B: 252}
	case 3:
		return RGB{R: 252, G: 252, B: 252}
	default:
		return RGB{R: 16, G: 44, B: 52}
	}
}

func (p *ppu) renderScanline() {
	var PatternOffset uint16 = 0
	if p.mem.register.BgPattern {
		PatternOffset = 4096
	}

	screenY := p.Scanline
	tileRow := screenY / 8
	pixelYWithinTile := screenY % 8

	var baseAddr uint16 = 0x2000

	for col := range 32 {
		ntIdx := uint16(tileRow*32 + col)
		tileID := uint16(p.read(baseAddr + ntIdx))

		attrRow := tileRow / 4
		attrCol := col / 4
		attrAddr := baseAddr + 960 + uint16(attrRow*8+attrCol)
		attrByte := p.read(attrAddr)

		qIDx := (col % 4) / 2
		qIDy := (tileRow % 4) / 2

		shift := uint8((qIDy * 4) + (qIDx * 2))
		block := (attrByte >> shift) & 0x03

		chrOffset := (tileID * 16) + PatternOffset + uint16(pixelYWithinTile)

		if int(chrOffset+8) <= len(p.mem.chrROM) {
			low := p.mem.chrROM[chrOffset]
			high := p.mem.chrROM[chrOffset+8]

			for x := range 8 {
				bit1 := (low >> (7 - x)) & 0x01
				bit2 := (high >> (7 - x)) & 0x01

				colorIndex := (bit2 << 1) | bit1

				ramAddr := p.getPalleteIDX(int(block), colorIndex)
				colorToken := p.read(ramAddr) & 0x3F

				rgb := NesPaletteLUT[colorToken]

				screenX := (col * 8) + x

				bufferIDX := (int(screenY)*256 + screenX) * 4
				pushcolor(rgb, p.backBuffer, bufferIDX)
			}

		}

	}

	var spriteOffset uint16 = 0

	if p.mem.register.SpritePattern {
		spriteOffset = 4096
	}

	for i := 63; i >= 0; i-- {
		oamIDX := i * 4

		spriteY := int(p.mem.oamData[oamIDX]) + 1

		if screenY >= spriteY && screenY < spriteY+8 {
			tileID := uint16(p.mem.oamData[oamIDX+1])

			attr := p.mem.oamData[oamIDX+2]
			palleteblock := 4 + attr&0x03
			// priority := getbitBool(attr, 5)
			flipX := getbitBool(attr, 6)
			flipY := getbitBool(attr, 7)

			spriteX := int(p.mem.oamData[oamIDX+3])

			rowOffset := screenY - spriteY
			if flipY {
				rowOffset = 7 - rowOffset
			}

			chrOffset := (tileID * 16) + spriteOffset + uint16(rowOffset)

			if int(chrOffset+8) <= len(p.mem.chrROM) {
				low := p.mem.chrROM[chrOffset]
				high := p.mem.chrROM[chrOffset+8]

				for x := range 8 {
					pixelX := spriteX + x

					if pixelX > 256 {
						continue
					}

					bitshift := 7 - x
					if flipX {
						bitshift = x
					}

					bit1 := (low >> uint8(bitshift)) & 0x01
					bit2 := (high >> uint8(bitshift)) & 0x01

					colorindex := (bit2 << 1) | bit1

					if colorindex == 0 {
						continue
					}

					addr := p.getPalleteIDX(int(palleteblock), colorindex)
					colorToken := p.read(addr) & 0x3F
					rgb := NesPaletteLUT[colorToken]
					bufferidx := (screenY*256 + pixelX) * 4

					pushcolor(rgb, p.backBuffer, bufferidx)

				}
			}
		}
	}
}

func (p *ppu) getPalleteIDX(palleteBlock int, colorindex uint8) uint16 {
	base := uint16(0x3F00)

	if palleteBlock >= 4 {
		base = uint16(0x3F10)
	}

	addr := base + (uint16(palleteBlock) * 4) + uint16(colorindex)

	if colorindex == 0 {
		return 0x3F00
	}

	return addr
}

func pushcolor(color RGB, buffer []uint8, start int) {
	buffer[start] = color.R
	buffer[start+1] = color.G
	buffer[start+2] = color.B
	buffer[start+3] = 255
}
