package Core

var lengthTable = [32]byte{
	10, 254, 20, 2, 40, 4, 80, 6,
	160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22,
	192, 24, 72, 26, 16, 28, 32, 30,
}

var dutyTable = [4][8]uint8{
	{0, 1, 0, 0, 0, 0, 0, 0}, //12.5
	{0, 1, 1, 0, 0, 0, 0, 0}, //25
	{0, 1, 1, 1, 1, 0, 0, 0}, // 50
	{1, 0, 0, 1, 1, 1, 1, 1}, //75
}

type PulseChannel struct {
	enabled bool

	timer       uint16
	timerPeriod uint16

	dutyMode byte
	dutyStep byte

	lengthCounter     byte
	lengthCounterHalt bool

	Sweep
	Envelope

	Pulse1 bool
}

func newPulseChannel(isPulse1 bool) *PulseChannel {
	return &PulseChannel{
		Pulse1: isPulse1,
	}
}

func (p *PulseChannel) setEnable(val bool) {
	p.enabled = val
	if !p.enabled {
		p.lengthCounter = 0
	}
}

type Sweep struct {
	enabled bool
	period  byte
	negate  bool
	shift   byte
	reload  bool
	divider byte
	muting  bool
}

type Envelope struct {
	start    bool
	loop     bool
	constant bool
	volume   byte
	decayVol byte
	divider  byte
}

func (p *PulseChannel) writeCtrl(val uint8) { // DD L C VVVV
	p.dutyMode = (val >> 6) & 0x03
	p.lengthCounterHalt = val&0x20 != 0
	p.Envelope.loop = p.lengthCounterHalt

	p.constant = val&0x10 != 0
	p.volume = val & 0x0F
}

func (p *PulseChannel) writeSweep(val uint8) { // E PPP N SSS
	p.Sweep.enabled = val&0x80 != 0
	p.Sweep.period = (val >> 4) & 0x07
	p.Sweep.negate = val&0x08 != 0
	p.Sweep.shift = val & 0x07
	p.Sweep.reload = true
}

func (p *PulseChannel) writeTimerLow(val uint8) { // LLLL LLLL
	p.timerPeriod = (p.timerPeriod & 0x0700) | uint16(val)
}

func (p *PulseChannel) writeTimerHigh(val uint8) { //LLLL LHHH

	p.timerPeriod = (p.timerPeriod & 0x00FF) | (uint16(val&0x07) << 8)
	if p.enabled {
		p.lengthCounter = lengthTable[val>>3]
	}

	p.Envelope.start = true
	p.dutyStep = 0

}

func (p *PulseChannel) StepTimer() {
	if p.timer == 0 {
		p.timer = p.timerPeriod
		p.dutyStep = (p.dutyStep + 1) % 8
	} else {
		p.timer--
	}
}

func (p *PulseChannel) StepEnvelope() {
	if p.Envelope.start {
		p.Envelope.start = false
		p.Envelope.decayVol = 15
		p.Envelope.divider = p.Envelope.volume
		return
	}

	if p.Envelope.divider == 0 {
		p.Envelope.divider = p.Envelope.volume
		if p.Envelope.decayVol == 0 {
			if p.Envelope.loop {
				p.Envelope.decayVol = 15
			}
		} else {
			p.Envelope.decayVol--
		}
	} else {
		p.Envelope.divider--
	}
}

func (p *PulseChannel) stepLengthCounter() {
	if !p.lengthCounterHalt && p.lengthCounter > 0 {
		p.lengthCounter--
	}
}

func (p *PulseChannel) StepSweep() {
	p.updateMuting()

	if p.Sweep.divider == 0 && p.Sweep.enabled && p.Sweep.shift != 0 && !p.Sweep.muting {
		p.timerPeriod = p.GetTargetPeriod()
	}

	if p.Sweep.divider == 0 || p.Sweep.reload {
		p.Sweep.divider = p.Sweep.period
		p.Sweep.reload = false
	} else {
		p.Sweep.divider--
	}
}

func (p *PulseChannel) GetTargetPeriod() uint16 {
	x := p.timerPeriod >> uint16(p.Sweep.shift)
	if p.Sweep.negate {
		if p.Pulse1 {
			return p.timerPeriod - x - 1
		} else {
			return p.timerPeriod - x
		}
	} else {
		return p.timerPeriod + x
	}
}

func (p *PulseChannel) updateMuting() {
	target := p.GetTargetPeriod()

	if p.timerPeriod < 8 || target > 0x07FF {
		p.muting = true
	} else {
		p.muting = false
	}
}

func (p *PulseChannel) Output() uint8 {
	//if you are contributing pls help me fix this its so bad

	if !p.enabled || p.lengthCounter == 0 {
		return 0
	}

	if p.timerPeriod < 8 { //fucking garbage dog level hearing
		return 0
	}

	if p.Sweep.muting {
		return 0
	}

	if dutyTable[p.dutyMode][p.dutyStep] == 0 {
		return 0
	}

	if p.Envelope.constant {
		return p.Envelope.volume
	}
	return p.Envelope.decayVol
}

var noisePeriodTable = [16]uint16{
	4, 8, 16, 32, 64, 96, 128, 160,
	202, 254, 380, 508, 762, 1016, 2034, 4068,
}

type NoiseChannel struct {
	enabled bool

	timer       uint16
	timerPeriod uint16

	shiftReg uint16

	mode bool // false = long, true = short (also not ai pls no kill)

	lengthCounter     uint8
	lengthCounterHalt bool

	envelope Envelope
}

func NewNoiseChannel() *NoiseChannel {
	return &NoiseChannel{
		shiftReg: 1,
	}
}

func (n *NoiseChannel) setEnabled(enabled bool) {
	n.enabled = enabled
	if !n.enabled {
		n.lengthCounter = 0
	}
}

func (n *NoiseChannel) WriteEnvelope(val uint8) { // --LC VVVV
	n.lengthCounterHalt = val&0x20 != 0
	n.envelope.loop = n.lengthCounterHalt

	n.envelope.constant = val&0x10 != 0
	n.envelope.volume = val & 0x0F
}

func (n *NoiseChannel) WriteMode(val uint8) { // M--- PPPP
	n.mode = val&0x80 != 0
	n.timerPeriod = noisePeriodTable[val&0x0F]
}

func (n *NoiseChannel) WriteLC(val byte) { // llll l---
	if n.enabled {
		n.lengthCounter = lengthTable[val>>3]
	}
	n.envelope.start = true
}

func (n *NoiseChannel) stepTimer() {
	if n.timer == 0 {
		n.timer = n.timerPeriod
		n.clockFLSR()
	} else {
		n.timer--
	}
}

func (n *NoiseChannel) clockFLSR() {
	var feed uint16

	if n.mode {
		feed = (n.shiftReg & 0x0001) ^ ((n.shiftReg >> 6) & 0x0001)
	} else {
		feed = (n.shiftReg & 0x0001) ^ ((n.shiftReg >> 1) & 0x0001)
	}

	n.shiftReg >>= 1
	n.shiftReg |= feed << 14
}

func (n *NoiseChannel) StepEnvelope() {
	if n.envelope.start {
		n.envelope.start = false
		n.envelope.decayVol = 15
		n.envelope.divider = n.envelope.volume
		return
	}

	if n.envelope.divider == 0 {
		n.envelope.divider = n.envelope.volume
		if n.envelope.decayVol == 0 {
			if n.envelope.loop {
				n.envelope.decayVol = 15
			}

		} else {
			n.envelope.decayVol--
		}

	} else {
		n.envelope.divider--
	}
}

func (n *NoiseChannel) StepLC() {
	if !n.lengthCounterHalt && n.lengthCounter > 0 {
		n.lengthCounter--
	}
}

func (n *NoiseChannel) Output() uint8 {
	if !n.enabled {
		return 0
	}

	if n.lengthCounter == 0 {
		return 0
	}

	if n.shiftReg&0x0001 == 1 {
		return 0
	}

	if n.envelope.constant {
		return n.envelope.volume
	}

	return n.envelope.decayVol
}

type Triangle struct {
	enabled bool

	timer       uint16
	timerPeriod uint16

	sequence uint8

	lengthCounter     uint8
	lengthCounterHalt bool

	linearCounter       uint8
	linearCounterReload uint8
	linearCounterFlag   bool
}

func newTriangle() *Triangle {
	return &Triangle{}
}

func (t *Triangle) setEnable(enabled bool) {
	t.enabled = enabled
	if !t.enabled {
		t.lengthCounter = 0
	}
}

func (t *Triangle) WriteLC(val uint8) { //CRRR RRRR
	t.lengthCounterHalt = val&0x80 != 0
	t.linearCounterReload = val & 0x7F
}

func (t *Triangle) WriteTimerLOW(val uint8) { // TTTT TTTT
	t.timerPeriod = (t.timerPeriod & 0x0700) | uint16(val)
}

func (t *Triangle) WriteTimerHIGH(val uint8) { // LLLL LHHH
	t.timerPeriod = (t.timerPeriod & 0x00FF) | (uint16(val&0x07) << 8)
	if t.enabled {
		t.lengthCounter = lengthTable[val>>3]
	}
	t.linearCounterFlag = true
}

func (t *Triangle) StepTimer() {
	if t.timer == 0 {
		t.timer = t.timerPeriod
		if t.lengthCounter > 0 && t.linearCounter > 0 {
			t.sequence = (t.sequence + 1) % 32
		}
	} else {
		t.timer--
	}
}

func (t *Triangle) stepLinC() {
	if t.linearCounterFlag {
		t.linearCounter = t.linearCounterReload
	} else if t.linearCounter > 0 {
		t.linearCounter--
	}

	if !t.lengthCounterHalt {
		t.linearCounterFlag = false
	}
}

func (t *Triangle) stepLenC() {
	if !t.lengthCounterHalt && t.lengthCounter > 0 {
		t.lengthCounter--
	}
}

var triangleSequence = [32]uint8{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

func (t *Triangle) Output() uint8 {
	if !t.enabled {
		return 0
	}

	// if t.lengthCounter == 0 {
	// 	return 0
	// }

	// if t.linearCounter == 0 {
	// 	return 0
	// }

	return triangleSequence[t.sequence]
}

type DMC struct {
	enabled    bool
	irqEnabled bool
	loop       bool
	IRGPending bool

	timer       uint16
	timerPeriod uint16

	output uint8

	sampleAddr uint16
	sampleLen  uint16

	currentAddr uint16
	excessBytes uint16

	sampleBuffer uint8
	BufferFull   bool

	shiftReg   uint8
	excessBits uint8
	mute       bool

	stall int
}

func NewDMC() *DMC {
	return &DMC{}
}

var dmcRateTable = [16]uint16{
	428, 380, 340, 320, 286, 254, 226, 214,
	190, 160, 142, 128, 106, 84, 72, 54,
}

func (d *DMC) WriteFlags(val uint8) { // IL-- RRRR
	d.irqEnabled = val&0x80 != 0
	d.loop = val&0x40 != 0
	d.timerPeriod = dmcRateTable[val&0x0F]

	if !d.irqEnabled {
		d.IRGPending = false
	}
}

func (d *DMC) WriteDL(val uint8) { // -DDD DDDD
	d.output = val & 0x7F
}

func (d *DMC) WriteSampleAddr(val uint8) {
	d.sampleAddr = 0xC000 + uint16(val)*64
}

func (d *DMC) WriteSampleLen(val uint8) {
	d.sampleLen = uint16(val)*16 + 1
}

func (d *DMC) setEnable(val bool) {
	d.enabled = val
	if !val {
		d.excessBytes = 0
		d.IRGPending = false
	} else if d.excessBytes == 0 {
		d.currentAddr = d.sampleAddr
		d.excessBytes = d.sampleLen
	}
}

func (d *DMC) LoadSample(val uint8) {
	d.sampleBuffer = val
	d.BufferFull = true
	d.stall = 0

	d.currentAddr++
	if d.currentAddr == 0 {
		d.currentAddr = 0x8000
	}

	d.excessBytes--
	if d.excessBytes == 0 {
		if d.loop {
			d.currentAddr = d.sampleAddr
			d.excessBytes = d.sampleLen
		} else if d.irqEnabled {
			d.IRGPending = true
		}
	}
}

func (d *DMC) stepTimer() {
	if d.timer == 0 {
		d.timer = d.timerPeriod
		d.stepOutput()
	} else {
		d.timer--
	}
}

func (d *DMC) stepOutput() {
	if !d.mute {
		if d.shiftReg&0x01 == 1 {
			if d.output <= 125 {
				d.output += 2
			}
		} else {
			if d.output >= 2 {
				d.output -= 2
			}
		}
	}

	d.shiftReg >>= 1

	if d.excessBits > 0 {
		d.excessBits--
	}

	if d.excessBits == 0 {
		d.excessBits = 8
		if d.BufferFull {
			d.mute = false
			d.shiftReg = d.sampleBuffer
			d.BufferFull = false

			if d.excessBytes > 0 {
				d.stall = 4
			}
		} else {
			d.mute = true
		}
	}
}

func (d *DMC) Output() uint8 {
	return d.output
}
