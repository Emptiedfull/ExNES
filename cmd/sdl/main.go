package main

/*
#cgo pkg-config: sdl2
#include <SDL2/SDL.h>
#include "_cgo_export.h"
extern void audioCallback(void *userdata, unsigned char *stream, int len);

static void* getAudioCallback() {
    return (void*)audioCallback;
}
*/

import "C"
import (
	"exnes/Core"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var console *Core.Console

func main() {

	console = Core.Quickstart("/Users/test/Projects/ExNES/games/Mapper2/contra.nes")

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
		Freq:     44100,
		Format:   sdl.AUDIO_S16LSB,
		Channels: 1,
		Samples:  1024,
		Callback: sdl.AudioCallback(C.getAudioCallback()),
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.WindowEvent:

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

func renderFrame(texture *sdl.Texture, renderer *sdl.Renderer, buffer []byte) {
	texture.Update(nil, unsafe.Pointer(&buffer[0]), 256*4)
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}
