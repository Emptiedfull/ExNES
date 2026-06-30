package Core

import (
	"fmt"
	"log"
)

type Mapper interface {
	ReadPRG(addr uint16) uint8
	WritePRG(addr uint16, val uint8)

	ReadCHR(addr uint16) uint8
	WriteCHR(addr uint16, val uint8)

	getMirroring() uint8
	TakeSnapshot() MapperScreenShot
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
			control:       0x0F,
			shiftRegister: 0x10,
			isRam:         c.Ppu.mem.CHR_isRam,
		}
	case 2:
		m = &Mapper2{
			PRGROM:     prgData,
			CHRROM:     chrData,
			Mirroring:  mirroring,
			BankSelect: 0,
		}
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

	c.Cpu.mem.mapper = m
	c.Ppu.mem.mapper = m
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

func (m *Mapper0) TakeSnapshot() MapperScreenShot {
	res := Mapper0SS{
		PRGROM:    make([]uint8, len(m.PRGROM)),
		CHRROM:    make([]uint8, len(m.CHRROM)),
		Mirroring: m.Mirrroring,
	}

	copy(res.CHRROM, m.CHRROM)
	copy(res.PRGROM, m.PRGROM)

	return res
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
	PRGRAM []uint8
	CHRROM []uint8

	BankSelect uint8
	Mirroring  uint8
}

func (m Mapper2) TakeSnapshot() MapperScreenShot {
	res := Mapper2SS{
		PRGRAM: make([]uint8, len(m.PRGROM)),
		CHRROM: make([]uint8, len(m.CHRROM)),

		BankSelect: m.BankSelect,
		Mirroring:  m.Mirroring,
	}

	copy(res.CHRROM, m.CHRROM)
	copy(res.PRGRAM, m.PRGROM)

	return res
}

func (s Mapper2SS) LoadSS(m Mapper) {

	safe, ok := m.(*Mapper2)

	if !ok {
		log.Fatal("doie")
	}

	copy(safe.CHRROM, s.CHRROM)
	copy(safe.PRGROM, s.PRGRAM)

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
	PRGRAM [0x2000]uint8

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
	PRGRAM [0x2000]uint8

	shiftRegister uint8

	control  uint8
	PrgBank  uint8
	ChrBank0 uint8
	ChrRank1 uint8

	isRam bool
}

func (m *Mapper1) TakeSnapshot() MapperScreenShot {

	res := Mapper1SS{
		PRGROM: make([]uint8, len(m.CHRROM)),
		CHRROM: make([]uint8, len(m.PRGROM)),
	}

	return res
}

func (s Mapper1SS) LoadSS(m Mapper) {
	safe, ok := m.(*Mapper1)
	if !ok {
		log.Fatal("im dead bro")
	}

	safe.CHRROM = s.CHRROM
	safe.PRGRAM = s.PRGRAM
	safe.PRGRAM = s.PRGRAM

	safe.shiftRegister = s.shiftRegister
	safe.control = s.control
	safe.PrgBank = s.PrgBank
	safe.ChrBank0 = s.ChrBank0
	safe.ChrBank1 = s.ChrRank1

	safe.isRam = s.isRam
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

//Mapper 4

type Mapper4 struct {
	PRGROM []uint8
	CHRROM []uint8
	PRGRAM [0x2000]uint8

	BankSelect uint8
	BankData   [8]uint8
	Mirroring  uint8

	PRGMode uint8
	CHRMode uint8

	IRQLatch   uint8
	IRQCounter uint8
	IRQReload  bool
	IRQEnabled bool
	IRQPending bool
}

func (m *Mapper4) WritePRG(addr uint16, val uint8) {

	if addr < 0x8000 {
		m.PRGRAM[addr-0x6000] = val
	}

	if addr >= 0x8000 && addr <= 0x9FFF {
		if addr&1 == 0 { //even
			m.BankSelect = val
			m.PRGMode = getbit(val, 6)
			m.CHRMode = getbit(val, 7)
		} else { // odd
			slot := m.BankSelect & 0x07
			m.BankData[slot] = val
		}
	}

	if addr >= 0xA000 && addr <= 0xBFFE {
		if addr&1 == 0 {
			m.Mirroring = val & 1
		} else {
			//intetionally left

			fmt.Println("hitting prg ram protectg")
		}
	}

	if addr >= 0xC000 && addr <= 0xDFFE {
		if addr&1 == 0 { //even
			m.IRQLatch = val
		} else {
			m.IRQCounter = 0
			m.IRQReload = true
		}
	}

	if addr >= 0xE000 && addr <= 0xFFFE {
		if addr&1 == 0 {
			m.IRQEnabled = false
			m.IRQPending = false
		} else {
			m.IRQEnabled = true
		}
	}
}

func (m *Mapper4) ReadPRG(addr uint16) uint8 {
	if addr < 0x8000 {
		return m.PRGRAM[addr-0x6000]
	}
	var bank uint8
	var offset uint32

	switch {
	case addr < 0xA000:
		if m.PRGMode == 0 {
			bank = m.BankData[6] //r6
		} else {
			bank = uint8((len(m.PRGROM) / 0x2000) - 2) //second last
		}
		offset = uint32(addr - 0x8000)
	case addr < 0xC000:
		bank = m.BankData[7] //R7
		offset = uint32(addr - 0xA000)
	case addr < 0xE000:
		if m.PRGMode == 0 {
			bank = uint8((len(m.PRGROM) / 0x2000) - 2) // secon last
		} else {
			bank = m.BankData[6] //r6
		}
		offset = uint32(addr - 0xC000)
	default: // $E000 - $FFFF
		bank = uint8((len(m.PRGROM) / 0x2000) - 1) //last
		offset = uint32(addr - 0xE000)
	}

	newaddr := uint32(bank)*0x2000 + offset

	if int(newaddr) > len(m.PRGROM) {
		fmt.Println(bank, offset, addr)
	}

	return m.PRGROM[uint32(bank)*0x2000+offset]

}

func (m *Mapper4) ReadCHR(addr uint16) uint8 {
	var bank uint8
	offset := addr % 0x400

	if addr <= 0x07FF {
		if m.CHRMode == 0 {
			bank = m.BankData[0] //R0
		} else {
			if addr <= 0x03FF {
				bank = m.BankData[2] //R2
			} else {
				bank = m.BankData[3] //R3
			}
		}
	}

	if addr >= 0x0800 && addr <= 0x0FFF {
		if m.CHRMode == 0 {
			bank = m.BankData[1] // R1
		} else {
			if addr <= 0x0BFF {
				bank = m.BankData[4] //R4
			} else {
				bank = m.BankData[5] //R5
			}
		}
	}

	if addr >= 0x1000 && addr <= 0x17FF {
		if m.CHRMode == 0 {
			if addr <= 0x13FF {
				bank = m.BankData[2] //R2
			} else {
				bank = m.BankData[3] //R3
			}
		} else {
			bank = m.BankData[0] //R0
		}
	}

	if addr >= 0x1800 && addr <= 0x1FFF {
		if m.CHRMode == 0 {
			if addr <= 0x1BFF {
				bank = m.BankData[4] //R4
			} else {
				bank = m.BankData[5] // R5
			}
		} else {
			bank = m.BankData[1]
		}
	}

	return m.CHRROM[uint32(bank)*0x400+uint32(offset)]

}

func (m *Mapper4) WriteCHR(addr uint16, val uint8) {

}

func (m *Mapper4) extractCHR() []uint8 {
	return m.CHRROM
}

func (m *Mapper4) loadCHR(data []uint8) {
	m.CHRROM = data
}

func (m *Mapper4) getMirroring() uint8 {
	return uint8(m.Mirroring)
}
