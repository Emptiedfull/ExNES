package main

import (
	"fmt"
	"net/http"
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
}

func main() {
	// err, size := compileWasm("../", "./static/nes.wasm")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Wasm compiled succesfully: %v", size)

	// err = CompressToFile("./static/nes.wasm", "./compressed/nes.wasm.br")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	bundleOutput()

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

	w.Header().Set("Content-type", mimeList[extension])

	if extension == ".png" {
		return
	}

	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

	if strings.Contains(accepts, "br") {
		handleCompressibles(w, r, path)
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
