package game

import "github.com/hajimehoshi/ebiten/v2"

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, ScreenHeight int) {
	return 720, 720
}

func (g Game) Update() error {
	return nil
}
