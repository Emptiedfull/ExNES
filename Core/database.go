package Core

import (
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"os"
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

type gameDir map[string]string

func ParseDb() Datafile {
	data, err := os.ReadFile("/Users/test/Projects/ExNES/Core/no-intro.xml")
	if err != nil {
		panic(err)
	}

	var file Datafile
	if err := xml.Unmarshal(data, &file); err != nil {
		panic(err)
	}

	// for _, rom := range file.Games {
	// 	fmt.Println(rom.Name)
	// }

	return file

	// fmt.Println(file)
}

func CreateGameDir(data Datafile, check string) gameDir {
	dir := make(gameDir)
	for _, game := range data.Games {
		dir[game.Name] = game.ROM.Crc
		if check == game.ROM.Crc {
			fmt.Println("found:", game.ROM.Name)
		}

	}

	return dir
}

func HashCrc(data []byte) string {

	hash := crc32.ChecksumIEEE(data)
	return fmt.Sprintf("%08x", hash)
}
