package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type controlMain struct {
	font      *ttf.Font
	smallFont *ttf.Font

	groups          map[string]*ControlGroup
	groupHoverItem  string
	buttonHoverItem string
	actionHoverItem string

	windowH int32
	windowW int32

	cache ControlCache
}

type ControlCache struct {
	textCache  map[string]textCache
	panelCache panelCache
}

type panelCache map[string]*sdl.Texture

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

		groupHoverItem:  "",
		buttonHoverItem: "",

		windowH: windowH,
		windowW: windowW,
	}

	win.controlMain.cache = ControlCache{
		textCache:  make(map[string]textCache),
		panelCache: make(panelCache),
	}

	win.controlMain.groups = map[string]*ControlGroup{"D-pad": {
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
	dpad := main.groups["D-pad"]
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

		w, h, _ = main.smallFont.SizeUTF8("TURBO")
		wx, hx := int32(w)/scale+25, int32(h)/scale+7

		group.mapTextRect = sdl.Rect{
			W: int32(wx),
			H: int32(hx),
			X: group.rect.X + 5,
			Y: group.rect.Y + 16,
		}

		group.mapButton = sdl.Rect{
			H: int32(hx),
			Y: group.rect.Y + 16,
			X: group.rect.X + int32(wx),
			W: group.rect.W - int32(wx) - 8,
		}

		group.TurboTextRect = sdl.Rect{
			W: int32(wx),
			H: int32(hx),
			X: group.rect.X + 5,
			Y: group.mapButton.Y + 10 + group.mapButton.H,
		}

		group.TurboButton = sdl.Rect{
			H: int32(hx),
			Y: group.mapButton.Y + 10 + group.mapButton.H,
			X: group.rect.X + int32(wx),
			W: group.rect.W - int32(wx) - 8,
		}

		dpad.buttons[key] = group

		Y += itemH + padding
	}

	main.groups["d-pad"] = dpad

}

func (main *controlMain) renderBoxes(r *sdl.Renderer) {

	for _, group := range main.groups {
		active := false
		bordercol := colPanelBorder
		if group.Title == main.groupHoverItem {
			active = true
			bordercol = colAccent
		}

		_ = active
		panel := group.rect

		main.cache.panelCache.drawRoundedRect(r, &panel, sdl.Color{R: 30, G: 30, B: 34, A: 200}, true)
		main.cache.panelCache.drawRoundedRect(r, &panel, bordercol, false)
		drawText(group.Title, group.TitleRect, r, 0, colText, main.cache.textCache, main.font)

		for _, button := range group.buttons {

			buttonsActive := false
			bordercol := sdl.Color{68, 64, 63, 255}
			if active && button.Title == main.buttonHoverItem {
				buttonsActive = true
				bordercol = colAccent
			}

			_ = buttonsActive

			main.cache.panelCache.drawRoundedRect(r, &button.rect, bordercol, false)

			drawText(button.Title, button.TitleRect, r, 0, colText, main.cache.textCache, main.smallFont)

			mapCol := colPanelBorder
			mapText := colTextDim
			if buttonsActive && main.actionHoverItem == "BIND" {
				mapCol = colAccent
				mapText = colText
			}

			main.cache.panelCache.drawRoundedRect(r, &button.mapButton, mapCol, false)
			drawText("BIND", button.mapTextRect, r, 12, mapText, main.cache.textCache, main.smallFont)

			turboCol := colPanelBorder
			turboText := colTextDim

			if buttonsActive && main.actionHoverItem == "TURBO" {

				turboCol = colAccent
				turboText = colText
			}

			drawText("TURBO", button.TurboTextRect, r, 12, turboText, main.cache.textCache, main.smallFont)
			main.cache.panelCache.drawRoundedRect(r, &button.TurboButton, turboCol, false)

		}
	}

}

func (main *controlMain) handleMouse(x, y int32) {
	main.groupHoverItem = ""
	main.buttonHoverItem = ""
	main.actionHoverItem = ""
	for _, group := range main.groups {

		if pointInRect(group.rect, x, y) {

			main.groupHoverItem = group.Title

			for _, button := range group.buttons {
				if pointInRect(button.rect, x, y) {
					main.buttonHoverItem = button.Title

					if pointInRect(button.TurboButton, x, y) {
						main.actionHoverItem = "TURBO"
					}

					if pointInRect(button.mapButton, x, y) {
						main.actionHoverItem = "BIND"
					}
				}
			}

		}

	}

}

func (main *controlMain) handleClick()
