package Core

import (
	"embed"
	"io/fs"
)

//go:embed NESpalette/ntsc.pal
var Pallete []byte

//go:embed NESpalette
var palleteSource embed.FS

type FPalDat [8][64][3]byte
type PalDat [64][3]byte

type PalleteData interface {
	getRGB(int, uint8) [3]byte
}

type PalleteEngine struct {
	Loaded  int
	ListPal []PalleteEntry
}

type PalleteEntry struct {
	Name string
	Pal  PalleteData
}

func (p *PalleteEngine) GetRGB(e int, col uint8) [3]uint8 {
	return p.ListPal[p.Loaded].Pal.getRGB(e, col)
}

func (p *PalleteEngine) Init() {
	fs.WalkDir(palleteSource, "", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		data, err := palleteSource.ReadFile(path)
		if err != nil {
			return nil
		}

		pal := LoadPal(data)
		if pal != nil {
			p.ListPal = append(p.ListPal, PalleteEntry{
				Name: path,
				Pal:  pal,
			})
		}

		return nil
	})
}

func LoadPal(data []byte) PalleteData {
	switch len(data) {
	case 8 * 64 * 3:
		var pal FPalDat
		offset := 0

		for e := range 8 {
			for i := range 64 {
				pal[e][i][0] = data[offset+0]
				pal[e][i][1] = data[offset+1]
				pal[e][i][2] = data[offset+2]
				offset += 3
			}
		}

		return &pal
	case 64 * 3:
		var pal PalDat
		offset := 0
		for i := range 64 {
			pal[i][0] = data[offset+0]
			pal[i][1] = data[offset+1]
			pal[i][2] = data[offset+2]
			offset += 3
		}

		return &pal
	default:
		return nil
	}

}

func (p *FPalDat) getRGB(e int, col uint8) [3]byte {
	return p[e][col]
}

func (p *PalDat) getRGB(e int, col uint8) [3]byte {
	return p[col]
}
