package main

import (
	"fmt"
)

type opCode struct {
	Name string

	Execute  func(c *cpu, arg operand, step int) bool
	opcode   uint8
	addrMode AddrMode
}

type operand struct {
	mode      uint8 //addr mode
	addr      AddrMode
	pageCross bool
}

var FetchTable = []opCode{

	0xEA: {
		Name: "NOP",
		Execute: func(c *cpu, arg operand, step int) bool {
			return true
		},
	},
	0xE8: {
		Name: "INX",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.X++
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0x01: {
		Name: "ORA (Indirect, X)",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x05: {
		Name: "ORA zeropage",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x06: {
		Name: "ASL Zero Page",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x08: {
		Name: "PHP impl",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x09: {
		Name: "ORA IMM",
		Execute: func(c *cpu, arg operand, step int) bool {
			val := c.fetchone()
			c.A |= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x0A: {
		Name: "ASL A",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x0D: {
		Name: "ORA ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x0E: {
		Name: "ASL ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x10: {
		Name: "BPL",
		Execute: func(c *cpu, arg operand, step int) bool {
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
				newPC := uint16(int32(32) + int32(offset))

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
		Name: "ORA (IND) Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
				c.pageCrossed = crossedPage(baseAddr, c.temp.addr)

				return false
			case 4:
				if c.pageCrossed {
					return false
				}

				fallthrough
			case 5:
				val := c.mem.Read(c.temp.addr)
				c.A |= val
				c.SetFlagNZ(c.A)
				return true
			default:
				return true
			}

		},
	}, 0x15: {
		Name: "ORA oper,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x16: {
		Name: "ASL Oper,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x18: {
		Name: "CLC",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.clearFlag(Carry)
			return true
		},
	},
	0x19: {
		Name: "ORA ABS,Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x1D: {
		Name: "ORA ABS,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x1E: {
		Name: "ASL ABS,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x20: {
		Name: "JSR",
		Execute: func(c *cpu, arg operand, step int) bool {
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
				c.temp.high = c.fetchone()
				c.PC = builduint16(c.temp.low, c.temp.high)
				return true
			default:
				fmt.Println("somethign very wrong happened")
				return true
			}
		},
	}, 0x21: {
		Name: "AND X,IND",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x24: {
		Name: "BITS ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x25: {
		Name: "AND ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x26: {
		Name: "ROL ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.val = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				var carry bool
				c.temp.val, carry = performROL(c.temp.val, c.getFlag(Carry))
				c.updateFlag(Carry, carry)
				return false
			case 3:
				c.mem.Write(uint16(c.temp.pointer), c.temp.val)
				return true
			default:
				fmt.Println("something bad happened at 0x26")
				return true
			}

		},
	}, 0x28: {
		Name: "PLP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "AND #",
		Execute: func(c *cpu, arg operand, step int) bool {
			val := c.fetchone()
			c.A &= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x2A: {
		Name: "ROL A",
		Execute: func(c *cpu, arg operand, step int) bool {
			var carry bool
			c.A, carry = performROL(c.A, c.getFlag(Carry))
			c.updateFlag(Carry, carry)
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x2C: {
		Name: "BITS ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x2D: {
		Name: "AND ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x2E: {
		Name: "ROL ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x30: {
		Name: "BMI ",
		Execute: func(c *cpu, arg operand, step int) bool {
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

				newPc := uint16(int32(oldPC) + int32(offset))

				if crossedPage(oldPC, newPc) {
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
	}, 0x31: {
		Name: "AND (IND),Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x35: {
		Name: "AND ZPG,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x38: {
		Name: "SEC",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.setFlag(Carry)
			return true
		},
	},
	0x39: {
		Name: "AND abs,Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x3D: {
		Name: "AND ABS,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ROL ABS,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x40: {
		Name: "RTI",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "EOR X,ind",
		Execute: func(c *cpu, arg operand, step int) bool {
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
				c.temp.high = c.mem.Read(uint16(c.temp.pointer))
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
	0x45: {
		Name: "EOR ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LSR ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x48: {
		Name: "PHA",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "EOR IMM",
		Execute: func(c *cpu, arg operand, step int) bool {
			val := c.fetchone()
			c.A ^= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x4A: {
		Name: "LSR A",
		Execute: func(c *cpu, arg operand, step int) bool {
			var carry bool
			c.A, carry = performLSR(c.A)
			c.updateFlag(Carry, carry)
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x4D: {
		Name: "EOR ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LSR ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x50: {
		Name: "BVC",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "EOR IND,Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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

				if c.pageCrossed {
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
	0x55: {
		Name: "EOR ZP,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LSR ZP,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x58: {
		Name: "CLI impl",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.clearFlag(Interrupt)
			return true
		},
	},
	0x59: {
		Name: "EOR abs,Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x5D: {
		Name: "EOR abs,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LSR ABS,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x60: {
		Name: "RTI IMPL",
		Execute: func(c *cpu, arg operand, step int) bool {
			switch step {
			case 0:
				return false
			case 1:
				return false
			case 2:
				c.temp.low = c.popStack()
				return false
			case 3:
				c.temp.low = c.popStack()
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
		Name: "ADC IND,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
				c.temp.val = c.mem.Read(builduint16(c.temp.low, c.temp.high))
				return false
			case 5:
				c.ADC(c.temp.val)
				return true
			default:
				fmt.Println("something bad at 0x61")
				return true
			}
		},
	},
	0x65: {
		Name: "ADC ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ROR ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x68: {
		Name: "PLA A",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ADC #",
		Execute: func(c *cpu, arg operand, step int) bool {
			val := c.fetchone()
			c.ADC(val)
			return true
		},
	},
	0x6A: {
		Name: "ROR A",
		Execute: func(c *cpu, arg operand, step int) bool {

			c.A = c.ROR(c.A)
			return true
		},
	},
	0x6C: {
		Name: "JMP IND",
		Execute: func(c *cpu, arg operand, step int) bool {
			switch step {
			case 0:
				c.temp.low = c.fetchone()
				return false
			case 1:
				c.temp.high = c.fetchone()
				return false
			case 2:
				pointer := builduint16(c.temp.low, c.temp.high)

				low := c.mem.Read(pointer)
				var high uint8
				if c.temp.low == 0xFF {
					high = c.mem.Read(pointer & 0xFF00)
				} else {
					high = c.mem.Read(pointer + 1)
				}

				c.PC = builduint16(low, high)
				return true
			default:
				fmt.Println("something wrong with 0x6C")
				return true
			}
		},
	},
	0x6D: {
		Name: "ADC ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ROR ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	}, 0x70: {
		Name: "BVS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ADC IND Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x75: {
		Name: "ADC ZP,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ROR ZP,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x78: {
		Name: "SEI IMPL",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.setFlag(Interrupt)
			return true
		},
	},
	0x79: {
		Name: "ADC ABS Y",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x7D: {
		Name: "ADC ABS X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "ROR ABS,X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x81: {
		Name: "STA IND X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x84: {
		Name: "STY ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "STA ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "STX ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x88: {
		Name: "DEY",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.Y--
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0x8A: {
		Name: "TXA",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.A = c.X
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x8C: {
		Name: "STY",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "STA",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "STX ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0x90: {
		Name: "BCC",
		Execute: func(c *cpu, arg operand, step int) bool {
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
				return false
			default:
				fmt.Println("something bad 0x90")
				return false
			}
		},
	},
	0x91: {
		Name: "STA IND Y",
		Execute: func(c *cpu, arg operand, step int) bool {
			switch step {
			case 0:
				c.temp.pointer = c.fetchone()
				return false
			case 1:
				c.temp.low = c.mem.Read(uint16(c.temp.pointer))
				return false
			case 2:
				c.temp.high = c.mem.Read(uint16(c.temp.pointer))
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
	0x94: {
		Name: "STY ZP X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "STA ZP X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Execute: func(c *cpu, arg operand, step int) bool {
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
				fmt.Print("something bad at 0x96")
				return true
			}
		},
	},
	0x98: {
		Name: "TYA",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.A = c.Y
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x99: {
		Name: "STA ABS Y ",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "TXS",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.S = c.X
			return true
		},
	},
	0x9D: {
		Name: "STA ABS X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0xA0: {
		Name: "LDY IMM",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.Y = c.fetchone()
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xA1: {
		Name: "LDA IND X",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LDX IMM",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.X = c.fetchone()
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xA4: {
		Name: "LDY ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LDA ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LDX ZP",
		Execute: func(c *cpu, arg operand, step int) bool {
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
	0xA8: {
		Name: "TAY",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.Y = c.A
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xA9: {
		Name: "LDA IMM",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.A = c.fetchone()
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0xAA: {
		Name: "TAX",
		Execute: func(c *cpu, arg operand, step int) bool {
			c.X = c.A

			c.SetFlagNZ(c.X)

			return true
		},
	},
	0xAC: {
		Name: "LDY ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LDA ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
		Name: "LDY ABS",
		Execute: func(c *cpu, arg operand, step int) bool {
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
}
