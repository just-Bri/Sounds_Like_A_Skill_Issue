package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DefaultPlayerSpeed  float32 = 200
	DefaultPlayerRadius float32 = 7
)

type Player struct {
	X                float32
	Y                float32
	Speed            float32
	Sprite           rl.Texture2D
	Radius           float32
	Alive            bool
	Width            float32
	Height           float32
	Name             string
	SpriteSheet      rl.Texture2D
	FrameWidth       int32
	FrameHeight      int32
	CurrentFrame     int32
	FramesPerRow     int32
	AnimationSpeed   float32
	AnimationTimer   float32
	CurrentDirection int
	Score            int
}

const (
	up        int = iota // 0
	down                 // 1
	left                 // 2
	right                // 3
	upRight              // 4
	downRight            // 5
	downLeft             // 6
	upLeft               // 7
)

func NewPlayer(x, y float32, spriteSheet rl.Texture2D) *Player {
	return &Player{
		X:                x,
		Y:                y,
		Speed:            DefaultPlayerSpeed,
		SpriteSheet:      spriteSheet,
		FrameWidth:       13,
		FrameHeight:      17,
		CurrentFrame:     0,
		FramesPerRow:     4,
		AnimationSpeed:   0.25,
		AnimationTimer:   0,
		CurrentDirection: down,
		Radius:           DefaultPlayerRadius,
		Alive:            true,
		Width:            13, // sprite width
		Height:           17, // sprite height
		Name:             "",
		Score:            0,
	}
}

func UpdatePlayer(player *Player) {
	if !player.Alive {
		return
	}

	dt := rl.GetFrameTime()

	playerX, playerY := getPlayerInput()

	if playerX != 0 || playerY != 0 {
		player.X += playerX * player.Speed * dt
		player.Y += playerY * player.Speed * dt

		player.AnimationTimer += dt
		if player.AnimationTimer >= player.AnimationSpeed {
			player.CurrentFrame = (player.CurrentFrame + 1) % player.FramesPerRow
			player.AnimationTimer = 0
		}

		// right and down are positive
		if playerX == 0 && playerY < 0 {
			player.CurrentDirection = up
		}
		if playerX == 0 && playerY > 0 {
			player.CurrentDirection = down
		}
		if playerX < 0 && playerY == 0 {
			player.CurrentDirection = left
		}
		if playerX > 0 && playerY == 0 {
			player.CurrentDirection = right
		}
		if playerX > 0 && playerY < 0 {
			player.CurrentDirection = upRight
		}
		if playerX > 0 && playerY > 0 {
			player.CurrentDirection = downRight
		}
		if playerX < 0 && playerY > 0 {
			player.CurrentDirection = downLeft
		}
		if playerX < 0 && playerY < 0 {
			player.CurrentDirection = upLeft
		}
	}

	// Keep player within screen bounds
	player.X = float32(math.Max(0, math.Min(float64(player.X), float64(ScreenWidth-player.Width))))
	player.Y = float32(math.Max(0, math.Min(float64(player.Y), float64(ScreenHeight-player.Height))))
}

func DrawPlayer(player *Player) {
	if player.Alive {
		var row int32
		switch player.CurrentDirection {
		case up:
			row = 0
		case down:
			row = 1
		case left:
			row = 2
		case right:
			row = 3
		case upRight:
			row = 4
		case downRight:
			row = 5
		case downLeft:
			row = 6
		case upLeft:
			row = 7
		}

		sourceRect := rl.NewRectangle(float32(player.CurrentFrame*player.FrameWidth), float32(row*player.FrameHeight), float32(player.FrameWidth), float32(player.FrameHeight))
		destRect := rl.NewRectangle(player.X, player.Y, player.Width, player.Height)
		rl.DrawTexturePro(player.SpriteSheet, sourceRect, destRect, rl.NewVector2(0, 0), 0, rl.White)
	}
}

var getPlayerInput = func() (float32, float32) {
	var dx, dy float32

	if !enteringName {
		if rl.IsKeyDown(rl.KeyW) {
			dy--
		}
		if rl.IsKeyDown(rl.KeyS) {
			dy++
		}
		if rl.IsKeyDown(rl.KeyA) {
			dx--
		}
		if rl.IsKeyDown(rl.KeyD) {
			dx++
		}

		if rl.IsGamepadAvailable(0) {
			leftX := rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftX)
			leftY := rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftY)

			if math.Abs(float64(leftX)) > 0.2 {
				dx += leftX
			}
			if math.Abs(float64(leftY)) > 0.2 {
				dy += leftY
			}

			// D-pad
			if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceUp) {
				dy--
			}
			if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceDown) {
				dy++
			}
			if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceLeft) {
				dx--
			}
			if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceRight) {
				dx++
			}
		}
	}

	// Normalize diagonal movement
	if dx != 0 && dy != 0 {
		length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		dx /= length
		dy /= length
	}

	return dx, dy
}
