package main

func (c *console) stepCycles(cycles int) {
	for range cycles {
		c.tick()

		debugConsole.Disassembly[c.Cpu.PC] = DisAssemble(c.Cpu.mem, c.Cpu.PC, c.Cpu.PC)
	}
}

func (c *console) DebugChrRom() {

	Rom := c.Ppu.mem.chrROM

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

				c.Ppu.debugBuffer[offset+(y*8+x)] = colorIndex

			}
		}
	}

}
