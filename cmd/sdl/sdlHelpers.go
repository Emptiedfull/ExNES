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
			res[i].label = "filled"
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
			option.label = "filled"
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
		// mb.menuFlags[saveAvailable] = true
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
