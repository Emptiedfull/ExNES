package Core

type Mapper interface {
	ReadPRG(addr uint16) uint8
	WritePRG(addr uint16, val uint8)

	ReadCHR(addr uint16) uint8
	WriteCHR(addr uint16, val uint8)

	extractCHR() []uint8
	loadCHR([]uint8)
}

func (c *console) assignMapper(id int, prgData []byte, chrData []byte) {
	var m Mapper
	switch id {
	case 0:
		m = &Mapper0{
			PRGROM: prgData,
			CHRROM: chrData,
		}
	case 2:
		m = &Mapper2{
			PRGROM:     prgData,
			CHRROM:     chrData,
			BankSelect: 0,
		}
	}

	c.Cpu.mem.mapper = m
	c.Ppu.mem.mapper = m
}

//MAPPER 0

type Mapper0 struct {
	PRGROM []uint8
	CHRROM []uint8
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

// MAPPER 1
type Mapper2 struct {
	PRGROM []uint8
	CHRROM []uint8

	BankSelect uint8
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
