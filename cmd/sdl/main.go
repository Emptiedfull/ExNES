package main

/*
#include <SDL2/SDL.h>
extern void audioCallback(void *userdata, Uint8 *stream, int len);
*/

import "C"
import (
	"exnes/Core"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

var console Core.Debugger

func main() {

	emulator := Core.Quickstart("/Users/test/Projects/ExNES/games/Mapper2/contra.nes")
	fmt.Println(emulator)

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

	c := make(chan bool)
	c <- true

}
