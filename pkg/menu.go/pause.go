package menu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func Paused(screen *ebiten.Image, windowWidth int, windowHeight int) {
	font := LoadFont()

	t := "PAUSED"
	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, (windowWidth/2)-(bounds.Dx()/2), (windowHeight/2)-(bounds.Dy()/2), color.White)
}
