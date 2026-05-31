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
	addr        uint16
	opcode      opCode
	disassembly string
}

func DisAssemble(mem typebus, addr uint16, pc uint16) AssemblyLine {

	line := AssemblyLine{
		addr: addr,
	}

	opcode := mem.Read(addr)
	info := FetchTable[opcode]
	line.opcode = info

	switch info.Size {
	case 1:
		line.disassembly = info.Name
	case 2:
		val := mem.Read(addr + 1)
		line.disassembly = formatOpcode(info.Name, info.AddressingMode, uint16(val), pc)
	case 3:
		val := builduint16(mem.Read(addr+1), mem.Read(addr+2))
		line.disassembly = formatOpcode(info.Name, info.AddressingMode, uint16(val), pc)
	}

	return line
}

func formatOpcode(mnemonic string, mode addressingMode, val uint16, PC uint16) string {
	switch mode {
	case Immediate:
		return fmt.Sprintf("%s #$%02X", mnemonic, val)
	case ZeroPage:
		return fmt.Sprintf("%s $%02x", mnemonic, val)
	case ZeroPageX:
		return fmt.Sprintf("%s $%02x,X", mnemonic, val)
	case ZeroPageY:
		return fmt.Sprintf("%s $%02x,Y", mnemonic, val)
	case Absolute:
		return fmt.Sprintf("%s $%04X", mnemonic, val)
	case AbsoluteX:
		return fmt.Sprintf("%s $%04X,X", mnemonic, val)
	case AbsoluteY:
		return fmt.Sprintf("%s $%04X,Y", mnemonic, val)
	case Indirect:
		return fmt.Sprintf("%s ($%04X)", mnemonic, val)
	case IndirectX:
		return fmt.Sprintf("%s ($%04X),X", mnemonic, val)
	case IndirectY:
		return fmt.Sprintf("%s ($%04X),Y", mnemonic, val)
	case Relative:
		offset := int8(uint8(val))
		targetAddr := uint16(int32(PC) + 2 + int32(offset))
		return fmt.Sprintf("%s $%04X", mnemonic, targetAddr)
	default:
		return fmt.Sprint("something wrong happened")
	}
}

func (c *cpu) lookAhead(pc uint16) []AssemblyLine {
	lines := make([]AssemblyLine, 0)
	inspectPC := pc

	for range 5 {
		line := DisAssemble(c.mem, inspectPC, inspectPC)

		inspectPC += uint16(line.opcode.Size)
		lines = append(lines, line)
	}

	return lines
}
