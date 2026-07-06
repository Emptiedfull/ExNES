package main

import (
	"encoding/json"
	"exnes/Core"
	"fmt"
	"os"
)

type localState struct {
	RecentFiles []recentRom `json:"recent"`
	running     bool
	snapshots   []*Core.Snapshot
}

func newState() *localState {
	return &localState{
		RecentFiles: make([]recentRom, 0),
	}
}

type recentRom struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func loadState() *localState {
	s := &localState{}

	data, err := os.ReadFile("state.json")
	if err != nil {
		fmt.Println("error loading state:", err)

	}

	err = json.Unmarshal(data, s)
	if err != nil {
		fmt.Println("error unmarshaling state:", err)

	}

	s.snapshots = make([]*Core.Snapshot, 10)

	return s
}

func (state *localState) saveState() {
	data, err := json.MarshalIndent(state, "", "")
	if err != nil {
		fmt.Println("error marshaling save state", err)
		return
	}

	os.WriteFile("state.json", data, 0644)
}

func (state *localState) addRecentRom(new recentRom) {
	for _, rom := range state.RecentFiles {
		if rom.Location == new.Location {
			return
		}
	}

	state.RecentFiles = append(state.RecentFiles, new)
}
