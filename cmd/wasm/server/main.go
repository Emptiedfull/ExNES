package main

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func handleStatic(w http.ResponseWriter, r *http.Request) {

	extension := filepath.Ext(r.URL.RawPath)
	fmt.Println(extension)

}
