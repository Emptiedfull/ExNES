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

	_ "net/http/pprof"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	game_width  = 256
	game_heigth = 240
	menu_height = 12
)

var scale = int32(2)

var (
	colBarBG       = sdl.Color{R: 24, G: 24, B: 28, A: 255}
	colHover       = sdl.Color{R: 88, G: 74, B: 168, A: 255}
	colSeparator   = sdl.Color{R: 12, G: 12, B: 15, A: 255}
	colText        = sdl.Color{R: 230, G: 228, B: 235, A: 255}
	colTextDim     = sdl.Color{R: 130, G: 128, B: 140, A: 255}
	colAccent      = sdl.Color{R: 130, G: 110, B: 220, A: 255}
	colPanelBG     = sdl.Color{R: 40, G: 40, B: 44, A: 240}
	colPanelBorder = sdl.Color{R: 68, G: 60, B: 86, A: 255}
)

type Windows map[uint32]Window

var windows = make(Windows)

func main() {
	Core.Parse()

	// Core.DecodeCheat("SXIOPO")
	// panic("hi")

	displayChannel := make(chan Core.ScreenInfo, 100)
	pauseChannel := make(chan bool)

	g := &game{
		fps:          60,
		pauseChannel: pauseChannel,
		TurboState:   make(map[Core.BUTTON]chan bool),
	}

	g.initConsole(displayChannel)

	go func() {
		time.Sleep(1 * time.Second)
		g.core.CheatEngine.AddCheat("SUNNIZVI")
		g.core.CheatEngine.Enabled = true
		fmt.Println("cheat applied")
	}()

	state := loadState()

	g.changeVolume(state.Settings.Current_volume)

	defer state.saveState()

	if err := ttf.Init(); err != nil {
		log.Fatal("fuck init failed:", err)
	}

	font, err := ttf.OpenFont("/System/Library/Fonts/SFNS.ttf", int(14*2))
	if err != nil {
		log.Fatal("something bad here:", err)
	}

	font.SetHinting(ttf.HINTING_LIGHT)

	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	gameWin, err := openGameWindow(font, g, state)
	if err != nil {
		log.Fatal(err)
	}

	windows[gameWin.getID()] = gameWin

	startLoop(state)

}

func startLoop(state *localState) {

	for state.running {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch e := event.(type) {
			case *sdl.QuitEvent:
				state.running = false
			case *sdl.WindowEvent:
				id := e.WindowID

				if windows[id] != nil {
					if e.Event == sdl.WINDOWEVENT_CLOSE {
						windows[id].close()

					}
				}

			case *sdl.MouseMotionEvent:
				id := e.WindowID
				if window, ok := windows[id]; ok {
					window.handleMouse(e)
				}

			case *sdl.MouseButtonEvent:
				id := e.WindowID
				if window, ok := windows[id]; ok {
					window.handleClick(e)
				}

			case *sdl.KeyboardEvent:
				id := e.WindowID
				if window, ok := windows[id]; ok {
					window.handleInput(e)
				}
			}

		}

		for _, window := range windows {
			if window != nil {
				window.render()
			}

		}

		sdl.Delay(1)
	}
}
