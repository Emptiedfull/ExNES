package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
)

type cache struct {
	data map[string][]byte
	mux  sync.RWMutex
}

func main() {
	err, size := compileWasm("../", "./static/nes.wasm")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wasm compiled succesfully: %v", size)

	err = compressFile("./static/nes.wasm", "./compressed/nes.wasm.br")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(getFileSize("./compressed/nes.wasm.br"))

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

var writerPool = sync.Pool{
	New: func() interface{} {
		return brotli.NewWriterLevel(io.Discard, 11)
	},
}

func compressFile(src, dst string) error {
	fmt.Println("compressing: ", src)

	input, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Error oepning src file:  %w", err)
	}

	output, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("Error creating dst file: %w", err)
	}

	writer := writerPool.Get().(*brotli.Writer)

	writer.Reset(output)
	defer writer.Close()

	_, err = io.Copy(writer, input)
	if err != nil {
		return fmt.Errorf("Error compressing file: %v", err)
	}
	// fmt.Println(getFileSize(src))
	// fmt.Println(getFileSize(dst))
	return nil
}

func getFileSize(src string) string {
	info, err := os.Stat(src)
	fmt.Println(info.Name())
	if err != nil {
		fmt.Println("error reading file")
	}

	return formatSize(info.Size())
}
