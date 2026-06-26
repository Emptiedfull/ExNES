//go:build js && wasm

package main

import (
	"encoding/binary"
	"exnes/Core"
	"fmt"
	"math"
	"syscall/js"
)

var currentSpeed js.Value

func main() {

	stallChan := make(chan bool)

	startFrameDriver()

	<-stallChan
}

func startFrameDriver() {
	var emu *Core.Console
	var jsScreen js.Value

	const S_size = 2048
	var S_Arr = make([]byte, S_size*4)

	var JS_Arr js.Value

	var Input_Arr js.Value //also a dirty fucking shared buffer
	var Speed_Arr js.Value

	js.Global().Set("initBuffer", js.FuncOf(func(this js.Value, args []js.Value) any {
		JS_Arr = args[0]
		return nil
	}))

	js.Global().Set("initInput", js.FuncOf(func(this js.Value, args []js.Value) any {
		Input_Arr = args[0]
		return nil
	}))

	js.Global().Set("initSpeed", js.FuncOf(func(this js.Value, args []js.Value) any {
		Speed_Arr = args[0]
		return nil
	}))

	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		emu = Core.InitializeConsole()
		fmt.Println("console initialized")

		jsScreen = js.Global().Get("Uint8ClampedArray").New(256 * 240 * 4)
		js.Global().Set("frameBuffer", jsScreen)
		return "x"
	}))

	js.Global().Set("initRom", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		Arr := args[0]
		romData := make([]byte, Arr.Get("length").Int())
		js.CopyBytesToGo(romData, Arr)

		reader, _ := Core.LoadRomData(romData)

		err := emu.InitRom(reader)
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}

		emu.Cpu.Reset()
		fmt.Println("rom loaded")

		return nil
	}))

	js.Global().Set("drive", js.FuncOf(func(this js.Value, args []js.Value) any {

		speed := getSpeed(Speed_Arr) / 1000.0

		samplesNeeded := int(float32(args[0].Int()) * speed)
		syncInputs(emu, Input_Arr)

		for i := range samplesNeeded {

			for !emu.Apu.HasSample() {
				emu.Apu.Console.TickNoAudio()
			}
			sample := emu.Apu.PopSample()

			bits := math.Float32bits(sample)

			S_Arr[i*4] = byte(bits)
			S_Arr[i*4+1] = byte(bits >> 8)
			S_Arr[i*4+2] = byte(bits >> 16)
			S_Arr[i*4+3] = byte(bits >> 24)
		}

		js.CopyBytesToJS(jsScreen, emu.Ppu.BackBuffer[:])
		js.CopyBytesToJS(JS_Arr, S_Arr)

		return nil
	}))

}

func getSpeed(speedBuf js.Value) float32 {
	bytes := make([]byte, 4)
	js.CopyBytesToGo(bytes, speedBuf)
	return float32(binary.LittleEndian.Uint32(bytes))
}

func syncInputs(emu *Core.Console, input js.Value) {
	if emu == nil || input.IsUndefined() {
		return
	}

	mask := input.Index(0).Int()
	emu.Player1.UpdateBtnBool(Core.ButtonA, mask&(1<<0) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonB, mask&(1<<1) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonSelect, mask&(1<<2) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonStart, mask&(1<<3) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonUp, mask&(1<<4) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonDown, mask&(1<<5) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonLeft, mask&(1<<6) != 0)
	emu.Player1.UpdateBtnBool(Core.ButtonRight, mask&(1<<7) != 0)

}
