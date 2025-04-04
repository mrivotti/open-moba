package game

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const maxProjectiles = 100
const collisionMargin = 5

type Game struct {
	Players            []Player
	Projectiles        []Projectile
	CurrentProjectiles int
}

func NewGame(nPlayers int) *Game {
	players := make([]Player, nPlayers)
	players[0] = Player{
		Walker:    Walker{Pos: rl.Vector2{X: 100, Y: 100}, Speed: 100},
		Radius:    20,
		Active:    true,
		Color:     rl.Blue,
		TeamID:    1,
		Health:    100,
		MaxHealth: 100,
	}
	for i := 1; i < nPlayers; i++ {
		x := rand.Float32() * 1000
		y := rand.Float32() * 600

		players[i] = Player{
			Walker:    Walker{Pos: rl.Vector2{X: x, Y: y}, Speed: 100},
			Radius:    20,
			Active:    true,
			Color:     rl.Yellow,
			TeamID:    2,
			Health:    100,
			MaxHealth: 100,
		}
	}

	projectiles := make([]Projectile, 0, 100)

	return &Game{
		Players:     players,
		Projectiles: projectiles,
	}
}

func (g *Game) HandleInput() {
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		g.Players[0].SetDestination(mousePos)
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		g.spawnProjectile(g.Players[0].Walker.Pos, mousePos, rl.Red, 15, 500, 5, 300)
	}
	if rl.IsKeyPressed(rl.KeyW) {
		g.spawnProjectile(g.Players[0].Walker.Pos, mousePos, rl.Orange, 50, 100, 30, 500)
	}

	if rl.IsKeyPressed(rl.KeyE) {
		var radius float32 = 7.0
		var speed float32 = 700
		var damage float32 = 5.0
		var projectileRange float32 = 150.0
		color := rl.Green
		var angle float32 = math.Pi / 4

		g.spawnProjectile(g.Players[0].Walker.Pos, mousePos, color, radius, speed, damage, projectileRange)

		leftProjectileDir := rl.Vector2Rotate(rl.Vector2Subtract(mousePos, g.Players[0].Walker.Pos), angle)
		g.spawnProjectile(g.Players[0].Walker.Pos, rl.Vector2Add(g.Players[0].Walker.Pos, leftProjectileDir), color, radius, speed, damage, projectileRange)

		rightProjectileDir := rl.Vector2Rotate(rl.Vector2Subtract(mousePos, g.Players[0].Walker.Pos), -angle)
		g.spawnProjectile(g.Players[0].Walker.Pos, rl.Vector2Add(g.Players[0].Walker.Pos, rightProjectileDir), color, radius, speed, damage, projectileRange)
	}

}

func (g *Game) UpdateState() {
	delta := rl.GetFrameTime()

	for i := range g.Players {
		if !g.Players[i].Active {
			continue
		}

		g.Players[i].Walker.UpdatePosition(delta)

		if g.Players[i].Walker.hasReachedDestination() {
			g.Players[i].Walker.Walking = false
		}
	}

	for i := range g.Projectiles {
		if !g.Projectiles[i].Active {
			continue
		}

		g.Projectiles[i].Walker.UpdatePosition(delta)

		if g.Projectiles[i].Walker.hasReachedDestination() {
			g.Projectiles[i].Walker.Walking = false
			g.Projectiles[i].Active = false
		}
	}

	g.handleProjectileCollisions()

	g.removeInactiveProjectiles()
}

func (game *Game) handleProjectileCollisions() {
	for i := range game.Players {
		if !game.Players[i].Active {
			continue
		}
		for j := 0; j < len(game.Projectiles); j++ {
			if !game.Projectiles[j].Active {
				continue
			}
			distance := rl.Vector2Distance(game.Players[i].Walker.Pos, game.Projectiles[j].Walker.Pos)
			if distance < game.Players[i].Radius+game.Projectiles[j].Radius-collisionMargin && game.Players[i].TeamID != game.Projectiles[j].TeamID {
				game.Players[i].Health -= game.Projectiles[j].Damage
				if game.Players[i].Health <= 0 {
					game.Players[i].Active = false
				}

				game.Projectiles[j].Active = false
			}
		}
	}
}

func (game *Game) spawnProjectile(position rl.Vector2, destination rl.Vector2, color rl.Color, radius float32, speed float32, damage float32, projectileRange float32) {
	if game.CurrentProjectiles >= maxProjectiles {
		game.CurrentProjectiles = 0
	}

	finalDestination := rl.Vector2Add(rl.Vector2Scale(rl.Vector2Normalize(rl.Vector2Subtract(destination, position)), float32(projectileRange)), position)

	game.Projectiles = append(game.Projectiles,
		Projectile{
			Walker: Walker{Pos: position, Speed: speed, Destination: finalDestination, Walking: true},
			Radius: radius,
			Active: true,
			Color:  color,
			Damage: damage,
			Range:  projectileRange,
			TeamID: game.Players[0].TeamID, // TODO remove hammer
		},
	)
}

func (game *Game) removeInactiveProjectiles() {
	for i := 0; i < len(game.Projectiles); i++ {
		if game.Projectiles[i].Active {
			continue
		}

		// Replace current inactive last item of array
		if i < len(game.Projectiles)-1 {
			game.Projectiles[i] = game.Projectiles[len(game.Projectiles)-1]
		}

		// Now just reslice remove last index
		game.Projectiles = game.Projectiles[:len(game.Projectiles)-1]

		// The same index now holds the previous last item, and it needs to be processed
		i--
	}
}
