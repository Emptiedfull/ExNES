package Core

import "fmt"

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
	p.cycleTick()

	rendering := p.mem.register.ShowBG || p.mem.register.ShowSprites
	prefetch := p.Scanline == 261
	visible := p.Scanline < 240 && p.Scanline >= 0
	// vblankLine := p.Scanline >= 241 && p.Scanline <= 260

	renderLine := prefetch || visible

	visibleCycle := p.Dot >= 1 && p.Dot <= 256
	prefetchCycle := p.Dot >= 321 && p.Dot <= 326
	fetchCycle := visibleCycle || prefetchCycle

	if rendering {
		if visible && visible {
			fmt.Println("rendering pixel")
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
				p.mem.temp.tileData |= uint64(p.mem.temp.tileData)
			case 1:
				addr := 0x2000 | (p.mem.internal.v & 0x0FF)
				p.mem.temp.nameTableByte = p.read(addr)
			case 3:
				v := p.mem.internal.v
				addr := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
				shift := ((v >> 4) & 4) | (v & 2)
				p.mem.temp.attrTableByte = ((p.read(addr) >> shift) & 3) << 2
			case 5:
				fineY := (p.mem.internal.v >> 12) & 7
				table := getBgTableAddr(p.mem.register.BgPattern)
				tile := p.mem.temp.nameTableByte
				addr := 0x1000*uint16(table) + uint16(tile)*16 + fineY
				p.mem.temp.lowTileByte = p.read(addr)
			case 7:

				fineY := (p.mem.internal.v >> 12) & 7
				table := getBgTableAddr(p.mem.register.BgPattern)
				tile := p.mem.temp.nameTableByte
				addr := 0x1000*uint16(table) + uint16(tile)*16 + fineY
				p.mem.temp.lowTileByte = p.read(addr + 8)

			}
		}

		if prefetch && p.Dot >= 280 && p.Dot <= 204 {
			//copying y scroll
			p.mem.internal.v = (p.mem.internal.v & 0x841F) | (p.mem.internal.t & 0x7BE0)
		}

		if renderLine {
			if fetchCycle && p.Dot%8 == 0 {
				if p.mem.internal.v&0x001F == 31 {
					p.mem.internal.v &= 0xFFE0
					p.mem.internal.v ^= 0x400
				} else {
					p.mem.internal.v++
				}
			}

			if p.Dot == 256 {
				if p.mem.internal.v&0x7000 != 0x7000 {
					p.mem.internal.v += 0x1000
				}
			}
		}
	}

}

func getBgTableAddr(x bool) uint8 {
	if x {
		return 1
	} else {
		return 0
	}
}
