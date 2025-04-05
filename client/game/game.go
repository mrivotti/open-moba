package game

import (
	"errors"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const collisionMargin = 5

type Game struct {
	Players     []Player
	Projectiles []Projectile
	PlayerIndex int
}

func NewGame(nPlayers int, playerIndex int) (*Game, error) {
	if playerIndex >= nPlayers || playerIndex < 0 {
		return nil, errors.New("invalid player index")
	}

	abilities := []Ability{
		{
			AbilityType:             _castSingleProjectile,
			Color:                   rl.Red,
			CastPositionType:        _playerPosition,
			DestinationPositionType: _mousePosition,
			Radius:                  15,
			Speed:                   500,
			Range:                   300,
			Damage:                  5,
			Cooldown:                time.Second * 3,
		},
		{
			AbilityType:             _castSingleProjectile,
			Color:                   rl.Orange,
			CastPositionType:        _playerPosition,
			DestinationPositionType: _mousePosition,
			Radius:                  50,
			Speed:                   100,
			Range:                   500,
			Damage:                  30,
			Cooldown:                time.Second * 10,
		},
		{
			AbilityType:             _castThreeProjectiles,
			Color:                   rl.Green,
			CastPositionType:        _playerPosition,
			DestinationPositionType: _mousePosition,
			Radius:                  10,
			Speed:                   700,
			Range:                   150,
			Damage:                  5,
			Cooldown:                time.Second * 4,
		},
		{
			AbilityType:             _teleport,
			Color:                   rl.Black,
			CastPositionType:        _playerPosition,
			DestinationPositionType: _mousePosition,
			Radius:                  0,
			Speed:                   0,
			Range:                   150,
			Damage:                  0,
			Cooldown:                time.Second * 3,
		},
	}

	players := make([]Player, nPlayers)
	players[playerIndex] = Player{
		Walker:    Walker{Pos: rl.Vector2{X: 100, Y: 100}, Speed: 100},
		Abilities: abilities,
		Radius:    20,
		Active:    true,
		Color:     rl.Blue,
		TeamID:    1,
		Health:    100,
		MaxHealth: 100,
	}

	for i := 0; i < nPlayers; i++ {
		if i == playerIndex {
			continue
		}
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

	return &Game{
		Players:     players,
		Projectiles: make([]Projectile, 0, 100),
		PlayerIndex: playerIndex,
	}, nil
}

func (g *Game) HandleInput() {
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		g.Players[0].SetDestination(mousePos)
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		g.castAbility(g.PlayerIndex, 0, mousePos)
	}
	if rl.IsKeyPressed(rl.KeyW) {
		g.castAbility(g.PlayerIndex, 1, mousePos)
	}
	if rl.IsKeyPressed(rl.KeyE) {
		g.castAbility(g.PlayerIndex, 2, mousePos)
	}
	if rl.IsKeyPressed(rl.KeyR) {
		g.castAbility(g.PlayerIndex, 3, mousePos)
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

func (g *Game) castAbility(playerIndex int, abilityIndex int, mousePos rl.Vector2) {
	if playerIndex > len(g.Players) || playerIndex < 0 {
		return
	}
	p := &g.Players[playerIndex]

	if abilityIndex > len(p.Abilities) || abilityIndex < 0 {
		return
	}
	a := &p.Abilities[abilityIndex]

	now := time.Now()

	if a.CastedAt.Add(a.Cooldown).After(now) {
		return
	}

	switch a.AbilityType {
	case _castSingleProjectile:
		g.spawnProjectile(p.Walker.Pos, mousePos, a.Color, a.Radius, a.Speed, a.Range, a.Damage)

	case _castThreeProjectiles:
		color := rl.Green
		radius := a.Radius
		speed := a.Speed
		rrange := a.Range
		damage := a.Damage

		var angle float32 = math.Pi / 4

		g.spawnProjectile(p.Walker.Pos, mousePos, color, radius, speed, rrange, damage)

		leftProjectileDir := rl.Vector2Rotate(rl.Vector2Subtract(mousePos, p.Pos), angle)
		g.spawnProjectile(p.Walker.Pos, rl.Vector2Add(p.Pos, leftProjectileDir), color, radius, speed, rrange, damage)

		rightProjectileDir := rl.Vector2Rotate(rl.Vector2Subtract(mousePos, p.Pos), -angle)
		g.spawnProjectile(p.Walker.Pos, rl.Vector2Add(p.Pos, rightProjectileDir), color, radius, speed, rrange, damage)

	case _teleport:
		g.Players[playerIndex].Pos = rl.Vector2MoveTowards(g.Players[playerIndex].Pos, mousePos, a.Range)
		g.Players[playerIndex].Walking = false
	}

	a.CastedAt = now
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

func (game *Game) spawnProjectile(position rl.Vector2, destination rl.Vector2, color rl.Color, radius float32, speed float32, projectileRange float32, damage float32) {
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
