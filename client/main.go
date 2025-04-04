package main

import (
	"open-moba/client/game"
	"open-moba/client/renderer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1200, 700, "Open Moba")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	game := game.NewGame(8)

	for !rl.WindowShouldClose() {
		// Update State

		game.HandleInput()
		game.UpdateState()

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		renderer.DrawGame(game)

		rl.EndDrawing()
	}
}
