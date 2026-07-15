package Core

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
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
