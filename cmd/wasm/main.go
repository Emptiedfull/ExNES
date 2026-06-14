//go:build js && wasm

package main

import (
	"exnes/Core"
	"fmt"
	"syscall/js"
)

var emu *Core.Console

func main() {
	emu = Core.InitializeConsole()
	fmt.Println("core initialized")
	js.Global().Set("startEmulator", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println(this, args)
		emu = Core.InitializeConsole()
		return nil
	}))

	select {}
}
