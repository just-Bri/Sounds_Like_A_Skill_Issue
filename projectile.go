package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/rand"
)

type Projectile struct {
	X      float32
	Y      float32
	DX     float32
	DY     float32
	Radius float32
}

var (
	ProjectileSpeed            float32 = 200
	InitialProjectileSpawnRate float32 = 2
	ProjectileSpawnRate        float32 = InitialProjectileSpawnRate
	SpawnTimer                 float32
	ProjectileSprite           rl.Texture2D
)

var Projectiles []Projectile

func InitProjectiles() {
	projectilePngData, err := gameAssets.ReadFile("assets/projectile.png")
	if err != nil {
		panic(err)
	}
	ProjectileSprite = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", projectilePngData, int32(len(projectilePngData))))

	ProjectileSpawnRate = InitialProjectileSpawnRate
}

func DrawProjectile() {
	for _, proj := range Projectiles {
		rl.DrawTexture(ProjectileSprite, int32(proj.X-proj.Radius), int32(proj.Y-proj.Radius), rl.White)
	}
}

func SpawnProjectile(player *Player) {
	side := rand.Intn(4)
	var x, y float32

	switch side {
	case 0: // Top
		x = float32(ScreenWidth)
		y = -50
	case 1: // Right
		x = float32(ScreenWidth) + 50
		y = float32(ScreenHeight)
	case 2: // Bottom
		x = float32(ScreenWidth)
		y = float32(ScreenHeight) + 50
	case 3: // Left
		x = -50
		y = float32(ScreenHeight)
	}

	dx := player.X - x
	dy := player.Y - y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	dx /= length
	dy /= length

	direction := rl.Vector2Normalize(rl.Vector2{X: player.X - x, Y: player.Y - y})
	Projectiles = append(Projectiles, Projectile{
		X:      x,
		Y:      y,
		DX:     direction.X * ProjectileSpeed,
		DY:     direction.Y * ProjectileSpeed,
		Radius: 7,
	})
}

func UpdateProjectiles(dt float32) {
	for i := 0; i < len(Projectiles); i++ {
		proj := &Projectiles[i]
		proj.X += proj.DX * dt
		proj.Y += proj.DY * dt

		playerRect := rl.NewRectangle(player.X, player.Y, player.Width, player.Height)

		// Check collision with player
		if rl.CheckCollisionCircleRec(rl.Vector2{X: proj.X, Y: proj.Y}, proj.Radius, playerRect) {
			player.Alive = false
			break
		}
	}
}

func RemoveOffscreenProjectiles(player *Player) {
	for i := len(Projectiles) - 1; i >= 0; i-- {
		proj := Projectiles[i]
		if proj.X < -50 || proj.X > float32(ScreenWidth)+50 || proj.Y < -50 || proj.Y > float32(ScreenHeight)+50 {
			Projectiles = append(Projectiles[:i], Projectiles[i+1:]...)
			if player.Alive {
				player.Score += int(float32(ProjectileSpawnRate) / InitialProjectileSpawnRate)
			}
		}
	}
}
