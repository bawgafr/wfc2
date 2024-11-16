package game

import (
	"fmt"
	"io/fs"
	"time"
)

type Game struct {
	Fs    fs.FS
	Cards []Card
	Rules BasicRules
	Board [][]Tile
	R     Rnd
}

type SeedTiles struct {
	X  int
	Y  int
	Id int
}

type BasicRules struct {
	ImageSize   int
	BoardWidth  int
	BoardHeight int
	BaseCards   []BaseCards
	SeedTiles   []SeedTiles
}

type Rnd interface {
	Intn(n int) int
}

// type BasicRandom struct {
// 	R *rand.Rand
// }

// func (b BasicRandom) Intn(n int) int {
// 	return b.R.Intn(n)
// }

func NewGame(fs fs.FS, seed uint64) *Game {

	rules := LoadRules("static/rules/basicRules.json", fs)
	cards := BuildCards(rules, fs)

	tiles := NewBoard(rules, cards)
	// tiles := make([][]Tile, rules.BoardWidth)
	// rows := make([]Tile, rules.BoardHeight*rules.BoardWidth)
	// for i, startRow := 0, 0; i < rules.BoardWidth; i, startRow = i+1, startRow+rules.BoardHeight {
	// 	endRow := startRow + rules.BoardHeight
	// 	tiles[i] = rows[startRow:endRow:endRow]
	// }

	// create the random number generator and seed it

	basicRandom := NewBasicRandom(seed)

	g := Game{
		Fs:    fs,
		Cards: cards,
		Rules: rules,
		Board: tiles,
		R:     basicRandom,
	}
	return &g
}

func (g *Game) CreateLandscape() {
	buildBoard := getInitialBuildBoard(g)

	cnt := 0
	startTime := time.Now()
	for i := 0; i < g.Rules.BoardHeight*g.Rules.BoardWidth; i++ {
		cnt++
		if !g.evolveBoard(&buildBoard) {
			break
		}
	}
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("%d evolutions of board in %v\n", cnt, elapsedTime)

}
