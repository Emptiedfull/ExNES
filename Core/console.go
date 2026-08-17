package Core

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"sync"
)

type Console struct {
	Cpu *Cpu
	Ppu *ppu
	Apu *APU

	CheatEngine    *GameGenieEngine
	DatabaseEngine *RomDB
	PalleteEngine  *PalleteEngine

	romHash uint32

	Player1 *joyPad
	Player2 *joyPad

	OpenBusVal uint8

	ScreenChannel chan []uint32

	Snapshots SnapshotBuffer

	Paused   bool
	pausedMu sync.Mutex

	Debug *DebugEngine

	mapper  Mapper
	LoadRam func(Name string, rom []uint8)
}

func (c *Console) GetRam() []uint8 {
	return c.mapper.GetPRGRAM()
}

func Init() {
	fmt.Print("Core V2")

}

func (c *Console) Step() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.TotalCycles++
	} else {
		c.Cpu.Tick()
	}

	for range 3 {
		c.Ppu.step()
	}

	c.Apu.tick()

	c.Cpu.irqLine = c.Apu.IRGPending || c.Apu.Dmc.IRGPending

	if m, ok := c.mapper.(IrqClocker); ok {

		if m.IRQPending() {

			c.Cpu.irqLine = true
		}

	}

}

func (c *Console) GetMapper() Mapper {
	return c.mapper
}

func (c *Console) Pause() {
	c.pausedMu.Lock()
	defer c.pausedMu.Unlock()
	c.Paused = true
}

func (c *Console) UnPause() {
	c.pausedMu.Lock()
	defer c.pausedMu.Unlock()
	c.Paused = false
}

func (c *Console) OpenRomFile(filepath string) (io.Reader, error) {
	file, err := os.Open(filepath)
	return file, err
}

func LoadRomData(data []uint8) (io.Reader, error) {
	reader := bytes.NewReader(data)
	return reader, nil
}

func (c *Console) InitRom(data io.Reader) error {
	c.CheatEngine.Reset()
	header := make([]byte, 16)
	if _, err := io.ReadFull(data, header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	mapper := getMapper(header)
	hasBattery := header[6]&0x02 != 0

	mirroring := 2
	if (header[6] & 0x01) == 0 {
		mirroring = 3
	}
	c.Ppu.mirroring = mirroring

	trainerByte := header[6] & 0x04
	if trainerByte != 0 {
		tainer := make([]byte, 512)
		io.ReadFull(data, tainer)
	}

	magic := header[:4]
	if !bytes.Equal(magic, []byte{'N', 'E', 'S', 0x1A}) {
		return fmt.Errorf("invalid nes header")
	}

	h := crc32.NewIEEE()
	data = io.TeeReader(data, h)

	prgBanks := int(header[4])
	prgSize := prgBanks * 16384
	prgData := make([]byte, prgSize)
	if _, err := io.ReadFull(data, prgData); err != nil {
		return fmt.Errorf("error reading prgData: %w", err)
	}

	chrBanks := int(header[5])

	var chrData []uint8
	if chrBanks == 0 {
		chrData = make([]byte, 8192)
		c.Ppu.Mem.CHR_isRam = true
	} else {
		chrSize := chrBanks * 8192
		chrData = make([]byte, chrSize)
		if _, err := io.ReadFull(data, chrData); err != nil {
			return fmt.Errorf("failed to read Mem: %w", err)
		}
	}

	io.Copy(io.Discard, data)
	c.romHash = h.Sum32()

	name, found := c.DatabaseEngine.Get(c.romHash)
	if found {
		c.CheatEngine.LoadCheats(name)

	}

	c.assignMapper(mapper, prgData, chrData, uint8(mirroring), hasBattery)

	if c.mapper == nil {
		return fmt.Errorf("Mapper not supported")
	}

	if c.LoadRam != nil {

		if hasBattery {
			c.LoadRam(c.GetHash(), c.mapper.GetPRGRAM())
		}

	}

	c.SetUpSnapshots()

	return nil
}

func (c *Console) GetHash() string {
	return fmt.Sprintf("%08x", c.romHash)
}

func (c *Console) GetName() string {
	name, found := c.DatabaseEngine.Get(c.romHash)
	if !found {
		return ""
	}
	return name[0 : len(name)-4]
}

func InitializeConsole() *Console {

	c := &Console{

		Ppu: &ppu{},
		Cpu: &Cpu{},

		OpenBusVal: 0,
		PalleteEngine: &PalleteEngine{
			ListPal: make([]PalleteEntry, 0),
			Loaded:  0,
		},
	}

	c.PalleteEngine.Init()

	db, err := LoadDB()
	if err != nil {
		fmt.Println("error intializing the db:", err)
	}

	c.DatabaseEngine = db
	c.CheatEngine = InitCheat()

	c.Ppu.NewBuffer = make([]uint32, 256*240)
	c.Ppu.ShowBuffer = make([]uint32, 256*240)

	c.Player1 = &joyPad{}
	c.Player2 = &joyPad{}

	c.Cpu.console = c
	c.Ppu.console = c
	c.Apu = NewApu(44100, c)

	c.OpenBusVal = 0
	c.Cpu.fetchNew = true
	c.ScreenChannel = make(chan []uint32, 100)

	c.Snapshots.Data = make([]Snapshot, 100)

	c.Cpu.Mem = &bus{
		Cpu: c.Cpu,
	}

	return c
}

func (c *Console) PowerCycle() {
	c.Cpu.Reset()
	c.Cpu.Mem = &bus{
		Cpu:      c.Cpu,
		internal: [2048]byte{},
	}

	c.Ppu.reset()
	c.Apu = NewApu(c.Apu.SampleRate, c)
}

func (c *Console) RunDisplayUpdates() {

	if c.Ppu.ScreenChanged {
		c.ScreenChannel <- c.Ppu.ShowBuffer
		c.Ppu.ScreenChanged = false
	}
}

func (c *Console) RunFrame() {
	target := c.Ppu.Frame + 1
	for target != c.Ppu.Frame {
		c.Step()
	}
}

func Quickstart(filepath string) *Console {
	c := InitializeConsole()
	file, err := c.OpenRomFile(filepath)
	fmt.Println(filepath)
	if err != nil {
		log.Fatalln("error reading file data", err)
	}
	c.InitRom(file)
	c.Cpu.Reset()

	return c
}
