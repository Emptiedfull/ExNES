package Core

type joyPad struct {
	strobe uint8
	index  uint8

	current [8]uint8
	latched [8]uint8

	debugflag bool
}

type BUTTON int

const (
	ButtonA BUTTON = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown

	ButtonLeft
	ButtonRight
)

func (b BUTTON) String() string {
	switch b {
	case ButtonA:
		return "Joypad-A"
	case ButtonB:
		return "Joypad-B"
	case ButtonSelect:
		return "Select"
	case ButtonStart:
		return "Start"
	case ButtonUp:
		return "Dpad-Up"
	case ButtonDown:
		return "Dpad-Down"
	case ButtonLeft:
		return "Dpad-Left"
	case ButtonRight:
		return "Dpad-Right"
	default:
		return "yeh bro pack it up"
	}
}

func GetActionByName(name string) BUTTON {

	switch name {
	case "Dpad-Right":
		return ButtonRight
	case "Dpad-Left":
		return ButtonLeft
	case "Dpad-Down":
		return ButtonDown
	case "Dpad-Up":
		return ButtonUp
	case "Start":
		return ButtonStart
	case "Select":
		return ButtonSelect
	case "Joypad-A":
		return ButtonA
	case "Joypad-B":
		return ButtonB
	}

	return 0
}

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
		return 0x40 | 1
	}

	if J.strobe == 1 {
		J.index = 0
		return 0x40 | J.current[ButtonA]
	}

	state := J.latched[J.index]

	J.index++

	return 0x40 | state

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

}

func (J *joyPad) UpdateBtnState(btn int, state uint8) {
	J.current[btn] = state
}

func (J *joyPad) UpdateBtnBool(btn BUTTON, state bool) {
	J.current[btn] = uint8(convertBoolToInt(state))
}
