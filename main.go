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

	"github.com/goshlanguage/bento/pkg/ball"
	"github.com/goshlanguage/bento/pkg/menu"
	"github.com/goshlanguage/bento/pkg/types"
)

var (
	ballSize      = 2
	width         = 640
	height        = 480
	maxVX         = 5.0
	maxVY         = 4.0
	particleLimit = 512
	playerX       = float64(width) / 2
	playerY       = float64(height) / 2
	playerVX      = 3.0
	playerVY      = 3.0
	trueTrue      = true

	audioContext *audio.Context
	sfxMap       []*audio.Player
)

type Game struct {
	entities []types.Entity
	m        types.Map
	toggles  map[string]bool
	state    map[string]interface{}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	ctx := audio.NewContext(44100)
	audioContext = ctx

	for _, file := range [...]string{
		"assets/sfx/shot.wav",
		"assets/sfx/pause.wav",
		"assets/sfx/menu.wav",
	} {
		player, err := LoadWav(file)
		if err != nil {
			fmt.Errorf("Crap, %v", err)
		}

		sfxMap = append(sfxMap, player)
	}
	// chill out with the menu volume
	sfxMap[2].SetVolume(0.3)
}

func toggle(state map[string]interface{}, key string) map[string]interface{} {
	value := *state[key].(*bool)
	value = !value

	state[key] = &value

	return state
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.state = toggle(g.state, "menu")

		if *g.state["menu"].(*bool) {
			sfxMap[1].Rewind()
			sfxMap[1].Play()
			sfxMap[2].Rewind()
			sfxMap[2].Play()
		} else {
			sfxMap[1].Rewind()
			sfxMap[1].Play()
			sfxMap[2].Pause()
		}
	}

	if !*g.state["started"].(*bool) {
		if !sfxMap[2].IsPlaying() {
			sfxMap[2].Rewind()
			sfxMap[2].Play()
		}

		if 1 == rand.Intn(50) {
			randX := rand.Float64() * float64(width)
			randY := rand.Float64() * float64(height)

			rand.Float64()
			randVX := 1 + rand.Float64()*maxVX
			randVY := 1 + rand.Float64()*maxVX

			ball := ball.NewBall(randX, randY, randVX, randVY)

			g.entities = append(g.entities, &ball)
		}

		for _, entity := range g.entities {
			entity.Update(g.m, g.state)
		}
	} else if !*g.state["menu"].(*bool) {
		if sfxMap[2].IsPlaying() {
			sfxMap[2].Pause()
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if *g.state["menu"].(*bool) {
			if !*g.state["started"].(*bool) {
				g.state = toggle(g.state, "started")
			}
			g.state = toggle(g.state, "menu")
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

			newBall := ball.NewBall(playerX, playerY, dx, dy)
			g.entities = append(g.entities, &newBall)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		playerY -= playerVY
		if playerY < 0.0 {
			playerY = 0.0
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerY += playerVY
		if playerY > float64(height) {
			playerY = float64(height)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		playerX -= playerVX
		if playerX < 0.0 {
			playerX = 0.0
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerX += playerVX
		if playerX > float64(width) {
			playerX = float64(width)
		}
	}

	counter := 0
	numEntities := len(g.entities)

	limit := numEntities / 10
	if limit > particleLimit/10 {
		limit = particleLimit / 10
	}

	for k, e := range g.entities {
		e.Update(g.m, g.state)
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

// Draw is the main draw portion of the Game loop, which triggers all subDraws using the screen
func (g *Game) Draw(screen *ebiten.Image) {
	for _, e := range g.entities {
		e.Draw(screen)
	}

	if *g.state["menu"].(*bool) && !*g.state["started"].(*bool) {
		menu.MainMenu(screen, width, height)
	}

	if *g.state["menu"].(*bool) && *g.state["started"].(*bool) {
		menu.Paused(screen, width, height)
	}

	if !*g.state["menu"].(*bool) && *g.state["started"].(*bool) {
		menu.HUD(screen, len(g.entities), width, height)
	}
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

// main loop for the game
// sets up the window, the initial game map and state, then starts
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("べんとう(弁当) Bento")

	m := types.Map{}
	for i := 0; i <= width/ballSize; i++ {
		m[i] = make(map[int]types.Entity)
	}

	state := make(map[string]interface{})

	state["width"] = width
	state["height"] = height

	state["maxVX"] = maxVX
	state["maxVY"] = maxVY

	menu := true
	started := false

	state["menu"] = &menu
	state["started"] = &started

	instance := &Game{
		entities: []types.Entity{},
		m:        m,
		toggles:  map[string]bool{},
		state:    state,
	}

	if err := ebiten.RunGame(instance); err != nil {
		log.Fatal(err)
	}
}
