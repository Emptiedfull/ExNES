package Core

type TempPPU struct {
	nameTableByte uint8
	attrTableByte uint8
	lowTileByte   uint8
	highTileByte  uint8
	tileData      uint64

	sprintCount      int
	spritePatterns   [8]uint32
	spritePos        [8]uint8
	spritePriorities [8]uint8
	spriteIdx        [8]uint8
}

type PPUInternal struct {
	v uint16 // Vram addr
	t uint16 // Temp Vram
	x uint8  // 3 bit
	w bool   // latch (0-first write,1-second write)
}

func (p *ppu) cycleTick() {
	p.Dot++
	if p.Dot > 340 {
		p.Dot = 0
		p.Scanline++

		if p.Scanline > 261 {
			p.Scanline = 0
			p.Frame++

			p.ScreenChanged = true

			// if m, ok := p.console.mapper.(*MMC3); ok {
			// 	fmt.Printf("frame: %v, edges %v \n", p.Frame, m.edges)
			// 	m.edges = 0
			// }
		}
	}

}

func (p *ppu) step() {

	rendering := p.Mem.register.ShowBG || p.Mem.register.ShowSprites
	if p.Scanline == 0 && p.Dot == 0 && p.Frame%2 == 1 && rendering {
		p.Dot = 1
	}
	prefetch := p.Scanline == 261
	visible := p.Scanline < 240 && p.Scanline >= 0

	renderLine := prefetch || visible

	visibleCycle := p.Dot >= 1 && p.Dot <= 256
	prefetchCycle := p.Dot >= 321 && p.Dot <= 336
	fetchCycle := visibleCycle || prefetchCycle

	if rendering {
		if visible && visibleCycle {
			p.renderPixel()
		}

		if renderLine && fetchCycle {
			p.Mem.temp.tileData <<= 4
			switch p.Dot % 8 {
			case 0:
				var data uint32
				for range 8 {
					a := p.Mem.temp.attrTableByte
					p1 := (p.Mem.temp.lowTileByte & 0x80) >> 7
					p2 := (p.Mem.temp.highTileByte & 0x80) >> 6
					p.Mem.temp.lowTileByte <<= 1
					p.Mem.temp.highTileByte <<= 1
					data <<= 4
					data |= uint32(a | p1 | p2)
				}
				p.Mem.temp.tileData |= uint64(data)
			case 1:
				addr := 0x2000 | (p.Mem.internal.v & 0x0FFF)

				p.Mem.temp.nameTableByte = p.read(addr)
			case 3:
				v := p.Mem.internal.v
				addr := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
				shift := ((v >> 4) & 4) | (v & 2)
				p.Mem.temp.attrTableByte = ((p.read(addr) >> shift) & 3) << 2
			case 5:
				fineY := (p.Mem.internal.v >> 12) & 7
				table := getTableAddr(p.Mem.register.BgPattern)
				tile := p.Mem.temp.nameTableByte
				addr := 0x1000*uint16(table) + uint16(tile)*16 + fineY

				p.Mem.temp.lowTileByte = p.read(addr)
				p.clockIRQ(addr)

			case 7:

				fineY := (p.Mem.internal.v >> 12) & 7
				table := getTableAddr(p.Mem.register.BgPattern)
				tile := p.Mem.temp.nameTableByte
				addr := 0x1000*uint16(table) + uint16(tile)*16 + fineY
				p.Mem.temp.highTileByte = p.read(addr + 8)
				p.clockIRQ(addr + 8)

			}
		}

		if prefetch && p.Dot >= 280 && p.Dot <= 304 {
			//copying y scroll

			p.Mem.internal.v = (p.Mem.internal.v & 0x841F) | (p.Mem.internal.t & 0x7BE0)
		}

		if renderLine {
			if fetchCycle && p.Dot%8 == 0 {
				// x increment
				if p.Mem.internal.v&0x001F == 31 {
					p.Mem.internal.v &= 0xFFE0
					p.Mem.internal.v ^= 0x400
				} else {
					p.Mem.internal.v++
				}
			}
			if p.Dot == 257 {
				p.Mem.internal.v = (p.Mem.internal.v & 0xFBE0) | (p.Mem.internal.t & 0x041F)
			}

			if p.Dot == 256 {
				// y increment
				if p.Mem.internal.v&0x7000 != 0x7000 {
					p.Mem.internal.v += 0x1000
				} else {
					p.Mem.internal.v &= 0x8FFF

					y := (p.Mem.internal.v & 0x03E0) >> 5

					switch y {
					case 29:
						y = 0
						p.Mem.internal.v ^= 0x0800

					case 31:
						y = 0
					default:
						y++
					}

					p.Mem.internal.v = (p.Mem.internal.v & 0xFC1F) | (y << 5)
				}

			}
		}
	}

	if rendering {
		if p.Dot == 257 {
			if visible {
				p.evalSprites()

			} else {
				p.Mem.temp.sprintCount = 0
			}
		}
	}

	if p.Scanline == 241 && p.Dot == 1 {
		p.Mem.Vblank_flag = true
		if p.Mem.register.NmiEnable {
			p.console.Cpu.NmiLine = true

		}
	}

	if prefetch && p.Dot == 1 {
		p.Mem.Vblank_flag = false
		p.Mem.register.Sprite0Hit = false
		p.Mem.register.sprietOverflow = false
	}

	p.cycleTick()
}

func (p *ppu) evalSprites() {
	var h int
	if !p.Mem.register.SpriteSize {
		h = 8
	} else {
		h = 16
	}

	count := 0

	for i := range 64 {
		y := p.Mem.oamData[i*4+0]
		a := p.Mem.oamData[i*4+2]
		x := p.Mem.oamData[i*4+3]

		row := p.Scanline - int(y)

		if row < 0 || row >= h {
			continue
		}

		if count < 8 {
			p.Mem.temp.spritePatterns[count] = p.fetchSpritePattern(i, row)
			p.Mem.temp.spritePos[count] = x
			p.Mem.temp.spritePriorities[count] = (a >> 5) & 1
			p.Mem.temp.spriteIdx[count] = uint8(i)
		}
		count++

	}

	if count > 8 {
		count = 8
		p.Mem.register.sprietOverflow = true
	}

	p.Mem.temp.sprintCount = count

	table := getTableAddr(p.Mem.register.SpritePattern)
	for i := count; i < 8; i++ {
		var dummyAddr uint16
		if !p.Mem.register.SpriteSize {
			dummyAddr = 0x1000*uint16(table) + 0xFF*16
		} else {
			dummyAddr = 0x1000*uint16(0xFF&1) + uint16(0xFE)*16
		}
		p.read(dummyAddr)
		p.clockIRQ(dummyAddr)
	}
}

func (p *ppu) fetchSpritePattern(i, row int) uint32 {
	tile := p.Mem.oamData[i*4+1]
	attr := p.Mem.oamData[i*4+2]

	var addr uint16
	if !p.Mem.register.SpriteSize { //small(8)
		if attr&0x80 == 0x80 {
			row = 7 - row
		}

		table := getTableAddr(p.Mem.register.SpritePattern)
		addr = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else { //big (16)
		if attr&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE

		if row > 7 {
			tile++
			row -= 8
		}

		addr = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}

	a := (attr & 3) << 2

	low := p.read(addr)
	high := p.read(addr + 8)

	p.clockIRQ(addr)

	var data uint32

	for range 8 {
		var p1, p2 byte
		if attr&0x40 == 0x40 {
			p1 = (low & 1) << 0
			p2 = (high & 1) << 1
			low >>= 1
			high >>= 1
		} else {
			p1 = (low & 0x80) >> 7
			p2 = (high & 0x80) >> 6
			low <<= 1
			high <<= 1
		}

		data <<= 4
		data |= uint32(a | p1 | p2)
	}

	return data
}

func getTableAddr(x bool) uint8 {
	if x {
		return 1
	} else {
		return 0
	}
}

func (p *ppu) renderPixel() {
	x := p.Dot - 1
	y := p.Scanline
	bg := p.getBgPixel()
	i, sprite := p.getSpritePixel()

	if x < 8 && !p.Mem.register.ShowLeftBG {
		bg = 0
	}

	if x < 8 && !p.Mem.register.ShowLeftSprite {
		sprite = 0
	}

	b := bg%4 != 0
	s := sprite%4 != 0

	var color uint8
	if !b && !s {
		color = 0
	} else if !b && s {
		color = sprite | 0x10
	} else if b && !s {
		color = bg
	} else {
		if p.Mem.temp.spriteIdx[i] == 0 && x < 255 {
			p.Mem.register.Sprite0Hit = true

		}

		if p.Mem.temp.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = bg
		}
	}

	idx := p.readPallete(uint16(color) % 64)
	col := p.console.PalleteEngine.getColPacked(p.Mem.register.emphasisIndex, idx)
	p.pushRGB(col, x, y)

	// p.pushRGB(p.console.palette[p.Mem.register.emphasisIndex][in], x, y)

}

func (p *ppu) pushRGB(col uint32, x, y int) {

	p.NewBuffer[(y<<8)+x] = col
}

func (p *ppu) readPallete(addr uint16) uint8 {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}

	return p.Mem.Pallete[addr]
}

func (p *ppu) getBgPixel() uint8 {
	if p.Mem.register.ShowBG {
		data := p.fetchTileData() >> ((7 - p.Mem.internal.x) * 4)
		return uint8(data & 0x0F)
	} else {
		return 0
	}
}

func (p *ppu) getSpritePixel() (uint8, uint8) {
	if !p.Mem.register.ShowSprites {
		return 0, 0
	}

	for i := range p.Mem.temp.sprintCount {
		offset := p.Dot - 1 - int(p.Mem.temp.spritePos[i])
		if offset < 0 || offset > 7 {
			continue
		}

		offset = 7 - offset

		color := uint8(p.Mem.temp.spritePatterns[i] >> uint8(offset*4) & 0x0F)
		if color%4 == 0 {
			continue
		}

		return uint8(i), color
	}
	return 0, 0
}

func (p *ppu) fetchTileData() uint32 {
	return uint32(p.Mem.temp.tileData >> 32)
}

func (p *ppu) clockIRQ(addr uint16) {
	if c, ok := p.console.mapper.(IrqClocker); ok {

		c.clockIrqCounter(addr)
	}
}
