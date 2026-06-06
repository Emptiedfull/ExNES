package Core

//for any readers, please stop here I barely understand what ive done

type ppu struct {
	mem     ppu_mem
	console *console

	Dot      int
	Scanline int
	Frame    int

	pallete Pallete

	backBuffer  []uint8
	frontBuffer []uint8

	DebugBuffer []uint8

	DrawFlg bool

	verticalMirroring bool
}

type ppu_mem struct {
	chrROM  []uint8
	Vram    [2048]uint8
	Pallete [32]uint8
	oamData [256]uint8
	Addr    uint16

	register registersFlags

	nmiDelay int
	nmiPrev  bool
	nmiOut   bool
	nmiOcc   bool
	Vblank   bool
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

	// internal registers
	VramAddr     uint16 // $2006 v
	TramAddr     uint16 // t
	AddressLatch bool   // w  0- high byte, 1-low byte
	XAddr        uint8  // x [3 bits]
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

func (p *ppu) GetScreenBuffer() []uint8 {
	return p.backBuffer
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
		p.mem.register.AddressLatch = false

		var result uint8
		result |= (openBusVal & 0x1F)

		result = AssignBit(result, 6, p.mem.register.Sprite0Hit)
		result = AssignBit(result, 5, p.mem.register.sprietOverflow)

		if p.mem.Vblank {
			result = AssignBit(result, 7, true)

		}

		p.mem.Vblank = false

		return result
	case 4:
		data := p.mem.oamData[p.mem.register.OAMADDR]

		if (p.mem.register.OAMADDR & 0x03) == 0x02 {
			data = data & 0xE3
		}

		return data
	case 7:
		val := p.read(p.mem.register.VramAddr)

		if p.mem.register.VramAddr%0x4000 < 0x3F00 {
			buffered := p.mem.register.bufferedData
			p.mem.register.bufferedData = val
			val = buffered
		} else {
			p.mem.register.bufferedData = p.read(p.mem.register.VramAddr - 0x1000)

			val = (val & 0x3F) | (openBusVal & 0xC0)
		}

		if p.mem.register.AddrIncrement {
			p.mem.register.VramAddr += 32
		} else {
			p.mem.register.VramAddr += 1
		}
		p.mem.register.VramAddr &= 0x3FFF

		return val
	default:
		return 0
	}
}

func (p *ppu) WriteReg(reg uint16, val uint8) {
	switch reg {
	case 0: //PPUCTRL

		p.mem.register.NmiEnable = getbitBool(val, 7)
		p.mem.nmiOut = p.mem.register.NmiEnable
		p.mem.register.SpriteSize = getbitBool(val, 5)
		p.mem.register.BgPattern = getbitBool(val, 4)
		p.mem.register.SpritePattern = getbitBool(val, 3)

		p.mem.register.AddrIncrement = getbitBool(val, 2)

		p.mem.register.BaseNameTable = val & 0x03

		p.mem.register.TramAddr = p.mem.register.TramAddr & 0xF3FF
		p.mem.register.TramAddr = p.mem.register.TramAddr | (uint16(val&0x03) << 10)

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
		p.mem.register.AddressLatch = !p.mem.register.AddressLatch
	case 6:
		if !p.mem.register.AddressLatch {
			p.mem.register.TramAddr = (p.mem.register.TramAddr & 0x00FF) | (uint16(val&0x3F) << 8)
			p.mem.register.AddressLatch = true
		} else {
			p.mem.register.TramAddr = (p.mem.register.TramAddr & 0xFF00) | uint16(val)
			p.mem.register.VramAddr = p.mem.register.TramAddr & 0x3FFF
			p.mem.register.AddressLatch = false

		}
	case 7:
		targetAddrBeforeWrite := p.mem.register.VramAddr
		isNametable0Grid := targetAddrBeforeWrite >= 0x2000 && targetAddrBeforeWrite <= 0x23BF

		p.Write(p.mem.register.VramAddr, val)
		if p.mem.register.AddrIncrement {
			p.mem.register.VramAddr += 32
		} else {
			p.mem.register.VramAddr += 1
		}
		p.mem.register.VramAddr &= 0x3FFF

		if isNametable0Grid {

			linearOffset := targetAddrBeforeWrite - 0x2000
			_ = linearOffset % 32
			_ = linearOffset / 32

			// fmt.Printf("[NAMETABLE 0 GRID WRITE] Tile ID: 0x%02X -> Addr: 0x%04X (Row: %d, Col: %d)\n",
			// 	val, targetAddrBeforeWrite, row, col)
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

	if console.Cpu.totalCycles%2 == 1 {
		cycles += 1
	}

	console.Cpu.Stall = cycles
}

func (p *ppu) Tick() {

	if p.mem.nmiOut {
		p.console.Cpu.nmiPending = true
		p.mem.nmiOut = false

	}

	p.Dot++
	if p.Dot > 340 {
		p.Dot = 0

		if p.Scanline < 240 {
			p.renderScanline()

		}

		p.Scanline++
		if p.Scanline > 261 {
			p.Scanline = 0
			p.Frame++

			copy(p.frontBuffer, p.backBuffer)
			p.DrawFlg = true
		}
	}

	if p.Scanline == 241 && p.Dot == 1 {

		p.mem.Vblank = true
		p.mem.nmiOcc = true

	}

	if p.Scanline == 261 && p.Dot == 1 {
		p.mem.nmiOcc = false
		p.mem.Vblank = false
		p.mem.register.Sprite0Hit = false
		p.mem.register.sprietOverflow = false

	}

}
