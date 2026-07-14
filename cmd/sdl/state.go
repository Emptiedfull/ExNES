package main

import (
	"encoding/json"
	"exnes/Core"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type localState struct {
	RecentFiles []recentRom `json:"recent"`
	running     bool
	saves       []romSave
	Settings    state_setting `json:"settings"`
}

type state_setting struct {
	Show_fps      bool   `json:"show_fps"`
	Current_speed string `json:"current_speed"`

	Muted          bool   `json:"muted"`
	Current_volume string `json:"current_volume"`

	Inputs      Inputs `json:"controlBinds"`
	TurboInputs Inputs `json:"turboBinds"`
}

type romSave struct {
	snapshot  *Core.Snapshot
	timestamp string
	filled    bool
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

	s.saves = make([]romSave, 10)

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

func (inp Inputs) MarshalJSON() ([]byte, error) {
	res := make(map[string]string)

	for key, action := range inp.ActionToKey {
		res[key.String()] = sdl.GetScancodeName(action)
	}

	x, err := json.Marshal(res)
	if err != nil {
		fmt.Println("unable to marshal json:", err)
		return make([]byte, 0), err
	}

	return x, nil
}

func (inp *Inputs) UnmarshalJSON(data []byte) error {
	raw := make(map[string]string)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	inp.KeyToAction = make(map[sdl.Scancode]Core.BUTTON)
	inp.ActionToKey = make(map[Core.BUTTON]sdl.Scancode)
	for buttonName, scancodeName := range raw {
		scancode := sdl.GetScancodeFromName(scancodeName)
		button := Core.GetActionByName(buttonName)

		inp.KeyToAction[scancode] = button
		inp.ActionToKey[button] = scancode
	}

	return nil
}
