package Core

import (
	"fmt"
	"time"
)

type addressingMode int

const (
	Implied addressingMode = iota
	Accumulator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
	Relative
)

type Debugger struct {
	Console     *console
	Disassembly map[uint16]AssemblyLine

	RecentHistory SnapshotBuffer
}

type SnapshotBuffer struct {
	Frame int

	Data  [200]snapshot
	Index int
}

type ScreenInfo struct {
	Buffer []uint8
}

type AssemblyLine struct {
	Opcode      opCode `json:"Opcode"`
	Disassembly string `json:"disassembly"`
	Val         uint8  `json:"val,omitempty"`
}
type cpustate struct {
	Pc     uint16    `json:"pc"`
	S      uint8     `json:"s"`
	A      uint8     `json:"a"`
	X      uint8     `json:"x"`
	Y      uint8     `json:"y"`
	P      uint8     `json:"p"`
	Flags  FlagState `json:"flags"`
	Cycles int       `json:"cycles"`
	Ram    [][]int   `json:"-"`
}

type FlagState struct {
	Carry     bool `json:"carry"`
	Overflow  bool `json:"overflow"`
	Interrupt bool `json:"interrupt"`
	Zero      bool `json:"zero"`
	Decimal   bool `json:"decimal"`
	Negative  bool `json:"negative"`
}

func (d *Debugger) StartDebugConsole() {

	targetTime := time.Now()

	var framecount = 0

	for {

		if d.Console.Paused {
			time.Sleep(100 * time.Millisecond)
			targetTime = time.Now()

			continue
		}

		now := time.Now()
		for now.After(targetTime) {
			framecount++

			for range 29781 {
				if d.Console.Paused {

					continue
				}
				d.DebugTick()
			}

			targetTime = targetTime.Add(time.Duration(nsPerFrame))

			d.Console.RunDisplayUpdates()
			d.AddSnapshot()
		}

		timeLeft := time.Until(targetTime)
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}
	}
}

func (d *Debugger) DebugTick() {
	d.Console.tick()
	d.DisAssemble(d.Console.Cpu.PC)

}

func (d *Debugger) DisAssemble(addr uint16) AssemblyLine {
	mem := d.Console.Cpu.mem
	if d.Console != nil {
		if _, ok := d.Disassembly[addr]; ok {
			return d.Disassembly[addr]
		}
	}

	line := AssemblyLine{}
	opcode := mem.Read(addr)
	if opcode == 255 {
		d.Console.Pause()
		fmt.Println("Invalid opcode recieved pausing state")
		return line
	}
	info := FetchTable[opcode]
	line.Opcode = info

	switch info.Size {
	case 1:
		line.Disassembly = info.Name
	case 2:
		operand := mem.Read(addr + 1)

		line.Disassembly, line.Val = formatOpcode(info.Name, info.AddressingMode, uint16(operand), d.Console.Cpu, mem)
		line.Val = mem.Read(uint16(operand))
	case 3:
		operand := builduint16(mem.Read(addr+1), mem.Read(addr+2))
		line.Disassembly, line.Val = formatOpcode(info.Name, info.AddressingMode, uint16(operand), d.Console.Cpu, mem)
		line.Val = mem.Read(operand)
	}

	if d.Console != nil {
		d.Disassembly[addr] = line
	}

	return line
}

func formatOpcode(mnemonic string, mode addressingMode, operand uint16, Cpu *cpu, mem typebus) (string, uint8) {
	switch mode {
	case Immediate:
		return fmt.Sprintf("%s #$%02X", mnemonic, operand), uint8(operand)

	case ZeroPage:
		return fmt.Sprintf("%s $%02x", mnemonic, operand), mem.Read(uint16(operand))
	case ZeroPageX:
		return fmt.Sprintf("%s $%02x,X", mnemonic, operand), mem.Read(uint16(operand + uint16(Cpu.X)))
	case ZeroPageY:
		return fmt.Sprintf("%s $%02x,Y", mnemonic, operand), mem.Read(uint16(operand) + uint16(Cpu.Y))
	case Absolute:
		return fmt.Sprintf("%s $%04X", mnemonic, operand), mem.Read(operand)
	case AbsoluteX:
		return fmt.Sprintf("%s $%04X,X", mnemonic, operand), mem.Read(operand + uint16(Cpu.X))
	case AbsoluteY:
		return fmt.Sprintf("%s $%04X,Y", mnemonic, operand), mem.Read(operand + uint16(Cpu.Y))
	case Indirect:

		low := mem.Read(operand)
		var high uint8
		if (operand & 0x00FF) == 0x00FF {
			high = mem.Read(operand & 0xFF00)
		} else {
			high = mem.Read(operand + 1)
		}

		target := builduint16(low, high)

		return fmt.Sprintf("%s ($%04X)", mnemonic, operand), mem.Read(target)
	case IndirectX:
		base := uint16(operand + uint16(Cpu.X))
		low := mem.Read(base)
		high := mem.Read(base + 1)

		return fmt.Sprintf("%s ($%04X),X", mnemonic, operand), mem.Read(builduint16(low, high))
	case IndirectY:
		base := uint16(operand + uint16(Cpu.Y))
		low := mem.Read(base)
		high := mem.Read(base + 1)

		return fmt.Sprintf("%s ($%04X),Y", mnemonic, operand), mem.Read(builduint16(low, high))
	case Relative:
		offset := int8(uint8(operand))
		targetAddr := uint16(int32(Cpu.PC) + 2 + int32(offset))
		return fmt.Sprintf("%s $%04X", mnemonic, targetAddr), 0
	default:
		return "something wrong happened", 0
	}
}

func (d *Debugger) LookAhead(pc uint16, size int) []AssemblyLine {
	lines := make([]AssemblyLine, 0)
	inspectPC := pc

	for range size {
		line := d.DisAssemble(inspectPC)

		inspectPC += uint16(line.Opcode.Size)
		lines = append(lines, line)
	}

	return lines
}

func (c *cpu) GetSate() cpustate {
	state := cpustate{}

	state.A = c.A
	state.X = c.X
	state.Y = c.Y
	state.Pc = c.PC
	state.P = c.P
	state.Cycles = c.totalCycles

	f := FlagState{
		Zero:      getbitBool(state.P, 1),
		Carry:     getbitBool(state.P, 0),
		Negative:  getbitBool(state.P, 7),
		Decimal:   getbitBool(state.P, 3),
		Interrupt: getbitBool(state.P, 2),
		Overflow:  getbitBool(state.P, 6),
	}
	state.Flags = f

	state.S = c.S
	return state
}
