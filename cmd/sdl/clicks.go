package main

import (
	"fmt"
	"path/filepath"

	"github.com/sqweek/dialog"
)

func openRom(console *game, state *localState) {
	filename, err := dialog.File().Filter("", ".nes").Title("Load Rom").Load()
	if err != nil {
		fmt.Println("error opening rom", err)
	}

	err = console.LoadRom(filename)
	if err != nil {
		fmt.Println("error loading rom", err)
	}

	if state != nil {

		state.addRecentRom(recentRom{
			Name:     filepath.Base(filename),
			Location: filename,
		})
	}

}
