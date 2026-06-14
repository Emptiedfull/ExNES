//go:build js && wasm

package main

import (
	"exnes/Core"
	"fmt"
	"syscall/js"
	"time"
)

var emu *Core.Console

func main() {

	stallChan := make(chan bool)

	fmt.Println("core initialized 4")
	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		emu = Core.InitializeConsole()
		fmt.Println("console initialized")

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
	var loop js.Func

	js.Global().Set("runFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		start := time.Now()
		for range 60 {
			emu.RunFrame()
		}

		duration := time.Since(start)

		return duration.Seconds()

	}))

	loop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		emu.RunFrame()

		js.Global().Call("requestAnimationFrame", loop)
		return nil
	})

	// js.Global().Call("requestAnimationFrame", loop)

	<-stallChan
}
