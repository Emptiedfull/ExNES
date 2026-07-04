// #cgo CFLAGS: -Wno-deprecated-declarations

package main

import (
	"fmt"
	"log"

	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func setUpMenu() {
	fmt.Println("redit")
}

type menuBar struct {
	Items     []menuItem
	Font      *ttf.Font
	textCache map[string]*sdl.Texture

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
	options []ItemOptions
	onClick func()
}

type ItemOptions struct {
	label   string
	Rect    sdl.Rect
	onClick func()
}

func createNewMenu(font *ttf.Font, menuH, menuW, dpiscale int32) *menuBar {
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
			options: []ItemOptions{
				{
					label:   "Open",
					onClick: func() { fmt.Println("YO BITCHES") },
				},
			},
		},
		{
			label: "Emulation",
			options: []ItemOptions{
				{
					label: "Pause",
				},
				{
					label: "Play",
				},
			},
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

func (mb *menuBar) positionLayout() {
	x := int32(0)
	for i := range mb.Items {
		item := &mb.Items[i]
		w, _, _ := mb.Font.SizeUTF8(item.label)
		padding := int32(24)
		itemW := int32(w)/mb.scale + padding
		fmt.Println(itemW)
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
			r.SetDrawColor(colHover.R, colBarBG.G, colBarBG.B, 255)
			r.FillRect(&item.Rect)
			r.SetDrawColor(0, 0, 0, 255)
			r.DrawRect(&item.Rect)
		}

		mb.drawText(item.label, item.Rect, r)
		// r.SetDrawColor(colSeparator.R, colSeparator.G, colSeparator.B, colSeparator.A)
		// r.DrawLine(item.Rect.X+item.Rect.W, 2, item.Rect.X+item.Rect.W, mb.H-2)
	}

	r.SetDrawColor(colSeparator.R, colSeparator.G, colSeparator.B, colSeparator.A)
	r.DrawLine(0, mb.H+1, mb.W, mb.H+1)
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
		Y: (rect.H - logicalH) / 2,
		W: logicalW,
		H: logicalH,
	}

	r.Copy(texture, nil, &dst)

}
