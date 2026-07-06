// #cgo CFLAGS: -Wno-deprecated-declarations

package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func setUpMenu() {
	fmt.Println("redit")
}

type menuBar struct {
	state *localState

	Items []menuItem
	Font  *ttf.Font

	dropdownIndex       int
	hoverIndex          int
	optionHoverIndex    int
	subOptionHoverIndex int

	cache TextureCache

	W     int32
	H     int32
	scale int32
}

type TextureCache struct {
	textCache  map[string]textCache
	hoverCache map[int32]*sdl.Texture
	arrowCache *sdl.Texture
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

	disabled bool

	Expandable      bool
	ExpandableItems []expandableOption

	isLine bool
}

type expandableOption struct {
	label   string
	Rect    sdl.Rect
	onClick func()
}

type textCache struct {
	texture *sdl.Texture
	W       int32
	H       int32
}

func createNewMenu(font *ttf.Font, console *game, state *localState, menuH, menuW, dpiscale int32) *menuBar {
	mb := &menuBar{

		dropdownIndex:       -1,
		hoverIndex:          -1,
		optionHoverIndex:    -1,
		subOptionHoverIndex: -1,
		Font:                font,
		W:                   menuW,
		H:                   menuH,
		scale:               dpiscale,
		state:               state,
	}

	cache := TextureCache{
		textCache:  map[string]textCache{},
		hoverCache: map[int32]*sdl.Texture{},
	}

	mb.cache = cache

	mb.Items = []menuItem{
		{
			label: "File",

			options: []ItemOption{
				{
					label:   "Load Rom",
					onClick: func() { openRom(console, state) },
				},
				{
					label:      "Recent files",
					Expandable: true,
					ExpandableItems: []expandableOption{
						{
							label: "hello",
						},
						{
							label: "bye",
						},
					},
				},
				{
					label:   "exit",
					onClick: func() { fmt.Println("exiting") },
				}},
		},
		{
			label: "Emulation",
			options: []ItemOption{
				{
					label: "Pause",
				},
				{
					label: "Play",
				},
			},
		},
		{
			label: "input",
			options: []ItemOption{
				{
					label: "Set controls",
				},
				{
					label: "pair controls",
				},
			},
		},
	}

	mb.positionLayout()

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

	if mb.optionHoverIndex != -1 {
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

		}

	}

}

func pointInRect(rect sdl.Rect, x, y int32) bool {
	return x >= rect.X && x <= rect.X+rect.W && y >= rect.Y && y <= rect.Y+rect.H
}

func (mb *menuBar) handleClick(x, y int32) {
	if mb.hoverIndex != -1 {
		item := mb.Items[mb.hoverIndex]
		for _, option := range item.options {
			if pointInRect(option.Rect, x, y) {
				option.onClick()
				return
			}
		}

		for _, item := range mb.Items {
			if pointInRect(item.Rect, x, y) {
				return
			}
		}
		mb.hoverIndex = -1
	}

}

const (
	optionH       = 24
	lineH         = 5
	panelPadding  = 8
	optionPadding = 26
	minW          = 140
	expandableW   = 120
	itemSpacer    = 6
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

					subOption.Rect = sdl.Rect{
						X: subX,
						Y: subY,
						H: optionH,
						W: expandableW,
					}

					subY += optionH
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
	r.SetDrawColor(colBarBG.R, colBarBG.G, colBarBG.B, 255)
	r.FillRect(&sdl.Rect{X: 0, Y: 0, W: mb.W, H: mb.H + 4})

	for i, item := range mb.Items {
		r.SetDrawColor(230, 210, 66, 255)

		if i == mb.hoverIndex {

			pill := mb.getHoverPill(r, item.Rect.W, item.Rect.H)
			r.Copy(pill, nil, &sdl.Rect{X: item.Rect.X, Y: item.Rect.Y, W: item.Rect.W, H: item.Rect.H})

			for j, option := range mb.Items[i].options {
				if j == mb.optionHoverIndex {
					pill = mb.getHoverPill(r, option.Rect.W, option.Rect.H)
					r.Copy(pill, nil, &option.Rect)

					mb.renderSupDropdown(r, &mb.Items[i].options[j])
				}
			}
			mb.renderDropdown(r, &mb.Items[i])

		}

		mb.drawText(item.label, item.Rect, r, 12)

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

		if i == mb.subOptionHoverIndex {
			pill := mb.getHoverPill(r, subOption.Rect.W, subOption.Rect.H)
			r.Copy(pill, nil, &option.ExpandableItems[i].Rect)
		}
		mb.drawText(subOption.label, subOption.Rect, r, 12)
	}

}

func (mb *menuBar) renderDropdown(r *sdl.Renderer, item *menuItem) {

	first := item.options[0].Rect
	last := item.options[len(item.options)-1].Rect

	panel := sdl.Rect{
		X: first.X,
		Y: first.Y,
		W: first.W,
		H: (last.Y + last.H) - first.Y,
	}

	gfx.RoundedBoxColor(r, panel.X, panel.Y, panel.X+panel.W, panel.Y+panel.H, 8, colPanelBG)

	for _, option := range item.options {
		if option.isLine {
			midY := option.Rect.Y + option.Rect.H/2
			r.SetDrawColor(colSeparator.R, colSeparator.G, colAccent.B, 255)
			r.DrawLine(option.Rect.X, midY, option.Rect.X+option.Rect.W, midY)
			continue
		}

		mb.drawText(option.label, option.Rect, r, 14)

		if option.Expandable {

			arrowRect := sdl.Rect{
				X: option.Rect.X + option.Rect.W - 24,
				H: option.Rect.H,
				Y: option.Rect.Y,
				W: 16,
			}

			r.Copy(mb.getArrow(r), nil, &arrowRect)
		}

	}

}

func (m *menuBar) drawText(text string, rect sdl.Rect, r *sdl.Renderer, offet int32) {
	entry, ok := m.cache.textCache[text]
	if !ok {
		entry = textCache{}
		surface, err := m.Font.RenderUTF8Blended(text, colText)
		if err != nil {
			log.Fatal("bad", err)
		}
		defer surface.Free()

		entry.W = int32(surface.W / m.scale)
		entry.H = int32(surface.H / m.scale)

		texture, err := r.CreateTextureFromSurface(surface)
		if err != nil {
			log.Fatal(err)
		}

		entry.texture = texture
		m.cache.textCache[text] = entry

	}

	dst := sdl.Rect{
		X: rect.X + offet,
		Y: rect.Y + (rect.H-entry.H)/2,
		W: entry.W,
		H: entry.H,
	}

	r.Copy(entry.texture, nil, &dst)

}

const upfactor = 4

func (mb *menuBar) getHoverPill(r *sdl.Renderer, w, h int32) *sdl.Texture {
	entry, ok := mb.cache.hoverCache[w]
	if ok {
		return entry
	}

	texture, err := createHover(r, colHover, w, h)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	mb.cache.hoverCache[w] = texture
	return texture
}

func createHover(r *sdl.Renderer, col sdl.Color, w, h int32) (*sdl.Texture, error) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	defer sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	bigW := w * upfactor
	bigH := h * upfactor

	tex, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, bigW, bigH)
	if err != nil {
		return nil, fmt.Errorf("unable to create texture fuck: %v", err)
	}

	tex.SetBlendMode(sdl.BLENDMODE_BLEND)

	r.SetRenderTarget(tex)
	r.SetDrawColor(0, 0, 0, 0)
	r.Clear()

	bigR := int32(10 * upfactor)

	gfx.RoundedBoxColor(r, 0, 0, bigW-1, bigH-1, bigR, col)

	r.SetRenderTarget(nil)
	return tex, nil
}
