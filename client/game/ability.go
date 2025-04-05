package game

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AbilityType string

const _castSingleProjectile = "CAST_SINGLE_PROJECTILE"
const _castThreeProjectiles = "CAST_THREE_PROJECTILES"
const _teleport = "TELEPORT"

type PositionType string

const _playerPosition = "PLAYER_POSITION"
const _mousePosition = "MOUSE_POSITION"

type Ability struct {
	AbilityType AbilityType
	Color       rl.Color

	CastPositionType        PositionType
	DestinationPositionType PositionType

	Cooldown time.Duration
	// CastingTime time.Duration // time that it takes since player activates ability until and it actually casts

	Radius float32
	Speed  float32
	Range  float32
	Damage float32

	CastedAt time.Time
}
