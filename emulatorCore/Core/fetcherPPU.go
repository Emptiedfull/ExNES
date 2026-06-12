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

func (p *ppu) cycleTick() {
	p.Dot++
	if p.Dot > 340 {
		p.Dot = 0
		p.Scanline++

		if p.Scanline > 261 {
			p.Scanline = 0
			p.Frame++

			p.screenChanged = true
			copy(p.frontBuffer, p.backBuffer)
		}
	}
}

func (p *ppu) step() {

	rendering := p.mem.register.ShowBG || p.mem.register.ShowSprites
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
			p.mem.temp.tileData <<= 4
			switch p.Dot % 8 {
			case 0:
				var data uint32
				for range 8 {
					a := p.mem.temp.attrTableByte
					p1 := (p.mem.temp.lowTileByte & 0x80) >> 7
					p2 := (p.mem.temp.highTileByte & 0x80) >> 6
					p.mem.temp.lowTileByte <<= 1
					p.mem.temp.highTileByte <<= 1
					data <<= 4
					data |= uint32(a | p1 | p2)
				}
				p.mem.temp.tileData |= uint64(data)
			case 1:
				addr := 0x2000 | (p.mem.internal.v & 0x0FFF)
				p.mem.temp.nameTableByte = p.read(addr)
			case 3:
				v := p.mem.internal.v
				addr := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
				shift := ((v >> 4) & 4) | (v & 2)
				p.mem.temp.attrTableByte = ((p.read(addr) >> shift) & 3) << 2
			case 5:
				fineY := (p.mem.internal.v >> 12) & 7
				table := getTableAddr(p.mem.register.BgPattern)
				tile := p.mem.temp.nameTableByte
				addr := 0x1000*uint16(table) + uint16(tile)*16 + fineY
				p.mem.temp.lowTileByte = p.read(addr)
			case 7:

				fineY := (p.mem.internal.v >> 12) & 7
				table := getTableAddr(p.mem.register.BgPattern)
				tile := p.mem.temp.nameTableByte
				addr := 0x1000*uint16(table) + uint16(tile)*16 + fineY
				p.mem.temp.highTileByte = p.read(addr + 8)

			}
		}

		if prefetch && p.Dot >= 280 && p.Dot <= 304 {
			//copying y scroll
			p.mem.internal.v = (p.mem.internal.v & 0x841F) | (p.mem.internal.t & 0x7BE0)
		}

		if renderLine {
			if fetchCycle && p.Dot%8 == 0 {
				// x increment
				if p.mem.internal.v&0x001F == 31 {
					p.mem.internal.v &= 0xFFE0
					p.mem.internal.v ^= 0x400
				} else {
					p.mem.internal.v++
				}
			}
			if p.Dot == 257 {
				p.mem.internal.v = (p.mem.internal.v & 0xFBE0) | (p.mem.internal.t & 0x041F)
			}

			if p.Dot == 256 {
				// y increment
				if p.mem.internal.v&0x7000 != 0x7000 {
					p.mem.internal.v += 0x1000
				} else {
					p.mem.internal.v &= 0x8FFF

					y := (p.mem.internal.v & 0x03E0) >> 5

					switch y {
					case 29:
						y = 0
						p.mem.internal.v ^= 0x0800

					case 31:
						y = 0
					default:
						y++
					}

					p.mem.internal.v = (p.mem.internal.v & 0xFC1F) | (y << 5)
				}

			}
		}
	}

	if rendering {
		if p.Dot == 257 {
			if visible {
				p.evalSprites()

			} else {
				p.mem.temp.sprintCount = 0
			}
		}
	}

	if p.Scanline == 241 && p.Dot == 1 {
		p.mem.Vblank_flag = true
		if p.mem.register.NmiEnable {
			p.console.Cpu.nmiPending = true
		}
	}

	if prefetch && p.Dot == 1 {
		p.mem.Vblank_flag = false
		p.mem.register.Sprite0Hit = false
		p.mem.register.sprietOverflow = false
	}

	p.cycleTick()
}

func (p *ppu) evalSprites() {
	var h int
	if !p.mem.register.SpriteSize {
		h = 8
	} else {
		h = 16
	}

	count := 0

	for i := range 64 {
		y := p.mem.oamData[i*4+0]
		a := p.mem.oamData[i*4+2]
		x := p.mem.oamData[i*4+3]

		row := p.Scanline - int(y)

		if row < 0 || row >= h {
			continue
		}

		if count < 8 {
			p.mem.temp.spritePatterns[count] = p.fetchSpritePattern(i, row)
			p.mem.temp.spritePos[count] = x
			p.mem.temp.spritePriorities[count] = (a >> 5) & 1
			p.mem.temp.spriteIdx[count] = uint8(i)
		}
		count++
	}

	if count > 8 {
		count = 8
		p.mem.register.sprietOverflow = true
	}

	p.mem.temp.sprintCount = count
}

func (p *ppu) fetchSpritePattern(i, row int) uint32 {
	tile := p.mem.oamData[i*4+1]
	attr := p.mem.oamData[i*4+2]

	var addr uint16
	if !p.mem.register.SpriteSize { //small(8)
		if attr&0x80 == 0x80 {
			row = 7 - row
		}

		table := getTableAddr(p.mem.register.SpritePattern)
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

	if x < 8 && !p.mem.register.ShowLeftBG {
		bg = 0
	}

	if x < 8 && !p.mem.register.ShowLeftSprite {
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
		if p.mem.temp.spriteIdx[i] == 0 && x < 255 {
			p.mem.register.Sprite0Hit = true

		}

		if p.mem.temp.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = bg
		}
	}

	c := NesPaletteLUT[p.readPallete(uint16(color)%64)]
	p.pushRGB(c, x, y)

}

func (p *ppu) pushRGB(c RGB, x, y int) {

	idx := ((y << 8) + x) << 2

	p.backBuffer[idx] = c.R
	p.backBuffer[idx+1] = c.G
	p.backBuffer[idx+2] = c.B
	p.backBuffer[idx+3] = 255
}

func (p *ppu) readPallete(addr uint16) uint8 {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}

	return p.mem.Pallete[addr]
}

func (p *ppu) getBgPixel() uint8 {
	if p.mem.register.ShowBG {
		data := p.fetchTileData() >> ((7 - p.mem.internal.x) * 4)
		return uint8(data & 0x0F)
	} else {
		return 0
	}
}

func (p *ppu) getSpritePixel() (uint8, uint8) {
	if !p.mem.register.ShowSprites {
		return 0, 0
	}

	for i := range p.mem.temp.sprintCount {
		offset := p.Dot - 1 - int(p.mem.temp.spritePos[i])
		if offset < 0 || offset > 7 {
			continue
		}

		offset = 7 - offset

		color := uint8(p.mem.temp.spritePatterns[i] >> uint8(offset*4) & 0x0F)
		if color%4 == 0 {
			continue
		}

		return uint8(i), color
	}
	return 0, 0
}

func (p *ppu) fetchTileData() uint32 {
	return uint32(p.mem.temp.tileData >> 32)
}
