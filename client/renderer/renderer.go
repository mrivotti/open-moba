package renderer

import (
	"open-moba/client/game"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawGame(g *game.Game) {
	for i := range g.Projectiles {
		drawProjectile(&g.Projectiles[i])
	}
	for i := range g.Players {
		drawPlayer(&g.Players[i])
	}
}

func drawPlayer(p *game.Player) {
	if !p.Active {
		return
	}

	rl.DrawCircle(int32(p.Walker.Pos.X), int32(p.Walker.Pos.Y), p.Radius, p.Color)
	drawHealthBar(p)
}

func drawHealthBar(p *game.Player) {
	healthPercentage := p.Health / p.MaxHealth
	maxBarWidth := p.Radius * 3
	barWidth := maxBarWidth * healthPercentage
	barHeight := int32(15)

	rl.DrawRectangle(int32(p.Walker.Pos.X)-int32(maxBarWidth/2), int32(p.Walker.Pos.Y+p.Radius*1.5), int32(barWidth), barHeight, rl.Red)
}

func drawProjectile(p *game.Projectile) {
	if !p.Active {
		return
	}

	rl.DrawCircle(int32(p.Walker.Pos.X), int32(p.Walker.Pos.Y), p.Radius, p.Color)
}
