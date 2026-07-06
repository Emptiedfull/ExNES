package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type localState struct {
	RecentFiles []string `json:"recent"`
}

func newState() *localState {
	return &localState{
		RecentFiles: make([]string, 0),
	}
}

func loadState() *localState {
	s := &localState{}

	data, err := os.ReadFile("state.json")
	if err != nil {
		fmt.Println("error loading state:", err)
		return newState()
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		fmt.Println("error unmarshaling state:", err)
		return newState()
	}

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
