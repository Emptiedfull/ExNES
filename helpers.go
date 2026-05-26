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

func builduint16(low uint8, high uint8) uint16 {
	combined := (uint16(high) << 8) | uint16(low)
	return combined
}

func performASL(val uint8) (shifted uint8, carry bool) {
	c := val&0x80 != 0
	val <<= 1
	return val, c
}

func (c *cpu) ASL(val uint8) (shifted uint8) {

	s, carry := performASL(val)

	c.updateFlag(Carry, carry)
	return s
}

func (c *cpu) ADC(val uint8) {

	var CAR uint8 = 0
	if c.getFlag(Carry) {
		CAR = 1
	}

	ACC := uint16(c.A)
	MEM := uint16(val)

	result := ACC + MEM + uint16(CAR)

	c.updateFlag(Carry, result > 0xFF)

	overflow := ((ACC ^ result) & (MEM ^ result) & 0x80) != 0
	c.updateFlag(oVerflow, overflow)

	c.A = uint8(result & 0xFF)

	c.SetFlagNZ(c.A)
}

func (c *cpu) SBC(val uint8) {
	inverted := ^val

	a := c.A
	carry := uint16(0)
	if c.getFlag(Carry) {
		carry = 1
	}

	sum := uint16(a) + uint16(inverted) + carry

	result8 := uint8(sum)

	c.updateFlag(Carry, sum > 0xFF)

	hasOverflow := ((a ^ result8) & (inverted ^ result8) & 0x80) != 0
	c.updateFlag(oVerflow, hasOverflow)

	c.A = result8
	c.SetFlagNZ(c.A)
}

func (c *cpu) ROR(val uint8) uint8 {
	var oldcarry uint8 = 0
	if c.getFlag(Carry) {
		oldcarry = 1
	}

	newCarry := (val & 0x01) != 0
	result := (val >> 1) | (oldcarry << 7)

	c.updateFlag(Carry, newCarry)
	c.SetFlagNZ(result)

	return result
}

func (c *cpu) COMPARE(A, B uint8) uint8 {
	result := A - B

	c.updateFlag(Carry, A >= B)
	c.updateFlag(Zero, A == B)
	c.updateFlag(Negative, (result&0x80) != 0)

	return result

}

func performROL(val uint8, carry bool) (uint8, bool) {
	var oldcarry uint8 = 0
	if carry {
		oldcarry = 1
	}
	newCarry := (val & 0x80) != 0
	res := (val << 1) | oldcarry
	return res, newCarry
}

func (c *cpu) ROL(val uint8) uint8 {
	r, co := performROL(val, c.getFlag(Carry))
	c.updateFlag(Carry, co)
	return r
}

func performLSR(val uint8) (shifted uint8, carry bool) {
	c := val&0x01 != 0
	val >>= 1
	return val, c
}

func (c *cpu) LSR(val uint8) uint8 {
	r, co := performLSR(val)
	c.updateFlag(Carry, co)
	return r
}

func getbit(val uint8, pos int8) uint8 {
	return (val >> pos) & 1
}

func getbitBool(val uint8, pos int8) bool {
	return getbit(val, pos) == 1
}

func AssignBit(val uint8, pos uint8, state bool) uint8 {
	if state {
		return val | (1 << pos)
	}
	return val &^ (1 << pos)
}
