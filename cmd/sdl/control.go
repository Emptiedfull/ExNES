package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type controlMain struct {
	font *ttf.Font

	groups           map[string]*ControlGroup
	groupHoverItem   string
	buttonHoverIndex int

	windowH int32
	windowW int32

	cache ControlCache
}

type ControlCache struct {
	textCache map[string]textCache
}

type ControlGroup struct {
	rect sdl.Rect

	Title     string
	TitleRect sdl.Rect

	buttons map[string]ButtonGroup
}

type ButtonGroup struct {
	rect  sdl.Rect
	label string
}

func (win *controlWindow) setUp(font *ttf.Font, windowH, windowW int32) {

	win.controlMain = controlMain{
		font: font,

		groupHoverItem:   "",
		buttonHoverIndex: -1,

		windowH: windowH,
		windowW: windowW,
	}

	win.controlMain.cache = ControlCache{
		textCache: make(map[string]textCache),
	}

	win.controlMain.groups = map[string]*ControlGroup{"d-pad": {
		Title: "d-pad",
		buttons: map[string]ButtonGroup{
			"up": {
				label: "Up",
			},
			"down": {
				label: "down",
			}, "right": {
				label: "right",
			}, "left": {
				label: "left",
			},
		},
	}}

	win.controlMain.layout(windowH, windowW)
}

func (main *controlMain) layout(windowH, windowW int32) {
	dpad := main.groups["d-pad"]
	dpad.rect = sdl.Rect{
		X: 20,
		Y: (windowH - int32(float32(windowH)/1.5)) / 2,
		W: int32(float32(windowW) / 3),
		H: int32(float32(windowH) / 1.5),
	}

	w, h, _ := main.font.SizeUTF8(dpad.Title)
	dpad.TitleRect = sdl.Rect{
		X: dpad.rect.X + 10,
		Y: dpad.rect.Y - int32(h)/4,
		W: int32(w) / 2,
		H: int32(h) / 2,
	}

	heightPerGroup := dpad.rect.H / 5

	Y := int32(dpad.rect.Y + 10)
	for key := range dpad.buttons {
		group := dpad.buttons[key]

		group.rect = sdl.Rect{
			X: dpad.rect.X,
			Y: Y,
			H: heightPerGroup - 10,
			W: dpad.rect.W,
		}

		Y += heightPerGroup

		dpad.buttons[key] = group

	}

	main.groups["d-pad"] = dpad

}

func (win *controlWindow) renderBoxes() {

	win.renderer.SetDrawColor(colPanelBorder.R, colPanelBorder.G, colPanelBorder.B, 255)

	for _, group := range win.controlMain.groups {
		win.renderer.DrawRect(&group.rect)
		drawText(group.Title, group.TitleRect, win.renderer, 0, colText, win.controlMain.cache.textCache, win.controlMain.font)

		for _, button := range group.buttons {
			fmt.Println("drawing rect", button.rect.Y)
			win.renderer.FillRect(&button.rect)
		}
	}

}
