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
	PlayerInstance     *Player
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

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, GameName)
	rl.SetExitKey(0)
	defer UnloadGame()

	rl.InitAudioDevice()

	rl.SetTargetFPS(60)

	InitMusic()
	InitProjectiles()
	InitPlayer()

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(backgroundMusic)
		update()
		draw()

		if isPaused && !showHighScores && rl.IsKeyPressed(rl.KeyQ) {
			break
		}
		if !PlayerInstance.Alive && !enteringName && rl.IsKeyPressed(rl.KeyQ) {
			break
		}
	}
}

func update() {
	if enteringName {
		HandleNameInput()
		return
	}

	if !PlayerInstance.Alive && rl.IsKeyReleased(rl.KeyL) && !ScoreLogged {
		enteringName = true
		PlayerInstance.Name = ""
		return
	}

	if rl.IsKeyReleased(rl.KeyR) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonMiddleLeft) && !enteringName && !showHighScores {
		ResetGame()
	}

	if (rl.IsKeyReleased(rl.KeyEscape) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonMiddleRight)) && !enteringName && !showHighScores {
		isPaused = !isPaused
	}

	if PlayerInstance.Alive && (!isPaused || showHighScores) && (rl.IsKeyReleased(rl.KeyH) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonRightFaceUp)) {
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

	if !PlayerInstance.Alive {
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
	if !isPaused && !showHighScores && PlayerInstance.Alive {
		UpdatePlayer(rl.GetFrameTime(), ScreenWidth, ScreenHeight, PlayerInstance.Alive, enteringName)
		UpdateProjectiles(dt)
		RemoveOffscreenProjectiles()
	}

	// Spawn new projectiles
	spawnTimer += dt
	if spawnTimer >= 1/ProjectileSpawnRate {
		SpawnProjectile()
		spawnTimer -= 1 / ProjectileSpawnRate
	}
}

func draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	DrawScore()
	DrawProjectile()
	DrawPlayer()

	// Display instruction for pause/highscore
	currentTime := rl.GetTime()
	if currentTime-gameStartTime < 5 {
		DrawIntro()
	}

	if !PlayerInstance.Alive {
		DrawGameOverScreen(CurrentScore)
	}

	if enteringName {
		DrawNameInputScreen()
	}

	if showHighScores {
		DrawHighScoresScreen()
	}

	if isPaused && !showHighScores {
		DrawPauseScreen()
	}
}
