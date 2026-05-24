package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type console struct {
	Cpu        *cpu
	Ppu        *ppu
	OpenBusVal uint8

	CanvasImage *ebiten.Image

	ready bool
}

func (c *console) loadROM(filepath string) error {
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

	prgSize := int(header[4]) * 16384
	//chrBanks := int(header[5])

	prgData := make([]byte, prgSize)
	if _, err := io.ReadFull(file, prgData); err != nil {
		return fmt.Errorf("error reading prgData: %w", err)
	}

	c.Ppu.mem.chrROM = make([]uint8, 8192)
	if _, err := io.ReadFull(file, c.Ppu.mem.chrROM); err != nil {
		return fmt.Errorf("failed to read mem: %w", err)
	}

	for i := 0; i < 16; i++ {
		fmt.Printf("Byte %02d: 0x%02X\n", i, c.Ppu.mem.chrROM[i])
	}

	//c.Cpu.mem.external = prgData
	return nil

}

func (g *console) Update() error {

	if !g.ready {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Ppu.DumpBackBufferToFile()
	}

	if !g.Ppu.DrawFlg {
		for range 89000 {
			g.tick()
		}
	}

	return nil
}

// func (g *console) Draw(screen *ebiten.Image) {
// 	if g.Ppu.DrawFlg {
// 		fmt.Println("drawing")
// 		screen.WritePixels(g.Ppu.backBuffer)
// 		g.Ppu.DrawFlg = false
// 	}
// }

func initializeConsole() *console {
	c := &console{
		CanvasImage: ebiten.NewImage(256, 240),
		Ppu: &ppu{
			backBuffer:  make([]byte, 256*240*4),
			frontBuffer: make([]byte, 256*240*4),
		},
		Cpu:        &cpu{},
		OpenBusVal: 0,
	}

	c.Cpu.console = c
	c.Ppu.console = c

	c.Cpu.mem = &bus{}
	//c.Cpu.mem.cpu = c.Cpu

	return c
}

func (b *bus) GetHistory() []cycleStep {
	return nil
}

func (b *bus) ClearHistory() {
	return
}

func (b *bus) Set(addr uint16, val uint8) {
	return
}

func (b *bus) Get(addr uint16) uint8 {
	return 0
}

func main() {
	c := initializeConsole()
	fmt.Println("running")
	err := c.loadROM("C:/Users/user/ExNES/games/dk.nes")
	//err := c.loadROM("C:/Users/user/ExNES/nestest.nes")
	c.Cpu.reset()

	fmt.Println(err)

	ebiten.SetWindowSize(768, 384)

	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}
