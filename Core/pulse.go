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
	//update mute (not ai dont kill me )

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
			return p.timer - x
		}
	} else {
		return p.timer + x
	}
}

func (p *PulseChannel) updateMuting() {
	target := p.GetTargetPeriod()

	if p.timerPeriod < 8 || target > 0x07FF {
		p.muting = true
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
