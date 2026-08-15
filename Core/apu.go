package Core

import "math"

type APU struct {
	Console *Console

	Pulse1 *PulseChannel
	Pulse2 *PulseChannel

	Triangle *Triangle

	Noise *NoiseChannel

	Dmc *DMC

	Counter frameCounter

	SampleRate      float64
	CycleAcc        float64
	CyclesPerSample float64

	timerTickLatch bool

	IRGPending bool

	SampleBuffer []float32
}

func (a *APU) HasSample() bool {
	return len(a.SampleBuffer) > 0
}

func (a *APU) PopSample() float32 {
	s := a.SampleBuffer[0]
	a.SampleBuffer = a.SampleBuffer[1:]
	return s
}

func NewApu(sampleRate float64, console *Console) *APU {
	return &APU{
		Pulse1:          newPulseChannel(true),
		Pulse2:          newPulseChannel(false),
		Console:         console,
		Triangle:        newTriangle(),
		Noise:           NewNoiseChannel(),
		Dmc:             NewDMC(),
		SampleRate:      sampleRate,
		CyclesPerSample: 1_789_773.0 / sampleRate,
		SampleBuffer:    make([]float32, 0, 4096),
	}

}

var pulseTable [31]float32
var tndTable [203]float32

func init() {
	for i := range 31 {
		if i == 0 {
			continue
		}
		pulseTable[i] = 95.52 / (8128.0/float32(i) + 100.0) //magic numbers wooo
	}

	for i := 1; i < 203; i++ {
		if i == 0 {
			continue
		}
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
	case 0x4008:
		A.Triangle.WriteLC(val)
	case 0x4009:
		// yeh fuck u cpu for trying to access this
	case 0x400A:
		A.Triangle.WriteTimerLOW(val)
	case 0x400B:
		A.Triangle.WriteTimerHIGH(val)
	case 0x4015:
		A.writeStatus(val)
	case 0x4017:
		A.writeFrameCounter(val)
	case 0x400C:
		A.Noise.WriteEnvelope(val)
	case 0x400D:
		// YEH FUCK U BRO DONT ACCESS THIS
	case 0x4010:
		A.Dmc.WriteFlags(val)
	case 0x4011:
		A.Dmc.WriteDL(val)
	case 0x4012:
		A.Dmc.WriteSampleAddr(val)
	case 0x4013:
		A.Dmc.WriteSampleLen(val)
	case 0x400E:
		A.Noise.WriteMode(val)
	case 0x400F:
		A.Noise.WriteLC(val)
	}
}

func (a *APU) writeStatus(val uint8) {
	a.Pulse1.setEnable(val&0x01 != 0)
	a.Pulse2.setEnable(val&0x02 != 0)
	a.Triangle.setEnable(val&0x04 != 0)
	a.Noise.setEnabled(val&0x08 != 0)
	a.Dmc.setEnable(val&0x10 != 0)

	a.Dmc.IRGPending = false

}

func (A *APU) readStatus() uint8 {
	var status uint8

	if A.Pulse1.lengthCounter > 0 {
		status |= 0x01
	}

	if A.Pulse2.lengthCounter > 0 {
		status |= 0x02
	}

	if A.Triangle.lengthCounter > 0 {
		status |= 0x04
	}

	if A.Noise.lengthCounter > 0 {
		status |= 0x08
	}

	if A.IRGPending {
		status |= 0x40
	}

	if A.Dmc.excessBytes > 0 {
		status |= 0x10
	}

	if A.Dmc.IRGPending {
		status |= 0x80
	}

	A.IRGPending = false
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

func (a *APU) tick() {
	a.Counter.cycles++
	a.stepFrameCounter()

	a.timerTickLatch = !a.timerTickLatch
	if a.timerTickLatch {
		a.Pulse1.StepTimer()
		a.Pulse2.StepTimer()
		a.Noise.stepTimer()
		a.Dmc.stepTimer()
	}

	a.Triangle.StepTimer()

	if a.Dmc.stall > 0 {
		val := a.Console.Cpu.Mem.Read(a.Dmc.currentAddr)
		a.Dmc.LoadSample(val)

		a.Console.Cpu.Stall += 4
	}

	a.CycleAcc++
	if a.CycleAcc >= a.CyclesPerSample {
		a.CycleAcc -= a.CyclesPerSample

		a.SampleBuffer = append(a.SampleBuffer, a.mix())

	}
}

func (a *APU) stepFrameCounter() {
	c := a.Counter.cycles

	if a.Counter.mode == 0 {
		switch c {
		case 7457:
			a.clockQuarterFrame()
		case 14913:
			a.clockQuarterFrame()
			a.clockHalfFrame()
		case 22371:
			a.clockQuarterFrame()
		case 29828:
			if !a.Counter.irqStop {
				a.IRGPending = true
			}
		case 29829:
			a.clockQuarterFrame()
			a.clockHalfFrame()
			if !a.Counter.irqStop {
				a.IRGPending = true
			}
		case 29830:
			a.Counter.cycles = 0
			if !a.Counter.irqStop {
				a.IRGPending = true
			}
		}
	} else {
		switch c {
		case 7457:
			a.clockQuarterFrame()
		case 14913:
			a.clockQuarterFrame()
			a.clockHalfFrame()
		case 22371:
			a.clockQuarterFrame()
		case 37281:
			a.clockQuarterFrame()
			a.clockHalfFrame()
		case 37282:
			a.Counter.cycles = 0
		}
	}
}

func (a *APU) clockQuarterFrame() {
	a.Pulse1.StepEnvelope()
	a.Pulse2.StepEnvelope()

	a.Triangle.stepLinC()

	a.Noise.StepEnvelope()

}

func (a *APU) clockHalfFrame() {
	a.Pulse1.stepLengthCounter()
	a.Pulse1.StepSweep()
	a.Pulse2.stepLengthCounter()
	a.Pulse2.StepSweep()

	a.Triangle.stepLenC()

	a.Noise.StepLC()
}

func (a *APU) mix() float32 {
	p1 := a.Pulse1.Output()
	p2 := a.Pulse2.Output()

	tri := a.Triangle.Output()
	noi := a.Noise.Output()

	dmc := a.Dmc.Output()

	pulseOut := pulseTable[p1+p2]
	tndOut := tndTable[3*uint(tri)+2*uint(noi)+uint(dmc)]

	return pulseOut + tndOut
}

func (a *APU) DrainSamples() []float32 {
	out := make([]float32, len(a.SampleBuffer))
	copy(out, a.SampleBuffer)

	a.SampleBuffer = a.SampleBuffer[:0]
	return out
}

func (a *APU) MalgoAdapter(output []byte, input []byte, framecount uint32) {
	a.DriveSamples(output, framecount)
}

func (a *APU) DriveSamples(output []byte, samplesNeeded uint32) {
	for i := range samplesNeeded {
		for !a.HasSample() {
			a.Console.Step()
		}

		sample := a.PopSample()

		bits := math.Float32bits(sample)
		output[i*4] = byte(bits)
		output[i*4+1] = byte(bits >> 8)
		output[i*4+2] = byte(bits >> 16)
		output[i*4+3] = byte(bits >> 24)
	}

	if a.Console.Ppu.ScreenChanged {
		a.Console.RunDisplayUpdates()

		if a.Console.Ppu.Frame%20 == 0 {

			a.Console.TakeSnapshot()
		}

	}

}
