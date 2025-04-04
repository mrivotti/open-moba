package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Walker struct {
	Pos         rl.Vector2
	Destination rl.Vector2
	Speed       float32
	Walking     bool
}

func (w *Walker) SetDestination(d rl.Vector2) {
	w.Destination = d
	w.Walking = true
}

func (w *Walker) UpdatePosition(delta float32) {
	if !w.Walking {
		return
	}

	dir := rl.Vector2Subtract(w.Destination, w.Pos)
	w.Pos = rl.Vector2Add(w.Pos, rl.Vector2Scale(rl.Vector2Normalize(dir), w.Speed*delta))
}

func (w *Walker) hasReachedDestination() bool {
	return rl.Vector2Distance(w.Destination, w.Pos) < 5
}
