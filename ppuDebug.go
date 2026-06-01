package main

func (c *console) stepCycles(cycles int) {
	for range cycles {
		c.tick()
	}
}

func (c *console) tick() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.totalCycles++
	} else {
		c.Cpu.tick()
	}

	c.Ppu.Tick()
	c.Ppu.Tick()
	c.Ppu.Tick()

}
