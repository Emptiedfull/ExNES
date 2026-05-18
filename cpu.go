package main

import "fmt"

type flags = uint8

const (
	Carry     flags = 1 << 0
	Zero            = 1 << 1
	Interrupt       = 1 << 2
	Decimal         = 1 << 3
	Break           = 1 << 4
	Unused          = 1 << 5
	oVerflow        = 1 << 6
	Negative        = 1 << 7
)

type cpu struct {
	PC uint16 //program counter
	S  uint8  //stack pointer
	P  flags  //processor status
	A  uint8  //Accumulator
	X  uint8
	Y  uint8

	mem *bus
	temp

	currentOp   uint8
	currentstep int
	totalCycles int
}

func (c *cpu) tick() {
	c.totalCycles++
	if c.currentstep == 0 {
		c.currentOp = c.fetchone()
		c.totalCycles++
	}

	opcode := FetchTable[c.currentOp]

	finished := opcode.Execute(c, c.currentstep)
	fmt.Println(opcode.Name, c.PC, c.totalCycles)
	if finished {
		c.currentstep = 0
		c.temp = temp{}
	} else {
		c.currentstep++
	}

}

type temp struct {
	high uint8
	low  uint8

	pointer uint8
	addr    uint16
	val     uint8
}

type bus struct {
	internal [2048]byte
	external []byte
}

func (b *bus) Read(addr uint16) uint8 {
	if addr <= 0x1FFF {
		// RAM Mirroring: $0000-$07FF is the real RAM.
		// $0800-$1FFF mirrors it.
		return b.internal[addr%0x0800]
	} else if addr >= 0x8000 {

		romAddr := addr - 0x8000
		if len(b.external) == 0x4000 {
			return b.external[romAddr%0x4000]
		}
		return b.external[romAddr]
	}

	return 0
}

func (b *bus) Write(addr uint16, val uint8) {
	//temp logic add handling later
	if addr <= 0x8000 {
		b.internal[addr] = val

	} else {
		b.external[addr-0x8000] = val
	}
}
