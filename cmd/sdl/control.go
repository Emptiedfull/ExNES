package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type controlMain struct {
	font      *ttf.Font
	smallFont *ttf.Font

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
	rect sdl.Rect

	Title     string
	TitleRect sdl.Rect

	mapButton   sdl.Rect
	mapText     string
	mapTextRect sdl.Rect

	TurboButton   sdl.Rect
	TurboText     string
	TurboTextRect sdl.Rect
}

func (win *controlWindow) setUp(font, smallfont *ttf.Font, windowH, windowW int32) {

	win.controlMain = controlMain{
		font:      font,
		smallFont: smallfont,

		groupHoverItem:   "",
		buttonHoverIndex: -1,

		windowH: windowH,
		windowW: windowW,
	}

	win.controlMain.cache = ControlCache{
		textCache: make(map[string]textCache),
	}

	win.controlMain.groups = map[string]*ControlGroup{"d-pad": {
		Title: "D-pad",
		buttons: map[string]ButtonGroup{
			"up": {
				Title: "Up",
			},
			"down": {
				Title: "Down",
			}, "right": {
				Title: "Right",
			}, "left": {
				Title: "Left",
			},
		},
	}}

	win.controlMain.layout(windowH, windowW)
}

func (main *controlMain) layout(windowH, windowW int32) {
	dpad := main.groups["d-pad"]
	dpad.rect = sdl.Rect{
		X: 20,
		Y: (windowH - int32(float32(windowH)/1.1)) / 2,
		W: int32(float32(windowW) / 3),
		H: int32(float32(windowH) / 1.1),
	}

	w, h, _ := main.font.SizeUTF8(dpad.Title)
	dpad.TitleRect = sdl.Rect{
		X: dpad.rect.X + 25,
		Y: dpad.rect.Y - int32(h)/4,
		W: int32(w) / 2,
		H: int32(h) / 2,
	}
	order := []string{"up", "down", "left", "right"}

	padding := int32(15)
	n := int32(len(order))

	totalPadding := padding * (n + 1)

	itemH := (dpad.rect.H - 8 - totalPadding) / n

	Y := dpad.rect.Y + padding + 10
	for _, key := range order {
		group := dpad.buttons[key]

		group.rect = sdl.Rect{
			X: dpad.rect.X + 10,
			Y: Y,
			H: itemH,
			W: dpad.rect.W - 20,
		}

		w, h, _ := main.smallFont.SizeUTF8(group.Title)

		group.TitleRect = sdl.Rect{
			X: group.rect.X + 10,
			Y: group.rect.Y - int32(h)/2,
			W: int32(w),
			H: int32(h),
		}

		w, h, _ = main.smallFont.SizeUTF8("BIND:")

		group.mapTextRect = sdl.Rect{
			W: int32(w),
			H: int32(h),
			X: group.rect.X + 5,
			Y: group.rect.Y + 8,
		}

		group.mapButton = sdl.Rect{
			H: int32(h),
			Y: group.rect.Y + 8,
			X: group.rect.X + int32(w) + 5,
			W: group.rect.W - int32(w) - 10,
		}

		group.TurboButton = sdl.Rect{
			H: int32(h),
			Y: group.rect.Y + 14 + group.mapButton.H,
			X: group.rect.X + 5,
			W: group.rect.W - 10,
		}

		dpad.buttons[key] = group

		Y += itemH + padding
	}

	main.groups["d-pad"] = dpad

}

func (win *controlWindow) renderBoxes() {

	for _, group := range win.controlMain.groups {

		panel := group.rect
		gfx.RoundedBoxColor(win.renderer, panel.X, panel.Y, panel.X+panel.W, panel.Y+panel.H, 8, sdl.Color{R: 30, G: 30, B: 34, A: 200})
		gfx.RoundedRectangleColor(win.renderer, panel.X, panel.Y, panel.X+panel.W, panel.Y+panel.H, 8, colPanelBorder)
		drawText(group.Title, group.TitleRect, win.renderer, 0, colText, win.controlMain.cache.textCache, win.controlMain.font)

		for _, button := range group.buttons {

			// gfx.RoundedRectangleColor(win.renderer, button.rect.X, button.rect.Y, button.rect.X+button.rect.W, button.rect.Y+button.rect.H, 8)
			drawRoundedRect(win.renderer, button.rect, sdl.Color{68, 64, 63, 255}, false)
			drawText(button.Title, button.TitleRect, win.renderer, 0, colText, win.controlMain.cache.textCache, win.controlMain.smallFont)

			drawText("BIND:", button.mapTextRect, win.renderer, 12, colTextDim, win.controlMain.cache.textCache, win.controlMain.smallFont)

			drawRoundedRect(win.renderer, button.mapButton, colAccent, false)
			drawRoundedRect(win.renderer, button.TurboButton, colAccent, false)

		}
	}

}
