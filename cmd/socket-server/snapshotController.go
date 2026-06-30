//go:build !js && !wasm

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type snapshotInfo struct {
	Index    int `json:"Index"`
	Frame_no int `json:"Frame_no"`
	Cycles   int `json:"Cycles"`
}

func loadSnapshot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	SnapIdx := r.URL.Query().Get("index")
	if SnapIdx == "" {
		SnapIdx = "0"
	}

	index, err := strconv.Atoi(SnapIdx)
	if err != nil {
		http.Error(w, "unable to load snapshot", http.StatusInternalServerError)
		return
	}

	debugConsole.Console.LoadSnapshot(debugConsole.Console.Snapshots.Data[index])

	w.WriteHeader(http.StatusOK)

}

func fetchSnapshots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	snapshots := make([]snapshotInfo, len(debugConsole.Console.Snapshots.Data))

	for i, snap := range debugConsole.Console.Snapshots.Data {
		snapshots[i].Frame_no = snap.Frame_no
		snapshots[i].Index = i
		snapshots[i].Cycles = snap.Cycles
	}

	err := json.NewEncoder(w).Encode(snapshots)
	if err != nil {
		fmt.Println("error parsing snapshots")
		http.Error(w, "something wrong", http.StatusInternalServerError)
	}
}
