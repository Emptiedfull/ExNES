package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Window interface {
	render()
	close()
	getID() uint32

	handleMouse(*sdl.MouseMotionEvent)
	handleClick(*sdl.MouseButtonEvent)
	handleInput(*sdl.KeyboardEvent)
	handleScroll(*sdl.MouseWheelEvent)
	handleTextInput(*sdl.TextInputEvent)
}

type cheatWindow struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	id       uint32
	open     bool

	console *game
	main    *cheatMain
}

func openCheatWindow(console *game) (Window, error) {

	win, err := sdl.CreateWindow("Cheats", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 350, 450, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_RESIZABLE|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetLogicalSize(350, 450)

	id, err := win.GetID()
	if err != nil {
		return nil, err
	}

	font := loadFont(14 * scale)

	CheatWindow := cheatWindow{
		id:       id,
		window:   win,
		renderer: renderer,
		open:     true,

		console: console,
		main: &cheatMain{
			windowW: 350,
			windowH: 450,
		},
	}

	CheatWindow.main.Setup(font, console.core.CheatEngine, console.core.GetName())

	return &CheatWindow, nil
}

func (win *cheatWindow) handleTextInput(e *sdl.TextInputEvent) {
	win.main.input.handleTextInput(e)

}

func (win *cheatWindow) getID() uint32 {
	return win.id
}

func (win *cheatWindow) handleClick(e *sdl.MouseButtonEvent) {
	pressed := e.State == sdl.PRESSED

	if pressed {
		win.main.handleMouseDown(e.X, e.Y)
	} else {
		win.main.handleMouseUp()
	}
}

func (win *cheatWindow) handleInput(e *sdl.KeyboardEvent) {
	win.main.input.handleKeyInput(e)
}

func (win *cheatWindow) handleMouse(e *sdl.MouseMotionEvent) {
	win.main.handleMouseMove(e.X, e.Y)
}

func (win *cheatWindow) handleScroll(e *sdl.MouseWheelEvent) {

	win.main.handleScroll(e.Y)
}

func (win *cheatWindow) render() {
	win.renderer.SetDrawColor(0, 0, 0, 255)
	win.renderer.Clear()
	win.main.render(win.renderer)
	win.renderer.Present()
}

func (win *cheatWindow) close() {
	windows[win.id] = nil
	win.window.Destroy()
	win.renderer.Destroy()

	win.open = false

}

type controlWindow struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	id       uint32
	open     bool
	state    *localState

	controlMain controlMain
}

func openControlWindow(state *localState) (Window, error) {

	font := loadFont(18 * scale)
	smallfont := loadFont(14 * scale)

	font.SetHinting(ttf.HINTING_LIGHT)
	win, err := sdl.CreateWindow("controls", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 750, 450, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_RESIZABLE|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {

		return nil, err
	}

	renderer, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err

	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetLogicalSize(750, 450)

	id, err := win.GetID()
	if err != nil {

		return nil, err
	}

	controlWin := &controlWindow{
		id:       id,
		window:   win,
		renderer: renderer,
		open:     true,
		state:    state,
	}

	w, h := renderer.GetLogicalSize()

	controlWin.setUp(font, smallfont, h, w)

	return controlWin, nil

}

func (win *controlWindow) handleInput(e *sdl.KeyboardEvent) {
	if win.controlMain.ListeningFor != nil {
		if e.State == sdl.PRESSED {
			win.controlMain.handleListen(e.Keysym.Scancode)
		}

	}
}

func (win *controlWindow) handleTextInput(e *sdl.TextInputEvent) {

}

func (win *controlWindow) render() {

	win.controlMain.renderBoxes(win.renderer)
	win.renderer.Present()
}

func (win *controlWindow) getID() uint32 {
	return win.id
}

func (win *controlWindow) handleMouse(e *sdl.MouseMotionEvent) {
	win.controlMain.handleMouse(e.X, e.Y)

}

func (win *controlWindow) handleScroll(e *sdl.MouseWheelEvent) {

}

func (win *controlWindow) handleClick(e *sdl.MouseButtonEvent) {
	if e.State == sdl.PRESSED {
		win.controlMain.handleClick()
	}

}

func (win *controlWindow) close() {
	windows[win.id] = nil
	win.renderer.Destroy()
	win.window.Destroy()
	win.controlMain.smallFont.Close()
	win.controlMain.font.Close()

	win.controlMain.cache.panelCache = make(panelCache)
	win.controlMain.cache.textCache = make(map[string]textCache)

	win.open = false
}

type gameWindow struct {
	window *sdl.Window

	renderer    *sdl.Renderer
	menuBar     *menuBar
	audioDevice sdl.AudioDeviceID

	gameTexture *sdl.Texture
	gameRect    *sdl.Rect

	id   uint32
	open bool
}

func openGameWindow(font *ttf.Font, console *game, state *localState) (Window, error) {

	windowW := int32(game_width * scale)
	windowH := int32((game_heigth + menu_height) * scale)

	win, err := sdl.CreateWindow("ExNES", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, windowW, windowH, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {

		return nil, err
	}

	id, err := win.GetID()
	if err != nil {
		return nil, err
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	renderer, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	renderer.SetLogicalSize(windowW, windowH)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	W_H, _ := win.GetSize()
	D_H, _, _ := renderer.GetOutputSize()

	dpi_scale := D_H / W_H

	mb := createNewMenu(font, console, state, (menu_height * scale), windowW, dpi_scale)

	gameTexture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		game_width,
		game_heigth,
	)

	if err != nil {
		return nil, err
	}

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

		return nil, err
	}

	state.running = true
	console.audioDevice = audioDevice

	sdl.StartTextInput()

	sdl.PauseAudioDevice(audioDevice, false)

	return &gameWindow{
		window: win,

		renderer:    renderer,
		audioDevice: audioDevice,
		menuBar:     mb,

		gameTexture: gameTexture,
		gameRect:    gameRect,

		open: true,
		id:   id,
	}, err
}

func (win *gameWindow) render() {
	win.renderer.Clear()

	select {
	case s := <-win.menuBar.console.screenChannel:
		renderFrame(win.gameTexture, win.renderer, s, win.gameRect)

	default:
		win.renderer.Copy(win.gameTexture, nil, win.gameRect)
	}

	win.menuBar.renderBar(win.renderer)
	win.menuBar.renderFps(win.renderer)

	renderToasts(win.renderer, win.menuBar.Font, win.menuBar.cache.textCache, *win.gameRect)

	win.renderer.Present()
}

func (win *gameWindow) handleTextInput(e *sdl.TextInputEvent) {

}

func (win *gameWindow) close() {

	delete(windows, win.id)
	win.renderer.Destroy()
	win.window.Destroy()
	win.menuBar.console.SaveRam()
	sdl.CloseAudioDevice(win.audioDevice)
	sdl.StopTextInput()

	win.menuBar.state.running = false
	win.open = false
}

func (win *gameWindow) getID() uint32 {
	return win.id
}

func (win *gameWindow) handleMouse(e *sdl.MouseMotionEvent) {
	win.menuBar.handleMouse(e.X, e.Y)
}

func (win *gameWindow) handleClick(e *sdl.MouseButtonEvent) {

	if e.Button == sdl.BUTTON_LEFT && e.Type == sdl.MOUSEBUTTONDOWN {
		win.menuBar.handleClick(e.X, e.Y)
	}

}

func (win *gameWindow) handleInput(e *sdl.KeyboardEvent) {

	handleInputs(win.menuBar.console, e, win.menuBar.state.Settings.Inputs, win.menuBar.state.Settings.TurboInputs)
}

func (win *gameWindow) handleScroll(e *sdl.MouseWheelEvent) {

}
