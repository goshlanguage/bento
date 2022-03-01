package player

import (
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/goshlanguage/bento/pkg/controller"
	"github.com/goshlanguage/bento/pkg/types"
)

const (
	VX = 3.0
	VY = 3.0
)

type Player struct {
	Controller   controller.Controller
	Sprite       *ebiten.Image
	HP           int
	LastX, LastY float64
	X, Y         float64
	Rotation     float64
	MaxX, MaxY   float64
	MapModifier  int
}

func (p *Player) Update(m types.Map) {
	if m[int(p.X)/p.MapModifier][int(p.Y)/p.MapModifier] != nil {
		p.HP -= m[int(p.X)/p.MapModifier][int(p.Y)/p.MapModifier].Damage()
		m[int(p.X)/p.MapModifier][int(p.Y)/p.MapModifier] = nil
	}

	if ebiten.IsKeyPressed(p.Controller.Up) {
		p.Y -= VY
		if p.Y < 0.0 {
			p.Y = 0.0
		}
	}

	if ebiten.IsKeyPressed(p.Controller.Down) {
		p.Y += VY
		if p.Y > float64(p.MaxY) {
			p.Y = float64(p.MaxY)
		}
	}

	if ebiten.IsKeyPressed(p.Controller.Left) {
		p.X -= VX
		if p.X < 0.0 {
			p.X = 0.0
		}
	}

	if ebiten.IsKeyPressed(p.Controller.Right) {
		p.X += VX
		if p.X > float64(p.MaxX) {
			p.X = float64(p.MaxX)
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	cursorX, cursorY := ebiten.CursorPosition()

	deltaX := p.X - float64(cursorX)
	deltaY := p.Y - float64(cursorY)

	theta := math.Atan2(deltaY, deltaX)
	opts.GeoM.Rotate(theta)
	opts.GeoM.Translate(p.X, p.Y)

	screen.DrawImage(p.Sprite, opts)
}

func Default(width, height, modifier int) *Player {
	img, _, err := ebitenutil.NewImageFromFile("assets/gfx/sprite.png")
	if err != nil {
		log.Fatal(err)
	}

	sprite := img.SubImage(image.Rect(0, 0, 8, 8)).(*ebiten.Image)

	return &Player{
		Controller:  controller.Default(),
		HP:          100,
		Sprite:      sprite,
		X:           float64(width / 2),
		Y:           float64(height / 2),
		LastX:       float64(width / 2),
		LastY:       float64(height / 2),
		MaxX:        float64(width),
		MaxY:        float64(height),
		MapModifier: modifier,
	}
}
