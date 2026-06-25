package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func startWatcher(restartChan chan bool) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		fmt.Println("error starting directory watcher")
		return
	}

	defer watcher.Close()

	time.Sleep(1 * time.Second)

	addPaths(watcher, "./")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				fmt.Println("detected changes on:", event.Name)
				restartChan <- true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			fmt.Println("error", err)
		}
	}
}

func addPaths(watcher *fsnotify.Watcher, directory string) error {
	err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.Contains(path, "dist") {
				return nil
			}
			return watcher.Add(path)
		}

		return nil
	})

	return err
}
