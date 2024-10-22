// game/ui.go
package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func (g *Game) drawPauseScreen(screen *ebiten.Image) {
	// Draw semi-transparent overlay
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)

	pausedText := "PAUSED"
	bound, _ := font.BoundString(g.fonts.Large, pausedText)
	x := (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	y := ScreenHeight/2 - 30
	text.Draw(screen, pausedText, g.fonts.Large, x, y, color.White)

	resumeText := "Press ESC to resume (Start on controller)"
	bound, _ = font.BoundString(g.fonts.Small, resumeText)
	x = (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	y = ScreenHeight/2 + 20
	text.Draw(screen, resumeText, g.fonts.Small, x, y, color.White)

	quitText := "Press Q to quit (B on controller)"
	bound, _ = font.BoundString(g.fonts.Small, quitText)
	x = (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	y = ScreenHeight/2 + 40
	text.Draw(screen, quitText, g.fonts.Small, x, y, color.White)
}

func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 178}) // ~0.7 opacity
	screen.DrawImage(overlay, nil)

	gameOverText := "Game Over!"
	bound, _ := font.BoundString(g.fonts.Large, gameOverText)
	x := (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	y := (ScreenHeight - (bound.Max.Y - bound.Min.Y).Ceil()) / 2
	text.Draw(screen, gameOverText, g.fonts.Large, x, y, color.White)

	scoreText := fmt.Sprintf("Final score: %d", g.player.score)
	bound, _ = font.BoundString(g.fonts.Medium, scoreText)
	x = (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	y = ScreenHeight*3/4 - 100
	text.Draw(screen, scoreText, g.fonts.Medium, x, y, color.White)

	restartText := "Press R to restart (Select/Back on controller)"
	bound, _ = font.BoundString(g.fonts.Medium, restartText)
	x = (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	text.Draw(screen, restartText, g.fonts.Medium, x, y+80, color.White)

	quitText := "Press Q to quit (B on controller)"
	bound, _ = font.BoundString(g.fonts.Medium, quitText)
	x = (ScreenWidth - (bound.Max.X - bound.Min.X).Ceil()) / 2
	text.Draw(screen, quitText, g.fonts.Medium, x, y+120, color.White)
}

func (g *Game) drawIntro(screen *ebiten.Image) {
	pauseText := "Press ESC to pause (Start on controller)"
	bound, _ := font.BoundString(g.fonts.Small, pauseText)
	x := ScreenWidth - (bound.Max.X - bound.Min.X).Ceil() - 10
	y := ScreenHeight - (bound.Max.Y - bound.Min.Y).Ceil() - 10
	text.Draw(screen, pauseText, g.fonts.Small, x, y, color.White)
}

func (g *Game) drawScore(screen *ebiten.Image) {
	timerText := fmt.Sprintf("Timer: %ds", int(g.timer))
	text.Draw(screen, timerText, g.fonts.Small, 0, 20, color.White)

	difficultyText := fmt.Sprintf("Difficulty: %dx", int(ProjectileSpawnRate/InitialProjectileSpawnRate))
	text.Draw(screen, difficultyText, g.fonts.Small, 0, 40, color.White)

	scoreText := fmt.Sprintf("Score: %d", g.player.score)
	text.Draw(screen, scoreText, g.fonts.Small, 0, 60, color.White)
}
