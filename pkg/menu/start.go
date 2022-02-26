package menu

import (
	_ "embed"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/goshlanguage/bento/assets"
)

func MainMenu(screen *ebiten.Image, windowWidth int, windowHeight int) {
	textColor := color.White
	font := assets.LoadFont()

	n := time.Now().Second()
	if n%2 > 0 {
		textColor = color.Black
	}

	t := "CLICK TO START"
	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, 10, windowHeight-bounds.Dy(), textColor)
}
