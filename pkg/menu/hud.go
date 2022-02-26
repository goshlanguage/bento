package menu

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/goshlanguage/bento/assets"
)

func HUD(screen *ebiten.Image, particles int, windowWidth int, windowHeight int) {
	font := assets.LoadFont()

	t := strconv.Itoa(particles)

	textColor := color.White

	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, windowWidth-bounds.Dx()-10, bounds.Dy()+10, textColor)
}
