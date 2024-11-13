package game

import (
	"io/fs"
)

type Game struct {
	Fs    fs.FS
	Cards []Card
	Rules BasicRules
	Board [][]Tile
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

func NewGame(fs fs.FS) *Game {

	rules := LoadRules("static/rules/basicRules.json", fs)
	tiles := make([][]Tile, rules.BoardWidth)
	rows := make([]Tile, rules.BoardHeight*rules.BoardWidth)
	for i, startRow := 0, 0; i < rules.BoardWidth; i, startRow = i+1, startRow+rules.BoardHeight {
		endRow := startRow + rules.BoardHeight
		tiles[i] = rows[startRow:endRow:endRow]
	}
	g := Game{
		Fs:    fs,
		Cards: BuildCards(rules, fs),
		Rules: rules,
		Board: tiles,
	}
	return &g
}
