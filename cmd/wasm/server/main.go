package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/evanw/esbuild/pkg/api"
)

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

	api.Build(api.BuildOptions{
		EntryPoints: []string{
			"static/scripts/nes.js",
		},
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Splitting:         true,
		Format:            api.FormatESModule,
		Outdir:            "static/dist",
		// Outfile:           "app.js",
		Write: true,
	})

	api.Build(api.BuildOptions{
		EntryPoints: []string{
			"static/scripts/emuWorker.js",
		},
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Outdir:            "static/dist",
		Write:             true,
	})

	results := api.Build(api.BuildOptions{
		EntryPoints: []string{
			"static/scripts/driverWorklet.js",
		},
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Outdir:            "static/dist",
		Write:             true,
	})

	for _, err := range results.Errors {
		fmt.Println(err)
	}

	res := api.Build(api.BuildOptions{
		EntryPoints: []string{
			"static/styles/all.css",
		},
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Outdir:            "static/dist",
		Write:             true,
	})

	for _, err := range res.Errors {
		fmt.Println(err)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {

	extension := filepath.Ext(r.URL.RawPath)
	fmt.Println(extension)

}

func compileWasm(src, dst string) (error, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "tinygo", "build", "-o", dst, "-target", "wasm", "-opt", "z", src)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error compiling wasm: %w", err), "0"
	}

	info, err := os.Stat(dst)
	if err != nil {
		return fmt.Errorf("Error finding compiled file: %w", err), "0"
	}

	sizeStr := formatSize(info.Size())

	return nil, sizeStr
}

func formatSize(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0

	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
