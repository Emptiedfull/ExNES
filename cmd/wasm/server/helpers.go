package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/evanw/esbuild/pkg/api"
)

func bundleOutput() {
	res := api.Build(api.BuildOptions{
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

	for _, err := range res.Errors {
		log.Fatalln(err)
	}

	res = api.Build(api.BuildOptions{
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

	for _, err := range res.Errors {
		log.Fatalln(err)
	}

	res = api.Build(api.BuildOptions{
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
	for _, err := range res.Errors {
		log.Fatalln(err)
	}

	res = api.Build(api.BuildOptions{
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
