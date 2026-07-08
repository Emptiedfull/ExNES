package main

import "github.com/veandco/go-sdl2/sdl"

func (win *controlWindow) setUp() {
	win.boxes = []sdl.Rect{
		{
			X: 20,
			Y: 20,
			H: 100,
			W: 100,
		},
	}
}

func (win *controlWindow) renderBoxes() {
	for _, box := range win.boxes {
		win.renderer.SetDrawColor(colPanelBorder.R, colPanelBorder.G, colPanelBorder.B, 255)
		win.renderer.FillRect(&box)
	}
}
