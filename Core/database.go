package Core

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Datafile struct {
	XMLName xml.Name `xml:"datafile"`
	Header  Header   `xml:"header"`
	Games   []Game   `xml:"game"`
}

type Header struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
}

type XMLGame struct {
	Name string
}

type Game struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`
	ROM         ROM    `xml:"rom"`
}

type ROM struct {
	Name   string `xml:"name,attr"`
	Size   int64  `xml:"size,attr"`
	CRC    string `xml:"crc,attr"`
	MD5    string `xml:"md5,attr"`
	SHA1   string `xml:"sha1,attr"`
	SHA256 string `xml:"sha256,attr"`
	Status string `xml:"status,attr"`
}

func Parse() {
	data, err := os.ReadFile("/Users/test/Projects/ExNES/Core/no-intro.xml")
	if err != nil {
		panic(err)
	}

	var file Datafile
	if err := xml.Unmarshal(data, &file); err != nil {
		panic(err)
	}

	fmt.Println(file)
}
