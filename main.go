package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	ballSize = 2
	maxVX    = 5.0
	maxVY    = 4.0
	width    = 640
	height   = 480
	friction = .005
)

type Map map[int]map[int]entity

type entity interface {
	Update(m Map)
	Draw(*ebiten.Image)
}

type Ball struct {
	X, Y           float64
	VX, VY         float64
	ScaleX, ScaleY float64
	Sprite         *ebiten.Image
}

func NewBall(x, y, vx, vy float64) Ball {
	img, _, err := ebitenutil.NewImageFromFile("ball.png")
	if err != nil {
		log.Fatal(err)
	}

	w, h := img.Size()

	scaleX := float64(ballSize / w)
	scaleY := float64(ballSize / h)

	return Ball{
		x, y,
		vx, vy,
		scaleX, scaleY,
		img,
	}
}

func (b *Ball) Update(m Map) {
	if m[int(b.X)/ballSize][int(b.Y)/ballSize] != nil {
		m[int(b.X)/ballSize][int(b.Y)/ballSize] = nil
	}

	b.X += b.VX
	b.Y += b.VY

	if m[int(b.X)/ballSize][int(b.Y)/ballSize] != nil {
		b.X -= b.VX - b.VX
		b.Y -= b.VY - b.VY
		b.VX /= 2
		b.VY /= 2

		collided := m[int(b.X)/ballSize][int(b.Y)/ballSize].(*Ball)
		collided.VX += b.VX
		collided.VY += b.VY
	}

	// If we're over limits, reset limits and positionals
	if b.X <= 0 || b.X >= 640 {
		b.VX = -1 * b.VX
	}

	if b.Y <= 0 || b.Y > 480 {
		b.VY = -1 * b.VY
	}

	if b.X < 0 {
		b.X = 0
	}

	if b.X > 640 {
		b.X = 640
	}

	if b.Y < 0 {
		b.Y = 0
	}

	if b.Y > 480 {
		b.Y = 480
	}

	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	cursorX, cursorY := ebiten.CursorPosition()
	// 	if b.VX >= -maxVX && b.VX <= maxVX {
	// 		if b.X-float64(cursorX) > 0 && float64(cursorX)-b.VX <= maxVX {
	// 			b.VX -= (b.X - float64(cursorX)) / 640
	// 		} else if b.VX >= -maxVX {
	// 			b.VX += (float64(cursorX) - b.X) / 640
	// 		}
	// 	}

	// 	if b.VY >= -maxVY && b.VY <= maxVY {
	// 		if b.Y-float64(cursorY) > 0 && float64(cursorY)-b.VY <= maxVY {
	// 			b.VY -= (b.Y - float64(cursorY)) / 480
	// 		} else if b.VY >= -maxVY {
	// 			b.VY += (float64(cursorY) - b.Y) / 480
	// 		}
	// 	}
	// }

	if b.VX > 0 {
		b.VX -= friction
	} else {
		b.VX += friction
	}

	if b.VY > 0 {
		b.VY -= friction
	} else {
		b.VY += friction
	}

	if b.VX > maxVX {
		b.VX = maxVX
	}

	if b.VX < -maxVX {
		b.VX = -maxVX
	}

	if b.VY > maxVY {
		b.VY = maxVY
	}

	if b.VY < -maxVY {
		b.VY = -maxVY
	}

	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
	// 	b.VX = 0
	// 	b.VY = 0
	// }

	m[int(b.X)/ballSize][int(b.Y)/ballSize] = b
}

func (b *Ball) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(b.X, b.Y)
	if b.ScaleX != 1 || b.ScaleY != 1 {
		if b.ScaleX == 0 {
			b.ScaleX = 1
		}

		if b.ScaleY == 0 {
			b.ScaleX = 1
		}

		opts.GeoM.Scale(b.ScaleX, b.ScaleY)
	}
	screen.DrawImage(b.Sprite, opts)
}

func (b *Ball) UpdateV(VX, VY float64) {
	b.VX = VX
	b.VY = VY
}

type Game struct {
	entities []entity
	m        Map
}

func init() {}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()

		if g.m[cursorX/ballSize][cursorY/ballSize] == nil {
			dx := (float64(cursorX) - (float64(width) / 2)) / float64(width) * 2 * maxVX
			dy := (float64(cursorY) - (float64(height) / 2)) / float64(height) * 2 * maxVY
			newBall := NewBall(float64(width)/2, float64(height)/2, dx, dy)
			g.entities = append(g.entities, &newBall)
		}
	}

	for _, e := range g.entities {
		e.Update(g.m)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, e := range g.entities {
		e.Draw(screen)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", len(g.entities)), 600, 400)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Sandbox")

	m := Map{}
	for i := 0; i <= width/ballSize; i++ {
		m[i] = make(map[int]entity)
	}

	instance := &Game{
		entities: []entity{},
		m:        m,
	}

	if err := ebiten.RunGame(instance); err != nil {
		log.Fatal(err)
	}
}
