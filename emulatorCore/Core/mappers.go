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
	switch id {
	case 0:
		m := Mapper0{
			PRGROM: prgData,
			CHRROM: chrData,
		}

		c.Cpu.mem.mapper = &m
		c.Ppu.mem.mapper = &m
	}
}

//MAPPER 0

type Mapper0 struct {
	PRGROM []uint8
	CHRROM []uint8
}

func (m *Mapper0) ReadPRG(addr uint16) uint8 {
	return m.PRGROM[addr&0x7FFF]
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

//MAPPER 1
