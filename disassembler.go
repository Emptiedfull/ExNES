package main

import "fmt"

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

type AssemblyLine struct {
	Opcode      opCode `json:"Opcode"`
	Disassembly string `json:"disassembly"`
	Val         uint8  `json:"val,omitempty"`
}

type cpustate struct {
	Pc  uint16 `json:"pc"`
	S   uint8  `json:"s"`
	A   uint8  `json:"a"`
	X   uint8  `json:"x"`
	Y   uint8  `json:"y"`
	P   uint8  `json:"p"`
	Ram [][]int
}

func DisAssemble(mem typebus, addr uint16, pc uint16) AssemblyLine {

	if debugConsole.console != nil {
		if _, ok := debugConsole.Disassembly[addr]; ok {
			return debugConsole.Disassembly[addr]
		}
	}

	line := AssemblyLine{}
	opcode := mem.Read(addr)
	info := FetchTable[opcode]
	line.Opcode = info

	switch info.Size {
	case 1:
		line.Disassembly = info.Name
	case 2:
		operand := mem.Read(addr + 1)

		line.Disassembly, line.Val = formatOpcode(info.Name, info.AddressingMode, uint16(operand), pc, mem)
		line.Val = mem.Read(uint16(operand))
	case 3:
		operand := builduint16(mem.Read(addr+1), mem.Read(addr+2))
		line.Disassembly, line.Val = formatOpcode(info.Name, info.AddressingMode, uint16(operand), pc, mem)
		line.Val = mem.Read(operand)
	}

	if debugConsole.console != nil {
		debugConsole.Disassembly[addr] = line
	}

	return line
}

func formatOpcode(mnemonic string, mode addressingMode, operand uint16, PC uint16, mem typebus) (string, uint8) {
	switch mode {
	case Immediate:
		return fmt.Sprintf("%s #$%02X", mnemonic, operand), uint8(operand)

	case ZeroPage:
		return fmt.Sprintf("%s $%02x", mnemonic, operand), mem.Read(uint16(operand))
	case ZeroPageX:
		return fmt.Sprintf("%s $%02x,X", mnemonic, operand), mem.Read(uint16(operand + uint16(debugConsole.console.Cpu.X)))
	case ZeroPageY:
		return fmt.Sprintf("%s $%02x,Y", mnemonic, operand), mem.Read(uint16(operand) + uint16(debugConsole.console.Cpu.Y))
	case Absolute:
		return fmt.Sprintf("%s $%04X", mnemonic, operand), mem.Read(operand)
	case AbsoluteX:
		return fmt.Sprintf("%s $%04X,X", mnemonic, operand), mem.Read(operand + uint16(debugConsole.console.Cpu.X))
	case AbsoluteY:
		return fmt.Sprintf("%s $%04X,Y", mnemonic, operand), mem.Read(operand + uint16(debugConsole.console.Cpu.Y))
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
		base := uint16(operand + uint16(debugConsole.console.Cpu.X))
		low := mem.Read(base)
		high := mem.Read(base + 1)

		return fmt.Sprintf("%s ($%04X),X", mnemonic, operand), mem.Read(builduint16(low, high))
	case IndirectY:
		base := uint16(operand + uint16(debugConsole.console.Cpu.Y))
		low := mem.Read(base)
		high := mem.Read(base + 1)

		return fmt.Sprintf("%s ($%04X),Y", mnemonic, operand), mem.Read(builduint16(low, high))
	case Relative:
		offset := int8(uint8(operand))
		targetAddr := uint16(int32(PC) + 2 + int32(offset))
		return fmt.Sprintf("%s $%04X", mnemonic, targetAddr), 0
	default:
		return "something wrong happened", 0
	}
}

func (c *cpu) lookAhead(pc uint16, size int) []AssemblyLine {
	lines := make([]AssemblyLine, 0)
	inspectPC := pc

	for range size {
		line := DisAssemble(c.mem, inspectPC, inspectPC)

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
	state.S = c.S
	return state
}
