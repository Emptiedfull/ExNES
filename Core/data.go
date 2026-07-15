package Core

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
)

//go:embed db.bin
var romData []byte

type RomDB struct {
	count  int
	hashes int
	Names  []string
}

func LoadDB() (*RomDB, error) {
	count, n := binary.Uvarint(romData)
	if n <= 0 {
		return nil, fmt.Errorf("fucking empty map")
	}

	hashesOff := n
	namesOff := hashesOff + int(count)*4

	names := make([]string, count)
	off := namesOff
	for i := range count {
		l, n := binary.Uvarint(romData[off:])
		if n <= 0 {
			return nil, fmt.Errorf("fucked name length at entry %v", i)
		}

		off += n

		names[i] = string(romData[off : off+int(l)])
		off += int(l)
	}

	db := &RomDB{
		count:  int(count),
		hashes: hashesOff,
		Names:  names,
	}

	return db, nil
}

func (db *RomDB) LookUp(hash [4]byte) (string, bool) {
	l, h := 0, db.count-1
	for l <= h {
		mid := (l + h) / 2
		offset := db.hashes + mid*4
		cmp := bytes.Compare(romData[offset:offset+4], hash[:])

		switch {
		case cmp == 0:
			return db.Names[mid], true
		case cmp < 0:
			l = mid + 1
		case cmp > 0:
			h = mid - 1
		}
	}

	return "", false
}

func (db *RomDB) Get(hash uint32) (string, bool) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], hash)
	return db.LookUp(b)
}

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
