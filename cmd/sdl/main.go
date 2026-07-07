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
	"time"
	"unsafe"

	_ "net/http/pprof"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {

	getSaveTimeStamp()
	displayChannel := make(chan Core.ScreenInfo, 100)
	pauseChannel := make(chan bool)

	g := &game{
		fps:          60,
		pauseChannel: pauseChannel,
	}

	g.initConsole(displayChannel)

	state := loadState()

	defer state.saveState()

	startWindow(g, displayChannel, state)

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
	colPanelBG     = sdl.Color{R: 40, G: 40, B: 44, A: 215}
	colPanelBorder = sdl.Color{R: 68, G: 60, B: 86, A: 255}
)

func startWindow(console *game, updateChan chan Core.ScreenInfo, state *localState) {

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
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
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

	mb := createNewMenu(font, console, state, (menu_height * 2), windowW, dpi_scale)
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

	sdl.PauseAudioDevice(audioDevice, true)
	console.audioDevice = audioDevice

	state.running = true
	visible := true

	const targetFPS = 60
	frameDuration := time.Second / targetFPS

	for state.running {

		start := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {

			case *sdl.QuitEvent:
				state.running = false

			case *sdl.WindowEvent:
				switch e.Event {
				case sdl.WINDOWEVENT_MINIMIZED, sdl.WINDOWEVENT_HIDDEN:
					visible = false

				case sdl.WINDOWEVENT_MAXIMIZED, sdl.WINDOWEVENT_SHOWN, sdl.WINDOWEVENT_EXPOSED:
					visible = true
				}

			case *sdl.KeyboardEvent:

				handleInputs(console.core, e)

			case *sdl.MouseMotionEvent:
				mb.handleMouse(e.X, e.Y)
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_LEFT && e.Type == sdl.MOUSEBUTTONDOWN {
					mb.handleClick(e.X, e.Y)
				}

			default:

			}
		}

		if visible {
			renderer.Clear()

			select {
			case s := <-updateChan:
				renderFrame(gameTexture, renderer, s.Buffer, gameRect)

			default:
				renderer.Copy(gameTexture, nil, gameRect)
			}
			mb.renderBar(renderer)
			mb.renderFps(renderer)

			renderer.Present()
		}

		elapsed := time.Since(start)
		if elapsed < frameDuration {
			sdl.Delay(uint32((frameDuration - elapsed).Milliseconds()))
		}

	}
}

func renderFrame(texture *sdl.Texture, renderer *sdl.Renderer, buffer []byte, gameRect *sdl.Rect) {
	texture.Update(nil, unsafe.Pointer(&buffer[0]), 256*4)
	renderer.Copy(texture, nil, gameRect)

}
