package main

import (
	"fmt"
	"log"
	"time"

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

func getRecentItems(roms []recentRom, console *game, mb *menuBar) []expandableOption {
	res := make([]expandableOption, len(roms))

	for i, rom := range roms {
		res[i] = expandableOption{
			label: rom.Name,
			onClick: func() {

				console.LoadRom(rom.Location, mb)
				mb.resetMenu()
			},
		}
	}

	return res
}

func (mb *menuBar) getSaveStateItems() []expandableOption {
	res := make([]expandableOption, len(mb.state.saves))

	for i, save := range mb.state.saves {
		if !save.filled {
			res[i].label = "empty"

		} else {
			res[i].label = save.timestamp
		}

		res[i].onClick = func() {

			mb.console.clickSnapshot(mb.state, i)
			mb.state.saves[i].filled = true
			mb.updateSavesMenus()

		}
	}

	return res
}

func (mb *menuBar) getLoadItems() []expandableOption {
	res := make([]expandableOption, 0)

	for i, save := range mb.state.saves {
		if save.filled {
			option := expandableOption{}
			option.label = save.timestamp
			option.onClick = func() {
				mb.console.loadSnapshot(mb.state, i)
			}
			res = append(res, option)
		}

	}

	return res
}

func (mb *menuBar) updateSavesMenus() {
	mb.Items[0].options[3].ExpandableItems = mb.getSaveStateItems()
	loadItems := mb.getLoadItems()
	mb.Items[0].options[4].ExpandableItems = loadItems

	if len(loadItems) > 0 {

		mb.setFlag(saveAvailable, true)
	}

	mb.positionLayout()
}

func (mb *menuBar) resetMenu() {
	mb.hoverIndex = -1
	mb.optionHoverIndex = -1
	mb.dropdownIndex = -1
	mb.subOptionHoverIndex = -1

}

func getSaveTimeStamp() string {
	now := time.Now()
	formatted := now.Format("01-02 3:04")

	return formatted
}

func convertColToString(col sdl.Color) string {
	return fmt.Sprintf("R: %v G: %v B:%v", col.R, col.G, col.B)
}

func (mb *menuBar) setupFlags() {
	for i := range mb.Items {
		for j := range mb.Items[i].options {
			option := &mb.Items[i].options[j]
			if option.affectedFlag != none {
				mb.affected[option.affectedFlag] = append(mb.affected[option.affectedFlag], option)
			}
		}
	}
}

func (mb *menuBar) updateMenuState() {

	if !mb.flagsUpdated {
		return
	}

	for flag, affected := range mb.affected {
		for _, option := range affected {
			option.enabled = mb.menuFlags[flag]
		}
	}
}

func (mb *menuBar) setFlag(flag Menuflag, state bool) {
	mb.flagsUpdated = true
	mb.menuFlags[flag] = state

}

func (m *menuBar) drawText(text string, rect sdl.Rect, r *sdl.Renderer, offet int32, col sdl.Color) {
	itemEntry := text + convertColToString(col)
	entry, ok := m.cache.textCache[itemEntry]
	if !ok {
		entry = textCache{}
		surface, err := m.Font.RenderUTF8Blended(text, col)
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
		m.cache.textCache[itemEntry] = entry

	}

	dst := sdl.Rect{
		X: rect.X + offet,
		Y: rect.Y + (rect.H-entry.H)/2,
		W: entry.W,
		H: entry.H,
	}

	r.Copy(entry.texture, nil, &dst)

}

const upfactor = 8

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
