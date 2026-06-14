package Core

import "fmt"

func (c *console) DebugChrRom() {

	Rom := c.Ppu.mem.mapper.extractCHR()

	tiles := min(len(Rom)/16, 512)

	for t := range tiles {
		tileBytes := Rom[(t * 16):((t + 1) * 16)]
		offset := t * 64

		for y := range 8 {
			low := tileBytes[y]
			high := tileBytes[y+8]

			for x := range 8 {
				bit1 := (low >> (7 - x)) & 0x01
				bit2 := (high >> (7 - x)) & 0x01

				colorIndex := (bit2 << 1) | bit1

				c.Ppu.DebugBuffer[offset+(y*8+x)] = colorIndex

			}
		}
	}

}

func (c *console) DebugNameTable() {
	fmt.Println("debugging namtable")
	Rom := c.Ppu.mem.mapper.extractCHR()
	c.Ppu.DebugBuffer = make([]uint8, 512*480)

	var PatternOffset uint16 = 0
	if c.Ppu.mem.register.BgPattern {
		PatternOffset = 4096
	}

	for quadrant := 0; quadrant < 4; quadrant++ {
		var baseAddr uint16 = 0x2000 + uint16(quadrant)*0x400

		quadXOffset := (quadrant % 2) * 256
		quadYOffset := (quadrant / 2) * 240

		for row := 0; row < 30; row++ {
			for col := 0; col < 32; col++ {

				ntIdx := uint16(row*32 + col)

				tileID := uint16(c.Ppu.read(baseAddr + ntIdx))

				chrOffset := (tileID * 16) + PatternOffset

				if int(chrOffset+16) <= len(Rom) {
					tileBytes := Rom[chrOffset : chrOffset+16]

					for y := 0; y < 8; y++ {
						low := tileBytes[y]
						high := tileBytes[y+8]

						for x := 0; x < 8; x++ {
							bit1 := (low >> (7 - x)) & 0x01
							bit2 := (high >> (7 - x)) & 0x01
							colorIndex := (bit2 << 1) | bit1

							pixelX := quadXOffset + (col * 8) + x
							pixelY := quadYOffset + (row * 8) + y

							c.Ppu.DebugBuffer[pixelY*512+pixelX] = colorIndex
						}
					}
				}
			}
		}
	}
}
