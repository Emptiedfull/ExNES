package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
}

func main() {
	fmt.Print("building... ")
	bundleOutput()
	fmt.Println("completed build")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleStatic)

	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		fmt.Println("ERROR STARTING THE FUCKASS SRERVER")
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

	// path := filepath.Join(staticDir, url)
	// data := globalCache.get(path)

	// w.Write(data)
}

func handleCompressibles(w http.ResponseWriter, r *http.Request, path string) {
	data := globalCache.get(path)

	w.Header().Set("Content-Encoding", "br")
	w.Header().Set("vary", "Accept-Encoding")

	w.Write(data)
}
