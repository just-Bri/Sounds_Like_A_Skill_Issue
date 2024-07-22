// File: utils.go

package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func InitMusic() {
	musicData, err := gameAssets.ReadFile("assets/Clement_Panchout_Revenge.wav") // or .ogg
	if err != nil {
		panic(err)
	}
	backgroundMusic := rl.LoadMusicStreamFromMemory(".wav", musicData, int32(len(musicData)))

	rl.PlayMusicStream(backgroundMusic)
}

func ResetGame() {
	PlayerInstance.X = 640
	PlayerInstance.Y = 400
	PlayerInstance.Alive = true
	Projectiles = nil
	CurrentScore = 0
	timer = 0
	difficultyTimer = 0
	ProjectileSpawnRate = InitialProjectileSpawnRate
	isPaused = false
	gameStartTime = rl.GetTime()
	ScoreLogged = false
}

func HandleNameInput() {
	key := rl.GetCharPressed()

	if IsCharacterValid(key) && len(PlayerInstance.Name) < 3 {
		PlayerInstance.Name += string(key)
	}

	if rl.IsKeyPressed(rl.KeyBackspace) && len(PlayerInstance.Name) > 0 {
		PlayerInstance.Name = PlayerInstance.Name[:len(PlayerInstance.Name)-1]
	}

	if rl.IsKeyReleased(rl.KeyEnter) && len(PlayerInstance.Name) == 3 {
		err := LogHighscore()
		if err != nil {
			fmt.Println("(TR) Error inserting highscore:", err)
		} else {
			fmt.Println("(TR) Highscore inserted successfully!")
		}
		enteringName = false
	}
}

func IsCharacterValid(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '_'
}

func UnloadGame() {
	if rl.IsMusicStreamPlaying(backgroundMusic) {
		rl.StopMusicStream(backgroundMusic)
	}

	if PlayerInstance != nil {
		PlayerInstance = nil
	}

	if ProjectileSprite.ID != 0 {
		rl.UnloadTexture(ProjectileSprite)
		ProjectileSprite = rl.Texture2D{}
	}

	if backgroundMusic.CtxData != nil {
		rl.UnloadMusicStream(backgroundMusic)
		backgroundMusic = rl.Music{}
	}

	Projectiles = nil
	time.Sleep(100 * time.Millisecond)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
