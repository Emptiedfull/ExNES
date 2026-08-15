package main

import (
	"exnes/Core"
	"fmt"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var errEngine errorEngine

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

	displayChannel := make(chan []uint32, 100)
	pauseChannel := make(chan bool)

	g := &game{
		fps:          60,
		pauseChannel: pauseChannel,
		TurboState:   make(map[Core.BUTTON]chan bool),
	}

	g.initConsole(displayChannel)

	state := loadState()

	g.changeVolume(state.Settings.Current_volume)

	defer state.saveState()

	if err := ttf.Init(); err != nil {
		log.Fatal("fuck init failed:", err)
	}

	font := loadFont(14 * scale)

	font.SetHinting(ttf.HINTING_LIGHT)

	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.EventState(sdl.DROPBEGIN, sdl.IGNORE)
	sdl.EventState(sdl.DROPCOMPLETE, sdl.IGNORE)
	sdl.EventState(sdl.DROPFILE, sdl.IGNORE)

	gameWin, err := openGameWindow(font, g, state)
	if err != nil {
		log.Fatal(err)
	}

	windows[gameWin.getID()] = gameWin

	if len(os.Args) > 1 {
		if w, ok := gameWin.(*gameWindow); ok {
			err := g.LoadRom(os.Args[1], w.menuBar)
			if err != nil {
				fmt.Println("error opening rom:", err)
			}
		}
	}

	startLoop(state)

}

func startLoop(state *localState) {

	for state.running {

		sdl.FlushEvents(sdl.DROPFILE, sdl.DROPCOMPLETE)

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
			case *sdl.MouseWheelEvent:
				id := e.WindowID
				if Window, ok := windows[id]; ok {
					Window.handleScroll(e)
				}
			case *sdl.DropEvent:
				if e.Type == sdl.DROPFILE {
					if w, ok := windows[e.WindowID].(*gameWindow); ok {
						err := w.menuBar.console.LoadRom(e.File, w.menuBar)
						if err != nil {
							fmt.Println("error loading file:", err)
						}
					}
				}
			case *sdl.TextInputEvent:
				id := e.WindowID
				if window, ok := windows[id]; ok {
					window.handleTextInput(e)
				}
			}

		}

		for _, window := range windows {
			if window != nil {
				window.render()
			}

		}

	}
}
