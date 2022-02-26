package assets

import (
	_ "embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed fonts/c64.ttf
var embeddedFont []byte

var Font font.Face

func LoadFont() font.Face {
	tt, _ := truetype.Parse(embeddedFont)
	return truetype.NewFace(tt, &truetype.Options{
		Size: 36,
		DPI:  44,
	})
}
