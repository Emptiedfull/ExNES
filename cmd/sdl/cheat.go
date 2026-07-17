package main

import (
	"exnes/Core"
	"fmt"
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

	ThumbRect       sdl.Rect
	TrackRect       sdl.Rect
	DraggingThumb   bool
	dragStartMouseY int32
	dragStartOffset int32

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

	ScrollSpeed    = 20
	ScrollEase     = 0.25
	ScrollBarWidth = 4
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
		main.updateScrollBar()
		return
	}

	main.scrollOffsetCurr += diff * ScrollEase

	main.updateScrollBar()
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

	const rowSpan = cheatItemHeight + cheatListPadding

	offset := int32(main.scrollOffsetCurr)

	firstIndex := max(offset/rowSpan, 0)

	Y := main.cheatsRect.Y - (offset % rowSpan) + 20

	for i := firstIndex; Y < main.cheatsRect.Y+main.cheatsRect.H && int(i) < len(main.cheats); i++ {
		rowRect := sdl.Rect{X: main.cheatsRect.X + 10, Y: Y, W: main.cheatsRect.W - 30, H: cheatItemHeight}

		main.drawRoundedRect(r, &rowRect, colFieldBG, true)
		main.drawRoundedRect(r, &rowRect, colAccent, false)

		drawText(main.cheats[i].Description, r, main.textcache, TextOptions{rect: rowRect, col: colText, font: main.font, clamped: true, centered: false, offset: 32})

		Y += rowSpan
	}

	r.SetClipRect(nil)

	r.SetDrawColor(colAccent.R, colAccent.G, colAccent.B, 180)
	r.FillRect(&main.ThumbRect)

}

func (main *cheatMain) updateScrollBar() {
	rowSpan := int32(cheatItemHeight + cheatListPadding)
	contentH := int32(len(main.cheats)) * rowSpan
	viewpoint := main.cheatsRect.H

	if contentH < viewpoint {
		return
	}

	maxS := contentH - viewpoint

	X := main.cheatsRect.X + main.cheatsRect.W - ScrollBarWidth - 6
	Y := main.cheatsRect.Y + 6
	H := main.cheatsRect.H - 12

	thumbH := int32(float64(H) * float64(viewpoint) / float64(contentH))
	if thumbH < 40 {
		thumbH = 40
	}

	progress := main.scrollOffsetCurr / float64(maxS)
	thumbY := Y + int32(progress*float64(H-thumbH))

	main.TrackRect = sdl.Rect{
		X: X - 4,
		Y: Y,
		W: ScrollBarWidth + 8,
		H: H,
	}

	thumbRect := sdl.Rect{
		X: X,
		Y: thumbY,
		W: ScrollBarWidth,
		H: thumbH,
	}

	main.ThumbRect = thumbRect

}

func (main *cheatMain) handleScroll(Y int32) {
	main.scrollOffsetTarget -= Y * ScrollSpeed
	main.clampScroll()
}

func (main *cheatMain) handleMouseDown(x, y int32) {
	fmt.Println("down")
	if pointInRect(main.ThumbRect, x, y) {
		main.DraggingThumb = true
		main.dragStartOffset = main.scrollOffsetTarget
		main.dragStartMouseY = y

		return
	}

	if pointInRect(main.TrackRect, x, y) {

		pageSize := main.cheatsRect.H

		if y < main.ThumbRect.Y {
			main.scrollOffsetTarget -= pageSize
		} else {
			main.scrollOffsetTarget += pageSize
		}
		main.clampScroll()
	}

}

func (main *cheatMain) handleMouseUp() {

	main.DraggingThumb = false
}

func (main *cheatMain) handleMouseMove(x, y int32) {
	if !main.DraggingThumb {
		return
	}

	rowSpan := int32(cheatItemHeight + cheatListPadding)
	contentH := int32(len(main.cheats)) * rowSpan

	H := main.TrackRect.H - main.ThumbRect.H

	if H < 0 {
		return
	}

	diffM := y - main.dragStartMouseY
	diffS := int32(float64(diffM) * float64(contentH-main.cheatsRect.H) / float64(H))

	main.scrollOffsetTarget = main.dragStartOffset + diffS
	main.clampScroll()
}

func (main *cheatMain) clampScroll() {
	contentHeight := int32(len(main.cheats)) * (cheatItemHeight + cheatListPadding)
	maxS := max(contentHeight-main.cheatsRect.H, 0)

	if main.scrollOffsetTarget < 0 {
		main.scrollOffsetTarget = 0
	}

	if main.scrollOffsetTarget > maxS {
		main.scrollOffsetTarget = maxS
	}
}
