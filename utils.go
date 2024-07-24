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

func ResetGame(player *Player) {
	player.X = 640
	player.Y = 400
	player.Alive = true
	player.Score = 0
	Projectiles = nil
	Timer = 0
	difficultyTimer = 0
	ProjectileSpawnRate = InitialProjectileSpawnRate
	isPaused = false
	gameStartTime = rl.GetTime()
	ScoreLogged = false
}

func HandleNameInput(player *Player) {
	key := rl.GetCharPressed()

	if IsCharacterValid(key) && len(player.Name) < 3 {
		player.Name += string(key)
	}

	if rl.IsKeyPressed(rl.KeyBackspace) && len(player.Name) > 0 {
		player.Name = player.Name[:len(player.Name)-1]
	}

	if rl.IsKeyReleased(rl.KeyEnter) && len(player.Name) == 3 {
		err := LogHighscore(player)
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

func UnloadGame(player *Player) {
	if rl.IsMusicStreamPlaying(backgroundMusic) {
		rl.StopMusicStream(backgroundMusic)
	}

	if player != nil {
		player = nil
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

func LoadSpritesheet(spriteFile string) rl.Texture2D {
	spriteSheetData, err := gameAssets.ReadFile(spriteFile)
	if err != nil {
		panic(err)
	}
	spriteSheet := rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", spriteSheetData, int32(len(spriteSheetData))))
	return spriteSheet
}
