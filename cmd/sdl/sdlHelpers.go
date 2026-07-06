package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

func drawArrowIcon(r *sdl.Renderer, rect sdl.Rect, col sdl.Color) {

	size := max(rect.H/3, 4)

	cx := rect.X + rect.W/2
	cy := rect.Y + rect.H/2

	x1 := cx - size/2
	y1 := cy - size/2
	x2 := cx - size/2
	y2 := cy + size/2
	x3 := cx + size/2
	y3 := cy

	gfx.FilledTrigonColor(r, x1, y1, x2, y2, x3, y3, col)
}

func (mb *menuBar) getArrow(r *sdl.Renderer) *sdl.Texture {
	if mb.cache.arrowCache != nil {
		return mb.cache.arrowCache
	}

	arrow, err := createArrow(r, 16, 16, colText)
	if err != nil {
		fmt.Println("soemthing fuckass:", err)
		return nil
	}

	mb.cache.arrowCache = arrow
	return arrow
}

func createArrow(r *sdl.Renderer, w, h int32, col sdl.Color) (*sdl.Texture, error) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	defer sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	bigW := w * upfactor
	bigH := h * upfactor

	texture, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, bigW, bigH)
	if err != nil {
		return nil, fmt.Errorf("error creating texture:%v", err)
	}

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	r.SetRenderTarget(texture)
	r.SetDrawColor(0, 0, 0, 0)
	r.Clear()

	drawArrowIcon(r, sdl.Rect{X: 0, Y: 0, W: bigW, H: bigH}, col)

	r.SetRenderTarget(nil)
	return texture, nil
}
