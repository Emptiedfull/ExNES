package Core

import "fmt"

// Contains utility for the snapshot function

type snapshot struct {
	Frame_no int
	Cycles   int
	CpuState CpuSnapshot
	PpuState PpuSnapshot
	ApuState ApuSnapshot

	mapperState MapperScreenShot
}

type SnapshotBuffer struct {
	Frame int
	// Data  [100]snapshot

	Data []snapshot

	Index int
}

type ApuSnapshot struct {
	Pulse1 PulseChannel
	Pulse2 PulseChannel

	Tri Triangle

	Noise NoiseChannel

	DMC DMC

	counter frameCounter

	CycleAcc       float64
	timerTickLatch bool

	IRGPending bool
}

type CpuSnapshot struct {
	PC uint16
	S  uint8
	P  uint8
	A  uint8
	X  uint8
	Y  uint8

	mem cpu_mem
	t   temp

	nmiPending   bool
	executingNmi bool
	nmiStep      int

	currentOp   uint8
	currentstep int
	totalCycles int
	Stall       int

	isJamming bool
}

type cpu_mem struct {
	internal [2048]byte
	external []byte
}

type PpuSnapshot struct {
	mem ppu_mem

	Dot      int
	Scanline int
	Frame    int

	mirroring       int
	mirroringChange bool

	BackBuffer []uint8

	screenChanged bool
}

func (c *Console) SetUpSnapshots() {
	for range len(c.Snapshots.Data) {
		fmt.Println("creating empty snapshots ")
		c.AddSnapshot(c.createEmptySnapshot())
	}

}

func (a *APU) TakeSnapshot(s *ApuSnapshot) {
	s.Pulse1 = *a.Pulse1
	s.Pulse2 = *a.Pulse2

	s.Tri = *a.Triangle
	s.Noise = *a.Noise
	s.DMC = *a.Dmc

	s.counter = a.Counter
	s.CycleAcc = a.CycleAcc

	s.timerTickLatch = a.timerTickLatch

	s.IRGPending = a.IRGPending
}

func (a *APU) LoadSnapshot(s *ApuSnapshot) {
	*a.Pulse1 = s.Pulse1
	*a.Pulse2 = s.Pulse2

	*a.Triangle = s.Tri

	*a.Noise = s.Noise
	*a.Dmc = s.DMC

	a.Counter = s.counter
	a.CycleAcc = s.CycleAcc
	a.timerTickLatch = s.timerTickLatch
	a.IRGPending = s.IRGPending
}

func (c *Console) TakeSnapshot() {
	c.PopulateSnapshot(&c.Snapshots.Data[c.Snapshots.Index])
	c.Snapshots.Index++
}

func (c *Console) createEmptySnapshot() snapshot {
	return snapshot{
		Frame_no:    0,
		Cycles:      0,
		mapperState: c.mapper.CreateEmptySnapshot(),
		CpuState:    CpuSnapshot{},
		PpuState:    PpuSnapshot{},
		ApuState:    ApuSnapshot{},
	}

}

func (p *ppu) TakePpuSnapshot(S *PpuSnapshot) {

	S.mem = p.mem

	S.Dot = p.Dot
	S.Scanline = p.Scanline
	S.Frame = p.Frame

	copy(S.BackBuffer, p.BackBuffer)

	S.screenChanged = p.ScreenChanged

}

func (b *SnapshotBuffer) AddSnapshot(shot snapshot) {
	b.Data[b.Index] = shot
	b.Index = (b.Index + 1) % len(b.Data)

}

func (c cpu) TakeCpuSnapshot(S *CpuSnapshot) {

	S.PC = c.PC
	S.S = c.S
	S.P = c.P
	S.X = c.X
	S.Y = c.Y

	S.mem = cpu_mem{
		internal: c.mem.returnInternal(),
	}

	S.t = c.console.Cpu.temp

	S.nmiPending = c.nmiPending
	S.executingNmi = c.executingNmi
	S.nmiStep = c.nmiStep

	S.currentOp = c.currentOp
	S.currentstep = c.currentstep
	S.totalCycles = c.TotalCycles

	S.Stall = c.Stall

	S.isJamming = c.isJamming

}

func (c *Console) PopulateSnapshot(s *snapshot) {

	if c.mapper == nil || c.Cpu == nil || c.Ppu == nil || c.Apu == nil {
		return
	}

	c.Cpu.TakeCpuSnapshot(&s.CpuState)
	c.Ppu.TakePpuSnapshot(&s.PpuState)
	c.Apu.TakeSnapshot(&s.ApuState)
	c.mapper.TakeSnapshot(&s.mapperState)
	s.Frame_no = c.Ppu.Frame
	s.Cycles = c.Cpu.TotalCycles
}

func (c *Console) AddSnapshot(s snapshot) {
	s.Frame_no = c.Snapshots.Frame
	c.Snapshots.Frame++
	c.Snapshots.AddSnapshot(s)
}

func (c *Console) LoadSnapshot(snap snapshot) {

	c.Cpu.LoadCpuSnapshot(snap.CpuState)
	c.Ppu.LoadPpuSnapshot(snap.PpuState)
	snap.mapperState.LoadSS(c.mapper)
}

func (p *ppu) LoadPpuSnapshot(snap PpuSnapshot) {
	p.mem = snap.mem

	p.Dot = snap.Dot
	p.Scanline = snap.Scanline
	p.Frame = snap.Frame

	copy(p.BackBuffer, snap.BackBuffer)

	p.ScreenChanged = snap.screenChanged
}

func (c *cpu) LoadCpuSnapshot(snap CpuSnapshot) {
	c.PC = snap.PC
	c.S = snap.S
	c.P = snap.P
	c.X = snap.X
	c.Y = snap.Y

	c.mem.loadInternal(snap.mem.internal)

	c.nmiPending = snap.nmiPending
	c.executingNmi = snap.executingNmi
	c.nmiStep = snap.nmiStep

	c.currentOp = snap.currentOp
	c.currentstep = snap.currentstep
	c.TotalCycles = snap.totalCycles
	c.Stall = snap.Stall

	c.temp = snap.t
}
