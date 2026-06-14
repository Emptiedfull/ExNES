package Core

// Contains utility for the snapshot function

type snapshot struct {
	Frame_no int
	Cycles   int
	CpuState CpuSnapshot
	PpuState PpuSnapshot
}

type SnapshotBuffer struct {
	Frame int

	Data  [200]snapshot
	Index int
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

	backBuffer  []uint8
	frontBuffer []uint8

	screenChanged bool
}

func (p *ppu) TakePpuSnapshot() PpuSnapshot {
	S := PpuSnapshot{}

	S.mem = p.mem

	orignal_chrRom := p.mem.mapper.extractCHR()

	S.mem.chrRom_WARNING = make([]uint8, len(orignal_chrRom))
	copy(S.mem.chrRom_WARNING, orignal_chrRom)

	S.Dot = p.Dot
	S.Scanline = p.Scanline
	S.Frame = p.Frame

	copy(S.backBuffer, p.backBuffer)
	copy(S.frontBuffer, p.frontBuffer)

	S.screenChanged = p.screenChanged

	return S
}

func (b *SnapshotBuffer) AddSnapshot(shot snapshot) {
	b.Data[b.Index] = shot
	b.Index = (b.Index + 1) % len(b.Data)

}

func (c cpu) TakeCpuSnapshot() CpuSnapshot {
	S := CpuSnapshot{}

	S.PC = c.PC
	S.S = c.S
	S.P = c.P
	S.X = c.X
	S.Y = c.Y

	S.mem = cpu_mem{
		internal: c.mem.returnInternal(),
		external: c.mem.returnExternal(),
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

	return S
}

func (d *Debugger) TakeSnapshot() snapshot {
	S := snapshot{
		CpuState: d.Console.Cpu.TakeCpuSnapshot(),
		PpuState: d.Console.Ppu.TakePpuSnapshot(),
		Cycles:   d.Console.Cpu.TotalCycles,
	}

	return S
}

func (d *Debugger) AddSnapshot() {
	s := d.TakeSnapshot()
	s.Frame_no = d.RecentHistory.Frame
	d.RecentHistory.Frame++
	d.RecentHistory.AddSnapshot(s)
}

func (d *Debugger) LoadSnapshot(snap snapshot) {
	d.Console.Pause()
	d.Console.Cpu.LoadCpuSnapshot(snap.CpuState)
	d.Console.Ppu.LoadPpuSnapshot(snap.PpuState)
}

func (p *ppu) LoadPpuSnapshot(snap PpuSnapshot) {
	p.mem = snap.mem

	p.mem.mapper.loadCHR(snap.mem.chrRom_WARNING)

	p.Dot = snap.Dot
	p.Scanline = snap.Scanline
	p.Frame = snap.Frame

	copy(p.backBuffer, snap.backBuffer)
	copy(p.frontBuffer, snap.frontBuffer)

	p.screenChanged = snap.screenChanged
}

func (c *cpu) LoadCpuSnapshot(snap CpuSnapshot) {
	c.PC = snap.PC
	c.S = snap.S
	c.P = snap.P
	c.X = snap.X
	c.Y = snap.Y

	c.mem.loadExternal(snap.mem.external)
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
