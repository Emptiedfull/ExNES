package main

import (
	"exnes/Core"

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

	ListeningFor    *ButtonGroup
	ListeningAction string

	State *localState

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

	buttons map[string]*ButtonGroup
}

type ButtonGroup struct {
	rect sdl.Rect

	actionButton Core.BUTTON

	Title     string
	TitleRect sdl.Rect

	mapButton    sdl.Rect
	mapBoundText string
	mapText      string
	mapTextRect  sdl.Rect

	TurboButton    sdl.Rect
	TurboBoundText string
	TurboText      string
	TurboTextRect  sdl.Rect
}

func (win *controlWindow) setUp(font, smallfont *ttf.Font, windowH, windowW int32) {

	win.controlMain = controlMain{
		font:      font,
		smallFont: smallfont,

		groupHoverItem:  "",
		buttonHoverItem: "",

		State: win.state,

		windowH: windowH,
		windowW: windowW,
	}

	win.controlMain.cache = ControlCache{
		textCache:  make(map[string]textCache),
		panelCache: make(panelCache),
	}

	win.controlMain.groups = map[string]*ControlGroup{"D-pad": {
		Title: "D-pad",
		buttons: map[string]*ButtonGroup{
			"Up": {
				Title:        "Up",
				actionButton: Core.ButtonUp,
			},
			"Down": {
				Title:        "Down",
				actionButton: Core.ButtonDown,
			}, "Right": {
				Title:        "Right",
				actionButton: Core.ButtonRight,
			}, "Left": {
				Title:        "Left",
				actionButton: Core.ButtonLeft,
			},
		},
	}, "Joypad": {
		Title: "Joypad",
		buttons: map[string]*ButtonGroup{
			"Joypad-A": {
				Title:        "Joypad-A",
				actionButton: Core.ButtonA,
			}, "Joypad-B": {
				Title:        "Joypad-B",
				actionButton: Core.ButtonB,
			},
		},
	}, "Special": {
		Title: "Special",
		buttons: map[string]*ButtonGroup{
			"Start": {
				Title:        "Start",
				actionButton: Core.ButtonStart,
			}, "Select": {
				Title:        "Select",
				actionButton: Core.ButtonSelect,
			},
		},
	}}

	win.controlMain.layout(windowH, windowW)
	win.controlMain.syncText()
}

func (main *controlMain) layout(windowH, windowW int32) {
	dpad := main.groups["D-pad"]

	dpad.rect = sdl.Rect{
		X: 20,
		Y: (windowH - int32(float32(windowH)/1.1)) / 2,
		W: int32(float32(windowW) / 3.5),
		H: int32(float32(windowH) / 1.1),
	}

	w, h, _ := main.font.SizeUTF8(dpad.Title)
	dpad.TitleRect = sdl.Rect{

		X: dpad.rect.X + 25,
		Y: dpad.rect.Y - int32(h)/4,
		W: int32(w) / 2,
		H: int32(h) / 2,
	}
	order := []string{"Up", "Down", "Left", "Right"}

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

	Joypad := main.groups["Joypad"]

	Joypad.rect = sdl.Rect{
		X: windowW - int32(float32(windowW)/3.5) - 20,
		W: int32(float32(windowW) / 3.5),
		H: (itemH+padding)*2 + 20,
		Y: (windowH - ((itemH+padding)*2 + 15)) / 2,
	}

	w, h, _ = main.font.SizeUTF8(Joypad.Title)
	Joypad.TitleRect = sdl.Rect{

		X: Joypad.rect.X + 25,
		Y: Joypad.rect.Y - int32(h)/4,
		W: int32(w) / 2,
		H: int32(h) / 2,
	}

	order = []string{"Joypad-A", "Joypad-B"}

	Y = Joypad.rect.Y + int32(25)
	for _, key := range order {
		group := Joypad.buttons[key]

		group.rect = sdl.Rect{
			W: Joypad.rect.W - 20,
			X: Joypad.rect.X + 10,
			H: itemH,
			Y: Y,
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

		Y += itemH + padding
	}

	special := main.groups["Special"]

	special.rect = sdl.Rect{
		X: dpad.rect.X + dpad.rect.W + 35,
		W: int32(float32(windowW) / 3.5),
		H: (itemH+padding)*2 + 20,
		Y: windowH - ((itemH+padding)*2 + 20) - 20,
	}

	w, h, _ = main.font.SizeUTF8(special.Title)
	special.TitleRect = sdl.Rect{

		X: special.rect.X + 25,
		Y: special.rect.Y - int32(h)/4,
		W: int32(w) / 2,
		H: int32(h) / 2,
	}

	order = []string{"Start", "Select"}

	for _, key := range order {
		group := Joypad.buttons[key]

		group.rect = sdl.Rect{
			W: Joypad.rect.W - 20,
			X: Joypad.rect.X + 10,
			H: itemH,
			Y: Y,
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

		Y += itemH + padding
	}

}

var (
	colControlPanelBG        = sdl.Color{R: 30, G: 30, B: 34, A: 200}
	colControlPanelBorder    = sdl.Color{R: 68, G: 60, B: 86, A: 255}
	colControlPanelBorderHov = sdl.Color{R: 130, G: 110, B: 220, A: 255}

	colButtonBG        = sdl.Color{R: 40, G: 38, B: 46, A: 150}
	colButtonBorder    = sdl.Color{R: 68, G: 64, B: 63, A: 255}
	colButtonBorderHov = sdl.Color{R: 130, G: 110, B: 220, A: 255}

	colFieldBG          = sdl.Color{R: 24, G: 24, B: 28, A: 255}
	colFieldBorder      = sdl.Color{R: 130, G: 128, B: 140, A: 255}
	colFieldBorderHov   = sdl.Color{R: 130, G: 110, B: 220, A: 255}
	colFieldBorderBound = sdl.Color{R: 130, G: 110, B: 220, A: 255}
	colFieldTextBound   = sdl.Color{R: 206, G: 203, B: 246, A: 255}

	colListeningBorder = sdl.Color{R: 226, G: 168, B: 63, A: 255}
	colListeningBG     = sdl.Color{R: 226, G: 168, B: 63, A: 80}
)

func (main *controlMain) renderBoxes(r *sdl.Renderer) {

	for _, group := range main.groups {
		active := false
		bordercol := colControlPanelBorder
		if group.Title == main.groupHoverItem {
			active = true
			bordercol = colControlPanelBorderHov
		}

		_ = active
		panel := group.rect

		main.cache.panelCache.drawRoundedRect(r, &panel, colControlPanelBG, true)
		main.cache.panelCache.drawRoundedRect(r, &panel, bordercol, false)
		drawText(group.Title, group.TitleRect, r, 0, colText, main.cache.textCache, main.font, false)

		for _, button := range group.buttons {

			buttonsActive := false
			bordercol := colButtonBorder
			if active && button.Title == main.buttonHoverItem {
				buttonsActive = true
				bordercol = colButtonBorderHov
			}

			main.cache.panelCache.drawRoundedRect(r, &button.rect, colButtonBG, true)
			main.cache.panelCache.drawRoundedRect(r, &button.rect, bordercol, false)

			drawText(button.Title, button.TitleRect, r, 0, colText, main.cache.textCache, main.smallFont, false)

			mapCol := colFieldBorder
			mapText := colTextDim
			if buttonsActive && main.actionHoverItem == "BIND" {
				mapCol = colFieldBorderHov
				mapText = colText
			}

			main.cache.panelCache.drawRoundedRect(r, &button.mapButton, colFieldBG, true)
			main.cache.panelCache.drawRoundedRect(r, &button.mapButton, mapCol, false)

			drawText(button.mapBoundText, button.mapButton, r, button.mapButton.W/4, mapText, main.cache.textCache, main.smallFont, true)
			drawText("BIND", button.mapTextRect, r, 12, mapText, main.cache.textCache, main.smallFont, false)

			turboCol := colFieldBorder
			turboText := colTextDim

			if buttonsActive && main.actionHoverItem == "TURBO" {

				turboCol = colFieldBorderHov
				turboText = colText
			}

			main.cache.panelCache.drawRoundedRect(r, &button.TurboButton, colFieldBG, true)
			main.cache.panelCache.drawRoundedRect(r, &button.TurboButton, turboCol, false)
			drawText("TURBO", button.TurboTextRect, r, 12, turboText, main.cache.textCache, main.smallFont, false)

		}
	}

	if main.ListeningFor != nil {

		if main.ListeningAction == "BIND" {

			main.cache.panelCache.drawRoundedRect(r, &main.ListeningFor.mapButton, colListeningBG, true)
			main.cache.panelCache.drawRoundedRect(r, &main.ListeningFor.mapButton, colListeningBorder, false)

			main.ListeningFor.mapBoundText = "Listening..."

		} else if main.ListeningAction == "TURBO" {
			main.cache.panelCache.drawRoundedRect(r, &main.ListeningFor.TurboButton, colListeningBG, true)
			main.cache.panelCache.drawRoundedRect(r, &main.ListeningFor.TurboButton, colListeningBorder, false)

			drawText("Listening...", main.ListeningFor.TurboButton, r, main.ListeningFor.TurboButton.W/4, colText, main.cache.textCache, main.smallFont, false)
		}
	}

}

func (main *controlMain) syncText() {
	for groupKey, group := range main.groups {
		for buttonKey, button := range group.buttons {
			if code, ok := main.State.Settings.Inputs.ActionToKey[button.actionButton]; ok {
				button.mapBoundText = sdl.GetScancodeName(code)
				main.groups[groupKey].buttons[buttonKey] = button
			}
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

func (main *controlMain) handleClick() {
	if (main.actionHoverItem != "") && (main.groupHoverItem != "") && (main.buttonHoverItem != "") {
		group := main.groups[main.groupHoverItem].buttons[main.buttonHoverItem]

		main.ListeningFor = group
		main.ListeningAction = main.actionHoverItem
	}
}

func (main *controlMain) handleListen(code sdl.Scancode) {

	if action := main.ListeningFor; action != nil {
		if main.ListeningAction == "BIND" {
			main.State.Settings.Inputs.AssignKey(code, main.ListeningFor.actionButton)
			main.syncText()
			main.ListeningFor = nil
		}

	}

}
