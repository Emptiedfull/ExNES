package Core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type console struct {
	Cpu        *cpu
	Ppu        *ppu
	JoyPad     *joyPad
	OpenBusVal uint8

	ScreenChannel chan ScreenInfo

	ready bool
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

	c.Cpu.mem.FillArr(0, prgData)

	c.Ppu.mem.chrROM = make([]uint8, 8192)
	if _, err := io.ReadFull(file, c.Ppu.mem.chrROM); err != nil {
		return fmt.Errorf("failed to read mem: %w", err)
	}

	return nil

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

	c.Cpu.mem = &bus{
		cpu: c.Cpu,
	}

	return c
}

var nsPerFrame = int64(float64(time.Second.Nanoseconds()) / 60.0988) // 2

func (c *console) StartConsoleCycle() {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	lastTime := time.Now()
	var lag int64 = 0

	var framecount = 0
	var fpsMilli = 0

	for range ticker.C {
		fpsMilli++
		elapsed := time.Since(lastTime).Nanoseconds()
		lastTime = time.Now()
		lag += int64(elapsed)

		for lag > nsPerFrame {
			framecount++
			for range 29781 {
				c.tick()
			}
			lag -= nsPerFrame

			c.runDisplayUpdates()
		}

		if fpsMilli > 1000 {
			fpsMilli = 0
			fmt.Println(framecount)
		}

	}

}

func (d *Debugger) StepCycles(cycles int) {
	for range cycles {
		d.Console.tick()

		d.Disassembly[d.Console.Cpu.PC] = d.DisAssemble(d.Console.Cpu.PC)
	}
}

func (c *console) runDisplayUpdates() {

	S := ScreenInfo{
		Buffer: c.Ppu.backBuffer,
	}

	c.ScreenChannel <- S

}

func (c *console) tick() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.totalCycles++
	} else {
		c.Cpu.tick()
	}

	c.Ppu.Tick()
	c.Ppu.Tick()
	c.Ppu.Tick()

}

func (b *bus) FillArr(addr uint16, data []byte) error {

	size := len(data)
	b.external = make([]byte, size)

	copy(b.external[addr:], data)
	return nil

}

func (b *bus) GetHistory() []cycleStep {
	return nil
}

func (b *bus) ClearHistory() {

}

func (b *bus) Set(addr uint16, val uint8) {

}

func (b *bus) Get(addr uint16) uint8 {
	return 0
}
