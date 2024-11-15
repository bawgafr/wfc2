package main

import (
	"embed"
	"wfc2/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed static/images/* static/rules/*
var embededStatic embed.FS

func main() {

	g := game.NewGame(embededStatic, 42)

	ebiten.SetWindowSize(720, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Wave function collapse 2")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}

}
