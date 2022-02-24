package menu

import (
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func Paused(screen *ebiten.Image, windowWidth int, windowHeight int) {
	textColor := color.White
	font := LoadFont()

	t := "PAUSED"

	n := time.Now().Second()
	if n%2 > 0 {
		textColor = color.Black
	}
	text.Draw(screen, strconv.Itoa(n), font, 300, 300, color.White)

	bounds := text.BoundString(font, t)
	text.Draw(screen, t, font, (windowWidth/2)-(bounds.Dx()/2), (windowHeight/2)-(bounds.Dy()/2), textColor)
}
