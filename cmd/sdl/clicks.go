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
	}

}

func (console *game) clickSnapshot(state *localState, index int) {

	state.saves[index].snapshot = console.core.SaveState()
	state.saves[index].timestamp = getSaveTimeStamp()

}

func (console *game) loadSnapshot(state *localState, index int) {
	console.core.LoadSnapshot(*state.saves[index].snapshot)

}

func (console *game) PauseGame() {
	console.pauseChannel <- true
}

func (console *game) UnPauseGame() {
	console.pauseChannel <- false
}
