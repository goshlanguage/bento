package controller

import "github.com/hajimehoshi/ebiten/v2"

type Controller struct {
	Up    ebiten.Key
	Down  ebiten.Key
	Left  ebiten.Key
	Right ebiten.Key

	A ebiten.Key
	B ebiten.Key

	Start  ebiten.Key
	Select ebiten.Key
}

func Default() Controller {
	return Controller{
		Up:    ebiten.KeyW,
		Left:  ebiten.KeyA,
		Down:  ebiten.KeyS,
		Right: ebiten.KeyD,

		A: ebiten.KeyE,
		B: ebiten.KeyE,

		Start:  ebiten.KeySpace,
		Select: ebiten.KeyEscape,
	}
}
