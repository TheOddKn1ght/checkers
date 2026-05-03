package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Checkers")
	if err := ebiten.RunGame(NewRunner(NewGame())); err != nil {
		log.Fatal(err)
	}
}
