package game

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Draw(screen *ebiten.Image) {
	for _, row := range g.Board {
		for _, tile := range row {

			tile.Draw(screen)

		}
	}
}

func (g *Game) Draw2(screen *ebiten.Image) {
	for i := 1; i < 13; i++ {
		card := g.Cards[i]
		x := 32
		y := i * 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-16.0, -16.0)

		op.GeoM.Rotate(card.Image.rotateAngle)
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(card.Image.img, op)
	}
}

func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, ScreenHeight int) {
	return 720, 720
}

func (g Game) Update() error {

	next := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if next {
		seed := rand.Uint64()
		g.NewSeed(seed)

	}

	return nil
}
