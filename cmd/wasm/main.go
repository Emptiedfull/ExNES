//go:build js && wasm

package main

import (
	"exnes/Core"
	"fmt"
	"syscall/js"
)

var currentSpeed js.Value

func main() {

	stallChan := make(chan bool)

	fmt.Println("core version 2")

	startAudioDriver()

	<-stallChan
}

func startFrameDriver() {
	var emu *Core.Console
	var jsScreen js.Value

	fmt.Println("core initialized 1")
	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		emu = Core.InitializeConsole()
		fmt.Println("console initialized")

		jsScreen = js.Global().Get("Uint8ClampedArray").New(256 * 240 * 4)
		js.Global().Set("frameBuffer", jsScreen)
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

	loop := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		emu.RunFrame()

		js.CopyBytesToJS(jsScreen, emu.Ppu.BackBuffer[:])

		return nil
	})

	js.Global().Set("nesFrame", loop)

}

func startAudioDriver() {

	var outputBuf []byte
	var emu *Core.Console
	var outputScreen js.Value

	js.Global().Set("initBuffer", js.FuncOf(func(this js.Value, args []js.Value) any {
		count := args[0].Int()
		outputBuf = make([]byte, count*4)

		return nil
	}))

	js.Global().Set("getSamples", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("tyo mna")
		samplesNeeded := args[0].Int()

		emu.Apu.DriveSamples(outputBuf, uint32(samplesNeeded))

		jsBuf := args[1]
		js.CopyBytesToJS(jsBuf, outputBuf)

		return nil
	}))

	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) any {
		emu = Core.InitializeConsole()

		outputScreen = js.Global().Get("Uint8ClampedArray").New(256 * 240 * 4)
		js.Global().Set("frameBuffer", outputScreen)
		return nil
	}))

	js.Global().Set("initRom", js.FuncOf(func(this js.Value, args []js.Value) any {

		romdat := args[0]
		romArr := make([]byte, romdat.Get("length").Int())
		js.CopyBytesToGo(romArr, romdat)

		reader, _ := Core.LoadRomData(romArr)

		err := emu.InitRom(reader)
		if err != nil {
			fmt.Println("Ah fuck here we go again", err)
		}

		emu.Cpu.Reset()
		fmt.Println("loaded")

		return nil
	}))

	loop := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.CopyBytesToJS(outputScreen, emu.Ppu.BackBuffer[:])
		return nil
	})

	js.Global().Set("nesFrame", loop)
}
