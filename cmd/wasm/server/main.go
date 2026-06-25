package main

import (
	"context"
	"encoding/json"
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

var mimeList = map[string]string{
	".html": "text/html",
	".js":   "text/javascript",
	".css":  "text/css",
	".wasm": "application/wasm",
	".png":  "image/png",
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

	compileWasm("../main.go", "./static/nes.wasm")

	restartChan := make(chan bool)

	go StartServer(restartChan)

	startWatcher(restartChan)
}

func StartServer(restartChan chan bool) {

	fmt.Println("building assests...")
	bundleOutput()
	fmt.Println("completed build")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleStatic)
	mux.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		games := PrepareGameList()

		err := json.NewEncoder(w).Encode(games)
		if err != nil {
			http.Error(w, "bad response", http.StatusBadRequest)
			return
		}

	})

	server := &http.Server{Handler: mux, Addr: ":8070"}

	go func() {
		for range restartChan {

			fmt.Println("building assests...")
			bundleOutput()
			globalCache.clear()
			fmt.Println("completed build")

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

func PrepareGameList() []GameInfo {
	result := make([]GameInfo, 0)
	err := filepath.WalkDir("./static/games", func(path string, d fs.DirEntry, err error) error {

		// if d.IsDir() {
		// 	fmt.Println(path)
		// 	return fmt.Errorf("FUCK FUCK FUCK THERES A DIRECTORY IN THE GAMES THE WORLD IS ENDING")
		// }

		s, found := strings.CutPrefix(path, "static/games/")

		if found {
			s2, found := strings.CutSuffix(s, ".nes")
			if found {

				g := GameInfo{
					ID:   s2,
					Name: strings.ToTitle(strings.ReplaceAll(s2, "_", " ")),
				}

				result = append(result, g)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return result
}
