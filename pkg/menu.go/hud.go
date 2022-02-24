package menu

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func HUD(screen *ebiten.Image, particles int, windowWidth int, windowHeight int) {
	font := LoadFont()

	t := strconv.Itoa(particles)

	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, windowWidth-bounds.Dx()-10, bounds.Dy()+10, color.White)
}
