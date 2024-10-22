// game/main.go
package game

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*
var gameAssets embed.FS

type Game struct {
	player           *Player
	projectiles      []*Projectile
	isPaused         bool
	timer            float64
	difficultyTimer  float64
	spawnTimer       float64
	showHighScores   bool
	enteringName     bool
	gameStartTime    float64
	backgroundMusic  *audio.Player
	scoreLogged      bool
	audioContext     *audio.Context
	spritesheet      *ebiten.Image
	projectileSprite *ebiten.Image
	fonts            *Fonts
}

func NewGame() (*Game, error) {
	g := &Game{}

	// Initialize audio context
	g.audioContext = audio.NewContext(44100)

	// Load fonts
	fonts, err := LoadFonts()
	if err != nil {
		return nil, err
	}
	g.fonts = fonts

	// Load spritesheet
	spriteSheetData, err := gameAssets.ReadFile("assets/player_spritesheet.png")
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(spriteSheetData))
	if err != nil {
		return nil, err
	}
	g.spritesheet = ebiten.NewImageFromImage(img)

	// Load projectile sprite
	projectileData, err := gameAssets.ReadFile("assets/projectile.png")
	if err != nil {
		return nil, err
	}
	projImg, _, err := image.Decode(bytes.NewReader(projectileData))
	if err != nil {
		return nil, err
	}
	g.projectileSprite = ebiten.NewImageFromImage(projImg)

	// Initialize player
	g.player = NewPlayer(float64(ScreenWidth/2), float64(ScreenHeight/2), g.spritesheet)
	g.player.game = g

	// Initialize music
	if err := g.initMusic(); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) isStartPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return true
	}

	gamepadID := getFirstGamepadID()
	if gamepadID >= 0 && ebiten.IsStandardGamepadLayoutAvailable(gamepadID) {
		return inpututil.IsStandardGamepadButtonJustPressed(gamepadID, ebiten.StandardGamepadButtonCenterRight)
	}
	return false
}

func (g *Game) isSelectPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return true
	}

	gamepadID := getFirstGamepadID()
	if gamepadID >= 0 && ebiten.IsStandardGamepadLayoutAvailable(gamepadID) {
		return inpututil.IsStandardGamepadButtonJustPressed(gamepadID, ebiten.StandardGamepadButtonCenterLeft)
	}
	return false
}

func (g *Game) isQuitPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return true
	}

	gamepadID := getFirstGamepadID()
	if gamepadID >= 0 && ebiten.IsStandardGamepadLayoutAvailable(gamepadID) {
		return inpututil.IsStandardGamepadButtonJustPressed(gamepadID, ebiten.StandardGamepadButtonRightRight)
	}
	return false
}

func (g *Game) isHighScorePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		return true
	}

	gamepadID := getFirstGamepadID()
	if gamepadID >= 0 && ebiten.IsStandardGamepadLayoutAvailable(gamepadID) {
		return inpututil.IsStandardGamepadButtonJustPressed(gamepadID, ebiten.StandardGamepadButtonFrontTopRight)
	}
	return false
}

func (g *Game) Update() error {
	if (g.isPaused || !g.player.alive) && g.isQuitPressed() {
		return fmt.Errorf("quit game")
	}

	if !g.player.alive && g.isSelectPressed() {
		g.resetGame()
	}

	if g.player.alive && g.isStartPressed() {
		g.isPaused = !g.isPaused
	}

	if g.isPaused {
		return nil
	}

	if !g.player.alive {
		return nil
	}

	dt := 1.0 / 60.0
	g.timer += dt

	// Difficulty increase logic
	g.difficultyTimer += dt
	if g.difficultyTimer >= 20 {
		ProjectileSpawnRate *= 2
		g.difficultyTimer -= 20
		if ProjectileSpawnRate > 32 {
			ProjectileSpawnRate = 32
		}
	}

	// Update game state
	if !g.isPaused && !g.showHighScores && g.player.alive {
		g.player.Update()
		g.updateProjectiles(dt)
		g.removeOffscreenProjectiles()
	}

	// Spawn new projectiles
	g.spawnTimer += dt
	if g.spawnTimer >= 1/ProjectileSpawnRate {
		g.spawnProjectile()
		g.spawnTimer -= 1 / ProjectileSpawnRate
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// The screen is automatically cleared in Ebitengine

	g.drawScore(screen)
	g.drawProjectiles(screen)
	g.player.Draw(screen)

	currentTime := float64(ebiten.ActualTPS())
	if currentTime-g.gameStartTime < 5 {
		g.drawIntro(screen)
	}

	if !g.player.alive {
		g.drawGameOverScreen(screen)
	}

	if g.isPaused && !g.showHighScores {
		g.drawPauseScreen(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
