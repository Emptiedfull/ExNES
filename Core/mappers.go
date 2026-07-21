package Core

import (
	"fmt"
	"log"
	"reflect"
)

type Mapper interface {
	ReadPRG(addr uint16) uint8
	WritePRG(addr uint16, val uint8)

	ReadCHR(addr uint16) uint8
	WriteCHR(addr uint16, val uint8)

	getMirroring() uint8
	TakeSnapshot(MapperScreenShot)
	CreateEmptySnapshot() MapperScreenShot

	HasBattery() bool
	GetPRGRAM() []uint8
}

type IrqClocker interface {
	clockIrqCounter(uint16)
	IRQPending() bool
}

type MapperScreenShot interface {
	LoadSS(m Mapper)
}

func getMapper(header []byte) int {
	high := header[7] & 0xF0
	low := (header[6] & 0xF0) >> 4

	return int(high | low)
}

func (c *Console) assignMapper(id int, prgData []byte, chrData []byte, mirroring uint8, hasBattery bool) {
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
			isRam:         c.Ppu.Mem.CHR_isRam,
			battery:       hasBattery,
		}
	case 2:
		m = &Mapper2{
			PRGROM:    prgData,
			CHRROM:    chrData,
			Mirroring: mirroring,

			BankSelect: 0,
		}
	case 4:
		fmt.Println("loading mmc3")
		m = &MMC3{
			PRGROM:     prgData,
			CHRROM:     chrData,
			PRGRAM:     make([]uint8, 0x2000),
			Mirroring:  mirroring,
			hasBattery: hasBattery,
		}
	case 163:
		m = &Mapper163{
			PRGROM:     prgData,
			PRGRAM:     make([]uint8, 0x2000),
			CHRRAM:     make([]uint8, 0x2000),
			hasBattery: hasBattery,
			Mirroring:  int(mirroring),
		}
	default:
		fmt.Println("unknown:", id)
	}

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

func (m *Mapper0) HasBattery() bool {
	return false
}

func (m *Mapper0) GetPRGRAM() []uint8 {
	return nil
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

	copy(safe.CHRROM, s.CHRROM)
	copy(safe.PRGROM, s.PRGROM)
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

func (m *Mapper2) HasBattery() bool {
	return false
}

func (m *Mapper2) GetPRGRAM() []uint8 {
	return nil
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

	isRam   bool
	battery bool
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

	isRam   bool
	battery bool
}

func (m *Mapper1) CreateEmptySnapshot() MapperScreenShot {
	res := Mapper1SS{
		PRGROM: make([]uint8, len(m.PRGROM)),
		CHRROM: make([]uint8, len(m.CHRROM)),
		PRGRAM: make([]uint8, 0x2000),

		isRam: m.isRam,
	}

	return &res
}

func (m *Mapper1) HasBattery() bool {
	return m.battery
}

func (m *Mapper1) GetPRGRAM() []uint8 {
	return m.PRGRAM
}

func (m *Mapper1) TakeSnapshot(s MapperScreenShot) {

	kind := reflect.TypeOf(s)

	snap, ok := s.(*Mapper1SS)
	if !ok {
		log.Fatal("wowowoow", kind)
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

type Mapper163 struct {
	PRGROM []uint8
	PRGRAM []uint8
	CHRRAM []uint8

	Mirroring int

	reg5000 uint8
	reg5200 uint8
	reg5300 uint8

	feedbackLatch bool

	hasBattery bool
}

func (m *Mapper163) ReadPRG(addr uint16) uint8 {
	switch {
	case addr >= 0x6000 && addr < 0x8000:
		return m.PRGRAM[addr-0x6000]
	case addr == 0x5500, addr == 0x5501:
		if m.feedbackLatch {
			return 0x00
		} else {
			return 0x04
		}
	case addr >= 0x8000:
		banks := uint32(len(m.PRGROM) / 0x8000)
		if banks == 0 {
			banks = 1
		}

		bankActive := m.getBank() % banks
		return m.PRGROM[bankActive*0x8000+uint32(addr-0x8000)]
	default:
		return 0
	}
}

func (m *Mapper163) WritePRG(addr uint16, val uint8) {
	switch {
	case addr >= 0x6000 && addr < 0x8000:
		m.PRGRAM[addr-0x6000] = val
		return
	case addr >= 0x5000 && addr <= 0x50FF: //its a fucking mask ik but idrc enough so u get a range fucksy
		m.reg5000 = val
		return
	case addr >= 0x5200 && addr <= 0x52FF:
		m.reg5200 = val
		return
	case addr >= 0x5300 && addr <= 0x53FF:
		m.reg5300 = val
		return
	case addr == 0x5100:
		m.feedbackLatch = (val & 0x04) != 0
		return
	case addr == 0x5101:
		if (val & 0x01) != 0 {
			m.feedbackLatch = !m.feedbackLatch
		}

		return

	}
}

func (m *Mapper163) ReadCHR(addr uint16) uint8 {
	return m.CHRRAM[addr&0x1FFF]
}

func (m *Mapper163) WriteCHR(addr uint16, val uint8) {
	m.CHRRAM[addr&0x1FFF] = val
}

func (m *Mapper163) getBank() uint32 {
	p := m.SwapBits(m.reg5000)
	hi := m.SwapBits(m.reg5200)

	var l2 uint8
	if m.reg5300&0x04 != 0 {
		l2 = p & 0x03
	} else {
		l2 = 0x03
	}

	m2 := (p >> 2) & 0x03
	h2 := hi & 0x03

	return (uint32(h2) << 4) | (uint32(m2) << 2) | uint32(l2)
}

func (m *Mapper163) SwapBits(val uint8) uint8 {
	if m.reg5300&0x01 == 0 {
		return val
	}

	b0 := val & 0x01
	b1 := (val >> 1) & 0x01
	return (val &^ 0x03) | b1 | b0<<1
}

func (m *Mapper163) CreateEmptySnapshot() MapperScreenShot {
	return Mapper2SS{}
}

func (m *Mapper163) TakeSnapshot(MapperScreenShot) {

}
func (m *Mapper163) getMirroring() uint8 {
	return uint8(m.Mirroring)
}

func (m *Mapper163) HasBattery() bool {
	return m.hasBattery
}

func (m *Mapper163) GetPRGRAM() []uint8 {
	return m.PRGRAM
}

type MMC3 struct {
	PRGROM []uint8
	PRGRAM []uint8
	CHRROM []uint8

	Mirroring uint8

	hasBattery bool

	BankSelect uint8
	BankReg    [8]uint8

	irqLatch   uint8
	irqCounter uint8
	irqEnable  bool
	irqReload  bool
	irqPending bool

	protect uint8

	prevA12 uint8

	edges int
}

type MMC3SS struct {
	PRGROM []uint8
	PRGRAM []uint8
	CHRROM []uint8

	Mirroring uint8

	BankSelect uint8
	BankReg    [8]uint8

	irqLatch   uint8
	irqCounter uint8
	irqEnable  bool
	irqReload  bool
	irqPending bool

	protect uint8

	prevA12 uint8
}

func (s MMC3SS) LoadSS(m Mapper) {
	safe, ok := m.(*MMC3)

	if !ok {
		log.Fatal("doie")
	}

	copy(safe.CHRROM, s.CHRROM)
	copy(safe.PRGROM, s.PRGROM)
	copy(safe.PRGRAM, s.PRGRAM)

	safe.Mirroring = s.Mirroring

	safe.BankSelect = s.BankSelect
	safe.BankReg = s.BankReg

	safe.irqLatch = s.irqLatch
	safe.irqCounter = s.irqCounter
	safe.irqEnable = s.irqEnable
	safe.irqPending = s.irqPending
	safe.irqReload = s.irqReload

	safe.protect = s.protect

	safe.prevA12 = s.prevA12

}

func (m *MMC3) CreateEmptySnapshot() MapperScreenShot {
	s := MMC3SS{
		CHRROM: make([]uint8, len(m.CHRROM)),
		PRGROM: make([]uint8, len(m.PRGROM)),
		PRGRAM: make([]uint8, len(m.PRGRAM)),
	}
	return &s
}

func (m *MMC3) WritePRG(addr uint16, val uint8) {
	even := addr%2 == 0

	switch {
	case addr >= 0x6000 && addr < 0x8000:
		m.PRGRAM[addr-0x6000] = val
		return

	case addr >= 0x8000 && addr < 0xA000:
		if even {
			m.BankSelect = val

		} else {
			m.BankReg[m.BankSelect&0x07] = val

		}
		return
	case addr >= 0xA000 && addr < 0xC000:
		if even {
			if val&0x01 == 0 {
				m.Mirroring = 3
			} else {
				m.Mirroring = 2
			}
		} else {
			m.protect = val
		}
		return
	case addr >= 0xC000 && addr < 0xE000:
		if even {
			m.irqLatch = val

		} else {
			m.irqReload = true

		}
		return
	case addr >= 0xE000:
		if even {
			m.irqEnable = false
			m.irqPending = false

		} else {
			m.irqEnable = true

		}

		return
	}

}

func (m *MMC3) ReadPRG(addr uint16) uint8 {
	if addr >= 0x6000 && addr < 0x8000 {
		return m.PRGRAM[addr-0x6000]
	}

	if addr < 0x8000 {
		return 0
	}

	numBanks := uint32(len(m.PRGROM) / 0x2000)
	last := numBanks - 1
	secondLast := numBanks - 2

	slot := (addr - 0x8000) / 0x2000
	offset := uint32(addr % 0x2000)

	prgMode := (m.BankSelect >> 6) & 0x01

	var bank uint32
	switch {
	case slot == 0:
		if prgMode == 0 {
			bank = uint32(m.BankReg[6]) % numBanks

		} else {
			bank = secondLast
		}
	case slot == 1:
		bank = uint32(m.BankReg[7]) % numBanks
	case slot == 2:
		if prgMode == 0 {
			bank = secondLast
		} else {
			bank = uint32(m.BankReg[6]) % numBanks
		}
	default:
		bank = last
	}

	return m.PRGROM[bank*0x2000+offset]
}

func (m *MMC3) ReadCHR(addr uint16) uint8 {
	numBanks := uint32(len(m.CHRROM) / 0x400)
	chrMode := (m.BankSelect >> 7) & 0x01

	var bank uint32
	var offset uint32

	region := uint32(addr / 0x400)
	offset = uint32(addr % 0x400)

	if chrMode == 0 {
		switch {
		case region < 2:
			bank = (uint32(m.BankReg[0]) &^ 1) + region
		case region < 4:
			bank = (uint32(m.BankReg[1]) &^ 1) + (region - 2)
		default:
			bank = uint32(m.BankReg[region-2])
		}
	} else {
		switch {
		case region < 4:
			bank = (uint32(m.BankReg[region+2]))
		case region < 6:
			bank = (uint32(m.BankReg[0] &^ 1)) + (region - 4)
		default:
			bank = (uint32(m.BankReg[1] &^ 1)) + (region - 6)
		}
	}

	bank %= numBanks
	return m.CHRROM[bank*0x400+offset]

}

func (m *MMC3) WriteCHR(addr uint16, val uint8) {

}

func (m *MMC3) getMirroring() uint8 {
	return m.Mirroring
}

func (m *MMC3) clockIrqCounter(addr uint16) {
	a12 := uint8((addr >> 12) & 0x01)

	if m.prevA12 == 0 && a12 == 1 {
		// fmt.Println("A12 raw edge")
		m.edges++
		m.clockScanlineCounter()
	}

	m.prevA12 = a12
}

func (m *MMC3) clockScanlineCounter() {
	if m.irqCounter == 0 || m.irqReload {
		m.irqCounter = m.irqLatch
		m.irqReload = false
	} else {
		m.irqCounter--
	}

	if m.irqCounter == 0 && m.irqEnable {
		// fmt.Println("MMC3 IRQ fired")
		m.irqPending = true
	}
}

func (m *MMC3) IRQPending() bool {
	return m.irqPending
}

func (m *MMC3) GetPRGRAM() []uint8 {
	return m.PRGRAM
}

func (m *MMC3) HasBattery() bool {
	return false
}

func (m *MMC3) TakeSnapshot(s MapperScreenShot) {
	res, ok := s.(*MMC3SS)

	if !ok {
		fmt.Println("FUCK FUCK FUCK FUCK")

		return
	}

	copy(res.CHRROM, m.CHRROM)
	copy(res.PRGROM, m.PRGROM)
	copy(res.PRGRAM, m.PRGRAM)

	res.BankReg = m.BankReg
	res.BankSelect = m.BankSelect

	res.Mirroring = m.Mirroring

	res.irqCounter = m.irqCounter
	res.irqReload = m.irqReload
	res.irqLatch = m.irqLatch
	res.irqEnable = m.irqEnable
	res.irqPending = m.irqPending

	res.protect = m.protect

	res.prevA12 = m.prevA12

}
