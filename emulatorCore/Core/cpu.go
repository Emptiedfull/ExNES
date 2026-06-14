package Core

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

	mem *bus
	temp

	nmiPending   bool
	executingNmi bool
	nmiStep      int
	nmiS         int

	currentOp   uint8
	currentstep int
	fetchNew    bool
	TotalCycles int
	Stall       int

	isJamming bool
}

type cycleStep struct {
	Addr uint16
	Val  uint8
	Mode string
}

func (c *cpu) executeNmiCycle() int {

	switch c.nmiStep {
	case 1:

		return 2
	case 2:

		return 3
	case 3:
		c.mem.Write(0x0100+uint16(c.S), uint8(c.PC>>8))
		c.S--
		return 4
	case 4:
		c.mem.Write(0x0100+uint16(c.S), uint8(c.PC&0xFF))
		c.S--
		return 5
	case 5:
		Status := c.P
		Status = AssignBit(Status, 4, false)
		Status = AssignBit(Status, 5, true)

		c.mem.Write(0x0100+uint16(c.S), Status)
		c.S--

		c.setFlag(Interrupt)
		return 6
	case 6:
		c.low = c.mem.Read(0xFFFA)
		return 7
	case 7:
		c.high = c.mem.Read(0xFFFB)

		c.PC = builduint16(c.low, c.high)

		c.executingNmi = false

		return 1
	default:
		return 1

	}

}

func (c *cpu) tick() {

	if c.isJamming {
		c.mem.Read(c.PC)
		return
	}

	c.TotalCycles++

	if c.currentstep == 0 && c.nmiPending {
		c.executingNmi = true
		c.nmiPending = false
		c.nmiS++
	}

	if c.executingNmi {

		c.nmiStep = c.executeNmiCycle()
		return
	}

	// if c.fetchNew {
	// 	c.currentOp = c.fetchone()

	// 	c.fetchNew = false
	// 	// c.TotalCycles++
	// 	c.temp = temp{}
	// 	return

	// }

	if c.currentstep == 0 {
		c.temp = temp{}
		c.currentOp = c.fetchone()
		c.currentstep = 1
		return
	}

	opcode := FetchTable[c.currentOp]
	finished := opcode.Execute(c, c.currentstep-1)

	if finished {
		c.currentstep = 0
		c.fetchNew = true

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

	mapper Mapper
}

func (b *bus) returnInternal() [2048]byte {
	return b.internal
}

func (b *bus) returnExternal() []byte {
	dst := make([]byte, len(b.external))
	copy(dst, b.external)
	return dst
}

func (b *bus) loadInternal(x [2048]byte) {
	b.internal = x
}

func (b *bus) loadExternal(x []byte) {
	copy(b.external, x)
}

func (b *bus) Read(addr uint16) uint8 {
	var val uint8 = 0
	switch {
	case addr <= 0x1FFF:
		val = b.internal[addr&0x07FF]
	case addr == 0x4016:

		return b.cpu.console.JoyPad.readState()
	case addr == 0x4017:
		// return b.cpu.console.JoyPad.readState()
	case 0x2000 <= addr && addr <= 0x3FFF:
		RegIndex := (addr - 0x2000) % 8
		val = b.cpu.console.Ppu.ReadReg(RegIndex, b.cpu.console.OpenBusVal)
	case addr >= 0x8000:
		val = b.mapper.ReadPRG(addr)
	}

	b.cpu.console.OpenBusVal = val

	return val
}

func (b *bus) Write(addr uint16, val uint8) {

	switch {
	case addr <= 0x1FFF:
		b.internal[addr&0x07FF] = val
	case addr == 0x4016:
		b.cpu.console.JoyPad.writeStrobe(val & 1)
	case 0x2000 <= addr && addr <= 0x3FFF:
		RegIndex := (addr - 0x2000) % 8
		b.cpu.console.Ppu.WriteReg(RegIndex, val)
	case addr == 0x4014:
		b.cpu.console.ExecuteOAMDMA(val)
	case addr >= 0x8000:
		b.mapper.WritePRG(addr, val)
	case addr <= 0x7FFF:
	default:

	}
}

func (c *cpu) Reset() {
	low := c.mem.Read(0xFFFC)
	high := c.mem.Read(0xFFFD)

	c.PC = builduint16(low, high)

	c.currentstep = 0
	c.currentOp = 0
	c.TotalCycles = 0
	c.S = 0xFD
	c.P = 0x24

	c.console.ready = true
	c.isJamming = false
	c.temp = temp{}
}
