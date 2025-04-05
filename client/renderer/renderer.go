package renderer

import (
	"fmt"
	"open-moba/client/game"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Renderer struct {
	ScreenWidth  int32
	ScreenHeight int32

	Background rl.Texture2D
}

func NewRenderer(screenWidth, screenHeight int32) Renderer {
	return Renderer{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,

		Background: rl.LoadTextureFromImage(rl.GenImagePerlinNoise(int(screenWidth), int(screenHeight), 50, 50, 0.5)),
	}
}

func (r *Renderer) DrawGame(g *game.Game) {
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

func (r *Renderer) DrawHUD(g *game.Game) {
	var abilitySquareLength int32 = 50
	var fontSize int32 = 32
	var padding int32 = 4
	backGroundColor := rl.Black
	squareCooldownColor := rl.LightGray
	squareActiveColor := rl.Green
	textColor := rl.Black

	abilitiesHUDWidth := abilitySquareLength * int32(len(g.Players[g.PlayerIndex].Abilities))

	posX := r.ScreenWidth/2 - abilitiesHUDWidth/2
	posY := r.ScreenHeight - abilitySquareLength

	rl.DrawRectangle(posX, posY, abilitiesHUDWidth, abilitySquareLength, backGroundColor)

	now := time.Now()

	for i, a := range g.Players[g.PlayerIndex].Abilities {
		squareColor := squareActiveColor

		remainingCooldown := now.Sub(a.CastedAt)
		isInCooldown := remainingCooldown < a.Cooldown
		if isInCooldown {
			squareColor = squareCooldownColor
		}

		rl.DrawRectangle(posX+int32(i)*abilitySquareLength+padding, posY+padding, abilitySquareLength-padding*2, abilitySquareLength-padding*2, squareColor)

		if isInCooldown {
			text := fmt.Sprintf("%d", int32((a.Cooldown-remainingCooldown).Seconds()+1))
			textWidth := rl.MeasureText(text, fontSize)

			rl.DrawText(text, posX+int32(i)*abilitySquareLength+abilitySquareLength/2-textWidth/2, r.ScreenHeight-abilitySquareLength/2-fontSize/2+padding/2, fontSize, textColor)
		}
	}
}

func (r *Renderer) DrawBackground() {
	rl.DrawTexture(r.Background, 0, 0, rl.NewColor(200, 255, 200, 255))
}
