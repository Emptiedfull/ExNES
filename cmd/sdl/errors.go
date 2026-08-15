package main

import (
	"sync"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type errorToast struct {
	msg    string
	begin  time.Time
	Severe bool
}

var (
	toasts   []errorToast
	toastMux sync.Mutex
)

func pushError(src string, err error, severe bool) {
	if err == nil {
		return
	}

	toastMux.Lock()
	defer toastMux.Unlock()

	toasts = append(toasts, errorToast{
		msg:    src + ":" + err.Error(),
		begin:  time.Now(),
		Severe: severe,
	})

	if len(toasts) > 5 {
		toasts = toasts[len(toasts)-5:]
	}

}

const fontSample = 2

func renderToasts(r *sdl.Renderer, font *ttf.Font, cache map[string]textCache, area sdl.Rect) {
	toastMux.Lock()
	defer toastMux.Unlock()

	y := area.Y + menu_height*2 + 10

	for _, t := range toasts {
		w, h, err := font.SizeUTF8(t.msg)
		if err != nil {
			continue
		}

		tw := max(int32(w)/fontSample+20, area.W-12)
		th := int32(h)/fontSample + 10

		y -= th
		x := area.X + (area.W-tw)/2

		bg := colPanelBG

		gfx.RoundedBoxColor(r, x, y, x+tw, y+th, 10, bg)
		gfx.RoundedRectangleColor(r, x, y, x+tw, y+th, 10, colPanelBorder)

		drawText(t.msg, r, cache, TextOptions{
			rect:     sdl.Rect{X: x, Y: y, W: tw, H: th},
			font:     font,
			col:      colText,
			centered: true,
			clamped:  true,
		})

		y -= 6

	}

	toasts = make([]errorToast, 0)
}
