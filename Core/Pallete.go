package Core

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed NESpalette/ntsc.pal
var Pallete []byte

//go:embed NESpalette
var palleteSource embed.FS

type packedPal [64 * 8]uint32

type PalleteEngine struct {
	Loaded  int
	ListPal []PalleteEntry
}

type PalleteEntry struct {
	Name string
	Pal  packedPal
}

func (p *PalleteEngine) Init() {
	fs.WalkDir(palleteSource, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		data, err := palleteSource.ReadFile(path)
		if err != nil {
			return nil
		}

		pal, err := LoadPal(data)
		if err != nil {
			return nil
		}

		name := strings.TrimSuffix(strings.TrimPrefix(path, "NESpalette/"), ".pal")

		if name == "ntsc" {
			p.Loaded = len(p.ListPal)
		}

		p.ListPal = append(p.ListPal, PalleteEntry{Name: name, Pal: pal})

		return nil
	})
}

func (p *PalleteEngine) GetPalletes() []string {
	list := make([]string, len(p.ListPal))

	for i, item := range p.ListPal {
		list[i] = item.Name
	}

	return list
}

func (p *PalleteEngine) LoadPallete(idx int) {
	p.Loaded = idx
}

func (p *PalleteEngine) getColPacked(e int, idx uint8) uint32 {
	return p.ListPal[p.Loaded].Pal[e*64+int(idx)]
}

func LoadPal(data []byte) (packedPal, error) {
	var out packedPal

	switch len(data) {
	case 8 * 64 * 3:

		offset := 0

		for e := range 8 {
			for i := range 64 {
				out[64*e+i] = packRGB(data[offset], data[offset+1], data[offset+2])
				offset += 3
			}
		}

		return out, nil
	case 64 * 3:
		var base [64]uint32
		offset := 0

		for i := range 64 {
			base[i] = packRGB(data[offset], data[offset+1], data[offset+2])
			offset += 3
		}

		for e := range 8 {
			copy(out[e*64:e*64+64], base[:])
		}

		return out, nil

	default:
		return out, fmt.Errorf("Invalid length: %v ", len(data))
	}

}

func packRGB(r, g, b uint8) uint32 {
	return uint32(r) | uint32(g)<<8 | uint32(b)<<16 | 0xFF000000
}
