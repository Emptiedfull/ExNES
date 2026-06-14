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
	Console     *console
	Disassembly map[uint16]AssemblyLine
	DMux        sync.Mutex

	RecentHistory SnapshotBuffer
}

type SnapshotBuffer struct {
	Frame int

	Data  [200]snapshot
	Index int
}

type ScreenInfo struct {
	Buffer []uint8
}

type AssemblyLine struct {
	Opcode      opCode `json:"Opcode"`
	Disassembly string `json:"disassembly"`
	Val         uint8  `json:"val,omitempty"`
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
	fmt.Println("console started")
	defer fmt.Println("console stopped for some reason")

	for {

		if d.Console.Paused {
			time.Sleep(100 * time.Millisecond)
			targetTime = time.Now()

			continue
		}

		now := time.Now()
		for now.After(targetTime) {
			d.RunDebugFrame()

			targetTime = targetTime.Add(time.Duration(nsPerFrame))

			d.Console.RunDisplayUpdates()
			d.AddSnapshot()
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

	d.Console.RunDisplayUpdates()
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
