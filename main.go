package main

import (
	"log"

	"SLASI/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Get current window size
	w, h := ebiten.Monitor().Size()
	if w <= game.ScreenWidth || h <= game.ScreenHeight {
		ebiten.SetFullscreen(true)
	} else {
		// For larger screens, optionally scale up
		scale := 1.0
		if w >= game.ScreenWidth*2 && h >= game.ScreenHeight*2 {
			scale = 2.0
		}
		ebiten.SetWindowSize(
			int(float64(game.ScreenWidth)*scale),
			int(float64(game.ScreenHeight)*scale),
		)
	}

	ebiten.SetWindowTitle("SLASI")

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle(game.GameName)

	// Initialize game
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	// Run game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
