package main

import (
	"fmt"

	"github.com/sqweek/dialog"
)

func openRom(console *game, state *localState) {
	filename, err := dialog.File().Filter("", ".nes").Title("Load Rom").Load()
	if err != nil {
		fmt.Println("error opening rom", err)
	}

	name, err := console.LoadRom(filename)
	if err != nil {
		fmt.Println("error loading rom", err)
	}

	state.RecentFiles = append(state.RecentFiles, recentRom{
		Name:     name,
		Location: filename,
	})

}
