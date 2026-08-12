//go:build js && wasm

package main

import (
	"encoding/binary"
	"encoding/json"
	"exnes/Core"
	"fmt"
	"log"
	"math"
	"syscall/js"
	"time"
	"unsafe"
)

var currentSpeed js.Value

func main() {

	stallChan := make(chan bool)

	startFrameDriver()

	<-stallChan
}

func startFrameDriver() {
	fmt.Println("new version loaded:", 30)
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

		fmt.Println("rom ready for playing")

		return nil
	}))

	js.Global().Set("getSnapshotList", js.FuncOf(func(this js.Value, args []js.Value) any {
		ptr := makeSnapshotList(emu.Snapshots)

		fmt.Println("length:", emu.Snapshots.GetLength())

		return ptr

	}))

	js.Global().Set("loadSnapshot", js.FuncOf(func(this js.Value, args []js.Value) any {
		index := args[0].Int()

		emu.LoadSnapshot(emu.Snapshots.Data[index])
		return nil
	}))

	js.Global().Set("reset", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("reset requested")
		start := time.Now()
		if emu != nil {
			emu.Cpu.Reset()

			fmt.Println("reset took:", time.Since(start))
			return 1
		}

		return 1
	}))

	js.Global().Set("runFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		syncInputs(emu, Input_Arr)
		frame := emu.Ppu.Frame
		for emu.Ppu.Frame < frame+1 {
			emu.Step()
			for emu.Apu.HasSample() {
				emu.Apu.PopSample()
			}

		}

		buf := emu.Ppu.NewBuffer
		bytes := unsafe.Slice((*byte)(unsafe.Pointer(&buf[0])), len(buf)*4)
		js.CopyBytesToJS(jsScreen, bytes)

		return nil
	}))

	js.Global().Set("drive", js.FuncOf(func(this js.Value, args []js.Value) any {

		speed := getSpeed(Speed_Arr) / 1000.0

		samplesNeeded := int(float32(args[0].Int()) * speed)
		syncInputs(emu, Input_Arr)

		for i := range samplesNeeded {

			for !emu.Apu.HasSample() {
				emu.Apu.Console.Step()
			}
			sample := emu.Apu.PopSample()

			bits := math.Float32bits(sample)

			S_Arr[i*4] = byte(bits)
			S_Arr[i*4+1] = byte(bits >> 8)
			S_Arr[i*4+2] = byte(bits >> 16)
			S_Arr[i*4+3] = byte(bits >> 24)
		}

		if emu.Ppu.ScreenChanged {
			buf := emu.Ppu.NewBuffer
			bytes := unsafe.Slice((*byte)(unsafe.Pointer(&buf[0])), len(buf)*4)
			js.CopyBytesToJS(jsScreen, bytes)

			if emu.Ppu.Frame%10 == 0 {

				emu.TakeSnapshot()
			}

		}

		js.CopyBytesToJS(JS_Arr, S_Arr)

		return nil
	}))

}

type SnapInfo struct {
	Frameno int `json:"frame"`
	Index   int `json:"index"`
}

var outHeader [3]uint32
var snapshotExport []uint8

func makeSnapshotList(B Core.SnapshotBuffer) uintptr {
	length := B.GetLength()

	res := make([]SnapInfo, length)
	snapshotExport = make([]uint8, length*size)

	for i := range length {
		res[i] = SnapInfo{
			Frameno: B.Data[i].Frame_no,
			Index:   i,
		}

		buf := B.Data[i].PpuState.NewBuffer
		bytes := unsafe.Slice((*byte)(unsafe.Pointer(&buf[0])), len(buf)*4)

		copy(snapshotExport[i*size:], bytes)
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("FUCK FUCK FUCK")
	}

	outHeader[0] = uint32(uintptr(unsafe.Pointer(&b[0])))
	outHeader[1] = uint32(len(b))
	outHeader[2] = uint32(uintptr(unsafe.Pointer(&snapshotExport[0])))

	return uintptr(unsafe.Pointer(&outHeader[0]))
}

const size = 256 * 240 * 4

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
