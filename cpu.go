package main

import (
	"log"
)

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
	console *console

	PC uint16 //program counter
	S  uint8  //stack pointer
	P  uint8  //processor status
	A  uint8  //Accumulator
	X  uint8
	Y  uint8

	mem typebus
	temp

	nmiPending   bool
	executingNmi bool
	nmiStep      int

	currentOp   uint8
	currentstep int
	totalCycles int
	Stall       int
}

type cycleStep struct {
	Addr uint16
	Val  uint8
	Mode string
}

type typebus interface {
	Read(uint16) uint8
	Write(uint16, uint8)
	Set(uint16, uint8)
	Get(uint16) uint8
	GetHistory() []cycleStep
	ClearHistory()
}

func (c *cpu) executeNmiCycle() {
	switch c.nmiStep {
	case 1:
		return
	case 2:
		return
	case 3:
		c.mem.Write(0x0100+uint16(c.S), uint8(c.PC>>8))
		c.S--
	case 4:
		c.mem.Write(0x0100+uint16(c.S), uint8(c.PC&0xFF))
		c.S--
	case 5:
		Status := c.P
		Status = AssignBit(Status, 4, false)
		Status = AssignBit(Status, 5, true)

		c.mem.Write(0x0100+uint16(c.S), Status)
		c.S--

		c.setFlag(Interrupt)
	case 6:
		c.low = c.mem.Read(0xFFFA)
	case 7:
		c.high = c.mem.Read(0xFFFB)

		c.PC = builduint16(c.low, c.high)

		c.executingNmi = false
		c.nmiStep = 0
		return

	}

	c.nmiStep++
}

func (c *cpu) tick() {
	c.totalCycles++

	if c.currentstep == 0 && c.nmiPending {
		c.executingNmi = true
		c.nmiPending = false
		c.nmiStep = 1
	}

	if c.executingNmi {
		c.executeNmiCycle()
		return
	}

	if c.currentstep == 0 {

		c.currentOp = c.fetchone()
		c.totalCycles++
	}

	opcode := FetchTable[c.currentOp]
	if opcode.Execute == nil {
		log.Fatalf("CRASH: Attempted to execute unmapped/nil opcode 0x%02X at PC: 0x%04X (Total Cycles: %d)",
			c.currentOp, c.PC-1, c.totalCycles)
	}
	finished := opcode.Execute(c, c.currentstep)

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
	cpu      *cpu
	internal [2048]byte
	external []byte
}

func (b *bus) Read(addr uint16) uint8 {
	var val uint8 = 0
	switch {
	case addr <= 0x1FFF:
		val = b.internal[addr&0x07FF]
	case 0x2000 <= addr && addr <= 0x3FFF:
		RegIndex := (addr - 0x2000) % 8
		val = b.cpu.console.Ppu.ReadReg(RegIndex, b.cpu.console.OpenBusVal)
	case addr >= 0x8000:
		romaddr := addr - 0x8000
		if len(b.external) == 0x4000 {
			val = b.external[romaddr%0x4000]
		} else {
			val = b.external[romaddr]
		}

	}

	b.cpu.console.OpenBusVal = val

	return val
}

func (b *bus) Write(addr uint16, val uint8) {

	switch {
	case addr <= 0x1FFF:
		b.internal[addr&0x07FF] = val
	case 0x2000 <= addr && addr <= 0x3FFF:
		RegIndex := (addr - 0x2000) % 8
		b.cpu.console.Ppu.WriteReg(RegIndex, val)
	case addr == 0x4014:
		b.cpu.console.ExecuteOAMDMA(val)
	case addr <= 0x7FFF:

	default:
		b.external[addr-0x8000] = val
	}
}

func (c *cpu) reset() {
	low := c.mem.Read(0xFFFC)
	high := c.mem.Read(0xFFFD)

	c.PC = builduint16(low, high)

	c.currentstep = 0
	c.currentOp = 0
	c.totalCycles = 0
	c.S = 0xFD
	c.P = 0x24

	c.console.ready = true
}
