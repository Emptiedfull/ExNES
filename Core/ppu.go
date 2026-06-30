package Core

//for any readers, please stop here I barely understand what ive done

type ppu struct {
	mem     ppu_mem
	console *Console

	Dot      int
	Scanline int
	Frame    int

	BackBuffer [245760]uint8

	mirroring       int
	mirroringChange bool

	screenChanged bool

	DebugBuffer []uint8
}

var NesPaletteLUT = [64]RGB{
	{0x7C, 0x7C, 0x7C}, {0x00, 0x00, 0xFC}, {0x00, 0x00, 0xBC}, {0x44, 0x28, 0xBC},
	{0x94, 0x00, 0x84}, {0xA8, 0x00, 0x20}, {0xA8, 0x10, 0x00}, {0x88, 0x14, 0x00},
	{0x50, 0x30, 0x00}, {0x00, 0x78, 0x00}, {0x00, 0x68, 0x00}, {0x00, 0x58, 0x00},
	{0x00, 0x40, 0x58}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},

	{0xBC, 0xBC, 0xBC}, {0x00, 0x78, 0xF8}, {0x00, 0x58, 0xF8}, {0x68, 0x44, 0xFC},
	{0xD8, 0x00, 0xCC}, {0xE4, 0x00, 0x58}, {0xF8, 0x38, 0x00}, {0xE4, 0x5C, 0x10},
	{0xAC, 0x7C, 0x00}, {0x00, 0xB8, 0x00}, {0x00, 0xA8, 0x00}, {0x00, 0xA8, 0x44},
	{0x00, 0x88, 0x88}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},

	{0xF8, 0x78, 0xF8}, {0x3C, 0xBC, 0xFC}, {0x68, 0xA0, 0xFC}, {0x98, 0x88, 0xF8},
	{0xF8, 0x78, 0xF8}, {0xFD, 0x54, 0x8C}, {0xF8, 0x78, 0x58}, {0xFC, 0xA0, 0x44},
	{0xFF, 0xBC, 0x3C}, {0xB4, 0xD4, 0x20}, {0x48, 0xCD, 0x48}, {0x58, 0xFA, 0x98},
	{0x58, 0xF8, 0xFC}, {0x78, 0x78, 0x78}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},

	{0xFC, 0xFC, 0xFC}, {0xA4, 0xE4, 0xFC}, {0xB8, 0xB8, 0xF8}, {0xD8, 0xB8, 0xF8},
	{0xF8, 0xB8, 0xF8}, {0xF8, 0xA4, 0xC0}, {0xF0, 0xD0, 0xB0}, {0xFC, 0xE0, 0xA8},
	{0xFD, 0xEA, 0x88}, {0xD8, 0xF8, 0x78}, {0xB8, 0xF8, 0xB8}, {0xB8, 0xF8, 0xD8},
	{0x00, 0xFC, 0xFC}, {0xD8, 0xD8, 0xD8}, {0x00, 0x00, 0x00}, {0x00, 0x00, 0x00},
}

type RGB struct {
	R, G, B uint8
}

type ppu_mem struct {
	mapper    Mapper
	CHR_isRam bool

	Vram    [2048]uint8
	Pallete [32]uint8
	oamData [256]uint8

	register registersFlags
	internal PPUInternal
	temp     TempPPU

	// chrRom_WARNING []uint8 //THIS IS ONLY TO BE USED FOR SNAPSHOT PURPOSES NO USE IN MAIN CPU I BEG YOU

	Vblank_flag bool
}

// I wrote all of these pain stakingly and THEY WERE NEEDED SO FUCK OFF AI LARPERS
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

func (p *ppu) MirrorNameTable(addr uint16) uint16 {
	addr = (addr - 0x2000) % 0x1000
	table := addr / 0x400
	offset := addr % 0x400

	m := p.mirroring

	case0 := offset & -uint16(BoolToUint16(m == 0))
	case1 := (0x400 + offset) & -uint16(BoolToUint16(m == 1))
	case2 := ((table%2)*0x400 + offset) & -uint16(BoolToUint16(m == 2))

	lt := BoolToUint16(table < 2)
	case3a := offset & -uint16(BoolToUint16(m == 3)&lt)
	case3b := (0x400 + offset) & -uint16(BoolToUint16(m == 3)&(1-lt))

	def := addr & -uint16(BoolToUint16(m > 3))
	return case0 | case1 | case2 | case3a | case3b | def
}

func (p *ppu) read(addr uint16) uint8 {
	addr &= 0x3FFF

	switch {
	case addr <= 0x1FFF:
		return p.mem.mapper.ReadCHR(addr)
	case addr <= 0x3EFF:

		mirrored := p.MirrorNameTable(addr)
		return p.mem.Vram[mirrored]
	case addr <= 0x3FFF:
		palleteAddr := (addr - 0x3F00) % 32

		isSpriteMirror := BoolToUint16(palleteAddr >= 16) & BoolToUint16(palleteAddr%4 == 0)
		palleteAddr -= 16 * isSpriteMirror

		return p.mem.Pallete[palleteAddr]
	}
	return 0
}

func (p *ppu) Write(addr uint16, val uint8) {

	addr &= 0x3FFF

	switch {
	case addr <= 0x1FFF:
		p.mem.mapper.WriteCHR(addr, val)
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
		p.mem.register.OAMADDR = val
	case 4: //PPU OAMDATA
		p.mem.oamData[p.mem.register.OAMADDR] = val
		p.mem.register.OAMADDR++
	case 5:
		if !p.mem.internal.w {
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

		p.mem.internal.v &= 0x3FFF
	}
}

func (console *Console) ExecuteOAMDMA(page uint8) {

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
