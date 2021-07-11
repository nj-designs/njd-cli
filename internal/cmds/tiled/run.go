package tiled

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type Map struct {
	XMLName        xml.Name `xml:"map"`
	Text           string   `xml:",chardata"`
	Version        string   `xml:"version,attr"`
	Tiledversion   string   `xml:"tiledversion,attr"`
	Orientation    string   `xml:"orientation,attr"`
	Renderorder    string   `xml:"renderorder,attr"`
	Width          string   `xml:"width,attr"`
	Height         string   `xml:"height,attr"`
	Tilewidth      string   `xml:"tilewidth,attr"`
	Tileheight     string   `xml:"tileheight,attr"`
	Infinite       string   `xml:"infinite,attr"`
	Nextlayerid    string   `xml:"nextlayerid,attr"`
	Nextobjectid   string   `xml:"nextobjectid,attr"`
	Editorsettings struct {
		Text   string `xml:",chardata"`
		Export struct {
			Text   string `xml:",chardata"`
			Target string `xml:"target,attr"`
			Format string `xml:"format,attr"`
		} `xml:"export"`
	} `xml:"editorsettings"`
	Tileset struct {
		Text     string `xml:",chardata"`
		Firstgid string `xml:"firstgid,attr"`
		Source   string `xml:"source,attr"`
	} `xml:"tileset"`
	Layer struct {
		Text   string `xml:",chardata"`
		ID     string `xml:"id,attr"`
		Name   string `xml:"name,attr"`
		Width  string `xml:"width,attr"`
		Height string `xml:"height,attr"`
		Data   struct {
			Text     string `xml:",chardata"`
			Encoding string `xml:"encoding,attr"`
		} `xml:"data"`
	} `xml:"layer"`
}

const tileMapAsmHeader = `
; Created by njd-cli tiled ...
; Width:{{.Width}} Height:{{.Height}}
`

func checkValidTileMapFileExt(fileExt string) bool {
	switch fileExt {
	case
		".tmx":
		return true

	}
	return false
}

func checkValidTileMapEncoding(mapEncoding string) bool {
	switch mapEncoding {
	case
		"csv":
		return true

	}
	return false
}

func Run(cmd *cobra.Command, args []string) {
	fmt.Printf("Tilemap file: %s\n", TilemapFileName)

	tileMapFile, err := ioutil.ReadFile(TilemapFileName)
	if err != nil {
		log.Fatalf("Failed to read tilemap file: %v", err)
	}

	tileMapFileExt := filepath.Ext(string(TilemapFileName))

	if !checkValidTileMapFileExt(tileMapFileExt) {
		log.Fatalf("%s is not a supported file extenstion", tileMapFileExt)
	}

	mapData := &Map{}

	err = xml.Unmarshal([]byte(tileMapFile), &mapData)
	if err != nil {
		log.Fatalf("Failed to Unmarshal XML tilemap : %v", err)
	}

	//fmt.Println(mapData)

	// Get map sizes
	mapWidth, err := strconv.ParseUint(mapData.Width, 10, 16)
	if err != nil {
		log.Fatalf("Failed to parse mapData.Width : %v", err)
	}

	mapHeight, err := strconv.ParseUint(mapData.Height, 10, 16)
	if err != nil {
		log.Fatalf("Failed to parse mapData.Height : %v", err)
	}

	fmt.Println(mapWidth, mapHeight)

	if !checkValidTileMapEncoding(mapData.Layer.Data.Encoding) {
		log.Fatalf("Unsupported tilemap encoding : %s", mapData.Layer.Data.Encoding)
	}

	// Now parse actually tileMapFile tile data (as csv)
	// r := csv.NewReader(strings.NewReader(mapData.Layer.Data.Text))
	// records, err := r.ReadAll()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(records)
	//fmt.Println(mapData.Layer.Data.Text)

	s := strings.Split(mapData.Layer.Data.Text, ",")

	var b strings.Builder

	// Insert header comment

	t := template.Must(template.New("tileMapAsmHeader").Parse(tileMapAsmHeader))
	err = t.Execute(&b, mapData)
	if err != nil {
		log.Println("executing template:", err)
	}
	for y := 0; y < int(mapHeight); y++ {
		fmt.Fprintf(&b, "    DB ")
		for x := 0; x < int(mapWidth); x++ {
			tileNumStr := strings.TrimSpace(s[(y*int(mapWidth))+x])
			tileNum, err := strconv.ParseUint(tileNumStr, 10, 8)
			if err != nil {
				log.Fatalf("Failed to parse tileNum %s : %v", tileNumStr, err)
			}
			tileNum -= 1
			if x > 0 {
				fmt.Fprintf(&b, ",$%02X", tileNum)
			} else {
				fmt.Fprintf(&b, "$%02X", tileNum)
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	//fmt.Println(b.String())
	err = os.WriteFile(OutputFileName, []byte(b.String()), 0666)
	if err != nil {
		log.Fatalf("Failed to write %s : %v", OutputFileName, err)
	}
}
