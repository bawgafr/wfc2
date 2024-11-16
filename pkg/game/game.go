package game

import (
	"io/fs"

	"golang.org/x/exp/rand"
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

type BasicRandom struct {
	R *rand.Rand
}

func (b BasicRandom) Intn(n int) int {
	return b.R.Intn(n)
}

func NewGame(fs fs.FS, seed uint64) *Game {

	rules := LoadRules("static/rules/basicRules.json", fs)
	cards := BuildCards(rules, fs)
	tiles := make([][]Tile, rules.BoardWidth)
	rows := make([]Tile, rules.BoardHeight*rules.BoardWidth)
	for i, startRow := 0, 0; i < rules.BoardWidth; i, startRow = i+1, startRow+rules.BoardHeight {
		endRow := startRow + rules.BoardHeight
		tiles[i] = rows[startRow:endRow:endRow]
	}
	// add in the seed tiles

	for _, seedTile := range rules.SeedTiles {
		tiles[seedTile.X][seedTile.Y] = Tile{
			X:    seedTile.X,
			Y:    seedTile.Y,
			Card: cards[seedTile.Id-1],
		}
	}

	// create the random number generator and seed it
	s := rand.NewSource(seed)
	r := rand.New(s)

	g := Game{
		Fs:    fs,
		Cards: cards,
		Rules: rules,
		Board: tiles,
		R:     BasicRandom{R: r},
	}
	return &g
}
