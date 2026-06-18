package Core

type APU struct {
	Pulse1 PulseChannel
	Pulse2 PulseChannel

	Counter frameCounter

	SampleRate      float64
	CycleAcc        float64
	CyclesPerSample float64

	IRGPending bool

	sampleBuffer []float32
}

func newApu(sampleRate float64) *APU {
	return &APU{
		Pulse1:          *newPulseChannel(true),
		Pulse2:          *newPulseChannel(false),
		SampleRate:      sampleRate,
		CyclesPerSample: 1_789_773.0 / sampleRate,
		sampleBuffer:    make([]float32, 0, 4096),
	}

}

var pulseTable [31]float32
var tndTable [203]float32

func init() {
	for i := range 31 {
		pulseTable[i] = 95.52 / (8128.0/float32(i) + 100.0) //magic numbers wooo
	}

	for i := 1; i < 203; i++ {
		tndTable[i] = 163.67 / (24329.0/float32(i) + 100.0)
	}
}

type frameCounter struct {
	mode    uint8
	irqStop bool
	cycles  int
}

type cycleInfo struct {
	cycle        int
	quarterFrame bool
	halfFrame    bool
	irq          bool
}

var FourSteps = [5]cycleInfo{
	{7457, true, false, false},
	{14913, true, true, false},
	{22371, true, false, false},
	{29828, false, false, true},
	{29829, true, true, true},
}

var FiveSteps = [4]cycleInfo{
	{7457, true, false, false},
	{14913, true, true, false},
	{22371, true, false, false},
	{37281, true, true, false},
}

func (A *APU) writeReg(addr uint16, val uint8) {
	switch addr {
	case 0x4000:
		A.Pulse1.writeCtrl(val)
	case 0x4001:
		A.Pulse1.writeSweep(val)
	case 0x4002:
		A.Pulse1.writeTimerLow(val)
	case 0x4003:
		A.Pulse1.writeTimerHigh(val)
	case 0x4004:
		A.Pulse2.writeCtrl(val)
	case 0x4005:
		A.Pulse2.writeSweep(val)
	case 0x4006:
		A.Pulse2.writeTimerLow(val)
	case 0x4007:
		A.Pulse2.writeTimerHigh(val)
	}
}

func (A *APU) readStatus() uint8 {
	var status uint8

	if A.Pulse1.lengthCounter > 0 {
		status |= 0x01
	}

	if A.Pulse2.lengthCounter > 0 {
		status |= 0x02
	}

	return status
}

func (a *APU) writeFrameCounter(val uint8) {
	a.Counter.mode = (val >> 7) & 0x01

	a.Counter.irqStop = val&0x40 != 0
	if a.Counter.irqStop {
		a.IRGPending = false
	}

	if a.Counter.mode == 1 { //5 step
		a.clockHalfFrame()
		a.clockQuarterFrame()
	}

	a.Counter.cycles = 0
}

func (a *APU) clockQuarterFrame() {
	a.Pulse1.StepEnvelope()
	a.Pulse2.StepEnvelope()
}

func (a *APU) clockHalfFrame() {
	a.Pulse1.stepLengthCounter()
	a.Pulse1.StepSweep()
	a.Pulse2.stepLengthCounter()
	a.Pulse2.StepSweep()
}
