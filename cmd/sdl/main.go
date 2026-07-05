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
	"unsafe"

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
	menu_height = 12
	scale       = 2
)

var (
	colBarBG       = sdl.Color{R: 24, G: 24, B: 28, A: 255}
	colHover       = sdl.Color{R: 88, G: 74, B: 168, A: 255}
	colSeparator   = sdl.Color{R: 12, G: 12, B: 15, A: 255}
	colText        = sdl.Color{R: 230, G: 228, B: 235, A: 255}
	colTextDim     = sdl.Color{R: 130, G: 128, B: 140, A: 255}
	colAccent      = sdl.Color{R: 130, G: 110, B: 220, A: 255}
	colPanelBG     = sdl.Color{R: 40, G: 40, B: 44, A: 235}
	colPanelBorder = sdl.Color{R: 68, G: 60, B: 86, A: 255}
)

func startWindow(console *Core.Console, updateChan chan Core.ScreenInfo) {

	windowW := int32(game_width * scale)
	windowH := int32((game_heigth + menu_height) * scale)

	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("ExNES", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, windowW, windowH, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI)

	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		log.Fatal(err)
	}

	renderer.SetLogicalSize(windowW, windowH)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

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

	mb := createNewMenu(font, console, (menu_height * 2), windowW, dpi_scale)
	mb.renderBar(renderer)

	renderer.Present()

	gameTexture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		game_width,
		game_heigth,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer gameTexture.Destroy()

	gameRect := &sdl.Rect{
		X: 0,
		Y: menu_height * scale,
		H: game_heigth * scale,
		W: game_width * scale,
	}

	renderer.Copy(gameTexture, nil, gameRect)

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
				if e.Button == sdl.BUTTON_LEFT && e.Type == sdl.MOUSEBUTTONDOWN {
					mb.handleClick()
				}
				fmt.Println("button:", e.Button)
			default:

			}
		}
		renderer.Clear()

		select {
		case s := <-updateChan:
			renderFrame(gameTexture, renderer, s.Buffer, gameRect)
			if console.Ppu.Frame%10 == 0 {
				console.TakeSnapshot()
			}
		default:
			renderer.Copy(gameTexture, nil, gameRect)
		}

		mb.renderBar(renderer)
		renderer.Present()

		sdl.Delay(1)
	}
}

func renderFrame(texture *sdl.Texture, renderer *sdl.Renderer, buffer []byte, gameRect *sdl.Rect) {
	texture.Update(nil, unsafe.Pointer(&buffer[0]), 256*4)

	renderer.Copy(texture, nil, gameRect)

}
