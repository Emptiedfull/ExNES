package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

func (mb *menuBar) getArrow(r *sdl.Renderer, col sdl.Color) *sdl.Texture {
	if mb.cache.arrowCache[col] != nil {
		return mb.cache.arrowCache[col]
	}

	arrow, err := createArrow(r, 16, 16, col)
	if err != nil {
		fmt.Println("soemthing fuckass:", err)
		return nil
	}

	mb.cache.arrowCache[col] = arrow
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
			enabled: true,
			label:   rom.Name,
			onClick: func() {
				if mb.menuFlags[gameRunning] {
					console.romPath = rom.Location
					console.reloadROM()
					return
				}

				console.LoadRom(rom.Location, mb)
				mb.resetMenu()
			},
		}
	}

	return res
}

// func drawRoundedRect(r *sdl.Renderer, rect sdl.Rect, col sdl.Color, filled bool) {
// 	if filled {
// 		gfx.RoundedBoxColor(r, rect.X, rect.Y, rect.X+rect.W, rect.Y+rect.H, 8, col)
// 	} else {
// 		gfx.RoundedRectangleColor(r, rect.X, rect.Y, rect.X+rect.W, rect.Y+rect.H, 8, col)
// 	}
// }

func (mb *menuBar) getSaveStateItems() []expandableOption {
	res := make([]expandableOption, len(mb.state.saves))

	for i, save := range mb.state.saves {
		res[i].enabled = true
		if !save.filled {
			res[i].label = "empty"
			res[i].Icon = "./icons/save-plus.svg"

		} else {
			res[i].label = save.timestamp
			res[i].Icon = "./icons/save-check.svg"
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
			option.enabled = true
			option.Icon = "./icons/save-check.svg"
			option.onClick = func() {
				mb.console.loadSnapshot(mb.state, i)
			}
			res = append(res, option)
		}

	}

	return res
}

func (mb *menuBar) setupMenus() {
	mb.updateSettingsMenu()
	mb.updateSoundMenu()
}

func (mb *menuBar) updateSettingsMenu() {
	state := mb.state.Settings

	fpsOption := &mb.Items[2].options[0].ExpandableItems[0]

	speedOptions := mb.Items[2].options[0].ExpandableItems[2:6]

	for i := range speedOptions {

		if speedOptions[i].label == state.Current_speed {
			speedOptions[i].Icon = "./icons/check.svg"

			label := speedOptions[i].label
			speedPerc, _ := strconv.Atoi(label[0 : len(label)-1])

			mb.console.changeSpeed(float64(speedPerc) / 100.0)

		} else {
			speedOptions[i].Icon = ""
		}

		speedOptions[i].onClick = func() {

			mb.state.Settings.Current_speed = speedOptions[i].label
			mb.updateSettingsMenu()
		}
	}

	if state.Show_fps {
		fpsOption.Icon = "./icons/check.svg"
	} else {
		fpsOption.Icon = ""
	}

	mb.positionLayout()
}

func (mb *menuBar) updateSoundMenu() {
	state := mb.state.Settings
	muteOption := &mb.Items[2].options[1].ExpandableItems[0]
	volumeOptions := mb.Items[2].options[1].ExpandableItems[2:6]

	if state.Muted {
		mb.console.changeVolume("0%")
		muteOption.Icon = "./icons/check.svg"
		muteOption.onClick = func() {
			mb.state.Settings.Muted = false
			mb.updateSoundMenu()
		}
	} else {
		muteOption.Icon = ""
		muteOption.onClick = func() {
			mb.state.Settings.Muted = true
			mb.updateSoundMenu()
		}
		mb.console.changeVolume(state.Current_volume)
	}

	for i := range volumeOptions {
		label := volumeOptions[i].label

		if !state.Muted {
			volumeOptions[i].enabled = true
		} else {
			volumeOptions[i].enabled = false
		}

		if label == state.Current_volume {
			volumeOptions[i].Icon = "./icons/check.svg"
		} else {
			volumeOptions[i].Icon = ""
		}

		volumeOptions[i].onClick = func() {
			mb.state.Settings.Current_volume = label
			mb.console.changeVolume(label)
			mb.updateSoundMenu()
		}

	}

	mb.positionLayout()
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

type TextOptions struct {
	rect     sdl.Rect
	offset   int32
	col      sdl.Color
	centered bool
	font     *ttf.Font
	clamped  bool
}

func drawText(text string, r *sdl.Renderer, cache map[string]textCache, options TextOptions) {
	var itemEntry string
	if options.clamped {
		itemEntry = text + convertColToString(options.col) + strconv.Itoa(int(options.rect.W))
	} else {
		itemEntry = text + convertColToString(options.col)
	}

	entry, ok := cache[itemEntry]
	if !ok {
		entry = textCache{}

		if options.clamped {
			text = truncateStr(options.font, text, int32(float32(options.rect.W)/1.2)-options.offset)
		}

		surface, err := options.font.RenderUTF8Blended(text, options.col)
		if err != nil {
			panic(err)
		}
		defer surface.Free()

		entry.W = int32(surface.W / 2)
		entry.H = int32(surface.H / 2)

		texture, err := r.CreateTextureFromSurface(surface)
		if err != nil {
			log.Fatal(err)
		}

		entry.texture = texture
		cache[itemEntry] = entry

	}

	x := options.rect.X + options.offset
	if options.centered {
		x = options.rect.X + (options.rect.W-entry.W)/2
	}

	dst := sdl.Rect{
		X: x,
		Y: options.rect.Y + (options.rect.H-entry.H)/2,
		W: entry.W,
		H: entry.H,
	}

	r.Copy(entry.texture, nil, &dst)

}

const filler = "..."

func truncateStr(font *ttf.Font, text string, MaxWidth int32) string {
	if textWdith(text, font) < MaxWidth {
		return text
	}

	fillerWidth := textWdith(filler, font)

	runes := []rune(text)
	for i := len(text) - 1; i > 0; i-- {
		check := string(runes[:i])
		if textWdith(check, font)+fillerWidth < MaxWidth {
			return check + filler
		}
	}

	return filler

}

func textWdith(text string, font *ttf.Font) int32 {
	w, _, err := font.SizeUTF8(text)
	if err != nil {
		fmt.Println("bad fucking text: ", err)
	}

	return int32(w) / upfactor
}

const upfactor = 2

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

func convertRectName(col sdl.Color, w, h int32, filled bool) string {

	var buf strings.Builder

	buf.WriteString(convertColToString(col))
	buf.WriteString(fmt.Sprintf("W:%v,h:%v", w, h))
	if filled {
		buf.WriteString("T")
	} else {
		buf.WriteString("F")
	}
	return buf.String()

}

func (p panelCache) drawRoundedRect(r *sdl.Renderer, rect *sdl.Rect, col sdl.Color, filled bool) {
	texture := p.getRoundedRect(r, col, rect.W, rect.H, filled)
	r.Copy(texture, nil, rect)
}

func (p panelCache) getRoundedRect(r *sdl.Renderer, col sdl.Color, w, h int32, filled bool) *sdl.Texture {
	key := convertRectName(col, w, h, filled)

	entry, ok := p[key]
	if ok {
		return entry

	}
	entry, err := createRoundedRect(r, col, w, h, filled)
	if err != nil {
		fmt.Println("no tex:", err)
		return nil
	}

	p[key] = entry

	return entry
}

func createRoundedRect(r *sdl.Renderer, col sdl.Color, w, h int32, filled bool) (*sdl.Texture, error) {
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

	bigR := int32(8 * upfactor)

	if filled {
		gfx.RoundedBoxRGBA(r, 0, 0, bigW-2, bigH-2, bigR, col.R, col.G, col.B, col.A)
	} else {
		gfx.RoundedRectangleColor(r, 0, 0, bigW-2, bigH-2, bigR, col)
	}

	r.SetRenderTarget(nil)

	return tex, nil
}

func (mb *menuBar) getIcon(r *sdl.Renderer, path string, col sdl.Color) *sdl.Texture {
	texture, ok := mb.cache.iconCache[path+convertColToString(col)]
	if ok {

		return texture
	}

	texture, err := img.LoadTexture(r, path)
	texture.SetColorMod(col.R, col.G, col.B)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	mb.cache.iconCache[path+convertColToString(col)] = texture
	return texture

}

func pointInRect(rect sdl.Rect, x, y int32) bool {
	return x >= rect.X && x <= rect.X+rect.W && y >= rect.Y && y <= rect.Y+rect.H
}
