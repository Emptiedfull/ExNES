package main

import (
	"exnes/Core"
	"fmt"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func startGame(filepath string, dev sdl.AudioDeviceID, screenChannel chan Core.ScreenInfo) *Core.Console {
	console := Core.Quickstart(filepath)
	console.Apu = Core.NewApu(44100, console)
	console.ScreenChannel = screenChannel

	go beginSampleLoop(dev, console)

	// go frameMonitor(console)

	return console

}

func frameMonitor(console *Core.Console) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	last := 0

	for range ticker.C {
		framesRun := console.Ppu.Frame - last

		fmt.Println("fps:", framesRun)

		last = console.Ppu.Frame
	}
}

const bufferMillis = 46

func backpressureThreshold(freq int32) uint32 {
	bytesPerSecond := uint32(freq) * 4
	return bytesPerSecond * bufferMillis / 1000
}

func beginSampleLoop(dev sdl.AudioDeviceID, console *Core.Console) {
	const samples = 512

	sampleBuf := make([]float32, samples)
	threshold := backpressureThreshold(44100)

	for {

		for i := range samples {
			for !console.Apu.HasSample() {
				console.TickNoAudio()

			}
			sample := console.Apu.PopSample()
			sampleBuf[i] = sample

		}
		for sdl.GetQueuedAudioSize(dev) > threshold {
			sdl.Delay(1)
		}
		console.RunDisplayUpdates()
		sdl.QueueAudio(dev, float32ToBytes(sampleBuf))
	}

}

func float32ToBytes(samples []float32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&samples[0])), len(samples)*4)
}

var controlMap = map[sdl.Keycode]int{
	sdl.K_z:      Core.ButtonA,
	sdl.K_x:      Core.ButtonB,
	sdl.K_UP:     Core.ButtonUp,
	sdl.K_DOWN:   Core.ButtonDown,
	sdl.K_LEFT:   Core.ButtonLeft,
	sdl.K_RIGHT:  Core.ButtonRight,
	sdl.K_LSHIFT: Core.ButtonSelect,
	sdl.K_RETURN: Core.ButtonStart,
}

func handleInputs(console *Core.Console, e *sdl.KeyboardEvent) {

	pressed := e.State == sdl.PRESSED

	console.Player1.UpdateBtnBool(controlMap[e.Keysym.Sym], pressed)
}
