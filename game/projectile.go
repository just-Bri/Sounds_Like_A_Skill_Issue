// game/projectile.go
package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	x      float64
	y      float64
	dx     float64
	dy     float64
	radius float64
}

var (
	ProjectileSpeed            float64 = 200
	InitialProjectileSpawnRate float64 = 2
	ProjectileSpawnRate        float64 = InitialProjectileSpawnRate
)

func (g *Game) drawProjectiles(screen *ebiten.Image) {
	for _, proj := range g.projectiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(proj.x-proj.radius, proj.y-proj.radius)
		screen.DrawImage(g.projectileSprite, op)
	}
}

func (g *Game) spawnProjectile() {
	side := rand.Intn(4)
	var x, y float64

	switch side {
	case 0: // Top
		x = float64(rand.Intn(ScreenWidth))
		y = -50
	case 1: // Right
		x = float64(ScreenWidth) + 50
		y = float64(rand.Intn(ScreenHeight))
	case 2: // Bottom
		x = float64(rand.Intn(ScreenWidth))
		y = float64(ScreenHeight) + 50
	case 3: // Left
		x = -50
		y = float64(rand.Intn(ScreenHeight))
	}

	dx := g.player.x - x
	dy := g.player.y - y
	length := math.Sqrt(dx*dx + dy*dy)
	dx /= length
	dy /= length

	g.projectiles = append(g.projectiles, &Projectile{
		x:      x,
		y:      y,
		dx:     dx * ProjectileSpeed,
		dy:     dy * ProjectileSpeed,
		radius: 7,
	})
}

func (g *Game) updateProjectiles(dt float64) {
	for _, proj := range g.projectiles {
		proj.x += proj.dx * dt
		proj.y += proj.dy * dt

		// Check collision with player using rectangular collision
		if g.checkCollision(proj) {
			g.player.alive = false
			break
		}
	}
}

func (g *Game) checkCollision(proj *Projectile) bool {
	// Simple circle vs rectangle collision
	closestX := math.Max(g.player.x, math.Min(proj.x, g.player.x+g.player.width))
	closestY := math.Max(g.player.y, math.Min(proj.y, g.player.y+g.player.height))

	distanceX := proj.x - closestX
	distanceY := proj.y - closestY

	distanceSquared := (distanceX * distanceX) + (distanceY * distanceY)
	return distanceSquared < (proj.radius * proj.radius)
}

func (g *Game) removeOffscreenProjectiles() {
	var activeProjectiles []*Projectile
	for _, proj := range g.projectiles {
		if proj.x >= -50 && proj.x <= float64(ScreenWidth)+50 &&
			proj.y >= -50 && proj.y <= float64(ScreenHeight)+50 {
			activeProjectiles = append(activeProjectiles, proj)
		} else if g.player.alive {
			g.player.score += int(ProjectileSpawnRate / InitialProjectileSpawnRate)
		}
	}
	g.projectiles = activeProjectiles
}
