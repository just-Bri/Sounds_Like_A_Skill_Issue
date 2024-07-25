package main

import (
	"embed"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/player_spritesheet.png assets/projectile.png assets/Clement_Panchout_Revenge.wav
var gameAssets embed.FS

var (
	isPaused           bool
	Timer              float32
	difficultyTimer    float32
	difficultyInterval float32 = 20
	spawnTimer         float32
	showHighScores     bool
	enteringName       bool
	gameStartTime      float64
	backgroundMusic    rl.Music
	ScoreLogged        bool
)

var player *Player
var Projectiles []Projectile

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, GameName)
	rl.SetExitKey(0)
	defer UnloadGame(player)

	rl.InitAudioDevice()

	rl.SetTargetFPS(60)

	InitMusic()
	InitProjectiles()

	playerOneSpriteSheet := LoadSpritesheet("assets/player_spritesheet.png")
	player = NewPlayer(ScreenWidth/2, ScreenHeight/2, playerOneSpriteSheet)

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(backgroundMusic)
		update()
		draw()

		if isPaused && !showHighScores && (rl.IsKeyPressed(rl.KeyQ) || rl.IsGamepadButtonPressed(0, rl.GamepadButtonRightFaceRight)) {
			break
		}
		if !player.Alive && !enteringName && (rl.IsKeyPressed(rl.KeyQ) || rl.IsGamepadButtonPressed(0, rl.GamepadButtonRightFaceRight)) {
			break
		}
	}
}

func update() {
	if enteringName {
		HandleNameInput(player)
		return
	}

	if !player.Alive && rl.IsKeyReleased(rl.KeyL) && !ScoreLogged {
		enteringName = true
		player.Name = ""
		return
	}

	if rl.IsKeyReleased(rl.KeyR) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonMiddleLeft) && !enteringName && !showHighScores {
		ResetGame(player)
	}

	if (rl.IsKeyReleased(rl.KeyEscape) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonMiddleRight)) && !enteringName && !showHighScores {
		isPaused = !isPaused
	}

	if player.Alive && (!isPaused || showHighScores) && (rl.IsKeyReleased(rl.KeyH) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonRightFaceUp)) {
		showHighScores = !showHighScores
		isPaused = !isPaused
		if showHighScores {
			var supaErr = FetchHighscores()
			if supaErr != nil {
				fmt.Println("Error fetching high scores:", supaErr)
			}
		}
	}

	if isPaused {
		return
	}

	if !player.Alive {
		return
	}

	dt := rl.GetFrameTime()
	Timer += dt

	// Difficulty increase logic
	difficultyTimer += dt
	if difficultyTimer >= difficultyInterval {
		ProjectileSpawnRate *= 2
		difficultyTimer -= difficultyInterval
		ProjectileSpawnRate = float32(math.Min(float64(ProjectileSpawnRate), 32))
	}

	// Only Update when player is alive and not paused/highscores
	if !isPaused && !showHighScores && player.Alive {
		UpdatePlayer(player)
		UpdateProjectiles(dt)
		RemoveOffscreenProjectiles(player)
	}

	// Spawn new Projectiles
	spawnTimer += dt
	if spawnTimer >= 1/ProjectileSpawnRate {
		SpawnProjectile(player)
		spawnTimer -= 1 / ProjectileSpawnRate
	}
}

func draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	DrawScore(player)
	DrawProjectile()
	DrawPlayer(player)

	// Display instruction for pause/highscore
	currentTime := rl.GetTime()
	if currentTime-gameStartTime < 5 {
		DrawIntro()
	}

	if !player.Alive {
		DrawGameOverScreen(player.Score)
	}

	if enteringName {
		DrawNameInputScreen(player)
	}

	if showHighScores {
		DrawHighScoresScreen()
	}

	if isPaused && !showHighScores {
		DrawPauseScreen()
	}
}
