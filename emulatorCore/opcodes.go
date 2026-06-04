package main

import (
	"fmt"
)

type opCode struct {
	Name           string                      `json:"name"`
	AddressingMode addressingMode              `json:"mode"`
	Size           uint8                       `json:"size"`
	Execute        func(c *cpu, step int) bool `json:"-"`
}

var FetchTable = []opCode{
	0x00: {
		Name:           "BRK",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.fetchone()
				return false
			case 1:
				c.pushStack(uint8(c.PC >> 8))
				return false
			case 2:
				c.pushStack(uint8(c.PC & 0x0FF))
				return false
			case 3:
				c.setFlag(Break)

				c.pushStack(c.P)
				c.updateFlag(Break, false)
				c.setFlag(Interrupt)
				return false
			case 4:
				c.temp.low = c.mem.Read(0xFFFE)

				return false
			case 5:
				c.temp.high = c.mem.Read(0xFFFF)
				c.PC = builduint16(c.low, c.high)
				return true
			default:
				fmt.Println("somethign wrong at 0x00")
				return true
			}
		},
	},

	0x01: {
		Name:           "ORA",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone() //zero-page base addr
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				val := c.mem.Read(baseAddr)
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				return true
			}

		},
	},
	0x02: {
		Name:           "JAM",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {

			c.isJamming = true

			return true
		},
	},

	0x03: {
		Name:           "SLO",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.temp.pointer + 1))
				c.temp.addr = builduint16(c.low, c.high)
				return false
			case 4:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 5:
				c.temp.val = c.ASL(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)
				return true
			default:

				fmt.Println("something really wrong at 0x03")
				return true

			}
		},
	},

	0x04: {
		Name:           "NOP",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.fetchone()
				return false
			case 1:
				return true
			default:
				fmt.Println("something wrong at 0x04")
				return false
			}
		},
	},

	0x05: {
		Name:           "ORA",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				return true
			}
		},
	},

	0x06: {
		Name:           "ASL",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				var f bool
				c.temp.val, f = performASL(c.temp.val)
				if f {
					c.setFlag(Carry)
				} else {
					c.clearFlag(Carry)
				}
				c.SetFlagNZ(c.temp.val)
				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				return true
			}
		},
	},

	0x07: {
		Name:           "SLO",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.val = c.ASL(c.val)
				return false
			case 3:
				c.mem.Write(uint16(c.pointer), c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("somethign wrong at 0x07")
				return true
			}
		},
	},

	0x08: {
		Name:           "PHP impl",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				tempStatus := c.P | 0x30
				c.pushStack(tempStatus)
				return true
			default:
				return true
			}
		},
	},
	0x09: {
		Name:           "ORA",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.A |= val
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x0A: {

		Name:           "ASL A",
		AddressingMode: Accumulator,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			var f bool
			c.A, f = performASL(c.A)
			if f {
				c.setFlag(Carry)
			} else {
				c.clearFlag(Carry)
			}
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x0B: {
		Name:           "ANC",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()

			c.A &= val
			c.updateFlag(Carry, getbitBool(c.A, 7))
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x0C: {
		Name:           "NOP",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				return true
			default:
				fmt.Println("something really wrong at 0x0C")
				return true
			}
		},
	},

	0x0D: {
		Name:           "ORA",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.A |= c.mem.Read(addr)
				c.SetFlagNZ(c.A)
				return true
			default:
				return true
			}
		},
	},
	0x0E: {
		Name:           "ASL",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.val = c.mem.Read(builduint16(c.temp.low, c.temp.high))
				return false
			case 3:
				var carry bool
				c.temp.val, carry = performASL(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(builduint16(c.temp.low, c.temp.high), c.temp.val)
				return true
			default:
				return true
			}

		},
	},

	0x0F: {
		Name:           "SLO",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.temp.low, c.temp.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 3:
				c.val = c.ASL(c.val)
				return false
			case 4:
				c.mem.Write(c.addr, c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something really wrong at 0x0F")
				return true
			}
		},
	},

	// 0x1X Seris- This is not ai pls dont kill me :<

	0x10: {
		Name:           "BPL",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()

				if c.getFlag(Negative) {
					return true
				}
				return false
			case 1:
				offset := int8(c.temp.val)
				oldPc := c.PC
				newPC := uint16(int32(oldPc) + int32(offset))

				c.PC = newPC

				if crossedPage(oldPc, newPC) {
					return false
				} else {
					return true
				}
			case 2:
				return true
			default:
				fmt.Println("soemthign wrong happened", "0x10")
				return true
			}
		},
	},
	0x11: {
		Name:           "ORA",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16((c.temp.pointer + 1) & 0xFF))
				return false
			case 3:
				baseAddr := (uint16(c.temp.high)<<8 | uint16(c.temp.low))
				c.temp.addr = baseAddr + uint16(c.Y)

				if crossedPage(baseAddr, c.temp.addr) {
					return false
				}

				fallthrough

			case 4:
				val := c.mem.Read(c.temp.addr)
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				return true
			}

		},
	},

	0x12: {
		Name:           "JAM",
		AddressingMode: Implied,
		Size:           10,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x13: {
		Name:           "SLO IND Y",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				baseAddr := builduint16(c.low, c.high)
				c.addr = baseAddr + uint16(c.Y)
				return false
			case 4:
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.val = c.ASL(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something wrong at 0x13")
				return true
			}
		},
	},

	0x14: {
		Name:           "NOP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				return true
			default:
				fmt.Println("something really wrong at 0x14")
				return true
			}
		},
	},

	0x15: {
		Name:           "ORA",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("somethign wrong happened at 0x15")
				return true
			}
		},
	},
	0x16: {
		Name:           "ASL Oper,X",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				var carry bool
				c.temp.val, carry = performASL(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				return true
			}

		},
	},
	0x17: {
		Name:           "SLO",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.val = c.ASL(c.val)
				return false
			case 4:
				c.mem.Write(uint16(c.pointer), c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("soemthign wrong at 0x17")
				return true
			}
		},
	},

	0x18: {
		Name:           "CLC",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(Carry)
			return true
		},
	},

	0x19: {
		Name:           "ORA ABS,Y",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				newAddr := baseAddr + uint16(c.Y)
				c.temp.addr = newAddr

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad happened at 0x19")
				return true

			}
		},
	},

	0x1A: {
		Name:           "NOP",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x1B: {
		Name:           "SLO",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {

			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.ASL(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("wrong at 0x1B")
				return true
			}
		},
	},

	0x1C: {
		Name:           "NOP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				newAddr := baseAddr + uint16(c.X)

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				return true
			default:
				fmt.Println("wrong at 0x1C")
				return true
			}
		},
	},

	0x1D: {
		Name:           "ORA",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				newAddr := baseAddr + uint16(c.X)
				c.temp.addr = newAddr

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad happened at 0x1D")
				return true

			}
		},
	},

	0x1E: {
		Name:           "ASL",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = baseAddr + uint16(c.X)
				return false
			case 3:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				var carry bool
				c.temp.val, carry = performASL(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 5:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0x1E")
				return true
			}
		},
	},

	0x1F: {
		Name:           "SLO",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.X)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.ASL(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)

				c.A |= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("something wrong at 0x1F")
				return true
			}
		},
	},

	//0x2X Series

	0x20: {
		Name:           "JSR",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()

				return false
			case 1:
				return false

			case 2:
				c.pushStack(uint8(c.PC >> 8))
				return false
			case 3:
				c.pushStack(uint8(c.PC & 0xFF))
				return false
			case 4:
				c.temp.high = c.fetchone()
				c.PC = builduint16(c.temp.low, c.temp.high)
				return true
			default:
				fmt.Println("somethign very wrong happened")
				return true
			}
		},
	},
	0x21: {
		Name:           "AND",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				val := c.mem.Read(baseAddr)
				c.A &= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad happened at 0x21")
				return true
			}
		},
	},

	0x22: {
		Name:           "JAM",
		AddressingMode: Implied,
		Size:           10,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x23: {
		Name:           "RLA",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.temp.pointer))
				c.temp.pointer += c.X
				return false
			case 2:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 4:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.val = c.ROL(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)
				c.A &= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("somethign wrong at 0x23")
				return true
			}
		},
	},

	0x24: {
		Name:           "BITS",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))

				c.updateFlag(Zero, (c.A&val) == 0)
				c.updateFlag(Negative, (val&0x80) != 0)
				c.updateFlag(oVerflow, (val&0x40) != 0)

				return true
			default:
				fmt.Println("something bad happened at 0x24")
				return true
			}
		},
	},
	0x25: {
		Name:           "AND",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A &= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad happened at 0x25")
				return true
			}
		},
	},
	0x26: {
		Name:           "ROL",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:

				c.temp.val = c.ROL(c.val)

				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("something bad happened at 0x26")
				return true
			}

		},
	},

	0x27: {
		Name:           "RLA",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.val = c.ROL(c.val)
				return false
			case 3:
				c.mem.Write(uint16(c.pointer), c.val)

				c.A &= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("soemthing wrong at 0x27")
				return true
			}
		},
	},

	0x28: {
		Name:           "PLP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				return false
			case 2:
				val := c.popStack()
				c.P = (val & 0xEF) | 0x20
				return true
			default:
				fmt.Println("somethign bad happened at 0x28")
				return true
			}
		},
	}, 0x29: {
		Name:           "AND",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.A &= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x2A: {
		Name:           "ROL",
		AddressingMode: Accumulator,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			var carry bool
			c.A, carry = performROL(c.A, c.getFlag(Carry))
			c.updateFlag(Carry, carry)
			c.SetFlagNZ(c.A)
			return true
		},
	},

	0x2B: {
		Name:           "ANC",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()

			c.A &= val
			c.updateFlag(Carry, getbitBool(c.A, 7))
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x2C: {
		Name:           "BITS",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))

				c.updateFlag(Zero, (c.A&val) == 0)
				c.updateFlag(Negative, (val&0x80) != 0)
				c.updateFlag(oVerflow, (val&0x40) != 0)

				return true
			default:
				fmt.Println("something bad happened at 0x2C")
				return true
			}
		},
	},

	0x2D: {
		Name:           "AND",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				val := c.mem.Read(addr)

				c.A &= val
				c.SetFlagNZ(c.A)
				return true

			default:
				fmt.Print("soemthing bad at 2D")
				return true
			}
		},
	},

	0x2E: {
		Name:           "ROL",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high)
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 3:
				var carry bool
				c.temp.val, carry = performROL(c.temp.val, c.getFlag(Carry))
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("soemthing gap at 2E")
				return true
			}
		},
	},

	0x2F: {
		Name:           "RLA",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 3:
				c.val = c.ROL(c.val)
				return false
			case 4:
				c.mem.Write(c.addr, c.val)

				c.A &= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("somethign wrong at 2F")
				return true
			}
		},
	},

	//0x3X Series

	0x30: {
		Name:           "BMI",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()

				if !c.getFlag(Negative) {
					return true
				}

				return false
			case 1:
				offset := int8(c.temp.val)
				oldPC := c.PC

				c.PC = uint16(int32(oldPC) + int32(offset))

				if crossedPage(oldPC, c.PC) {
					return false
				} else {
					return true
				}
			case 2:
				return true
			default:
				fmt.Println("bad somethign at 0x30")
				return true
			}
		},
	},
	0x31: {
		Name:           "AND",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				base := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = base + uint16(c.Y)

				if crossedPage(base, c.temp.addr) {
					return false
				}

				fallthrough
			case 4:
				val := c.mem.Read(c.temp.addr)
				c.A &= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad at 0x31")
				return true
			}
		},
	},

	0x32: {
		Name:           "JAM",
		AddressingMode: Implied,
		Size:           10,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x33: {
		Name:           "RLA IND Y",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 4:
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.val = c.ROL(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)

				c.A &= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something wrong at 0x33")
				return true
			}
		},
	},

	0x34: {
		Name:           "NOP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				return true
			default:
				fmt.Println("something wrong at 0x34")
				return true
			}
		},
	},

	0x35: {
		Name:           "AND",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A &= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("someting bad at 35")
				return true
			}
		},
	},

	0x36: {
		Name:           "ROL",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.val = c.ROL(c.val)
				return false
			case 4:
				c.mem.Write(uint16(c.pointer), c.val)
				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("something wrong at 0x36")
				return true

			}
		},
	},

	0x37: {
		Name:           "RLA",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.val = c.ROL(c.val)
				return false
			case 4:
				c.mem.Write(uint16(c.pointer), c.val)

				c.A &= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("something wrong at 0x37")
				return true
			}
		},
	},

	0x38: {
		Name:           "SEC",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.setFlag(Carry)
			return true
		},
	},
	0x39: {
		Name:           "AND",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				base := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = base + uint16(c.Y)
				if crossedPage(base, c.temp.addr) {
					return false
				}
				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.A &= val
				c.SetFlagNZ(c.A)
				return true

			default:
				return false
			}
		},
	},

	0x3A: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x3B: {
		Name:           "RLA",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.ROL(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)

				c.A &= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("something wrong at 3B")
				return true
			}
		},
	},

	0x3C: {
		Name:           "NOP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				newAddr := baseAddr + uint16(c.X)

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				return true
			default:
				fmt.Println("wrong at 0x1C")
				return true
			}
		},
	},

	0x3D: {
		Name:           "AND",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				base := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = base + uint16(c.X)
				if crossedPage(base, c.temp.addr) {
					return false
				}
				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.A &= val
				c.SetFlagNZ(c.A)
				return true

			default:
				return false
			}
		},
	},
	0x3E: {
		Name:           "ROL ABS,X",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high)
				c.temp.addr += uint16(c.X)
				return false
			case 3:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				var carry bool
				c.temp.val, carry = performROL(c.temp.val, c.getFlag(Carry))
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 5:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0x3E")
				return true
			}

		},
	},
	0x3F: {
		Name:           "RLA",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.X)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.ROL(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)

				c.A &= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("something wrong at 3B")
				return true
			}
		},
	},

	//0x4X Series

	0x40: {
		Name:           "RTI",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				return false
			case 2:
				val := c.popStack()

				c.P = (val & 0xEF) | 0x20
				return false
			case 3:
				c.temp.low = c.popStack()
				return false
			case 4:
				c.temp.high = c.popStack()
				c.PC = builduint16(c.temp.low, c.temp.high)
				return true
			default:
				fmt.Println("bad at 0x40")
				return true
			}
		},
	},
	0x41: {
		Name:           "EOR",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad at 0x41")
				return true
			}
		},
	},

	0x42: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x43: {
		Name:           "SRE",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 4:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.val = c.LSR(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("somethign wrong at 0x43")
				return true
			}
		},
	},

	0x44: {
		Name:           "NOP",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.fetchone()
				return false
			case 1:
				return true
			default:
				fmt.Println("something wrong at 0x04")
				return false
			}
		},
	},

	0x45: {
		Name:           "EOR",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("smethign bad at 0x45")
				return true
			}
		},
	}, 0x46: {
		Name:           "LSR",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				var carry bool
				c.temp.val, carry = performLSR(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0x46")
				return true
			}
		},
	},

	0x47: {
		Name:           "SRE",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.val = c.LSR(c.val)
				return false
			case 3:
				c.mem.Write(uint16(c.pointer), c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("something wrong at 0x47")
				return true
			}
		},
	},

	0x48: {
		Name:           "PHA",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				c.pushStack(c.A)
				return true
			default:
				fmt.Println("bad at 0x48")
				return true
			}
		},
	}, 0x49: {
		Name:           "EOR",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.A ^= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x4A: {
		Name:           "LSR",
		AddressingMode: Accumulator,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			var carry bool
			c.A, carry = performLSR(c.A)
			c.updateFlag(Carry, carry)
			c.SetFlagNZ(c.A)
			return true
		},
	},

	0x4B: {
		Name:           "ALR",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			mid := c.A & val

			c.A = c.LSR(mid)
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x4C: {
		Name:           "JMP",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()

				targetArr := builduint16(c.temp.low, c.temp.high)

				c.PC = targetArr
				return true
			default:
				fmt.Println("something wrong at 0x4C")
				return true

			}
		},
	},
	0x4D: {
		Name:           "EOR",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad 0x4D")
				return true
			}
		},
	}, 0x4E: {
		Name:           "LSR",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high)
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 3:
				var carry bool
				c.temp.val, carry = performLSR(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Print("bad 0x4E")
				return true
			}
		},
	},

	0x4F: {
		Name:           "SRE",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 3:
				c.val = c.LSR(c.val)
				return false
			case 4:
				c.mem.Write(c.addr, c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("soemthing wrong at 0x4f")
				return true
			}
		},
	},

	//0x5X Series

	0x50: {
		Name:           "BVC",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()
				if c.getFlag(oVerflow) {
					return true
				}
				return false
			case 1:
				offset := int8(c.temp.val)
				oldPC := c.PC

				c.PC = uint16(int32(oldPC) + int32(offset))

				if !crossedPage(oldPC, c.PC) {
					return true
				} else {
					return false
				}

			case 2:
				return true
			default:
				fmt.Println("bad at 0x50")
				return true
			}
		},
	},
	0x51: {
		Name:           "EOR",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}
				fallthrough
			case 4:
				val := c.mem.Read(c.temp.addr)
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0x51")
				return true
			}
		},
	},

	0x52: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0x53: {
		Name:           "SRE",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 4:
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.val = c.LSR(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("bad at 0x54")
				return true
			}
		},
	},

	0x54: {
		Name:           "NOP",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.fetchone()
				return false
			case 1:
				return false
			case 2:
				return true
			default:
				fmt.Println("something wrong at 0x54")
				return true
			}
		},
	},

	0x55: {
		Name:           "EOR",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Print("bad at 0x55")
				return true
			}
		},
	},
	0x56: {
		Name:           "LSR",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				var carry bool
				c.temp.val, carry = performLSR(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("bad at 0x56")
				return true
			}
		},
	},

	0x57: {
		Name:           "SRE",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.val = c.LSR(c.val)
				return false
			case 4:
				c.mem.Write(uint16(c.pointer), c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)

				return true
			default:
				fmt.Println("something wrong at 0x57")
				return true
			}
		},
	},

	0x58: {
		Name:           "CLI",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(Interrupt)
			return true
		},
	},
	0x59: {
		Name:           "EOR",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0x59")
				return true
			}
		},
	},

	0x5A: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x5b: {
		Name:           "SRE ABS Y",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.LSR(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("somethign wrong at 0x5b")
				return true
			}
		},
	},

	0x5C: {
		Name:           "NOP ABS X",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				newAddr := baseAddr + uint16(c.X)

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				return true
			default:
				fmt.Println("wrong at 0x5C")
				return true
			}
		},
	},

	0x5D: {
		Name:           "EOR",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.A ^= val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0x5D")
				return true
			}
		},
	}, 0x5E: {
		Name:           "LSR",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)
				return false
			case 3:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				var carry bool
				c.temp.val, carry = performLSR(c.temp.val)
				c.updateFlag(Carry, carry)
				c.SetFlagNZ(c.temp.val)
				return false
			case 5:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("bad at 0x5E")
				return true
			}
		},
	},

	0x5F: {
		Name:           "SRE",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.X)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.LSR(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)

				c.A ^= c.val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("somethign wrong at 0x5F")
				return true
			}
		},
	},

	// 0x6X Series

	0x60: {
		Name:           "RTS",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				return false
			case 2:
				c.temp.low = c.popStack()
				return false
			case 3:
				c.temp.high = c.popStack()
				c.PC = builduint16(c.temp.low, c.temp.high)
				return false
			case 4:
				c.PC++
				return true
			default:
				fmt.Println("somethign bat at 0x60")
				return true
			}
		},
	},
	0x61: {
		Name:           "ADC",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.temp.pointer += c.X

				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				c.temp.val = c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.ADC(c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0x61")
				return true
			}
		},
	},

	0x62: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x63: {
		Name:           "RRA",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 4:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				c.val = c.ROR(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)
				c.ADC(c.val)
				return true
			default:
				fmt.Println("something bad at 0x63")
				return true
			}
		},
	},

	0x64: {
		Name:           "NOP",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.fetchone()
				return false
			case 1:
				return true
			default:
				fmt.Println("something wrong at 0x04")
				return false
			}
		},
	},

	0x65: {
		Name:           "ADC",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.ADC(val)
				return true
			default:
				fmt.Print("somethign bad at 0x65")
				return true
			}
		},
	},
	0x66: {
		Name:           "ROR",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.val = c.ROR(c.temp.val)
				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("bad at 0x66")
				return true
			}
		},
	},

	0x67: {
		Name:           "RRA",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.val = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.val = c.ROR(c.val)
				return false
			case 3:
				c.mem.Write(uint16(c.pointer), c.val)

				c.ADC(c.val)
				return true
			default:
				fmt.Println("something wrong at 0x67")
				return true

			}
		},
	},

	0x68: {
		Name:           "PLA",
		AddressingMode: Accumulator,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				return false
			case 2:
				c.A = c.popStack()
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Print("something wrong at 0x68")
				return true
			}
		},
	},
	0x69: {
		Name:           "ADC",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.ADC(val)
			return true
		},
	},
	0x6A: {
		Name:           "ROR",
		AddressingMode: Accumulator,
		Size:           1,
		Execute: func(c *cpu, step int) bool {

			c.A = c.ROR(c.A)
			return true
		},
	},

	0x6B: {
		Name:           "ARR",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			mid := c.A & val

			c.A = c.ROR(mid)

			c.updateFlag(Carry, getbitBool(c.A, 6))
			c.updateFlag(oVerflow, (getbit(c.A, 6)^getbit(c.A, 5)) == 1)
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x6C: {
		Name:           "JMP IND",
		AddressingMode: Indirect,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high)

				c.temp.val = c.mem.Read(c.temp.addr)

				return false
			case 3:
				var high uint8
				if c.temp.low == 0xFF {
					high = c.mem.Read(c.temp.addr & 0xFF00)
				} else {
					high = c.mem.Read(c.temp.addr + 1)
				}

				c.PC = builduint16(c.temp.val, high)
				return true
			default:
				fmt.Println("something wrong with 0x6C")
				return true
			}
		},
	},
	0x6D: {
		Name:           "ADC",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.ADC(val)
				return true
			default:
				fmt.Println("something bad at 0x6D")
				return true
			}
		},
	},
	0x6E: {
		Name:           "ROR",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high)
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 3:
				c.temp.val = c.ROR(c.temp.val)
				return false
			case 4:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0x6E")
				return true
			}
		},
	},

	0x6F: {
		Name:           "RRA",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				return false
			case 3:
				c.val = c.ROR(c.val)
				return false
			case 4:
				c.mem.Write(c.addr, c.val)

				c.ADC(c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0x6F")
				return true
			}
		},
	},

	//0x7X Series

	0x70: {
		Name:           "BVS",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()
				if !c.getFlag(oVerflow) {
					return true
				}
				return false
			case 1:
				offset := int8(c.temp.val)
				oldPC := c.PC

				c.PC = uint16(int32(oldPC) + int32(offset))

				if !crossedPage(oldPC, c.PC) {
					return true
				} else {
					return false
				}

			case 2:
				return true
			default:
				fmt.Println("bad at 0x70")
				return true
			}
		},
	},
	0x71: {
		Name:           "ADC",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				baseaddr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = baseaddr + uint16(c.Y)

				if crossedPage(baseaddr, c.temp.addr) {
					return false
				}

				fallthrough
			case 4:
				c.temp.val = c.mem.Read(c.temp.addr)
				c.ADC(c.temp.val)
				return true

			default:
				fmt.Println("something bad at 0x71")
				return true
			}
		},
	},

	0x72: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0x73: {
		Name:           "RRA IND Y",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 4:
				c.val = c.mem.Read(c.addr)
				return false
			case 5:
				c.val = c.ROR(c.val)
				return false
			case 6:
				c.mem.Write(c.addr, c.val)
				c.ADC(c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0x73")
				return true
			}
		},
	},

	0x74: {
		Name:           "NOP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				return true
			default:
				fmt.Println("something really wrong at 0x14")
				return true
			}
		},
	},

	0x75: {
		Name:           "ADC",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.ADC(val)
				return true
			default:
				fmt.Println("bad at 0x75")
				return true
			}
		},
	},
	0x76: {
		Name:           "ROR",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.val = c.ROR(c.temp.val)
				return false
			case 4:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("bad at 0x76")
				return true
			}
		},
	},

	0x77: {
		Name:           "RRA",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.val = c.ROR(c.val)
				return false
			case 4:
				c.mem.Write(uint16(c.pointer), c.val)
				c.ADC(c.val)
				return true
			default:
				fmt.Println("something really wrong at 0x77")
				return true
			}
		},
	},

	0x78: {
		Name:           "SEI",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.setFlag(Interrupt)
			return true
		},
	},
	0x79: {
		Name:           "ADC",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = baseAddr + uint16(c.Y)

				if crossedPage(baseAddr, c.temp.addr) {
					return false

				}
				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.ADC(val)
				return true
			default:
				fmt.Println("bad at 0x79")
				return true
			}
		},
	},

	0x7A: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x7B: {
		Name:           "RRA",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 3:
				c.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				c.val = c.ROR(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				c.ADC(c.val)
				return true
			default:
				fmt.Println("something wrong at 0x7B")
				return true
			}
		},
	},

	0x7C: {
		Name:           "NOP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				newAddr := baseAddr + uint16(c.X)

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				return true
			default:
				fmt.Println("wrong at 0x7C")
				return true
			}
		},
	},

	0x7D: {
		Name:           "ADC",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = baseAddr + uint16(c.X)

				if crossedPage(baseAddr, c.temp.addr) {
					return false

				}
				fallthrough
			case 3:
				val := c.mem.Read(c.temp.addr)
				c.ADC(val)
				return true
			default:
				fmt.Println("bad at 0x7D")
				return true
			}
		},
	},
	0x7E: {
		Name:           "ROR",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high) + uint16(c.X)
				return false
			case 3:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				c.temp.val = c.ROR(c.temp.val)
				return false
			case 5:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("somethign bad at 0x7E")
				return true
			}
		},
	},

	0x7F: {
		Name:           "RRA",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.X)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				return false
			case 4:
				c.val = c.ROR(c.val)
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				c.ADC(c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0x7F")
				return true
			}
		},
	},

	0x80: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.fetchone()
			return true
		},
	},

	0x81: {
		Name:           "STA",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				adrr := builduint16(c.temp.low, c.temp.high)
				c.mem.Write(adrr, c.A)
				return true
			default:
				fmt.Println("something bad at 0x81")
				return true
			}
		},
	},

	0x82: {
		Name:           "NOP",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.fetchone()
			return true
		},
	},

	0x83: {
		Name:           "SAX",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.pointer += c.X

				return false
			case 2:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				c.val = c.A & c.X
				return false
			case 4:
				c.mem.Write(builduint16(c.low, c.high), c.val)
				return true
			default:
				fmt.Println("something wrong at 0x83")
				return true
			}
		},
	},

	0x84: {
		Name:           "STY",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Write(uint16(c.temp.pointer), c.Y)
				return true
			default:
				fmt.Println("something bad at 0x84")
				return true
			}
		},
	},
	0x85: {
		Name:           "STA",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Write(uint16(c.temp.pointer), c.A)
				return true
			default:
				fmt.Println("something bad at 0x85")
				return true
			}
		},
	},
	0x86: {
		Name:           "STX",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Write(uint16(c.temp.pointer), c.X)
				return true
			default:
				fmt.Println("something bad at 0x86")
				return true
			}
		},
	},

	0x87: {
		Name:           "SAX",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.val = c.A & c.X
				c.mem.Write(uint16(c.pointer), c.val)
				return true
			default:
				fmt.Println("something wrong at 0x87")
				return true
			}
		},
	},

	0x88: {
		Name:           "DEY",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.Y--
			c.SetFlagNZ(c.Y)
			return true
		},
	},

	0x89: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.fetchone()
			return true
		},
	},

	0x8A: {
		Name:           "TXA",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.A = c.X
			c.SetFlagNZ(c.A)
			return true
		},
	},

	0x8B: {
		Name:           "ANE",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {

			imm := c.fetchone()

			const magic uint8 = 0xEE

			c.A = (c.A | magic) & c.X & imm
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x8C: {
		Name:           "STY",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.mem.Write(addr, c.Y)
				return true
			default:
				fmt.Println("Something bad at 0x8C")
				return true
			}
		},
	},
	0x8D: {
		Name:           "STA",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.mem.Write(addr, c.A)
				return true
			default:
				fmt.Println("something bad at 0x8D")
				return true
			}

		},
	},
	0x8E: {
		Name:           "STX",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.mem.Write(addr, c.X)
				return true
			default:
				fmt.Println("something bad at 0x8E")
				return true
			}

		},
	},
	0x8F: {
		Name:           "SAX",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.val = c.A & c.X
				c.mem.Write(c.addr, c.val)
				return true
			default:
				fmt.Println("soemthing wrong at 0x8F")
				return true
			}
		},
	},

	//0x9X Series

	0x90: {
		Name:           "BCC",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()

				if c.getFlag(Carry) {
					return true
				}

				return false

			case 1:
				offset := int8(c.temp.val)
				oldPc := c.PC

				c.PC = uint16(int32(oldPc) + int32(offset))

				if crossedPage(oldPc, c.PC) {
					return false
				} else {
					return true
				}
			case 2:
				return true
			default:
				fmt.Println("something bad 0x90")
				return false
			}
		},
	},
	0x91: {
		Name:           "STA",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				base := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = base + uint16(c.Y)
				return false
			case 4:
				c.mem.Write(c.temp.addr, c.A)
				return true

			default:
				fmt.Println("something bad at 0x91")
				return true
			}
		},
	},

	0x92: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0x93: {
		Name:           "SHA",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				c.mem.Read(builduint16(c.low, c.high))
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 4:
				highbyteModifier := uint8(c.high + 1)
				c.temp.val = c.A & c.X & highbyteModifier

				if (uint16(c.low) + uint16(c.Y)) > 0xFF {
					c.temp.addr = (uint16(c.temp.val) << 8) | (c.temp.addr & 0x00FF)
				}
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				return true
			default:
				fmt.Println("something wrong at 0x93")
				return true
			}
		},
	},

	0x94: {
		Name:           "STY",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.mem.Write(uint16(c.temp.pointer), c.Y)
				return true
			default:
				fmt.Print("something bad at 0x94")
				return true
			}
		},
	},
	0x95: {
		Name:           "STA",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.mem.Write(uint16(c.temp.pointer), c.A)
				return true
			default:
				fmt.Print("something bad at 0x95")
				return true
			}
		},
	},
	0x96: {
		Name: "STX ZP Y",
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.Y
				return false
			case 2:
				c.mem.Write(uint16(c.temp.pointer), c.X)
				return true
			default:
				fmt.Print("something bad at 0x96")
				return true
			}
		},
	},

	0x97: {
		Name:           "SAX",
		AddressingMode: ZeroPageY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.Y
				return false
			case 2:
				c.val = c.A & c.X
				c.mem.Write(uint16(c.pointer), c.val)
				return true
			default:
				fmt.Println("something wrong at 0x97")
				return true
			}
		},
	},
	0x98: {
		Name:           "TYA",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.A = c.Y
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x99: {
		Name:           "STA",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)
				return false
			case 3:
				c.mem.Write(c.temp.addr, c.A)
				return true
			default:
				fmt.Println("bad at 0x99")
				return true
			}
		},
	},
	0x9A: {
		Name:           "TXS",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.S = c.X
			return true
		},
	},

	0x9B: {
		Name:           "TAS",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.mem.Read(c.addr)
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				c.S = c.A & c.X
				highbyteModifier := c.high + 1
				c.temp.val = c.S & highbyteModifier

				if (uint16(c.temp.low) + uint16(c.Y)) > 0xFF {
					c.addr = (uint16(c.temp.val) << 8) | (c.addr & 0x00FF)
				}

				return false
			case 3:
				c.mem.Write(c.addr, c.val)
				return true
			default:
				fmt.Println("something wrong at 0x9B")
				return true
			}
		},
	},

	0x9C: {
		Name:           "SHY",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				c.addr = builduint16(c.low, c.high)
				return false
			case 2:
				eff := c.low + c.X
				pageCross := eff < c.low

				var high uint16 = uint16(c.high)

				if pageCross {
					high = uint16(c.high) + uint16(c.X)
				}

				dummyAddr := (high << 8) | uint16(eff)
				c.mem.Read(dummyAddr)
				return false
			case 3:
				effLow := c.low + c.X

				valToStore := c.Y & (c.high + 1)

				var targetHigh uint16 = uint16(c.high)
				if effLow < c.low {
					targetHigh = uint16(valToStore)
				}

				targetAddr := (targetHigh << 8) | uint16(effLow)

				c.mem.Write(targetAddr, valToStore)
				return true
			default:
				fmt.Println("soemthign wrong at 0x9C")
				return true
			}
		},
	},

	0x9D: {
		Name:           "STA",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)
				return false
			case 3:
				c.mem.Write(c.temp.addr, c.A)
				return true
			default:
				fmt.Println("bad at 0x9D")
				return true
			}
		},
	},

	0x9E: {
		Name:           "SHX",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				c.addr = builduint16(c.low, c.high)
				return false
			case 2:
				eff := c.low + c.Y
				pageCross := eff < c.low

				var high uint16 = uint16(c.high)

				if pageCross {
					high = uint16(c.high) + uint16(c.Y)
				}

				dummyAddr := (high << 8) | uint16(eff)
				c.mem.Read(dummyAddr)
				return false
			case 3:
				effLow := c.low + c.Y

				valToStore := c.X & (c.high + 1)

				var targetHigh uint16 = uint16(c.high)
				if effLow < c.low {
					targetHigh = uint16(valToStore)
				}

				targetAddr := (targetHigh << 8) | uint16(effLow)

				c.mem.Write(targetAddr, valToStore)
				return true
			default:
				fmt.Println("soemthign wrong at 0x9C")
				return true
			}
		},
	},

	0x9F: {
		Name:           "SHA",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				c.addr = builduint16(c.low, c.high)
				return false
			case 2:
				eff := c.low + c.Y
				pageCross := eff < c.low

				var high uint16 = uint16(c.high)

				if pageCross {
					high = uint16(c.high) + uint16(c.Y)
				}

				dummyAddr := (high << 8) | uint16(eff)
				c.mem.Read(dummyAddr)
				return false
			case 3:
				effLow := c.low + c.Y

				valToStore := c.A & (c.high + 1)

				var targetHigh uint16 = uint16(c.high)
				if effLow < c.low {
					targetHigh = uint16(valToStore)
				}

				targetAddr := (targetHigh << 8) | uint16(effLow)

				c.mem.Write(targetAddr, valToStore)
				return true
			default:
				fmt.Println("soemthign wrong at 0x9C")
				return true
			}
		},
	},

	0xA0: {
		Name:           "LDY",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.Y = c.fetchone()
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xA1: {
		Name:           "LDA IND X",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				addr := builduint16(c.temp.low, c.temp.high)
				val := c.mem.Read(addr)

				c.A = val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0xA1")
				return true
			}
		},
	},
	0xA2: {
		Name:           "LDX",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.X = c.fetchone()
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xA3: {
		Name:           "LAX",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.pointer += c.X
				return false
			case 2:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 4:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)

				c.A = c.val
				c.X = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("something wrontt at 0xA3")
				return true
			}
		},
	},
	0xA4: {
		Name:           "LDY",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.Y = val
				c.SetFlagNZ(c.Y)
				return true
			default:
				fmt.Println("bad at 0xA4")
				return true
			}
		},
	},
	0xA5: {
		Name:           "LDA",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.A = val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0xA5")
				return true
			}
		},
	},
	0xA6: {
		Name:           "LDX",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.X = val
				c.SetFlagNZ(c.X)
				return true
			default:
				fmt.Println("bad at 0xA6")
				return true
			}
		},
	},

	0xA7: {
		Name:           "LAX",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.val = c.mem.Read(uint16(c.pointer))

				c.A = c.val
				c.X = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("somethign wrong at A7")
				return true
			}
		},
	},

	0xA8: {
		Name:           "TAY",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.Y = c.A
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xA9: {
		Name:           "LDA",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.A = c.fetchone()
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0xAA: {
		Name:           "TAX",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.X = c.A

			c.SetFlagNZ(c.X)

			return true
		},
	},

	0xAB: {
		Name:           "LXA",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			imm := c.fetchone()

			const magic uint8 = 0xEE

			res := (c.A | magic) & imm
			c.A = res
			c.X = res

			c.SetFlagNZ(res)
			return true

		},
	},

	0xAC: {
		Name:           "LDY",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.Y = val
				c.SetFlagNZ(c.Y)
				return true
			default:
				fmt.Println("something bad at 0xAC")
				return true
			}
		},
	},
	0xAD: {
		Name:           "LDA",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.A = val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("something bad at 0xAD")
				return true
			}
		},
	},
	0xAE: {
		Name:           "LDY",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.X = val
				c.SetFlagNZ(c.X)
				return true
			default:
				fmt.Println("something bad at 0xAE")
				return true
			}
		},
	},

	0xAF: {
		Name:           "LAX",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.val = c.mem.Read(builduint16(c.low, c.high))

				c.A = c.val
				c.X = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("something wrong at AF")
				return true
			}
		},
	},

	0xB0: {
		Name:           "BCS",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()

				if !c.getFlag(Carry) {
					return true
				}
				return false
			case 1:
				offset := int8(c.temp.val)
				oldPC := c.PC

				c.PC = uint16(int32(oldPC) + int32(offset))

				if crossedPage(oldPC, c.PC) {
					return false
				} else {
					return true
				}
			case 2:
				return true
			default:
				fmt.Println("something wrong at 0xB0")
				return true
			}
		},
	},
	0xB1: {
		Name:           "LDA",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 4:
				val := c.mem.Read(c.temp.addr)
				c.A = val
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("seomthng bat at 0xB1")
				return true
			}
		},
	},
	0xB2: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0xB3: {
		Name:           "LAX",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:

				baseAddr := builduint16(c.low, c.high)
				c.mem.Read(baseAddr)
				c.addr = baseAddr + uint16(c.Y)

				if crossedPage(baseAddr, c.addr) {
					return false
				}

				fallthrough
			case 4:
				c.val = c.mem.Read(c.addr)

				c.A = c.val
				c.X = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("someting wrong at 0xb3")
				return true
			}
		},
	},

	0xB4: {
		Name:           "LDY",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.Y = c.mem.Read(uint16(c.temp.pointer))
				c.SetFlagNZ(c.Y)
				return true
			default:
				fmt.Print("something bad at b4")
				return true
			}
		},
	},
	0xB5: {
		Name:           "LDA",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.A = c.mem.Read(uint16(c.temp.pointer))
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Print("something bad at b5")
				return true
			}
		},
	},
	0xB6: {
		Name:           "LDX",
		AddressingMode: ZeroPageY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.Y
				return false
			case 2:
				c.X = c.mem.Read(uint16(c.temp.pointer))
				c.SetFlagNZ(c.X)
				return true
			default:
				fmt.Print("something bad at b6")
				return true
			}
		},
	},

	0xB7: {
		Name:           "LAX",
		AddressingMode: ZeroPageY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.pointer += c.Y
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))

				c.A = c.val
				c.X = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("somthing wrong at b7")
				return true
			}
		},
	},

	0xB8: {
		Name:           "CLV",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(oVerflow)
			return true
		},
	},
	0xB9: {
		Name:           "LDA",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 3:
				c.A = c.mem.Read(c.temp.addr)
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0xb9")
				return true
			}
		},
	},
	0xBA: {
		Name:           "TSX",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.X = c.S
			c.SetFlagNZ(c.X)
			return true
		},
	},

	0xBB: {
		Name:           "LAS",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				c.mem.Read(baseAddr)
				c.addr = baseAddr + uint16(c.Y)

				if crossedPage(baseAddr, c.addr) {
					return false

				}

				fallthrough
			case 3:
				c.val = c.mem.Read(c.addr)

				c.val &= c.S

				c.A = c.val
				c.X = c.val
				c.S = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("something wrong at 0xbb")
				return true
			}
		},
	},

	0xBC: {
		Name:           "LDY",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 3:
				c.Y = c.mem.Read(c.temp.addr)
				c.SetFlagNZ(c.Y)
				return true
			default:
				fmt.Println("bad at 0xbC")
				return true
			}
		},
	},
	0xBD: {
		Name:           "LDA",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 3:
				c.A = c.mem.Read(c.temp.addr)
				c.SetFlagNZ(c.A)
				return true
			default:
				fmt.Println("bad at 0xb9")
				return true
			}
		},
	},
	0xBE: {
		Name:           "LDX",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 3:
				c.X = c.mem.Read(c.temp.addr)
				c.SetFlagNZ(c.X)
				return true
			default:
				fmt.Println("bad at 0xb9")
				return true
			}
		},
	},
	0xBF: {
		Name:           "LAX",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				c.mem.Read(baseAddr)
				c.addr = baseAddr + uint16(c.Y)

				if crossedPage(baseAddr, c.addr) {
					return false
				}

				fallthrough
			case 3:
				c.val = c.mem.Read(c.addr)

				c.A = c.val
				c.X = c.val

				c.SetFlagNZ(c.val)
				return true
			default:
				fmt.Println("something wrong at 0xBF")
				return true
			}
		},
	},
	0xC0: {
		Name:           "CMY",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.COMPARE(c.Y, val)
			return true
		},
	},
	0xC1: {
		Name:           "CMP",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.COMPARE(c.A, val)
				return true
			default:
				fmt.Println("bad at 0xC1")
				return true
			}
		},
	},
	0xC2: {
		Name:           "NOP",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.fetchone()
			return true
		},
	},

	0xC3: {
		Name:           "DCP",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.pointer += c.X
				return false
			case 2:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 4:
				c.addr = builduint16(c.low, c.high)
				c.mem.Read(c.addr)
				c.val = c.mem.Read(c.addr)

				c.val--
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				return false
			case 6:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("something wrong at 0xC3")
				return true
			}
		},
	},

	0xC4: {
		Name:           "CPY",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.COMPARE(c.Y, val)
				return true
			default:
				fmt.Println("bad at 0xc4")
				return true
			}
		},
	},
	0xC5: {
		Name:           "CPM",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				val := c.mem.Read(uint16(c.temp.pointer))
				c.COMPARE(c.A, val)
				return true
			default:
				fmt.Println("bad at 0xc5")
				return true
			}
		},
	},
	0xC6: {
		Name:           "DEC",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))

				return false
			case 2:
				c.temp.val--
				c.SetFlagNZ(c.temp.val)
				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("bad at 0xc6")
				return true
			}
		},
	},

	0xC7: {
		Name:           "DCP",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.val = c.mem.Read(uint16(c.pointer))
				c.val--
				return false
			case 2:
				c.mem.Write(uint16(c.pointer), c.val)
				return false
			case 3:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("soemethign wrong at 0xC7")
				return true
			}
		},
	},

	0xC8: {
		Name:           "INY",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.Y++
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xC9: {
		Name:           "CMP",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.COMPARE(c.A, c.fetchone())
			return true
		},
	},
	0xCA: {
		Name:           "DEX",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.X--
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xCB: {
		Name:           "SBX",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()

			temp := c.A & c.X

			res := uint16(temp) - uint16(val)

			c.X = uint8(res)

			c.updateFlag(Carry, temp >= val)
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xCC: {
		Name:           "CPY",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.COMPARE(c.Y, val)
				return true
			default:
				fmt.Println("seomtgin bad at 0xCC")
				return true
			}
		},
	},
	0xCD: {
		Name:           "CPM",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.COMPARE(c.A, val)
				return true
			default:
				fmt.Println("seomtgin bad at 0xCD")
				return true
			}
		},
	},
	0xCE: {
		Name:           "DEC",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.addr = builduint16(c.temp.low, c.temp.high)
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 3:
				c.temp.val--
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("seomtgin bad at 0xCE")
				return true
			}
		},
	},
	0xCF: {
		Name:           "DCP",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				c.val--
				return false
			case 3:
				c.mem.Write(c.addr, c.val)
				return false
			case 4:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0xCF")
				return true
			}
		},
	},
	0xD0: {
		Name:           "BNE",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()

				if c.getFlag(Zero) {
					return true
				}

				return false
			case 1:
				offset := int8(c.temp.val)
				oldPc := c.PC

				c.PC = uint16(int32(oldPc) + int32(offset))

				if crossedPage(oldPc, c.PC) {
					return false
				} else {
					return true
				}
			case 2:
				return true
			default:
				fmt.Println("something bad at 0xD0")
				return true
			}
		},
	},
	0xD1: {
		Name:           "CMP",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}

				fallthrough
			case 4:
				val := c.mem.Read(c.temp.addr)
				c.COMPARE(c.A, val)

				return true
			default:
				fmt.Println("something bad at 0xD1")
				return true
			}
		},
	},

	0xD2: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0xD3: {
		Name:           "DCP",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				c.mem.Read(c.addr)
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 4:
				c.val = c.mem.Read(c.addr)
				c.val--
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				return false
			case 6:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0xD3")
				return true
			}
		},
	},

	0xD4: {
		Name:           "NOP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.mem.Read(uint16(c.pointer))
				return true
			default:
				fmt.Println("something wrong at 0xD4")
				return true

			}
		},
	},

	0xD5: {
		Name:           "CMP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.COMPARE(c.A, c.mem.Read(uint16(c.temp.pointer)))
				return true
			default:
				fmt.Println("somethign bad at 0xd5")
				return true
			}
		},
	},
	0xD6: {
		Name:           "DEC",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.val--
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("somethign bad at 0xd6")
				return true
			}
		},
	},

	0xD7: {
		Name:           "DCP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				c.val--
				return false
			case 3:
				c.mem.Write(uint16(c.pointer), c.val)
				return false
			case 4:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("something wrong at 0xd7")
				return true
			}
		},
	},

	0xD8: {
		Name:           "CLD",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(Decimal)
			return true
		},
	},
	0xD9: {
		Name:           "CMP",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}
				fallthrough
			case 3:
				c.COMPARE(c.A, c.mem.Read(c.temp.addr))
				return true
			default:
				fmt.Println("somethign bad at 0xd9")
				return true
			}
		},
	},
	0xDA: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {

			return true
		},
	},

	0xDB: {
		Name:           "DCP",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				c.val--
				return false
			case 4:
				c.mem.Write(c.addr, c.val)
				return false
			case 5:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0xDB")
				return true
			}
		},
	},

	0xDC: {
		Name:           "NOP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				newAddr := baseAddr + uint16(c.X)

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				return true
			default:
				fmt.Println("wrong at 0xDC")
				return true
			}
		},
	},

	0xDD: {
		Name:           "CMP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)

				if crossedPage(addr, c.temp.addr) {
					return false
				}
				fallthrough
			case 3:
				c.COMPARE(c.A, c.mem.Read(c.temp.addr))
				return true
			default:
				fmt.Println("somethign bad at 0xdd")
				return true
			}
		},
	},
	0xDE: {
		Name:           "DEC",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)
				return false
			case 3:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				c.temp.val--
				c.SetFlagNZ(c.temp.val)
				return false
			case 5:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("somethign bad at 0xd9")
				return true
			}
		},
	},

	0xDF: {
		Name:           "DCP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.X)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				c.val--
				return false
			case 4:
				c.mem.Write(c.addr, c.val)
				return false
			case 5:
				c.COMPARE(c.A, c.val)
				return true
			default:
				fmt.Println("something wrong at 0xDF")
				return false
			}
		},
	},

	0xE0: {
		Name:           "CPX",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.COMPARE(c.X, c.fetchone())
			return true
		},
	},
	0xE1: {
		Name:           "SBC",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 4:
				val := c.mem.Read(builduint16(c.temp.low, c.temp.high))
				c.SBC(val)
				return true
			default:
				fmt.Println("bad at 0xE1")
				return true
			}
		},
	},
	0xE2: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.fetchone()
			return true
		},
	},

	0xE3: {
		Name:           "ISC",
		AddressingMode: IndirectX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.mem.Read(uint16(c.pointer))
				c.pointer += c.X
				return false
			case 2:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 3:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 4:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				c.val++
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				return false
			case 6:
				c.SBC(c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0xE3")
				return true
			}
		},
	},
	0xE4: {
		Name:           "CPX",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.COMPARE(c.X, c.mem.Read(uint16(c.temp.pointer)))
				return true
			default:
				fmt.Println("soemthin bad at 0xE4")
				return true
			}
		},
	},
	0xE5: {
		Name:           "SBC",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.SBC(c.mem.Read(uint16(c.temp.pointer)))
				return true
			default:
				fmt.Println("soemthin bad at 0xE4")
				return true
			}
		},
	},
	0xE6: {
		Name:           "INC",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.val++
				c.SetFlagNZ(c.temp.val)
				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("soemthin bad at 0xE4")
				return true
			}
		},
	},

	0xE7: {
		Name:           "ISC",
		AddressingMode: ZeroPage,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.val = c.mem.Read(uint16(c.pointer))
				c.val++
				return false
			case 2:
				c.mem.Write(uint16(c.pointer), c.val)
				return false
			case 3:
				c.SBC(c.val)
				return true
			default:
				fmt.Println("something wrong at 0xE7")
				return true
			}
		},
	},

	0xE8: {
		Name:           "INX",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.X++
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xE9: {
		Name:           "SBC",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.SBC(c.fetchone())
			return true
		},
	}, 0xEA: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0xEB: {
		Name:           "USBC",
		AddressingMode: Immediate,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			c.SBC(c.fetchone())
			return true
		},
	},

	0xEC: {
		Name:           "CPX ABS",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.COMPARE(c.X, c.mem.Read(builduint16(c.temp.low, c.temp.high)))
				return true
			default:
				fmt.Println("bad at 0xEC")
				return true
			}
		},
	},
	0xED: {
		Name:           "SBC",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.SBC(c.mem.Read(builduint16(c.temp.low, c.temp.high)))
				return true
			default:
				fmt.Println("bad at 0xED")
				return true
			}
		},
	},
	0xEE: {
		Name:           "INC",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				c.temp.val = c.mem.Read(builduint16(c.temp.low, c.temp.high))
				return false
			case 3:
				c.temp.val++
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(builduint16(c.temp.low, c.temp.high), c.temp.val)
				return true
			default:
				fmt.Println("bad at 0xED")
				return true
			}
		},
	},
	0xEF: {
		Name:           "ISC",
		AddressingMode: Absolute,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high)
				c.val = c.mem.Read(c.addr)
				c.val++
				return false
			case 3:
				c.mem.Write(c.addr, c.val)
				return false
			case 4:
				c.SBC(c.val)
				return true
			default:
				fmt.Println("somethign wrong at 0xEF")
				return true
			}
		},
	},
	0xF0: {
		Name:           "BEQ",
		AddressingMode: Relative,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.val = c.fetchone()

				if !c.getFlag(Zero) {
					return true
				}

				return false
			case 1:
				offset := int8(c.temp.val)
				oldPc := c.PC

				c.PC = uint16(int32(oldPc) + int32(offset))

				if crossedPage(oldPc, c.PC) {
					return false
				} else {
					return true
				}
			case 2:
				return true
			default:
				fmt.Println("something bad at 0xF0")
				return true
			}
		},
	},
	0xF1: {
		Name:           "SBC",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer + 1))
				return false
			case 3:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}
				fallthrough
			case 4:
				c.SBC(c.mem.Read(c.temp.addr))
				return true
			default:
				fmt.Println("something bad at 0xf1")
				return true
			}
		},
	},
	0xF2: {
		Name:           "JAM",
		AddressingMode: Implied,
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0xF3: {
		Name:           "ISC",
		AddressingMode: IndirectY,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.low = c.mem.Read(uint16(c.pointer))
				return false
			case 2:
				c.high = c.mem.Read(uint16(c.pointer + 1))
				return false
			case 3:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 4:
				c.val = c.mem.Read(c.addr)
				c.val++
				return false
			case 5:
				c.mem.Write(c.addr, c.val)
				return false
			case 6:
				c.SBC(c.val)
				return true
			default:
				fmt.Println("something wrong at 0xF3")
				return true
			}
		},
	},

	0xF4: {
		Name:           "NOP",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				return true
			default:
				fmt.Println("something wrong at 0xF4")
				return true
			}
		},
	},

	0xF5: {
		Name:           "SBC",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.SBC(c.mem.Read(uint16(c.pointer)))
				return true
			default:
				fmt.Println("something bad at 0xF5")
				return true
			}
		},
	},
	0xF6: {
		Name:           "INC",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.pointer += c.X
				return false
			case 2:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 3:
				c.temp.val++
				c.SetFlagNZ(c.temp.val)
				return false
			case 4:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0xF6")
				return true
			}
		},
	},

	0xF7: {
		Name:           "ISC",
		AddressingMode: ZeroPageX,
		Size:           2,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.pointer = c.fetchone()
				return false
			case 1:
				c.pointer += c.X
				return false
			case 2:
				c.val = c.mem.Read(uint16(c.pointer))
				c.val++
				return false
			case 3:
				c.mem.Write(uint16(c.pointer), c.val)
				return false
			case 4:
				c.SBC(c.val)
				return true
			default:
				fmt.Println("something wrong at 0xF7")
				return true
			}
		},
	},

	0xF8: {
		Name:           "SED",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			c.setFlag(Decimal)
			return true
		},
	},
	0xF9: {
		Name:           "SBC",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.Y)

				if crossedPage(addr, c.temp.addr) {
					return false
				}
				fallthrough
			case 3:
				c.SBC(c.mem.Read(c.temp.addr))
				return true
			default:
				fmt.Println("something bad at 0xF9")
				return true
			}
		},
	},

	0xFA: {
		Name:           "NOP",
		AddressingMode: Implied,
		Size:           1,
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0xFB: {
		Name:           "ISC",
		AddressingMode: AbsoluteY,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				c.addr = builduint16(c.low, c.high) + uint16(c.Y)
				return false
			case 3:
				c.val = c.mem.Read(c.addr)
				c.val++
				return false
			case 4:
				c.mem.Write(c.addr, c.val)
				return false
			case 5:
				c.SBC(c.val)
				return true
			default:
				fmt.Println("something wrong at 0xFB")
				return true
			}
		},
	},

	0xFC: {
		Name:           "NOP",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.high = c.fetchone()
				return false
			case 2:
				baseAddr := builduint16(c.low, c.high)
				newAddr := baseAddr + uint16(c.X)

				if crossedPage(baseAddr, newAddr) {
					return false
				}

				fallthrough
			case 3:
				return true
			default:
				fmt.Println("wrong at 0xFC")
				return true
			}
		},
	},

	0xFD: {
		Name:           "SBC",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:

				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)

				if crossedPage(addr, c.temp.addr) {
					return false
				}
				fallthrough
			case 3:
				c.SBC(c.mem.Read(c.temp.addr))
				return true
			default:
				fmt.Println("something bad at 0xFD")
				return true
			}
		},
	},
	0xFE: {
		Name:           "INC",
		AddressingMode: AbsoluteX,
		Size:           3,
		Execute: func(c *cpu, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				addr := builduint16(c.temp.low, c.temp.high)
				c.temp.addr = addr + uint16(c.X)
				return false
			case 3:
				c.temp.val = c.mem.Read(c.temp.addr)
				return false
			case 4:
				c.temp.val++
				c.SetFlagNZ(c.temp.val)
				return false
			case 5:
				c.mem.Write(c.temp.addr, c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0xF9")
				return true
			}
		},
	},
}
