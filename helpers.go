package main

func (c *cpu) SetFlagNZ(val uint8) {
	if val == 0 {
		c.P |= Zero
	} else {
		c.P &^= Zero
	}

	if val&0x80 != 0 {
		c.P |= Negative
	} else {
		c.P &^= Negative
	}
}

func (c *cpu) pushStack(val uint8) {
	stackAdrr := 0x0100 + uint16(c.S)
	c.mem.Write(stackAdrr, val)
	c.S--
}

func (c *cpu) popStack() uint8 {
	c.S++
	stackAddr := 0x0100 + uint16(c.S)
	return c.mem.Read(stackAddr)
}

func (c *cpu) clearFlag(flag flags) {
	c.P &^= flag
}

func (c *cpu) setFlag(flag flags) {
	c.P |= flag
}

func (c *cpu) getFlag(flag flags) bool {
	return (c.P & flag) != 0
}

func (c *cpu) updateFlag(flag flags, setting bool) {
	if setting {
		c.setFlag(flag)
	} else {
		c.clearFlag(flag)
	}
}

func (c *cpu) fetchone() uint8 {
	val := c.mem.Read(c.PC)
	c.PC++
	return uint8(val)
}

func crossedPage(prev uint16, next uint16) bool {
	return (prev & 0xFF00) != (next & 0xFF00)
}
