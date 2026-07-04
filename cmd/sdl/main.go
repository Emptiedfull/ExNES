package main

/*
#cgo pkg-config: sdl2
#include <SDL2/SDL.h>


*/

import "C"
import (
	"exnes/Core"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	var console *Core.Console

	displayChannel := make(chan Core.ScreenInfo, 100)

	startWindow(console, displayChannel)

}

const (
	game_width  = 256
	game_heigth = 240
	menu_height = 10
	scale       = 2
)

var (
	colBarBG     = sdl.Color{R: 240, G: 240, B: 241, A: 255}
	colHover     = sdl.Color{R: 210, G: 225, B: 250, A: 255}
	colSeparator = sdl.Color{R: 200, G: 200, B: 202, A: 255}
	colText      = sdl.Color{R: 20, G: 20, B: 20, A: 255}
)

func startWindow(console *Core.Console, updateChan chan Core.ScreenInfo) {

	windowW := int32(game_width * scale)
	windowH := int32((game_heigth + menu_height) * scale)

	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	window, err := sdl.CreateWindow("ExNES", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, windowW, windowH, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI)

	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}

	renderer.SetLogicalSize(windowW, windowH)
	defer renderer.Destroy()

	W_H, _ := window.GetSize()
	D_H, _, _ := renderer.GetOutputSize()

	fmt.Println(W_H, D_H)
	dpi_scale := D_H / W_H
	fmt.Println("dpi sclae:", dpi_scale)

	if err := ttf.Init(); err != nil {
		log.Fatal("fuck init failed:", err)
	}

	font, err := ttf.OpenFont("/System/Library/Fonts/SFNS.ttf", int(14*dpi_scale))
	if err != nil {
		log.Fatal("something bad here:", err)
	}

	font.SetHinting(ttf.HINTING_LIGHT)

	renderer.SetDrawColor(colBarBG.R, colBarBG.G, colBarBG.B, 255)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: windowW, H: windowH})

	mb := createNewMenu(font, (menu_height * 2), windowW, dpi_scale)
	mb.renderBar(renderer)

	renderer.Present()

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
			case *sdl.MouseMotionEvent:
				mb.handleMouse(e.X, e.Y)
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_LEFT {
					mb.handleClick()
				}
				fmt.Println("button:", e.Button)
			default:

			}
		}

		renderer.Clear()
		mb.renderBar(renderer)
		renderer.Present()

		sdl.Delay(1)
	}
}
