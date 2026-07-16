package main

import (
	"exnes/Core"
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type cheatMain struct {
	font *ttf.Font

	windowW int32
	windowH int32

	TitleRect sdl.Rect
	Title     string

	cheatsRect sdl.Rect
	cheats     []Core.Cheat
	HoverIndex int

	scrollOffsetTarget int32
	scrollOffsetCurr   float64

	cheatCache
}

type cheatCache struct {
	textcache map[string]textCache
	panelCache
}

type cheatEntry struct {
	cheat *Core.Cheat
}

func (main *cheatMain) Setup(font *ttf.Font) {

	main.font = font

	main.panelCache = make(panelCache)
	main.cheatCache.textcache = make(map[string]textCache)

	main.cheats = Core.CreateDemoCheats()

	main.Layout()
}

const (
	cheatPanelPadding = 10
	cheatItemHeight   = 35
	cheatListPadding  = 10

	ScrollSpeed = 10
	ScrollEase  = 0.25
)

func (main *cheatMain) Layout() {
	Y := int32(cheatPanelPadding)

	main.TitleRect = sdl.Rect{
		W: main.windowW - 20,
		H: 30,
		X: 10,
		Y: Y,
	}

	Y += main.TitleRect.Y + main.TitleRect.H + cheatPanelPadding

	main.cheatsRect = sdl.Rect{
		W: main.windowW - 20,
		H: main.windowH - Y - 10,
		Y: Y,
		X: 10,
	}
}

func (main *cheatMain) updateScroll() {
	diff := float64(main.scrollOffsetTarget) - main.scrollOffsetCurr

	if math.Abs(diff) < 0.5 {
		main.scrollOffsetCurr = float64(main.scrollOffsetTarget)
		return
	}

	main.scrollOffsetCurr += diff * ScrollEase
}

func (main *cheatMain) render(r *sdl.Renderer) {
	main.updateScroll()
	r.SetDrawColor(colBarBG.R, colBarBG.G, colAccent.B, 255)

	main.panelCache.drawRoundedRect(r, &main.TitleRect, colControlPanelBG, true)
	main.panelCache.drawRoundedRect(r, &main.TitleRect, colControlPanelBorder, false)

	main.panelCache.drawRoundedRect(r, &main.cheatsRect, colControlPanelBG, true)
	main.panelCache.drawRoundedRect(r, &main.cheatsRect, colControlPanelBorder, false)

	main.renderList(r)
}

func (main *cheatMain) renderList(r *sdl.Renderer) {
	r.SetClipRect(&main.cheatsRect)

	offset := int32(main.scrollOffsetCurr)

	firstIndex := max(offset/cheatItemHeight, 0)

	Y := main.cheatsRect.Y - (offset % (cheatItemHeight + cheatListPadding)) + 20

	for i := firstIndex; Y < main.cheatsRect.Y+main.cheatsRect.H && int(i) < len(main.cheats); i++ {
		// entry := main.cheats[i]

		rowRect := sdl.Rect{X: main.cheatsRect.X + 15, Y: Y, W: main.cheatsRect.W - 30, H: cheatItemHeight}

		main.drawRoundedRect(r, &rowRect, colFieldBG, true)
		main.drawRoundedRect(r, &rowRect, colAccent, false)

		drawText(main.cheats[i].Description, rowRect, r, 0, colText, main.textcache, main.font, true)

		Y += cheatItemHeight + cheatListPadding

	}

	r.SetClipRect(nil)

}

func (main *cheatMain) handleScroll(Y int32) {
	main.scrollOffsetTarget -= Y * ScrollSpeed

	contentHeight := int32(len(main.cheats)) * (cheatItemHeight + cheatListPadding)
	maxS := max(contentHeight-main.cheatsRect.H, 0)

	if main.scrollOffsetTarget < 0 {
		main.scrollOffsetTarget = 0
	}

	if main.scrollOffsetTarget > maxS {
		main.scrollOffsetTarget = maxS
	}
}
