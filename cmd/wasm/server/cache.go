package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/andybalholm/brotli"
)

type cache struct {
	data    map[string][]byte
	mux     sync.RWMutex
	logging bool
}

func newCache(log bool) cache {
	return cache{
		data:    make(map[string][]byte),
		logging: log,
	}
}

func (c *cache) clear() {
	c.mux.Lock()
	defer c.mux.Unlock()

	for key := range c.data {
		fmt.Println(key)
	}
	wasmEntry := c.data["static/nes.wasm"]
	c.data = make(map[string][]byte)
	c.data["static/nes.wasm"] = wasmEntry
}

func (c *cache) get(file string) []byte {

	if val, ok := c.data[file]; ok {
		return val
	} else {

		c.CompressToMemory(file)

		return c.data[file]
	}
}

var writerPool = sync.Pool{
	New: func() interface{} {
		return brotli.NewWriterLevel(io.Discard, 11)
	},
}

func (c *cache) CompressToMemory(src string) (string, error) {

	input, err := os.Open(src)

	if err != nil {
		return "0", err
	}

	var buf bytes.Buffer

	writer := writerPool.Get().(*brotli.Writer)
	defer writerPool.Put(writer)

	writer.Reset(&buf)

	_, err = io.Copy(writer, input)

	if err != nil {
		return "0", err
	}
	writer.Close()
	c.mux.Lock()
	c.data[src] = buf.Bytes()
	size := int64(len(buf.Bytes()))
	c.mux.Unlock()

	return formatSize(size), nil
}

func CompressToFile(src, dst string) error {
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

	_, err = io.Copy(writer, input)
	if err != nil {
		return fmt.Errorf("Error compressing file: %v", err)
	}

	writer.Close()
	fmt.Println(getFileSize(src))
	fmt.Println(getFileSize(dst))
	return nil
}

func getFileSize(src string) string {
	info, err := os.Stat(src)

	if err != nil {
		// log.Fatalf("error reading file: %v", err)
	}

	return formatSize(info.Size())
}
