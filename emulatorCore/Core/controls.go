package Core

import "fmt"

type joyPad struct {
	strobe uint8
	index  uint8

	current [8]uint8
	latched [8]uint8

	debugflag bool
}

const (
	ButtonA      = 0
	ButtonB      = 1
	ButtonSelect = 2
	ButtonStart  = 3
	ButtonUp     = 4
	ButtonDown   = 5
	ButtonLeft   = 6
	ButtonRight  = 7
)

func (J *joyPad) writeStrobe(val uint8) {
	J.strobe = val
	if val == 1 {
		J.index = 0
	} else {
		J.latched = J.current
	}
}

func (J *joyPad) readState() uint8 {

	if J.index > 7 {
		return 1
	}

	if J.strobe == 1 {
		return J.current[ButtonA]
	}

	var state uint8 = 0
	state = J.current[J.index]

	if J.debugflag {
		if state == 1 {
			fmt.Println("Reading control state:", J.current)
			J.debugflag = false
		}
	}

	J.index++

	return state

}

type ControlState struct {
	A      bool `json:"a"`
	B      bool `json:"b"`
	Select bool `json:"select"`
	Start  bool `json:"start"`
	Up     bool `json:"up"`
	Down   bool `json:"down"`
	Left   bool `json:"left"`
	Right  bool `json:"right"`
}

func (J *joyPad) UpdateState(c ControlState) {
	J.current[ButtonA] = uint8(convertBoolToInt(c.A))
	J.current[ButtonB] = uint8(convertBoolToInt(c.B))
	J.current[ButtonSelect] = uint8(convertBoolToInt(c.Select))
	J.current[ButtonStart] = uint8(convertBoolToInt(c.Start))
	J.current[ButtonUp] = uint8(convertBoolToInt(c.Up))
	J.current[ButtonDown] = uint8(convertBoolToInt(c.Down))
	J.current[ButtonLeft] = uint8(convertBoolToInt(c.Left))
	J.current[ButtonRight] = uint8(convertBoolToInt(c.Right))

	fmt.Println("updating state", J.current)
	J.debugflag = true

}

func convertBoolToInt(b bool) int {
	var val int = 0
	if b {
		val = 1
	}
	return val
}
