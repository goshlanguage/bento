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
	"github.com/goshlanguage/bento/pkg/player"
	"github.com/goshlanguage/bento/pkg/types"
)

var (
	ballSize      = 2
	width         = 640
	height        = 480
	maxVX         = 5.0
	maxVY         = 4.0
	particleLimit = 512

	audioContext *audio.Context
	sfxMap       []*audio.Player
	gameover     bool
)

type Game struct {
	entities []types.Entity
	m        types.Map
	players  []*player.Player
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
		"assets/sfx/rip.wav",
	} {
		player, err := LoadWav(file)
		if err != nil {
			fmt.Printf("Crap, %v\n", err)
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
	if !gameover {
		for _, player := range g.players {
			player.Update(g.m)

			if player.HP <= 0 {
				sfxMap[3].Rewind()
				sfxMap[3].Play()

				g.entities = []types.Entity{}
				gameover = true
			}
		}

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

			dx := (float64(cursorX) - (float64(width) / 2)) / float64(width) * 2 * maxVX
			dy := (float64(cursorY) - (float64(height) / 2)) / float64(height) * 2 * maxVY

			newBall := ball.NewBall(g.players[0].X+4, g.players[0].Y, dx, dy)
			g.entities = append(g.entities, &newBall)
		}

		counter := 0
		numEntities := len(g.entities)

		limit := numEntities / 10
		if limit > particleLimit/10 {
			limit = particleLimit / 10
		}

		for k, e := range g.entities {
			e.Update(g.m, g.state)

			// TODO refactor entity reaping into individual entity Update methods
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
	} else {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustReleased(g.players[0].Controller.Start) || inpututil.IsTouchJustReleased(0) {
			gameover = false

			g.players[0].HP = 100
			newGame := InitGame()
			g = newGame
		}
	}

	return nil
}

// Draw is the main draw portion of the Game loop, which triggers all subDraws using the screen
func (g *Game) Draw(screen *ebiten.Image) {
	if !gameover {
		for _, e := range g.entities {
			e.Draw(screen)
		}

		for _, player := range g.players {
			player.Draw(screen)
		}

		if *g.state["menu"].(*bool) && !*g.state["started"].(*bool) {
			menu.MainMenu(screen, width, height)
		}

		if *g.state["menu"].(*bool) && *g.state["started"].(*bool) {
			menu.Paused(screen, width, height)
		}

		if !*g.state["menu"].(*bool) && *g.state["started"].(*bool) {
			menu.HUD(screen, g.players[0].HP, width, height)
		}
	} else {
		menu.GameOver(screen, width, height)
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

func InitGame() *Game {
	p := player.Default(width, height, ballSize)
	players := []*player.Player{
		p,
	}

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

	return &Game{
		entities: []types.Entity{},
		m:        m,
		players:  players,
		toggles:  map[string]bool{},
		state:    state,
	}
}

// main loop for the game
// sets up the window, the initial game map and state, then starts
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("弁当 Bento")

	instance := InitGame()

	if err := ebiten.RunGame(instance); err != nil {
		log.Fatal(err)
	}
}
