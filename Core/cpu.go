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

type interuptKind int

const (
	intNone interuptKind = iota
	intNMI
	intIRQ
	intBRK
)

type Cpu struct {
	console *Console

	PC uint16 //program counter
	S  uint8  //stack pointer
	P  uint8  //processor status
	A  uint8  //Accumulator
	X  uint8
	Y  uint8

	Mem *bus
	temp

	intPresent interuptKind
	intStep    int
	intVector  uint16

	NmiLine bool
	irqLine bool

	currentOp   uint8
	currentstep int
	fetchNew    bool
	TotalCycles int
	Stall       int

	isJamming bool
}

type CycleStep struct {
	Addr int
	Val  int
	Mode string
}

// func (c *Cpu) executeNmiCycle() int {

// 	switch c.nmiStep {
// 	case 1:

// 		return 2
// 	case 2:

// 		return 3
// 	case 3:
// 		c.Mem.Write(0x0100+uint16(c.S), uint8(c.PC>>8))
// 		c.S--
// 		return 4
// 	case 4:
// 		c.Mem.Write(0x0100+uint16(c.S), uint8(c.PC&0xFF))
// 		c.S--
// 		return 5
// 	case 5:
// 		Status := c.P
// 		Status = AssignBit(Status, 4, false)
// 		Status = AssignBit(Status, 5, true)

// 		c.Mem.Write(0x0100+uint16(c.S), Status)
// 		c.S--

// 		c.setFlag(Interrupt)
// 		return 6
// 	case 6:
// 		c.low = c.Mem.Read(0xFFFA)
// 		return 7
// 	case 7:
// 		c.high = c.Mem.Read(0xFFFB)

// 		c.PC = builduint16(c.low, c.high)

// 		c.executingNmi = false

// 		return 1
// 	default:
// 		return 1

// 	}

// }

func (c *Cpu) ResetTemp() {
	c.temp = temp{}

}

func (c *Cpu) Tick() {

	if c.isJamming {
		c.Mem.Read(c.PC)
		return
	}

	c.TotalCycles++

	if c.intPresent != intNone {
		c.stepInt()
		return
	}

	if c.currentstep == 0 {
		if c.NmiLine {
			c.NmiLine = false
			c.intStep = 1
			c.intPresent = intNMI
			c.stepInt()
			return
		}

		if c.irqLine && !c.getFlag(Interrupt) {
			c.intStep = 1
			c.intPresent = intIRQ
			c.stepInt()
			return
		}
	}

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

func (c *Cpu) stepInt() {
	switch c.intStep {
	case 1:
		c.Mem.Read(c.PC)
	case 2:
		c.Mem.Read(c.PC)
		if c.intPresent == intBRK {
			c.PC++
		}
	case 3:
		c.pushStack((uint8(c.PC >> 8)))
	case 4:
		c.pushStack(uint8(c.PC & 0xFF))
	case 5:
		if c.intPresent != intNMI && c.NmiLine {
			c.intPresent = intNMI
			c.NmiLine = false
		}

		c.intVector = 0xFFFE
		if c.intPresent == intNMI {
			c.intVector = 0xFFFA
		}

		c.low = c.Mem.Read(c.intVector)
	case 6:
		c.high = c.Mem.Read(c.intVector + 1)
		c.PC = builduint16(c.low, c.high)

		c.intPresent = intNone
		c.intStep = 0

		return
	}

	c.intStep++
}

type temp struct {
	high uint8
	low  uint8

	pointer uint8
	addr    uint16
	val     uint8
}

type bus struct {
	Cpu      *Cpu
	internal [2048]byte

	FlatMode bool
	FlatMem  []uint8
	Log      []CycleStep
}

func GetBus() *bus {
	return &bus{}
}

func (b *bus) returnInternal() [2048]byte {
	return b.internal
}

func (b *bus) loadInternal(x [2048]byte) {
	b.internal = x
}

func (b *bus) Read(addr uint16) uint8 {

	if b.FlatMode {
		val := b.FlatMem[addr]
		b.Log = append(b.Log, CycleStep{Mode: "read", Addr: int(addr), Val: int(val)})
		return val
	}

	var val uint8 = 0
	if b.Cpu.console.CheatEngine.Enabled {
		b.Cpu.console.CheatEngine.TableMutex.Lock()
		defer b.Cpu.console.CheatEngine.TableMutex.Unlock()
		if cheat, ok := b.Cpu.console.CheatEngine.cheatTable[addr]; ok {

			// fmt.Println("cheat found")
			if cheat.compare {
				if cheat.compareVal == val {
					return cheat.val
				}
			} else {
				// fmt.Println("returning cheat")
				return cheat.val
			}

		}
	}

	switch {
	case addr <= 0x1FFF:
		val = b.internal[addr&0x07FF]
	case addr == 0x4015:
		return (b.Cpu.console.Apu.readStatus() & 0xDF) | (b.Cpu.console.OpenBusVal & 0x20)
	case addr == 0x4016:
		val = (b.Cpu.console.Player1.readState() & 0x01) | (b.Cpu.console.OpenBusVal & 0xE0)
		b.Cpu.console.OpenBusVal = val
		return val
	case addr == 0x4017:
		val = (b.Cpu.console.Player2.readState() & 0x01) | (b.Cpu.console.OpenBusVal & 0xE0)
		b.Cpu.console.OpenBusVal = val
		return val

	case 0x2000 <= addr && addr <= 0x3FFF:
		RegIndex := (addr - 0x2000) % 8
		val = b.Cpu.console.Ppu.ReadReg(RegIndex, b.Cpu.console.OpenBusVal)
	case addr >= 0x5000:
		val = b.Cpu.console.mapper.ReadPRG(addr)
	default:

		return b.Cpu.console.OpenBusVal
	}

	b.Cpu.console.OpenBusVal = val

	return val
}

func (b *bus) Write(addr uint16, val uint8) {

	if b.FlatMode {
		b.FlatMem[addr] = val
		b.Log = append(b.Log, CycleStep{Mode: "write", Addr: int(addr), Val: int(val)})
		return
	}
	b.Cpu.console.OpenBusVal = val

	switch {
	case addr <= 0x1FFF:
		b.internal[addr&0x07FF] = val
	case addr == 0x4014:
		b.Cpu.console.ExecuteOAMDMA(val)

	case addr == 0x4016:
		b.Cpu.console.Player1.writeStrobe(val & 1)
		b.Cpu.console.Player2.writeStrobe(val & 1)
	case addr >= 0x4000 && addr <= 0x4017:
		b.Cpu.console.Apu.writeReg(addr, val)
	case 0x2000 <= addr && addr <= 0x3FFF:
		RegIndex := (addr - 0x2000) % 8
		b.Cpu.console.Ppu.WriteReg(RegIndex, val)

	case addr >= 0x5000:
		b.Cpu.console.mapper.WritePRG(addr, val)
	case addr <= 0x7FFF:
	default:

	}
}

func (c *Cpu) Reset() {
	low := c.Mem.Read(0xFFFC)
	high := c.Mem.Read(0xFFFD)

	c.PC = builduint16(low, high)

	c.currentstep = 0
	c.currentOp = 0
	c.TotalCycles = 0
	c.S = 0xFD
	c.P = 0x24

	c.isJamming = false
	c.temp = temp{}

}
