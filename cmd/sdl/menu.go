package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/gfx"
	_ "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Menuflag int

const (
	none Menuflag = iota
	gameRunning
	saveAvailable
	gamePlaying
	gamePaused
)

type MenuIcon int

const (
	pause MenuIcon = iota
	play
)

type menuBar struct {
	state   *localState
	console *game

	Items []menuItem
	Font  *ttf.Font

	dropdownIndex       int
	hoverIndex          int
	optionHoverIndex    int
	subOptionHoverIndex int

	cache TextureCache

	flagsUpdated bool
	menuFlags    map[Menuflag]bool
	affected     map[Menuflag][]*ItemOption // [flag]affected FUCK U THIS ISNT AI

	W     int32
	H     int32
	scale int32
}

type TextureCache struct {
	textCache  map[string]textCache
	hoverCache map[int32]*sdl.Texture
	arrowCache map[sdl.Color]*sdl.Texture
	iconCache  map[string]*sdl.Texture
}

type menuItem struct {
	label   string
	Rect    sdl.Rect
	options []ItemOption
}

type ItemOption struct {
	label   string
	Rect    sdl.Rect
	onClick func()

	affectedFlag Menuflag

	enabled bool

	Expandable      bool
	ExpandableItems []expandableOption

	isLine bool

	Icon string
}

type expandableOption struct {
	label   string
	Rect    sdl.Rect
	onClick func()

	Icon    string
	enabled bool

	isLine bool
}

type textCache struct {
	texture *sdl.Texture
	W       int32
	H       int32
	fontW   int32
}

func createNewMenu(font *ttf.Font, console *game, state *localState, menuH, menuW, dpiscale int32) *menuBar {

	mb := &menuBar{
		console:             console,
		dropdownIndex:       -1,
		hoverIndex:          -1,
		optionHoverIndex:    -1,
		subOptionHoverIndex: -1,
		Font:                font,
		W:                   menuW,
		H:                   menuH,
		scale:               dpiscale,
		state:               state,

		menuFlags: map[Menuflag]bool{},
		affected:  map[Menuflag][]*ItemOption{},
	}

	cache := TextureCache{
		textCache:  map[string]textCache{},
		hoverCache: map[int32]*sdl.Texture{},
		iconCache:  map[string]*sdl.Texture{},
		arrowCache: map[sdl.Color]*sdl.Texture{},
	}

	mb.cache = cache

	mb.Items = []menuItem{
		{
			label: "File",

			options: []ItemOption{
				{
					label:   "Load Rom",
					onClick: func() { openRom(console, state, mb) },
					enabled: true,
					Icon:    "icons/file-down.svg",
				},
				{
					label:           "Recent files",
					Expandable:      true,
					enabled:         true,
					ExpandableItems: getRecentItems(state.RecentFiles, console, mb),
				},
				{
					isLine: true,
				},
				{
					label:           "Save state",
					Expandable:      true,
					Icon:            "icons/save-pen.svg",
					ExpandableItems: mb.getSaveStateItems(),
					enabled:         false,
					affectedFlag:    gameRunning,
				},
				{
					label:           "Load state",
					Icon:            "icons/save.svg",
					Expandable:      true,
					ExpandableItems: make([]expandableOption, 0),
					enabled:         false,
					affectedFlag:    saveAvailable,
				},
				{
					isLine: true,
				},
				{
					label:   "Exit",
					enabled: true,
					Icon:    "icons/x.svg",
					onClick: func() { state.running = false },
				}},
		},
		{
			label: "Game",
			options: []ItemOption{
				{
					enabled:      false,
					Icon:         "icons/flag.svg",
					label:        "Cheats",
					affectedFlag: gameRunning,
					onClick: func() {
						win, _ := openCheatWindow(console)

						windows[win.getID()] = win
					},
				}, {
					isLine: true,
				},

				{
					enabled:      false,
					affectedFlag: gameRunning,
					Icon:         "icons/reset.svg",
					label:        "Reset",
					onClick: func() {
						console.core.Cpu.Reset()
					},
				},
				{
					enabled:      false,
					affectedFlag: gameRunning,
					Icon:         "icons/refresh.svg",
					label:        "Power cycle",
					onClick: func() {
						console.core.PowerCycle()
					},
				},
				{
					label:        "Reload Rom",
					affectedFlag: gameRunning,
					Icon:         "icons/download.svg",
					enabled:      false,
					onClick: func() {
						console.reloadROM()
					},
				},
				{
					isLine: true,
				},
				{
					enabled:      false,
					affectedFlag: gamePlaying,
					label:        "Pause",
					Icon:         "icons/pause.svg",
					onClick: func() {
						go console.PauseGame()
						mb.setFlag(gamePlaying, false)
						mb.setFlag(gamePaused, true)
					},
				},
				{
					enabled:      false,
					affectedFlag: gamePaused,
					label:        "Play",
					Icon:         "icons/play.svg",
					onClick: func() {
						go console.UnPauseGame()
						mb.setFlag(gamePlaying, true)
						mb.setFlag(gamePaused, false)
					},
				},
			},
		},
		{
			label: "Settings",
			options: []ItemOption{
				{
					label:        "Speed",
					Icon:         "icons/gauge.svg",
					enabled:      false,
					affectedFlag: gameRunning,
					Expandable:   true,
					ExpandableItems: []expandableOption{
						{
							label:   "Show Fps",
							enabled: true,
							onClick: func() {
								mb.state.Settings.Show_fps = !mb.state.Settings.Show_fps
								mb.updateSettingsMenu()
							},
						},

						{
							isLine: true,
						},
						{
							enabled: true,
							label:   "50%",
						}, {
							enabled: true,
							label:   "75%",
						}, {
							enabled: true,
							label:   "100%",
						}, {
							enabled: true,
							label:   "200%",
						},
					},
				},
				{
					label: "Sound",

					affectedFlag: gameRunning,
					Icon:         "icons/headset.svg",

					Expandable: true,
					ExpandableItems: []expandableOption{
						{
							label:   "Muted",
							enabled: true,
						},
						{
							isLine: true,
						},
						{
							label: "25%",
						}, {
							label: "50%",
						},
						{
							label: "75%",
						},
						{
							label: "100%",
						},
					},
				},
				{
					isLine: true,
				},
				{
					label:   "Inputs",
					enabled: true,
					Icon:    "icons/gamepad.svg",
					onClick: func() {
						for _, window := range windows {
							if _, ok := window.(*controlWindow); ok {

								return
							}
						}
						win, err := openControlWindow(state)
						if err != nil {
							panic(err)
						}
						windows[win.getID()] = win
					},
				}, {
					label:           "Palettes",
					enabled:         true,
					Icon:            "icons/palette.svg",
					Expandable:      true,
					ExpandableItems: make([]expandableOption, 0),
				},
			},
		},
	}
	mb.setupFlags()
	mb.setupMenus()

	return mb
}

func (mb *menuBar) handleMouse(x, y int32) {

	if mb.hoverIndex != -1 {
		item := mb.Items[mb.hoverIndex]

		for i, option := range item.options {
			if pointInRect(option.Rect, x, y) {
				mb.optionHoverIndex = i

			}
		}

	}

	if mb.optionHoverIndex != -1 && mb.hoverIndex != -1 {
		option := mb.Items[mb.hoverIndex].options[mb.optionHoverIndex]

		if len(option.ExpandableItems) == 0 {
			mb.subOptionHoverIndex = -1
		}

		for i, suboption := range option.ExpandableItems {
			if pointInRect(suboption.Rect, x, y) {
				mb.subOptionHoverIndex = i
			}
		}
	} else {
		mb.subOptionHoverIndex = -1
	}

	for i, item := range mb.Items {
		if pointInRect(item.Rect, x, y) {
			mb.hoverIndex = i
			mb.optionHoverIndex = -1
		}

	}

}

func (mb *menuBar) handleClick(x, y int32) {
	if mb.hoverIndex != -1 {
		item := mb.Items[mb.hoverIndex]

		for i, option := range item.options {
			if !option.enabled || option.isLine {
				continue
			}
			if option.Expandable {
				if i == mb.optionHoverIndex {
					for _, subOption := range option.ExpandableItems {
						if pointInRect(subOption.Rect, x, y) {
							if subOption.onClick != nil {
								subOption.onClick()
							}

							return
						}
					}
				}
			} else if pointInRect(option.Rect, x, y) {
				if option.onClick != nil && option.enabled {
					option.onClick()
				}

				return
			}

		}

		for _, item := range mb.Items {
			if pointInRect(item.Rect, x, y) {
				return
			}
		}

		mb.resetMenu()

	}

}

const (
	optionH       = 24
	subOptionH    = 24
	lineH         = 8
	panelPadding  = 8
	optionPadding = 0
	minW          = 140
	expandableW   = 120
	itemSpacer    = 6

	iconGutter = 24
	iconSize   = 12
)

func (mb *menuBar) positionLayout() {
	x := int32(6)
	for i := range mb.Items {
		item := &mb.Items[i]
		w, _, _ := mb.Font.SizeUTF8(item.label)
		padding := int32(24)
		itemW := int32(w)/mb.scale + padding

		dropdownW := int32(minW)
		for _, option := range item.options {
			if option.isLine {
				continue
			}

			ow, _, _ := mb.Font.SizeUTF8(option.label)
			padded := int32(ow)/mb.scale + optionPadding
			if padded > dropdownW {
				dropdownW = padded
			}

		}

		y := mb.H + 6
		for j := range item.options {
			option := &item.options[j]
			h := int32(optionH)
			if option.isLine {
				h = lineH
			}

			option.Rect = sdl.Rect{
				X: x,
				Y: y,
				H: h,
				W: dropdownW,
			}

			if option.Expandable && len(option.ExpandableItems) > 0 {
				subX := option.Rect.X + option.Rect.W + 4
				subY := option.Rect.Y
				for k := range option.ExpandableItems {
					subOption := &option.ExpandableItems[k]

					opH := int32(subOptionH)

					if subOption.isLine {
						opH = lineH
					}

					subOption.Rect = sdl.Rect{
						X: subX,
						Y: subY,
						H: opH,
						W: expandableW,
					}

					subY += opH
				}
			}

			y += h

		}

		item.Rect = sdl.Rect{
			X: x,
			Y: 0,
			W: itemW,
			H: mb.H,
		}
		x += itemW + itemSpacer
	}
}

func (mb *menuBar) renderBar(r *sdl.Renderer) {

	mb.updateMenuState()

	r.SetDrawColor(colBarBG.R, colBarBG.G, colBarBG.B, 255)
	r.FillRect(&sdl.Rect{X: 0, Y: 0, W: mb.W, H: mb.H + 4})

	for i, item := range mb.Items {
		r.SetDrawColor(230, 210, 66, 255)

		if i == mb.hoverIndex {

			pill := mb.getHoverPill(r, item.Rect.W, item.Rect.H)
			r.Copy(pill, nil, &sdl.Rect{X: item.Rect.X, Y: item.Rect.Y, W: item.Rect.W, H: item.Rect.H})

			for j, option := range mb.Items[i].options {
				if j == mb.optionHoverIndex && !option.isLine && option.enabled {
					pill = mb.getHoverPill(r, option.Rect.W, option.Rect.H)
					r.Copy(pill, nil, &option.Rect)

					mb.renderSupDropdown(r, &mb.Items[i].options[j])
				}
			}
			mb.renderDropdown(r, &mb.Items[i])

		}

		// drawText(item.label, item.Rect, r, 12, colText, mb.cache.textCache, mb.Font, false)

		drawText(item.label, r, mb.cache.textCache, TextOptions{rect: item.Rect, offset: 12, font: mb.Font, col: colText})

	}

}

func (mb *menuBar) renderSupDropdown(r *sdl.Renderer, option *ItemOption) {
	if len(option.ExpandableItems) == 0 {
		return
	}

	first := option.ExpandableItems[0]
	last := option.ExpandableItems[len(option.ExpandableItems)-1].Rect

	panel := sdl.Rect{
		X: first.Rect.X,
		Y: first.Rect.Y,
		H: (last.Y + last.H) - first.Rect.Y,
		W: first.Rect.W,
	}

	gfx.RoundedBoxColor(r, panel.X, panel.Y, panel.X+panel.W, panel.Y+panel.H, 8, colPanelBG)

	for i, subOption := range option.ExpandableItems {
		color := colText

		if !subOption.enabled {
			color = colTextDim
		}
		if subOption.isLine {
			midY := subOption.Rect.Y + subOption.Rect.H/2
			r.SetDrawColor(colHover.R, colHover.G, colHover.B, 255)
			r.DrawLine(subOption.Rect.X+10, midY, subOption.Rect.X+subOption.Rect.W-10, midY)
			continue

		}

		if i == mb.subOptionHoverIndex {
			pill := mb.getHoverPill(r, subOption.Rect.W, subOption.Rect.H)
			r.Copy(pill, nil, &option.ExpandableItems[i].Rect)
		}
		textRect := subOption.Rect

		if subOption.Icon != "" {
			icon := getIcon(r, subOption.Icon, colText, mb.cache.iconCache)

			if icon != nil {
				rect := sdl.Rect{
					X: subOption.Rect.X + 6,
					Y: subOption.Rect.Y + (subOptionH-iconSize)/2,
					H: iconSize,
					W: iconSize,
				}

				r.Copy(icon, nil, &rect)
			}
		}

		// drawText(subOption.label, textRect, r, iconGutter, color, mb.cache.textCache, mb.Font, false)

		drawText(subOption.label, r, mb.cache.textCache, TextOptions{col: color, font: mb.Font, offset: iconGutter, rect: textRect, clamped: true})
	}

}

func (mb *menuBar) renderDropdown(r *sdl.Renderer, item *menuItem) {

	first := item.options[0].Rect
	last := item.options[len(item.options)-1].Rect

	panel := sdl.Rect{
		X: first.X,
		Y: first.Y,
		W: first.W,
		H: (last.Y + last.H) - first.Y + 5,
	}

	gfx.RoundedBoxColor(r, panel.X, panel.Y, panel.X+panel.W, panel.Y+panel.H, 8, colPanelBG)

	for _, option := range item.options {
		if option.isLine {
			midY := option.Rect.Y + option.Rect.H/2
			r.SetDrawColor(colHover.R, colHover.G, colHover.B, 255)
			r.DrawLine(option.Rect.X+10, midY, option.Rect.X+option.Rect.W-10, midY)
			continue
		}

		color := colText
		if !option.enabled {
			color = colTextDim
		}

		if option.Icon != "" {
			icon := getIcon(r, option.Icon, color, mb.cache.iconCache)
			if icon != nil {
				iconRect := sdl.Rect{
					X: option.Rect.X + 6,
					Y: option.Rect.Y + (option.Rect.H-iconSize)/2,
					H: iconSize,
					W: iconSize,
				}

				r.Copy(icon, nil, &iconRect)
			}
		}

		// drawText(option.label, option.Rect, r, iconGutter, color, mb.cache.textCache, mb.Font, false)

		drawText(option.label, r, mb.cache.textCache, TextOptions{rect: option.Rect, font: mb.Font, col: color, offset: iconGutter})

		if option.Expandable {

			arrowRect := sdl.Rect{
				X: option.Rect.X + option.Rect.W - 24,
				H: option.Rect.H,
				Y: option.Rect.Y,
				W: 16,
			}

			r.Copy(mb.getArrow(r, color), nil, &arrowRect)
		}

	}

}

func (mb *menuBar) renderFps(r *sdl.Renderer) {
	if mb.state.Settings.Show_fps {
		label := fmt.Sprintf("%d Fps (%v)", mb.console.fps, mb.state.Settings.Current_speed)

		rect := sdl.Rect{
			X: mb.W - 100,
			Y: mb.H + menu_height,
			W: 100,
			H: 30,
		}

		// drawText(label, rect, r, 0, colText, mb.cache.textCache, mb.Font, false)

		drawText(label, r, mb.cache.textCache, TextOptions{rect: rect, font: mb.Font, col: colText})
	}
}
