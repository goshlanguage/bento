package types

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Map map[int]map[int]Entity

type Entity interface {
	Damage() int
	Update(m Map, s map[string]interface{})
	Draw(*ebiten.Image)
	TimesUp() bool
}
