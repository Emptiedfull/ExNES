package Core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type console struct {
	Cpu        *cpu
	Ppu        *ppu
	JoyPad     *joyPad
	OpenBusVal uint8

	ScreenChannel chan ScreenInfo

	ready bool

	Paused   bool
	pausedMu sync.Mutex
}

func (c *console) Pause() {
	c.pausedMu.Lock()
	defer c.pausedMu.Unlock()
	c.Paused = true
}

func (c *console) UnPause() {
	c.pausedMu.Lock()
	defer c.pausedMu.Unlock()
	c.Paused = false
}

func (c *console) LoadROM(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	header := make([]byte, 16)
	if _, err := io.ReadFull(file, header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	mapper := getMapper(header)
	fmt.Println(mapper)

	mirroring := 2
	if (header[6] & 0x01) == 0 {
		mirroring = 3
	}

	trainerByte := header[6] & 0x04
	if trainerByte != 0 {
		tainer := make([]byte, 512)
		io.ReadFull(file, tainer)
	}

	magic := header[:4]
	if !bytes.Equal(magic, []byte{'N', 'E', 'S', 0x1A}) {
		return fmt.Errorf("invalid nes header")
	}

	prgBanks := int(header[4])
	prgSize := prgBanks * 16384
	prgData := make([]byte, prgSize)
	if _, err := io.ReadFull(file, prgData); err != nil {
		return fmt.Errorf("error reading prgData: %w", err)
	}

	chrBanks := int(header[5])

	var chrData []uint8
	if chrBanks == 0 {
		chrData = make([]byte, 8192)
		c.Ppu.mem.CHR_isRam = true
	} else {
		chrSize := chrBanks * 8192
		chrData = make([]byte, chrSize)
		if _, err := io.ReadFull(file, chrData); err != nil {
			return fmt.Errorf("failed to read mem: %w", err)
		}
	}
	fmt.Println("banks", chrBanks)

	c.assignMapper(mapper, prgData, chrData, uint8(mirroring))

	return nil
}

func getMapper(header []byte) int {
	high := header[7] & 0xF0
	low := (header[6] & 0xF0) >> 4

	return int(high | low)
}

func InitializeConsole() *console {
	c := &console{

		Ppu: &ppu{
			backBuffer:  make([]byte, 256*240*4),
			frontBuffer: make([]byte, 256*240*4),
		},
		Cpu:        &cpu{},
		OpenBusVal: 0,
	}
	c.JoyPad = &joyPad{}
	c.Cpu.console = c
	c.Ppu.console = c
	c.OpenBusVal = 0
	c.Cpu.fetchNew = true
	c.ScreenChannel = make(chan ScreenInfo, 100)

	c.Cpu.mem = &bus{
		cpu: c.Cpu,
	}

	return c
}

var nsPerFrame = int64(float64(time.Second.Nanoseconds()) / 60.0988)

func (c *console) StartConsoleCycle() {

	targetTime := time.Now()
	defer fmt.Println("console stopped for some reason")

	var framecount = 0

	for {

		if c.Paused {
			time.Sleep(100 * time.Millisecond)
			targetTime = time.Now()
			continue
		}

		now := time.Now()
		for now.After(targetTime) {
			framecount++

			for range 29781 {
				c.tick()
			}

			targetTime = targetTime.Add(time.Duration(nsPerFrame))

			c.RunDisplayUpdates()
		}

		timeLeft := time.Until(targetTime)
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}
	}

}

func (d *Debugger) StepCycles(cycles int) {
	target := d.Console.Cpu.TotalCycles + cycles

	for d.Console.Cpu.TotalCycles < target {
		if d.Console.Cpu.currentstep == 0 {
			d.Disassembly[d.Console.Cpu.PC] = d.DisAssemble(d.Console.Cpu.PC)
		}
		d.Console.tick()

	}
}

func (c *console) RunDisplayUpdates() {
	if c.Ppu.screenChanged {
		S := ScreenInfo{
			Buffer: c.Ppu.backBuffer,
		}
		c.Ppu.screenChanged = false
		c.ScreenChannel <- S
	}

}

func (c *console) tick() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.TotalCycles++
	} else {
		c.Cpu.tick()
	}
	for range 3 {
		c.Ppu.step()
	}

}

func (c *console) runFrame() {
	targetFrame := c.Ppu.Frame + 1
	for targetFrame != c.Ppu.Frame {
		c.tick()
	}
}
