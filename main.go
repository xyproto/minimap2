package main

import (
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

var (
	keywords       = []string{"func", "if", "else", "for", "return", "struct", "enum", "match", "use", "mod", "const", "pub", "def", "print", "#include", "int", "float", "char", "double"}
	commentMarkers = []string{"//", "/*", "*/", "#", "///"}
	stringMarkers  = []string{"\"", "'"}
)

var (
	colorKeyword    = color.RGBA{255, 0, 0, 255}
	colorString     = color.RGBA{0, 255, 0, 255}
	colorComment    = color.RGBA{0, 0, 255, 255}
	colorNormal     = color.RGBA{255, 255, 255, 255}
	colorWhitespace = color.RGBA{0, 0, 0, 0}
)

func getColorForChar(ch, nextCh string, inString, inComment *bool) color.RGBA {
	if *inString {
		if ch == "\"" || ch == "'" {
			*inString = false
		}
		return colorString
	} else if *inComment {
		if ch == "*" && nextCh == "/" {
			*inComment = false
		}
		return colorComment
	} else {
		if ch == "/" && (nextCh == "/" || nextCh == "*") {
			*inComment = true
			return colorComment
		} else if ch == "\"" || ch == "'" {
			*inString = true
			return colorString
		} else if strings.ContainsAny(ch, " \t\n") {
			return colorWhitespace
		} else if contains(keywords, ch) {
			return colorKeyword
		} else {
			return colorNormal
		}
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func main() {
	filename := os.Args[1]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	content := string(data)
	dc := gg.NewContext(800, 1000)
	dc.SetRGB(0.15, 0.15, 0.15)
	dc.Clear()

	fontPath := "SometypeMono-SemiBold.ttf"
	fontSize := 20.0
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}

	inString := false
	inComment := false
	x, y := 10.0, fontSize
	for i := 0; i < len(content); i++ {
		ch := string(content[i])
		nextCh := ""
		if i+1 < len(content) {
			nextCh = string(content[i+1])
		}

		col := getColorForChar(ch, nextCh, &inString, &inComment)

		if ch == "\n" {
			y += fontSize
			x = 10.0
			continue
		}

		dc.SetColor(col)
		dc.DrawString(ch, x, y)
		width, _ := dc.MeasureString(ch)
		x += width
	}

	dc.SavePNG("output.png")
	log.Println("Image saved as output.png")
}
