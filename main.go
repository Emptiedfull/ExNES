package main

type console struct {
	Cpu        cpu
	Ppu        ppu
	OpenBusVal uint8
}

func (c *console) tick() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.totalCycles++
	} else {
		c.Cpu.tick()
	}

}
