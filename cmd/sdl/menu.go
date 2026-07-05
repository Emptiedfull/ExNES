// #cgo CFLAGS: -Wno-deprecated-declarations

package main

import (
	"exnes/Core"
	"fmt"
	"log"

	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func setUpMenu() {
	fmt.Println("redit")
}

type menuBar struct {
	Items []menuItem
	Font  *ttf.Font

	textCache  map[string]*sdl.Texture
	hoverCache map[int32]*sdl.Texture

	dropdownIndex    int
	hoverIndex       int
	optionHoverIndex int

	W     int32
	H     int32
	scale int32
}

type menuItem struct {
	label   string
	Rect    sdl.Rect
	options []ItemOption
	onClick func()
}

type ItemOption struct {
	label   string
	Rect    sdl.Rect
	onClick func()

	disabled bool

	isLine bool
}

func createNewMenu(font *ttf.Font, console *Core.Console, menuH, menuW, dpiscale int32) *menuBar {
	mb := &menuBar{
		textCache:        map[string]*sdl.Texture{},
		dropdownIndex:    -1,
		hoverIndex:       -1,
		optionHoverIndex: -1,
		Font:             font,
		W:                menuW,
		H:                menuH,
		scale:            dpiscale,
	}

	mb.Items = []menuItem{
		{
			label: "File",
			onClick: func() {
				filename, err := dialog.File().Filter("NES ROM", "nes").Load()
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("game:", filename)
			},
			options: []ItemOption{
				{
					label:   "Open",
					onClick: func() { fmt.Println("YO BITCHES") },
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
		},
	}

	mb.positionLayout()

	return mb
}

func (mb *menuBar) handleMouse(x, y int32) {
	mb.hoverIndex = -1

	for i, item := range mb.Items {
		if pointInRect(item.Rect, x, y) {
			mb.hoverIndex = i

			return
		}

	}
}

func pointInRect(rect sdl.Rect, x, y int32) bool {
	return x >= rect.X && x <= rect.X+rect.W && y >= rect.Y && y <= rect.Y+rect.H
}

func (mb *menuBar) handleClick() {
	if mb.hoverIndex != -1 {
		mb.Items[mb.hoverIndex].onClick()
	}
}

const (
	optionH       = 20
	lineH         = 5
	panelPadding  = 8
	optionPadding = 26
	minW          = 100
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
			y += h

		}

		item.Rect = sdl.Rect{
			X: x,
			Y: 0,
			W: itemW,
			H: mb.H,
		}
		x += itemW
	}
}

func (mb *menuBar) renderBar(r *sdl.Renderer) {
	r.SetDrawColor(colBarBG.R, colBarBG.G, colBarBG.B, 255)
	r.FillRect(&sdl.Rect{X: 0, Y: 0, W: mb.W, H: mb.H})

	for i, item := range mb.Items {
		r.SetDrawColor(230, 210, 66, 255)

		if i == mb.hoverIndex {

			gfx.RoundedBoxColor(r, item.Rect.X+2, item.Rect.Y+2, item.Rect.X+item.Rect.W-4, item.Rect.Y+item.Rect.H-2, 10, colHover)
			mb.renderDropdown(r, &mb.Items[i])
		}

		mb.drawText(item.label, item.Rect, r)

	}

	r.SetDrawColor(colSeparator.R, colSeparator.G, colSeparator.B, colSeparator.A)
	r.DrawLine(0, mb.H+1, mb.W, mb.H+1)
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

	gfx.RoundedRectangleColor(r, panel.X, panel.Y, panel.X+panel.W, panel.Y+panel.H, 8, colPanelBorder)

	for _, option := range item.options {
		if option.isLine {
			midY := option.Rect.Y + option.Rect.H/2
			r.SetDrawColor(colSeparator.R, colSeparator.G, colAccent.B, 255)
			r.DrawLine(option.Rect.X, midY, option.Rect.X+option.Rect.W, midY)
			continue
		}

		mb.drawText(option.label, option.Rect, r)

	}

}
func (m *menuBar) drawText(text string, rect sdl.Rect, r *sdl.Renderer) {
	surface, err := m.Font.RenderUTF8Blended(text, colText)
	if err != nil {
		log.Fatal("bad", err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}
	defer texture.Destroy()

	logicalW := int32(surface.W / m.scale)
	logicalH := int32(surface.H / m.scale)

	dst := sdl.Rect{
		X: rect.X + 12,
		Y: rect.Y + (rect.H-logicalH)/2,
		W: logicalW,
		H: logicalH,
	}

	r.Copy(texture, nil, &dst)

}
