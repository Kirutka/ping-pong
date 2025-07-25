package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"ping-pong/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Ping-Pong Game")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}