package main

import (
	"open-moba/client/game"
	"open-moba/client/renderer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var screenWidth int32 = 1200
	var screenHeight int32 = 700

	rl.InitWindow(screenWidth, screenHeight, "Open Moba")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	renderer := renderer.NewRenderer(screenWidth, screenHeight)

	game, err := game.NewGame(8, 0)
	if err != nil {
		panic(err)
	}

	for !rl.WindowShouldClose() {
		// Update State

		game.HandleInput()
		game.UpdateState()

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		renderer.DrawBackground()
		renderer.DrawGame(game)
		renderer.DrawHUD(game)

		rl.EndDrawing()
	}
}
