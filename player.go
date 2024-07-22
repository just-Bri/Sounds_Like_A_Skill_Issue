// File: player.go

package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DefaultSpeed  float32 = 200
	DefaultRadius float32 = 7
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
}

const (
	Up        int = iota // 0
	Down                 // 1
	Left                 // 2
	Right                // 3
	UpRight              // 4
	DownRight            // 5
	DownLeft             // 6
	UpLeft               // 7
)

type PlayerMovement struct {
	Speed float32
}

var CurrentScore int = 0

func NewPlayer(x, y float32, spriteSheet rl.Texture2D) *Player {
	return &Player{
		X:                x,
		Y:                y,
		Speed:            DefaultSpeed,
		SpriteSheet:      spriteSheet,
		FrameWidth:       13,
		FrameHeight:      17,
		CurrentFrame:     0,
		FramesPerRow:     4,
		AnimationSpeed:   0.25,
		AnimationTimer:   0,
		CurrentDirection: Down,
		Radius:           DefaultRadius,
		Alive:            true,
		Width:            13, // sprite width
		Height:           17, // sprite height
		Name:             "",
	}
}

func InitPlayer() {
	spriteSheetData, err := gameAssets.ReadFile("assets/player_spritesheet.png")
	if err != nil {
		panic(err)
	}
	spriteSheet := rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", spriteSheetData, int32(len(spriteSheetData))))

	PlayerInstance = NewPlayer(ScreenWidth/2, ScreenHeight/2, spriteSheet)
}

func UpdatePlayer(dt, screenWidth float32, screenHeight float32, alive bool, enteringName bool) {
	if !PlayerInstance.Alive {
		return
	}

	playerX, playerY := getInput(alive, enteringName)

	if playerX != 0 || playerY != 0 {
		PlayerInstance.X += playerX * PlayerInstance.Speed * dt
		PlayerInstance.Y += playerY * PlayerInstance.Speed * dt

		PlayerInstance.AnimationTimer += dt
		if PlayerInstance.AnimationTimer >= PlayerInstance.AnimationSpeed {
			PlayerInstance.CurrentFrame = (PlayerInstance.CurrentFrame + 1) % PlayerInstance.FramesPerRow
			PlayerInstance.AnimationTimer = 0
		}

		// right and down are positive
		if playerX == 0 && playerY < 0 {
			PlayerInstance.CurrentDirection = Up
		}
		if playerX == 0 && playerY > 0 {
			PlayerInstance.CurrentDirection = Down
		}
		if playerX < 0 && playerY == 0 {
			PlayerInstance.CurrentDirection = Left
		}
		if playerX > 0 && playerY == 0 {
			PlayerInstance.CurrentDirection = Right
		}
		if playerX > 0 && playerY < 0 {
			PlayerInstance.CurrentDirection = UpRight
		}
		if playerX > 0 && playerY > 0 {
			PlayerInstance.CurrentDirection = DownRight
		}
		if playerX < 0 && playerY > 0 {
			PlayerInstance.CurrentDirection = DownLeft
		}
		if playerX < 0 && playerY < 0 {
			PlayerInstance.CurrentDirection = UpLeft
		}
	}

	// Keep player within screen bounds
	PlayerInstance.X = float32(math.Max(0, math.Min(float64(PlayerInstance.X), float64(screenWidth-PlayerInstance.Width))))
	PlayerInstance.Y = float32(math.Max(0, math.Min(float64(PlayerInstance.Y), float64(screenHeight-PlayerInstance.Height))))
}

func DrawPlayer() {
	if PlayerInstance.Alive {
		var row int32
		switch PlayerInstance.CurrentDirection {
		case Up:
			row = 0
		case Down:
			row = 1
		case Left:
			row = 2
		case Right:
			row = 3
		case UpRight:
			row = 4
		case DownRight:
			row = 5
		case DownLeft:
			row = 6
		case UpLeft:
			row = 7
		}

		sourceRect := rl.NewRectangle(float32(PlayerInstance.CurrentFrame*PlayerInstance.FrameWidth), float32(row*PlayerInstance.FrameHeight), float32(PlayerInstance.FrameWidth), float32(PlayerInstance.FrameHeight))
		destRect := rl.NewRectangle(PlayerInstance.X, PlayerInstance.Y, PlayerInstance.Width, PlayerInstance.Height)
		rl.DrawTexturePro(PlayerInstance.SpriteSheet, sourceRect, destRect, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func getInput(alive bool, enteringName bool) (float32, float32) {
	var dx, dy float32

	if alive && !enteringName {
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
