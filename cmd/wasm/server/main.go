// go: build dev
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const distDir = "./dist"
const staticDir = "./static"

var globalCache = newCache(true)
var games = make([]GameInfo, 0)

var mimeList = map[string]string{
	".html": "text/html",
	".js":   "text/javascript",
	".css":  "text/css",
	".wasm": "application/wasm",
	".png":  "image/png",
	".svg":  "image/svg+xml",
}

var compressibleList = map[string]bool{
	".wasm": true,
	".html": true,
	".js":   true,
	".css":  true,
	".png":  false,
	".nes":  true,
}

func main() {
	noWatch := flag.Bool("no-watch", false, "disable the file watcher")
	flag.Parse()

	PrepareGameList()
	err, _ := compileWasm("../wasm.go", "./static/nes.wasm")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Wasm compiled")

	restartChan := make(chan int)

	if *noWatch {

		StartServer(restartChan)
	} else {
		go StartServer(restartChan)
		startWatcher(restartChan)
	}
}
func StartServer(restartChan chan int) {

	fmt.Println("building assests...")
	bundleOutput()
	fmt.Println("completed build")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleStatic)
	mux.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {

		err := json.NewEncoder(w).Encode(games)
		if err != nil {
			http.Error(w, "bad response", http.StatusBadRequest)
			return
		}

	})

	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(PrepareDocList())
		if err != nil {
			http.Error(w, "FUCKF UHRQIORHIOQHRIOQ", http.StatusBadRequest)
			return
		}
	})

	server := &http.Server{Handler: mux, Addr: ":8070"}

	go func() {
		for x := range restartChan {

			switch x {
			case 1:
				fmt.Print("building assests...")
				bundleOutput()
				globalCache.clear()
				fmt.Println(" completed build")
			case 2:
				fmt.Print("compiling wasm...")
				compileWasm("../wasm.go", "./static/nes.wasm")
				delete(globalCache.data, "../wasm.go")
				globalCache.get("../wasm.go")
				fmt.Println("completed")

			}

		}
	}()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("unable to start server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	fmt.Println("quitting server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	url := filepath.Clean(r.URL.Path)

	if url == "/" || url == "" {
		url = "/index.html"
	}

	accepts := r.Header.Get("Accept-Encoding")
	extension := strings.ToLower(filepath.Ext(url))

	path := filepath.Join(staticDir, url)

	_, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	w.Header().Set("Content-type", mimeList[extension])

	if strings.Contains(accepts, "br") && compressibleList[extension] {
		handleCompressibles(w, r, path)
	} else {
		file, _ := os.Open(path)
		io.Copy(w, file)
		defer file.Close()
	}

}

func handleCompressibles(w http.ResponseWriter, r *http.Request, path string) {
	data := globalCache.get(path)

	w.Header().Set("Content-Encoding", "br")
	w.Header().Set("vary", "Accept-Encoding")

	w.Write(data)
}

type GameInfo struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func PrepareGameList() {

	err := filepath.WalkDir("./static/games", func(path string, d fs.DirEntry, err error) error {

		s, found := strings.CutPrefix(path, "static/games/")
		globalCache.get(path)

		if found {
			s2, found := strings.CutSuffix(s, ".nes")
			if found {

				g := GameInfo{
					ID:   s2,
					Name: strings.ToTitle(strings.ReplaceAll(s2, "_", " ")),
				}

				games = append(games, g)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

}

func PrepareDocList() []string {
	list := make([]string, 0)

	err := filepath.WalkDir("./static/docs", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)

		s, found := strings.CutPrefix(path, "static/docs/")
		if !found {
			return nil
		}

		list = append(list, s)

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return list

}
