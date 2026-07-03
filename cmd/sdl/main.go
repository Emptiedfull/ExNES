package main

/*
#cgo pkg-config: sdl2
#include <SDL2/SDL.h>
#include "_cgo_export.h"

*/

import "C"
import (
	"exnes/Core"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	console := Core.Quickstart("/Users/test/Projects/ExNES/games/Mapper2/contra.nes")

	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("ExNES", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 256, 240, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)

	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}

	defer renderer.Destroy()

	screenTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, 256, 240)

	if err != nil {
		log.Fatal(err)
	}

	defer screenTexture.Destroy()

	audioSpec := sdl.AudioSpec{
		Freq:     44100 * 2,
		Format:   sdl.AUDIO_F32SYS,
		Channels: 1,
		Samples:  512,
	}

	audioDevice, err := sdl.OpenAudioDevice("", false, &audioSpec, nil, 0)
	if err != nil {
		log.Fatal("unable to begin audio device")
	}

	defer sdl.CloseAudioDevice(audioDevice)

	sdl.PauseAudioDevice(audioDevice, false)
	go beginSampleLoop(audioDevice, console)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.WindowEvent:

			case *sdl.KeyboardEvent:
				handleInputs(console, e)
			default:

			}
		}

		select {
		case s := <-console.ScreenChannel:
			renderFrame(screenTexture, renderer, s.Buffer)
		default:

		}

		sdl.Delay(1)
	}

}

const bufferMillis = 46

func backpressureThreshold(freq int32) uint32 {
	bytesPerSecond := uint32(freq) * 4
	return bytesPerSecond * bufferMillis / 1000
}

func renderFrame(texture *sdl.Texture, renderer *sdl.Renderer, buffer []byte) {
	texture.Update(nil, unsafe.Pointer(&buffer[0]), 256*4)
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func beginSampleLoop(dev sdl.AudioDeviceID, console *Core.Console) {
	const samples = 512

	sampleBuf := make([]float32, samples)
	threshold := backpressureThreshold(44100 * 2)

	for {

		for i := range samples {
			for !console.Apu.HasSample() {
				console.TickNoAudio()

			}
			sample := console.Apu.PopSample()
			sampleBuf[i] = sample
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
