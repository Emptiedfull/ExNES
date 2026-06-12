package Core

//for any readers, please stop here I barely understand what ive done

type ppu struct {
	mem     ppu_mem
	console *console

	Dot      int
	Scanline int
	Frame    int

	backBuffer  []uint8
	frontBuffer []uint8

	screenChanged bool

	DebugBuffer []uint8

	verticalMirroring bool
}

type ppu_mem struct {
	chrROM  []uint8
	Vram    [2048]uint8
	Pallete [32]uint8
	oamData [256]uint8
	Addr    uint16

	register registersFlags
	internal PPUInternal
	temp     TempPPU

	Vblank_flag bool
}

type registersFlags struct {

	// $2000 WRITE PPUCONTROL bit, meaning
	NmiEnable     bool  // 7 0-off, 1-on
	SpriteSize    bool  // 5 0-8x8, 1-8x16
	BgPattern     bool  // 4 0- $0000, 1- $1000
	SpritePattern bool  // 3 0-$0000, 1- $1000
	AddrIncrement bool  // 2 0-1 going across, 1-32 going down
	BaseNameTable uint8 // 0,1 bit

	// $2001 WRITE PPUMASK
	EmpBlue        bool // 7
	EmpGreen       bool // 6
	EmpRed         bool // 5
	ShowSprites    bool // 4
	ShowBG         bool // 3
	ShowLeftSprite bool // 2
	ShowLeftBG     bool // 1
	GreyScale      bool // 0

	// $2002 READ PPUSTATUS bits 0-4 are useless
	VBlank         bool // 7, set at vblank start
	Sprite0Hit     bool // 6
	sprietOverflow bool // 5

	// $2003 OAMDDR READ
	OAMADDR uint8

	// $2004 OAMDATA WRITE/READ
	OAMDATA uint8

	// $2005
	ScrollY   bool  //Tracks first or second write
	ScrollVal uint8 // first-scrollx, second-scrolly

	// $2007
	bufferedData uint8
}

type PPUInternal struct {
	v uint16 // Vram addr
	t uint16 // Temp Vram
	x uint8  // 3 bit
	w bool   // latch (0-first write,1-second write)
}

func (p *ppu) MirrorNameTable(addr uint16) uint16 {

	addr = (addr - 0x2000) % 0x1000

	if p.verticalMirroring {
		return addr % 0x0800
	} else {

		if addr < 0x0400 {
			return addr
		} else if addr < 0x0800 {
			return addr - 0x0400
		} else if addr < 0x0C00 {
			return addr - 0x0400
		} else {
			return addr - 0x0800
		}
	}
}

func (p *ppu) read(addr uint16) uint8 {
	addr &= 0x3FFF

	switch {
	case addr <= 0x1FFF:
		return p.mem.chrROM[addr]
	case addr <= 0x3EFF:

		mirrored := p.MirrorNameTable(addr)
		return p.mem.Vram[mirrored]
	case addr <= 0x3FFF:
		palleteAddr := (addr - 0x3F00) % 32

		if palleteAddr >= 16 && palleteAddr%4 == 0 {
			palleteAddr -= 16
		}

		return p.mem.Pallete[palleteAddr]
	}
	return 0
}

func (p *ppu) GetScreenBuffer() []uint8 {
	return p.backBuffer
}

func (p *ppu) Write(addr uint16, val uint8) {
	addr &= 0x3FFF

	switch {
	case addr <= 0x1FFF:
		p.mem.chrROM[addr] = val
	case addr <= 0x3EFF:
		mirroredaddr := p.MirrorNameTable(addr)
		p.mem.Vram[mirroredaddr] = val
	case addr <= 0x3FFF:
		palleteAddr := (addr - 0x3F00) % 32

		if palleteAddr >= 16 && (palleteAddr&0x03) == 0 {
			palleteAddr -= 16
		}

		p.mem.Pallete[palleteAddr] = val
	}
}

func (p *ppu) ReadReg(reg uint16, openBusVal uint8) uint8 {
	switch reg {
	case 2:

		var result uint8
		result |= (openBusVal & 0x1F)

		result = AssignBit(result, 6, p.mem.register.Sprite0Hit)
		result = AssignBit(result, 5, p.mem.register.sprietOverflow)

		if p.mem.Vblank_flag {
			result = AssignBit(result, 7, true)

		}

		p.mem.Vblank_flag = false
		p.mem.internal.w = false

		return result
	case 4:
		data := p.mem.oamData[p.mem.register.OAMADDR]

		if (p.mem.register.OAMADDR & 0x03) == 0x02 {
			data = data & 0xE3
		}

		return data

	case 7:
		val := p.read(p.mem.internal.v)

		if p.mem.internal.v%0x4000 < 0x3F00 {
			buffered := p.mem.register.bufferedData
			p.mem.register.bufferedData = val
			val = buffered
		} else {
			p.mem.register.bufferedData = p.read(p.mem.internal.v - 0x1000)

			val = (val & 0x3F) | (openBusVal & 0xC0)
		}

		if p.mem.register.AddrIncrement {
			p.mem.internal.v += 32
		} else {
			p.mem.internal.v += 1
		}
		p.mem.internal.v &= 0x3FFF

		return val
	default:
		return 0
	}
}

func (p *ppu) WriteReg(reg uint16, val uint8) {
	switch reg {
	case 0: //PPUCTRL

		p.mem.register.NmiEnable = getbitBool(val, 7)
		p.mem.register.SpriteSize = getbitBool(val, 5)
		p.mem.register.BgPattern = getbitBool(val, 4)
		p.mem.register.SpritePattern = getbitBool(val, 3)

		p.mem.register.AddrIncrement = getbitBool(val, 2)

		p.mem.register.BaseNameTable = val & 0x03

		p.mem.internal.t = (p.mem.internal.t & 0xF3FF) | ((uint16(val) & 0x03) << 10)

	case 1: //PPU MASK
		p.mem.register.EmpBlue = getbitBool(val, 7)
		p.mem.register.EmpGreen = getbitBool(val, 6)
		p.mem.register.EmpRed = getbitBool(val, 5)
		p.mem.register.ShowSprites = getbitBool(val, 4)
		p.mem.register.ShowBG = getbitBool(val, 3)
		p.mem.register.ShowLeftSprite = getbitBool(val, 2)
		p.mem.register.ShowLeftBG = getbitBool(val, 1)
		p.mem.register.GreyScale = getbitBool(val, 0)
	case 3: //PPU OAMADDR
		p.mem.Addr = uint16(val)
	case 4: //PPU OAMDATA
		p.mem.oamData[p.mem.Addr] = val
		p.mem.Addr++
	case 5:
		if p.mem.internal.w {
			p.mem.internal.t = (p.mem.internal.t & 0xFFE0) | (uint16(val) >> 3)
			p.mem.internal.x = val & 0x07

		} else {
			p.mem.internal.t = (p.mem.internal.t & 0x8FFF) | ((uint16(val) & 0x07) << 12)
			p.mem.internal.t = (p.mem.internal.t & 0xFC1F) | ((uint16(val) & 0xF8) << 2)

		}
		p.mem.internal.w = !p.mem.internal.w
	case 6:
		if !p.mem.internal.w {
			p.mem.internal.t = (p.mem.internal.t & 0x00FF) | (uint16(val&0x3F) << 8)

		} else {
			p.mem.internal.t = (p.mem.internal.t & 0xFF00) | uint16(val)
			p.mem.internal.v = p.mem.internal.t

		}
		p.mem.internal.w = !p.mem.internal.w
	case 7:
		p.Write(p.mem.internal.v, val)

		if p.mem.register.AddrIncrement {
			p.mem.internal.v += 32
		} else {
			p.mem.internal.v += 1
		}
	}
}

func (console *console) ExecuteOAMDMA(page uint8) {

	cpuSrc := uint(page) << 8

	for i := uint16(0); i < 256; i++ {
		dat := console.Cpu.mem.Read(uint16(cpuSrc) + i)

		console.Ppu.mem.oamData[console.Ppu.mem.register.OAMADDR] = dat
		console.Ppu.mem.register.OAMADDR++
	}

	cycles := 513

	if console.Cpu.TotalCycles%2 == 1 {
		cycles += 1
	}

	console.Cpu.Stall = cycles
}

func (p *ppu) Tick() {

	p.Dot++
	if p.Dot > 340 {
		p.Dot = 0

		if p.Scanline < 240 {
			p.renderScanline()
			p.screenChanged = true
		}
		p.Scanline++
		if p.Scanline > 261 {
			p.Scanline = 0
			p.Frame++

			copy(p.frontBuffer, p.backBuffer)
		}
	}

	if p.Scanline == 241 && p.Dot == 1 {

		p.mem.Vblank_flag = true
		if p.mem.register.NmiEnable {
			p.console.Cpu.nmiPending = true
		}

	}

	if p.Scanline == 261 && p.Dot == 1 {
		p.mem.Vblank_flag = false
		p.mem.register.Sprite0Hit = false
		p.mem.register.sprietOverflow = false

	}

}
