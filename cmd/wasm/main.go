//go:build js && wasm

package main

import (
	"exnes/Core"
	"fmt"
	"syscall/js"
)

var emu *Core.Console
var jsScreen js.Value
var currentSpeed js.Value

func main() {

	stallChan := make(chan bool)

	fmt.Println("core initialized 1")
	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		emu = Core.InitializeConsole()
		fmt.Println("console initialized")

		jsScreen = js.Global().Get("Uint8ClampedArray").New(256 * 240 * 4)
		js.Global().Set("frameBuffer", jsScreen)
		js.Global().Set("currentSpeed", currentSpeed)
		return "x"
	}))

	js.Global().Set("update", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		emu.Player1.UpdateBtnBool(args[0].Int(), args[1].Bool())
		emu.Player2.UpdateBtnBool(args[0].Int(), args[1].Bool())

		return nil
	}))

	js.Global().Set("pause", js.FuncOf(func(this js.Value, args []js.Value) any {
		if emu != nil {
			emu.Pause()
		}
		return nil
	}))

	js.Global().Set("unpause", js.FuncOf(func(this js.Value, args []js.Value) any {
		if emu != nil {
			emu.UnPause()
		}
		return nil
	}))

	js.Global().Set("takeSnapshot", js.FuncOf(func(this js.Value, args []js.Value) any {
		emu.AddSnapshot()
		return nil
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
		fmt.Println("console ready")

		return nil
	}))

	js.Global().Set("reset", js.FuncOf(func(this js.Value, args []js.Value) any {
		return nil
	}))

	currentSpeed = js.ValueOf(1)

	loop := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		for range 2 {
			emu.RunFrame()
		}
		js.CopyBytesToJS(jsScreen, emu.Ppu.BackBuffer[:])

		return nil
	})

	js.Global().Set("nesFrame", loop)

	// loop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	emu.RunFrame()
	// 	framesRun++
	// 	// fmt.Println("frame")

	// 	if framesRun%60 == 0 {
	// 		now := js.Global().Get("performance").Call("now").Float()
	// 		elapsed := now - lasttime
	// 		fmt.Printf("FPS: %.1f\n", 60000/elapsed)
	// 		lasttime = now

	// 	}

	// 	js.Global().Call("requestAnimationFrame", loop)
	// 	return nil
	// })

	<-stallChan
}
