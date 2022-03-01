package menu

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/goshlanguage/bento/assets"
)

func HUD(screen *ebiten.Image, HP int, windowWidth int, windowHeight int) {
	font := assets.LoadFont()

	t := strconv.Itoa(HP)

	textColor := color.White

	bounds := text.BoundString(font, fmt.Sprintf("HP: %v", t))
	text.Draw(screen, t, font, windowWidth-bounds.Dx()-10, bounds.Dy()+10, textColor)
}
