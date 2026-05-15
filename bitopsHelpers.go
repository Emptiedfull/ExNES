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
