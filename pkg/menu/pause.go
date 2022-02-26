package menu

import (
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/goshlanguage/bento/assets"
)

func Paused(screen *ebiten.Image, windowWidth int, windowHeight int) {
	textColor := color.White
	font := assets.LoadFont()

	t := "PAUSED"

	s := time.Now().Second()
	n := time.Now().Nanosecond()
	if n > 2*(999999999/3) {
		textColor = color.Black
	}

	text.Draw(screen, strconv.Itoa(s), font, 300, 300, textColor)

	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, (windowWidth/2)-(bounds.Dx()/2), (windowHeight/2)-(bounds.Dy()/2), textColor)
}
