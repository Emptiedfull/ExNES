package Core

type Mapper interface {
	ReadPRG(addr uint16) uint8
	WritePRG(addr uint16, val uint8)

	ReadCHR(addr uint16) uint8
	WriteCHR(addr uint16, val uint8)

	extractCHR() []uint8
	loadCHR([]uint8)

	getMirroring() uint8
}

func (c *console) assignMapper(id int, prgData []byte, chrData []byte, mirroring uint8) {
	var m Mapper
	switch id {
	case 0:
		m = &Mapper0{
			PRGROM:     prgData,
			CHRROM:     chrData,
			Mirrroring: mirroring,
		}

	case 2:
		m = &Mapper2{
			PRGROM:     prgData,
			CHRROM:     chrData,
			Mirroring:  mirroring,
			BankSelect: 0,
		}
	case 1:
		m = &Mapper1{
			PRGROM:        prgData,
			CHRROM:        chrData,
			control:       0x0F,
			shiftRegister: 0x10,
			isRam:         c.Ppu.mem.CHR_isRam,
		}
	}

	c.Cpu.mem.mapper = m
	c.Ppu.mem.mapper = m
}

//MAPPER 0

type Mapper0 struct {
	PRGROM []uint8
	CHRROM []uint8

	Mirrroring uint8
}

func (m *Mapper0) ReadPRG(addr uint16) uint8 {
	mask := uint16(len(m.PRGROM) - 1)
	return m.PRGROM[(addr-0x8000)&mask]
}

func (m *Mapper0) WritePRG(addr uint16, val uint8) {

}

func (m *Mapper0) ReadCHR(addr uint16) uint8 {
	return m.CHRROM[addr]
}

func (m *Mapper0) WriteCHR(addr uint16, val uint8) {

}

func (m *Mapper0) extractCHR() []uint8 {
	return m.CHRROM
}

func (m *Mapper0) loadCHR(data []uint8) {
	m.CHRROM = data
}

func (m *Mapper0) getMirroring() uint8 {
	return m.Mirrroring
}

// MAPPER 2
type Mapper2 struct {
	PRGROM []uint8
	CHRROM []uint8

	BankSelect uint8
	Mirroring  uint8
}

func (m *Mapper2) ReadPRG(addr uint16) uint8 {
	if addr < 0xC000 {
		bankOffset := uint32(m.BankSelect) * 0x4000
		return m.PRGROM[bankOffset+uint32(addr-0x8000)]
	} else {
		lastBank := uint32((len(m.PRGROM)/0x4000)-1) * 0x4000
		return m.PRGROM[lastBank+uint32(addr-0xC000)]
	}
}

func (m *Mapper2) WritePRG(addr uint16, val uint8) {
	if addr >= 0x8000 {
		bankCount := uint8(len(m.PRGROM) / 0x4000)
		m.BankSelect = val % bankCount
	}
}

func (m *Mapper2) ReadCHR(addr uint16) uint8 {
	return m.CHRROM[addr]
}

func (m *Mapper2) WriteCHR(addr uint16, val uint8) {
	m.CHRROM[addr] = val
}

func (m *Mapper2) extractCHR() []uint8 {
	return m.CHRROM
}

func (m *Mapper2) loadCHR(data []uint8) {
	m.CHRROM = data
}

func (m *Mapper2) getMirroring() uint8 {
	return m.Mirroring
}

// Mapper 1
type Mapper1 struct {
	PRGROM []uint8
	CHRROM []uint8
	PRGRAM [0x2000]uint8

	shiftRegister uint8

	control  uint8
	PrgBank  uint8
	ChrBank0 uint8
	ChrBank1 uint8

	isRam bool
}

func (m *Mapper1) WritePRG(addr uint16, val uint8) {

	if addr < 0x8000 {

		if addr >= 0x6000 {

			if m.PrgBank&0x10 == 0 {
				m.PRGRAM[addr-0x6000] = val
			}
		}
		return
	}

	if val&0x80 != 0 {
		//reset
		m.shiftRegister = 0x10
		m.control |= 0x0C
	} else {

		complete := (m.shiftRegister & 1) == 1

		m.shiftRegister >>= 1
		m.shiftRegister |= (val & 1) << 4

		if complete {
			regId := (addr >> 13) & 0x03
			data := m.shiftRegister & 0x1F

			switch regId {
			case 0:
				m.control = data
			case 1:
				m.ChrBank0 = data
			case 2:
				m.ChrBank1 = data
			case 3:
				m.PrgBank = data
			}

			m.shiftRegister = 0x10

		}
	}
}

func (m *Mapper1) ReadPRG(addr uint16) uint8 {

	if addr < 0x8000 {
		if addr >= 0x6000 {

			if m.PrgBank&0x10 == 0 {
				return m.PRGRAM[addr-0x6000]
			}
		}
		return 0
	}
	mode := (m.control >> 2) & 0x03

	switch mode {
	case 0, 1:
		numBanks := uint32(len(m.PRGROM) / 0x8000)
		bank := uint32((m.PrgBank&0x0E)>>1) % numBanks
		return m.PRGROM[bank*0x8000+uint32(addr-0x8000)]
	case 2:
		if addr < 0xC000 {
			return m.PRGROM[addr-0x8000]
		}

		numBanks := uint8(len(m.PRGROM) / 0x4000)
		bank := (m.PrgBank & 0x0F) % numBanks
		return m.PRGROM[uint32(bank)*0x4000+uint32(addr-0xC000)]
	case 3:
		if addr < 0xC000 {
			numBanks := uint8(len(m.PRGROM) / 0x4000)
			bank := (m.PrgBank & 0x0F) % numBanks
			return m.PRGROM[uint32(bank)*0x4000+uint32(addr-0x8000)]
		}

		lastbank := uint32((len(m.PRGROM) / 0x4000) - 1)
		return m.PRGROM[lastbank*0x4000+uint32(addr-0xC000)]
	}

	return 0
}

func (m *Mapper1) ReadCHR(addr uint16) uint8 {
	if m.isRam {

		return m.CHRROM[addr]
	}

	if (m.control & 0x10) != 0 {

		if addr < 0x1000 {
			bank := uint32(m.ChrBank0)
			return m.CHRROM[(bank*0x1000)+uint32(addr)]
		} else {
			bank := uint32(m.ChrBank1)
			return m.CHRROM[(bank*0x1000)+uint32(addr-0x1000)]
		}
	} else {
		bank := uint32(m.ChrBank0&0x1E) >> 1
		return m.CHRROM[(bank*0x2000)+uint32(addr)]
	}
}

func (m *Mapper1) WriteCHR(addr uint16, val uint8) {
	if m.isRam {
		m.CHRROM[addr] = val
		// if val != 0 {
		// 	fmt.Printf("First CHR RAM write: addr=%04X val=%02X\n", addr, val)

		// }
		return
	}

	if (m.control & 0x10) != 0 {
		if addr < 0x1000 {
			bank := uint32(m.ChrBank0)
			m.CHRROM[(bank*0x1000)+uint32(addr)] = val
		} else {
			bank := uint32(m.ChrBank1)
			m.CHRROM[(bank*0x1000)+uint32(addr-0x1000)] = val
		}
	} else {

		bank := uint32(m.ChrBank0&0x1E) >> 1
		m.CHRROM[(bank*0x2000)+uint32(addr)] = val
	}
}

func (m *Mapper1) extractCHR() []uint8 {
	return m.CHRROM
}

func (m *Mapper1) loadCHR(data []uint8) {
	m.CHRROM = data
}

func (m *Mapper1) getMirroring() uint8 {
	mode := m.control & 0x03

	return mode
}
