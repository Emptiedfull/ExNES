package Core

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
)

type Datafile struct {
	XMLName xml.Name `xml:"datafile"`
	Games   []Game   `xml:"game"`
}

type Game struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`
	ROM         ROM    `xml:"rom"`
}

type ROM struct {
	Name string `xml:"name,attr"`
	Crc  string `xml:"crc,attr"`
}

type packgedEntry struct {
	Name string
	Hash [4]byte
}

func ParseDb() Datafile {
	data, err := os.ReadFile("/Users/test/Projects/ExNES/Core/no-intro.xml")
	if err != nil {
		panic(err)
	}

	var file Datafile
	if err := xml.Unmarshal(data, &file); err != nil {
		panic(err)
	}

	return file

}

func WriteBin(datafile Datafile) {
	entries := make([]packgedEntry, 0, len(datafile.Games))

	for _, game := range datafile.Games {
		hash, err := HashCrc(game.ROM.Crc)
		if err != nil {
			panic(err)
		}
		entries = append(entries, packgedEntry{Hash: hash, Name: game.ROM.Name})
	}

	sort.Slice(entries, func(i, j int) bool {
		return (bytes.Compare(entries[i].Hash[:], entries[j].Hash[:])) < 0
	})

	file, err := os.Create("./db.bin")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	var lenbuf [binary.MaxVarintLen64]byte

	l := binary.PutUvarint(lenbuf[:], uint64(len(datafile.Games)))
	if _, err := w.Write(lenbuf[:l]); err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if _, err := w.Write(entry.Hash[:]); err != nil {
			panic(err)
		}
	}

	for _, entry := range entries {
		l := binary.PutUvarint(lenbuf[:], uint64(len(entry.Name)))
		if _, err := w.Write(lenbuf[:l]); err != nil {
			panic(err)
		}
		if _, err := w.WriteString(entry.Name); err != nil {
			panic(err)
		}
	}
}

func HashCrc(hash string) ([4]byte, error) {
	var dst [4]byte

	slice, err := hex.DecodeString(hash)
	if err != nil {
		return dst, err
	}
	if len(slice) != 4 {
		return dst, fmt.Errorf("bad hashing ")
	}

	copy(dst[:], slice)
	return dst, nil
}
