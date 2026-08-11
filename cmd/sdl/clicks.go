package main

import (
	"fmt"
	"path/filepath"

	"github.com/sqweek/dialog"
)

func openRom(console *game, state *localState, mb *menuBar) {
	filename, err := dialog.File().Filter("", ".nes").Title("Load Rom").Load()
	if err != nil {
		fmt.Println("error opening rom", err)
	}

	err = console.LoadRom(filename, mb)
	if err != nil {
		fmt.Println("error loading rom", err)
	}

	if state != nil {

		state.addRecentRom(recentRom{
			Name:     filepath.Base(filename),
			Location: filename,
		})

		fmt.Println(getRecentItems(state.RecentFiles, mb.console, mb))

		mb.Items[0].options[2].ExpandableItems = getRecentItems(state.RecentFiles, mb.console, mb)
		mb.positionLayout()
	}

}

func (console *game) clickSnapshot(save *romSave) {

	save.timestamp = getSaveTimeStamp()
	save.snapshot = console.core.SaveState()

}

func (console *game) loadSnapshot(save romSave) {
	console.core.LoadSnapshot(*save.snapshot)

}

func (console *game) PauseGame() {
	console.pauseChannel <- true
}

func (console *game) UnPauseGame() {
	console.pauseChannel <- false
}
