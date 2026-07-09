package Core

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Console struct {
	Cpu *Cpu
	Ppu *ppu
	Apu *APU

	Player1 *joyPad
	Player2 *joyPad

	OpenBusVal uint8

	ScreenChannel chan ScreenInfo

	Snapshots SnapshotBuffer

	ready bool

	Paused   bool
	pausedMu sync.Mutex

	DebugMode int
	Debugger  Debugger

	mapper Mapper

	palette FullPalette
}

func Init() {
	fmt.Print("Core V2")

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

func (c *Console) LoadRom(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("unable to read rom file: %v", err)
	}

	err = c.InitRom(file)
	if err != nil {
		return fmt.Errorf("Invalid rom data: %v", err)
	}

	c.Cpu.Reset()

	return nil
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

	header := make([]byte, 16)
	if _, err := io.ReadFull(data, header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	mapper := getMapper(header)

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

	c.assignMapper(mapper, prgData, chrData, uint8(mirroring))

	c.SetUpSnapshots()

	return nil
}

func InitializeConsole() *Console {

	c := &Console{

		Ppu: &ppu{},
		Cpu: &Cpu{},

		OpenBusVal: 0,
		palette:    loadFPal("/Users/test/Projects/ExNES/Core/ntsc.pal"),
	}

	c.Ppu.BackBuffer = make([]uint8, 245760)

	c.Player1 = &joyPad{}
	c.Player2 = &joyPad{}

	c.Cpu.console = c
	c.Ppu.console = c
	c.Apu = NewApu(44100, c)

	c.OpenBusVal = 0
	c.Cpu.fetchNew = true
	c.ScreenChannel = make(chan ScreenInfo, 100)

	c.Snapshots.Data = make([]Snapshot, 100)

	c.Cpu.Mem = &bus{
		Cpu: c.Cpu,
	}

	return c
}

var nsPerFrame = int64(float64(time.Second.Nanoseconds()) / 60.0988)

func (c *Console) StartConsoleCycle() {

	targetTime := time.Now()
	defer fmt.Println("console stopped for some reason")

	var framecount = 0
	start := time.Now()

	for {

		if c.Paused {
			time.Sleep(100 * time.Millisecond)
			targetTime = time.Now()
			continue
		}

		now := time.Now()
		for now.After(targetTime) {
			framecount++

			c.RunFrame()
			c.RunDisplayUpdates()

			targetTime = targetTime.Add(time.Duration(nsPerFrame))
		}

		timeLeft := time.Until(targetTime)
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}

		if framecount == 60 {
			fmt.Println(time.Since(start))
		}
	}

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
		S := ScreenInfo{
			Buffer: c.Ppu.BackBuffer,
		}
		c.Ppu.ScreenChanged = false
		c.ScreenChannel <- S
	}

}

func (c *Console) tick() {
	// c.RunDisplayUpdates()
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.TotalCycles++
	} else {
		c.Cpu.tick()
	}
	for range 3 {
		c.Ppu.step()
	}

	c.Apu.tick()

	if c.Apu.IRGPending || c.Apu.Dmc.IRGPending {
		c.Cpu.triggerIRQ()
	}

}

func (c *Console) RunFrame() {

	targetFrame := c.Ppu.Frame + 1
	for targetFrame != c.Ppu.Frame {
		c.tick()
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
