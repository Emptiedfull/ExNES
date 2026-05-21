package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func (g *console) Draw(screen *ebiten.Image) {

	//buffer := make([]byte, 256*128*4)

	g.CanvasImage.WritePixels(g.Ppu.frontBuffer)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(3, 3)
	screen.DrawImage(g.CanvasImage, opts)
}

func (g *console) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256 * 3, 256 * 3
}

func debug() {
	c := &console{
		Cpu:         &cpu{},
		Ppu:         &ppu{},
		CanvasImage: ebiten.NewImage(256, 128),
	}
	dummyTileBytes := []byte{
		0xFF, 0x81, 0x90, 0x99, 0x99, 0x81, 0x81, 0xFF,
		0xFF, 0x00, 0x3C, 0x24, 0x24, 0x3C, 0x00, 0xFF,
	}

	c.Ppu.mem = ppu_mem{}
	c.Ppu.mem.chrROM = make([]byte, 8192)

	for i, b := range dummyTileBytes {

		c.Ppu.mem.chrROM[i] = b
	}

	ebiten.SetWindowSize(768, 384)

	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}
