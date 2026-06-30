package Core

import (
	"fmt"
	"sync"
	"time"
)

type addressingMode int

const (
	Implied addressingMode = iota
	Accumulator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
	Relative
)

type Debugger struct {
	Console     *Console
	Disassembly map[uint16]AssemblyLine
	DMux        sync.Mutex
}

type ScreenInfo struct {
	Buffer *[245760]uint8
}

type cpustate struct {
	Pc     uint16    `json:"pc"`
	S      uint8     `json:"s"`
	A      uint8     `json:"a"`
	X      uint8     `json:"x"`
	Y      uint8     `json:"y"`
	P      uint8     `json:"p"`
	Flags  FlagState `json:"flags"`
	Cycles int       `json:"cycles"`
	Frames int       `json:"frames"`
	Ram    [][]int   `json:"-"`
}

type FlagState struct {
	Carry     bool `json:"carry"`
	Overflow  bool `json:"overflow"`
	Interrupt bool `json:"interrupt"`
	Zero      bool `json:"zero"`
	Decimal   bool `json:"decimal"`
	Negative  bool `json:"negative"`
}

func (d *Debugger) StartDebugConsole() {

	targetTime := time.Now()
	defer fmt.Println("console stopped for some reason")

	for {

		if d.Console.Paused {
			time.Sleep(100 * time.Millisecond)
			targetTime = time.Now()

			continue
		}

		now := time.Now()
		for now.After(targetTime) {

			d.Console.RunFrame()
			targetTime = targetTime.Add(time.Duration(nsPerFrame))

		}

		timeLeft := time.Until(targetTime)
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}

	}
}

func (d *Debugger) RunDebugFrame() {
	Tagretframe := d.Console.Ppu.Frame + 1
	for d.Console.Ppu.Frame != Tagretframe {
		d.DebugTick()

	}

	// d.Console.RunDisplayUpdates()
}

func (d *Debugger) DebugTick() {

	d.Console.tick()
	// if d.Console.Cpu.currentstep == 0 {
	// 	d.DisAssemble(d.Console.Cpu.PC)
	// }

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

// func (c *Console) DebugChrRom() {

// 	Rom := c.Ppu.mem.mapper.extractCHR()

// 	tiles := min(len(Rom)/16, 512)

// 	for t := range tiles {
// 		tileBytes := Rom[(t * 16):((t + 1) * 16)]
// 		offset := t * 64

// 		for y := range 8 {
// 			low := tileBytes[y]
// 			high := tileBytes[y+8]

// 			for x := range 8 {
// 				bit1 := (low >> (7 - x)) & 0x01
// 				bit2 := (high >> (7 - x)) & 0x01

// 				colorIndex := (bit2 << 1) | bit1

// 				c.Ppu.DebugBuffer[offset+(y*8+x)] = colorIndex

// 			}
// 		}
// 	}

// }

// func (c *Console) DebugNameTable() {
// 	fmt.Println("debugging namtable")
// 	Rom := c.Ppu.mem.mapper.extractCHR()
// 	c.Ppu.DebugBuffer = make([]uint8, 512*480)

// 	var PatternOffset uint16 = 0
// 	if c.Ppu.mem.register.BgPattern {
// 		PatternOffset = 4096
// 	}

// 	for quadrant := 0; quadrant < 4; quadrant++ {
// 		var baseAddr uint16 = 0x2000 + uint16(quadrant)*0x400

// 		quadXOffset := (quadrant % 2) * 256
// 		quadYOffset := (quadrant / 2) * 240

// 		for row := 0; row < 30; row++ {
// 			for col := 0; col < 32; col++ {

// 				ntIdx := uint16(row*32 + col)

// 				tileID := uint16(c.Ppu.read(baseAddr + ntIdx))

// 				chrOffset := (tileID * 16) + PatternOffset

// 				if int(chrOffset+16) <= len(Rom) {
// 					tileBytes := Rom[chrOffset : chrOffset+16]

// 					for y := 0; y < 8; y++ {
// 						low := tileBytes[y]
// 						high := tileBytes[y+8]

// 						for x := 0; x < 8; x++ {
// 							bit1 := (low >> (7 - x)) & 0x01
// 							bit2 := (high >> (7 - x)) & 0x01
// 							colorIndex := (bit2 << 1) | bit1

// 							pixelX := quadXOffset + (col * 8) + x
// 							pixelY := quadYOffset + (row * 8) + y

// 							c.Ppu.DebugBuffer[pixelY*512+pixelX] = colorIndex
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
