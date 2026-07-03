package Core

import (
	"log"
)

type Mapper interface {
	ReadPRG(addr uint16) uint8
	WritePRG(addr uint16, val uint8)

	ReadCHR(addr uint16) uint8
	WriteCHR(addr uint16, val uint8)

	getMirroring() uint8
	TakeSnapshot(MapperScreenShot)
	CreateEmptySnapshot() MapperScreenShot
}

type MapperScreenShot interface {
	LoadSS(m Mapper)
}

func getMapper(header []byte) int {
	high := header[7] & 0xF0
	low := (header[6] & 0xF0) >> 4

	return int(high | low)
}

func (c *Console) assignMapper(id int, prgData []byte, chrData []byte, mirroring uint8) {
	var m Mapper
	switch id {
	case 0:
		m = &Mapper0{
			PRGROM:     prgData,
			CHRROM:     chrData,
			Mirrroring: mirroring,
		}
	case 1:
		m = &Mapper1{
			PRGROM:        prgData,
			CHRROM:        chrData,
			PRGRAM:        make([]uint8, 0x2000),
			control:       0x0F,
			shiftRegister: 0x10,
			isRam:         c.Ppu.mem.CHR_isRam,
		}
		// case 2:
		// 	m = &Mapper2{
		// 		PRGROM:     prgData,
		// 		CHRROM:     chrData,
		// 		Mirroring:  mirroring,
		// 		BankSelect: 0,
		// 	}
		// case 4:
		// 	m = &Mapper4{
		// 		PRGROM:    prgData,
		// 		CHRROM:    chrData,
		// 		Mirroring: mirroring,
		// 	}

	}

	// case 4:
	// 	m = &Mapper4{
	// 		PRGROM:    prgData,
	// 		CHRROM:    chrData,
	// 		Mirroring: mirroring,
	// 	}

	// }

	c.mapper = m
}

//MAPPER 0

type Mapper0 struct {
	PRGROM []uint8
	CHRROM []uint8

	Mirrroring uint8
}

type Mapper0SS struct {
	PRGROM    []uint8
	CHRROM    []uint8
	Mirroring uint8
}

func (m *Mapper0) CreateEmptySnapshot() MapperScreenShot {
	res := Mapper0SS{
		CHRROM: make([]uint8, len(m.CHRROM)),
		PRGROM: make([]uint8, len(m.PRGROM)),
	}
	return &res
}

func (m *Mapper0) TakeSnapshot(s MapperScreenShot) {

	res, ok := s.(*Mapper0SS)
	if !ok {
		log.Fatalf("babababa")
	}

	copy(res.CHRROM, m.CHRROM)
	copy(res.PRGROM, m.PRGROM)

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

func (s Mapper0SS) LoadSS(m Mapper) {
	safe, ok := m.(*Mapper0)
	if !ok {
		log.Fatalf("FUCK FUCK SMTH RLLY WRONG HAPPENED SHUT DOWN")
	}

	copy(s.CHRROM, safe.CHRROM)
	copy(s.PRGROM, safe.PRGROM)
	safe.Mirrroring = s.Mirroring
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

type Mapper2SS struct {
	PRGROM []uint8
	CHRROM []uint8

	BankSelect uint8
	Mirroring  uint8
}

func (m *Mapper2) CreateEmptySnapshot() MapperScreenShot {
	res := Mapper2SS{
		PRGROM: make([]uint8, len(m.PRGROM)),
		CHRROM: make([]uint8, len(m.CHRROM)),
	}

	return &res
}

func (m Mapper2) TakeSnapshot(s MapperScreenShot) {

	res, ok := s.(*Mapper2SS)

	if !ok {
		log.Fatalf("invalid mapper for screenshot")
	}

	res.BankSelect = m.BankSelect
	res.Mirroring = m.Mirroring

	copy(res.CHRROM, m.CHRROM)
	copy(res.PRGROM, m.PRGROM)

}

func (s Mapper2SS) LoadSS(m Mapper) {

	safe, ok := m.(*Mapper2)

	if !ok {
		log.Fatal("doie")
	}

	copy(safe.CHRROM, s.CHRROM)
	copy(safe.PRGROM, s.PRGROM)

	safe.Mirroring = s.Mirroring
	safe.BankSelect = s.BankSelect
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

func (m *Mapper2) getMirroring() uint8 {
	return m.Mirroring
}

// Mapper 1
type Mapper1 struct {
	PRGROM []uint8
	CHRROM []uint8
	PRGRAM []uint8

	shiftRegister uint8

	control  uint8
	PrgBank  uint8
	ChrBank0 uint8
	ChrBank1 uint8

	isRam bool
}

type Mapper1SS struct {
	PRGROM []uint8
	CHRROM []uint8
	PRGRAM []uint8

	shiftRegister uint8

	control  uint8
	PrgBank  uint8
	ChrBank0 uint8
	ChrBank1 uint8

	isRam bool
}

func (m *Mapper1) CreateEmptySnapshot() MapperScreenShot {
	res := Mapper1SS{
		PRGROM: make([]uint8, len(m.PRGROM)),
		CHRROM: make([]uint8, len(m.CHRROM)),
		PRGRAM: make([]uint8, 0x2000),

		isRam: m.isRam,
	}

	return res
}

func (m *Mapper1) TakeSnapshot(s MapperScreenShot) {
	snap, ok := s.(*Mapper1SS)
	if !ok {
		log.Fatal("wowowoow")
	}

	copy(snap.CHRROM, m.CHRROM)
	copy(snap.PRGROM, m.PRGROM)
	copy(snap.PRGRAM, m.PRGRAM)

	snap.shiftRegister = m.shiftRegister

	snap.control = m.control
	snap.PrgBank = m.PrgBank
	snap.ChrBank0 = m.ChrBank0
	snap.ChrBank1 = m.ChrBank1

}

func (s Mapper1SS) LoadSS(m Mapper) {
	safe, ok := m.(*Mapper1)
	if !ok {
		log.Fatal("im dead bro")
	}

	copy(safe.CHRROM, s.CHRROM)
	copy(safe.PRGROM, s.PRGROM)
	copy(safe.PRGRAM, s.PRGRAM)

	safe.shiftRegister = s.shiftRegister

	safe.control = s.control
	safe.PrgBank = s.PrgBank
	safe.ChrBank0 = s.ChrBank0
	safe.ChrBank1 = s.ChrBank1

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
