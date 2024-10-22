// game/player.go
package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DefaultPlayerSpeed  = 200
	DefaultPlayerRadius = 7
	deadZone            = 0.15 // Ignore small analog stick movements
)

type Player struct {
	x                float64
	y                float64
	speed            float64
	spritesheet      *ebiten.Image
	radius           float64
	alive            bool
	width            float64
	height           float64
	name             string
	frameWidth       int
	frameHeight      int
	currentFrame     int
	framesPerRow     int
	animationSpeed   float64
	animationTimer   float64
	currentDirection int
	score            int
	game             *Game
}

const (
	dirUp int = iota
	dirDown
	dirLeft
	dirRight
	dirUpRight
	dirDownRight
	dirDownLeft
	dirUpLeft
)

// getFirstGamepadID returns the first connected gamepad ID or -1 if none
func getFirstGamepadID() ebiten.GamepadID {
	ids := make([]ebiten.GamepadID, 0)
	ids = ebiten.AppendGamepadIDs(ids)
	if len(ids) > 0 {
		return ids[0]
	}
	return -1
}

func NewPlayer(x, y float64, spritesheet *ebiten.Image) *Player {
	return &Player{
		x:                x,
		y:                y,
		speed:            DefaultPlayerSpeed,
		spritesheet:      spritesheet,
		frameWidth:       13,
		frameHeight:      17,
		currentFrame:     0,
		framesPerRow:     4,
		animationSpeed:   0.25,
		animationTimer:   0,
		currentDirection: dirDown,
		radius:           DefaultPlayerRadius,
		alive:            true,
		width:            13,
		height:           17,
		name:             "",
		score:            0,
	}
}

func (p *Player) Update() {
	if !p.alive {
		return
	}

	dt := 1.0 / 60.0

	dx, dy := p.getInput()

	if dx != 0 || dy != 0 {
		p.x += dx * p.speed * dt
		p.y += dy * p.speed * dt

		p.animationTimer += dt
		if p.animationTimer >= p.animationSpeed {
			p.currentFrame = (p.currentFrame + 1) % p.framesPerRow
			p.animationTimer = 0
		}

		// Update direction based on movement
		p.updateDirection(dx, dy)
	}

	// Keep player within screen bounds
	p.x = math.Max(0, math.Min(float64(ScreenWidth-p.width), p.x))
	p.y = math.Max(0, math.Min(float64(ScreenHeight-p.height), p.y))
}

func (p *Player) Draw(screen *ebiten.Image) {
	if !p.alive {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Calculate source rectangle based on current frame and direction
	sx := float64(p.currentFrame * p.frameWidth)
	sy := float64(p.currentDirection * p.frameHeight)

	// Set destination position
	op.GeoM.Translate(p.x, p.y)

	// Draw the current frame
	screen.DrawImage(p.spritesheet.SubImage(image.Rect(
		int(sx),
		int(sy),
		int(sx)+p.frameWidth,
		int(sy)+p.frameHeight,
	)).(*ebiten.Image), op)
}

func (p *Player) getInput() (float64, float64) {
	var dx, dy float64

	// Keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx++
	}

	// Gamepad input
	gamepadID := getFirstGamepadID()
	if gamepadID >= 0 && ebiten.IsStandardGamepadLayoutAvailable(gamepadID) {
		// Left analog stick
		stickX := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal)
		stickY := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickVertical)

		// Apply deadzone
		if math.Abs(stickX) > deadZone {
			dx = stickX
		}
		if math.Abs(stickY) > deadZone {
			dy = stickY
		}

		// D-pad support
		if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftTop) {
			dy--
		}
		if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftBottom) {
			dy++
		}
		if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftLeft) {
			dx--
		}
		if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftRight) {
			dx++
		}
	}

	// Normalize diagonal movement
	if dx != 0 && dy != 0 {
		length := math.Sqrt(dx*dx + dy*dy)
		dx /= length
		dy /= length
	}

	return dx, dy
}

func (p *Player) updateDirection(dx, dy float64) {
	switch {
	case dx == 0 && dy < 0:
		p.currentDirection = dirUp
	case dx == 0 && dy > 0:
		p.currentDirection = dirDown
	case dx < 0 && dy == 0:
		p.currentDirection = dirLeft
	case dx > 0 && dy == 0:
		p.currentDirection = dirRight
	case dx > 0 && dy < 0:
		p.currentDirection = dirUpRight
	case dx > 0 && dy > 0:
		p.currentDirection = dirDownRight
	case dx < 0 && dy > 0:
		p.currentDirection = dirDownLeft
	case dx < 0 && dy < 0:
		p.currentDirection = dirUpLeft
	}
}
