package Core

// Contains utility for the snapshot function

type Snapshot struct {
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

	Data []Snapshot

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

	Mem Cpu_Mem
	t   temp

	intKind interuptKind
	intStep int

	NmiLine bool
	IrqLine bool

	currentOp   uint8
	currentstep int
	totalCycles int
	Stall       int

	isJamming bool
}

type Cpu_Mem struct {
	internal [2048]byte
	external []byte
}

type PpuSnapshot struct {
	Mem ppu_Mem

	Dot      int
	Scanline int
	Frame    int

	mirroring       int
	mirroringChange bool

	// BackBuffer []uint8

	NewBuffer []uint32

	screenChanged bool
}

func (b *SnapshotBuffer) GetLength() int {
	accum := 0
	for _, snap := range b.Data {
		if snap.Frame_no != 0 {

			accum += 1

		}
	}

	return accum
}

func (c *Console) SetUpSnapshots() {
	for range len(c.Snapshots.Data) {

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

	c.Snapshots.Index = (c.Snapshots.Index + 1) % (len(c.Snapshots.Data) - 1)
}

func (c *Console) SaveState() *Snapshot {
	s := c.createEmptySnapshot()
	c.PopulateSnapshot(&s)
	return &s
}

func (c *Console) createEmptySnapshot() Snapshot {
	s := Snapshot{
		Frame_no:    0,
		Cycles:      0,
		mapperState: c.mapper.CreateEmptySnapshot(),
		CpuState:    CpuSnapshot{},
		PpuState:    PpuSnapshot{},
		ApuState:    ApuSnapshot{},
	}

	s.PpuState.NewBuffer = make([]uint32, 256*240)

	return s

}

func (p *ppu) TakePpuSnapshot(S *PpuSnapshot) {

	S.Mem = p.Mem

	S.Dot = p.Dot
	S.Scanline = p.Scanline
	S.Frame = p.Frame

	copy(S.NewBuffer, p.NewBuffer)

	S.screenChanged = p.ScreenChanged

}

func (b *SnapshotBuffer) AddSnapshot(shot Snapshot) {
	b.Data[b.Index] = shot
	b.Index = (b.Index + 1) % len(b.Data)

}

func (c Cpu) TakeCpuSnapshot(S *CpuSnapshot) {

	S.PC = c.PC
	S.S = c.S
	S.P = c.P
	S.X = c.X
	S.A = c.A
	S.Y = c.Y

	S.Mem = Cpu_Mem{
		internal: c.Mem.returnInternal(),
	}

	S.t = c.console.Cpu.temp

	S.intKind = c.intPresent
	S.intStep = c.intStep

	S.NmiLine = c.NmiLine
	S.IrqLine = c.irqLine

	S.currentOp = c.currentOp
	S.currentstep = c.currentstep
	S.totalCycles = c.TotalCycles

	S.Stall = c.Stall

	S.isJamming = c.isJamming

}

func (c *Console) PopulateSnapshot(s *Snapshot) {

	if c.mapper == nil || c.Cpu == nil || c.Ppu == nil || c.Apu == nil {
		return
	}

	c.Cpu.TakeCpuSnapshot(&s.CpuState)
	c.Ppu.TakePpuSnapshot(&s.PpuState)
	c.Apu.TakeSnapshot(&s.ApuState)
	c.mapper.TakeSnapshot(s.mapperState)
	s.Frame_no = c.Ppu.Frame
	s.Cycles = c.Cpu.TotalCycles
}

func (c *Console) AddSnapshot(s Snapshot) {
	s.Frame_no = c.Snapshots.Frame

	c.Snapshots.AddSnapshot(s)
}

func (c *Console) LoadSnapshot(snap Snapshot) {

	c.Cpu.LoadCpuSnapshot(snap.CpuState)
	c.Ppu.LoadPpuSnapshot(snap.PpuState)
	snap.mapperState.LoadSS(c.mapper)
}

func (p *ppu) LoadPpuSnapshot(snap PpuSnapshot) {
	p.Mem = snap.Mem

	p.Dot = snap.Dot
	p.Scanline = snap.Scanline
	p.Frame = snap.Frame

	copy(p.NewBuffer, snap.NewBuffer)

	p.ScreenChanged = snap.screenChanged
}

func (c *Cpu) LoadCpuSnapshot(snap CpuSnapshot) {
	c.PC = snap.PC
	c.S = snap.S
	c.P = snap.P
	c.X = snap.X
	c.A = snap.A
	c.Y = snap.Y

	c.Mem.loadInternal(snap.Mem.internal)

	c.intPresent = snap.intKind
	c.intStep = snap.intStep

	c.NmiLine = snap.NmiLine
	c.irqLine = snap.IrqLine

	c.currentOp = snap.currentOp
	c.currentstep = snap.currentstep
	c.TotalCycles = snap.totalCycles
	c.Stall = snap.Stall

	c.temp = snap.t
}
