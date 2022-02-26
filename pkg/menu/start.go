package menu

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var Font font.Face

func LoadFont() font.Face {
	fontPath := "./assets/fonts/c64.ttf"

	fontData, err := ioutil.ReadFile(fontPath)
	if err != nil {
		buffer := ""
		filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				buffer += fmt.Sprintf("%s\n", path)
				return nil
			})

		panic(fmt.Sprintf("failed to load font face, womp womp\n%s", buffer))
	}

	tt, _ := truetype.Parse(fontData)
	return truetype.NewFace(tt, &truetype.Options{
		Size: 36,
		DPI:  44,
	})
}

func MainMenu(screen *ebiten.Image, windowWidth int, windowHeight int) {
	textColor := color.White
	font := LoadFont()

	n := time.Now().Second()
	if n%2 > 0 {
		textColor = color.Black
	}

	t := "CLICK TO START"
	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, 10, windowHeight-bounds.Dy(), textColor)
}
