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

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	var console *Core.Console

	displayChannel := make(chan Core.ScreenInfo, 100)

	startWindow(console, displayChannel)

}

const (
	game_width  = 256
	game_heigth = 240
	menu_height = 20
)

func startWindow(console *Core.Console, updateChan chan Core.ScreenInfo) {
	setUpMenu()
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

	renderer.SetViewport(nil)
	renderer.SetDrawColor(200, 200, 200, 255)
	// renderer.DrawRect(&sdl.Rect{X: 0, Y: 0, W: 256, H: 40})
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: 256, H: 20})

	// screenTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, 256, 240)

	if err != nil {
		log.Fatal(err)
	}

	renderer.Present()
	// defer screenTexture.Destroy()

	audioSpec := sdl.AudioSpec{
		Freq:     44100,
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

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.WindowEvent:

			case *sdl.KeyboardEvent:
				if console != nil {
					handleInputs(console, e)
				} else {
					console = startGame("/Users/test/Projects/ExNES/games/NROM/mario.nes", audioDevice, updateChan)
				}

			default:

			}
		}

		// select {
		// case  := <-updateChan:

		// 	// renderFrame(screenTexture, renderer, s.Buffer)
		// default:

		// }

		sdl.Delay(1)
	}
}
