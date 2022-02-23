package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ballSize      = 2
	maxVX         = 5.0
	maxVY         = 4.0
	width         = 640
	height        = 480
	friction      = .001
	menu          = true
	menuTxt       = "Click Me..."
	menuTxtX      = width/2 - 50
	menuTxtY      = height / 2
	particleLimit = 512
	started       = false
	wrapAround    = true

	audioContext *audio.Context
	sfxMap       []*audio.Player
)

type Map map[int]map[int]entity

type entity interface {
	Update(m Map)
	Draw(*ebiten.Image)
	TimesUp() bool
}

type Ball struct {
	Expired        bool
	X, Y           float64
	VX, VY         float64
	ScaleX, ScaleY float64
	Sprite         *ebiten.Image
	Spawned        time.Time
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

func (b *Ball) Update(m Map) {
	if !menu || !started {
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

type Game struct {
	entities []entity
	m        Map
}

func init() {
	rand.Seed(time.Now().UnixNano())

	ctx := audio.NewContext(44100)
	audioContext = ctx

	for _, file := range [...]string{"bounce.wav", "pause.wav", "unpause.wav"} {
		player, err := LoadWav(file)
		if err != nil {
			fmt.Errorf("Crap, %v", err)
		}

		sfxMap = append(sfxMap, player)
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		menu = !menu
		if menu {
			sfxMap[1].Rewind()
			sfxMap[1].Play()
		} else {
			sfxMap[2].Rewind()
			sfxMap[2].Play()
		}
	}

	if !started {
		if 1 == rand.Intn(50) {
			randX := rand.Float64() * float64(width)
			randY := rand.Float64() * float64(height)

			rand.Float64()
			randVX := 1 + rand.Float64()*maxVX
			randVY := 1 + rand.Float64()*maxVX

			ball := NewBall(randX, randY, randVX, randVY)

			g.entities = append(g.entities, &ball)
		}

		for _, entity := range g.entities {
			entity.Update(g.m)
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if menu {
			if !started {
				menuTxt = "Paused. Press space to continue"
				menuTxtX = width/2 - 100
				started = true
			}
			menu = !menu
		}

		if !sfxMap[0].IsPlaying() {
			sfxMap[0].Rewind()
			sfxMap[0].SetVolume(rand.Float64())
			sfxMap[0].Play()
		}

		cursorX, cursorY := ebiten.CursorPosition()

		if g.m[cursorX/ballSize][cursorY/ballSize] == nil {
			dx := (float64(cursorX) - (float64(width) / 2)) / float64(width) * 2 * maxVX
			dy := (float64(cursorY) - (float64(height) / 2)) / float64(height) * 2 * maxVY
			newBall := NewBall(float64(width)/2, float64(height)/2, dx, dy)
			g.entities = append(g.entities, &newBall)
		}
	}

	counter := 0
	numEntities := len(g.entities)

	limit := numEntities / 10
	if limit > particleLimit/10 {
		limit = particleLimit / 10
	}

	for k, e := range g.entities {
		e.Update(g.m)
		if e.TimesUp() {
			if len(g.entities) > particleLimit {
				if k+1 >= len(g.entities) {
					g.entities = g.entities[:k]
				} else {
					g.entities = append(g.entities[:k], g.entities[k+1:]...)
				}

			} else {
				// only allow limit deletions per run, should be a small number relative to the overall size of the entity slice
				if counter < limit {
					if k+1 >= len(g.entities) {
						g.entities = g.entities[:k]
					} else {
						g.entities = append(g.entities[:k], g.entities[k+1:]...)
					}
				}

				counter++
			}
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, e := range g.entities {
		e.Draw(screen)
	}

	if menu {
		ebitenutil.DebugPrintAt(screen, menuTxt, menuTxtX, menuTxtY)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", len(g.entities)), 600, 440)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func LoadWav(filepath string) (*audio.Player, error) {
	var errs error

	f, err := ebitenutil.OpenFile(filepath)
	if err != nil {
		return nil, err
	}

	d, err := wav.Decode(audioContext, f)
	if err != nil {
		return nil, err
	}

	audioPlayer, err := audio.NewPlayer(audioContext, d)
	if err != nil {
		return nil, err
	}

	return audioPlayer, errs
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("べんとう(弁当) Bento")

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
