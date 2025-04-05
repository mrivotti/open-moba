package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	Walker
	Active bool

	Abilities []Ability

	TeamID    int
	Health    float32
	MaxHealth float32

	Radius float32
	Color  rl.Color

	Damage int
}
