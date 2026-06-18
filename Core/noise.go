package Core

//none of the code or comments here are ai istg dont ban for ai

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
