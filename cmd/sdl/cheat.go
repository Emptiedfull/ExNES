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
	iconCache map[string]*sdl.Texture
}

type cheatEntry struct {
	cheat *Core.Cheat
}

func (main *cheatMain) Setup(font *ttf.Font) {

	main.font = font

	main.panelCache = make(panelCache)
	main.cheatCache.textcache = make(map[string]textCache)
	main.iconCache = make(map[string]*sdl.Texture)

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

type cheatRowColScheme struct {
	bg     sdl.Color
	border sdl.Color
	text   sdl.Color
	icon   sdl.Color
}

var (
	CheatNormal = cheatRowColScheme{
		bg:     colFieldBG,
		border: colPanelBorder,
		text:   colTextDim,
		icon:   colTextDim,
	}

	cheatHover = cheatRowColScheme{
		bg:     colHover,
		border: colAccent,
		text:   colText,
		icon:   colText,
	}

	cheatActive = cheatRowColScheme{
		bg:     sdl.Color{R: 46, G: 40, B: 74, A: 255},
		border: colAccent,
		text:   colText,
		icon:   colAccent,
	}

	cheatActiveHover = cheatRowColScheme{
		bg:     sdl.Color{R: 58, G: 50, B: 92, A: 255},
		border: colAccent,
		text:   colText,
		icon:   colAccent,
	}
)

func get_Scheme(active, hover bool) cheatRowColScheme {
	switch {
	case active && hover:
		return cheatActiveHover
	case active:
		return cheatActive
	case hover:
		return cheatHover
	default:
		return CheatNormal
	}
}

func (main *cheatMain) renderList(r *sdl.Renderer) {
	r.SetClipRect(&main.cheatsRect)

	const rowSpan = cheatItemHeight + cheatListPadding

	offset := int32(main.scrollOffsetCurr)

	firstIndex := max(offset/rowSpan, 0)

	Y := main.cheatsRect.Y - (offset % rowSpan) + 20

	for i := firstIndex; Y < main.cheatsRect.Y+main.cheatsRect.H && int(i) < len(main.cheats); i++ {
		cheat := main.cheats[i]
		rowRect := sdl.Rect{X: main.cheatsRect.X + 10, Y: Y, W: main.cheatsRect.W - 30, H: cheatItemHeight}
		innerRect := sdl.Rect{X: rowRect.X + 1, Y: rowRect.Y + 1, W: rowRect.W - 2, H: rowRect.H - 2}

		hovered := int32(main.HoverIndex) == i

		scheme := get_Scheme(cheat.Enabled, hovered)

		main.drawRoundedRect(r, &rowRect, scheme.border, true)
		main.drawRoundedRect(r, &innerRect, scheme.bg, true)

		checkboxRect := sdl.Rect{
			X: rowRect.X + 6,
			H: cheatItemHeight - 16,
			Y: rowRect.Y + 8,
			W: cheatItemHeight - 16,
		}

		if cheat.Enabled {
			drawIcon(r, "./icons/check.svg", scheme.icon, main.iconCache, &checkboxRect)
		}

		r.SetDrawColor(scheme.border.R, scheme.border.G, scheme.border.B, 255)
		r.DrawLine(rowRect.X+32, checkboxRect.Y-2, rowRect.X+32, checkboxRect.Y+checkboxRect.H+2)

		drawText(cheat.Description, r, main.textcache, TextOptions{rect: rowRect, col: scheme.text, font: main.font, clamped: true, centered: false, offset: 44})

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

	if idx := main.getIdx(x, y); idx != -1 {

		main.cheats[idx].Enabled = !main.cheats[idx].Enabled
	}

}

func (main *cheatMain) handleMouseUp() {

	main.DraggingThumb = false
}

func (main *cheatMain) handleMouseMove(x, y int32) {
	if main.DraggingThumb {

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

	main.HoverIndex = main.getIdx(x, y)

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

func (main *cheatMain) getIdx(x, y int32) int {
	if !pointInRect(main.cheatsRect, x, y) {
		return -1
	}

	const rowSpan = cheatItemHeight + cheatListPadding
	offset := int32(main.scrollOffsetCurr)

	r := y - main.cheatsRect.Y - 20 + offset
	if r < 0 {
		return -1
	}

	idx := r / rowSpan
	inRow := r % rowSpan

	if inRow >= cheatItemHeight {
		return -1
	}

	if int(idx) >= len(main.cheats) {
		return -1
	}

	return int(idx)

}
