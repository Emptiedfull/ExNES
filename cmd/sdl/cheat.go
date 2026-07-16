package main

import (
	"exnes/Core"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type cheatMain struct {
	font *ttf.Font

	TitleRect sdl.Rect
	Title     string

	cheatsRect   sdl.Rect
	cheats       []*cheatEntry
	scrollOffset int32
	HoverIndex   int

	cheatCache
}

type cheatCache struct {
	textCache
}

type cheatEntry struct {
	cheat *Core.Cheat
}

func (cheatMain) Setup(font *ttf.Font) {

}
