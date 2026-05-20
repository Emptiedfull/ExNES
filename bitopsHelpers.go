package main

func builduint16(low uint8, high uint8) uint16 {
	combined := (uint16(high) << 8) | uint16(low)
	return combined
}

func performASL(val uint8) (shifted uint8, carry bool) {
	c := val&0x80 != 0
	val <<= 1
	return val, c
}

func (c *cpu) ADC(val uint8) {
	ACC := uint16(c.A)
	MEM := uint16(val)

	var CAR uint16 = 0
	if c.getFlag(Carry) {
		CAR = 1
	}

	result := ACC + MEM + CAR

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

func performLSR(val uint8) (shifted uint8, carry bool) {
	c := val&0x01 != 0
	val >>= 1
	return val, c
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
