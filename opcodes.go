package main

import (
	"fmt"
)

type opCode struct {
	Name    string
	Execute func(c *cpu, step int) bool
}

var FetchTable = []opCode{
	0x00: {
		Name: "BRK",
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
		Name: "ORA (Indirect, X)",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {

			c.isJamming = true

			return true
		},
	},

	0x03: {
		Name: "SLO",
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
		Name: "NOP ZP",
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
		Name: "ORA zeropage",
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
	}, 0x06: {
		Name: "ASL Zero Page",
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
		Name: "SLO ZPG",
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
		Name: "PHP impl",
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
	}, 0x09: {
		Name: "ORA IMM",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.A |= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x0A: {
		Name: "ASL A",
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
	}, 0x0B: {
		Name: "ANC ",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()

			c.A &= val
			c.updateFlag(Carry, getbitBool(c.A, 7))
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x0C: {
		Name: "NOP ABS",
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
		Name: "ORA ABS",
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
	}, 0x0E: {
		Name: "ASL ABS",
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
		Name: "SLO ABS",
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

	0x10: {
		Name: "BPL",
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
		Name: "ORA (IND) Y",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x13: {
		Name: "SLO IND Y",
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
		Name: "NOP ZP, X",
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
		Name: "ORA oper,X",
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
	}, 0x16: {
		Name: "ASL Oper,X",
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
		Name: "SLO ZP X",
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
		Name: "CLC",
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(Carry)
			return true
		},
	},
	0x19: {
		Name: "ORA ABS,Y",
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
		Name: "NOP IMPL",
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x1B: {
		Name: "SLO ABS Y",
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
		Name: "NOP ABS X",
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
		Name: "ORA ABS,X",
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
	}, 0x1E: {
		Name: "ASL ABS,X",
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
		Name: "SLO ABS X",
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

	0x20: {
		Name: "JSR",
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
	}, 0x21: {
		Name: "AND X,IND",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x23: {
		Name: "RLA X IND ",
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
		Name: "BITS ZP",
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
	}, 0x25: {
		Name: "AND ZP",
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
	}, 0x26: {
		Name: "ROL ZP",
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
		Name: "RLA ZP",
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
		Name: "PLP",
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
		Name: "AND #",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.A &= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x2A: {
		Name: "ROL A",
		Execute: func(c *cpu, step int) bool {
			var carry bool
			c.A, carry = performROL(c.A, c.getFlag(Carry))
			c.updateFlag(Carry, carry)
			c.SetFlagNZ(c.A)
			return true
		},
	},

	0x2B: {
		Name: "ANC ",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()

			c.A &= val
			c.updateFlag(Carry, getbitBool(c.A, 7))
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x2C: {
		Name: "BITS ABS",
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
	}, 0x2D: {
		Name: "AND ABS",
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
	}, 0x2E: {
		Name: "ROL ABS",
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
		Name: "RLA ABS",
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

	0x30: {
		Name: "BMI ",
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
	}, 0x31: {
		Name: "AND (IND),Y",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x33: {
		Name: "RLA IND Y",
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
		Name: "NOP ",
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
		Name: "AND ZPG,X",
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
		Name: "ROL ZPG,X",
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
		Name: "RLA ZPG X",
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
		Name: "SEC",
		Execute: func(c *cpu, step int) bool {
			c.setFlag(Carry)
			return true
		},
	},
	0x39: {
		Name: "AND abs,Y",
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
		Name: "NOP IMPL",
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x3B: {
		Name: "RLA ABS Y",
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
		Name: "NOP ABS X",
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
		Name: "AND ABS,X",
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
		Name: "ROL ABS,X",
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
		Name: "RLA ABS X",
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
	0x40: {
		Name: "RTI",
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
		Name: "EOR X,ind",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x43: {
		Name: "SRE X IND ",
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
		Name: "NOP ZP",
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
		Name: "EOR ZP",
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
		Name: "LSR ZP",
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
		Name: "SRE ZP",
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
		Name: "PHA",
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
		Name: "EOR IMM",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.A ^= val
			c.SetFlagNZ(c.A)
			return true
		},
	}, 0x4A: {
		Name: "LSR A",
		Execute: func(c *cpu, step int) bool {
			var carry bool
			c.A, carry = performLSR(c.A)
			c.updateFlag(Carry, carry)
			c.SetFlagNZ(c.A)
			return true
		},
	},

	0x4B: {
		Name: "ALR",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			mid := c.A & val

			c.A = c.LSR(mid)
			c.SetFlagNZ(c.A)

			return true
		},
	},

	0x4C: {
		Name: "JMP ABS",
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
		Name: "EOR ABS",
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
		Name: "LSR ABS",
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
		Name: "SRE ABS",
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

	0x50: {
		Name: "BVC",
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
		Name: "EOR IND,Y",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true
			return true
		},
	},

	0x53: {
		Name: "SRE IND Y",
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
		Name: "NOP ZP",
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
		Name: "EOR ZP,X",
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
		Name: "LSR ZP,X",
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
		Name: "SRE ZP X",
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
		Name: "CLI impl",
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(Interrupt)
			return true
		},
	},
	0x59: {
		Name: "EOR abs,Y",
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
		Name: "NOP IMPL",
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0x5b: {
		Name: "SRE ABS Y",
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
		Name: "NOP ABS X",
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
		Name: "EOR abs,X",
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
		Name: "LSR ABS,X",
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
		Name: "SRE ABS X",
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

	0x60: {
		Name: "RTS IMPL",
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
		Name: "ADC IND,X",
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
		Name: "JAM",
		Execute: func(c *cpu, step int) bool {
			c.isJamming = true

			return true
		},
	},

	0x63: {
		Name: "RRA X IND",
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
		Name: "NOP ZP",
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
		Name: "ADC ZP",
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
		Name: "ROR ZP",
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
		Name: "RRA ZP",
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
		Name: "PLA A",
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
		Name: "ADC #",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.ADC(val)
			return true
		},
	},
	0x6A: {
		Name: "ROR A",
		Execute: func(c *cpu, step int) bool {

			c.A = c.ROR(c.A)
			return true
		},
	},

	0x6B: {
		Name: "ARR",
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
		Name: "JMP IND",
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
		Name: "ADC ABS",
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
		Name: "ROR ABS",
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
		Name: "RRA ABS",
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

	0x70: {
		Name: "BVS",
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
		Name: "ADC IND Y",
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
	0x75: {
		Name: "ADC ZP,X",
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
		Name: "ROR ZP,X",
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
	0x78: {
		Name: "SEI IMPL",
		Execute: func(c *cpu, step int) bool {
			c.setFlag(Interrupt)
			return true
		},
	},
	0x79: {
		Name: "ADC ABS Y",
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
	0x7D: {
		Name: "ADC ABS X",
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
		Name: "ROR ABS,X",
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
	0x81: {
		Name: "STA IND X",
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
	0x84: {
		Name: "STY ZP",
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
		Name: "STA ZP",
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
		Name: "STX ZP",
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
	0x88: {
		Name: "DEY",
		Execute: func(c *cpu, step int) bool {
			c.Y--
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0x8A: {
		Name: "TXA",
		Execute: func(c *cpu, step int) bool {
			c.A = c.X
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x8C: {
		Name: "STY",
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
		Name: "STA",
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
		Name: "STX ABS",
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
	0x90: {
		Name: "BCC",
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
				return false
			default:
				fmt.Println("something bad 0x90")
				return false
			}
		},
	},
	0x91: {
		Name: "STA IND Y",
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
	0x94: {
		Name: "STY ZP X",
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
		Name: "STA ZP X",
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
		Execute: func(c *cpu, step int) bool {
			c.A = c.Y
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0x99: {
		Name: "STA ABS Y ",
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
		Name: "TXS",
		Execute: func(c *cpu, step int) bool {
			c.S = c.X
			return true
		},
	},
	0x9D: {
		Name: "STA ABS X",
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
	0xA0: {
		Name: "LDY IMM",
		Execute: func(c *cpu, step int) bool {
			c.Y = c.fetchone()
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xA1: {
		Name: "LDA IND X",
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
		Name: "LDX IMM",
		Execute: func(c *cpu, step int) bool {
			c.X = c.fetchone()
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xA4: {
		Name: "LDY ZP",
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
		Name: "LDA ZP",
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
		Name: "LDX ZP",
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
	0xA8: {
		Name: "TAY",
		Execute: func(c *cpu, step int) bool {
			c.Y = c.A
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xA9: {
		Name: "LDA IMM",
		Execute: func(c *cpu, step int) bool {
			c.A = c.fetchone()
			c.SetFlagNZ(c.A)
			return true
		},
	},
	0xAA: {
		Name: "TAX",
		Execute: func(c *cpu, step int) bool {
			c.X = c.A

			c.SetFlagNZ(c.X)

			return true
		},
	},
	0xAC: {
		Name: "LDY ABS",
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
		Name: "LDA ABS",
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
		Name: "LDY ABS",
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
	0xB0: {
		Name: "BCS",
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
		Name: "LDA IND Y",
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
	0xB4: {
		Name: "LDY ZPG X",
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
		Name: "LDA ZPG X",
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
		Name: "LDX ZPG Y",
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
	0xB8: {
		Name: "CLV",
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(oVerflow)
			return true
		},
	},
	0xB9: {
		Name: "LDA ABS Y",
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
		Name: "TSX",
		Execute: func(c *cpu, step int) bool {
			c.X = c.S
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xBC: {
		Name: "LDY ABS X",
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
		Name: "LDA ABS X",
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
		Name: "LDX ABS Y",
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
	0xC0: {
		Name: "CMY",
		Execute: func(c *cpu, step int) bool {
			val := c.fetchone()
			c.COMPARE(c.Y, val)
			return true
		},
	},
	0xC1: {
		Name: "CMP IND X",
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
	0xC4: {
		Name: "CPY ZP",
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
		Name: "CPM ZP",
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
		Name: "DEC ZP",
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
	0xC8: {
		Name: "INY",
		Execute: func(c *cpu, step int) bool {
			c.Y++
			c.SetFlagNZ(c.Y)
			return true
		},
	},
	0xC9: {
		Name: "CMP #",
		Execute: func(c *cpu, step int) bool {
			c.COMPARE(c.A, c.fetchone())
			return true
		},
	},
	0xCA: {
		Name: "DEX",
		Execute: func(c *cpu, step int) bool {
			c.X--
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xCC: {
		Name: "CPY ABS",
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
		Name: "CPM ABS",
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
		Name: "DEC ABS",
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
	0xD0: {
		Name: "BNE",
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
		Name: "CMP IND Y ",
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
	0xD5: {
		Name: "CMP ZP X",
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
		Name: "DEC ZP X",
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
	0xD8: {
		Name: "CLD",
		Execute: func(c *cpu, step int) bool {
			c.clearFlag(Decimal)
			return true
		},
	},
	0xD9: {
		Name: "CMP ABS Y",
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
	0xDD: {
		Name: "CMP ABS X",
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
		Name: "DEC ABS X",
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
	0xE0: {
		Name: "CPX #",
		Execute: func(c *cpu, step int) bool {
			c.COMPARE(c.X, c.fetchone())
			return true
		},
	},
	0xE1: {
		Name: "SBC IND X",
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
	0xE4: {
		Name: "CPX ZP",
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
		Name: "SBC ZP",
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
		Name: "INC ZP",
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
	0xE8: {
		Name: "INX IMPL",
		Execute: func(c *cpu, step int) bool {
			c.X++
			c.SetFlagNZ(c.X)
			return true
		},
	},
	0xE9: {
		Name: "SBC #",
		Execute: func(c *cpu, step int) bool {
			c.SBC(c.fetchone())
			return true
		},
	}, 0xEA: {
		Name: "NOP",
		Execute: func(c *cpu, step int) bool {
			return true
		},
	},

	0xEC: {
		Name: "CPX ABS",
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
		Name: "SBC ABS",
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
		Name: "INC ABS",
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
	0xF0: {
		Name: "BEQ",
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
		Name: "SBC IND Y ",
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
	0xF5: {
		Name: "SBC ZP X",
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
		Name: "INC ZP X",
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
	0xF8: {
		Name: "SED IMPL",
		Execute: func(c *cpu, step int) bool {
			c.setFlag(Decimal)
			return true
		},
	},
	0xF9: {
		Name: "SBC ABS Y",
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
	0xFD: {
		Name: "SBC ABS X",
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
		Name: "INC ABS X",
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
				fmt.Println("something bad at 0xF9")
				return true
			}
		},
	},
}
