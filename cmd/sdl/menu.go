package main

import (
	"fmt"
	"log"

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
}

type menuItem struct {
	label   string
	Rect    sdl.Rect
	options []ItemOptions
}

type ItemOptions struct {
	label   string
	Rect    sdl.Rect
	onClick func()
}

func createNewMenu() *menuBar {
	mb := &menuBar{
		textCache:        map[string]*sdl.Texture{},
		dropdownIndex:    -1,
		hoverIndex:       -1,
		optionHoverIndex: -1,
	}

	mb.Items = []menuItem{
		{
			label: "FILE",
			options: []ItemOptions{
				{
					label:   "Open",
					onClick: openRom,
				},
			},
		},
		{
			label: "EMULATION",
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
	fmt.Println(mb)

	return mb
}

func (mb *menuBar) positionLayout() {
	x := int32(0)
	for _, item := range mb.Items {
		w := int32(game_width / len(mb.Items))
		item.Rect = sdl.Rect{X: x, Y: 0, H: menu_height, W: w}

		for j, option := range item.options {
			optW := int32(len(option.label)*3 + 30)
			if optW < 90 {
				optW = 90
			}

			option.Rect = sdl.Rect{
				X: x,
				Y: menu_height + int32(j)*menu_height,
			}
		}

		x += w
	}
}

func (mb *menuBar) renderBar(r *sdl.Renderer) {
	r.SetDrawColor(35, 35, 35, 255)
	r.FillRect(&sdl.Rect{X: 0, Y: 0, W: game_width, H: game_heigth + menu_height})

	for _, item := range mb.Items {

	}
}

func (item *menuItem) drawText(r *sdl.Renderer, font *ttf.Font) {
	surface, err := font.RenderUTF8Blended(item.label, sdl.Color{230, 230, 230, 255})
	if err != nil {
		log.Fatal("FUCK FUCK", err)
	}

	text, err := r.CreateTextureFromSurface(surface)
	surface.Free()
	if err != nil {
		log.Fatal("FUCK NO", err)
	}

	_, _, w, h, _ := text.Query()
	r.Copy(text, nil, &sdl.Rect{X: item.Rect.X + 10, Y: item.Rect.Y + 5, W: w, H: h})

}
