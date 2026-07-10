package Core

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

type Cpu struct {
	console *Console
	dma     DMAUnit

	PC uint16 //program counter
	S  uint8  //stack pointer
	P  uint8  //processor status
	A  uint8  //Accumulator
	X  uint8
	Y  uint8

	Mem *bus
	temp

	nmiPending   bool
	executingNmi bool
	nmiStep      int
	nmiS         int

	currentOp   uint8
	currentstep int

	TotalCycles int
	Stall       int

	isJamming bool

	dmaHijackRequested bool
	hijackedThisTick   bool
	hijacked           bool
}

type CycleStep struct {
	Addr int
	Val  int
	Mode string
}

func (c *Cpu) triggerIRQ() {
	if c.getFlag(Interrupt) {
		return
	}

	c.pushStack(uint8(c.PC >> 8))
	c.pushStack(uint8(c.PC & 0xFF))
	c.pushStack(c.P &^ Break)

	c.setFlag(Interrupt)

	l := c.Mem.Read(0xFFFE)
	h := c.Mem.Read(0xFFFF)

	c.PC = builduint16(l, h)

	c.Stall += 7

}

var opcodesTouched = map[uint8]struct{}{}

func (c *Cpu) executeNmiCycle() int {

	switch c.nmiStep {
	case 1:

		return 2
	case 2:

		return 3
	case 3:
		c.Mem.Write(0x0100+uint16(c.S), uint8(c.PC>>8))
		c.S--
		return 4
	case 4:
		c.Mem.Write(0x0100+uint16(c.S), uint8(c.PC&0xFF))
		c.S--
		return 5
	case 5:
		Status := c.P
		Status = AssignBit(Status, 4, false)
		Status = AssignBit(Status, 5, true)

		c.Mem.Write(0x0100+uint16(c.S), Status)
		c.S--

		c.setFlag(Interrupt)
		return 6
	case 6:
		c.low = c.Mem.Read(0xFFFA)
		return 7
	case 7:
		c.high = c.Mem.Read(0xFFFB)

		c.PC = builduint16(c.low, c.high)

		c.executingNmi = false

		return 1
	default:
		return 1

	}

}

func (c *Cpu) ResetTemp() {
	c.temp = temp{}

}

func (c *Cpu) Tick() {

	if c.isJamming {
		c.Mem.Read(c.PC)
		return
	}

	if c.currentstep == 0 && c.nmiPending {
		c.executingNmi = true
		c.nmiPending = false
		c.nmiS++
	}

	hijacked := c.runProtected(func() {
		if c.executingNmi {
			c.nmiStep = c.executeNmiCycle()
			return
		}

		if c.currentstep == 0 {
			c.temp = temp{}
			c.currentOp = c.fetchone()
			c.currentstep = 1
			return
		}

		opCode := FetchTable[c.currentOp]
		finished := opCode.Execute(c, c.currentstep-1)

		if finished {
			c.currentstep = 0

		} else {
			c.currentstep++
		}
	})

	if hijacked {

		if _, ok := opcodesTouched[c.currentOp]; !ok {
			opcodesTouched[c.currentOp] = struct{}{}
			fmt.Printf("code: %x, cycles: %v \n", c.currentOp, c.TotalCycles)

		}

		c.console.busHijacked = true
		return
	}

}

type hijackSignal struct{}

func (c *Cpu) runProtected(f func()) (hijacked bool) {
	defer func() {

		if r := recover(); r != nil {
			if _, ok := r.(hijackSignal); ok {
				hijacked = true
				return
			}
			panic(r)
		}

	}()
	f()
	return false
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

	if b.Cpu.dmaHijackRequested && !b.Cpu.console.busHijacked {

		b.Cpu.dmaHijackRequested = false

		b.Cpu.hijackedThisTick = true

		panic(hijackSignal{})

	}

	var val uint8 = 0
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
	case addr >= 0x6000:
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

	case addr >= 0x6000:
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

	c.console.ready = true
	c.isJamming = false
	c.temp = temp{}

}

type DMA_TYPE int

const (
	dmc DMA_TYPE = iota
)

type DMAUnit struct {
	cpu  *Cpu
	step int
	kind DMA_TYPE
}

func (d *DMAUnit) tick() bool {

	if d.step == 0 {
		d.step++
		return false
	}

	if d.step == 1 {
		d.step++
		if d.kind == dmc {

			return false
		}
	}

	if d.step == 2 {
		d.step++
		if d.cpu.TotalCycles%2 == 0 {
			// fmt.Println("passing", d.cpu.TotalCycles)
			return false
		}
	}

	if d.step == 3 {
		d.step = 0
		if d.kind == dmc {
			val := d.cpu.Mem.Read(d.cpu.console.Apu.Dmc.currentAddr)
			d.cpu.console.Apu.Dmc.LoadSample(val)

			return true
		}

	}

	return true
}

func (c *Cpu) runStep(op opCode) (finished bool) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(hijackSignal); ok {
				finished = false
				return
			}
			panic(r)
		}

	}()
	return op.Execute(c, c.currentstep-1)
}
