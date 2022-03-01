package menu

import (
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/goshlanguage/bento/assets"
)

func GameOver(screen *ebiten.Image, windowWidth int, windowHeight int) {
	textColor := color.White
	font := assets.LoadFont()

	t := "Gameover\nInsert Coin"
	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, windowWidth/2-bounds.Dx(), windowHeight/2, textColor)
}
