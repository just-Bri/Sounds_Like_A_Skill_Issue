// File: ui.go

package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawPauseScreen() {
	rl.DrawRectangle(0, 0, ScreenWidth, ScreenHeight, rl.ColorAlpha(rl.Black, 0.5))

	pausedText := "PAUSED"
	pausedFontSize := int32(40)
	pausedWidth := rl.MeasureText(pausedText, pausedFontSize)
	rl.DrawText(pausedText, (ScreenWidth-int32(pausedWidth))/2, ScreenHeight/2-30, pausedFontSize, rl.White)

	resumeText := "Press ESC to resume (Start on controller)"
	resumeFontSize := int32(20)
	resumeWidth := rl.MeasureText(resumeText, resumeFontSize)
	rl.DrawText(resumeText, (ScreenWidth-int32(resumeWidth))/2, ScreenHeight/2+20, resumeFontSize, rl.White)

	quitText := "Press Q to quit (B on controller)"
	quitFontSize := int32(20)
	quitWidth := rl.MeasureText(quitText, quitFontSize)
	rl.DrawText(quitText, (ScreenWidth-int32(quitWidth))/2, ScreenHeight/2+40, quitFontSize, rl.White)
}

func DrawHighScoresScreen() {
	rl.DrawRectangle(0, 0, ScreenWidth, ScreenHeight, rl.ColorAlpha(rl.Black, 0.5))

	highscoreText := "High Scores"
	highscoreFontSize := int32(40)
	highscoreWidth := rl.MeasureText(highscoreText, highscoreFontSize)
	rl.DrawText(highscoreText, (ScreenWidth-int32(highscoreWidth))/2, 100, highscoreFontSize, rl.White)

	columnWidth := int32(ScreenWidth / 6)
	startX := columnWidth + 80

	headerFontSize := int32(25)
	rl.DrawText("Rank", startX, 160, headerFontSize, rl.White)
	rl.DrawText("Name", startX+columnWidth, 160, headerFontSize, rl.White)
	rl.DrawText("Score", startX+columnWidth*2, 160, headerFontSize, rl.White)
	rl.DrawText("Version", startX+columnWidth*3, 160, headerFontSize, rl.White)

	scoresFontSize := int32(20)
	for i, highscore := range Highscores {
		y := int32(200 + i*40)

		rankText := fmt.Sprintf("%d.", i+1)
		rl.DrawText(rankText, startX, y, scoresFontSize, rl.White)

		rl.DrawText(highscore.Name, startX+columnWidth, y, scoresFontSize, rl.White)

		scoreText := fmt.Sprintf("%d", highscore.Score)
		rl.DrawText(scoreText, startX+columnWidth*2, y, scoresFontSize, rl.White)

		rl.DrawText(highscore.Version, startX+columnWidth*3, y, scoresFontSize, rl.White)
	}

	instructionText := "Press H to close (Y on controller)"
	instructionFontSize := int32(20)
	instructionWidth := rl.MeasureText(instructionText, instructionFontSize)
	rl.DrawText(instructionText, (ScreenWidth-int32(instructionWidth))/2, ScreenHeight-50, instructionFontSize, rl.White)
}

func DrawNameInputScreen() {
	rl.DrawRectangle(0, 0, ScreenWidth, ScreenHeight, rl.ColorAlpha(rl.Black, 0.7))

	promptText := "Enter Your Name"
	promptFontSize := int32(40)
	promptWidth := rl.MeasureText(promptText, promptFontSize)
	rl.DrawText(promptText, (ScreenWidth-int32(promptWidth))/2, ScreenHeight/3, promptFontSize, rl.White)

	helperText := "(Only alphanumeric, and _ allowed)"
	helperFontSize := int32(20)
	helperWidth := rl.MeasureText(helperText, helperFontSize)
	rl.DrawText(helperText, (ScreenWidth-int32(helperWidth))/2, ScreenHeight/3+40, helperFontSize, rl.White)

	nameText := "---"
	if len(PlayerInstance.Name) > 0 {
		nameText = PlayerInstance.Name + strings.Repeat("-", 3-len(PlayerInstance.Name))
	}
	nameFontSize := int32(30)
	nameWidth := rl.MeasureText(nameText, nameFontSize)
	rl.DrawText(nameText, (ScreenWidth-int32(nameWidth))/2, ScreenHeight/2, nameFontSize, rl.White)

	instructionText := "Press Enter to submit"
	instructionFontSize := int32(20)
	instructionWidth := rl.MeasureText(instructionText, instructionFontSize)
	rl.DrawText(instructionText, (ScreenWidth-int32(instructionWidth))/2, ScreenHeight*2/3, instructionFontSize, rl.White)
}

func DrawGameOverScreen(score int) {
	rl.DrawRectangle(0, 0, ScreenWidth, ScreenHeight, rl.ColorAlpha(rl.Black, 0.7))

	smallFontSize := int32(30)

	gameOverText := "Game Over!"
	gameOverFontSize := int32(60)
	gameOverWidth := rl.MeasureText(gameOverText, gameOverFontSize)

	rl.DrawText(gameOverText,
		(ScreenWidth-int32(gameOverWidth))/2,
		(ScreenHeight-gameOverFontSize)/2,
		gameOverFontSize, rl.White)

	finalScoreText := fmt.Sprintf("Final score: %d", score)
	finalScoreWidth := rl.MeasureText(finalScoreText, smallFontSize)

	rl.DrawText(finalScoreText,
		(ScreenWidth-int32(finalScoreWidth))/2,
		ScreenHeight*3/4-100,
		smallFontSize, rl.White)

	if !ScoreLogged {
		logScoreText := "Press L to log your score (not supported on controller)"
		logScoreWidth := rl.MeasureText(logScoreText, smallFontSize)
		rl.DrawText(logScoreText,
			(ScreenWidth-int32(logScoreWidth))/2,
			ScreenHeight*3/4+smallFontSize+10-100,
			smallFontSize, rl.White)
	}

	quitText := "Press Q to quit (B on controller)"
	quitWidth := rl.MeasureText(quitText, smallFontSize)
	rl.DrawText(quitText,
		(ScreenWidth-int32(quitWidth))/2,
		ScreenHeight*3/4+smallFontSize*2+90-100,
		smallFontSize, rl.White)

	restartText := "Press R to restart (Select/Back on controller)"
	restartWidth := rl.MeasureText(restartText, smallFontSize)
	rl.DrawText(restartText,
		(ScreenWidth-int32(restartWidth))/2,
		ScreenHeight*3/4+smallFontSize*2+50-100,
		smallFontSize, rl.White)
}

func DrawIntro() {
	instructionText := "Press H to view high scores (Y on controller)"
	instructionFontSize := int32(20)
	instructionTextWidth := rl.MeasureText(instructionText, instructionFontSize)
	instructionTextX := int32(ScreenWidth - instructionTextWidth - 10) // 10 pixels from the right edge
	instructionTextY := int32(ScreenHeight - instructionFontSize - 10) // 10 pixels from the bottom edge
	rl.DrawText(instructionText, instructionTextX, instructionTextY, instructionFontSize, rl.White)

	pauseText := "Press ESC to pause (Start on controller)"
	pauseFontSize := int32(20)
	pauseTextWidth := rl.MeasureText(pauseText, pauseFontSize)
	pauseTextX := int32(ScreenWidth - pauseTextWidth - 10) // 10 pixels from the right edge
	pauseTextY := int32(ScreenHeight - pauseFontSize - 40) // 40 pixels from the bottom edge
	rl.DrawText(pauseText, pauseTextX, pauseTextY, pauseFontSize, rl.White)
}

func DrawScore() {
	rl.DrawText(fmt.Sprintf("Timer: %ds", int(timer)), 0, 0, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Difficulty: %dx", int(ProjectileSpawnRate/InitialProjectileSpawnRate)), 0, 20, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Score: %d", CurrentScore), 0, 40, 20, rl.White)
}
