package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	Walker
	Active bool

	TeamID int
	Radius float32
	Color  rl.Color

	Range  float32
	Damage float32
}

func (p *Projectile) UpdateState(delta float32) {
	p.Walker.UpdatePosition(delta)

	destinationReached := p.Walker.hasReachedDestination()
	if destinationReached {
		fmt.Println("destination reached!")
		p.Walker.Walking = false
		p.Active = false
	}
}
