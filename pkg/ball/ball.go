package ball

import (
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/goshlanguage/bento/pkg/types"
)

const (
	ballSize   = 2
	friction   = .001
	wrapAround = true
)

type Ball struct {
	Expired        bool
	X, Y           float64
	VX, VY         float64
	ScaleX, ScaleY float64
	Sprite         *ebiten.Image
	Spawned        time.Time
}

func NewBall(x, y, vx, vy float64) Ball {
	img, _, err := ebitenutil.NewImageFromFile("assets/gfx/ball.png")
	if err != nil {
		log.Fatal(err)
	}

	w, h := img.Size()

	scaleX := float64(ballSize / w)
	scaleY := float64(ballSize / h)

	return Ball{
		Expired: false,
		X:       x,
		Y:       y,
		VX:      vx,
		VY:      vy,
		ScaleX:  scaleX,
		ScaleY:  scaleY,
		Sprite:  img,
		Spawned: time.Now().UTC(),
	}
}

func (b *Ball) Update(m types.Map, s map[string]interface{}) {
	width := s["width"].(int)
	height := s["height"].(int)

	maxVX := s["maxVX"].(float64)
	maxVY := s["maxVY"].(float64)

	menu := s["menu"].(*bool)
	started := s["started"].(*bool)

	if !*menu || !*started {
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

		if wrapAround {
			if b.X < 0 {
				b.X = float64(width)
			}

			if b.X > 640 {
				b.X = 0
			}

			if b.Y < 0 {
				b.Y = float64(height)
			}

			if b.Y > 480 {
				b.Y = 0
			}
		} else {
			// If we're over limits, bounce the ball back into the viewport
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

		expiry := 3 * time.Second
		now := time.Now().UTC()

		if now.Sub(b.Spawned) > expiry {
			b.Expired = true
		}
	}
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

func (b *Ball) TimesUp() bool {
	return b.Expired
}

func (b *Ball) UpdateV(VX, VY float64) {
	b.VX = VX
	b.VY = VY
}
