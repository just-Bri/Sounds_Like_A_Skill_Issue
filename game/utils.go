// game/utils.go
package game

import (
	"bytes"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

func (g *Game) initMusic() error {
	musicData, err := gameAssets.ReadFile("assets/Clement_Panchout_Revenge.wav")
	if err != nil {
		return fmt.Errorf("failed to load music: %v", err)
	}

	d, err := wav.Decode(g.audioContext, bytes.NewReader(musicData))
	if err != nil {
		return fmt.Errorf("failed to decode music: %v", err)
	}

	g.backgroundMusic, err = g.audioContext.NewPlayer(d)
	if err != nil {
		return fmt.Errorf("failed to create music player: %v", err)
	}

	g.backgroundMusic.Play()
	return nil
}

func (g *Game) resetGame() {
	g.player.x = float64(ScreenWidth / 2)
	g.player.y = float64(ScreenHeight / 2)
	g.player.alive = true
	g.player.score = 0
	g.projectiles = nil
	g.timer = 0
	g.difficultyTimer = 0
	ProjectileSpawnRate = InitialProjectileSpawnRate
	g.isPaused = false
	g.gameStartTime = float64(ebiten.ActualTPS())
	g.scoreLogged = false
}

func (g *Game) cleanup() {
	if g.backgroundMusic != nil {
		g.backgroundMusic.Close()
	}
}

func loadSpritesheet(path string) (*ebiten.Image, error) {
	spriteSheetData, err := gameAssets.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load spritesheet: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(spriteSheetData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode spritesheet: %v", err)
	}

	return ebiten.NewImageFromImage(img), nil
}
