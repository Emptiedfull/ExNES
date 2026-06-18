package Core

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

	if t.lengthCounter == 0 {
		return 0
	}

	if t.linearCounter == 0 {
		return 0
	}

	return triangleSequence[t.sequence]
}
