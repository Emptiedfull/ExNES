//go:build js && wasm

package main

import (
	"exnes/Core"
	"fmt"
	"math"
	"syscall/js"
	"unsafe"
)

var currentSpeed js.Value

func main() {

	stallChan := make(chan bool)

	fmt.Println("core version 20")

	startFrameDriver()

	<-stallChan
}

func startFrameDriver() {
	var emu *Core.Console
	var jsScreen js.Value

	const S_size = 2048
	var S_Arr = make([]byte, S_size*4)

	var JS_Arr js.Value

	js.Global().Set("initBuffer", js.FuncOf(func(this js.Value, args []js.Value) any {
		JS_Arr = args[0]
		return nil
	}))

	js.Global().Set("getSamples", js.FuncOf(func(this js.Value, args []js.Value) any {

		// size := int32(len(sampleBuf))
		if emu == nil {
			return js.ValueOf(0)
		}

		samples := emu.Apu.DrainSamples()
		if len(samples) == 0 {
			return js.ValueOf(0)
		}

		count := len(samples)
		need := count * 4

		view := unsafe.Slice(
			(*byte)(unsafe.Pointer(&samples[0])), need,
		)

		copy(S_Arr, view)

		js.CopyBytesToJS(JS_Arr, S_Arr[:need])

		return js.ValueOf(count)
	}))

	fmt.Println("core initialized 2")
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

	js.Global().Set("drive", js.FuncOf(func(this js.Value, args []js.Value) any {

		samplesNeeded := args[0].Int()

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

		fmt.Println(emu.Ppu.Frame)

		return nil
	}))

	loop := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		// emu.RunFrame()

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
		samplesNeeded := args[0].Int()
		jsBuf := args[1]

		for i := range samplesNeeded {
			var sample float32
			if emu.Apu.HasSample() {
				sample = emu.Apu.PopSample()
			} else {
				sample = 1
			}

			bits := math.Float32bits(sample)
			outputBuf[i*4] = byte(bits)
			outputBuf[i*4+1] = byte(bits >> 8)

			outputBuf[i*4+2] = byte(bits >> 16)
			outputBuf[i*4+3] = byte(bits >> 24)
		}
		js.CopyBytesToJS(jsBuf, outputBuf)
		return nil
	}))

	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) any {
		emu = Core.InitializeConsole()

		outputScreen = js.Global().Get("Uint8ClampedArray").New(256 * 240 * 4)
		js.Global().Set("frameBuffer", outputScreen)
		return nil
	}))

	js.Global().Set("update", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		emu.Player1.UpdateBtnBool(args[0].Int(), args[1].Bool())

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
