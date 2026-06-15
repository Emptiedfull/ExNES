//go:build js && wasm

package main

import (
	"exnes/Core"
	"fmt"
	"syscall/js"
)

var emu *Core.Console
var jsScreen js.Value

func main() {

	stallChan := make(chan bool)

	fmt.Println("core initialized 51")
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
		}

		emu.Cpu.Reset()
		fmt.Println("console ready")

		return nil
	}))

	var framesRun int

	loop := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		emu.RunFrame()
		js.CopyBytesToJS(jsScreen, emu.Ppu.BackBuffer[:])
		framesRun++

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

func initCanvas() {
	canvas := js.Global().Get("document").Call("getElementById", "screen")
	ctx := canvas.Call("getContext", "2d")

	imageData := ctx.Call("createImageData", 256, 240)
	jsScreen = imageData.Get("Data")
}

func updateDisplay(this js.Value, args []js.Value) interface{} {
	fmt.Println("somethign bad happened")

	return nil
}
