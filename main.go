package main

import (
    "log"
    
    "github.com/hajimehoshi/ebiten/v2"
    "ping-pong/game"
)

func main() {
    g := game.NewGame()
    
    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("Ping-Pong")
    
    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}